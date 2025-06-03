using ChessEngine.Application.Interfaces;
using ChessEngine.Application.Fen;
using ChessEngine.Domain.Models;
using ChessEngine.Enums;
// using ChessEngine.Application.Minimax;
using ChessEngine.Application.MCTS;
using ChessEngine.Application.Minimax;

namespace ChessEngine.Application.Services
{
    public class GameService : IGameService
    {
        public object GetMove(string fen, int from, int to, List<string> selectedVariations, int time)
        {
            Position position = FenGenerator.GetPositionFromFen(fen: fen);
            (PieceEnum piece, PieceColor color) = position.GetPieceAndColorBySquare(square: from);
            MoveDetail userMove = new MoveDetail { Color = color, Piece = piece, Square = from, Move = to };

            bool isLegal = false;
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
            if (isLegal)
            {
                Console.WriteLine("Move is legal");

                position.ApplyMove(userMove);
                
                if (position.IsCheckmate())
                {
                    return new
                    {
                        outcome = PositionOutcome.LOSS
                    };
                }

                //MinimaxMoveFinder mmf = new MinimaxMoveFinder();
                //(int eval, MoveDetail bestMove) = mmf.GetBestMove(position: position, depth: 3);
                Console.WriteLine(time);
                MCTSMoveFinder mmf = new MCTSMoveFinder(
                    @"D:\LNU\8 семестр\Дипломна\onnx\puzzles 3e6.onnx",
                    timeLimitSeconds: time / 20
                );
                // (int eval, MoveDetail bestMove) = mmf.GetBestMove(rootPos: position, _: 0);
                (int eval, MoveDetail bestMove) = mmf.GetBestMove(position: position, unused: 0);

                position.ApplyMove(bestMove);
                object response = null;
                if (position.IsCheckmate() ) {
                    response = new
                    {
                        outcome = PositionOutcome.WIN
                    };
                } 
                else if (position.IsStalemate())
                {
                    response = new
                    {
                        outcome = PositionOutcome.DRAW
                    };
                }
                else
                {
                    Dictionary<int, List<int>> possibleMovesDict = GetPossibleMoves(position);

                    string newFen = FenGenerator.GetFenFromPosition(position);

                    response = new
                    {
                        possibleMoves = possibleMovesDict,
                        fen = newFen,
                        isLegal = true,
                        outcome = PositionOutcome.ONGOING
                    };
                }

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
