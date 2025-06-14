using ChessEngine.Enums;
using ChessEngine.Helpers;
using ChessEngine.MoveGeneration;
using ChessEngine.Domain.Models;
using System.Numerics;

namespace ChessEngine.Domain.Models
{
    public class Position
    {
        public List<Bitboard> Bitboards { get; set; } // 12 bitboards
        public PieceColor SideToMove { get; set; }
        public Dictionary<CastleEnum, bool> CastlingRights { get; set; }
        public int? EnPassantSquare { get; set; }
        public int HalfMoves { get; set; }
        public int CurrentTurn { get; set; }

        private static class ZKeys
        {
            internal static readonly ulong[,] Piece = new ulong[12, 64];
            internal static readonly ulong[] Castle = new ulong[4];
            internal static readonly ulong[] EnPassantFile = new ulong[8];
            internal static readonly ulong Side;

            static ZKeys()
            {
                var rng = new Random(0x5EED_2025);

                ulong NextU64()
                {
                    Span<byte> buf = stackalloc byte[8];
                    rng.NextBytes(buf);
                    return BitConverter.ToUInt64(buf);
                }

                for (int p = 0; p < 12; ++p)
                    for (int sq = 0; sq < 64; ++sq)
                        Piece[p, sq] = NextU64();

                for (int i = 0; i < 4; ++i) Castle[i] = NextU64();
                for (int f = 0; f < 8; ++f) EnPassantFile[f] = NextU64();
                Side = NextU64();
            }
        }

        private ulong? _cachedZobrist;
        public ulong Zobrist
        {
            get
            {
                if (_cachedZobrist.HasValue) return _cachedZobrist.Value;

                ulong h = 0UL;

                // pieces
                for (int pieceIdx = 0; pieceIdx < 12; ++pieceIdx)
                {
                    ulong bb = Bitboards[pieceIdx].BitboardValue;
                    while (bb != 0)
                    {
                        int sq = BitOperations.TrailingZeroCount(bb);
                        h ^= ZKeys.Piece[pieceIdx, sq];
                        bb &= bb - 1;
                    }
                }

                // side-to-move
                if (SideToMove == PieceColor.BLACK)
                    h ^= ZKeys.Side;

                // castling rights
                if (CastlingRights[CastleEnum.WhiteShortCastle]) h ^= ZKeys.Castle[0];
                if (CastlingRights[CastleEnum.WhiteLongCastle]) h ^= ZKeys.Castle[1];
                if (CastlingRights[CastleEnum.BlackShortCastle]) h ^= ZKeys.Castle[2];
                if (CastlingRights[CastleEnum.BlackLongCastle]) h ^= ZKeys.Castle[3];

                // en-passant 
                if (EnPassantSquare is int epSq)
                {
                    int file = epSq % 8;
                    h ^= ZKeys.EnPassantFile[file];
                }

                _cachedZobrist = h;
                return h;
            }
        }

        // every method, that update position (ApplyMove) — should call InvalidateZobrist
        private void InvalidateZobrist() => _cachedZobrist = null;

        public Position()
        {
            Bitboards = new List<Bitboard>();
            for (int i = 0; i < 12; i++)
            {
                Bitboards.Add(new Bitboard());
            }

            SideToMove = PieceColor.WHITE;

            CastlingRights = new Dictionary<CastleEnum, bool>
            {
                { CastleEnum.WhiteShortCastle, false },
                { CastleEnum.WhiteLongCastle, false },
                { CastleEnum.BlackShortCastle, false },
                { CastleEnum.BlackLongCastle, false } 
            };

            EnPassantSquare = null;
            HalfMoves = 0;
            CurrentTurn = 1;
        }

        /// <summary>
        /// Creates a deep copy of the current Position.
        /// </summary>
        /// <returns>A cloned Position object.</returns>
        public Position Clone()
        {
            Position cloned = new Position
            {
                SideToMove = this.SideToMove,
                HalfMoves = this.HalfMoves,
                CurrentTurn = this.CurrentTurn,
                EnPassantSquare = this.EnPassantSquare
            };

            cloned.CastlingRights = new Dictionary<CastleEnum, bool>(this.CastlingRights);
            cloned.Bitboards = new List<Bitboard>();
            foreach (var bb in this.Bitboards)
            {
                cloned.Bitboards.Add(new Bitboard(bb.BitboardValue));
            }

            return cloned;
        }

