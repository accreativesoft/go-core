package coredto

type Autenticacion struct {
	Id            uint   `json:"id,omitempty"`
	NombreUsuario string `json:"nombreUsuario,omitempty"`
	Email         string `json:"email,omitempty"`
	Clave         string `json:"clave,omitempty"`
}
