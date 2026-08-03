package engine

import (
	"testing"
)

func TestCastlingRights_Has(t *testing.T) {
	tests := []struct {
		name   string
		rights CastlingRights
		check  CastlingRights
		want   bool
	}{
		{
			name:   "Has WhiteShort",
			rights: WhiteShort | BlackLong,
			check:  WhiteShort,
			want:   true,
		},
		{
			name:   "Has BlackLong",
			rights: WhiteShort | BlackLong,
			check:  BlackLong,
			want:   true,
		},
		{
			name:   "Does not have WhiteLong",
			rights: WhiteShort | BlackLong,
			check:  WhiteLong,
			want:   false,
		},
		{
			name:   "Does not have BlackShort",
			rights: WhiteShort | BlackLong,
			check:  BlackShort,
			want:   false,
		},
		{
			name:   "Empty rights has nothing",
			rights: 0,
			check:  WhiteShort,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rights.Has(tt.check); got != tt.want {
				t.Errorf("CastlingRights(%v).Has(%v) = %v, want %v", tt.rights, tt.check, got, tt.want)
			}
		})
	}
}

func TestCastlingRights_Remove(t *testing.T) {
	tests := []struct {
		name       string
		initial    CastlingRights
		remove     CastlingRights
		wantRights CastlingRights
	}{
		{
			name:       "Remove right from available rights",
			initial:    WhiteShort | WhiteLong | BlackShort | BlackLong,
			remove:     WhiteShort,
			wantRights: WhiteLong | BlackShort | BlackLong,
		},
		{
			name:       "Remove multiple rights",
			initial:    WhiteShort | WhiteLong | BlackShort | BlackLong,
			remove:     BlackShort | BlackLong,
			wantRights: WhiteShort | WhiteLong,
		},
		{
			name:       "Remove right that doesn't exist",
			initial:    WhiteShort | WhiteLong,
			remove:     BlackShort,
			wantRights: WhiteShort | WhiteLong,
		},
		{
			name:       "Remove all remaining rights",
			initial:    WhiteShort | BlackLong,
			remove:     WhiteShort | BlackLong,
			wantRights: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rights := tt.initial
			rights.Remove(tt.remove)
			if rights != tt.wantRights {
				t.Errorf("CastlingRights(%v).Remove(%v) = %v, want %v", tt.initial, tt.remove, rights, tt.wantRights)
			}
		})
	}
}
