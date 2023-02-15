package coremodel

import (
	"time"
)

type Entidad struct {
	Id                  int        `gorm:"primary_key;auto_increment;not_null" json:"id,omitempty"`
	FechaCreacion       *time.Time `json:"fechaCreacion,omitempty"`
	FechaModificacion   *time.Time `json:"fechaModificacion,omitempty"`
	UsuarioCreacion     string     `gorm:"type:varchar(50)" json:"usuarioCreacion,omitempty"`
	UsuarioModificacion string     `gorm:"type:varchar(50)" json:"usuarioModificacion,omitempty"`
	IdEntidadOrigen     *int       `json:"idEntidadOrigen,omitempty"`
	Eliminado           bool       `gorm:"-:all" json:"eliminado"`
}
