package usecase

import (
	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

// TypesCount is the number of distinct move planes in the AlphaZero policy head.
const (
	TypesCount = 76
	PolicySize = 64 * TypesCount // 4864
)

var (
	typeToIndex map[string]int
)

func init() {
	typeToIndex = make(map[string]int)
	idx := 0

	// sliding
	for _, dir := range []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"} {
		for s := 1; s <= 7; s++ {
			key := makeSlideKey(dir, s)
			typeToIndex[key] = idx
			idx++
		}
	}

	// knights
	knightDirs := []struct{ d2, d1 string }{
		{"N", "E"}, {"N", "W"}, {"S", "E"}, {"S", "W"},
		{"E", "N"}, {"E", "S"}, {"W", "N"}, {"W", "S"},
	}
	for _, k := range knightDirs {
		key := makeKnightKey(k.d2, k.d1[0])
		typeToIndex[key] = idx
		idx++
	}

	// promotions
	for _, dir := range []string{"NW", "N", "NE"} {
		for _, promo := range []byte{'Q', 'R', 'B', 'N'} {
			key := makePromoKey(dir, promo)
			typeToIndex[key] = idx
			idx++
		}
	}
}

func makeSlideKey(dir string, steps int) string {
	return "slide|" + dir + "|" + string(rune(steps+'0')) + "|\x00"
}

func makeKnightKey(dir string, p byte) string {
	return "knight|" + dir + "|0|" + string([]byte{p})
}

func makePromoKey(dir string, p byte) string {
	return "promo|" + dir + "|0|" + string([]byte{p})
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// EncodeMove converts an engine.Move into a unique index for policy vector.
func EncodeMove(m engine.Move, mirror bool) int {
	fromSq := int(m.From())
	toSq := int(m.To())

	if mirror {
		fromSq = (7-(fromSq/8))*8 + (fromSq % 8)
		toSq = (7-(toSq/8))*8 + (toSq % 8)
	}

	fx, fy := fromSq%8, fromSq/8
	tx, ty := toSq%8, toSq/8
	dx, dy := tx-fx, ty-fy

	var key string

	// Check Promotion
	if m.Flags() == engine.FlagPromoQueen || m.Flags() == engine.FlagPromoRook ||
		m.Flags() == engine.FlagPromoBishop || m.Flags() == engine.FlagPromoKnight ||
		m.Flags() == engine.FlagPromoQueenCapture || m.Flags() == engine.FlagPromoRookCapture ||
		m.Flags() == engine.FlagPromoBishopCapture || m.Flags() == engine.FlagPromoKnightCapture {

		dir := ""
		if dy > 0 {
			dir += "N"
		}
		if dy < 0 {
			dir += "S"
		}
		if dx > 0 {
			dir += "E"
		}
		if dx < 0 {
			dir += "W"
		}

		var promoByte byte
		switch m.Flags() {
		case engine.FlagPromoQueen, engine.FlagPromoQueenCapture:
			promoByte = 'Q'
		case engine.FlagPromoRook, engine.FlagPromoRookCapture:
			promoByte = 'R'
		case engine.FlagPromoBishop, engine.FlagPromoBishopCapture:
			promoByte = 'B'
		case engine.FlagPromoKnight, engine.FlagPromoKnightCapture:
			promoByte = 'N'
		}
		key = makePromoKey(dir, promoByte)
	} else if (abs(dx) == 2 && abs(dy) == 1) || (abs(dx) == 1 && abs(dy) == 2) {
		// Knight
		big := ""
		if abs(dx) == 2 {
			if dx > 0 {
				big = "E"
			} else {
				big = "W"
			}
		} else {
			if dy > 0 {
				big = "N"
			} else {
				big = "S"
			}
		}

		small := ""
		if big == "N" || big == "S" {
			if dx > 0 {
				small = "E"
			} else {
				small = "W"
			}
		} else {
			if dy > 0 {
				small = "N"
			} else {
				small = "S"
			}
		}
		key = makeKnightKey(big, small[0])
	} else {
		// Slide / Pawn Push / King
		dir := ""
		if dy > 0 {
			dir += "N"
		}
		if dy < 0 {
			dir += "S"
		}
		if dx > 0 {
			dir += "E"
		}
		if dx < 0 {
			dir += "W"
		}
		steps := max(abs(dx), abs(dy))

		key = makeSlideKey(dir, steps)
	}

	t := typeToIndex[key]
	return fromSq*TypesCount + t
}