        /// <summary>
        /// Adds a piece to the bitboards based on piece type, color, and square.
        /// </summary>
        /// <param name="piece">The type of the piece.</param>
        /// <param name="color">The color of the piece.</param>
        /// <param name="square">The square index (0-63).</param>
        public void AddPiece(PieceEnum piece, PieceColor color, int square)
        {
            Bitboards[(int)piece + (int)color].AddPiece(square);
        }

        /// <summary>
        /// Adds a piece to the bitboards based on internal piece code and square.
        /// </summary>
        /// <param name="intPiece">The internal piece code.</param>
        /// <param name="square">The square index (0-63).</param>
        public void AddPieceByInt(int intPiece, int square)
        {
            Bitboards[intPiece].AddPiece(square);
        }

        /// <summary>
        /// Retrieves the bitboard for a specific piece and color.
        /// </summary>
        /// <param name="piece">The type of the piece.</param>
        /// <param name="color">The color of the piece.</param>
        /// <returns>The corresponding Bitboard.</returns>
        public Bitboard GetPieceBitboard(PieceEnum piece, PieceColor color)
        {
            return Bitboards[(int)piece + (int)color];
        }

        /// <summary>
        /// Retrieves the combined bitboard for all white pieces.
        /// </summary>
        /// <returns>The white bitboard.</returns>
        public Bitboard GetWhiteBitboard()
        {
            Bitboard whiteBitboard = new Bitboard();
            for (int i = 6; i < 12; i++)
            {
                whiteBitboard += Bitboards[i];
            }
            return whiteBitboard;
        }

        /// <summary>
        /// Retrieves the combined bitboard for all black pieces.
        /// </summary>
        /// <returns>The black bitboard.</returns>
        public Bitboard GetBlackBitboard()
        {
            Bitboard blackBitboard = new Bitboard();
            for (int i = 0; i < 6; i++)
            {
                blackBitboard += Bitboards[i];
            }
            return blackBitboard;
        }

        /// <summary>
        /// Retrieves the combined bitboard for all pieces.
        /// </summary>
        /// <returns>The combined bitboard.</returns>
        public Bitboard GetAllBitboard()
        {
            return GetWhiteBitboard() + GetBlackBitboard();
        }

        /// <summary>
        /// Retrieves all pieces of a specific color or both colors.
        /// </summary>
        /// <param name="color">The color to filter by. Null for both.</param>
        /// <returns>A list of bitboards representing the pieces.</returns>
        public List<ulong> GetAllPiecesOfOneColor(PieceColor? color = null)
        {
            Bitboard bitboard;
            if (color == null)
            {
                bitboard = GetAllBitboard();
            }
            else if (color == PieceColor.WHITE)
            {
                bitboard = GetWhiteBitboard();
            }
            else
            {
                bitboard = GetBlackBitboard();
            }

            List<ulong> pieces = new List<ulong>();
            ulong mask = 1UL;
            for (int i = 0; i < 64; i++)
            {
                if ((bitboard.BitboardValue & mask) != 0)
                {
                    pieces.Add(mask);
                }
                mask <<= 1;
            }

            return pieces;
        }

        /// <summary>
        /// Retrieves separate bitboards for a specific piece and color.
        /// </summary>
        /// <param name="piece">The type of the piece.</param>
        /// <param name="color">The color of the piece.</param>
        /// <returns>A list of bitboards representing each piece instance.</returns>
        public List<ulong> GetSeparatePiece(PieceEnum piece, PieceColor color)
        {
            List<ulong> separatePieces = new List<ulong>();
            List<ulong> allPiecesOfColor = GetAllPiecesOfOneColor(color);

            Bitboard pieceBitboard = GetPieceBitboard(piece, color);

            foreach (var pieceMask in allPiecesOfColor)
            {
                if ((pieceMask & pieceBitboard.BitboardValue) != 0)
                {
                    separatePieces.Add(pieceMask);
                }
            }

            return separatePieces;
        }

