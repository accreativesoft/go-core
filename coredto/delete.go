package coredto

import "strconv"

type Delete struct {
	Filtros  []Filtro `json:"filtros,omitempty"`
	Contador int      `json:"contador,omitempty"`
}

func (delete *Delete) AddFiltro(campo string, operador string, valor interface{}) {
	filtroTmp := NewFiltro(campo, operador, valor)
	delete.Filtros = append(delete.Filtros, filtroTmp)
}

func (delete *Delete) AddFiltroGrupoAnd(campo string, operador string, valor interface{}, grupo string) {
	filtroTmp := NewFiltroGrupo(campo, operador, valor, "AND~"+grupo+"~"+strconv.Itoa(delete.Contador))
	delete.Filtros = append(delete.Filtros, filtroTmp)
	delete.Contador++
}

func (delete *Delete) AddFiltroGrupoOr(campo string, operador string, valor interface{}, grupo string) {
	filtroTmp := NewFiltroGrupo(campo, operador, valor, "OR~"+grupo+"~"+strconv.Itoa(delete.Contador))
	delete.Filtros = append(delete.Filtros, filtroTmp)
	delete.Contador++
}
