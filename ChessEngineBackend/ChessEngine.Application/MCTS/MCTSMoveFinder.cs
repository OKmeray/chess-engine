using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using ChessEngine.Application.Fen;
using ChessEngine.Domain.Interfaces;
using ChessEngine.Domain.Models;
using ChessEngine.Enums;
using Microsoft.ML.OnnxRuntime;
using Microsoft.ML.OnnxRuntime.Tensors;

namespace ChessEngine.Application.MCTS
{
    public class MCTSMoveFinder : IMoveFinder, IDisposable
    {
        private readonly InferenceSession _session;
        private readonly int _timeLimitMs;
        private int _simulationsCount = 0;
        private const float C_PUCT = 2.5f;
        // Add this method to your class

        public MCTSMoveFinder(string onnxModelPath, int timeLimitSeconds = 2)
        {
            var sessionOptions = new SessionOptions();

            // Set the GPU device ID (0 for primary GPU)
            int gpuDeviceId = 0;

            // Create DirectML execution provider
            sessionOptions.AppendExecutionProvider_DML(gpuDeviceId);

            // Important performance optimizations
            sessionOptions.EnableMemoryPattern = false;
            sessionOptions.ExecutionMode = ExecutionMode.ORT_SEQUENTIAL;

            // For better throughput with batch processing
            sessionOptions.GraphOptimizationLevel = GraphOptimizationLevel.ORT_ENABLE_ALL;

            // Disable memory arena (can help with memory fragmentation)
            sessionOptions.EnableCpuMemArena = false;

            _session = new InferenceSession(onnxModelPath, sessionOptions);
            _timeLimitMs = timeLimitSeconds * 1000;
        }

        public (int evaluation, MoveDetail bestMove) GetBestMove(Position position, int unused = 0)
        {
            var root = new Node(null, position, null, prior: 1f);
            var sw = Stopwatch.StartNew();
            int batchSize = 32; // Experiment with this value
            var batchPositions = new List<Position>();
            var batchNodes = new List<Node>();

            batchNodes.Clear();

            while (sw.ElapsedMilliseconds < _timeLimitMs)
            {
                // Selection phase (can also be parallelized)
                var rootNode = root;
                while (rootNode.Children.Count > 0)
                    rootNode = rootNode.SelectChild();

                if (rootNode.State.IsGameOver())
                {
                    float leafValue = rootNode.State.DetermineOutcome() == PositionOutcome.WIN ? 1f : 0f;
                    rootNode.Backpropagate(leafValue);
                    continue;
                }

                bool isBlack = false;
                if (rootNode.State.SideToMove == PieceColor.BLACK)
                {
                    isBlack = true;
                }

                batchNodes.Add(rootNode);
                batchPositions.Add(rootNode.State);

                // When batch is full or time is running out
                if (batchNodes.Count >= batchSize ||
                    (_timeLimitMs - sw.ElapsedMilliseconds) < 50)
                {
                    _simulationsCount += batchNodes.Count;
                    var (policies, values) = EvaluateNetworkBatch(batchPositions.ToArray());

                    for (int i = 0; i < batchNodes.Count; i++)
                    {
                        var node = batchNodes[i];
                        var policy = policies[i];
                        var value = values[i];
                        float leafValue = value > 0.34f ? 1f : 0f;

                        // Expansion
                        var legal = node.State.GetAllMoves();
                        bool mirror = node.State.SideToMove == PieceColor.BLACK;
                        foreach (var move in legal)
                        {
                            var next = node.State.Clone();
                            next.ApplyMove(move);
                            
                            int moveIndex = MoveEncoder.Encode(move, mirror);
                            int falseMoveIndex = MoveEncoder.Encode(move, !mirror);
                            float prior = policy[moveIndex];
                            float falsePrior = policy[falseMoveIndex];
                            Console.WriteLine($"{move.Square} {mirror} {moveIndex} {falseMoveIndex}");

                            Console.WriteLine($"{move.Square}, {move.Move}, {moveIndex} -> {prior}; {falseMoveIndex} -> {falsePrior}");
                            node.Children.Add(new Node(node, next, move, prior));
                        }
                        // Backpropagate
                        node.Backpropagate(leafValue);
                    }

                    batchNodes.Clear();
                    batchPositions.Clear();
                }
            }

            // Choose best move
            var best = root.Children.OrderByDescending(c => c.Visits).First();

            Console.WriteLine(_simulationsCount);
            _simulationsCount = 0;
            Console.WriteLine($"Policy for selected move {best.Prior}");

            return (best.Visits, best.PrevMove);
        }