        /// <summary>
        /// Applies a move to the current position.
        /// </summary>
        /// <param name="moveDetail">Details of the move to apply.</param>
        public void ApplyMove(MoveDetail moveDetail)
        {
            InvalidateZobrist();

            EnPassantSquare = null;

            ulong oldSquare = SquareHelper.GetBitboardFromNum(moveDetail.Square);
            ulong newSquare = SquareHelper.GetBitboardFromNum(moveDetail.Move);

            // Handling castling moves
            if (moveDetail.Piece == PieceEnum.KING && Math.Abs(moveDetail.Square - moveDetail.Move) == 2)
            {
                if (SideToMove == PieceColor.WHITE)
                {
                    if (moveDetail.Square - moveDetail.Move == -2 && CastlingRights[CastleEnum.WhiteShortCastle]) // short castle
                    {
                        // Rook move
                        ApplyMove(new MoveDetail
                        {
                            Piece = PieceEnum.ROOK,
                            Color = PieceColor.WHITE,
                            Square = 63, // h1
                            Move = 61    // f1
                        });

                        // King move
                        Bitboards[(int)PieceEnum.KING + (int)PieceColor.WHITE].BitboardValue = newSquare;
                    }
                    else if (moveDetail.Square - moveDetail.Move == 2 && CastlingRights[CastleEnum.WhiteLongCastle]) // long castle
                    {
                        // Rook move
                        ApplyMove(new MoveDetail
                        {
                            Piece = PieceEnum.ROOK,
                            Color = PieceColor.WHITE,
                            Square = 56, // a1
                            Move = 59    // d1
                        });

                        // King move
                        Bitboards[(int)PieceEnum.KING + (int)PieceColor.WHITE].BitboardValue = newSquare;
                    }

                    // Update castling rights
                    CastlingRights[CastleEnum.WhiteShortCastle] = false;
                    CastlingRights[CastleEnum.WhiteLongCastle] = false;
                }
                else // Black
                {
                    if (moveDetail.Square - moveDetail.Move == -2 && CastlingRights[CastleEnum.BlackShortCastle]) // short castle
                    {
                        // Rook move
                        ApplyMove(new MoveDetail
                        {
                            Piece = PieceEnum.ROOK,
                            Color = PieceColor.BLACK,
                            Square = 7, // h8
                            Move = 5    // f8
                        });

                        // King move
                        Bitboards[(int)PieceEnum.KING + (int)PieceColor.BLACK].BitboardValue = newSquare;
                    }
                    else if (moveDetail.Square - moveDetail.Move == 2 && CastlingRights[CastleEnum.BlackLongCastle]) // long castle
                    {
                        // Rook move
                        ApplyMove(new MoveDetail
                        {
                            Piece = PieceEnum.ROOK,
                            Color = PieceColor.BLACK,
                            Square = 0, // a8
                            Move = 3    // d8
                        });

                        // King move
                        Bitboards[(int)PieceEnum.KING + (int)PieceColor.BLACK].BitboardValue = newSquare;
                    }

                    // Update castling rights
                    CastlingRights[CastleEnum.BlackShortCastle] = false;
                    CastlingRights[CastleEnum.BlackLongCastle] = false;
                }

                return;
            }

            // Handle pawn moves
            if (moveDetail.Piece == PieceEnum.PAWN)
            {
                if (moveDetail.Color == PieceColor.WHITE && moveDetail.Square - moveDetail.Move == 16)
                {
                    EnPassantSquare = moveDetail.Square - 8;
                }
                if (moveDetail.Color == PieceColor.BLACK && moveDetail.Square - moveDetail.Move == -16)
                {
                    EnPassantSquare = moveDetail.Square + 8;
                }

                // TODO: TODO: TODO: TODO: !!!
                // En passant capture
                if (moveDetail.Color == PieceColor.WHITE && (moveDetail.Square - moveDetail.Move == 7 || moveDetail.Square - moveDetail.Move == 9))
                {
                    int capturePawnSquareNum = moveDetail.Move + 8;
                    if (capturePawnSquareNum == EnPassantSquare)
                    {
                        Bitboards[(int)PieceEnum.PAWN + (int)PieceColor.BLACK].BitboardValue &= ~SquareHelper.GetBitboardFromNum(capturePawnSquareNum);
                    }
                   
                }
                if (moveDetail.Color == PieceColor.BLACK && (moveDetail.Square - moveDetail.Move == -7 || moveDetail.Square - moveDetail.Move == -9))
                {
                    int capturePawnSquareNum = moveDetail.Move - 8;
                    if (capturePawnSquareNum == EnPassantSquare)
                    {
                        Bitboards[(int)PieceEnum.PAWN + (int)PieceColor.WHITE].BitboardValue &= ~SquareHelper.GetBitboardFromNum(capturePawnSquareNum);
                    }
                }
            }

            // Remove any captured piece
            bool isCapture = false;
            for (int colorNum = 0; colorNum < 2; colorNum++)
            {
                int baseIndex = colorNum * 6;
                for (int pieceNum = 0; pieceNum < 6; pieceNum++)
                {
                    int index = baseIndex + pieceNum;
                    if ((Bitboards[index].BitboardValue & newSquare) != 0)
                    {
                        Bitboards[index].BitboardValue &= ~newSquare;
                        isCapture = true;
                        if (moveDetail.Move == 0)
                        {
                            CastlingRights[CastleEnum.BlackLongCastle] = false;
                        }
                        else if (moveDetail.Move == 7)
                        {
                            CastlingRights[CastleEnum.BlackShortCastle] = false;
                        }
                        else if (moveDetail.Move == 56)
                        {
                            CastlingRights[CastleEnum.WhiteLongCastle] = false;
                        }
                        else if (moveDetail.Move == 63)
                        {
                            CastlingRights[CastleEnum.WhiteShortCastle] = false;
                        }
                    }
                }
            }

            // Move the piece
            Bitboards[(int)moveDetail.Piece + (int)moveDetail.Color].BitboardValue =
                (Bitboards[(int)moveDetail.Piece + (int)moveDetail.Color].BitboardValue & ~oldSquare) | newSquare;

            // Handle pawn promotion
            if (moveDetail.Piece == PieceEnum.PAWN)
            {
                if (moveDetail.Color == PieceColor.WHITE && moveDetail.Move >= 0 && moveDetail.Move <= 7)
                {
                    ulong promotionBitboard = SquareHelper.GetBitboardFromNum(moveDetail.Move);
                    Bitboards[(int)PieceEnum.PAWN + (int)PieceColor.WHITE].BitboardValue &= ~promotionBitboard;
                    
                    // for checking legal moves (assuming promotion to a queen)
                    if (moveDetail.Promotion == PieceEnum.NONE) 
                    {
                        Bitboards[(int)PieceEnum.QUEEN + (int)PieceColor.WHITE].BitboardValue |= promotionBitboard;
                    }
                    else
                    {
                        PieceEnum newPiece = moveDetail.Promotion;
                        Bitboards[(int)newPiece + (int)PieceColor.WHITE].BitboardValue |= promotionBitboard;
                    }   
                }

                if (moveDetail.Color == PieceColor.BLACK && moveDetail.Move >= 56 && moveDetail.Move <= 63)
                {
                    ulong promotionBitboard = SquareHelper.GetBitboardFromNum(moveDetail.Move);
                    Bitboards[(int)PieceEnum.PAWN + (int)PieceColor.BLACK].BitboardValue &= ~promotionBitboard;

                    // for checking legal moves (assuming promotion to a queen)
                    if (moveDetail.Promotion == PieceEnum.NONE)
                    {
                        Bitboards[(int)PieceEnum.QUEEN + (int)PieceColor.BLACK].BitboardValue |= promotionBitboard;
                    }
                    else
                    {
                        PieceEnum newPiece = moveDetail.Promotion;
                        Bitboards[(int)newPiece + (int)PieceColor.BLACK].BitboardValue |= promotionBitboard;
                    }
                }
            }

            // Update castling rights if King or Rook moves or Rook is captured
            if (moveDetail.Piece == PieceEnum.KING)
            {
                if (moveDetail.Color == PieceColor.WHITE)
                {
                    CastlingRights[CastleEnum.WhiteShortCastle] = false;
                    CastlingRights[CastleEnum.WhiteLongCastle] = false;
                }
                else if (moveDetail.Color == PieceColor.BLACK)
                {
                    CastlingRights[CastleEnum.BlackShortCastle] = false;
                    CastlingRights[CastleEnum.BlackLongCastle] = false;
                }
            }

            if (moveDetail.Piece == PieceEnum.ROOK)
            {
                if (moveDetail.Color == PieceColor.WHITE)
                {
                    if ((Bitboards[(int)PieceEnum.ROOK + (int)PieceColor.WHITE].BitboardValue & SquareHelper.GetBitboardFromNum(63)) == 0)
                    {
                        CastlingRights[CastleEnum.WhiteShortCastle] = false;
                    }
                    if ((Bitboards[(int)PieceEnum.ROOK + (int)PieceColor.WHITE].BitboardValue & SquareHelper.GetBitboardFromNum(56)) == 0)
                    {
                        CastlingRights[CastleEnum.WhiteLongCastle] = false;
                    }
                }
                else if (moveDetail.Color == PieceColor.BLACK)
                {
                    if ((Bitboards[(int)PieceEnum.ROOK + (int)PieceColor.BLACK].BitboardValue & SquareHelper.GetBitboardFromNum(7)) == 0)
                    {
                        CastlingRights[CastleEnum.BlackShortCastle] = false;
                    }
                    if ((Bitboards[(int)PieceEnum.ROOK + (int)PieceColor.BLACK].BitboardValue & SquareHelper.GetBitboardFromNum(0)) == 0)
                    {
                        CastlingRights[CastleEnum.BlackLongCastle] = false;
                    }
                }
            }

            // Change side to move
            ChangeSideToMove();

            // Update fullmove number
            if (SideToMove == PieceColor.WHITE)
            {
                CurrentTurn += 1;
            }

            // Update halfmove counter
            if (isCapture || moveDetail.Piece == PieceEnum.PAWN)
            {
                HalfMoves = 0;
            }
            else
            {
                HalfMoves += 1;
            }
        }

