using ChessEngine.Enums;
using ChessEngine.Helpers;
using ChessEngine.Domain.Models;

namespace ChessEngine.Application.Fen
{
    public static class FenGenerator
    {
        /// <summary>
        /// Sets the piece placement on the board based on the FEN string.
        /// </summary>
        /// <param name="position">The Position object to modify.</param>
        /// <param name="piecesPlacementFen">The piece placement part of the FEN string.</param>
        public static void SetPiecePlacement(Position position, string piecesPlacementFen)
        {
            int square = 0; // Start from square 0 (a8)

            foreach (char c in piecesPlacementFen)
            {
                if (char.IsDigit(c))
                {
                    square += (int)char.GetNumericValue(c); // Skip empty squares
                }
                else if (c == '/')
                {
                    continue; // Move to the next rank
                }
                else
                {
                    int? pieceCode = FenCharToPieceCode(c);
                    if (pieceCode.HasValue)
                    {
                        position.AddPieceByInt(pieceCode.Value, square);
                        square += 1; // Move to the next square
                    }
                    else
                    {
                        throw new ArgumentException($"Invalid FEN character: {c}");
                    }
                }
            }
        }

        /// <summary>
        /// Sets the side to move based on the FEN string.
        /// </summary>
        /// <param name="position">The Position object to modify.</param>
        /// <param name="sideToMoveFen">The side to move part of the FEN string ("w" or "b").</param>
        public static void SetSideToMove(Position position, string sideToMoveFen)
        {
            if (sideToMoveFen == "w")
            {
                position.SideToMove = PieceColor.WHITE;
            }
            else if (sideToMoveFen == "b")
            {
                position.SideToMove = PieceColor.BLACK;
            }
            else
            {
                throw new ArgumentException($"Invalid side to move value in FEN: {sideToMoveFen}");
            }
        }

        /// <summary>
        /// Sets the castling rights based on the FEN string.
        /// </summary>
        /// <param name="position">The Position object to modify.</param>
        /// <param name="castlingRightsFen">The castling rights part of the FEN string.</param>
        public static void SetCastlingRights(Position position, string castlingRightsFen)
        {
            // Reset all castling rights
            foreach (CastleEnum castling in Enum.GetValues(typeof(CastleEnum)))
            {
                position.CastlingRights[castling] = false;
            }

            foreach (char c in castlingRightsFen)
            {
                switch (c)
                {
                    case 'K':
                        position.CastlingRights[CastleEnum.WhiteShortCastle] = true;
                        break;
                    case 'Q':
                        position.CastlingRights[CastleEnum.WhiteLongCastle] = true;
                        break;
                    case 'k':
                        position.CastlingRights[CastleEnum.BlackShortCastle] = true;
                        break;
                    case 'q':
                        position.CastlingRights[CastleEnum.BlackLongCastle] = true;
                        break;
                    default:
                        // Ignore unexpected characters
                        break;
                }
            }
        }

        /// <summary>
        /// Sets the en passant target square based on the FEN string.
        /// </summary>
        /// <param name="position">The Position object to modify.</param>
        /// <param name="enPassantFen">The en passant part of the FEN string.</param>
        public static void SetEnPassant(Position position, string enPassantFen)
        {
            if (enPassantFen != "-")
            {
                position.EnPassantSquare = SquareHelper.GetNumBySquareName(enPassantFen);
            }
            else
            {
                position.EnPassantSquare = null;
            }
        }

        /// <summary>
        /// Parses a FEN string and returns the corresponding Position object.
        /// </summary>
        /// <param name="fen">The full FEN string.</param>
        /// <returns>A Position object representing the chessboard state.</returns>
        public static Position GetPositionFromFen(string fen)
        {
            if (string.IsNullOrWhiteSpace(fen))
                throw new ArgumentException("FEN string cannot be null or empty.", nameof(fen));

            string[] fenParts = fen.Split(' ');
            if (fenParts.Length != 6)
                throw new ArgumentException("Invalid FEN string format.", nameof(fen));

            Position position = new Position();

            // Set piece placement
            SetPiecePlacement(position, fenParts[0]);

            // Set side to move
            SetSideToMove(position, fenParts[1]);

            // Set castling rights
            SetCastlingRights(position, fenParts[2]);

            // Set en passant target square
            SetEnPassant(position, fenParts[3]);

            // Set halfmove clock
            if (!int.TryParse(fenParts[4], out int halfMoves))
                throw new ArgumentException("Invalid halfmove clock in FEN.", nameof(fen));
            position.HalfMoves = halfMoves;

            // Set fullmove number
            if (!int.TryParse(fenParts[5], out int currentTurn))
                throw new ArgumentException("Invalid fullmove number in FEN.", nameof(fen));
            position.CurrentTurn = currentTurn;

            return position;
        }

