using ChessEngine.Domain.Models;
using ChessEngine.Enums;

/// <summary>
/// Fully describes one of the 76 “move‐types”:
///  • kind = “slide” / “knight” / “promo”
///  • dir  = N, NE, E, …
///  • steps = for slides
///  • promoPiece = for promotions
/// </summary>
public sealed record MoveType
(
    string Kind,        // "slide", "knight" or "promo"
    string Dir,         // "N", "NE", ...
    int Steps,       // 1..7 for slide, unused for knight/promo
    char PromoPiece   // 'Q','R','B','N' for promo, else '\0'
)
{
    public static MoveType FromMoveDetail(MoveDetail m)
    {
        // extract from‐square index, to‐square index
        int from = m.Square, to = m.Move;
        int fx = from % 8, fy = from / 8;
        int tx = to % 8, ty = to / 8;
        int dx = tx - fx, dy = ty - fy;

        // promotion
        if (m.Promotion != PieceEnum.NONE)
        {
            // dx,dy must be one step
            var dir =
              (dy > 0 ? "N" : "") +
              (dy < 0 ? "S" : "") +
              (dx > 0 ? "E" : "") +
              (dx < 0 ? "W" : "");

            return new MoveType(
                Kind: "promo",
                Dir: dir,
                Steps: 0,
                PromoPiece: m.Promotion switch
                {
                    PieceEnum.QUEEN => 'Q',
                    PieceEnum.ROOK => 'R',
                    PieceEnum.BISHOP => 'B',
                    PieceEnum.KNIGHT => 'N',
                    _ => throw new InvalidOperationException("Unexpected promotion")
                }
            );
        }

        // knight
        if (Math.Abs(dx) == 2 && Math.Abs(dy) == 1 ||
            Math.Abs(dx) == 1 && Math.Abs(dy) == 2)
        {
            // order the two moves so they match your Python logic:
            string big = Math.Abs(dx) == 2
                         ? (dx > 0 ? "E" : "W")
                         : (dy > 0 ? "N" : "S");
            string small = big == ("N") || big == ("S")
                         ? (dx > 0 ? "E" : "W")
                         : (dy > 0 ? "N" : "S");

            return new MoveType(
                Kind: "knight",
                Dir: big,
                Steps: 0,
                PromoPiece: small[0]
            );
        }

        // 4) slide
        var sdir =
          (dy > 0 ? "N" : "") +
          (dy < 0 ? "S" : "") +
          (dx > 0 ? "E" : "") +
          (dx < 0 ? "W" : "");
        int steps = Math.Max(Math.Abs(dx), Math.Abs(dy));
        return new MoveType("slide", sdir, steps, '\0');
    }
}
