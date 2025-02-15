using ChessEngine.API.Controllers;
using ChessEngine.Application.Interfaces;
using ChessEngine.Application.Services;
using System.Net.WebSockets;
using System.Text;
using System.Text.Json;

namespace ChessEngine.API.WebSockets
{
    public class MoveWebSocketHandler
    {
        //private readonly ILogger<GameController> _logger;
        private readonly IGameService _gameService;

        public MoveWebSocketHandler(IGameService gameService) // ILogger<MoveWebSocketHandler> logger,
        {
            //_logger = logger;
            _gameService = gameService;
        }

        public async Task HandleAsync(WebSocket webSocket)
        {

            var buffer = new byte[1024 * 4];

            try
            {
                while (webSocket.State == WebSocketState.Open)
                {
                    var result = await webSocket.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);

                    if (result.MessageType == WebSocketMessageType.Text)
                    {
                        var receivedMessage = Encoding.UTF8.GetString(buffer, 0, result.Count);
                        Console.WriteLine($"Received: {receivedMessage}");

                        var moveData = JsonSerializer.Deserialize<MoveRequest>(receivedMessage);

                        int fromSquare = int.Parse(moveData.from);
                        int toSquare = int.Parse(moveData.to);

                        var responseObject = _gameService.GetMove(moveData.fen, fromSquare, toSquare, moveData.selectedVariations);
                        string jsonResponse = JsonSerializer.Serialize(responseObject);
                        Console.WriteLine($"Sent: {jsonResponse}");
                        
                        byte[] responseBytes = Encoding.UTF8.GetBytes(jsonResponse);
                        await webSocket.SendAsync(
                            new ArraySegment<byte>(responseBytes),
                            WebSocketMessageType.Text,
                            true,
                            CancellationToken.None
                        );
                    }
                    else if (result.MessageType == WebSocketMessageType.Close)
                    {
                        await webSocket.CloseAsync(WebSocketCloseStatus.NormalClosure, "Closing", CancellationToken.None);
                    }
                }
            }
            catch (WebSocketException ex)
            {
                Console.WriteLine($"WebSocket error: {ex.Message}");
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Unexpected error: {ex.Message}");
            }
        }
    }

    public class MoveRequest
    {
        public string fen { get; set; } = string.Empty;
        public string from { get; set; } = string.Empty; 
        public string to { get; set; } = string.Empty;
        public List<string> selectedVariations { get; set; } = new();
    }
}
