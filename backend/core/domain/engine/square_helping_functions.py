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


def get_square_name_by_num(num: int) -> str:
    match num % 8:
        case 0:
            vertical = "a"
        case 1:
            vertical = "b"
        case 2:
            vertical = "c"
        case 3:
            vertical = "d"
        case 4:
            vertical = "e"
        case 5:
            vertical = "f"
        case 6:
            vertical = "g"
        case 7:
            vertical = "h"
        case _:
            vertical = ""

    match num // 8:
        case 0:
            horizontal = "8"
        case 1:
            horizontal = "7"
        case 2:
            horizontal = "6"
        case 3:
            horizontal = "5"
        case 4:
            horizontal = "4"
        case 5:
            horizontal = "3"
        case 6:
            horizontal = "2"
        case 7:
            horizontal = "1"
        case _:
            horizontal = ""

    return vertical + horizontal


def get_num_by_square_name(square: str) -> int:
    match square[0]:
        case "a":
            vertical = 0
        case "b":
            vertical = 1
        case "c":
            vertical = 2
        case "d":
            vertical = 3
        case "e":
            vertical = 4
        case "f":
            vertical = 5
        case "g":
            vertical = 6
        case "h":
            vertical = 7
        case _:
            vertical = ""

    match square[1]:
        case "8":
            horizontal = 0
        case "7":
            horizontal = 1
        case "6":
            horizontal = 2
        case "5":
            horizontal = 3
        case "4":
            horizontal = 4
        case "3":
            horizontal = 5
        case "2":
            horizontal = 6
        case "1":
            horizontal = 7
        case _:
            horizontal = ""

    return vertical + horizontal * 8


def get_bitboard_from_num(num):
    try:
        return 1 << num
    except Exception as e:
        raise Exception(str(num))


def get_num_from_bitboard(bitboard_number):
    return int.bit_length(bitboard_number) - 1


def get_nums_from_bit_nums(bitboard_numbers):
    return [get_num_from_bitboard(bitboard_num) for bitboard_num in bitboard_numbers]
