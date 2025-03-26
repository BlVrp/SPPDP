package jwt_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"one-help/internal/jwt"
)

type TestPayload struct {
	Value string `json:"value"`
}

func TestToken(t *testing.T) {
	header := jwt.Header{Alg: "alg", Typ: "jwt"}
	tokenValue := TestPayload{Value: "hello-test"}

	t.Run("to string", func(t *testing.T) {
		t.Run("without signature", func(t *testing.T) {
			expected := "eyJhbGciOiJhbGciLCJ0eXAiOiJqd3QifQ==.eyJ2YWx1ZSI6ImhlbGxvLXRlc3QifQ=="
			token := &jwt.Token[TestPayload]{
				Header:    header,
				Payload:   tokenValue,
				Signature: nil,
			}

			encoded, err := token.String()
			require.NoError(t, err)
			require.EqualValues(t, expected, encoded)

			encodedPreSign, err := token.PreSignString()
			require.NoError(t, err)
			require.EqualValues(t, expected, encodedPreSign)
		})

		t.Run("with signature", func(t *testing.T) {
			expected := "eyJhbGciOiJhbGciLCJ0eXAiOiJqd3QifQ==.eyJ2YWx1ZSI6ImhlbGxvLXRlc3QifQ==.c29tZS1zaWduYXR1cmUtZXhhbXBsZQ=="
			token := &jwt.Token[TestPayload]{
				Header:    header,
				Payload:   tokenValue,
				Signature: []byte("some-signature-example"),
			}

			encoded, err := token.String()
			require.NoError(t, err)
			require.EqualValues(t, expected, encoded)

			encodedPreSign, err := token.PreSignString()
			require.NoError(t, err)
			require.Contains(t, expected, encodedPreSign)
		})
	})

	t.Run("parse", func(t *testing.T) {
		t.Run("without signature", func(t *testing.T) {
			encoded := "eyJhbGciOiJhbGciLCJ0eXAiOiJqd3QifQ==.eyJ2YWx1ZSI6ImhlbGxvLXRlc3QifQ=="
			expected := &jwt.Token[TestPayload]{
				Header:    header,
				Payload:   tokenValue,
				Signature: nil,
			}

			token := new(jwt.Token[TestPayload])
			err := token.Parse(encoded)
			require.NoError(t, err)
			require.EqualValues(t, expected, token)
		})

		t.Run("with signature", func(t *testing.T) {
			encoded := "eyJhbGciOiJhbGciLCJ0eXAiOiJqd3QifQ==.eyJ2YWx1ZSI6ImhlbGxvLXRlc3QifQ==.c29tZS1zaWduYXR1cmUtZXhhbXBsZQ=="
			expected := &jwt.Token[TestPayload]{
				Header:    header,
				Payload:   tokenValue,
				Signature: []byte("some-signature-example"),
			}

			token := new(jwt.Token[TestPayload])
			err := token.Parse(encoded)
			require.NoError(t, err)
			require.EqualValues(t, expected, token)
		})
	})
}
