package coredto

type Autorizacion struct {
	Uri   string `json:"uri,omitempty"`
	Token string `json:"token,omitempty"`
	Exp   int64  `json:"exp,omitempty"`
}
