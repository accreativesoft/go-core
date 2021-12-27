package coredto

import "strconv"

type Query struct {
	PrimerResultado int      `json:"campo,omitempty"`
	ResultadoMaximo int      `json:"username,omitempty"`
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

func (query *Query) AddFiltro(campo string, operador string, valor string) {
	filtroTmp := *new(Filtro)
	filtroTmp.NewFiltro(campo, operador, valor)
	query.Filtros = append(query.Filtros, filtroTmp)
}

func (query *Query) NewFiltroGrupoAnd(campo string, operador string, valor string, grupo string) *Filtro {
	filtro := new(Filtro)
	query.Contador++
	filtro.NewFiltroGrupoAnd(campo, operador, valor, "AND~"+grupo+"~"+strconv.Itoa(query.Contador)+"~")
	return filtro
}

func (query *Query) NewFiltroGrupoOr(campo string, operador string, valor string, grupo string) *Filtro {
	filtro := new(Filtro)
	query.Contador++
	filtro.NewFiltroGrupoAnd(campo, operador, valor, "AND~"+grupo+"~"+strconv.Itoa(query.Contador)+"~")
	return filtro
}
