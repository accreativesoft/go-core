package coredto

import (
	"time"
)

type EntidadDto struct {
	Id                  int        `json:"id,omitempty"`
	FechaCreacion       *time.Time `json:"fechaCreacion,omitempty"`
	FechaModificacion   *time.Time `json:"fechaModificacion,omitempty"`
	UsuarioCreacion     string     `json:"usuarioCreacion,omitempty"`
	UsuarioModificacion string     `json:"usuarioModificacion,omitempty"`
	IdEntidadOrigen     *int       `json:"idEntidadOrigen,omitempty"`
	Eliminado           bool       `json:"eliminado"`
}
