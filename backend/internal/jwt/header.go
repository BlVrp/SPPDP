package jwt

// Header holds JWT header data.
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}