        /// <summary>
        /// Retrieves all separate moves from the current position.
        /// </summary>
        /// <returns>A list of MoveDetail objects.</returns>
        public List<MoveDetail> GetAllSeparateMoves()
        {
            List<MoveDetail> separatedMoves = new List<MoveDetail>();

            List<MoveDetail> allMoves = GetAllMoves();

            foreach (var move in allMoves)
            {
                if (move.Piece == PieceEnum.PAWN && (
                    (move.Color == PieceColor.WHITE && new[] { 8, 9, 10, 11, 12, 13, 14, 15 }.Contains(move.Square)) ||
                    (move.Color == PieceColor.BLACK && new[] { 48, 49, 50, 51, 52, 53, 54, 55 }.Contains(move.Square)))
                    )
                {
                    foreach (var piece in new[] { PieceEnum.QUEEN, PieceEnum.ROOK, PieceEnum.BISHOP, PieceEnum.KNIGHT })
                    {
                        separatedMoves.Add(new MoveDetail
                        {
                            Piece = move.Piece,
                            Color = move.Color,
                            Square = move.Square,
                            Move = move.Move,
                            Promotion = piece,
                            Priority = move.Priority,
                        });
                    }
                }
                else
                {
                    separatedMoves.Add(move);
                } 
            }

            return separatedMoves;
        }

