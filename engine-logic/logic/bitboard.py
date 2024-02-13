from enum import IntEnum


# # ------------------------
# WHITE_PAWN = 0b1001  # 9
# WHITE_KNIGHT = 0b1010  # 10
# WHITE_BISHOP = 0b1011  # 11
# WHITE_ROOK = 0b1100  # 12
# WHITE_QUEEN = 0b1101  # 13
# WHITE_KING = 0b1110  # 14
# # ------------------------
# BLACK_PAWN = 0b0001  # 1
# BLACK_KNIGHT = 0b0010  # 2
# BLACK_BISHOP = 0b0011  # 3
# BLACK_ROOK = 0b0100  # 4
# BLACK_QUEEN = 0b0101  # 5
# BLACK_KING = 0b0110  # 6
# # ------------------------
# <<  ->  +
# >>  ->  -

class PieceEnum(IntEnum):
    NONE = 0
    PAWN = 1
    KNIGHT = 2
    BISHOP = 3
    ROOK = 4
    QUEEN = 5
    KING = 6


class PieceColor(IntEnum):
    BLACK = 0
    WHITE = 8


class CastleEnum:
    BlackShortCastle = 0
    BlackLongCastle = 1
    WhiteShortCastle = 8
    WhiteLongCastle = 9


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


class Position:
    def __init__(self):
        self.bitboard_position = {
            PieceEnum.PAWN + PieceColor.BLACK: Bitboard(),
            PieceEnum.KNIGHT + PieceColor.BLACK: Bitboard(),
            PieceEnum.BISHOP + PieceColor.BLACK: Bitboard(),
            PieceEnum.ROOK + PieceColor.BLACK: Bitboard(),
            PieceEnum.QUEEN + PieceColor.BLACK: Bitboard(),
            PieceEnum.KING + PieceColor.BLACK: Bitboard(),
            PieceEnum.PAWN + PieceColor.WHITE: Bitboard(),
            PieceEnum.KNIGHT + PieceColor.WHITE: Bitboard(),
            PieceEnum.BISHOP + PieceColor.WHITE: Bitboard(),
            PieceEnum.ROOK + PieceColor.WHITE: Bitboard(),
            PieceEnum.QUEEN + PieceColor.WHITE: Bitboard(),
            PieceEnum.KING + PieceColor.WHITE: Bitboard(),
        }

        self.enpassant_square = None  # for example if e pawn just moved two squares then I should set
        # self.enpassant_file = 8 * 2 + 4 and then when I am doing check for pawn moves I have to check if
        # current pawn square +7 or +9 (or for the other color -9 and -7) equals enpassant_square after either
        # scenario I have to set enpassant_square to either None or to some other file if another color also made
        # two square pawn move

        self.castle = {
            CastleEnum.BlackShortCastle: True,
            CastleEnum.BlackLongCastle: True,
            CastleEnum.WhiteShortCastle: True,
            CastleEnum.WhiteLongCastle: True,
        }  # the castle item should change if:
        # 1) king moves
        # 2) rook moves
        # 3) rook is captured
        # also you shouldn't forget, that castle squares can be
        # attacked, should check it while validating for legal moves

    def add_piece(self, piece: PieceEnum, color: PieceColor, square: int):
        self.bitboard_position[piece + color].add_piece(square)

    def add_piece_by_int(self, int_piece: int, square: int):
        self.bitboard_position[int_piece].add_piece(square)

    def get_piece_bitboard(self, piece: PieceEnum, color: PieceColor):
        return self.bitboard_position[piece + color]

    def get_white_bitboard(self):
        white_bitboard = Bitboard()
        for i in range(8, 16):
            if i in self.bitboard_position.keys():
                white_bitboard += self.bitboard_position[i]

        return white_bitboard

    def get_black_bitboard(self):
        black_bitboard = Bitboard()
        for i in range(0, 8):
            if i in self.bitboard_position.keys():
                black_bitboard += self.bitboard_position[i]

        return black_bitboard

    def get_all_bitboard(self):
        return self.get_white_bitboard() + self.get_black_bitboard()


