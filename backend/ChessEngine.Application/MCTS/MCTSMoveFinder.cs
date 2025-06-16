using System.Diagnostics;
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
        private int _simulationsCount = 0;
        private const float C_PUCT = 8.0f;
        private const float C_PUCT_ROOT = 12.0f;

        public MCTSMoveFinder(string onnxModel)
        {
            if (onnxModel == String.Empty)
            {
                onnxModel = "4d51";
            }
            Console.WriteLine($"Model here: {onnxModel}");

            var sessionOptions = new SessionOptions();
            try // trying to access GPU
            {
                sessionOptions.AppendExecutionProvider_DML(0); 
            }
            catch
            {
                sessionOptions.AppendExecutionProvider_CPU();
            }
            sessionOptions.EnableMemoryPattern = false;
            sessionOptions.ExecutionMode = ExecutionMode.ORT_SEQUENTIAL;
            sessionOptions.GraphOptimizationLevel = GraphOptimizationLevel.ORT_ENABLE_ALL;
            sessionOptions.EnableCpuMemArena = false;

            string onnxModelPath = Path.Combine($"AIModels/{onnxModel}.onnx");

            _session = new InferenceSession(onnxModelPath, sessionOptions);
        }

        public (int evaluation, MoveDetail bestMove) GetBestMove(Position position, int timeLimitMs = 50)
        {
            var root = new Node(null, position, null, prior: 1f, treeDepth: 0);
            var sw = Stopwatch.StartNew();

            while (sw.ElapsedMilliseconds < timeLimitMs)
            {
                // Selection
                var leaf = root.SelectLeaf();

                if (leaf.State.IsGameOver())
                {
                    float terminalValue = leaf.State.DetermineOutcome() switch
                    {
                        PositionOutcome.WIN => 1f,
                        PositionOutcome.LOSS => -1f,
                        _ => 0f,
                    };
                    leaf.Backpropagate(terminalValue);
                    continue;
                }

                // Expansion
                var (policyBatch, valueBatch) = EvaluateNetworkBatch(new[] { leaf.State });
                float[] policy = policyBatch[0];
                float value = valueBatch[0];

                leaf.Expand(policy);
                leaf.Backpropagate(value);
               
                _simulationsCount++;
            }

            var bestNode = root.GetMostVisitedNode();
            Console.WriteLine($"Prior of best move: {bestNode.Prior}");
            Console.WriteLine($"Simulations: {_simulationsCount}");
            Console.WriteLine($"Selected Move UCT: {bestNode.GetPUCT():f3} | Visits: {bestNode.Visits}");
            _simulationsCount = 0;

            return (bestNode.Visits, bestNode.PrevMove);
        }

        private (float[][] policies, float[] values) EvaluateNetworkBatch(Position[] positions)
        {
            int batchSize = positions.Length;
            var inputData = new float[batchSize * 8 * 8 * 21];

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

            var policyTensor = results.First(x => x.Name == "policy").AsTensor<float>();
            var valueTensor = results.First(x => x.Name == "value").AsTensor<float>();

            float[][] policies = new float[batchSize][];
            for (int b = 0; b < batchSize; b++)
            {
                policies[b] = new float[policyTensor.Dimensions[1]];
                for (int i = 0; i < policies[b].Length; i++)
                    policies[b][i] = policyTensor[b, i];
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
            public bool IsExpanded { get; private set; }
            public int Visits { get; private set; }
            public float TotalValue { get; private set; }
            public float Prior { get; }
            public int TreeDepth { get; set; }

            public Node(Node parent, Position state, MoveDetail prevMove, float prior, int treeDepth)
            {
                Parent = parent;
                State = state;
                PrevMove = prevMove;
                Prior = prior;
                TreeDepth = treeDepth;
            }

            public Node SelectLeaf()
            {
                var current = this;
                while (current.IsExpanded && current.Children.Count > 0)
                {
                    current = current.SelectChild(limit: 1);
                }
                return current;
            }

            public void Expand(float[] policy)
            {
                var legalMoves = State.GetAllMoves();
                bool mirror = State.SideToMove == PieceColor.WHITE;
                foreach (var move in legalMoves)
                {
                    var next = State.Clone();
                    next.ApplyMove(move);
                    int moveIndex = MoveEncoder.Encode(move, mirror);
                    float prior = policy[moveIndex];
                    Children.Add(new Node(this, next, move, prior, TreeDepth + 1));
                }
                IsExpanded = true;
            }

            public void Backpropagate(float value)
            {
                Visits++;
                if (TreeDepth % 2 == 1)
                {
                    TotalValue -= value;
                }
                else
                {
                    TotalValue += value;
                }
                Parent?.Backpropagate(value);

            }

            public Node SelectChild(int limit = 1)
            {
                if (Children.Count == 0)
                    return null;

                Node best = null;
                float bestScore = float.NegativeInfinity;

                foreach (var c in Children)
                {
                    float s = c.GetPUCT();
                    if (s > bestScore)
                    {
                        bestScore = s;
                        best = c;
                    }
                }

                return best;
            }

            public float GetPUCT()
            {
                float parentN = Parent?.Visits ?? Visits;
                float cpuct = Parent == null ? C_PUCT_ROOT : C_PUCT;
                float q = Visits > 0 ? TotalValue / Visits : 0;
                float u = Prior * cpuct * (float)Math.Sqrt(parentN) / (1 + Visits);

                return q + u;
            }

            public Node GetMostVisitedNode()
            {
                return Children.OrderByDescending(c => c.Visits).FirstOrDefault();
            }
        }
    }
}