        /// <summary>
        /// Retrieves all legal moves from the current position.
        /// </summary>
        /// <returns>A list of MoveDetail objects.</returns>
        public List<MoveDetail> GetAllMoves()
        {
            IEnumerable<MoveDetail> pseudoLegalMoves = GenerateAllPseudoLegalMoves();
            IEnumerable<MoveDetail> legalMoves = FilterIllegalMoves(pseudoLegalMoves);

            return legalMoves.Reverse().ToList();
        }

        /// <summary>
        /// Retrieves all legal moves sorted by priority.
        /// </summary>
        /// <returns>A list of MoveDetail objects sorted by priority descending.</returns>
        public List<MoveDetail> GetMovesSortedByPriority() // TODO: redo this to MVV-LVA it
        {
            List<MoveDetail> separateMoves = GetAllSeparateMoves();
            List<MoveDetail> prioritizedMoves = new List<MoveDetail>();

            foreach (var move in separateMoves)
            {
                Position positionCopy = this.Clone();
                positionCopy.ApplyMove(move);
                bool isCheck = positionCopy.IsKingInCheck();

                bool isCapture = this.IsCapture(move);

                if (isCheck && isCapture)
                {
                    move.Priority = 14;
                }
                else if (isCheck)
                {
                    move.Priority = 10;
                }
                else if (isCapture)
                {
                    move.Priority = 4;
                }
                else
                {
                    move.Priority = 0;
                }

                prioritizedMoves.Add(move);
            }

            return prioritizedMoves.OrderByDescending(m => m.Priority).ToList();
        }

        /// <summary>
        /// Filters pseudo-legal moves to retain only legal moves.
        /// </summary>
        /// <param name="pseudoLegalMoves">An enumerable of pseudo-legal MoveDetail objects.</param>
        /// <returns>An enumerable of legal MoveDetail objects.</returns>
        private IEnumerable<MoveDetail> FilterIllegalMoves(IEnumerable<MoveDetail> pseudoLegalMoves)
        {
            List<MoveDetail> separateLegalMoves = new List<MoveDetail>();

            foreach (var move in pseudoLegalMoves)
            {
                if (IsMoveLegal(move))
                {
                    separateLegalMoves.Add(move);
                }
            }

            return separateLegalMoves;
        }

