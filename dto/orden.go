package dto

type Orden struct {
	Campo string `json:"campo,omitempty"`
	Orden int    `json:"orden,omitempty"`
}

func (o *Orden) NewOrden(campo string, orden int) {
	o.Campo = campo
	o.Orden = orden
}
