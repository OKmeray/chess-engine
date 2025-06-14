using ChessEngine.Application.Interfaces;
using ChessEngine.Application.Fen;
using ChessEngine.Domain.Models;
using ChessEngine.Enums;
using ChessEngine.Application.MCTS;
using ChessEngine.Application.Minimax;
using ChessEngine.Domain.Interfaces;


namespace ChessEngine.Application.Services
{
    public class GameService : IGameService
    {
        public object GetMove(string fen, int from, int to, float time, string promotion, bool requestFirstMove = false, string modelVersion = "", string searchMethod = "mcts")
        {
            Position position = FenGenerator.GetPositionFromFen(fen: fen);

            (PieceEnum piece, PieceColor color) = position.GetPieceAndColorBySquare(square: from);
            PieceEnum promotionPiece = GetPromotionPiece(promotion);
            MoveDetail userMove = new MoveDetail { Color = color, Piece = piece, Square = from, Move = to, Promotion = promotionPiece };

            bool isLegal = false;
            if (requestFirstMove)
            {
                isLegal = true;
            }
            else
            {
                var allLegalMoves = position.GetAllMoves();
                foreach (var legalMove in allLegalMoves)
                {
                    if (
                        legalMove.Piece == piece &&
                        legalMove.Color == color &&
                        legalMove.Square == from &&
                        legalMove.Move == to
                        )
                    {
                        isLegal = true;
                    }
                }
            }

            if (isLegal)
            {
                Console.WriteLine("Move is legal");
                if (!requestFirstMove)
                {
                    position.ApplyMove(userMove);
                }

                if (position.IsCheckmate())
                {
                    return new
                    {
                        outcome = PositionOutcome.LOSS
                    };
                }
                   
                Console.WriteLine($"Given time: {time} s.");
                Console.WriteLine($"Fen is: {fen}");
                IMoveFinder mmf;
                if (searchMethod == "minimax")
                {
                    Console.WriteLine("Minimax works.");
                    mmf = new MinimaxMoveFinder();
                }
                else
                {
                    Console.WriteLine("MCTS works.");
                    mmf = new MCTSMoveFinder(modelVersion);
                }

                (int eval, MoveDetail bestMove) = mmf.GetBestMove(position: position, timeLimitMs: (int)(time * 1000) / 20);

                position.ApplyMove(bestMove);
                object response = null;
                Dictionary<int, List<int>> possibleMovesDict = new Dictionary<int, List<int>>();
                PositionOutcome outcome = PositionOutcome.ONGOING;

                if (position.IsCheckmate() ) {
                    outcome = PositionOutcome.WIN;
                } 
                else if (position.IsStalemate())
                {
                    outcome = PositionOutcome.DRAW;
                }
                else
                {
                    outcome = PositionOutcome.ONGOING;
                    possibleMovesDict = GetPossibleMoves(position);
                }
                string newFen = FenGenerator.GetFenFromPosition(position);
                response = new
                {
                    possibleMoves = possibleMovesDict,
                    fen = newFen,
                    isLegal = true,
                    outcome = outcome,
                };
                return response;
            }
            else
            {
                Console.WriteLine("Illgeal");
                Dictionary<int, List<int>> possibleMovesDict = GetPossibleMoves(position);
                return new
                {
                    possibleMoves = possibleMovesDict,
                    fen = fen,
                    isLegal = false,
                    outcome = PositionOutcome.ONGOING
                };
            }
        }

        private PieceEnum GetPromotionPiece(string promotion)
        {
            if (promotion == null)
            {
                return PieceEnum.NONE;
            }
            else if (promotion.ToLower() == "q")
            {
                return PieceEnum.QUEEN;
            }
            else if (promotion.ToLower() == "r")
            {
                return PieceEnum.ROOK;
            }
            else if (promotion.ToLower() == "b")
            {
                return PieceEnum.BISHOP;
            }
            else if (promotion.ToLower() == "n")
            {
                return PieceEnum.KNIGHT;
            }
            else
            {
                return PieceEnum.NONE;
            }
        }

        private Dictionary<int, List<int>> GetPossibleMoves(Position position)
        {
            List<MoveDetail> possibleMoves = position.GetAllMoves();
            Dictionary<int, List<int>> possibleMovesDict = new Dictionary<int, List<int>>();

            foreach (var move in possibleMoves)
            {
                if (possibleMovesDict.ContainsKey(move.Square))
                {
                    possibleMovesDict[move.Square].Add(move.Move);
                }
                else
                {
                    possibleMovesDict[move.Square] = new List<int> { move.Move };
                }
            }

            return possibleMovesDict;
        }

        public Dictionary<int, List<int>> GetPossibleMovesByFen(string fen)
        {
            Position position = FenGenerator.GetPositionFromFen(fen: fen);
            return GetPossibleMoves(position: position);
        }
    }
}