        /// <summary>
        /// Determines if a move is legal.
        /// </summary>
        /// <param name="move">The move to evaluate.</param>
        /// <returns>True if the move is legal; otherwise, false.</returns>
        public bool IsMoveLegal(MoveDetail move)
        {
            // Handle castling moves
            if (move.Piece == PieceEnum.KING && Math.Abs(move.Square - move.Move) == 2)
            {
                if (
                    move.Color == PieceColor.WHITE && move.Square - move.Move == -2 && CastlingRights[CastleEnum.WhiteShortCastle] ||
                    move.Color == PieceColor.WHITE && move.Square - move.Move == 2 && CastlingRights[CastleEnum.WhiteLongCastle] ||
                    move.Color == PieceColor.BLACK && move.Square - move.Move == -2 && CastlingRights[CastleEnum.BlackShortCastle] ||
                    move.Color == PieceColor.BLACK && move.Square - move.Move == 2 && CastlingRights[CastleEnum.BlackLongCastle]
                    ) {
                    int intermediateMove = (move.Square + move.Move) / 2;
                    MoveDetail intermediateMoveDetail = new MoveDetail
                    {
                        Piece = move.Piece,
                        Color = move.Color,
                        Square = move.Square,
                        Move = intermediateMove
                    };

                    // Check if squares between king and rook are empty
                    if ((GetAllBitboard().BitboardValue & (SquareHelper.GetBitboardFromNum(intermediateMove) | SquareHelper.GetBitboardFromNum(move.Move))) != 0)
                    {
                        return false;
                    }

                    // Check for king being in check before castling
                    Position positionCopy = Clone();

                    // Check for king going through check during castling
                    Position positionCopy2 = Clone();
                    positionCopy2.ApplyMove(intermediateMoveDetail);
                    positionCopy2.ChangeSideToMove();

                    // Check for king being in check after castling
                    Position positionCopy3 = Clone();
                    positionCopy3.ApplyMove(move);
                    positionCopy3.ChangeSideToMove();

                    return !positionCopy.IsKingInCheck() && !positionCopy2.IsKingInCheck() && !positionCopy3.IsKingInCheck();
                }
                return false;
            }
            else
            {
                // All other moves
                Position positionCopy = Clone();
                positionCopy.ApplyMove(move);
                positionCopy.ChangeSideToMove();

                return !positionCopy.IsKingInCheck();
            }
        }

        /// <summary>
        /// Determines the outcome of the game based on the current position.
        /// </summary>
        /// <returns>A string representing the game outcome: "win", "loss", or "draw".</returns>
        public PositionOutcome DetermineOutcome()
        {
            // TODO: Implement evaluation logic if needed

            if (IsCheckmate())
            {
                return PositionOutcome.LOSS;
            }
            else if (IsStalemate() || IsDraw())
            {
                return PositionOutcome.DRAW;
            }
            else
            {
                return PositionOutcome.WIN;
            }
        }

        /// <summary>
        /// Checks if the current player's king is in check.
        /// </summary>
        /// <returns>True if in check; otherwise, false.</returns>
        public bool IsKingInCheck()
        {
            List<ulong> kings = GetSeparatePiece(PieceEnum.KING, SideToMove);
            if (kings.Count == 0 || kings[0] == 0)
            {
                // No king found, invalid position
                return false;
            }

            ulong kingBitboard = kings[0];

            // Switch side to move to check attacks from the opponent
            ChangeSideToMove();
            List<MoveDetail> allPseudoLegalMoves = GenerateAllPseudoLegalMoves().ToList();
            ChangeSideToMove();

            foreach (var move in allPseudoLegalMoves)
            {
                if (SquareHelper.GetBitboardFromNum(move.Move) == kingBitboard)
                {
                    return true;
                }
            }

            return false;
        }

        /// <summary>
        /// Determines if a move is a capture.
        /// </summary>
        /// <param name="move">The move to evaluate.</param>
        /// <returns>True if the move captures an opponent's piece; otherwise, false.</returns>
        public bool IsCapture(MoveDetail move)
        {
            var (capturedPiece, capturedColor) = GetPieceAndColorBySquare(move.Move);
            return capturedPiece != PieceEnum.NONE && capturedColor != move.Color;
        }

        /// <summary>
        /// Checks if the current position is a checkmate.
        /// </summary>
        /// <returns>True if checkmate; otherwise, false.</returns>
        public bool IsCheckmate()
        {
            return IsKingInCheck() && GetAllMoves().Count == 0;
        }

        /// <summary>
        /// Checks if the current position is a stalemate.
        /// </summary>
        /// <returns>True if stalemate; otherwise, false.</returns>
        public bool IsStalemate()
        {
            return !IsKingInCheck() && GetAllMoves().Count == 0;
        }

        /// <summary>
        /// Checks if the game is a draw.
        /// </summary>
        /// <returns>True if the game is a draw; otherwise, false.</returns>
        public bool IsDraw()
        {
            return IsStalemate() || IsThreefoldRepetition() || IsFiftyMoveRuleReached() || IsInsufficientMaterial();
        }

