using ChessEngine.Domain.Models;
using ChessEngine.Enums;

namespace ChessEngine.Application.MCTS
{
    /// <summary>
    /// Build the 8×8×119 input tensor:
    ///   • last 8 positions × 14 planes  (12 piece + 2 repetition)  = 112  
    ///   • 4 castling planes                = 116  
    ///   • 1 fullmove-zero plane            = 117  
    ///   • 1 halfmove-count plane           = 118  
    ///   • 1 side-to-move plane             = 119  
    /// </summary>
    public static class FeatureExtractor
    {
        public const int HistoryLength = 8;
        public const int PlanesPerPosition = 14;
        public const int TotalPlanes = HistoryLength * PlanesPerPosition + 4 + 1 + 1 + 1; // 119

        public static float[,,] Extract(Position pos)
        {
            var full119 = ExtractFeatures(new List<Position> { pos }, repetitionCount: 1);  // float[8,8,119]

            var mirrored = new float[8, 8, TotalPlanes];
            if (pos.SideToMove == PieceColor.BLACK)
            {
                // planes: 0 -> 6; 1 -> 7; 2 -> 8
                // horizontals mirror i = 7 - i
                for (int i = 0; i < 8; i++)
                {
                    for (int j = 0; j < 8; j++)
                    {
                        for (int k = 0; k < TotalPlanes; k++)
                        {
                            int newk = k;
                            if (k < 12)
                            {
                                newk = (k + 6) % 12;
                            }
                            mirrored[7 - i, j, newk] = full119[i, j, k];
                        }
                    }
                }
                full119 = mirrored;
            }

            // allocate the 21‐plane result
            const int OUT_PLANES = 14 + 7;
            var F21 = new float[8, 8, OUT_PLANES];

            // copy planes 0–13
            for (int r = 0; r < 8; r++)
                for (int c = 0; c < 8; c++)
                    for (int k = 0; k < 14; k++)
                        F21[r, c, k] = full119[r, c, k];

            // copy planes 112–118
            int startOfLast7 = HistoryLength * PlanesPerPosition; // 112
            for (int r = 0; r < 8; r++)
                for (int c = 0; c < 8; c++)
                    for (int k = 0; k < 7; k++)
                        F21[r, c, 14 + k] = full119[r, c, startOfLast7 + k];

            return F21;
        }

        public static float[,,] ExtractFeatures(
            List<Position> history,
            int repetitionCount = 1
        )
        {
            // [row, col, plane]
            var F = new float[8, 8, TotalPlanes];
            int p = 0;

            for (int i = 0; i < HistoryLength; i++, p += PlanesPerPosition)
            {
                if (i < history.Count)
                {
                    var posPlanes = GetPositionPlanes(history[i], repetitionCount);
                    for (int k = 0; k < PlanesPerPosition; k++)
                        for (int r = 0; r < 8; r++)
                            for (int c = 0; c < 8; c++)
                                F[r, c, p + k] = posPlanes[k, r, c];
                }
                // else leave zeros
            }

            var cur = history[0];

            // castling rights 
            foreach (CastleEnum ct in new[] {
                CastleEnum.WhiteShortCastle,
                CastleEnum.WhiteLongCastle,
                CastleEnum.BlackShortCastle,
                CastleEnum.BlackLongCastle })
            {
                float v = cur.CastlingRights[ct] ? 1f : 0f;
                for (int r = 0; r < 8; r++)
                    for (int c = 0; c < 8; c++)
                        F[r, c, p] = v;
                p++;
            }

            // fullmove-zero plane (all zeros)
            p++;

            // halfmove-count plane (one 8×8 bitboard)
            string bits = Convert.ToString(cur.HalfMoves, 2).PadLeft(64, '0');
            for (int i = 0; i < 64; i++)
            {
                int r = i / 8, c = i % 8;
                F[r, c, p] = bits[i] == '1' ? 1f : 0f;
            }
            p++;

            // side‐to‐move plane 
            float stm = cur.SideToMove == PieceColor.WHITE ? 1f : 0f;
            for (int r = 0; r < 8; r++)
                for (int c = 0; c < 8; c++)
                    F[r, c, p] = stm;

            return F;
        }

        /// <summary>
        /// Returns a [14,8,8] block for one position:
        ///   • 12 piece-presence bitboards  
        ///   •  2 repetition planes  
        /// </summary>
        private static float[,,] GetPositionPlanes(Position pos, int rep)
        {
            var M = new float[PlanesPerPosition, 8, 8];
            int idx = 0;
            
            // 11 -> WhiteKing, ..., 0 -> BlackPawn
            for (int bi = 11; bi >= 0; bi--, idx++)
            {
                ulong bb = pos.Bitboards[bi].BitboardValue;
                for (int sq = 0; sq < 64; sq++)
                {
                    int r = sq / 8, c = sq % 8;
                    M[idx, r, c] = ((bb >> sq) & 1UL) != 0 ? 1f : 0f;
                }
            }

            // repetition planes
            float r1 = rep >= 2 ? 1f : 0f;
            float r2 = rep >= 3 ? 1f : 0f;
            for (int r = 0; r < 8; r++)
                for (int c = 0; c < 8; c++)
                {
                    M[12, r, c] = r1;
                    M[13, r, c] = r2;
                }

            return M;
        }
    }
}
