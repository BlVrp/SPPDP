package jwt

import (
	"github.com/zeebo/errs"
)

// Error defines wrapper for errors produced by jwt package.
var Error = errs.Class("jwt")

// Tokenizer holds JWT creation/verification operations.
// INFO: PT - payload type. Must a flat struct type with all public fields.
type Tokenizer[PT any] interface {
	// Token creates a Token from provided payload.
	Token(PT) (*Token[PT], error)
	// Sign signs token payload, updated provided token signature field.
	Sign(*Token[PT]) error
	// Verify returns error if signature value is not derived from payload data.
	Verify(*Token[PT]) error
}