class GenerateMove:
    rank_1 = 18374686479671623680  # 1111111100000000000000000000000000000000000000000000000000000000
    rank_2 = 71776119061217280  # 0000000011111111000000000000000000000000000000000000000000000000
    rank_7 = 65280  # 0000000000000000000000000000000000000000000000001111111100000000
    rank_8 = 255  # 0000000000000000000000000000000000000000000000000000000011111111

    file_a = 72340172838076673  # 0000000100000001000000010000000100000001000000010000000100000001
    file_b = 144680345676153346  # 0000001000000010000000100000001000000010000000100000001000000010
    file_g = 4629771061636907072  # 0100000001000000010000000100000001000000010000000100000001000000
    file_h = 9259542123273814144  # 1000000010000000100000001000000010000000100000001000000010000000

    @classmethod
    def is_a_file(cls, square):
        return cls.file_a & square == square

    @classmethod
    def is_b_file(cls, square):
        return cls.file_b & square == square

    @classmethod
    def is_g_file(cls, square):
        return cls.file_g & square == square

    @classmethod
    def is_h_file(cls, square):
        return cls.file_h & square == square

    @classmethod
    def is_rank_1(cls, square):
        return cls.rank_1 & square == square

    @classmethod
    def is_rank_2(cls, square):
        return cls.rank_2 & square == square

    @classmethod
    def is_rank_7(cls, square):
        return cls.rank_7 & square == square

    @classmethod
    def is_rank_8(cls, square):
        return cls.rank_8 & square == square

    # TODO:
    @classmethod
    def generate_pawn_move(cls):
        pass

    @classmethod
    def generate_knight_move(cls, own_bitboard: Bitboard, square: int) -> list[int]:
        moves = []
        if not (
                cls.is_a_file(square) or
                cls.is_b_file(square) or
                cls.is_rank_8(square)
        ):
            attack_square = square >> 10  # -10

            # check for own pieces on target square
            if not own_bitboard.bitboard & attack_square == attack_square:
                moves.append(attack_square)

        if not (
                cls.is_a_file(square) or
                cls.is_rank_7(square) or
                cls.is_rank_8(square)
        ):
            attack_square = square >> 17  # -17

            # check for own pieces on target square
            if not own_bitboard.bitboard & attack_square == attack_square:
                moves.append(attack_square)

        if not (
                cls.is_h_file(square) or
                cls.is_g_file(square) or
                cls.is_rank_8(square)
        ):
            attack_square = square >> 6  # -6

            # check for own pieces on target square
            if not own_bitboard.bitboard & attack_square == attack_square:
                moves.append(attack_square)

        if not (
                cls.is_h_file(square) or
                cls.is_rank_7(square) or
                cls.is_rank_8(square)
        ):
            attack_square = square >> 15  # -15

            # check for own pieces on target square
            if not own_bitboard.bitboard & attack_square == attack_square:
                moves.append(attack_square)

        if not (
                cls.is_a_file(square) or
                cls.is_b_file(square) or
                cls.is_rank_1(square)
        ):
            attack_square = square << 6  # +6

            # check for own pieces on target square
            if not own_bitboard.bitboard & attack_square == attack_square:
                moves.append(attack_square)

        if not (
                cls.is_a_file(square) or
                cls.is_rank_2(square) or
                cls.is_rank_1(square)
        ):
            attack_square = square << 15  # +15

            # check for own pieces on target square
            if not own_bitboard.bitboard & attack_square == attack_square:
                moves.append(attack_square)

        if not (
                cls.is_h_file(square) or
                cls.is_g_file(square) or
                cls.is_rank_1(square)
        ):
            attack_square = square << 10  # +10
            
            # check for own pieces on target square
            if not own_bitboard.bitboard & attack_square == attack_square:
                moves.append(attack_square)

        if not (
                cls.is_h_file(square) or
                cls.is_rank_2(square) or
                cls.is_rank_1(square)
        ):
            attack_square = square << 17  # +17
            
            # check for own pieces on target square
            if not own_bitboard.bitboard & attack_square == attack_square:
                moves.append(attack_square)

        return moves

    @classmethod
    def generate_diagonal_up(cls, square, own_bitboard, enemy_bitboard, steps: int = 8):  # -9
        moves = []
        while steps > 0:
            if not (
                    cls.is_a_file(square) or
                    cls.is_rank_8(square)
            ):
                square >>= 9

                # check for own or enemy pieces on target square
                if square & own_bitboard.bitboard == square:
                    return moves
                elif square & enemy_bitboard.bitboard == square:
                    moves.append(square)
                    return moves
                else:
                    moves.append(square)
                    steps -= 1
            else:
                return moves
        return moves

    @classmethod
    def generate_diagonal_down(cls, square, own_bitboard, enemy_bitboard, steps: int = 8):  # +9
        moves = []
        while steps > 0:
            if not (
                    cls.is_h_file(square) or
                    cls.is_rank_1(square)
            ):
                square <<= 9

                # check for own or enemy pieces on target square
                if square & own_bitboard.bitboard == square:
                    return moves
                elif square & enemy_bitboard.bitboard == square:
                    moves.append(square)
                    return moves
                else:
                    moves.append(square)
                    steps -= 1
            else:
                return moves
        return moves

    @classmethod
    def generate_antidiagonal_up(cls, square, own_bitboard, enemy_bitboard, steps: int = 8):  # -7
        moves = []
        while steps > 0:
            if not (
                    cls.is_h_file(square) or
                    cls.is_rank_8(square)
            ):
                square >>= 7

                # check for own or enemy pieces on target square
                if square & own_bitboard.bitboard == square:
                    return moves
                elif square & enemy_bitboard.bitboard == square:
                    moves.append(square)
                    return moves
                else:
                    moves.append(square)
                    steps -= 1
            else:
                return moves
        return moves

    @classmethod
    def generate_antidiagonal_down(cls, square, own_bitboard, enemy_bitboard, steps: int = 8):  # +7
        moves = []
        while steps > 0:
            if not (
                    cls.is_a_file(square) or
                    cls.is_rank_1(square)
            ):
                square <<= 7

                # check for own or enemy pieces on target square
                if square & own_bitboard.bitboard == square:
                    return moves
                elif square & enemy_bitboard.bitboard == square:
                    moves.append(square)
                    return moves
                else:
                    moves.append(square)
                    steps -= 1
            else:
                return moves
        return moves

    @classmethod
    def generate_vertical_up(cls, square, own_bitboard, enemy_bitboard, steps: int = 8):
        moves = []
        while steps > 0:
            if not cls.is_rank_8(square):
                square >>= 8

                # check for own or enemy pieces on target square
                if square & own_bitboard.bitboard == square:
                    return moves
                elif square & enemy_bitboard.bitboard == square:
                    moves.append(square)
                    return moves
                else:
                    moves.append(square)
                    steps -= 1
            else:
                return moves
        return moves

    @classmethod
    def generate_vertical_down(cls, square, own_bitboard, enemy_bitboard, steps: int = 8):
        moves = []
        while steps > 0:
            if not cls.is_rank_1(square):
                square <<= 8

                # check for own or enemy pieces on target square
                if square & own_bitboard.bitboard == square:
                    return moves
                elif square & enemy_bitboard.bitboard == square:
                    moves.append(square)
                    return moves
                else:
                    moves.append(square)
                    steps -= 1
            else:
                return moves
        return moves

    @classmethod
    def generate_horizontal_left(cls, square, own_bitboard, enemy_bitboard, steps: int = 8):
        moves = []
        while steps > 0:
            if not cls.is_a_file(square):
                square >>= 1

                # check for own or enemy pieces on target square
                if square & own_bitboard.bitboard == square:
                    return moves
                elif square & enemy_bitboard.bitboard == square:
                    moves.append(square)
                    return moves
                else:
                    moves.append(square)
                    steps -= 1
            else:
                return moves
        return moves

    @classmethod
    def generate_horizontal_right(cls, square, own_bitboard, enemy_bitboard, steps: int = 8):
        moves = []
        while steps > 0:
            if not cls.is_h_file(square):
                square <<= 1

                # check for own or enemy pieces on target square
                if square & own_bitboard.bitboard == square:
                    return moves
                elif square & enemy_bitboard.bitboard == square:
                    moves.append(square)
                    return moves
                else:
                    moves.append(square)
                    steps -= 1
            else:
                return moves
        return moves

    @classmethod
    def generate_rook_move(cls, square, own_bitboard, enemy_bitboard, steps: int = 8):
        return [
            *cls.generate_vertical_up(square, own_bitboard, enemy_bitboard, steps),
            *cls.generate_vertical_down(square, own_bitboard, enemy_bitboard, steps),
            *cls.generate_horizontal_left(square, own_bitboard, enemy_bitboard, steps),
            *cls.generate_horizontal_right(square, own_bitboard, enemy_bitboard, steps)
        ]

    @classmethod
    def generate_bishop_move(cls, square, own_bitboard, enemy_bitboard, steps: int = 8):
        return [
            *cls.generate_diagonal_up(square, own_bitboard, enemy_bitboard, steps),
            *cls.generate_diagonal_down(square, own_bitboard, enemy_bitboard, steps),
            *cls.generate_antidiagonal_up(square, own_bitboard, enemy_bitboard, steps),
            *cls.generate_antidiagonal_down(square, own_bitboard, enemy_bitboard, steps)
        ]

    @classmethod
    def generate_queen_move(cls, square, own_bitboard, enemy_bitboard):
        return [
            *cls.generate_rook_move(square, own_bitboard, enemy_bitboard),
            *cls.generate_bishop_move(square, own_bitboard, enemy_bitboard)
        ]

    @classmethod
    def generate_king_move(cls, square, own_bitboard, enemy_bitboard):
        return [
            *cls.generate_rook_move(square, own_bitboard, enemy_bitboard, steps=1),
            *cls.generate_bishop_move(square, own_bitboard, enemy_bitboard, steps=1)
        ]
