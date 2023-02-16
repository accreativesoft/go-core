package coreapi

import (
	"encoding/json"
	"reflect"

	"github.com/accreativesoft/go-core/corecons"
	"github.com/accreativesoft/go-core/coredto"
	"github.com/accreativesoft/go-core/coreerror"
	"github.com/accreativesoft/go-core/coremsg"
	"github.com/accreativesoft/go-core/coresrv"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Api struct {
	EntidadListaRef interface{}
	EntidadRef      interface{}
	Uri             string
	Trn             *gorm.DB
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

func NewApi(uri string, trn *gorm.DB, entidadRef interface{}, entidadListaRef interface{}) *Api {
	return &Api{Uri: uri, Trn: trn, EntidadRef: entidadRef, EntidadListaRef: entidadListaRef}
}

func (api *Api) crear(ctx *fiber.Ctx) error {

	// Recupero la referencia del objeto
	objectRef, e := api.GetObjectRef(ctx)
	if e != nil {
		return e
	}

	//Ejecuto el servicio
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.Crear(objectRef); e != nil {
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
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.Insertar(objectRef); e != nil {
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
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.Eliminar(objectRef); e != nil {
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
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.Actualizar(objectRef); e != nil {
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
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.Guardar(objectRef); e != nil {
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
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.ActualizarLista(update); e != nil {
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
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.EliminarLista(delete); e != nil {
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
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	count, e := srv.NumeroRegistros(filtros)
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
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.BuscarPorId(objectRef); e != nil {
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
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.CargarDetalle(objectRef); e != nil {
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

	if e := json.Unmarshal(ctx.Body(), &query); e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.GetEntidad(objectRef, query); e != nil {
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

	if e := json.Unmarshal(ctx.Body(), &query); e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.GetEntidadList(objectRef, query); e != nil {
		return e
	}

	return ctx.JSON(objectRef)
}

func (api *Api) getObjetoList(ctx *fiber.Ctx) error {

	//Creo objeto principal para llenar listado
	listaRef := make([]interface{}, 0)

	//Cast del objeto
	query := coredto.Query{}

	if e := json.Unmarshal(ctx.Body(), &query); e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.GetObjetoList(&listaRef, query); e != nil {
		return e
	}

	return ctx.JSON(listaRef)

}

func (api *Api) getObjeto(ctx *fiber.Ctx) error {

	//Creo objeto principal para llenar listado
	listaRef := make([]interface{}, 0)

	//Cast del objeto
	query := coredto.Query{}

	if e := json.Unmarshal(ctx.Body(), &query); e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Ejecuto el servicio
	var srv coresrv.Service = coresrv.NewService(api.Trn, api.EntidadRef)

	if e := srv.GetObjeto(&listaRef, query); e != nil {
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

	if e := json.Unmarshal(ctx.Body(), objectRef); e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return nil, coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Cast del objeto
	return objectRef, nil

}
