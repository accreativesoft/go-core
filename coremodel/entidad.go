package coremodel

import (
	"reflect"
	"time"

	"github.com/accreativesoft/go-core/corecons"
	"github.com/accreativesoft/go-core/coredto"
	"github.com/accreativesoft/go-core/coreerror"
	"github.com/accreativesoft/go-core/coremsg"
	"github.com/accreativesoft/go-core/corereflect"
	"github.com/accreativesoft/go-core/coresql"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (entidad *Entidad) Insertar(trn *gorm.DB, entidadRef interface{}) error {
	e := trn.Create(entidadRef).Error
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}
	return nil
}

func (entidad *Entidad) Actualizar(trn *gorm.DB, entidadRef interface{}) error {
	e := trn.Save(entidadRef).Error
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}
	return nil
}

func (entidad *Entidad) Eliminar(trn *gorm.DB, entidadRef interface{}) error {
	e := trn.Omit(clause.Associations).Delete(entidadRef).Error
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
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
			e := trn.Delete(entidadRef, ref).Error
			if e != nil {
				log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
				return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
			}
		}
	} else {
		if ent.FieldByName("Id").Int() == 0 {
			//Si el id es igual a cero inserta el registro
			//entidad.Insertar(trn, entidadRef)
			e := trn.Create(entidadRef).Error
			if e != nil {
				log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
				return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
			}
		} else {
			//Si el id es diferente a cero actualiza el registro
			//entidad.Actualizar(trn, entidadRef)
			e := trn.Save(entidadRef).Error
			if e != nil {
				log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
				return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
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
	//e := trn.First(entidadRef, id).Error
	query := coredto.Query{}
	query.AddCampos(corereflect.GetColumns(entidadRef)...)
	query.AddFiltro("id", corecons.EQUALS, id)
	e := coresql.GetEntidad(trn, entidadRef, query)
	if e != nil {
		v := reflect.ValueOf(entidadRef).Elem()
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	return nil
}

func (entidad *Entidad) CargarDetalle(trn *gorm.DB, entidadRef interface{}) error {
	id, _ := corereflect.GetField(entidadRef, "Id")
	//e := trn.First(entidadRef, id).Error
	query := coredto.Query{}
	query.AddCampos(corereflect.GetColumns(entidadRef)...)
	query.AddFiltro("id", corecons.EQUALS, id)
	e := coresql.GetEntidad(trn, entidadRef, query)
	if e != nil {
		v := reflect.ValueOf(entidadRef).Elem()
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	return nil
}

func (entidad *Entidad) GetEntidad(trn *gorm.DB, entidadRef interface{}, query coredto.Query) error {
	return coresql.GetEntidad(trn, entidadRef, query)
}

func (entidad *Entidad) GetEntidadList(trn *gorm.DB, entidadRef interface{}, query coredto.Query, listaRef interface{}) error {
	return coresql.GetEntidadList(trn, entidadRef, query, listaRef)
}

func (entidad *Entidad) GetObjetoList(trn *gorm.DB, entidadRef interface{}, query coredto.Query, listaRef *[]interface{}) error {
	return coresql.GetObjetoList(trn, entidadRef, query, listaRef)
}

func (entidad *Entidad) GetObjeto(trn *gorm.DB, entidadRef interface{}, query coredto.Query, listaRef *[]interface{}) error {
	return coresql.GetObjeto(trn, entidadRef, query, listaRef)
}
