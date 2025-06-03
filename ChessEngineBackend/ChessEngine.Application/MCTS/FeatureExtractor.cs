using System;
using System.Collections.Generic;
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
        public const int TotalPlanes = HistoryLength * PlanesPerPosition + 4 + 1 + 1 + 1; // ==119

        public static float[,,] Extract(Position pos)
        {
            // 1) call your full 119‐plane extractor over just this one position
            //    (history of length 1, zero-padding for the other 7 positions)
            var full119 = ExtractFeatures(new List<Position> { pos }, repetitionCount: 1);
            //    full119 is float[8,8,119]

            // 2) allocate the 21‐plane result
            const int OUT_PLANES = 14 + 7;
            var F21 = new float[8, 8, OUT_PLANES];

            // 3) copy planes 0–13
            for (int r = 0; r < 8; r++)
                for (int c = 0; c < 8; c++)
                    for (int k = 0; k < 14; k++)
                        F21[r, c, k] = full119[r, c, k];

            // 4) copy planes 112–118 into F21[...,14..20]
            //    (112 = HistoryLength*PlanesPerPosition, i.e. 8*14)
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

            // --- 1) last 8 positions (or pad with zeros) ---
            for (int i = 0; i < HistoryLength; i++, p += PlanesPerPosition)
            {
                if (i < history.Count)
                {
                    var posPlanes = GetPositionPlanes(history[i], repetitionCount);
                    // copy [14,8,8] into F[..., p..p+13]
                    for (int k = 0; k < PlanesPerPosition; k++)
                        for (int r = 0; r < 8; r++)
                            for (int c = 0; c < 8; c++)
                                F[r, c, p + k] = posPlanes[k, r, c];
                }
                // else leave zeros
            }

            var cur = history[0];

            // --- 2) castling rights (4 planes) ---
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

            // --- 3) fullmove-zero plane (all zeros) ---
            p++;

            // --- 4) halfmove-count plane (one 8×8 bitboard) ---
            string bits = Convert.ToString(cur.HalfMoves, 2).PadLeft(64, '0');
            for (int i = 0; i < 64; i++)
            {
                int r = i / 8, c = i % 8;
                F[r, c, p] = bits[i] == '1' ? 1f : 0f;
            }
            p++;

            // --- 5) side‐to‐move plane ---
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

            //  a) 12 piece‐planes (we assume pos.Bitboards[0]=WhiteKing,...[11]=BlackPawn)
            for (int bi = 0; bi < 12; bi++, idx++)
            {
                ulong bb = pos.Bitboards[bi].BitboardValue;
                for (int sq = 0; sq < 64; sq++)
                {
                    int r = sq / 8, c = sq % 8;
                    M[idx, r, c] = ((bb >> sq) & 1UL) != 0 ? 1f : 0f;
                }
            }

            //  b) repetition planes
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
