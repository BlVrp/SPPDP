package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
)

var hs256Header = Header{Alg: "HS256", Typ: "JWT"}

// HS256 is a tokenizer over SHA-256 hash function.
type HS256[PT any] struct {
	secret []byte
}

// NewHS256 creates new tokenizer over SHA-256 hash function.
func NewHS256[PT any](secret []byte) Tokenizer[PT] {
	return &HS256[PT]{secret: secret}
}

// Token creates a Token from provided payload.
func (hs *HS256[PT]) Token(payload PT) (*Token[PT], error) {
	token := &Token[PT]{
		Header:  hs256Header,
		Payload: payload,
	}

	return token, hs.Sign(token)
}

// Sign signs token payload, updated provided token signature field.
func (hs *HS256[PT]) Sign(token *Token[PT]) error {
	mac := hmac.New(sha256.New, hs.secret)

	encoded, err := token.PreSignString()
	if err != nil {
		return Error.Wrap(err)
	}

	_, err = mac.Write([]byte(encoded))
	if err != nil {
		return Error.Wrap(err)
	}

	token.Signature = mac.Sum(nil)

	return nil
}

// Verify returns error if signature value is not derived from payload data.
func (hs *HS256[PT]) Verify(token *Token[PT]) error {
	signature := token.Signature
	if err := hs.Sign(token); err != nil {
		return err
	}

	if subtle.ConstantTimeCompare(signature, token.Signature) != 1 {
		return Error.New("invalid signature")
	}

	return nil
}
