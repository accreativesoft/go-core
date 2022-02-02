package coremodel

import (
	"reflect"
	"time"

	"github.com/accreativesoft/go-core/coredto"
	"github.com/accreativesoft/go-core/coreerror"
	"github.com/accreativesoft/go-core/coremsg"
	"github.com/accreativesoft/go-core/corereflect"
	"github.com/accreativesoft/go-core/coresql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Entidad struct {
	Id                  uint       `gorm:"primary_key;auto_increment;not_null" json:"id,omitempty"`
	FechaCreacion       *time.Time `json:"fechaCreacion,omitempty"`
	FechaModificacion   *time.Time `json:"fechaModificacion,omitempty"`
	UsuarioCreacion     string     `gorm:"type:varchar(50)" json:"usuarioCreacion,omitempty"`
	UsuarioModificacion string     `gorm:"type:varchar(50)" json:"usuarioModificacion,omitempty"`
	IdEntidadOrigen     *uint      `json:"idEntidadOrigen,omitempty"`
	Eliminado           bool       `gorm:"not null" json:"eliminado,omitempty"`
}

func (entidad *Entidad) Insertar(trn *gorm.DB, entidadRef interface{}) error {
	e := trn.Omit(clause.Associations).Create(entidadRef).Error
	if e != nil {
		return coreerror.NewError(coremsg.MSG_ERROR_SQL, e.Error())
	}
	return nil
}

func (entidad *Entidad) Actualizar(trn *gorm.DB, entidadRef interface{}) error {
	e := trn.Omit(clause.Associations).Save(entidadRef).Error
	if e != nil {
		return coreerror.NewError(coremsg.MSG_ERROR_SQL, e.Error())
	}
	return nil
}

func (entidad *Entidad) Eliminar(trn *gorm.DB, entidadRef interface{}) error {
	e := trn.Omit(clause.Associations).Delete(entidadRef).Error
	if e != nil {
		return coreerror.NewError(coremsg.MSG_ERROR_SQL, e.Error())
	}
	return nil
}

func (entidad *Entidad) Guardar(trn *gorm.DB, entidadRef interface{}) error {

	//Recupero mi campo entidad de la entidadRef
	ent := reflect.ValueOf(entidadRef).Elem().FieldByName("Entidad")

	//Si es eliminado guarda la informacion
	if ent.FieldByName("Eliminado").Bool() {
		if ent.FieldByName("Id").Int() != 0 {
			//Si el id es diferente de 0 elimina el registro
			//entidad.Eliminar(trn, entidadRef)
			ref, _ := corereflect.GetField(entidadRef, "Id")
			e := trn.Omit(clause.Associations).Delete(entidadRef, ref).Error
			if e != nil {
				return coreerror.NewError(coremsg.MSG_ERROR_SQL, e.Error())
			}
		}
	} else {
		if ent.FieldByName("Id").Uint() == 0 {
			//Si el id es igual a cero inserta el registro
			//entidad.Insertar(trn, entidadRef)
			e := trn.Omit(clause.Associations).Create(entidadRef).Error
			if e != nil {
				return coreerror.NewError(coremsg.MSG_ERROR_SQL, e.Error())
			}
		} else {
			//Si el id es diferente a cero actualiza el registro
			//entidad.Actualizar(trn, entidadRef)
			e := trn.Omit(clause.Associations).Save(entidadRef).Error
			if e != nil {
				return coreerror.NewError(coremsg.MSG_ERROR_SQL, e.Error())
			}
		}
	}
	return nil
}

func (entidad *Entidad) ActualizarLista(trn *gorm.DB, entidadRef interface{}, update coredto.Update) error {
	return coresql.ActualizarLista(trn, entidadRef, update)
}

func (entidad *Entidad) EliminarLista(trn *gorm.DB, entidadRef interface{}, delete coredto.Delete) error {
	return coresql.EliminarLista(trn, entidadRef, delete)
}

func (entidad *Entidad) NumeroRegistros(trn *gorm.DB, entidadRef interface{}, filtros []coredto.Filtro) (int, error) {
	return coresql.NumeroRegistros(trn, entidadRef, filtros)
}

func (entidad *Entidad) BuscarPorId(trn *gorm.DB, entidadRef interface{}) error {
	id, _ := corereflect.GetField(entidadRef, "Id")
	e := trn.First(entidadRef, id).Error
	if e != nil {
		return coreerror.NewError(coremsg.MSG_ERROR_SQL, e.Error())
	}
	return nil
}

func (entidad *Entidad) Get(trn *gorm.DB, entidadRef interface{}, query coredto.Query) error {
	return coresql.Get(trn, entidadRef, query)
}

func (entidad *Entidad) GetLista(trn *gorm.DB, entidadRef interface{}, query coredto.Query, listaRef interface{}) error {
	return coresql.GetLista(trn, entidadRef, query, listaRef)
}
