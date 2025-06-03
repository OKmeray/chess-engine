using ChessEngine.Helpers;
using System.Collections.Generic;
using ChessEngine.Enums;
using ChessEngine.Domain.Models;

namespace ChessEngine.MoveGeneration
{
    public static class GenerateMove
    {
        public static readonly ulong Rank1 = 0xFF00000000000000UL;
        public static readonly ulong Rank2 = 0x00FF000000000000UL;
        public static readonly ulong Rank7 = 0x000000000000FF00UL;
        public static readonly ulong Rank8 = 0x00000000000000FFUL;

        public static readonly ulong FileA = 0x0101010101010101UL;
        public static readonly ulong FileB = 0x0202020202020202UL;
        public static readonly ulong FileG = 0x4040404040404040UL;
        public static readonly ulong FileH = 0x8080808080808080UL;
        public static readonly ulong FileE = 0x1010101010101010UL;

        public static bool IsAFile(ulong square) => (square & FileA) != 0;
        public static bool IsBFile(ulong square) => (square & FileB) != 0;
        public static bool IsEFile(ulong square) => (square & FileE) != 0;
        public static bool IsGFile(ulong square) => (square & FileG) != 0;
        public static bool IsHFile(ulong square) => (square & FileH) != 0;

        public static bool IsRank1(ulong square) => (square & Rank1) != 0;
        public static bool IsRank2(ulong square) => (square & Rank2) != 0;
        public static bool IsRank7(ulong square) => (square & Rank7) != 0;
        public static bool IsRank8(ulong square) => (square & Rank8) != 0;

        /// <summary>
        /// Generates all possible pawn moves for a given square.
        /// </summary>
        public static List<ulong> GeneratePawnMove(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, PieceColor color, ulong? enPassantSquare)
        {
            List<ulong> moves = new List<ulong>();

            if (color == PieceColor.WHITE)
            {
                ulong target = square >> 9;
                if (!IsAFile(square) && ((target & enemyBitboard.BitboardValue) != 0 || (enPassantSquare.HasValue && target == enPassantSquare.Value)))
                {
                    moves.Add(target);
                }

                target = square >> 7;
                if (!IsHFile(square) && ((target & enemyBitboard.BitboardValue) != 0 || (enPassantSquare.HasValue && target == enPassantSquare.Value)))
                {
                    moves.Add(target);
                }

                target = square >> 8;
                if ((target & (enemyBitboard.BitboardValue | ownBitboard.BitboardValue)) == 0)
                {
                    moves.Add(target);

                    if (IsRank2(square))
                    {
                        ulong doubleForward = square >> 16;
                        if ((doubleForward & (enemyBitboard.BitboardValue | ownBitboard.BitboardValue)) == 0)
                        {
                            moves.Add(doubleForward);
                        }
                    }
                }
            }
            else // BLACK
            {
                ulong target = square << 9;
                if (!IsHFile(square) && ((target & enemyBitboard.BitboardValue) != 0 || (enPassantSquare.HasValue && target == enPassantSquare.Value)))
                {
                    moves.Add(target);
                }

                target = square << 7;
                if (!IsAFile(square) && ((target & enemyBitboard.BitboardValue) != 0 || (enPassantSquare.HasValue && target == enPassantSquare.Value)))
                {
                    moves.Add(target);
                }

                target = square << 8;
                if ((target & (enemyBitboard.BitboardValue | ownBitboard.BitboardValue)) == 0)
                {
                    moves.Add(target);

                    if (IsRank7(square))
                    {
                        ulong doubleForward = square << 16;
                        if ((doubleForward & (enemyBitboard.BitboardValue | ownBitboard.BitboardValue)) == 0)
                        {
                            moves.Add(doubleForward);
                        }
                    }
                }
            }

            return moves;
        }

        /// <summary>
        /// Generates all possible knight moves for a given square.
        /// </summary>
        public static List<ulong> GenerateKnightMove(ulong square, Bitboard ownBitboard)
        {
            List<ulong> moves = new List<ulong>();

            var knightMoves = new (ulong Shift, ulong EdgeMask)[]
            {
                (square >> 17, FileA), // Up-Left
                (square >> 15, FileH), // Up-Right
                (square >> 10, FileA | FileB), // Left-Up
                (square >> 6, FileH | FileG),  // Right-Up
                (square << 6, FileA | FileB),  // Left-Down
                (square << 10, FileH | FileG), // Right-Down
                (square << 15, FileA), // Down-Left
                (square << 17, FileH)  // Down-Right
            };

            foreach (var (shiftedMove, edgeMask) in knightMoves)
            {
                if ((square & edgeMask) == 0 && (shiftedMove & ownBitboard.BitboardValue) == 0)
                {
                    moves.Add(shiftedMove);
                }
            }

            return moves;
        }

