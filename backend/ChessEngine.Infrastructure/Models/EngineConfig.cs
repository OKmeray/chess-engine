using MongoDB.Bson.Serialization.Attributes;
using MongoDB.Bson;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ChessEngine.Persistance.Models
{
    public class EngineConfig
    {
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string? Id { get; set; }

        public string HEXID { get; set; }
        public int ID { get; set; }
        public int Filters { get; set; }
        public int ResBlocks { get; set; }
        public int BatchSize { get; set; }
        public double LabelSmoothing { get; set; }
        public double LR { get; set; }
        public bool SEBlocks { get; set; }
        public string LossWeights { get; set; }
        public double DropoutRate { get; set; }
        public double L2 { get; set; }
        public int Epochs { get; set; }
        public int SimplePositions { get; set; }
        public int PuzzlePositions { get; set; }
        public double test_policy { get; set; }
        public double test_value { get; set; }
        public double mae { get; set; }
        public double test_top1 { get; set; }
        public double test_top3 { get; set; }
        public double test_top5 { get; set; }
    }
}
