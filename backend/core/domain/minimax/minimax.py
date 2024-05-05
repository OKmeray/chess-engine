from core.domain.engine.Position import Position
from core.domain.evaluation.evaluation import evaluate_position


def minimax(position: Position, alpha, beta, depth=5, maximizing_player=True, is_root=True):
    if position.is_game_over():

        outcome = position.determine_outcome()
        if outcome == "win":
            return (float('inf'), None) if maximizing_player else (-float('inf'), None)

        elif outcome == "loss":
            return (-float('inf'), None) if maximizing_player else (float('inf'), None)

        else:
            return 0, None  # Assuming a draw

    if depth == 0:
        return evaluate_position(position, position.side_to_move), None

    best_move = None

    if maximizing_player:
        max_eval = -float('inf')
        all_moves = position.get_moves_sorted_by_priority()

        for move in all_moves:
            position_copy = position.clone()
            position_copy.apply_move(move)

            eval, _ = minimax(position_copy, alpha, beta, depth - 1, False, False)

            if eval > max_eval:
                max_eval = eval
                best_move = move if is_root else None  # Store the move only on the root call

            alpha = max(alpha, eval)
            if beta <= alpha:
                break
        return max_eval, best_move
    else:
        min_eval = float('inf')
        all_moves = position.get_moves_sorted_by_priority()

        for move in all_moves:
            position_copy = position.clone()
            position_copy.apply_move(move)

            eval, _ = minimax(position_copy, alpha, beta, depth - 1, True, False)

            if eval < min_eval:
                min_eval = eval
                best_move = move if is_root else None  # Store the move only on the root call

            beta = min(beta, eval)
            if beta <= alpha:
                break
        return min_eval, best_move
