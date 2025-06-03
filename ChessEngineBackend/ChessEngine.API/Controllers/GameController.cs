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

        [HttpGet(Name = "possibleMoves")]
        public object Get(string fen)
        {
            Dictionary<int, List<int>> possibleMoves = _gameService.GetPossibleMovesByFen(fen: fen);
            return new
            {
                possibleMoves = possibleMoves,
            };
        }
    }
}
