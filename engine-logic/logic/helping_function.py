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
from logic.bitboard import Position, PieceEnum, PieceColor, GenerateMove


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

def get_board_arr(fen_string: str) -> list[int]:
    pieces_placement = fen_string.split()[0]

    board = []
    for char in pieces_placement:
        if char.isdigit():
            for i in range(int(char)):
                board.append(0b0000)
        elif char == "/":
            continue
        else:
            match char:
                case "P":
                    board.append(0b1001)

                case "N":
                    board.append(0b1010)

                case "B":
                    board.append(0b1011)

                case "R":
                    board.append(0b1100)

                case "Q":
                    board.append(0b1101)

                case "K":
                    board.append(0b1110)

                case "p":
                    board.append(0b0001)

                case "n":
                    board.append(0b0010)

                case "b":
                    board.append(0b0011)

                case "r":
                    board.append(0b0100)

                case "q":
                    board.append(0b0101)

                case "k":
                    board.append(0b0110)

                case _:
                    raise Exception()

    return board


def get_position_from_fen(fen: str) -> Position:
    board = get_board_arr(fen)

    position = Position()
    for index, piece in enumerate(board):
        if piece != 0:
            position.add_piece_by_int(piece, index)

    return position


def get_all_pieces_of_one_color(position: Position, color: PieceColor = None):
    if color is None:
        bitboard = position.get_all_bitboard()
    elif color == PieceColor.WHITE:
        bitboard = position.get_white_bitboard()
    else:
        bitboard = position.get_black_bitboard()

    j = 1
    arr = []
    while j <= 9223372036854775808:  # 2 ** 63
        arr.append(bitboard.bitboard & j) if bitboard.bitboard & j == j else None
        j <<= 1

    return arr


def get_separate_piece(position: Position, piece: PieceEnum, color: PieceColor):
    separate_pieces = []

    all_pieces_of_one_color = get_all_pieces_of_one_color(position=position, color=color)
    piece_bitboard = position.get_piece_bitboard(piece, color)

    for piece in all_pieces_of_one_color:
        separate_pieces.append(piece & piece_bitboard.bitboard) if piece & piece_bitboard.bitboard != 0 else None

    return separate_pieces

# def get_all_moves_of_piece(position: Position, piece: PieceEnum, color: PieceColor):
#     pass


def get_num_from_bitboard(bitboard_number):
    return int.bit_length(bitboard_number) - 1


def get_nums_from_bit_nums(bitboard_numbers):
    return [get_num_from_bitboard(bitboard_num) for bitboard_num in bitboard_numbers]
