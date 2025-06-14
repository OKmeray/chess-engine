using System;
using System.Collections.Generic;
using ChessEngine.Enums;
using ChessEngine.Domain.Models;
using ChessEngine.Helpers;
using System.Linq;

namespace ChessEngine.Evaluation
{
    public static class Evaluator
    {
        private static readonly int[] WhitePawnSquareTable = new int[]
        {
            0,  0,  0,  0,  0,  0,  0,  0,
            50, 50, 50, 50, 50, 50, 50, 50,
            10, 10, 20, 30, 30, 20, 10, 10,
             5,  5, 10, 25, 25, 10,  5,  5,
             0,  0,  0, 20, 20,  0,  0,  0,
             5, -5,-10,  0,  0,-10, -5,  5,
             5, 10, 10,-20,-20, 10, 10,  5,
             0,  0,  0,  0,  0,  0,  0,  0
        };

        private static readonly int[] BlackPawnSquareTable = WhitePawnSquareTable
    .Reverse()
    .Select(item => -item)
    .ToArray();

        private static readonly int[] WhiteKnightSquareTable = new int[]
        {
            -50,-40,-30,-30,-30,-30,-40,-50,
            -40,-20,  0,  0,  0,  0,-20,-40,
            -30,  0, 10, 15, 15, 10,  0,-30,
            -30,  5, 15, 20, 20, 15,  5,-30,
            -30,  0, 15, 20, 20, 15,  0,-30,
            -30,  5, 10, 15, 15, 10,  5,-30,
            -40,-20,  0,  5,  5,  0,-20,-40,
            -50,-40,-30,-30,-30,-30,-40,-50,
        };
        private static readonly int[] BlackKnightSquareTable = WhiteKnightSquareTable
    .Reverse()
    .Select(item => -item)
    .ToArray();

        private static readonly int[] WhiteBishopSquareTable = new int[]
        {
            -20,-10,-10,-10,-10,-10,-10,-20,
            -10,  0,  0,  0,  0,  0,  0,-10,
            -10,  0,  5, 10, 10,  5,  0,-10,
            -10,  5,  5, 10, 10,  5,  5,-10,
            -10,  0, 10, 10, 10, 10,  0,-10,
            -10, 10, 10, 10, 10, 10, 10,-10,
            -10,  5,  0,  0,  0,  0,  5,-10,
            -20,-10,-10,-10,-10,-10,-10,-20,
        };
        private static readonly int[] BlackBishopSquareTable = WhiteBishopSquareTable
    .Reverse()
    .Select(item => -item)
    .ToArray();

        private static readonly int[] WhiteRookSquareTable = new int[]
        {
              0,  0,  0,  0,  0,  0,  0,  0,
              5, 10, 10, 10, 10, 10, 10,  5,
             -5,  0,  0,  0,  0,  0,  0, -5,
             -5,  0,  0,  0,  0,  0,  0, -5,
             -5,  0,  0,  0,  0,  0,  0, -5,
             -5,  0,  0,  0,  0,  0,  0, -5,
             -5,  0,  0,  0,  0,  0,  0, -5,
              0,  0,  0,  5,  5,  0,  0,  0
        };
        private static readonly int[] BlackRookSquareTable = WhiteRookSquareTable
    .Reverse()
    .Select(item => -item)
    .ToArray();

        private static readonly int[] WhiteQueenSquareTable = new int[]
        {
            -20,-10,-10, -5, -5,-10,-10,-20,
            -10,  0,  0,  0,  0,  0,  0,-10,
            -10,  0,  5,  5,  5,  5,  0,-10,
             -5,  0,  5,  5,  5,  5,  0, -5,
              0,  0,  5,  5,  5,  5,  0, -5,
            -10,  5,  5,  5,  5,  5,  0,-10,
            -10,  0,  5,  0,  0,  0,  0,-10,
            -20,-10,-10, -5, -5,-10,-10,-20
        };
        private static readonly int[] BlackQueenSquareTable = WhiteQueenSquareTable
    .Reverse()
    .Select(item => -item)
    .ToArray();

        private static readonly int[] WhiteKingSquareMiddlegameTable = new int[]
        {
            -30,-40,-40,-50,-50,-40,-40,-30,
            -30,-40,-40,-50,-50,-40,-40,-30,
            -30,-40,-40,-50,-50,-40,-40,-30,
            -30,-40,-40,-50,-50,-40,-40,-30,
            -20,-30,-30,-40,-40,-30,-30,-20,
            -10,-20,-20,-20,-20,-20,-20,-10,
             20, 20,  0,  0,  0,  0, 20, 20,
             20, 30, 10,  0,  0, 10, 30, 20
        };
        private static readonly int[] BlackKingSquareMiddlegameTable = WhiteKingSquareMiddlegameTable
    .Reverse()
    .Select(item => -item)
    .ToArray();

