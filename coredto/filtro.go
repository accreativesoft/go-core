package coredto

type Filtro struct {
	Campo      string      `json:"campo,omitempty"`
	Valor      interface{} `json:"valor,omitempty"`
	Operador   string      `json:"operador,omitempty"`
	GrupoAndOr string      `json:"grupoAndOr,omitempty"`
}

func NewFiltro(campo string, operador string, valor interface{}) Filtro {
	filtro := Filtro{}
	filtro.Campo = campo
	filtro.Operador = operador
	filtro.Valor = valor
	return filtro
}

func NewFiltroGrupo(campo string, operador string, valor interface{}, grupo string) Filtro {
	filtro := Filtro{}
	filtro.Campo = campo
	filtro.Operador = operador
	filtro.Valor = valor
	filtro.GrupoAndOr = grupo
	return filtro
}
