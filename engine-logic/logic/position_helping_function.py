from logic.bitboard import Position, PieceEnum, PieceColor, GenerateMove, CastleEnum, get_num_by_square_name


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