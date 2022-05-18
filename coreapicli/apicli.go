package coreapicli

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/accreativesoft/go-core/coredto"
	"github.com/accreativesoft/go-core/coreerror"
	"github.com/accreativesoft/go-core/coremsg"
	"github.com/rs/zerolog/log"
)

type ApiClient struct {
	EntidadListaRef interface{}
	EntidadRef      interface{}
	Token           string
	Uri             string
}

func (apiClient *ApiClient) Crear(entidadRef interface{}) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(entidadRef)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/crear", jsonBytes)
	if e != nil {
		return e
	}

	//Transformo el objeto
	e = apiClient.GetObjectRef(entidadRef, bodyBytes)
	if e != nil {
		return e
	}

	return nil
}

func (apiClient *ApiClient) Insertar(entidadRef interface{}) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(entidadRef)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/insertar", jsonBytes)
	if e != nil {
		return e
	}

	//Transformo el objeto
	e = apiClient.GetObjectRef(entidadRef, bodyBytes)
	if e != nil {
		return e
	}

	return nil
}

func (apiClient *ApiClient) Eliminar(entidadRef interface{}) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(entidadRef)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/eliminar", jsonBytes)
	if e != nil {
		return e
	}

	//Transformo el objeto
	e = apiClient.GetObjectRef(entidadRef, bodyBytes)
	if e != nil {
		return e
	}

	return nil
}

func (apiClient *ApiClient) Actualizar(entidadRef interface{}) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(entidadRef)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/actualizar", jsonBytes)
	if e != nil {
		return e
	}

	//Transformo el objeto
	e = apiClient.GetObjectRef(entidadRef, bodyBytes)
	if e != nil {
		return e
	}

	return nil
}

func (apiClient *ApiClient) Guardar(entidadRef interface{}) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(entidadRef)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPost, apiClient.Uri+"/guardar", jsonBytes)
	if e != nil {
		return e
	}

	//Transformo el objeto
	e = apiClient.GetObjectRef(entidadRef, bodyBytes)
	if e != nil {
		return e
	}

	return nil
}

func (apiClient *ApiClient) ActualizarLista(update coredto.Update) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(update)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/actualizarList", jsonBytes)
	if e != nil {
		return e
	}

	//Transformo el objeto
	e = json.Unmarshal(bodyBytes, &update)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	return nil
}

func (apiClient *ApiClient) EliminarLista(delete coredto.Delete) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(delete)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/eliminarList", jsonBytes)
	if e != nil {
		return e
	}

	//Transformo el objeto
	e = json.Unmarshal(bodyBytes, &delete)
	if e != nil {
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	return nil
}

func (apiClient *ApiClient) NumeroRegistros(filtros []coredto.Filtro) (int, error) {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(filtros)
	if e != nil {
		return 0, e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/numeroRegistros", jsonBytes)
	if e != nil {
		return 0, e
	}

	//Transformo el objeto
	var num int
	e = json.Unmarshal(bodyBytes, &num)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return 0, coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	return num, nil
}

func (apiClient *ApiClient) BuscarPorId(entidadRef interface{}) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(entidadRef)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/buscarPorId", jsonBytes)
	if e != nil {
		return e
	}

	//Transformo el objeto
	e = apiClient.GetObjectRef(entidadRef, bodyBytes)
	if e != nil {
		return e
	}

	return nil
}

func (apiClient *ApiClient) CargarDetalle(entidadRef interface{}) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(entidadRef)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/cargarDetalle", jsonBytes)
	if e != nil {
		return e
	}

	//Transformo el objeto
	e = apiClient.GetObjectRef(entidadRef, bodyBytes)
	if e != nil {
		return e
	}

	return nil
}

func (apiClient *ApiClient) GetEntidad(entidadRef interface{}, query coredto.Query) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(query)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/getEntidad", jsonBytes)
	if e != nil {
		return e
	}

	//Transformo el objeto
	e = apiClient.GetObjectRef(entidadRef, bodyBytes)
	if e != nil {
		return e
	}

	return nil
}

