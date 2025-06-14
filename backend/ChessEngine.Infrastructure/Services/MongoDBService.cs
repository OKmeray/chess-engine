using ChessEngine.Persistance.Models;
using MongoDB.Driver;
using Microsoft.Extensions.Options;


namespace ChessEngine.Persistance.Services
{
    public class MongoDBService
    {
        private readonly IMongoCollection<EngineConfig> _configsCollection;

        public MongoDBService(
            IOptions<MongoDBSettings> mongoDBSettings)
        {
            var client = new MongoClient(mongoDBSettings.Value.ConnectionString);
            var database = client.GetDatabase(mongoDBSettings.Value.DatabaseName);
            _configsCollection = database.GetCollection<EngineConfig>(mongoDBSettings.Value.CollectionName);
        }

        public async Task<List<EngineConfig>> GetAsync() =>
            await _configsCollection.Find(_ => true).ToListAsync();

        public async Task<EngineConfig?> GetAsync(string hexId) =>
            await _configsCollection.Find(x => x.HEXID == hexId).FirstOrDefaultAsync();

        public async Task CreateAsync(EngineConfig config) =>
            await _configsCollection.InsertOneAsync(config);

        public async Task UpdateAsync(string hexId, EngineConfig config) =>
            await _configsCollection.ReplaceOneAsync(x => x.HEXID == hexId, config);

        public async Task RemoveAsync(string hexId) =>
            await _configsCollection.DeleteOneAsync(x => x.HEXID == hexId);
    }

    public class MongoDBSettings
    {
        public string ConnectionString { get; set; }
        public string DatabaseName { get; set; }
        public string CollectionName { get; set; }
    }
}
