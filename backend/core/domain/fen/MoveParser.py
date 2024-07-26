import re
from core.domain.engine.enums import PieceEnum, PieceColor
from core.domain.fen.fen_generator import get_position_from_fen


piece_map = {
    'P': PieceEnum.PAWN,
    'N': PieceEnum.KNIGHT,
    'B': PieceEnum.BISHOP,
    'R': PieceEnum.ROOK,
    'Q': PieceEnum.QUEEN,
    'K': PieceEnum.KING
}


class MoveParser:
    def convert_san_to_move_detail(self, prev_fen, san):
        position = get_position_from_fen(prev_fen)
        if san == 'O-O':
            if position.side_to_move == PieceColor.WHITE:
                return {
                    'piece': PieceEnum.KING,
                    'color': PieceColor.WHITE,
                    'square': 60,
                    'move': 62
                }
            else:
                return {
                    'piece': PieceEnum.KING,
                    'color': PieceColor.BLACK,
                    'square': 4,
                    'move': 6
                }

        elif san == 'O-O-O':
            if position.side_to_move == PieceColor.WHITE:
                return {
                    'piece': PieceEnum.KING,
                    'color': PieceColor.WHITE,
                    'square': 60,
                    'move': 58
                }
            else:
                return {
                    'piece': PieceEnum.KING,
                    'color': PieceColor.BLACK,
                    'square': 4,
                    'move': 2
                }

        else:
            parsed = self.parse_san(san)

            piece_char = parsed.get('piece', 'P')
            if piece_char is None:
                piece_char = 'P'
            piece_enum = self.get_piece_enum(piece_char)


            color = position.side_to_move

            to_file = parsed['to_file']
            to_rank = parsed['to_rank']
            to_square = self.convert_file_rank_to_square(to_file, to_rank)

            from_file = parsed.get('from_file')
            from_rank = parsed.get('from_rank')
            from_square = self.convert_file_rank_to_square(from_file, from_rank)

            move_detail = {
                'piece': piece_enum,
                'color': color,
                'square': from_square,
                'move': to_square
            }

        return move_detail

    def parse_san(self, san):
        san_pattern = re.compile(
            r'^(?P<piece>[KQRBN])?(?P<from_file>[a-h])?(?P<from_rank>[1-8])?(?P<capture>x)?(?P<to_file>[a-h])(?P<to_rank>[1-8])(?P<promotion>=Q|=R|=B|=N)?(?P<check>[+#])?$')
        match = san_pattern.match(san)
        if not match:
            raise ValueError(f"Invalid SAN move: {san}")
        return match.groupdict()

    def convert_file_rank_to_square(self, file, rank):
        file_index = ord(file) - ord('a')
        rank_index = 8 - int(rank)
        return rank_index * 8 + file_index

    def get_piece_enum(self, piece_char):
        return piece_map[piece_char]

    def get_file_from_square(self, square):
        return chr((square % 8) + ord('a'))

    def get_rank_from_square(self, square):
        return str((square // 8) + 1)