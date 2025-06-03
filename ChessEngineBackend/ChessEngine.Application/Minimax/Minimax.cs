using System;
using System.Collections.Generic;
using ChessEngine.Enums;
using ChessEngine.Domain.Models;
using ChessEngine.Evaluation;
using System.Diagnostics;
using ChessEngine.FEN;
using ChessEngine.Domain.Interfaces;


namespace ChessEngine.Application.Minimax
{
    public class MinimaxMoveFinder : IMoveFinder
    {
        public (int evaluation, MoveDetail bestMove) GetBestMove(Position position, int depth)
        {
            int alpha, beta;
            if (position.SideToMove == PieceColor.WHITE)
            {
                alpha = int.MaxValue;
                beta = int.MinValue;
            }
            else
            {
                alpha = int.MinValue;
                beta = int.MaxValue;
            }
            Console.WriteLine(depth + " " + alpha + " " + beta);
            return RunMinimax(position, alpha: alpha, beta: beta, depth: depth, isRoot: true);
        }

        private static (int, MoveDetail) RunMinimax(Position position, int alpha, int beta, int depth = 5, bool isRoot = true)  //  int branchesToCheck = 10
        {
            if (position.IsGameOver())
            {
                PositionOutcome outcome = position.DetermineOutcome();
                if (outcome == PositionOutcome.WIN)
                {
                    if (position.SideToMove == PieceColor.WHITE)
                    {
                        return (int.MaxValue, null);
                    }
                    else
                    {
                        return (int.MinValue, null);
                    }
                }
                else if (outcome == PositionOutcome.LOSS)
                {
                    if (position.SideToMove == PieceColor.WHITE)
                    {
                        return (int.MinValue, null);
                    }
                    else
                    {
                        return (int.MaxValue, null);
                    }
                }
                else
                {
                    return (0, null); // Draw
                }
            }

            if (depth == 0)
            {
                int evaluation = Evaluator.EvaluatePosition(position);
                return (evaluation, null);
            }

            MoveDetail bestMove = null;

            if (position.SideToMove == PieceColor.WHITE)
            {
                int maxEval = int.MinValue;

                // TODO: test here allMovesSorted
                List<MoveDetail> allMoves = position.GetAllMoves();
                allMoves.Sort((a, b) =>
                {
                    Position posA = position.Clone();
                    Position posB = position.Clone();
                    posA.ApplyMove(a);
                    posB.ApplyMove(b);
                    return Evaluator.EvaluatePosition(posA).CompareTo(Evaluator.EvaluatePosition(posB));
                });
                // allMoves = allMoves.GetRange(0, Math.Min(branchesToCheck, allMoves.Count));

                foreach (var move in allMoves)
                {
                    Position positionCopy = position.Clone();
                    positionCopy.ApplyMove(move);

                    (int eval, _) = RunMinimax(positionCopy, alpha, beta, depth - 1, isRoot: false); //branchesToCheck - 1,
                    if (eval > maxEval)
                    {
                        maxEval = eval;
                        if (isRoot)
                        {
                            bestMove = move;
                        }
                    }

                    alpha = Math.Max(alpha, eval);
                    if (beta <= alpha)
                    {
                        break;
                    }
                }
                return (maxEval, bestMove);
            }
            else // Black's turn
            {
                int minEval = int.MaxValue;

                // TODO: the same 
                List<MoveDetail> allMoves = position.GetAllMoves();
                allMoves.Sort((a, b) =>
                {
                    Position posA = position.Clone();
                    Position posB = position.Clone();
                    posA.ApplyMove(a);
                    posB.ApplyMove(b);
                    return Evaluator.EvaluatePosition(posA).CompareTo(Evaluator.EvaluatePosition(posB));
                });
                // allMoves = allMoves.GetRange(0, Math.Min(branchesToCheck, allMoves.Count));

                foreach (var move in allMoves)
                {
                    Position positionCopy = position.Clone();
                    positionCopy.ApplyMove(move);

                    (int eval, _) = RunMinimax(positionCopy, alpha, beta, depth - 1, isRoot: false); //  branchesToCheck - 1,
                    if (eval < minEval)
                    {
                        minEval = eval;
                        if (isRoot)
                        {
                            bestMove = move;
                        }
                    }

                    beta = Math.Min(beta, eval);
                    if (beta <= alpha)
                    {
                        break;
                    }
                }
                return (minEval, bestMove);
            }
        }

        //public static int CountMoves(Position position, int depth = 5, bool isRoot = false)
        //{
        //    if (depth == 0)
        //    {
        //        return 1;
        //    }

        //    int count = 0;
        //    List<MoveDetail> moves = position.GetAllSeparateMoves();
        //    foreach (var move in moves)
        //    {
        //        Position copyPosition = position.Clone();
        //        if (isRoot)
        //        {
        //            Console.Write($"{move.Piece} {move.Color} {ConvertToAlgebraic(move.Square, move.Move)}   ");
        //        }
        //        copyPosition.ApplyMove(move);
        //        int copyPositionCount = CountMoves(copyPosition, depth - 1);
        //        count += copyPositionCount;
        //        if (isRoot)
        //        {
        //            Console.WriteLine(copyPositionCount);
        //        }
        //    }

        //    return count;
        //}
    }
}
