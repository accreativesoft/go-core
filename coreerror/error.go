package coreerror

type Error struct {
	Codigo  string
	Mensaje string
	Args    []string
}

func NewError(codigo string, mensaje string, args ...string) *Error {
	return &Error{
		Codigo:  codigo,
		Mensaje: mensaje,
		Args:    args,
	}
}

func (e *Error) Error() string {
	return e.Mensaje
}
