using ChessEngine.Persistance.Models;
using ChessEngine.Persistance.Services;
using Microsoft.AspNetCore.Mvc;

namespace Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class EngineConfigsController : ControllerBase
    {
        private readonly MongoDBService _mongoDBService;

        public EngineConfigsController(MongoDBService mongoDBService)
        {
            _mongoDBService = mongoDBService;
        }

        [HttpGet]
        public async Task<List<EngineConfig>> Get() =>
            await _mongoDBService.GetAsync();

        [HttpGet("{hexId}")]
        public async Task<ActionResult<EngineConfig>> Get(string hexId)
        {
            var config = await _mongoDBService.GetAsync(hexId);
            if (config is null)
                return NotFound();
            return config;
        }

        [HttpPost]
        public async Task<IActionResult> Post(EngineConfig config)
        {
            await _mongoDBService.CreateAsync(config);
            return CreatedAtAction(nameof(Get), new { hexId = config.HEXID }, config);
        }

        [HttpPut("{hexId}")]
        public async Task<IActionResult> Update(string hexId, EngineConfig config)
        {
            var existingConfig = await _mongoDBService.GetAsync(hexId);
            if (existingConfig is null)
                return NotFound();

            config.Id = existingConfig.Id;
            await _mongoDBService.UpdateAsync(hexId, config);
            return NoContent();
        }

        [HttpDelete("{hexId}")]
        public async Task<IActionResult> Delete(string hexId)
        {
            var config = await _mongoDBService.GetAsync(hexId);
            if (config is null)
                return NotFound();

            await _mongoDBService.RemoveAsync(hexId);
            return NoContent();
        }
    }
}
