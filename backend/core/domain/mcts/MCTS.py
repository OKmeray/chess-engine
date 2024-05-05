# import math
# import random
# from logic.engine.Position import Position
#
#
# class MCTSNode:
#     def __init__(self, position: Position, move=None, parent=None):
#         self.position = position
#         self.move = move
#         self.parent = parent
#         self.children = []
#         self.wins = 0
#         self.visits = 0
#
#     def is_leaf(self):
#         return len(self.children) == 0
#
#     def select_child(self):
#         # if there are unvisited nodes, then return one of them
#         unvisited = [child for child in self.children if child.visits == 0]
#         if unvisited:
#             return random.choice(unvisited)
#
#         # Calculation of UCB1 for visited nodes
#         C = math.sqrt(2)
#         highest_UCB1 = float('-inf')
#         node_with_highest_UCB1 = None
#
#         for child in self.children:
#             UCB1 = child.wins / child.visits + C * math.sqrt(math.log(child.parent.visits) / child.visits)
#             if UCB1 > highest_UCB1:
#                 highest_UCB1 = UCB1
#                 node_with_highest_UCB1 = child
#
#         return node_with_highest_UCB1
#
#     def expand(self):
#         possible_moves = self.position.get_all_separate_moves()
#
#         for move in possible_moves:
#             new_position = self.position.clone()
#             new_position.apply_move(move)
#             child_node = MCTSNode(new_position, move, self)
#             self.children.append(child_node)
#
#         if not self.children:
#             print("No children added. Possible issue with move generation or application.")
#
#     def simulate(self):
#         temp_position = self.position.clone()
#
#         while not temp_position.is_game_over():
#             possible_moves = temp_position.get_all_separate_moves()
#             # TODO: choose checks and captures
#             move = random.choice(possible_moves)
#
#             temp_position.apply_move(move)
#
#         return temp_position.determine_outcome()
#
#     def backpropagate(self, result):
#         # Update this node's stats and propagate the result up the tree
#         self.visits += 1
#         if result == "win":
#             self.wins += 1
#
#         # Recursively update all ancestors
#         if self.parent is not None:
#             self.parent.backpropagate(result)
#
# class MCTS:
#     def __init__(self, root_position):
#         self.root = MCTSNode(root_position)
#
#     def run(self, iterations):
#         for _ in range(iterations):
#             node = self.root
#             path = []
#             while not node.is_leaf():
#                 node = node.select_child()
#                 path.append(node)
#
#             if node.visits == 0:
#                 result = node.simulate()
#             else:
#                 node.expand()
#                 if node.children:  # Check if children were successfully added
#                     node = random.choice(node.children)
#                     path.append(node)
#                 result = node.simulate()
#
#             node.backpropagate(result)
#
#     def best_move(self):
#         # Return the move from root with the highest win rate
#         highest_UCB1 = float('-inf')
#         highest_UCB1_node = None
#
#         for i in range(0, len(self.root.children)):
#
#             if self.root.children[i].visits != 0:
#                 UCB1 = self.root.children[i].wins / self.root.children[i].visits + math.sqrt(
#                     math.log(self.root.children[i].parent.visits) / (2 * self.root.children[i].visits))
#                 # print("UCB1", UCB1)
#                 if UCB1 > highest_UCB1:
#                     print(UCB1, ">", highest_UCB1, "\t", self.root.children[i].move)
#                     highest_UCB1 = UCB1
#                     highest_UCB1_node = self.root.children[i]
#
#         return highest_UCB1_node
#

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
        self.ucb1_constant = math.sqrt(2)  # You can tune this value

    def is_leaf(self):
        return len(self.children) == 0

    def select_child(self):
        # if there are unvisited nodes, then return one of them
        unvisited = [child for child in self.children if child.visits == 0]
        if unvisited:
            return random.choice(unvisited)

        return max(self.children, key=lambda child: child.ucb1())

    def ucb1(self):
        if self.visits == 0:
            return float('inf')  # Infinite UCB1 value to ensure all nodes are visited at least once
        exploration = self.ucb1_constant * math.sqrt(math.log(self.parent.visits) / self.visits)
        exploitation = self.wins / self.visits
        return exploitation + exploration

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
        self.visits += 1
        if result == "win":
            self.wins += 1
        if self.parent:
            self.parent.backpropagate(result)

class MCTS:
    def __init__(self, root_position):
        self.root = MCTSNode(root_position)

    def run(self, iterations):
        for _ in range(iterations):
            node = self.root
            while not node.is_leaf():
                node = node.select_child()
            node.expand()
            if node.children:
                node = random.choice(node.children)
            result = node.simulate()
            node.backpropagate(result)

    def best_move(self):
        best_node = max(self.root.children, key=lambda n: n.wins / n.visits if n.visits > 0 else -1)
        return best_node.move if best_node else None
