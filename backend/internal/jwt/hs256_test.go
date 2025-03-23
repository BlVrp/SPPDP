package jwt_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"one-help/internal/jwt"
)

func TestHS256(t *testing.T) {
	tokenizer := jwt.NewHS256[TestPayload]([]byte("some-secret"))

	t.Run("create token", func(t *testing.T) {
		expected := &jwt.Token[TestPayload]{
			Header:    jwt.Header{Alg: "HS256", Typ: "JWT"},
			Payload:   TestPayload{Value: "test"},
			Signature: []byte{0x5, 0x21, 0x59, 0x41, 0x99, 0x2c, 0xa1, 0x65, 0x99, 0x42, 0x2f, 0x1e, 0xa7, 0xc3, 0xb9, 0x6e, 0xca, 0x83, 0x79, 0x55, 0x5a, 0x2f, 0x13, 0x19, 0x11, 0x61, 0xff, 0x18, 0xf, 0xda, 0xd3, 0xe8},
		}

		token, err := tokenizer.Token(TestPayload{Value: "test"})
		require.NoError(t, err)
		require.EqualValues(t, expected, token)
	})

	t.Run("verify token", func(t *testing.T) {
		token, err := tokenizer.Token(TestPayload{Value: "test"})
		require.NoError(t, err)
		require.NoError(t, tokenizer.Verify(token))

		// INFO: Change signature.
		token.Signature = []byte{}
		require.Error(t, tokenizer.Verify(token))
	})
}
