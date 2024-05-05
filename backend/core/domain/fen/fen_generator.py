
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


from core.domain.engine.enums import PieceColor, CastleEnum
from core.domain.engine.Position import Position
from core.domain.engine.square_helping_functions import get_num_by_square_name


def set_piece_placement(position: Position, pieces_placement_fen: str):
    pieces = []

    for char in pieces_placement_fen:
        if char.isdigit():
            for i in range(int(char)):
                pieces.append(0b0000)
        elif char == "/":
            continue
        else:
            match char:
                case "P":
                    pieces.append(0b1001)

                case "N":
                    pieces.append(0b1010)

                case "B":
                    pieces.append(0b1011)

                case "R":
                    pieces.append(0b1100)

                case "Q":
                    pieces.append(0b1101)

                case "K":
                    pieces.append(0b1110)

                case "p":
                    pieces.append(0b0001)

                case "n":
                    pieces.append(0b0010)

                case "b":
                    pieces.append(0b0011)

                case "r":
                    pieces.append(0b0100)

                case "q":
                    pieces.append(0b0101)

                case "k":
                    pieces.append(0b0110)

                case _:
                    raise Exception()

    for index, piece in enumerate(pieces):
        if piece != 0:
            position.add_piece_by_int(piece, index)


def set_side_to_move(position: Position, side_to_move_fen: str):
    if side_to_move_fen == "w":
        position.side_to_move = PieceColor.WHITE
    else:
        position.side_to_move = PieceColor.BLACK


def set_castling_rights(position: Position, castling_rights_fen: str):
    if "k" in castling_rights_fen:
        position.castling_rights[CastleEnum.BlackShortCastle] = True
    if "q" in castling_rights_fen:
        position.castling_rights[CastleEnum.BlackLongCastle] = True
    if "K" in castling_rights_fen:
        position.castling_rights[CastleEnum.WhiteShortCastle] = True
    if "Q" in castling_rights_fen:
        position.castling_rights[CastleEnum.WhiteLongCastle] = True


def set_en_passant(position: Position, en_passant_fen: str):
    if en_passant_fen != "-":
        position.en_passant_square = get_num_by_square_name(en_passant_fen)


def get_position_from_fen(fen: str) -> Position:
    fen_splited = fen.split()
    position = Position()

    set_piece_placement(position, fen_splited[0])
    set_side_to_move(position, fen_splited[1])
    set_castling_rights(position, fen_splited[2])
    set_en_passant(position, fen_splited[3])
    position.half_moves = int(fen_splited[4])
    position.current_turn = int(fen_splited[5])

    return position

