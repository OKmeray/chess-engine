using System;

namespace ChessEngine.Helpers
{
    public static class SquareHelper
    {
        public static string GetSquareNameByNum(int num)
        {
            if (num < 0 || num > 63)
                throw new ArgumentOutOfRangeException(nameof(num), "Square number must be between 0 and 63.");

            char file = (char)('a' + (num % 8));
            char rank = (char)('8' - (num / 8));
            return $"{file}{rank}";
        }

        public static int GetNumBySquareName(string square)
        {
            if (square.Length != 2)
                throw new ArgumentException("Square name must be exactly 2 characters.", nameof(square));

            char fileChar = square[0];
            char rankChar = square[1];

            if (fileChar < 'a' || fileChar > 'h')
                throw new ArgumentException("Invalid file character.", nameof(square));
            if (rankChar < '1' || rankChar > '8')
                throw new ArgumentException("Invalid rank character.", nameof(square));

            int file = fileChar - 'a';
            int rank = '8' - rankChar;
            return rank * 8 + file;
        }

        public static ulong GetBitboardFromNum(int num)
        {
            if (num < 0 || num > 63)
                throw new ArgumentOutOfRangeException(nameof(num), "Square number must be between 0 and 63.");

            return 1UL << num;
        }

        public static int GetNumFromBitboard(ulong bitboardNumber)
        {
            if (bitboardNumber == 0)
                throw new ArgumentException("Bitboard number cannot be zero.", nameof(bitboardNumber));

            return (int)(BitHelper.TrailingZeroCount(bitboardNumber));
        }

        public static int[] GetNumsFromBitboards(ulong[] bitboardNumbers)
        {
            int[] nums = new int[bitboardNumbers.Length];
            for (int i = 0; i < bitboardNumbers.Length; i++)
            {
                nums[i] = GetNumFromBitboard(bitboardNumbers[i]);
            }
            return nums;
        }
    }
}