        /// <summary>
        /// Generates all possible bishop moves for a given square.
        /// </summary>
        public static List<ulong> GenerateBishopMove(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8)
        {
            List<ulong> moves = new List<ulong>();

            moves.AddRange(GenerateDiagonalUp(square, ownBitboard, enemyBitboard, steps));
            moves.AddRange(GenerateDiagonalDown(square, ownBitboard, enemyBitboard, steps));
            moves.AddRange(GenerateAntidiagonalUp(square, ownBitboard, enemyBitboard, steps));
            moves.AddRange(GenerateAntidiagonalDown(square, ownBitboard, enemyBitboard, steps));

            return moves;
        }

        /// <summary>
        /// Generates all possible rook moves for a given square.
        /// </summary>
        public static List<ulong> GenerateRookMove(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8)
        {
            List<ulong> moves = new List<ulong>();

            moves.AddRange(GenerateVerticalUp(square, ownBitboard, enemyBitboard, steps));
            moves.AddRange(GenerateVerticalDown(square, ownBitboard, enemyBitboard, steps));
            moves.AddRange(GenerateHorizontalLeft(square, ownBitboard, enemyBitboard, steps));
            moves.AddRange(GenerateHorizontalRight(square, ownBitboard, enemyBitboard, steps));

            return moves;
        }

        /// <summary>
        /// Generates all possible queen moves for a given square.
        /// </summary>
        public static List<ulong> GenerateQueenMove(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8)
        {
            List<ulong> moves = new List<ulong>();

            moves.AddRange(GenerateRookMove(square, ownBitboard, enemyBitboard, steps));
            moves.AddRange(GenerateBishopMove(square, ownBitboard, enemyBitboard, steps));

            return moves;
        }

        /// <summary>
        /// Generates all possible king moves for a given square.
        /// </summary>
        public static List<ulong> GenerateKingMove(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, bool isWhite = true, int steps = 1)
        {
            List<ulong> moves = new List<ulong>();

            //ulong[] kingOffsets = {
            //    square >> 9, square >> 8, square >> 7,
            //    square >> 1, square << 1,
            //    square << 7, square << 8, square << 9
            //};
            moves.AddRange(GenerateRookMove(square, ownBitboard, enemyBitboard, steps));
            moves.AddRange(GenerateBishopMove(square, ownBitboard, enemyBitboard, steps));

            moves.AddRange(GenerateCastleMoves(square, ownBitboard, enemyBitboard, isWhite));

            return moves;
        }

        /// <summary>
        /// Generates all possible castling moves for a given king.
        /// </summary>
        public static List<ulong> GenerateCastleMoves(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, bool isWhite = true)
        {
            List<ulong> moves = new List<ulong>();

            if ((isWhite && square != 0x1000000000000000UL) || (!isWhite && square != 0x0000000000000010UL))
                return moves;

            bool canCastleKingside = true; // Replace with actual castling rights check
            if (canCastleKingside)
            {
                ulong emptySquares = isWhite ? ((1UL << 61) | (1UL << 62)) : ((1UL << 5) | (1UL << 6));
                if (((ownBitboard.BitboardValue | enemyBitboard.BitboardValue) & emptySquares) == 0)
                {
                    moves.Add(isWhite ? 1UL << 62 : 1UL << 6);
                }
            }

            bool canCastleQueenside = true; // Replace with actual castling rights check
            if (canCastleQueenside)
            {
                ulong emptySquares = isWhite ? ((1UL << 59) | (1UL << 58) | (1UL << 57)) : ((1UL << 3) | (1UL << 2) | (1UL << 1));
                if (((ownBitboard.BitboardValue | enemyBitboard.BitboardValue) & emptySquares) == 0)
                {
                    moves.Add(isWhite ? 1UL << 58 : 1UL << 2);
                }
            }

            return moves;
        }

        private static List<ulong> GenerateDiagonalUp(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8) // -9
        {
            List<ulong> moves = new List<ulong>();
            while (steps > 0)
            {
                if (!(IsAFile(square) || IsRank8(square)))
                {
                    square >>= 9;

                    // Check for own or enemy pieces on target square
                    if ((square & ownBitboard.BitboardValue) == square)
                    {
                        return moves;
                    }
                    else if ((square & enemyBitboard.BitboardValue) == square)
                    {
                        moves.Add(square);
                        return moves;
                    }
                    else
                    {
                        moves.Add(square);
                        steps--;
                    }
                }
                else
                {
                    return moves;
                }
            }
            return moves;
        }

