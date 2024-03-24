import copy

from logic.engine.bitboard import Bitboard
from logic.engine.square_helping_functions import get_bitboard_from_num, get_nums_from_bit_nums
from logic.evaluation.evaluation import evaluate_position, PieceEnum, PieceColor, CastleEnum
from logic.moves.generate_move import GenerateMove


def get_num_from_bitboard(bitboard_number):
    return int.bit_length(bitboard_number) - 1


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

        self.side_to_move = PieceColor.WHITE

        self.castling_rights = {
            CastleEnum.BlackShortCastle: False,
            CastleEnum.BlackLongCastle: False,
            CastleEnum.WhiteShortCastle: False,
            CastleEnum.WhiteLongCastle: False,
        }  # the castle item should change if:
        # 1) king moves
        # 2) rook moves
        # 3) rook is captured
        # also you shouldn't forget, that castle squares can be
        # attacked, should check it while validating for legal moves

        self.en_passant_square = None  # for example if e pawn just moved two squares then I should set
        # self.enpassant_file = 8 * 2 + 4 and then when I am doing check for pawn moves I have to check if
        # current pawn square +7 or +9 (or for the other color -9 and -7) equals enpassant_square after either
        # scenario I have to set enpassant_square to either None or to some other file if another color also made
        # two square pawn move

        self.half_moves = 0
        self.current_turn = 1

    def clone(self):
        # Create a new Position object
        cloned_position = Position()

        # Deep copy complex attributes
        cloned_position.bitboard_position = copy.deepcopy(self.bitboard_position)
        cloned_position.side_to_move = self.side_to_move
        cloned_position.castling_rights = copy.deepcopy(self.castling_rights)
        cloned_position.en_passant_square = self.en_passant_square
        cloned_position.half_moves = self.half_moves
        cloned_position.current_turn = self.current_turn

        return cloned_position

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

    def get_all_pieces_of_one_color(self, color: PieceColor = None):
        if color is None:
            bitboard = self.get_all_bitboard()
        elif color == PieceColor.WHITE:
            bitboard = self.get_white_bitboard()
        else:
            bitboard = self.get_black_bitboard()

        j = 1
        arr = []
        while j <= 9223372036854775808:  # 2 ** 63
            arr.append(bitboard.bitboard & j) if bitboard.bitboard & j == j else None
            j <<= 1

        return arr

    def get_separate_piece(self, piece: PieceEnum, color: PieceColor):
        separate_pieces = []

        all_pieces_of_one_color = self.get_all_pieces_of_one_color(color=color)
        piece_bitboard = self.get_piece_bitboard(piece, color)

        for piece in all_pieces_of_one_color:
            separate_pieces.append(piece & piece_bitboard.bitboard) if piece & piece_bitboard.bitboard != 0 else None

        return separate_pieces

    # {'piece': <PieceEnum.BISHOP: 3>, 'color': <PieceColor.WHITE: 8>, 'square': 49, 'move': 40}
    def apply_move(self, move_detail):
        old_square = get_bitboard_from_num(move_detail["square"])
        new_square = get_bitboard_from_num(move_detail["move"])

        for num_color in range(0, 9, 8):  # for all the pieces clear the new_square
            # (only one piece can have this square)
            for num_piece in range(1, 7):
                self.bitboard_position[num_piece + num_color].bitboard &= ~new_square

        self.bitboard_position[move_detail["piece"] + move_detail["color"]].bitboard = \
            self.bitboard_position[move_detail["piece"] + move_detail["color"]].bitboard & ~old_square | new_square

        self.change_side_to_move()
        # TODO: change .half_moves (check for captures, pawn moves)
        # TODO: change .current_turn

    def get_all_moves(self):
        pseudo_legal_moves = self.generate_all_pseudo_legal_moves()

        legal_moves = self.filter_illegal_moves(pseudo_legal_moves)
        return legal_moves

    def generate_all_pseudo_legal_moves(self):
        result = []

        # for num_color in range(0, 9, 8):
        num_color = self.side_to_move
        own_bitboard = self.get_white_bitboard() if num_color else self.get_black_bitboard()
        enemy_bitboard = self.get_white_bitboard() if not num_color else self.get_black_bitboard()

        for num_piece in range(1, 7):
            separate_bit_pieces = self.get_separate_piece(PieceEnum(num_piece), PieceColor(num_color))

            for bit_piece in separate_bit_pieces:
                separate_pieces_moves = []
                if num_piece == PieceEnum.PAWN:
                    bitboard_en_passant_square = None
                    if self.en_passant_square is not None:
                        bitboard_en_passant_square = 2 ** self.en_passant_square
                    separate_pieces_moves = GenerateMove.generate_pawn_move(bit_piece, own_bitboard, enemy_bitboard, color=num_color, en_passant_square=bitboard_en_passant_square)
                elif num_piece == PieceEnum.KNIGHT:
                    separate_pieces_moves = GenerateMove.generate_knight_move(bit_piece, own_bitboard)
                elif num_piece == PieceEnum.BISHOP:
                    separate_pieces_moves = GenerateMove.generate_bishop_move(bit_piece, own_bitboard, enemy_bitboard)
                elif num_piece == PieceEnum.ROOK:
                    separate_pieces_moves = GenerateMove.generate_rook_move(bit_piece, own_bitboard, enemy_bitboard)
                elif num_piece == PieceEnum.QUEEN:
                    separate_pieces_moves = GenerateMove.generate_queen_move(bit_piece, own_bitboard, enemy_bitboard)
                elif num_piece == PieceEnum.KING:
                    separate_pieces_moves = GenerateMove.generate_king_move(bit_piece, own_bitboard, enemy_bitboard)

                piece_moves_num = get_nums_from_bit_nums(separate_pieces_moves)
                local_dictionary = {
                    "piece": PieceEnum(num_piece),
                    "color": PieceColor(num_color),
                    "square": get_num_from_bitboard(bit_piece),
                    "possible_moves": piece_moves_num
                }
                result.append(local_dictionary)

        return result

    def filter_illegal_moves(self, pseudo_legal_moves):
        separate_legal_moves = []

        for piece in pseudo_legal_moves:
            for move in piece['possible_moves']:
                move_detail = {
                    'piece': piece['piece'],
                    'color': piece['color'],
                    'square': piece['square'],
                    'move': move
                }
                if self.is_move_legal(move_detail):
                    separate_legal_moves.append(move_detail)

        helping_dictionary = {}
        for separate_move in separate_legal_moves:
            key = separate_move['piece'] + separate_move['color']
            if key not in helping_dictionary.keys():
                helping_dictionary[key] = {
                    'piece': separate_move['piece'],
                    'color': separate_move['color'],
                    'square': separate_move['square'],
                    'possible_moves': [separate_move['move']]
                }
            else:
                helping_dictionary[key]['possible_moves'].append(separate_move['move'])

        legal_moves = helping_dictionary.values()

        return legal_moves

    def is_move_legal(self, move):
        position_copy = self.clone()
        position_copy.apply_move(move)
        position_copy.change_side_to_move()  # changing the side to check for the side which just moved
        return not position_copy.is_king_in_check()

    def determine_outcome(self):
        evaluation = evaluate_position(position=self, side=self.side_to_move)

        if self.is_checkmate():
            return 'loss'  # If the current player is in checkmate, it's a loss for them.
        elif self.is_stalemate() or self.is_draw():
            return 'draw'  # Stalemate or other draw conditions.
        # TODO: correct the evaluation
        # elif evaluation >= 20:  # TODO: maybe if the evaluation is 20+ then the win for the side
        #     return 'win'
        else:
            # If the game isn't over, or if implementing specific win conditions isn't feasible,
            # you might default to a draw or continue the simulation until a definitive outcome is reached.
            return 'draw'

    # accepts:
    def is_king_in_check(self):
        king = self.get_separate_piece(piece=PieceEnum.KING, color=self.side_to_move)

        self.change_side_to_move()
        all_pseudo_legal_moves = self.generate_all_pseudo_legal_moves()
        self.change_side_to_move()

        for piece in all_pseudo_legal_moves:
            for move in piece['possible_moves']:
                if get_bitboard_from_num(move) == king[0]:
                    return True
        return False

    def is_checkmate(self):
        return self.is_king_in_check() and len(self.get_all_moves()) == 0

    def is_stalemate(self):
        return not self.is_king_in_check() and len(self.get_all_moves()) == 0

    def is_draw(self):
        return self.is_threefold_repetition() or self.is_fifty_move_rule_reached() or self.is_insufficient_material()

    def is_threefold_repetition(self):
        return False  # TODO: implement

    def is_fifty_move_rule_reached(self):
        return self.half_moves >= 100

    def is_insufficient_material(self):
        return False  # TODO: implement

    def change_side_to_move(self):
        if self.side_to_move == PieceColor.WHITE:
            self.side_to_move = PieceColor.BLACK
        else:
            self.side_to_move = PieceColor.WHITE
