package coresrv

import (
	"github.com/accreativesoft/go-core/coredto"
	"github.com/accreativesoft/go-core/corereflect"
	"gorm.io/gorm"
)

type Service struct {
	EntidadRef interface{}
	Trn        *gorm.DB
}

func (service *Service) Iniciar(entidadRef interface{}) error {
	return nil
}

func (service *Service) Crear(entidadRef interface{}) error {
	return nil
}

func (service *Service) Insertar(entidadRef interface{}) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "Insertar", service.Trn, entidadRef)
}

func (service *Service) Actualizar(entidadRef interface{}) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "Actualizar", service.Trn, entidadRef)
}

func (service *Service) Eliminar(entidadRef interface{}) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "Eliminar", service.Trn, entidadRef)
}

func (service *Service) Guardar(entidadRef interface{}) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "Guardar", service.Trn, entidadRef)
}

func (service *Service) ActualizarLista(update coredto.Update) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "ActualizarLista", service.Trn, service.EntidadRef, update)
}

func (service *Service) EliminarLista(delete coredto.Delete) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "EliminarLista", service.Trn, service.EntidadRef, delete)
}

func (service *Service) NumeroRegistros(filtros []coredto.Filtro) (int, error) {
	v, e := corereflect.InvokeFuncReturnValueAndError(service.EntidadRef, "NumeroRegistros", service.Trn, service.EntidadRef, filtros)
	return v.(int), e
}

func (service *Service) BuscarPorId(entidadRef interface{}) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "BuscarPorId", service.Trn, entidadRef)
}

func (service *Service) CargarDetalle(entidadRef interface{}) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "CargarDetalle", service.Trn, entidadRef)
}

func (service *Service) GetEntidad(entidadRef interface{}, query coredto.Query) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "GetEntidad", service.Trn, entidadRef, query)
}

func (service *Service) GetEntidadList(listaRef interface{}, query coredto.Query) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "GetEntidadList", service.Trn, service.EntidadRef, query, listaRef)
}

func (service *Service) GetObjetoList(listaRef interface{}, query coredto.Query) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "GetObjetoList", service.Trn, service.EntidadRef, query, listaRef)
}

func (service *Service) GetObjeto(listaRef interface{}, query coredto.Query) error {
	return corereflect.InvokeFuncReturnError(service.EntidadRef, "GetObjeto", service.Trn, service.EntidadRef, query, listaRef)
}
