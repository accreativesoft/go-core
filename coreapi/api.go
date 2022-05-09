package coreapi

import (
	"encoding/json"
	"reflect"

	"github.com/accreativesoft/go-core/corecons"
	"github.com/accreativesoft/go-core/coredto"
	"github.com/accreativesoft/go-core/coreerror"
	"github.com/accreativesoft/go-core/coremsg"
	"github.com/accreativesoft/go-core/corereflect"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/rs/zerolog/log"
)

type Api struct {
	EntidadListaRef interface{}
	EntidadRef      interface{}
	ServiceRef      interface{}
	Uri             string
}

func (api *Api) InitRoutes(app *fiber.App) {
	private := app.Group(api.Uri)
	private.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(corecons.JWT_SECRET),
	}))
	private.Put("/crear", api.crear)
	private.Put("/insertar", api.insertar)
	private.Put("/eliminar", api.eliminar)
	private.Put("/actualizar", api.actualizar)
	private.Post("/guardar", api.guardar)
	private.Put("/actualizarList", api.actualizarLista)
	private.Put("/eliminarList", api.eliminarLista)
	private.Put("/numeroRegistros", api.numeroRegistros)
	private.Put("/buscarPorId", api.buscarPorId)
	private.Put("/cargarDetalle", api.cargarDetalle)
	private.Put("/getEntidad", api.getEntidad)
	private.Put("/getEntidadList", api.getEntidadList)
	private.Put("/getObjetoList", api.getObjetoList)
	private.Put("/getObjeto", api.getObjeto)
}

/*
func GenerarToken(autenticacion Autenticacion) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = autenticacion.Id
	claims["exp"] = exp
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}*/

func (api *Api) crear(ctx *fiber.Ctx) error {

	// Recupero la referencia del objeto
	objectRef, e := api.GetObjectRef(ctx)
	if e != nil {
		return e
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "Crear", objectRef)
	if e != nil {
		return e
	}

	return ctx.JSON(objectRef)
}

func (api *Api) insertar(ctx *fiber.Ctx) error {

	// Recupero la referencia del objeto
	objectRef, e := api.GetObjectRef(ctx)
	if e != nil {
		return e
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "Insertar", objectRef)
	if e != nil {
		return e
	}

	return ctx.JSON(objectRef)
}

func (api *Api) eliminar(ctx *fiber.Ctx) error {

	// Recupero la referencia del objeto
	objectRef, e := api.GetObjectRef(ctx)
	if e != nil {
		return e
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "Eliminar", objectRef)
	if e != nil {
		return e
	}

	return ctx.JSON(objectRef)
}

func (api *Api) actualizar(ctx *fiber.Ctx) error {

	// Recupero la referencia del objeto
	objectRef, e := api.GetObjectRef(ctx)
	if e != nil {
		return e
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "Actualizar", objectRef)
	if e != nil {
		return e
	}

	return ctx.JSON(objectRef)
}

func (api *Api) guardar(ctx *fiber.Ctx) error {

	// Recupero la referencia del objeto
	objectRef, e := api.GetObjectRef(ctx)
	if e != nil {
		return e
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "Guardar", objectRef)
	if e != nil {
		return e
	}

	return ctx.JSON(objectRef)
}

func (api *Api) actualizarLista(ctx *fiber.Ctx) error {

	//Cast del objeto
	update := coredto.Update{}
	e := json.Unmarshal(ctx.Body(), &update)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "ActualizarLista", update)
	if e != nil {
		return e
	}

	return ctx.JSON(update)
}

func (api *Api) eliminarLista(ctx *fiber.Ctx) error {

	//Cast del objeto
	delete := coredto.Delete{}
	e := json.Unmarshal(ctx.Body(), &delete)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "EliminarLista", delete)
	if e != nil {
		return e
	}

	return ctx.JSON(delete)
}

func (api *Api) numeroRegistros(ctx *fiber.Ctx) error {

	//Cast del objeto
	var filtros []coredto.Filtro
	e := json.Unmarshal(ctx.Body(), &filtros)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	count, e := corereflect.InvokeFuncReturnValueAndError(api.ServiceRef, "NumeroRegistros", filtros)
	if e != nil {
		return e
	}

	return ctx.JSON(count)
}

func (api *Api) buscarPorId(ctx *fiber.Ctx) error {

	//Ejecuto
	objectRef, e := api.GetObjectRef(ctx)
	if e != nil {
		return e
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "BuscarPorId", objectRef)
	if e != nil {
		return e
	}

	return ctx.JSON(objectRef)
}

func (api *Api) cargarDetalle(ctx *fiber.Ctx) error {

	//Ejecuto
	objectRef, e := api.GetObjectRef(ctx)
	if e != nil {
		return e
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "CargarDetalle", objectRef)
	if e != nil {
		return e
	}

	return ctx.JSON(objectRef)
}

func (api *Api) getEntidad(ctx *fiber.Ctx) error {

	//Recupero el tipo elemento
	typeObject := reflect.TypeOf(api.EntidadRef).Elem()

	//Creo objeto principal para llenar listado
	objectRef := reflect.New(typeObject).Interface()

	//Cast del objeto
	query := coredto.Query{}
	e := json.Unmarshal(ctx.Body(), &query)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "GetEntidad", objectRef, query)
	if e != nil {
		return e
	}

	return ctx.JSON(objectRef)
}

func (api *Api) getEntidadList(ctx *fiber.Ctx) error {

	//Creo objeto principal para llenar listado
	elemType := reflect.TypeOf(api.EntidadListaRef).Elem()
	objectRef := reflect.New(elemType).Interface()

	//Cast del objeto
	query := coredto.Query{}
	e := json.Unmarshal(ctx.Body(), &query)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "GetEntidadList", objectRef, query)
	if e != nil {
		return e
	}

	return ctx.JSON(objectRef)
}

func (api *Api) getObjetoList(ctx *fiber.Ctx) error {

	//Creo objeto principal para llenar listado
	listaRef := make([]interface{}, 0)

	//Cast del objeto
	query := coredto.Query{}
	e := json.Unmarshal(ctx.Body(), &query)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "GetObjetoList", &listaRef, query)
	if e != nil {
		return e
	}

	return ctx.JSON(listaRef)
}

func (api *Api) getObjeto(ctx *fiber.Ctx) error {

	//Creo objeto principal para llenar listado
	listaRef := make([]interface{}, 0)

	//Cast del objeto
	query := coredto.Query{}
	e := json.Unmarshal(ctx.Body(), &query)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	e = corereflect.InvokeFuncReturnError(api.ServiceRef, "GetObjeto", &listaRef, query)
	if e != nil {
		return e
	}

	return ctx.JSON(listaRef)
}

func (api *Api) GetObjectRef(ctx *fiber.Ctx) (interface{}, error) {

	//Recupero el tipo elemento
	typeObject := reflect.TypeOf(api.EntidadRef).Elem()

	//Creo objeto principal para llenar listado
	objectRef := reflect.New(typeObject).Interface()

	//Cast del objeto
	e := json.Unmarshal(ctx.Body(), objectRef)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return nil, coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Cast del objeto
	return objectRef, nil

}
