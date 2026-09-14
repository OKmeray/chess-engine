package engine

import "testing"

const shortTestingMaxDepth = 4

// Perft recursively generates and counts all leaf nodes up to a certain depth.
func Perft(pos *Position, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	nodes := uint64(0)
	moves := make([]Move, 0, 256)
	moves = pos.GenerateMoves(moves)

	for _, move := range moves {
		undo := pos.MakeMove(move)
		nodes += Perft(pos, depth-1)
		pos.UnmakeMove(move, undo)
	}

	return nodes
}

func TestPerftPositions(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		depth    int
		expected uint64
	}{
		// Position 1 (Starting Position)
		{"Pos1_Depth1", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 1, 20},
		{"Pos1_Depth2", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 2, 400},
		{"Pos1_Depth3", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 3, 8902},
		{"Pos1_Depth4", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 4, 197_281},
		{"Pos1_Depth5", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 5, 4_865_609},

		// Position 2 (Kiwipete)
		{"Pos2_Depth1", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 1, 48},
		{"Pos2_Depth2", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 2, 2039},
		{"Pos2_Depth3", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 3, 97862},
		{"Pos2_Depth4", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 4, 4_085_603},
		{"Pos2_Depth5", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 5, 193_690_690},

		// Position 3
		{"Pos3_Depth1", "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1", 1, 14},
		{"Pos3_Depth2", "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1", 2, 191},
		{"Pos3_Depth3", "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1", 3, 2812},
		{"Pos3_Depth4", "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1", 4, 43238},
		{"Pos3_Depth5", "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1", 5, 674_624},

		// Position 4
		{"Pos4_Depth1", "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", 1, 6},
		{"Pos4_Depth2", "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", 2, 264},
		{"Pos4_Depth3", "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", 3, 9467},
		{"Pos4_Depth4", "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", 4, 422333},
		{"Pos4_Depth5", "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", 5, 15_833_292},

		// Position 5
		{"Pos5_Depth1", "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8", 1, 44},
		{"Pos5_Depth2", "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8", 2, 1486},
		{"Pos5_Depth3", "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8", 3, 62379},
		{"Pos5_Depth4", "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8", 4, 2_103_487},
		{"Pos5_Depth5", "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8", 5, 89_941_194},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if testing.Short() && tt.depth > shortTestingMaxDepth {
				t.Skip("Skipping high depth perft in short mode")
			}

			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}

			// Ensure ByColor is fully populated from Pieces bitboards
			for pt := Pawn; pt <= King; pt++ {
				pos.ByColor[1] |= pos.Pieces[int(pt)+int(White)]
				pos.ByColor[0] |= pos.Pieces[int(pt)+int(Black)]
			}

			nodes := Perft(pos, tt.depth)
			if nodes != tt.expected {
				t.Errorf("Perft() = %d, want %d", nodes, tt.expected)
			}
		})
	}
}
