package coresrv

import (
	"strings"

	"github.com/accreativesoft/go-core/coredao"
	"github.com/accreativesoft/go-core/coredto"
	"github.com/accreativesoft/go-core/coreerror"
	"github.com/accreativesoft/go-core/coremsg"
	"github.com/accreativesoft/go-core/corereflect"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service interface {
	Iniciar(entidadRef interface{}) error
	Crear(entidadRef interface{}) error
	Insertar(entidadRef interface{}) error
	Actualizar(entidadRef interface{}) error
	Eliminar(entidadRef interface{}) error
	Guardar(entidadRef interface{}) error
	ActualizarLista(update coredto.Update) error
	EliminarLista(delete coredto.Delete) error
	NumeroRegistros(filtros []coredto.Filtro) (int, error)
	BuscarPorId(entidadRef interface{}) error
	CargarDetalle(entidadRef interface{}) error
	GetEntidad(entidadRef interface{}, query coredto.Query) error
	GetEntidadList(listaRef interface{}, query coredto.Query) error
	GetObjetoList(listaRef *[]interface{}, query coredto.Query) error
	GetObjeto(listaRef *[]interface{}, query coredto.Query) error
}

type ServiceImpl struct {
	Dao coredao.Dao
}

func NewService(trn *gorm.DB, entidadRef interface{}) *ServiceImpl {
	var dao coredao.Dao = coredao.NewDao(trn, entidadRef)
	return &ServiceImpl{Dao: dao}
}

func (service *ServiceImpl) Iniciar(entidadRef interface{}) error {
	return nil
}

func (service *ServiceImpl) Crear(entidadRef interface{}) error {
	//v := reflect.ValueOf(entidadRef).Elem()
	//v.Set(reflect.Zero(v.Type()))
	ok, e := corereflect.HasField(entidadRef, "Activo")
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}

	datatype, e := corereflect.GetFieldType(entidadRef, "Activo")
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}

	if ok && strings.Compare(datatype, "bool") == 0 {
		e := corereflect.SetField(entidadRef, "Activo", true)
		if e != nil {
			log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
			return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
		}
	}
	return nil
}

func (service *ServiceImpl) Insertar(entidadRef interface{}) error {
	return service.Dao.Insertar(entidadRef)
}

func (service *ServiceImpl) Actualizar(entidadRef interface{}) error {
	return service.Dao.Actualizar(entidadRef)
}

func (service *ServiceImpl) Eliminar(entidadRef interface{}) error {
	return service.Dao.Eliminar(entidadRef)
}

func (service *ServiceImpl) Guardar(entidadRef interface{}) error {
	return service.Dao.Guardar(entidadRef)
}

func (service *ServiceImpl) ActualizarLista(update coredto.Update) error {
	return service.Dao.ActualizarLista(service.Dao.GetEntidadRef(), update)
}

func (service *ServiceImpl) EliminarLista(delete coredto.Delete) error {
	return service.Dao.EliminarLista(service.Dao.GetEntidadRef(), delete)
}

func (service *ServiceImpl) NumeroRegistros(filtros []coredto.Filtro) (int, error) {
	return service.Dao.NumeroRegistros(filtros)
}

func (service *ServiceImpl) BuscarPorId(entidadRef interface{}) error {
	return service.Dao.BuscarPorId(entidadRef)
}

func (service *ServiceImpl) CargarDetalle(entidadRef interface{}) error {
	return service.Dao.CargarDetalle(entidadRef)
}

func (service *ServiceImpl) GetEntidad(entidadRef interface{}, query coredto.Query) error {
	return service.Dao.GetEntidad(entidadRef, query)
}

func (service *ServiceImpl) GetEntidadList(listaRef interface{}, query coredto.Query) error {
	return service.Dao.GetEntidadList(service.Dao.GetEntidadRef(), query, listaRef)
}

func (service *ServiceImpl) GetObjetoList(listaRef *[]interface{}, query coredto.Query) error {
	return service.Dao.GetObjetoList(service.Dao.GetEntidadRef(), query, listaRef)
}

func (service *ServiceImpl) GetObjeto(listaRef *[]interface{}, query coredto.Query) error {
	return service.Dao.GetObjeto(service.Dao.GetEntidadRef(), query, listaRef)
}
