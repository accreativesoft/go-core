package coredto

type Campo struct {
	Campo string      `json:"campo,omitempty"`
	Valor interface{} `json:"valor,omitempty"`
}

func (o *Campo) NewCampo(campo string, valor interface{}) {
	o.Campo = campo
	o.Valor = valor
}
