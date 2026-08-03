package engine

// CastlingRights is an 8-bit mask holding all 4 castling states
type CastlingRights uint8

const (
	WhiteShort CastlingRights = 1 << 0 // 0001
	WhiteLong  CastlingRights = 1 << 1 // 0010
	BlackShort CastlingRights = 1 << 2 // 0100
	BlackLong  CastlingRights = 1 << 3 // 1000
)

// Has returns true if the specific castling right is still available
func (c CastlingRights) Has(right CastlingRights) bool {
	return (c & right) != 0
}

// Remove clears the specific castling rights
func (c *CastlingRights) Remove(rights CastlingRights) {
	*c &^= rights
}
