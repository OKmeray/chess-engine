namespace ChessEngine.Domain.Models
{
    public class Bitboard
    {
        public ulong BitboardValue { get; set; }

        public Bitboard(ulong bitboardValue = 0)
        {
            BitboardValue = bitboardValue;
        }

        public void AddPiece(int squareNo)
        {
            BitboardValue |= 1UL << squareNo;
        }

        public override string ToString()
        {
            return Convert.ToString((long)BitboardValue, 2).PadLeft(64, '0');
        }

        public string Get8By8Board()
        {
            string bitString = Convert.ToString((long)BitboardValue, 2).PadLeft(64, '0');
            string formatted = "";
            for (int i = 0; i < 64; i += 8)
            {
                formatted += bitString.Substring(i, 8) + "\n";
            }
            return formatted;
        }

        public static Bitboard operator +(Bitboard a, Bitboard b)
        {
            return new Bitboard(a.BitboardValue | b.BitboardValue);
        }

        public static Bitboard operator |(Bitboard a, Bitboard b)
        {
            return new Bitboard(a.BitboardValue | b.BitboardValue);
        }

        public static Bitboard operator ^(Bitboard a, Bitboard b)
        {
            return new Bitboard(a.BitboardValue ^ b.BitboardValue);
        }

        public static Bitboard operator &(Bitboard a, Bitboard b)
        {
            return new Bitboard(a.BitboardValue & b.BitboardValue);
        }
    }
}