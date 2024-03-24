class Bitboard:
    def __init__(self, *args):
        if len(args) == 1:
            self.bitboard = args[0]
        else:
            self.bitboard = 0

    def add_piece(self, square_no):
        self.bitboard |= 1 << square_no

    def __str__(self):
        return format(self.bitboard, '064b')

    def get_8by8_board(self):
        bit_string = format(self.bitboard, '064b')
        formatted_string = '\n'.join(bit_string[i:i + 8] for i in range(0, 64, 8))
        return formatted_string

    def __add__(self, other):
        return Bitboard(self.bitboard | other.bitboard)

    def __or__(self, other):
        return Bitboard(self.bitboard | other.bitboard)

    def __xor__(self, other):
        return Bitboard(self.bitboard ^ other.bitboard)

    def __and__(self, other):
        return Bitboard(self.bitboard & other.bitboard)