func (apiClient *ApiClient) GetEntidadList(listaRef interface{}, query coredto.Query) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(query)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/getEntidadList", jsonBytes)
	if e != nil {
		return e
	}

	//Creo objeto principal para llenar listado
	elemType := reflect.TypeOf(listaRef).Elem()
	objectRef := reflect.New(elemType).Interface()

	//Cast del objeto
	e = json.Unmarshal(bodyBytes, objectRef)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	return nil
}

func (apiClient *ApiClient) GetObjetoList(listaRef interface{}, query coredto.Query) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(query)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/getObjetoList", jsonBytes)
	if e != nil {
		return e
	}

	//Cast del objeto
	e = json.Unmarshal(bodyBytes, listaRef)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	return nil
}

func (apiClient *ApiClient) GetObjeto(listaRef interface{}, query coredto.Query) error {

	//Transformo mi entidad  a un json de bytes
	jsonBytes, e := apiClient.GetJsonBytes(query)
	if e != nil {
		return e
	}

	//Consumo del servicio
	bodyBytes, e := apiClient.ConsumeApi(http.MethodPut, apiClient.Uri+"/getObjeto", jsonBytes)
	if e != nil {
		return e
	}

	//Cast del objeto
	e = json.Unmarshal(bodyBytes, listaRef)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	return nil
}

func (apiClient *ApiClient) GetObjectRef(entidadRef interface{}, jsonBytes []byte) error {

	//Cast del objeto
	e := json.Unmarshal(jsonBytes, entidadRef)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO)
		return coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_JSON_A_OBJECTO, "")
	}

	//Cast del objeto
	return nil

}

func (apiClient *ApiClient) GetJsonBytes(entidadRef interface{}) ([]byte, error) {

	jsonRequest, e := json.Marshal(entidadRef)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_OBJETO_A_JSON)
		return nil, coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_OBJETO_A_JSON, "")

	}

	//Cast del objeto
	return jsonRequest, nil

}

func (apiClient *ApiClient) ConsumeApi(httpMethod string, url string, jsonBytes []byte) ([]byte, error) {

	//Creacion del request y cabeceras
	request, e := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonBytes))
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONSTRUIR_REQUEST_API)
		return nil, coreerror.NewError(coremsg.MSG_ERROR_CONSTRUIR_REQUEST_API, "")
	}
	request.Header.Add("Accept", "application/json; charset=utf-8")
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("Authorization", "Bearer "+apiClient.Token)

	//Creacion de cliente rest y ejecucion
	client := &http.Client{}
	response, e := client.Do(request)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONSUMIR_API + " " + url)
		return nil, coreerror.NewError(coremsg.MSG_ERROR_CONSUMIR_API, "", url)
	}
	defer response.Body.Close()

	//Recupero el body de la respuest
	bodyBytes, e := ioutil.ReadAll(response.Body)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_OBJETO_A_JSON)
		return nil, coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_OBJETO_A_JSON, "")
	}

	// Convert response body to string
	//bodyString := string(bodyBytes)
	//fmt.Println(bodyString)

	//Envio respuesta de json en bytes
	return bodyBytes, nil

}

func (apiClient *ApiClient) Consume(httpMethod string, url string, header http.Header, jsonBytes []byte) ([]byte, error) {

	//Creacion del request y cabeceras
	request, e := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonBytes))
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONSTRUIR_REQUEST_API)
		return nil, coreerror.NewError(coremsg.MSG_ERROR_CONSTRUIR_REQUEST_API, "")
	}

	//Set de header
	request.Header = header

	//Creacion de cliente rest y ejecucion
	client := &http.Client{}
	response, e := client.Do(request)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONSUMIR_API + " " + url)
		return nil, coreerror.NewError(coremsg.MSG_ERROR_CONSUMIR_API, "", url)
	}
	defer response.Body.Close()

	//Recupero el body de la response
	bodyBytes, e := ioutil.ReadAll(response.Body)
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_CONVERTIR_OBJETO_A_JSON)
		return nil, coreerror.NewError(coremsg.MSG_ERROR_CONVERTIR_OBJETO_A_JSON, "")
	}

	// Convert response body to string
	//bodyString := string(bodyBytes)
	//fmt.Println(bodyString)

	//Envio respuesta de json en bytes
	return bodyBytes, nil

}
