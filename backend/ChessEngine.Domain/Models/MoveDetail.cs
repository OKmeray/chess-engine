using ChessEngine.Enums;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ChessEngine.Domain.Models
{
    public class MoveDetail
    {
        public PieceEnum Piece { get; set; }
        public PieceColor Color { get; set; }
        public int Square { get; set; }
        public int Move { get; set; }
        public PieceEnum Promotion { get; set; } = PieceEnum.NONE;
        public int Priority { get; set; } = 0;
    }
}
