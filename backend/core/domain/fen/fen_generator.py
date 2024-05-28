
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
from core.domain.engine.square_helping_functions import get_num_by_square_name, get_bitboard_from_num


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


def get_fen_from_position(position: Position) -> str:
    # Convert the board to FEN piece placement
    fen_rows = []
    empty_count = 0

    for i in range(64):
        piece_code = sum([int(i) for i in position.get_piece_and_color_by_square(i)])

        if piece_code == 0:
            empty_count += 1
        else:
            if empty_count > 0:
                fen_rows.append(str(empty_count))
                empty_count = 0
            fen_rows.append(piece_code_to_fen_char(piece_code))

        if (i + 1) % 8 == 0:
            if empty_count > 0:
                fen_rows.append(str(empty_count))
                empty_count = 0
            if i != 63:
                fen_rows.append("/")

    # Active color
    active_color = 'w' if position.side_to_move == PieceColor.WHITE else 'b'

    # Castling availability
    castling_rights = []
    if position.castling_rights[CastleEnum.WhiteShortCastle]:
        castling_rights.append('K')
    if position.castling_rights[CastleEnum.WhiteLongCastle]:
        castling_rights.append('Q')
    if position.castling_rights[CastleEnum.BlackShortCastle]:
        castling_rights.append('k')
    if position.castling_rights[CastleEnum.BlackLongCastle]:
        castling_rights.append('q')
    castling_rights = ''.join(castling_rights) if castling_rights else '-'

    # En passant target square
    en_passant = '-' if not position.en_passant_square else square_to_name(position.en_passant_square)

    # Halfmove clock and fullmove number
    half_moves = str(position.half_moves)
    full_moves = str(position.current_turn)

    return ' '.join([
        ''.join(fen_rows),
        active_color,
        castling_rights,
        en_passant,
        half_moves,
        full_moves
    ])


def piece_code_to_fen_char(piece_code):
    # Maps internal piece codes to FEN characters
    return {
        0b1001: 'P', 0b1010: 'N', 0b1011: 'B', 0b1100: 'R', 0b1101: 'Q', 0b1110: 'K',
        0b0001: 'p', 0b0010: 'n', 0b0011: 'b', 0b0100: 'r', 0b0101: 'q', 0b0110: 'k'
    }.get(piece_code, None)


def square_to_name(square):
    # Convert square index to algebraic notation (e.g., 0 -> a8)
    file = square % 8
    rank = 8 - (square // 8)
    return f"{chr(file + 97)}{rank}"

