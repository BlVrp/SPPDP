package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Token holds payload-parametrized JWT token data.
// INFO: PT - payload type. Must a flat struct type with all public fields.
type Token[PT any] struct {
	Header    Header
	Payload   PT
	Signature []byte
}

// String returns string form of the token.
func (t *Token[PT]) String() (string, error) {
	data := make([]string, 0, 2)

	preSign, err := t.PreSignString()
	if err != nil {
		return "", err
	}

	data = append(data, preSign)

	if len(t.Signature) > 0 {
		signature := base64.URLEncoding.EncodeToString(t.Signature)
		data = append(data, signature)
	}

	return strings.Join(data, "."), nil
}

// PreSignString returns string to be signed.
func (t *Token[PT]) PreSignString() (string, error) {
	headerData, err := json.Marshal(t.Header)
	if err != nil {
		return "", Error.Wrap(err)
	}

	payloadData, err := json.Marshal(t.Payload)
	if err != nil {
		return "", Error.Wrap(err)
	}

	header := base64.URLEncoding.EncodeToString(headerData)
	payload := base64.URLEncoding.EncodeToString(payloadData)

	return strings.Join([]string{header, payload}, "."), nil
}

// Parse parses JWT token from string with parametrized payload.
func (t *Token[PT]) Parse(tkn string) error {
	if t == nil {
		return Error.New("receiver is nil")
	}

	parts := strings.Split(tkn, ".")
	if len(parts) != 3 && len(parts) != 2 {
		return Error.New("invalid token format")
	}

	// INFO: header.
	headerData, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return Error.Wrap(err)
	}

	err = json.Unmarshal(headerData, &t.Header)
	if err != nil {
		return Error.Wrap(err)
	}

	// Payload.
	payloadData, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return Error.Wrap(err)
	}

	err = json.Unmarshal(payloadData, &t.Payload)
	if err != nil {
		return Error.Wrap(err)
	}

	if len(parts) == 3 {
		t.Signature, err = base64.URLEncoding.DecodeString(parts[2])
		if err != nil {
			return Error.Wrap(err)
		}
	}

	return nil
}
