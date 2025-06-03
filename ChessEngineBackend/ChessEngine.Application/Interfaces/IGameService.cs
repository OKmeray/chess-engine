using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ChessEngine.Application.Interfaces
{
    public interface IGameService
    {
        public object GetMove(string fen, int from, int to, List<string> selectedVariations, int time);
        public Dictionary<int, List<int>> GetPossibleMovesByFen(string fen);
    }
}
