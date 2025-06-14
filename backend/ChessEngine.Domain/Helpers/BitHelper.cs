namespace ChessEngine.Helpers
{
    public static class BitHelper
    {
        /// <summary>
        /// Counts the number of set bits (1s) in the given ulong.
        /// </summary>
        public static int PopCount(ulong value)
        {
            int count = 0;
            while (value != 0)
            {
                count += (int)(value & 1);
                value >>= 1;
            }
            return count;
        }

        /// <summary>
        /// Finds the index of the least significant set bit (0-based).
        /// Returns -1 if no bits are set.
        /// </summary>
        public static int TrailingZeroCount(ulong value)
        {
            if (value == 0)
                return -1;

            int count = 0;
            while ((value & 1) == 0)
            {
                count++;
                value >>= 1;
            }
            return count;
        }

        /// <summary>
        /// Finds the index of the most significant set bit (0-based).
        /// Returns -1 if no bits are set.
        /// </summary>
        public static int LeadingZeroCount(ulong value)
        {
            if (value == 0)
                return -1;

            int count = 0;
            ulong mask = 1UL << 63;
            while ((value & mask) == 0)
            {
                count++;
                mask >>= 1;
            }
            return count;
        }
    }
}
