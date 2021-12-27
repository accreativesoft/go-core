package coremodel

import (
	"time"

	"gorm.io/gorm"
)

type Entidad struct {
	Id                  uint      `gorm:"primary_key;auto_increment;not_null" json:"id,omitempty"`
	FechaCreacion       time.Time `json:"fechaCreacion,omitempty"`
	FechaModificacion   time.Time `json:"fechaModificacion,omitempty"`
	UsuarioCreacion     string    `gorm:"type:varchar(50)" json:"usuarioCreacion,omitempty"`
	UsuarioModificacion string    `gorm:"type:varchar(50)" json:"usuarioModificacion,omitempty"`
	IdEntidadOrigen     uint      `json:"idEntidadOrigen,omitempty"`
	Eliminado           bool      `gorm:"not null" json:"eliminado,omitempty"`
}

func (entidad *Entidad) Guardar(tx *gorm.DB, T interface{}) interface{} {
	return tx.Create(T)
}
