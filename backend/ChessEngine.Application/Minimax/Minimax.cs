using System.Diagnostics;
using ChessEngine.Domain.Interfaces;
using ChessEngine.Domain.Models;
using ChessEngine.Enums;
using ChessEngine.Evaluation;

namespace ChessEngine.Application.Minimax
{
    public class MinimaxMoveFinder : IMoveFinder
    {
        /// <summary>
        /// Повертає найкращий хід, витрачаючи не більше timeLimitMs.
        /// </summary>
        public (int evaluation, MoveDetail bestMove) GetBestMove(
            Position position,
            int timeLimitMs = 50)                // ← 50 мс – типовий «один кадр» у UI
        {
            var sw = Stopwatch.StartNew();
            const int SAFETY_MARGIN = 3;          // залишимо кілька мс «запасу»

            MoveDetail bestMoveSoFar = null;
            int bestEvalSoFar = 0;

            // Ітеративне поглиблення: 1, 2, 3, …
            for (int depth = 1; ; depth++)
            {
                // якщо часу вже майже нема – припиняємо
                if (sw.ElapsedMilliseconds + SAFETY_MARGIN >= timeLimitMs)
                    break;

                bool completed = true;

                // Звичайний α-β
                (int eval, MoveDetail move) =
                    AlphaBeta(position,
                              alpha: int.MinValue,
                              beta: int.MaxValue,
                              depth: depth,
                              sw,
                              timeLimitMs,
                              ref completed,
                              isRoot: true);

                if (completed)
                {
                    bestMoveSoFar = move;
                    bestEvalSoFar = eval;
                }
                else
                {
                    // не встигли завершити цю глибину – більше не пробуємо
                    break;
                }
            }

            return (bestEvalSoFar, bestMoveSoFar);
        }

        // ────────── α-β з перевіркою тайм-ауту ──────────────────────────────
        private static (int, MoveDetail) AlphaBeta(
            Position pos,
            int alpha,
            int beta,
            int depth,
            Stopwatch sw,
            int timeLimitMs,
            ref bool completed,
            bool isRoot = false)
        {
            // Перериваємо рекурсію, якщо час вийшов
            if (sw.ElapsedMilliseconds >= timeLimitMs)
            {
                completed = false;
                return (0, null);
            }

            // 1. Фініш позиції / глибини
            if (depth == 0 || pos.IsGameOver())
                return (Evaluator.EvaluatePosition(pos), null);

            var moves = pos.GetAllMoves();
            // Невеличке сортування «найкращі спочатку» – прискорює відсікання
            moves.Sort((a, b) =>
            {
                var pa = pos.Clone(); pa.ApplyMove(a);
                var pb = pos.Clone(); pb.ApplyMove(b);
                return Evaluator.EvaluatePosition(pb)
                     .CompareTo(Evaluator.EvaluatePosition(pa));
            });

            MoveDetail bestMove = null;

            if (pos.SideToMove == PieceColor.WHITE)
            {
                int maxEval = int.MinValue;
                foreach (var mv in moves)
                {
                    var next = pos.Clone();
                    next.ApplyMove(mv);

                    (int eval, _) = AlphaBeta(
                        next, alpha, beta, depth - 1,
                        sw, timeLimitMs, ref completed);

                    if (!completed) return (0, null);   // далі вже не йдемо

                    if (eval > maxEval)
                    {
                        maxEval = eval;
                        if (isRoot) bestMove = mv;
                    }
                    alpha = Math.Max(alpha, eval);
                    if (alpha >= beta) break;          // β-відсікання
                }
                return (maxEval, bestMove);
            }
            else
            {
                int minEval = int.MaxValue;
                foreach (var mv in moves)
                {
                    var next = pos.Clone();
                    next.ApplyMove(mv);

                    (int eval, _) = AlphaBeta(
                        next, alpha, beta, depth - 1,
                        sw, timeLimitMs, ref completed);

                    if (!completed) return (0, null);

                    if (eval < minEval)
                    {
                        minEval = eval;
                        if (isRoot) bestMove = mv;
                    }
                    beta = Math.Min(beta, eval);
                    if (beta <= alpha) break;          // α-відсікання
                }
                return (minEval, bestMove);
            }
        }
    }
}