        /// <summary>
        /// Generates a FEN string from a Position object.
        /// </summary>
        /// <param name="position">The Position object representing the chessboard state.</param>
        /// <returns>A FEN string representing the current position.</returns>
        public static string GetFenFromPosition(Position position)
        {
            if (position == null)
                throw new ArgumentNullException(nameof(position));

            List<string> fenRows = new List<string>();
            int emptyCount = 0;

            for (int i = 0; i < 64; i++)
            {
                var (piece, color) = position.GetPieceAndColorBySquare(i);
                if (piece == PieceEnum.NONE)
                {
                    emptyCount++;
                }  
                else
                {
                    int pieceCode = GetInternalPieceCode(piece, color);
                    if (emptyCount > 0)
                    {
                        fenRows.Add(emptyCount.ToString());
                        emptyCount = 0;
                    }

                    string fenChar = PieceCodeToFenChar(pieceCode);
                    if (fenChar == null)
                        throw new InvalidOperationException($"Invalid piece code at square {i}: {pieceCode}");

                    fenRows.Add(fenChar);
                }

                if ((i + 1) % 8 == 0)
                {
                    if (emptyCount > 0)
                    {
                        fenRows.Add(emptyCount.ToString());
                        emptyCount = 0;
                    }

                    if (i != 63)
                        fenRows.Add("/");
                }
            }

            // Active color
            string activeColor = position.SideToMove == PieceColor.WHITE ? "w" : "b";

            // Castling availability
            List<char> castlingRights = new List<char>();
            foreach (var castling in position.CastlingRights)
            {
                if (castling.Value)
                {
                    switch (castling.Key)
                    {
                        case CastleEnum.WhiteShortCastle:
                            castlingRights.Add('K');
                            break;
                        case CastleEnum.WhiteLongCastle:
                            castlingRights.Add('Q');
                            break;
                        case CastleEnum.BlackShortCastle:
                            castlingRights.Add('k');
                            break;
                        case CastleEnum.BlackLongCastle:
                            castlingRights.Add('q');
                            break;
                    }
                }
            }
            string castlingAvailability = castlingRights.Count > 0 ? new string(castlingRights.ToArray()) : "-";

            // En passant target square
            string enPassant = position.EnPassantSquare.HasValue
                ? SquareHelper.GetSquareNameByNum(position.EnPassantSquare.Value)
                : "-";

            // Halfmove clock and fullmove number
            string halfMoves = position.HalfMoves.ToString();
            string fullMoves = position.CurrentTurn.ToString();

            // Combine all parts
            string fen = string.Join(" ", new string[]
            {
                string.Concat(fenRows),
                activeColor,
                castlingAvailability,
                enPassant,
                halfMoves,
                fullMoves
            });

            return fen;
        }

        /// <summary>
        /// Converts an internal piece code to its corresponding FEN character.
        /// </summary>
        /// <param name="pieceCode">The internal piece code.</param>
        /// <returns>The corresponding FEN character, or null if invalid.</returns>
        private static string PieceCodeToFenChar(int pieceCode)
        {
            var mapping = new Dictionary<int, string>
            {
                {6, "P"}, {7, "N"}, {8, "B"}, {9, "R"}, {10, "Q"}, {11, "K"},
                {0, "p"}, {1, "n"}, {2, "b"}, {3, "r"}, {4, "q"}, {5, "k"}
            };

            return mapping.ContainsKey(pieceCode) ? mapping[pieceCode] : null;
        }

        /// <summary>
        /// Converts a FEN character to its corresponding internal piece code.
        /// </summary>
        /// <param name="c">The FEN character.</param>
        /// <returns>The corresponding internal piece code, or null if invalid.</returns>
        private static int? FenCharToPieceCode(char c)
        {
            var mapping = new Dictionary<char, int>
            {
                {'P', 6}, {'N', 7}, {'B', 8}, {'R', 9}, {'Q', 10}, {'K', 11},
                {'p', 0}, {'n', 1}, {'b', 2}, {'r', 3}, {'q', 4}, {'k', 5}
            };

            return mapping.ContainsKey(c) ? (int?)mapping[c] : null;
        }

        /// <summary>
        /// Converts a piece and its color to the internal piece code.
        /// </summary>
        /// <param name="piece">The type of the piece.</param>
        /// <param name="color">The color of the piece.</param>
        /// <returns>The internal piece code.</returns>
        private static int GetInternalPieceCode(PieceEnum piece, PieceColor color)
        {
            // Assuming the internal piece codes are defined as follows: +++1
            // WHITE_PAWN = 6, WHITE_KNIGHT = 7, WHITE_BISHOP = 8, WHITE_ROOK = 9,
            // WHITE_QUEEN = 10, WHITE_KING = 11
            // BLACK_PAWN = 0, BLACK_KNIGHT = 1, BLACK_BISHOP = 2, BLACK_ROOK = 3,
            // BLACK_QUEEN = 4, BLACK_KING = 5

            if (color == PieceColor.WHITE)
            {
                switch (piece)
                {
                    case PieceEnum.PAWN:
                        return 6;
                    case PieceEnum.KNIGHT:
                        return 7;
                    case PieceEnum.BISHOP:
                        return 8;
                    case PieceEnum.ROOK:
                        return 9;
                    case PieceEnum.QUEEN:
                        return 10;
                    case PieceEnum.KING:
                        return 11;
                    default:
                        throw new ArgumentException($"Invalid piece type: {piece}");
                }
            }
            else // BLACK
            {
                switch (piece)
                {
                    case PieceEnum.PAWN:
                        return 0;
                    case PieceEnum.KNIGHT:
                        return 1;
                    case PieceEnum.BISHOP:
                        return 2;
                    case PieceEnum.ROOK:
                        return 3;
                    case PieceEnum.QUEEN:
                        return 4;
                    case PieceEnum.KING:
                        return 5;
                    default:
                        throw new ArgumentException($"Invalid piece type: {piece}");
                }
            }
        }
    }
}