        private (float[][] policies, float[] values) EvaluateNetworkBatch(Position[] positions)
        {
            int batchSize = positions.Length;
            var inputData = new float[batchSize * 8 * 8 * 21];

            // Prepare batch input
            for (int b = 0; b < batchSize; b++)
            {
                var planes = FeatureExtractor.Extract(positions[b]);
                for (int i = 0; i < 8; i++)
                    for (int j = 0; j < 8; j++)
                        for (int c = 0; c < 21; c++)
                            inputData[b * 8 * 8 * 21 + i * 8 * 21 + j * 21 + c] = planes[i, j, c];
            }

            var tensor = new DenseTensor<float>(inputData, new[] { batchSize, 8, 8, 21 });

            using var results = _session.Run(new[] {
                NamedOnnxValue.CreateFromTensor(_session.InputMetadata.Keys.First(), tensor)
            });

            // Process batch output
            var policyTensor = results.First(x => x.Name == "policy").AsTensor<float>();
            var valueTensor = results.First(x => x.Name == "value").AsTensor<float>();

            float[][] policies = new float[batchSize][];
            for (int b = 0; b < batchSize; b++)
            {
                policies[b] = new float[policyTensor.Dimensions[1]];
                float maxPolicy = 0;
                float maxIndex = -1;
                for (int i = 0; i < policies[b].Length; i++)
                {
                    policies[b][i] = policyTensor[b, i];
                    maxPolicy = maxPolicy > policies[b][i] ? maxPolicy : policies[b][i];
                    maxIndex = i;
                }
                Console.WriteLine($"MaxPolicy {maxPolicy}; maxIndex {maxIndex}");
            }

            float[] values = new float[batchSize];
            for (int b = 0; b < batchSize; b++)
                values[b] = valueTensor[b];

            return (policies, values);
        }

        public void Dispose() => _session.Dispose();

        class Node
        {
            public Node Parent { get; }
            public Position State { get; }
            public MoveDetail PrevMove { get; }
            public List<Node> Children { get; } = new();
            public int Visits { get; private set; }
            public float TotalValue { get; private set; }
            public float Prior { get; set; }

            public Node(Node parent, Position state, MoveDetail prevMove, float prior)
            {
                Parent = parent;
                State = state;
                PrevMove = prevMove;
                Prior = prior;
            }

            public Node SelectChild()
            {
                float parentN = Parent?.Visits ?? Visits;
                return Children
                  .OrderByDescending(c =>
                  {
                      float q = c.Visits > 0 ? c.TotalValue / c.Visits : 0f;
                      float u = c.Prior * C_PUCT * (float)Math.Sqrt(parentN) / (1 + c.Visits); // 
                      return q + u;
                  })
                  .First();
            }

            public void Backpropagate(float leafValue)
            {
                Visits++;
                TotalValue += leafValue;
                Parent?.Backpropagate(leafValue);
            }
        }
    }
}


//using System;
//using System.Collections.Concurrent;
//using System.Collections.Generic;
//using System.Diagnostics;
//using System.Linq;
//using ChessEngine.Domain.Interfaces;
//using ChessEngine.Domain.Models;
//using ChessEngine.Enums;
//using Microsoft.ML.OnnxRuntime;
//using Microsoft.ML.OnnxRuntime.Tensors;

//namespace ChessEngine.Application.MCTS
//{
//    public class MCTSMoveFinder : IMoveFinder, IDisposable
//    {
//        private readonly InferenceSession _session;
//        private readonly ConcurrentDictionary<ulong, (float[] P, float V)> _cache = new();
//        private readonly int _timeLimitMs;
//        private int _simulationsCount = 0;
//        private const int BATCH_SIZE = 64;
//        private const float C_PUCT = 1.2f;
//        private const float C_PUCT_ROOT = 2.00f;
//        private const float FPU_REDUCTION = 0.49f;
//        private const float FPU_REDUCTION_ROOT = 1.00f;
//        private const float DIRICHLET_EPSILON = 0.10f;
//        private const float DIRICHLET_ALPHA = 0.12f;
//        // Add this method to your class

//        public MCTSMoveFinder(string onnxModelPath, int timeLimitSeconds = 2)
//        {
//            var sessionOptions = new SessionOptions();

//            // Set the GPU device ID (0 for primary GPU)
//            int gpuDeviceId = 0;

//            // Create DirectML execution provider
//            sessionOptions.AppendExecutionProvider_DML(gpuDeviceId);

//            // Important performance optimizations
//            sessionOptions.EnableMemoryPattern = false;
//            sessionOptions.ExecutionMode = ExecutionMode.ORT_SEQUENTIAL;

