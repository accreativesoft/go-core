package coredto

type Filtro struct {
	Campo      string      `json:"campo,omitempty"`
	Valor      interface{} `json:"valor,omitempty"`
	Operador   string      `json:"operador,omitempty"`
	GrupoAndOr string      `json:"grupoAndOr,omitempty"`
}

func (filtro *Filtro) NewFiltro(campo string, operador string, valor interface{}) {
	filtro.Campo = campo
	filtro.Operador = operador
	filtro.Valor = valor
}

func (filtro *Filtro) NewFiltroGrupoAnd(campo string, operador string, valor interface{}, grupo string) {
	filtro.Campo = campo
	filtro.Operador = operador
	filtro.Valor = valor
	filtro.GrupoAndOr = grupo
}

func (filtro *Filtro) NewFiltroGrupoOr(campo string, operador string, valor interface{}, grupo string) {
	filtro.Campo = campo
	filtro.Operador = operador
	filtro.Valor = valor
	filtro.GrupoAndOr = grupo
}