        /// <summary>
        /// Checks if the threefold repetition rule has been met.
        /// </summary>
        /// <returns>True if threefold repetition; otherwise, false.</returns>
        public bool IsThreefoldRepetition()
        {
            // TODO: Implement threefold repetition logic
            return false;
        }

        /// <summary>
        /// Checks if the fifty-move rule has been reached.
        /// </summary>
        /// <returns>True if fifty-move rule is reached; otherwise, false.</returns>
        public bool IsFiftyMoveRuleReached()
        {
            return HalfMoves >= 100;
        }

        /// <summary>
        /// Checks if there is insufficient material to continue the game.
        /// </summary>
        /// <returns>True if insufficient material; otherwise, false.</returns>
        public bool IsInsufficientMaterial()
        {
            // TODO: Implement insufficient material logic
            // Example implementation for basic cases

            ulong whitePieces = GetWhiteBitboard().BitboardValue;
            ulong blackPieces = GetBlackBitboard().BitboardValue;
            ulong allPieces = whitePieces | blackPieces;

            // Only kings
            if (BitHelper.PopCount(allPieces) == 2)
            {
                return true;
            }

            // King and bishop vs King
            if (BitHelper.PopCount(allPieces) == 3)
            {
                if ((whitePieces != 0 && BitHelper.PopCount(whitePieces) == 2 && GetPieceBitboard(PieceEnum.BISHOP, PieceColor.WHITE).BitboardValue != 0) ||
                    (blackPieces != 0 && BitHelper.PopCount(blackPieces) == 2 && GetPieceBitboard(PieceEnum.BISHOP, PieceColor.BLACK).BitboardValue != 0))
                {
                    return true;
                }
            }

            // King and knight vs King
            if (BitHelper.PopCount(allPieces) == 3)
            {
                if ((whitePieces != 0 && BitHelper.PopCount(whitePieces) == 2 && GetPieceBitboard(PieceEnum.KNIGHT, PieceColor.WHITE).BitboardValue != 0) ||
                    (blackPieces != 0 && BitHelper.PopCount(blackPieces) == 2 && GetPieceBitboard(PieceEnum.KNIGHT, PieceColor.BLACK).BitboardValue != 0))
                {
                    return true;
                }
            }

            // King and bishop vs King and bishop with both bishops on the same color
            if (BitHelper.PopCount(allPieces) == 4)
            {
                bool whiteHasBishop = GetPieceBitboard(PieceEnum.BISHOP, PieceColor.WHITE).BitboardValue != 0;
                bool blackHasBishop = GetPieceBitboard(PieceEnum.BISHOP, PieceColor.BLACK).BitboardValue != 0;

                if (whiteHasBishop && blackHasBishop)
                {
                    int whiteBishopSquare = SquareHelper.GetNumFromBitboard(GetPieceBitboard(PieceEnum.BISHOP, PieceColor.WHITE).BitboardValue);
                    int blackBishopSquare = SquareHelper.GetNumFromBitboard(GetPieceBitboard(PieceEnum.BISHOP, PieceColor.BLACK).BitboardValue);

                    bool whiteBishopColor = ((whiteBishopSquare / 8) + (whiteBishopSquare % 8)) % 2 == 0;
                    bool blackBishopColor = ((blackBishopSquare / 8) + (blackBishopSquare % 8)) % 2 == 0;

                    return whiteBishopColor == blackBishopColor;
                }
            }

            return false;
        }

        /// <summary>
        /// Changes the side to move.
        /// </summary>
        public void ChangeSideToMove()
        {
            SideToMove = SideToMove == PieceColor.WHITE ? PieceColor.BLACK : PieceColor.WHITE;
        }

        /// <summary>
        /// Checks if the game is over (checkmate or draw).
        /// </summary>
        /// <returns>True if the game is over; otherwise, false.</returns>
        public bool IsGameOver()
        {
            return IsDraw() || IsCheckmate();
        }

        /// <summary>
        /// Retrieves the piece and its color at a given square.
        /// </summary>
        /// <param name="square">The square index (0-63).</param>
        /// <returns>A tuple containing the PieceEnum and PieceColor.</returns>
        public (PieceEnum, PieceColor) GetPieceAndColorBySquare(int square)
        {
            ulong bitSquare = SquareHelper.GetBitboardFromNum(square);
            for (int i = 0; i < Bitboards.Count; i++)
            {
                if ((Bitboards[i].BitboardValue & bitSquare) != 0)
                {
                    PieceEnum piece = GetPieceEnumFromIndex(i);
                    PieceColor color = GetColorFromIndex(i);
                    return (piece, color);
                }
            }

            return (PieceEnum.NONE, PieceColor.BLACK); // No piece on the square
        }

