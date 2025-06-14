using ChessEngine.Domain.Models;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ChessEngine.Domain.Interfaces
{
    public interface IMoveFinder
    {
        (int evaluation, MoveDetail bestMove) GetBestMove(Position position, int timeLimitMs);
    }
}