//            // For better throughput with batch processing
//            sessionOptions.GraphOptimizationLevel = GraphOptimizationLevel.ORT_ENABLE_ALL;

//            // Disable memory arena (can help with memory fragmentation)
//            sessionOptions.EnableCpuMemArena = false;

//            _session = new InferenceSession(onnxModelPath, sessionOptions);
//            _timeLimitMs = timeLimitSeconds * 1000;
//        }

//        public (int evaluation, MoveDetail bestMove) GetBestMove(Position rootPos, int _ = 0)
//        {
//            var root = new Node(null, rootPos, null, prior: 1f, isRoot: true);
//            AddDirichletNoise(root);

//            var sw = Stopwatch.StartNew();
//            var batchPositions = new List<Position>(BATCH_SIZE);
//            var batchNodes = new List<Node>(BATCH_SIZE);

//            while (sw.ElapsedMilliseconds < _timeLimitMs)
//            {
//                // Selection phase (can also be parallelized)
//                var node = root;
//                while (node.Expanded)
//                    node = node.SelectChild();

//                if (node.State.IsGameOver())
//                {
//                    float leaf = ResultToValue(node.State.DetermineOutcome());
//                    node.Backpropagate(leaf);
//                    continue;
//                }

//                batchNodes.Add(node);
//                batchPositions.Add(node.State);

//                // When batch is full or time is running out
//                if (batchNodes.Count < BATCH_SIZE && (_timeLimitMs - sw.ElapsedMilliseconds) > 50)
//                {
//                    continue;
//                }

//                EvaluateAndExpand(batchNodes, batchPositions);
//                _simulationsCount += batchNodes.Count;
//                batchNodes.Clear(); 
//                batchPositions.Clear();

//            }

//            // Choose best move
//            var best = root.Children.OrderByDescending(c => c.Visits).First();

//            Console.WriteLine(_simulationsCount);
//            _simulationsCount = 0;
//            Console.WriteLine($"Policy for selected move {best.Prior}");

//            return (best.Visits, best.Move);
//        }

//        private void EvaluateAndExpand(IReadOnlyList<Node> nodes,
//                                      IReadOnlyList<Position> poss)
//        {
//            // 1) NN inference (batched, с кешем)
//            var policies = new float[nodes.Count][];
//            var values = new float[nodes.Count];

//            // prepare tensor
//            var tensor = new DenseTensor<float>(new Memory<float>(
//                               new float[nodes.Count * 8 * 8 * 21]), // NHWC
//                               new[] { nodes.Count, 8, 8, 21 });

//            int cursor = 0;
//            var inferIndices = new List<int>(nodes.Count); // which need NN

//            for (int i = 0; i < nodes.Count; ++i)
//            {
//                ulong key = poss[i].Zobrist;               // Position must expose Zobrist hash
//                if (_cache.TryGetValue(key, out var hit))
//                {
//                    policies[i] = hit.P;
//                    values[i] = hit.V;
//                }
//                else
//                {
//                    inferIndices.Add(i);

//                    // --- flatten 8×8×21 у такий самий порядок, як у старій версії -------------
//                    var planes = FeatureExtractor.Extract(poss[i]);     // float[8,8,21]
//                    const int S = 8 * 8 * 21;                           // 1344
//                    var dest = tensor.Buffer.Span.Slice(cursor, S);     // куди пишемо

//                    int d = 0;                                          // позиція в dest
//                    for (int r = 0; r < 8; ++r)                         // row    (i  у старому коді)
//                        for (int c = 0; c < 8; ++c)                     // column (j)
//                            for (int ch = 0; ch < 21; ++ch)             // channel(c)
//                                dest[d++] = planes[r, c, ch];

//                    cursor += S;
//                }
//            }

//            if (inferIndices.Count > 0)
//            {
//                // shrink if needed
//                var subTensor = tensor;
//                if (inferIndices.Count != nodes.Count)
//                    subTensor = new DenseTensor<float>(
//                        tensor.Buffer,
//                        new[] { inferIndices.Count, 8, 8, 21 });

//                using var results = _session.Run(new[]
//                {
//                    NamedOnnxValue.CreateFromTensor(
//                        _session.InputMetadata.Keys.First(), subTensor)
//                });

//                var pT = results.First(r => r.Name == "policy").AsTensor<float>();
//                var vT = results.First(r => r.Name == "value").AsTensor<float>();

//                for (int k = 0; k < inferIndices.Count; ++k)
//                {
//                    int i = inferIndices[k];

//                    int len = (int)pT.Dimensions[1];
//                    var buf = new float[len];
//                    for (int j = 0; j < len; ++j)
//                        buf[j] = pT[k, j];
//                    policies[i] = buf;