        /// <summary>
        /// Converts a bitboard index to its corresponding PieceEnum.
        /// </summary>
        /// <param name="index">The bitboard index (0-11).</param>
        /// <returns>The corresponding PieceEnum.</returns>
        private PieceEnum GetPieceEnumFromIndex(int index)
        {
            if (index >= 0 && index <= 5)
                return (PieceEnum)index;
            else if (index >= 6 && index <= 11)
                return (PieceEnum)(index - 6);
            else
                throw new ArgumentException($"Invalid bitboard index: {index}");
        }

        /// <summary>
        /// Determines the color of a piece based on its bitboard index.
        /// </summary>
        /// <param name="index">The bitboard index (0-11).</param>
        /// <returns>The corresponding PieceColor.</returns>
        private PieceColor GetColorFromIndex(int index)
        {
            return index >= 6 ? PieceColor.WHITE : PieceColor.BLACK;
        }

        /// <summary>
        /// Generates all pseudo-legal moves from the current position.
        /// </summary>
        /// <returns>An enumerable of MoveDetail objects representing pseudo-legal moves.</returns>
        public IEnumerable<MoveDetail> GenerateAllPseudoLegalMoves()
        {
            List<MoveDetail> result = new List<MoveDetail>();

            // Determine which color is moving
            PieceColor currentColor = SideToMove;
            Bitboard ownBitboard = currentColor == PieceColor.WHITE ? GetWhiteBitboard() : GetBlackBitboard();
            Bitboard enemyBitboard = currentColor == PieceColor.WHITE ? GetBlackBitboard() : GetWhiteBitboard();

            // Iterate through all piece types
            for (int numPiece = 0; numPiece < 6; numPiece++)
            {
                PieceEnum pieceEnum = (PieceEnum)numPiece;
                List<ulong> separateBitPieces = GetSeparatePiece(pieceEnum, currentColor);

                foreach (var bitPiece in separateBitPieces)
                {
                    if (bitPiece == 0) continue; // Skip empty bitboards
                    List<ulong> separatePiecesMoves = new List<ulong>();
                    switch (pieceEnum)
                    {
                        case PieceEnum.PAWN:
                            ulong? enPassant = EnPassantSquare.HasValue ? SquareHelper.GetBitboardFromNum(EnPassantSquare.Value) : (ulong?)null;
                            separatePiecesMoves = GenerateMove.GeneratePawnMove(bitPiece, new Bitboard(ownBitboard.BitboardValue), new Bitboard(enemyBitboard.BitboardValue), currentColor, enPassant).ToList();
                            break;
                        case PieceEnum.KNIGHT:
                            separatePiecesMoves = GenerateMove.GenerateKnightMove(bitPiece, new Bitboard(ownBitboard.BitboardValue)).ToList();
                            break;
                        case PieceEnum.BISHOP:
                            separatePiecesMoves = GenerateMove.GenerateBishopMove(bitPiece, new Bitboard(ownBitboard.BitboardValue), new Bitboard(enemyBitboard.BitboardValue)).ToList();
                            break;
                        case PieceEnum.ROOK:
                            separatePiecesMoves = GenerateMove.GenerateRookMove(bitPiece, new Bitboard(ownBitboard.BitboardValue), new Bitboard(enemyBitboard.BitboardValue)).ToList();
                            break;
                        case PieceEnum.QUEEN:
                            separatePiecesMoves = GenerateMove.GenerateQueenMove(bitPiece, new Bitboard(ownBitboard.BitboardValue), new Bitboard(enemyBitboard.BitboardValue)).ToList();
                            break;
                        case PieceEnum.KING:
                            separatePiecesMoves = GenerateMove.GenerateKingMove(bitPiece, new Bitboard(ownBitboard.BitboardValue), new Bitboard(enemyBitboard.BitboardValue), currentColor == PieceColor.WHITE).ToList();
                            break;
                    }

                    //List<int> pieceMovesNum = separatePiecesMoves.Select(SquareHelper.GetNumFromBitboard).ToList();
                    List<int> pieceMovesNum = separatePiecesMoves
                .Where(m => m != 0) // Make sure there are no zero moves
                .Select(SquareHelper.GetNumFromBitboard)
                .ToList();

                    foreach (var move in pieceMovesNum)
                    {
                        result.Add(new MoveDetail
                        {
                            Piece = pieceEnum,
                            Color = currentColor,
                            Square = SquareHelper.GetNumFromBitboard(bitPiece),
                            Move = move
                        });
                    }
                }
            }

            return result;
        }
    }
}
