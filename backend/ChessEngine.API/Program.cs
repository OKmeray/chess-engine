using ChessEngine.API.WebSockets;
using ChessEngine.Application.Interfaces;
using ChessEngine.Application.Services;
using ChessEngine.Persistance.Services;
using System.Text.Json;

var builder = WebApplication.CreateBuilder(args);
//builder.WebHost.UseUrls("http://localhost:5000");
builder.WebHost.UseUrls("http://*:80");

builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowAll", policyBuilder =>
    {
        policyBuilder
            .AllowAnyOrigin()
            .AllowAnyMethod()
            .AllowAnyHeader();
    });
});


// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddScoped<IGameService, GameService>();
builder.Services.AddScoped<MoveWebSocketHandler>();
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

builder.Services.Configure<MongoDBSettings>(
    builder.Configuration.GetSection("MongoDB"));
builder.Services.AddSingleton<MongoDBService>();
builder.Services.AddControllers()
    .AddJsonOptions(options =>
    {
        options.JsonSerializerOptions.PropertyNamingPolicy = JsonNamingPolicy.CamelCase;
        options.JsonSerializerOptions.DictionaryKeyPolicy = JsonNamingPolicy.CamelCase;
    });

//builder.Services.AddSingleton<IMoveFinderProvider, MoveFinderProvider>();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

var webSocketOptions = new WebSocketOptions
{
    KeepAliveInterval = TimeSpan.FromMinutes(2)
};

app.UseWebSockets(webSocketOptions);
app.UseCors("AllowAll");

app.Use(async (context, next) =>
{
    if (context.Request.Path == "/ws/move")
    {
        if (context.WebSockets.IsWebSocketRequest)
        {
            var webSocketHandler = context.RequestServices.GetRequiredService<MoveWebSocketHandler>();
            using var webSocket = await context.WebSockets.AcceptWebSocketAsync();
            await webSocketHandler.HandleAsync(webSocket);
        }
        else
        {
            context.Response.StatusCode = StatusCodes.Status400BadRequest;
        }
    }
    else
    {
        await next(context);
    }

});


// app.UseHttpsRedirection();

app.UseAuthorization();

app.MapControllers();

app.Run();