        private static readonly int[] WhiteKingSquareEndgameTable = new int[]
        {
            -50,-40,-30,-20,-20,-30,-40,-50,
            -30,-20,-10,  0,  0,-10,-20,-30,
            -30,-10, 20, 30, 30, 20,-10,-30,
            -30,-10, 30, 40, 40, 30,-10,-30,
            -30,-10, 30, 40, 40, 30,-10,-30,
            -30,-10, 20, 30, 30, 20,-10,-30,
            -30,-30,  0,  0,  0,  0,-30,-30,
            -50,-30,-30,-30,-30,-30,-30,-50
        };
        private static readonly int[] BlackKingSquareEndgameTable = WhiteKingSquareEndgameTable
    .Reverse()
    .Select(item => -item)
    .ToArray();

        /// <summary>
        /// Evaluates the current position from White's perspective.
        /// Positive values favor White, negative favor Black.
        /// </summary>
        public static int EvaluatePosition(Position position)
        {
            int evaluation = 0;

            // Iterate through all squares
            for (int square = 0; square < 64; square++)
            {
                var (piece, color) = position.GetPieceAndColorBySquare(square);
                if (piece == PieceEnum.NONE)
                    continue;

                int sideSign = color == PieceColor.WHITE ? 1 : -1;

                switch (piece)
                {
                    case PieceEnum.PAWN:
                        evaluation += (int)PiecePrice.PAWN * sideSign;
                        evaluation += color == PieceColor.WHITE ? WhitePawnSquareTable[square] : BlackPawnSquareTable[square];
                        break;
                    case PieceEnum.KNIGHT:
                        evaluation += (int)PiecePrice.KNIGHT * sideSign;
                        evaluation += color == PieceColor.WHITE ? WhiteKnightSquareTable[square] : BlackKnightSquareTable[square];
                        break;
                    case PieceEnum.BISHOP:
                        evaluation += (int)PiecePrice.BISHOP * sideSign;
                        evaluation += color == PieceColor.WHITE ? WhiteBishopSquareTable[square] : BlackBishopSquareTable[square];
                        break;
                    case PieceEnum.ROOK:
                        evaluation += (int)PiecePrice.ROOK * sideSign;
                        evaluation += color == PieceColor.WHITE ? WhiteRookSquareTable[square] : BlackRookSquareTable[square];
                        break;
                    case PieceEnum.QUEEN:
                        evaluation += (int)PiecePrice.QUEEN * sideSign;
                        evaluation += color == PieceColor.WHITE ? WhiteQueenSquareTable[square] : BlackQueenSquareTable[square];
                        break;
                    case PieceEnum.KING:
                        evaluation += (int)PiecePrice.KING * sideSign;

                        // Determine game phase
                        bool isEndgame = IsEndgame(position);
                        evaluation += isEndgame
                                      ? (color == PieceColor.WHITE ? WhiteKingSquareEndgameTable[square] : BlackKingSquareEndgameTable[square])
                                      : (color == PieceColor.WHITE ? WhiteKingSquareMiddlegameTable[square] : BlackKingSquareMiddlegameTable[square]);
                        break;
                }
            }

            return evaluation;
        }

        /// <summary>
        /// Determines if the game is in the endgame phase.
        /// </summary>
        private static bool IsEndgame(Position position)
        {
            int queensCount = 0;
            int piecesCount = 0;

            foreach (var bitboard in position.Bitboards)
            {
                if (bitboard.BitboardValue != 0)
                {
                    // Assuming PieceEnum.QUEEN is at index 4
                    // Adjust based on your actual enum definitions
                    if ((int)PieceEnum.QUEEN == position.Bitboards.IndexOf(bitboard) % 6)
                        queensCount++;

                    // Count pieces excluding kings and queens
                    if ((int)PieceEnum.QUEEN != position.Bitboards.IndexOf(bitboard) % 6 &&
                        (int)PieceEnum.KING != position.Bitboards.IndexOf(bitboard) % 6)
                        piecesCount++;
                }
            }

            return queensCount <= 1 && piecesCount <= 4;
        }
    }
}