//                    values[i] = vT[k];
//                    _cache.TryAdd(poss[i].Zobrist, (policies[i], values[i]));
//                }
//            }

//            // 2) Expansion + backup
//            for (int i = 0; i < nodes.Count; ++i)
//            {
//                var n = nodes[i];
//                var P = policies[i];
//                var V = values[i];

//                n.Expand(P);                          // create children
//                n.Backpropagate(-V);                  // raw value, invert for parent
//            }
//        }

//        private static float ResultToValue(PositionOutcome o) =>
//            o switch
//            {
//                PositionOutcome.WIN => 1f,
//                PositionOutcome.LOSS => -1f,
//                _ => 0f
//            };

//        private static readonly Random _rng = new();
//        private static void AddDirichletNoise(Node root)
//        {
//            var α = DIRICHLET_ALPHA;
//            var ε = DIRICHLET_EPSILON;

//            int k = root.Children.Count;
//            if (k == 0) return;

//            // simple γ-sample
//            var dir = new double[k];
//            double sum = 0;
//            for (int i = 0; i < k; ++i)
//            {
//                // Γ(α,1) via Marsaglia & Tsang
//                double d, c;
//                if (α < 1)
//                {
//                    d = α + (1 / 3.0) - 1;
//                    c = 1 / Math.Sqrt(9 * d);
//                }
//                else
//                {
//                    d = α - 1 / 3.0;
//                    c = 1 / Math.Sqrt(9 * d);
//                }
//                double x, v, u;
//                do
//                {
//                    do { x = SampleNormal(); v = 1 + c * x; } while (v <= 0);
//                    v = v * v * v;
//                    u = _rng.NextDouble();
//                } while (u > 1 - 0.331 * Math.Pow(x, 4) &&
//                         Math.Log(u) > 0.5 * x * x + d * (1 - v + Math.Log(v)));
//                dir[i] = d * v;
//                sum += dir[i];
//            }
//            for (int i = 0; i < k; ++i)
//            {
//                float η = (float)(dir[i] / sum);
//                root.Children[i].Prior = (1 - ε) * root.Children[i].Prior + ε * η;
//            }

//            static double SampleNormal()
//            {
//                // Box-Muller
//                double u1 = 1 - _rng.NextDouble();
//                double u2 = 1 - _rng.NextDouble();
//                return Math.Sqrt(-2 * Math.Log(u1)) * Math.Cos(2 * Math.PI * u2);
//            }
//        }

//        public void Dispose() => _session.Dispose();

//        class Node
//        {
//            public readonly Node Parent;
//            public readonly Position State;
//            public readonly MoveDetail Move;
//            public readonly bool IsRoot;

//            public readonly List<Node> Children = new();
//            public int Visits;
//            public float ValueSum;
//            public float Prior; // P(s,a)

//            public bool Expanded => Children.Count > 0;
//            public float Q => Visits == 0 ? 0 : ValueSum / Visits;

//            public Node(Node parent, Position state, MoveDetail move, float prior, bool isRoot = false)
//            {
//                Parent = parent;
//                State = state;
//                Move = move;
//                Prior = prior;
//                IsRoot = isRoot;
//            }

//            public void Expand(float[] policy)
//            {
//                if (Children.Count > 0)
//                    return;
//                var legal = State.GetAllMoves();
//                Console.WriteLine($"Legal moves: {legal.Count}");
//                foreach (var m in legal)
//                {
//                    var next = State.Clone();
//                    next.ApplyMove(m);
//                    int idx = MoveEncoder.Encode(m);
//                    Children.Add(new Node(this, next, m, policy[idx]));
//                }
//            }

//            public Node SelectChild()
//            {
//                float sqrtParent = (float)Math.Sqrt(Math.Max(1, Visits));
//                float parentQ = Q;

//                float bestScore = float.NegativeInfinity;
//                Node bestNode = null!;

//                foreach (var c in Children)
//                {
//                    float q = c.Visits > 0
//                              ? c.Q
//                              : parentQ - (IsRoot ? FPU_REDUCTION_ROOT : FPU_REDUCTION);

//                    float cpuct = IsRoot ? C_PUCT_ROOT : C_PUCT;
//                    float u = cpuct * c.Prior * sqrtParent / (1 + c.Visits);
//                    float score = q + u;

//                    if (score > bestScore)
//                    {
//                        bestScore = score;
//                        bestNode = c;
//                    }
//                }
//                return bestNode!;
//            }

//            public void Backpropagate(float leafValue)
//            {
//                Visits++;
//                ValueSum += leafValue;
//                Parent?.Backpropagate(-leafValue); // TODO: test this
//            }
//        }
//    }
//}