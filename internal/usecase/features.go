package usecase

import (
	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

// ExtractTransformerFeatures populates a pre-allocated slice of 1280 floats (64 * 20)
func ExtractTransformerFeatures(pos *engine.Position, out []float32) {
	_ = out[1279] // ensures out has at least a length of 1280, otherwise panics

	mirror := pos.SideToMove == engine.Black

	epSq := int(pos.EnPassantSquare)
	if epSq != int(engine.NoSquare) && mirror {
		epSq = (7-epSq/8)*8 + epSq%8
	}

	wsc := float32(0)
	wlc := float32(0)
	bsc := float32(0)
	blc := float32(0)

	if pos.CastlingRights.Has(engine.WhiteShort) {
		wsc = 1.0
	}
	if pos.CastlingRights.Has(engine.WhiteLong) {
		wlc = 1.0
	}
	if pos.CastlingRights.Has(engine.BlackShort) {
		bsc = 1.0
	}
	if pos.CastlingRights.Has(engine.BlackLong) {
		blc = 1.0
	}

	if mirror {
		wsc, wlc, bsc, blc = bsc, blc, wsc, wlc
	}

	halfMoves := uint64(pos.HalfMoves)

	// Count position repetition
	repCount := 0
	for i := 0; i < pos.HalfMoves; i++ {
		if pos.RepetitionHistory[i] == pos.Hash {
			repCount++
		}
	}
	f12 := float32(0.0)
	f13 := float32(0.0)
	if repCount == 1 {
		f12 = 1.0
	} else if repCount >= 2 {
		f13 = 1.0
	}

	for sq := 0; sq < 64; sq++ {
		srcSq := sq
		if mirror {
			srcSq = (7-sq/8)*8 + sq%8
		}

		offset := sq * 20

		// features 0-11: pieces
		for bi := 0; bi < 12; bi++ {
			srcBi := bi
			if mirror {
				srcBi = (bi + 6) % 12
			}

			bb := uint64(pos.Pieces[srcBi])
			if ((bb >> srcSq) & 1) != 0 {
				out[offset+bi] = 1.0
			} else {
				out[offset+bi] = 0.0
			}
		}

		// features 12-13: position repetition
		out[offset+12] = f12
		out[offset+13] = f13

		// features 14-17: castling
		out[offset+14] = wsc
		out[offset+15] = wlc
		out[offset+16] = bsc
		out[offset+17] = blc

		// feature 18: en passant
		if epSq != int(engine.NoSquare) && sq == epSq {
			out[offset+18] = 1.0
		} else {
			out[offset+18] = 0.0
		}

		// feature 19: halfmove binary bit
		if ((halfMoves >> (63 - sq)) & 1) != 0 {
			out[offset+19] = 1.0
		} else {
			out[offset+19] = 0.0
		}
	}
}
