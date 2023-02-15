package coredao

import (
	"reflect"

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

type Dao interface {
	Insertar(entidadRef interface{}) error
	Actualizar(entidadRef interface{}) error
	Eliminar(entidadRef interface{}) error
	Guardar(entidadRef interface{}) error
	ActualizarLista(entidadRef interface{}, update coredto.Update) error
	EliminarLista(entidadRef interface{}, delete coredto.Delete) error
	NumeroRegistros(filtros []coredto.Filtro) (int, error)
	BuscarPorId(entidadRef interface{}) error
	CargarDetalle(entidadRef interface{}) error
	GetEntidad(entidadRef interface{}, query coredto.Query) error
	GetEntidadList(entidadRef interface{}, query coredto.Query, listaRef interface{}) error
	GetObjetoList(entidadRef interface{}, query coredto.Query, listaRef *[]interface{}) error
	GetObjeto(entidadRef interface{}, query coredto.Query, listaRef *[]interface{}) error
	GetEntidadRef() interface{}
	GetTrn() *gorm.DB
}

type DaoImpl struct {
	EntidadRef interface{}
	Trn        *gorm.DB
}

func NewDao(trn *gorm.DB, entidadRef interface{}) *DaoImpl {
	return &DaoImpl{Trn: trn, EntidadRef: entidadRef}
}

func (daoImpl *DaoImpl) GetEntidadRef() interface{} {
	return daoImpl.EntidadRef
}

func (daoImpl *DaoImpl) GetTrn() *gorm.DB {
	return daoImpl.Trn
}

func (daoImpl *DaoImpl) Insertar(entidadRef interface{}) error {
	e := daoImpl.Trn.Create(entidadRef).Error
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}
	return nil
}

func (daoImpl *DaoImpl) Actualizar(entidadRef interface{}) error {
	e := daoImpl.Trn.Save(entidadRef).Error
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}
	return nil
}

func (daoImpl *DaoImpl) Eliminar(entidadRef interface{}) error {
	e := daoImpl.Trn.Omit(clause.Associations).Delete(entidadRef).Error
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}
	return nil
}

func (daoImpl *DaoImpl) Guardar(entidadRef interface{}) error {

	//Recupero mi campo daoImpl de la entidadRef
	ent := reflect.ValueOf(entidadRef).Elem().FieldByName("Entidad")

	//Si es eliminado guarda la informacion
	if ent.FieldByName("Eliminado").Bool() {
		if ent.FieldByName("Id").Int() != 0 {
			//Si el id es diferente de 0 elimina el registro
			//daoImpl.Eliminar(daoImpl.Trn entidadRef)
			ref, _ := corereflect.GetField(entidadRef, "Id")
			e := daoImpl.Trn.Delete(entidadRef, ref).Error
			if e != nil {
				log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
				return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
			}
		}
	} else {
		if ent.FieldByName("Id").Int() == 0 {
			//Si el id es igual a cero inserta el registro
			//daoImpl.Insertar(daoImpl.Trn entidadRef)
			e := daoImpl.Trn.Create(entidadRef).Error
			if e != nil {
				log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
				return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
			}
		} else {
			//Si el id es diferente a cero actualiza el registro
			//daoImpl.Actualizar(daoImpl.Trn entidadRef)
			e := daoImpl.Trn.Save(entidadRef).Error
			if e != nil {
				log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
				return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
			}
		}
	}
	return nil
}

func (daoImpl *DaoImpl) ActualizarLista(entidadRef interface{}, update coredto.Update) error {
	return coresql.ActualizarLista(daoImpl.Trn, daoImpl.EntidadRef, update)
}

func (daoImpl *DaoImpl) EliminarLista(entidadRef interface{}, delete coredto.Delete) error {
	return coresql.EliminarLista(daoImpl.Trn, daoImpl.EntidadRef, delete)
}

func (daoImpl *DaoImpl) NumeroRegistros(filtros []coredto.Filtro) (int, error) {
	return coresql.NumeroRegistros(daoImpl.Trn, daoImpl.EntidadRef, filtros)
}

func (daoImpl *DaoImpl) BuscarPorId(entidadRef interface{}) error {
	id, _ := corereflect.GetField(entidadRef, "Id")
	//e := trn.First(entidadRef, id).Error
	query := coredto.Query{}
	query.AddCampos(corereflect.GetColumns(entidadRef)...)
	query.AddFiltro("id", corecons.EQUALS, id)
	e := coresql.GetEntidad(daoImpl.Trn, entidadRef, query)
	if e != nil {
		v := reflect.ValueOf(entidadRef).Elem()
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	return nil
}

func (daoImpl *DaoImpl) CargarDetalle(entidadRef interface{}) error {
	id, _ := corereflect.GetField(entidadRef, "Id")
	//e := trn.First(entidadRef, id).Error
	query := coredto.Query{}
	query.AddCampos(corereflect.GetColumns(entidadRef)...)
	query.AddFiltro("id", corecons.EQUALS, id)
	e := coresql.GetEntidad(daoImpl.Trn, entidadRef, query)
	if e != nil {
		v := reflect.ValueOf(entidadRef).Elem()
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	return nil
}

func (daoImpl *DaoImpl) GetEntidad(entidadRef interface{}, query coredto.Query) error {
	return coresql.GetEntidad(daoImpl.Trn, entidadRef, query)
}

func (daoImpl *DaoImpl) GetEntidadList(entidadRef interface{}, query coredto.Query, listaRef interface{}) error {
	return coresql.GetEntidadList(daoImpl.Trn, daoImpl.EntidadRef, query, listaRef)
}

func (daoImpl *DaoImpl) GetObjetoList(entidadRef interface{}, query coredto.Query, listaRef *[]interface{}) error {
	return coresql.GetObjetoList(daoImpl.Trn, daoImpl.EntidadRef, query, listaRef)
}

func (daoImpl *DaoImpl) GetObjeto(entidadRef interface{}, query coredto.Query, listaRef *[]interface{}) error {
	return coresql.GetObjeto(daoImpl.Trn, daoImpl.EntidadRef, query, listaRef)
}
