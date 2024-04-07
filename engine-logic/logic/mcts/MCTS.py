import math
import random

from logic.engine.Position import Position


class MCTSNode:
    def __init__(self, position: Position, move=None, parent=None):
        self.position = position
        self.move = move
        self.parent = parent
        self.children = []
        self.wins = 0
        self.visits = 0

    def is_leaf(self):
        return len(self.children) == 0

    def select_child(self):
        # if there are unvisited nodes, then return one of them
        unvisited = [child for child in self.children if child.visits == 0]
        if unvisited:
            return random.choice(unvisited)

        # Calculation of UCB1 for visited nodes
        C = math.sqrt(2)
        highest_UCB1 = float('-inf')
        node_with_highest_UCB1 = None
        for child in self.children:
            UCB1 = child.wins / child.visits + C * math.sqrt(math.log(child.parent.visits) / child.visits)
            if UCB1 > highest_UCB1:
                highest_UCB1 = UCB1
                node_with_highest_UCB1 = child
        return node_with_highest_UCB1

    def expand(self):
        possible_moves = self.position.get_all_separate_moves()

        for move in possible_moves:
            new_position = self.position.clone()
            new_position.apply_move(move)
            self.children.append(MCTSNode(new_position, move, self))

    def simulate(self):
        temp_position = self.position.clone()

        while not temp_position.is_game_over():
            possible_moves = temp_position.get_all_separate_moves()
            move = random.choice(possible_moves)
            temp_position.apply_move(move)

        return temp_position.determine_outcome()

    def backpropagate(self, result):
        # Update this node's stats and propagate the result up the tree
        self.visits += 1
        if result == "win":
            self.wins += 1

class MCTS:
    def __init__(self, root_position):
        self.root = MCTSNode(root_position)

    def run(self, iterations):
        for _ in range(iterations):
            node = self.root
            # Selection
            while not node.is_leaf():
                node = node.select_child()
            # Expansion
            if node.visits == 0:
                result = node.simulate()
            else:
                node.expand()
                node = node.select_child()
                result = node.simulate()
            # Backpropagation
            node.backpropagate(result)

    def best_move(self):
        # Return the move from root with the highest win rate
        highest_UCB1 = float('-inf')
        highest_UCB1_node = None
        for i in range(0, len(self.root.children)):
            if self.root.children[i].visits != 0:
                UCB1 = self.root.children[i].wins / self.root.children[i].visits + math.sqrt(
                    2 * math.log(self.root.children[i].parent.visits) / self.root.children[i].visits)
                if UCB1 > highest_UCB1:
                    highest_UCB1_node = self.root.children[i]
        return highest_UCB1_node

