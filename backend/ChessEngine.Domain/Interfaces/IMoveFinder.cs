using ChessEngine.Domain.Models;

namespace ChessEngine.Domain.Interfaces
{
    public interface IMoveFinder
    {
        (int evaluation, MoveDetail bestMove) GetBestMove(Position position, int timeLimitMs);
    }
}