        private static List<ulong> GenerateDiagonalDown(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8) // +9
        {
            List<ulong> moves = new List<ulong>();
            while (steps > 0)
            {
                if (!(IsHFile(square) || IsRank1(square)))
                {
                    square <<= 9;

                    // Check for own or enemy pieces on target square
                    if ((square & ownBitboard.BitboardValue) == square)
                    {
                        return moves;
                    }
                    else if ((square & enemyBitboard.BitboardValue) == square)
                    {
                        moves.Add(square);
                        return moves;
                    }
                    else
                    {
                        moves.Add(square);
                        steps--;
                    }
                }
                else
                {
                    return moves;
                }
            }
            return moves;
        }

        private static List<ulong> GenerateAntidiagonalUp(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8) // -7
        {
            List<ulong> moves = new List<ulong>();
            while (steps > 0)
            {
                if (!(IsHFile(square) || IsRank8(square)))
                {
                    square >>= 7;

                    // Check for own or enemy pieces on target square
                    if ((square & ownBitboard.BitboardValue) == square)
                    {
                        return moves;
                    }
                    else if ((square & enemyBitboard.BitboardValue) == square)
                    {
                        moves.Add(square);
                        return moves;
                    }
                    else
                    {
                        moves.Add(square);
                        steps--;
                    }
                }
                else
                {
                    return moves;
                }
            }
            return moves;
        }

        private static List<ulong> GenerateAntidiagonalDown(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8) // +7
        {
            List<ulong> moves = new List<ulong>();
            while (steps > 0)
            {
                if (!(IsAFile(square) || IsRank1(square)))
                {
                    square <<= 7;

                    // Check for own or enemy pieces on target square
                    if ((square & ownBitboard.BitboardValue) == square)
                    {
                        return moves;
                    }
                    else if ((square & enemyBitboard.BitboardValue) == square)
                    {
                        moves.Add(square);
                        return moves;
                    }
                    else
                    {
                        moves.Add(square);
                        steps--;
                    }
                }
                else
                {
                    return moves;
                }
            }
            return moves;
        }

        private static List<ulong> GenerateVerticalUp(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8) // -8
        {
            List<ulong> moves = new List<ulong>();
            while (steps > 0)
            {
                if (!IsRank8(square))
                {
                    square >>= 8;

                    if ((square & ownBitboard.BitboardValue) == square) return moves;
                    if ((square & enemyBitboard.BitboardValue) == square)
                    {
                        moves.Add(square);
                        return moves;
                    }

                    moves.Add(square);
                    steps--;
                }
                else
                {
                    return moves;
                }
            }
            return moves;
        }

        private static List<ulong> GenerateVerticalDown(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8) // +8
        {
            List<ulong> moves = new List<ulong>();
            while (steps > 0)
            {
                if (!IsRank1(square))
                {
                    square <<= 8;

                    if ((square & ownBitboard.BitboardValue) == square) return moves;
                    if ((square & enemyBitboard.BitboardValue) == square)
                    {
                        moves.Add(square);
                        return moves;
                    }

                    moves.Add(square);
                    steps--;
                }
                else
                {
                    return moves;
                }
            }
            return moves;
        }

        private static List<ulong> GenerateHorizontalLeft(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8) // -1
        {
            List<ulong> moves = new List<ulong>();
            while (steps > 0)
            {
                if (!IsAFile(square))
                {
                    square >>= 1;

                    if ((square & ownBitboard.BitboardValue) == square) return moves;
                    if ((square & enemyBitboard.BitboardValue) == square)
                    {
                        moves.Add(square);
                        return moves;
                    }

                    moves.Add(square);
                    steps--;
                }
                else
                {
                    return moves;
                }
            }
            return moves;
        }

        private static List<ulong> GenerateHorizontalRight(ulong square, Bitboard ownBitboard, Bitboard enemyBitboard, int steps = 8) // +1
        {
            List<ulong> moves = new List<ulong>();
            while (steps > 0)
            {
                if (!IsHFile(square))
                {
                    square <<= 1;

                    if ((square & ownBitboard.BitboardValue) == square) return moves;
                    if ((square & enemyBitboard.BitboardValue) == square)
                    {
                        moves.Add(square);
                        return moves;
                    }

                    moves.Add(square);
                    steps--;
                }
                else
                {
                    return moves;
                }
            }
            return moves;
        }
    }
}
