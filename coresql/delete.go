package coresql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/accreativesoft/go-core/coredto"
	"github.com/accreativesoft/go-core/coreerror"
	"github.com/accreativesoft/go-core/coremsg"
	"github.com/accreativesoft/go-core/corereflect"
	"github.com/elliotchance/orderedmap"
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Eliminar(trn *gorm.DB, entityRef interface{}) error {

	//Recupero el nombre del dialector
	//dialector := trn.Dialector.Name()

	//Recupero la sentencia insert
	sql := GetDeleteSql(entityRef)

	//Recupero los valores para insert
	vat := GetValoresDelete(entityRef)

	//Lleno los valores
	values := make([]interface{}, 0)
	for el := vat.Front(); el != nil; el = el.Next() {
		values = append(values, el.Value)
	}

	//Ejecuto la consulta
	e := trn.Exec(sql, values...).Error
	if e != nil {
		log.Error().Msg(coremsg.MSG_ERROR_BACKEND)
		return coreerror.NewError(coremsg.MSG_ERROR_BACKEND, "")
	}
	return nil
}

func EliminarLista(trn *gorm.DB, entidadRef interface{}, delete coredto.Delete) error {

	//Recupero el nombre del dialector
	dialector := trn.Dialector.Name()

	//Recupero la sentencia select
	sql := GetDeleteFromSql(dialector, entidadRef, delete)

	//Set de los parametros del sql
	values := make([]interface{}, 0)
	for _, filtro := range delete.Filtros {
		values = append(values, filtro.Valor)
	}

	//Ejecuto la consulta
	e := trn.Exec(sql, values...).Error
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_ERROR_BACKEND)
		return coreerror.NewError(coremsg.MSG_ERROR_BACKEND, "")
	}
	return nil
}

func GetDeleteSql(entityRef interface{}) string {

	var sql strings.Builder

	sign := "?"
	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])
	sql.WriteString("DELETE FROM ")
	sql.WriteString(model)
	sql.WriteString(" WHERE id = ")
	sql.WriteString(sign)

	return sql.String()

}

func GetDeleteFromSql(dialector string, entityRef interface{}, delete coredto.Delete) string {
	switch dialector {
	case "postgres":
		return GetDeleteFromPostgres(entityRef, delete)
	default:
		return GetDeleteFromMysql(entityRef, delete)
	}
}

func GetDeleteFromMysql(entityRef interface{}, delete coredto.Delete) string {

	//Formo Query
	query := coredto.Query{}
	query.AddCampo("id")
	query.Filtros = delete.Filtros

	//Recupero los joins de la consulta
	joins := GetJoins(entityRef, query)

	//Declaracion de variables
	var sql strings.Builder

	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])
	j, _ := joins.Get(model)
	join := j.(*Join)

	//Cosntuyo el update
	sql.WriteString("DELETE ")
	sql.WriteString(join.Alias)

	//Formo los joins
	for el := joins.Front(); el != nil; el = el.Next() {
		j := el.Value
		join := j.(*Join)
		sql.WriteString(join.Sql)
	}

	//Formo where
	sql.WriteString(GetWhereSql(entityRef, query, joins))

	//fmt.Println("sql--->", sql.String())

	return sql.String()

}

func GetDeleteFromPostgres(entityRef interface{}, delete coredto.Delete) string {

	//Formo Query
	query := coredto.Query{}
	query.AddCampo("id")
	query.Filtros = delete.Filtros

	//Declaracion de variables
	var sql strings.Builder

	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])

	//Cosntuyo el update
	sql.WriteString("DELETE FROM ")
	sql.WriteString("\n")
	sql.WriteString(model)
	sql.WriteString("\n")
	sql.WriteString("WHERE id IN ( ")
	sql.WriteString("\n")
	sql.WriteString(GetSql("", entityRef, query))
	sql.WriteString(")")

	return sql.String()

}

func GetValoresDelete(entidadRef interface{}) *orderedmap.OrderedMap {

	//Valores
	valores := orderedmap.NewOrderedMap()

	//Recupero el campo
	ref, _ := corereflect.GetField(entidadRef, "Id")
	valores.Set("Id", ref)

	return valores

}
