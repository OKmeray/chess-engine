from core.domain.engine.Position import Position
from core.domain.engine.enums import PieceColor
from core.domain.evaluation.evaluation import evaluate_position


def minimax(position: Position, alpha, beta, depth=5, branches_to_check=10, is_root=True):
    global visited_nodes
    visited_nodes += 1
    if position.is_game_over():

        outcome = position.determine_outcome()
        if outcome == "win":
            return (float('inf'), None) if position.side_to_move == PieceColor.WHITE else (-float('inf'), None)

        elif outcome == "loss":
            return (-float('inf'), None) if position.side_to_move == PieceColor.WHITE else (float('inf'), None)

        else:
            return 0, None

    if depth == 0:
        return evaluate_position(position), None

    best_move = None

    if position.side_to_move == PieceColor.WHITE:
        max_eval = -float('inf')

        # TODO: when MVVLVA will be ready, then no such thing needed --> [:branches_to_check]
        all_moves = position.get_moves_sorted_by_priority()[:branches_to_check]
        # take :branches_to_check with priority 0 and all the rest with bigger priority

        for move in all_moves:
            position_copy = position.clone()
            position_copy.apply_move(move)

            eval, _ = minimax(position_copy, alpha, beta, depth - 1, branches_to_check-1, is_root=False)
            if depth == 4:
                print(eval, move)
            if eval > max_eval:
                max_eval = eval
                best_move = move if is_root else None  # Store the move only on the root call

            alpha = max(alpha, eval)
            if beta <= alpha:
                break
        return max_eval, best_move
    else:  # position.side_to_move == PieceColor.BLACK
        min_eval = float('inf')

        all_moves = position.get_moves_sorted_by_priority()[:branches_to_check]

        for move in all_moves:
            position_copy = position.clone()
            position_copy.apply_move(move)

            eval, _ = minimax(position_copy, alpha, beta, depth - 1, branches_to_check-1, is_root=False)

            if eval < min_eval:
                min_eval = eval
                best_move = move if is_root else None  # Store the move only on the root call

            beta = min(beta, eval)
            if beta <= alpha:
                break
        return min_eval, best_move

