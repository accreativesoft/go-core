package coredto

import "strconv"

type Update struct {
	Campos   []Campo  `json:"campos,omitempty"`
	Filtros  []Filtro `json:"filtros,omitempty"`
	Contador int      `json:"contador,omitempty"`
}

func (update *Update) AddCampo(campo string, valor interface{}) {
	campoTmp := *new(Campo)
	campoTmp.NewCampo(campo, valor)
	update.Campos = append(update.Campos, campoTmp)
}

func (update *Update) AddFiltro(campo string, operador string, valor interface{}) {
	filtroTmp := *new(Filtro)
	filtroTmp.NewFiltro(campo, operador, valor)
	update.Filtros = append(update.Filtros, filtroTmp)
}

func (update *Update) AddFiltroGrupoAnd(campo string, operador string, valor interface{}, grupo string) {
	filtroTmp := *new(Filtro)
	filtroTmp.NewFiltroGrupoAnd(campo, operador, valor, "AND~"+grupo+"~"+strconv.Itoa(update.Contador))
	update.Filtros = append(update.Filtros, filtroTmp)
	update.Contador++
}

func (update *Update) AddFiltroGrupoOr(campo string, operador string, valor interface{}, grupo string) {
	filtroTmp := *new(Filtro)
	filtroTmp.NewFiltroGrupoOr(campo, operador, valor, "OR~"+grupo+"~"+strconv.Itoa(update.Contador))
	update.Filtros = append(update.Filtros, filtroTmp)
	update.Contador++
}
