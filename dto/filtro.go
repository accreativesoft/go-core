package dto

type Filtro struct {
	Campo      string `json:"campo,omitempty"`
	Valor      string `json:"valor,omitempty"`
	Operador   string `json:"operador,omitempty"`
	GrupoAndOr string `json:"grupoAndOr,omitempty"`
}

func (filtro *Filtro) NewFiltro(campo string, operador string, valor string) {
	filtro.Campo = campo
	filtro.Operador = operador
	filtro.Valor = valor
}

func (filtro *Filtro) NewFiltroGrupoAnd(campo string, operador string, valor string, grupo string) {
	filtro.Campo = campo
	filtro.Operador = operador
	filtro.Valor = valor
	filtro.GrupoAndOr = grupo
}

func (filtro *Filtro) NewFiltroGrupoOr(campo string, operador string, valor string, grupo string) {
	filtro.Campo = campo
	filtro.Operador = operador
	filtro.Valor = valor
	filtro.GrupoAndOr = grupo
}
