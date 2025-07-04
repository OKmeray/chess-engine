﻿
namespace ChessEngine.Application.Interfaces
{
    public interface IGameService
    {
        public object GetMove(string fen, int from, int to, float time, string promotion, bool requestFirstMove = false, string modelVersion = "", string searchMethod = "mcts");
        public Dictionary<int, List<int>> GetPossibleMovesByFen(string fen);
    }
}
