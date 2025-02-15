using Microsoft.AspNetCore.Mvc;
using ChessEngine.Application.Interfaces;

namespace ChessEngine.API.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class GameController : ControllerBase
    {
        private readonly ILogger<GameController> _logger;
        private readonly IGameService _gameService;

        public GameController(ILogger<GameController> logger, IGameService gameService)
        {
            _logger = logger;
            _gameService = gameService;
        }

        [HttpGet(Name = "move")]
        public string Get()
        {
            // string res = _gameService.ToString();
            return "Nc6";
            // return 0;
        }
    }
}
