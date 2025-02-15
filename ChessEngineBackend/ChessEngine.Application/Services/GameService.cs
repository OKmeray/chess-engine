using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using ChessEngine.Application.Interfaces;
using ChessEngine.Application.Fen;
using ChessEngine.Domain.Models;
using ChessEngine.Enums;

namespace ChessEngine.Application.Services
{
    public class GameService : IGameService
    {
        public object GetMove(string fen, int from, int to, List<string> selectedVariations)
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
                // TODO: currently alpha and beta for black side only! change it
                (int eval, MoveDetail bestMove) = Minimax.Minimax.RunMinimax(position: position, alpha: Int32.MinValue, beta: Int32.MaxValue, depth: 3, isRoot: true);
                position.ApplyMove(bestMove);

                Dictionary<int, List<int>> possibleMovesDict = GetPossibleMoves(position);


                string newFen = FenGenerator.GetFenFromPosition(position);

                object response = new
                {
                    possibleMoves = possibleMovesDict,
                    fen = newFen,
                    isLegal = true,
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
                    isLegal = false
                };
            }
        }

        public Dictionary<int, List<int>> GetPossibleMoves(Position position)
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
    }
}
