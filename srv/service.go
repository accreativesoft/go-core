package srv

import (
	"go-core/util"

	"gorm.io/gorm"
)

type Service struct {
	EntidadRef interface{}
	Trn        *gorm.DB
}

func (service *Service) Guardar(EntidadRef interface{}) {
	util.Invoke(service.EntidadRef, "Guardar", service.Trn, EntidadRef)
}
