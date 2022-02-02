package coredto

import "strconv"

type Query struct {
	PrimerResultado int      `json:"primerResultado,omitempty"`
	ResultadoMaximo int      `json:"resultadoMaximo,omitempty"`
	FiltroGlobal    string   `json:"filtroGlobal,omitempty"`
	Campos          []string `json:"campos,omitempty"`
	Ordenamientos   []Orden  `json:"ordenamientos,omitempty"`
	Filtros         []Filtro `json:"filtros,omitempty"`
	Contador        int      `json:"contador,omitempty"`
}

func (query *Query) AddCampos(campos ...string) {
	query.Campos = campos
}

func (query *Query) AddCampo(campo string) {
	query.Campos = append(query.Campos, campo)
}

func (query *Query) AddOrden(campo string, orden int) {
	ordenTmp := *new(Orden)
	ordenTmp.NewOrden(campo, orden)
	query.Ordenamientos = append(query.Ordenamientos, ordenTmp)
}

func (query *Query) AddFiltro(campo string, operador string, valor interface{}) {
	filtroTmp := *new(Filtro)
	filtroTmp.NewFiltro(campo, operador, valor)
	query.Filtros = append(query.Filtros, filtroTmp)
}

func (query *Query) AddFiltroGrupoAnd(campo string, operador string, valor interface{}, grupo string) {
	filtroTmp := *new(Filtro)
	filtroTmp.NewFiltroGrupoAnd(campo, operador, valor, "AND~"+grupo+"~"+strconv.Itoa(query.Contador))
	query.Filtros = append(query.Filtros, filtroTmp)
	query.Contador++
}

func (query *Query) AddFiltroGrupoOr(campo string, operador string, valor interface{}, grupo string) {
	filtroTmp := *new(Filtro)
	filtroTmp.NewFiltroGrupoOr(campo, operador, valor, "OR~"+grupo+"~"+strconv.Itoa(query.Contador))
	query.Filtros = append(query.Filtros, filtroTmp)
	query.Contador++
}
