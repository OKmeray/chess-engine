using ChessEngine.Domain.Models;
using ChessEngine.Helpers;
using static System.Runtime.InteropServices.JavaScript.JSType;

public static class MoveEncoder
{
    public const int TypesCount = 56 + 8 + 12;      // 76
    public const int PolicySize = 64 * TypesCount;  // 4864

    private static readonly Dictionary<MoveType, int> _typeToIndex;
    private static readonly MoveType[] _indexToType;

    static MoveEncoder()
    {
        _typeToIndex = new();
        int idx = 0;

        // sliding
        foreach (var dir in new[] { "N", "NE", "E", "SE", "S", "SW", "W", "NW" })
            for (int s = 1; s <= 7; s++)
                _typeToIndex[new MoveType("slide", dir, s, '\0')] = idx++;

        // knights
        foreach (var (d2, d1) in new[]{
            ("N","E"),("N","W"),("S","E"),("S","W"),
            ("E","N"),("E","S"),("W","N"),("W","S")
        })
            _typeToIndex[new MoveType("knight", d2, 0, d1[0])] = idx++;

        // promotions
        foreach (var dir in new[] { "NW", "N", "NE" })
            foreach (var promo in new[] { 'Q', 'R', 'B', 'N' })
                _typeToIndex[new MoveType("promo", dir, 0, promo)] = idx++;

        // build reverse lookup if you need it
        _indexToType = new MoveType[PolicySize];
        foreach (var kv in _typeToIndex)
            _indexToType[kv.Value] = kv.Key;
    }

    public static int Encode(MoveDetail m, bool mirror = false)
    {
        if (mirror)
            m = SquareHelper.MirrorMove(m);

        int fromSq = m.Square;                   // 0..63
        var kind = MoveType.FromMoveDetail(m);
        int t = _typeToIndex[kind];             // 0..75
        return fromSq * TypesCount + t;
    }

    //public static float[] MaskAndNormalize(IEnumerable<MoveDetail> legal, float[] rawPolicy)
    //{
    //    var mask = new float[PolicySize];
    //    float sum = 0;

    //    foreach (var m in legal)
    //    {
    //        int i = Encode(m);
    //        mask[i] = rawPolicy[i];
    //        sum += mask[i];
    //    }

    //    if (sum > 0f)
    //        for (int i = 0; i < mask.Length; i++)
    //            mask[i] /= sum;

    //    return mask;
    //}
}
