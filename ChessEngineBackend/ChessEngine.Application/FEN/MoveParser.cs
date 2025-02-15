using System;
using System.Collections.Generic;
using System.Text.RegularExpressions;
using ChessEngine.Enums;
using ChessEngine.Helpers;
using ChessEngine.Domain.Models;
using ChessEngine.FEN;
using ChessEngine.Application.Fen;
using ChessEngine.Domain.Models;

namespace ChessEngine.FEN
{
    public class MoveParser
    {
        private static readonly Dictionary<char, PieceEnum> PieceMap = new Dictionary<char, PieceEnum>
        {
            { 'P', PieceEnum.PAWN },
            { 'N', PieceEnum.KNIGHT },
            { 'B', PieceEnum.BISHOP },
            { 'R', PieceEnum.ROOK },
            { 'Q', PieceEnum.QUEEN },
            { 'K', PieceEnum.KING }
        };

        public MoveDetail ConvertSanToMoveDetail(string prevFen, string san)
        {
            Position position = FenGenerator.GetPositionFromFen(prevFen);
            if (san == "O-O")
            {
                if (position.SideToMove == PieceColor.WHITE)
                {
                    return new MoveDetail
                    {
                        Piece = PieceEnum.KING,
                        Color = PieceColor.WHITE,
                        Square = 60,
                        Move = 62
                    };
                }
                else
                {
                    return new MoveDetail
                    {
                        Piece = PieceEnum.KING,
                        Color = PieceColor.BLACK,
                        Square = 4,
                        Move = 6
                    };
                }
            }
            else if (san == "O-O-O")
            {
                if (position.SideToMove == PieceColor.WHITE)
                {
                    return new MoveDetail
                    {
                        Piece = PieceEnum.KING,
                        Color = PieceColor.WHITE,
                        Square = 60,
                        Move = 58
                    };
                }
                else
                {
                    return new MoveDetail
                    {
                        Piece = PieceEnum.KING,
                        Color = PieceColor.BLACK,
                        Square = 4,
                        Move = 2
                    };
                }
            }
            else
            {
                var parsed = ParseSan(san);

                char pieceChar = parsed.Piece;
                PieceEnum pieceEnum = GetPieceEnum(pieceChar);

                PieceColor color = position.SideToMove;

                int toSquare = ConvertFileRankToSquare(parsed.ToFile, parsed.ToRank);

                // TODO: revise char? and char in ParseSan
                int fromSquare = parsed.FromFile != null && parsed.FromRank != null
                                 ? ConvertFileRankToSquare(parsed.FromFile, parsed.FromRank)
                                 : FindFromSquare(position, pieceEnum, color, toSquare);

                return new MoveDetail
                {
                    Piece = pieceEnum,
                    Color = color,
                    Square = fromSquare,
                    Move = toSquare
                };
            }
        }

        private int FindFromSquare(Position position, PieceEnum piece, PieceColor color, int toSquare)
        {
            // Implement logic to find the originating square if not specified
            // This can involve generating all possible moves and matching the destination
            // For simplicity, return 0 here.
            return 0;
        }

        private ParsedSan ParseSan(string san)
        {
            string pattern = @"^(?<piece>[KQRBN])?(?<from_file>[a-h])?(?<from_rank>[1-8])?(?<capture>x)?(?<to_file>[a-h])(?<to_rank>[1-8])(?<promotion>=[QRBN])?(?<check>[+#])?$";
            var match = Regex.Match(san, pattern);
            if (!match.Success)
                throw new ArgumentException($"Invalid SAN move: {san}");

            return new ParsedSan
            {
                // Piece = match.Groups["piece"].Value.Length > 0 ? match.Groups["piece"].Value[0] : (char?)null,
                Piece = match.Groups["piece"].Value[0],
                // FromFile = match.Groups["from_file"].Value.Length > 0 ? match.Groups["from_file"].Value[0] : (char?)null,
                FromFile = match.Groups["from_file"].Value[0],
                // FromRank = match.Groups["from_rank"].Value.Length > 0 ? match.Groups["from_rank"].Value[0] : (char?)null,
                FromRank = match.Groups["from_rank"].Value[0],
                ToFile = match.Groups["to_file"].Value[0],
                ToRank = match.Groups["to_rank"].Value[0],
                Promotion = match.Groups["promotion"].Value.Length > 0 ? match.Groups["promotion"].Value[0] : (char?)null,
                Check = match.Groups["check"].Value.Length > 0 ? match.Groups["check"].Value[0] : (char?)null
            };
        }

        private PieceEnum GetPieceEnum(char pieceChar)
        {
            if (PieceMap.TryGetValue(pieceChar, out PieceEnum pieceEnum))
                return pieceEnum;
            else
                throw new ArgumentException($"Invalid piece character: {pieceChar}");
        }

        private int ConvertFileRankToSquare(char file, char rank)
        {
            return SquareHelper.GetNumBySquareName($"{file}{rank}");
        }

        private class ParsedSan
        {
            public char Piece { get; set; }
            public char FromFile { get; set; }
            public char FromRank { get; set; }
            public char ToFile { get; set; }
            public char ToRank { get; set; }
            public char? Promotion { get; set; }
            public char? Check { get; set; }
        }
    }
}
