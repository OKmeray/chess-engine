from core.domain.engine.Bitboard import Bitboard
from core.domain.engine.enums import PieceColor


class GenerateMove:
    rank_1 = 18374686479671623680  # 1111111100000000000000000000000000000000000000000000000000000000
    rank_2 = 71776119061217280  # 0000000011111111000000000000000000000000000000000000000000000000
    rank_7 = 65280  # 0000000000000000000000000000000000000000000000001111111100000000
    rank_8 = 255  # 0000000000000000000000000000000000000000000000000000000011111111

    file_a = 72340172838076673  # 0000000100000001000000010000000100000001000000010000000100000001
    file_b = 144680345676153346  # 0000001000000010000000100000001000000010000000100000001000000010
    file_g = 4629771061636907072  # 0100000001000000010000000100000001000000010000000100000001000000
    file_h = 9259542123273814144  # 1000000010000000100000001000000010000000100000001000000010000000

    file_e = 1157442765409226768  # 0001000000010000000100000001000000010000000100000001000000010000

    @classmethod
    def is_a_file(cls, square):
        return cls.file_a & square == square

    @classmethod
    def is_b_file(cls, square):
        return cls.file_b & square == square

    @classmethod
    def is_e_file(cls, square):
        return cls.file_e & square == square

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

    @classmethod
    def handle_pawn_promotion(cls):
        pass

    @classmethod
    def generate_pawn_move(cls, square, own_bitboard, enemy_bitboard, color, en_passant_square):
        moves = []

        if color == PieceColor.WHITE:
            # capture and en passant
            if ((square >> 9) & enemy_bitboard.bitboard == square >> 9) or (en_passant_square is not None and square >> 9 == en_passant_square) and \
                    not cls.is_h_file(square):
                moves.append(square >> 9)
            # capture and en passant
            if ((square >> 7) & enemy_bitboard.bitboard == square >> 7) or (en_passant_square is not None and square >> 7 == en_passant_square) and \
                    not cls.is_a_file(square):
                moves.append(square >> 7)
            # one square move
            if (square >> 8) & (enemy_bitboard.bitboard | own_bitboard.bitboard) != square >> 8:
                moves.append(square >> 8)

            # two square move
            if cls.is_rank_2(square) and ((square >> 16) & (enemy_bitboard.bitboard | own_bitboard.bitboard) != square >> 16):
                moves.append(square >> 16)

        elif color == PieceColor.BLACK:
            # capture and en passant
            if ((square << 9) & enemy_bitboard.bitboard == square << 9) or (en_passant_square is not None and square << 9 == en_passant_square) and \
                    not cls.is_h_file(square):
                moves.append(square << 9)
            # capture and en passant
            if ((square << 7) & enemy_bitboard.bitboard == square << 7) or (en_passant_square is not None and square << 9 == en_passant_square) and \
                    not cls.is_a_file(square):
                moves.append(square << 7)
            # one square move
            if (square << 8) & (enemy_bitboard.bitboard | own_bitboard.bitboard) != square << 8:
                moves.append(square << 8)

            # two square move
            if cls.is_rank_7(square) and ((square << 16) & (enemy_bitboard.bitboard | own_bitboard.bitboard) != square << 16):
                moves.append(square << 16)
        return moves

    @classmethod
    def generate_knight_move(cls, square: int, own_bitboard: Bitboard) -> list[int]:
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
    def generate_castle_moves(cls, square, own_bitboard, enemy_bitboard, is_white=True):
        moves = []
        if (is_white and cls.is_rank_1(square) and cls.is_e_file(square)) or\
           (not is_white and cls.is_rank_8(square) and cls.is_e_file(square)):

            # if not ((square << 1 | square << 2) & (own_bitboard.bitboard | enemy_bitboard.bitboard)):
            #     moves.append(square << 2)
            # if not ((square >> 1 | square >> 2) & (own_bitboard.bitboard | enemy_bitboard.bitboard)):
            #     moves.append(square >> 2)
            if len(cls.generate_horizontal_left(square, own_bitboard, enemy_bitboard, steps=2)):
                moves.append(square >> 2)
            if len(cls.generate_horizontal_right(square, own_bitboard, enemy_bitboard, steps=2)):
                moves.append(square << 2)

        return moves

    @classmethod
    def generate_king_move(cls, square, own_bitboard, enemy_bitboard, is_white=True):
        # TODO: catsle moves (check target squares) (if enemy is on g1 then king can't castle)
        return [
            *cls.generate_castle_moves(square, own_bitboard, enemy_bitboard, is_white),
            *cls.generate_rook_move(square, own_bitboard, enemy_bitboard, steps=1),
            *cls.generate_bishop_move(square, own_bitboard, enemy_bitboard, steps=1)
        ]

