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

func Actualizar(trn *gorm.DB, entityRef interface{}) error {

	//Recupero el nombre del dialector
	//dialector := trn.Dialector.Name()

	//Recupero la sentencia insert
	sql := GetUpdateSql(entityRef)

	//Recupero los valores para insert
	vat := GetValoresUpdate(entityRef)

	//Lleno los valores
	values := make([]interface{}, 0)
	for el := vat.Front(); el != nil; el = el.Next() {
		values = append(values, el.Value)
	}

	//Ejecuto la consulta
	e := trn.Exec(sql, values...).Error
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}
	return nil
}

func ActualizarLista(trn *gorm.DB, entidadRef interface{}, update coredto.Update) error {

	//Recupero el nombre del dialector
	dialector := trn.Dialector.Name()

	//Recupero la sentencia select
	sql, e := GetUpdateFrom(dialector, entidadRef, update)
	if e != nil {
		return e
	}

	//Set de los parametros del sql
	values := make([]interface{}, 0)
	for _, campo := range update.Campos {
		values = append(values, campo.Valor)
	}
	for _, filtro := range update.Filtros {
		values = append(values, filtro.Valor)
	}

	//Ejecuto la consulta
	e = trn.Exec(sql, values...).Error
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}
	return nil
}

func GetUpdateSql(entityRef interface{}) string {

	var sql strings.Builder

	sign := "?"
	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])
	sql.WriteString("UPDATE ")
	sql.WriteString(model)
	sql.WriteString(" SET")

	//Obtengo el valor
	v := reflect.ValueOf(entityRef).Elem()

	//Recorro los campos
	for i := 0; i < v.NumField(); i++ {

		//Recupero nombre de campo
		f := reflect.TypeOf(entityRef).Elem().Field(i).Name

		//Recupero el campo
		ref, _ := corereflect.GetField(entityRef, f)

		//Recupero el kinf
		k := reflect.TypeOf(ref).Kind()

		if k != reflect.Struct {
			//Si el campo es entidad diferente a un struct
			f = strcase.ToSnake(f)
			sql.WriteString("\n")
			sql.WriteString(f)
			sql.WriteString(" = ")
			sql.WriteString(sign)
			sql.WriteString(", ")

		} else {

			//Recupero valor de la referecia referencia Entidad
			ve := reflect.ValueOf(ref)

			//Si es una estructura
			if f == "Entidad" {

				for j := 0; j < ve.NumField(); j++ {
					f = reflect.TypeOf(ref).Field(j).Name
					if f != "Id" {
						f = strcase.ToSnake(f)
						sql.WriteString("\n")
						sql.WriteString(f)
						sql.WriteString(" = ")
						sql.WriteString(sign)
						sql.WriteString(", ")
					}
				}
			}
		}
	}

	sqltmp := sql.String()
	sql.Reset()
	sql.WriteString(sqltmp[0 : len(sqltmp)-2])

	//Escribo los valores
	sql.WriteString("\n")
	sql.WriteString(" WHERE id = ")
	sql.WriteString(sign)

	return sql.String()

}

func GetUpdateFrom(dialector string, entityRef interface{}, update coredto.Update) (string, error) {
	switch dialector {
	case "postgres":
		sql, e := GetUpdateFromPostgres(entityRef, update)
		if e != nil {
			return "", e
		}
		return sql, nil
	default:
		sql, e := GetUpdateFromMysql(entityRef, update)
		if e != nil {
			return "", e
		}
		return sql, nil
	}
}

func GetUpdateFromMysql(entityRef interface{}, update coredto.Update) (string, error) {

	//Formo Query
	query := coredto.Query{}
	query.AddCampo("id")
	query.Filtros = update.Filtros

	//Recupero los joins de la consulta
	joins, e := GetJoins(entityRef, query)
	if e != nil {
		return "", e
	}

	//Declaracion de variables
	var sql strings.Builder
	sign := "?"

	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])
	j, _ := joins.Get(model)
	join := j.(*Join)

	//Cosntuyo el update
	sql.WriteString("UPDATE ")
	sql.WriteString(model)
	sql.WriteString(" ")
	sql.WriteString(join.Alias)

	//Formo los joins
	for el := joins.Front(); el != nil; el = el.Next() {
		j := el.Value
		jn := j.(*Join)
		if !strings.Contains(jn.Sql, "FROM") {
			sql.WriteString(jn.Sql)
		}
	}

	//Formo los sets
	sql.WriteString("\n")
	sql.WriteString("SET")
	for _, c := range update.Campos {

		//Split del campo
		propiedad := strcase.ToSnake(c.Campo)

		//Formo mi Update
		sql.WriteString("\n")
		sql.WriteString(join.Alias)
		sql.WriteString(".")
		sql.WriteString(propiedad)
		sql.WriteString(" = ")
		sql.WriteString(sign)
		sql.WriteString(",")

	}

	sqltmp := sql.String()[0 : len(sql.String())-1]
	sql.Reset()
	sql.WriteString(sqltmp)

	//Formo where
	sql.WriteString(GetWhereSql(entityRef, query, joins))

	//fmt.Println("sql--->", sql.String())

	return sql.String(), e

}

func GetUpdateFromPostgres(entityRef interface{}, update coredto.Update) (string, error) {

	//Formo Query
	query := coredto.Query{}
	query.AddCampo("id")
	query.Filtros = update.Filtros

	//Recupero los joins de la consulta
	joins, e := GetJoins(entityRef, query)
	if e != nil {
		return "", e
	}

	//Declaracion de variables
	var sql strings.Builder
	sign := "?"

	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])
	j, _ := joins.Get(model)
	join := j.(*Join)

	//Cosntuyo el update
	sql.WriteString("UPDATE ")
	sql.WriteString(model)

	//Formo los sets
	sql.WriteString("\n")
	sql.WriteString("SET")
	for _, c := range update.Campos {

		//Split del campo
		propiedad := strcase.ToSnake(c.Campo)

		//Formo mi Update
		sql.WriteString("\n")
		sql.WriteString(propiedad)
		sql.WriteString(" = ")
		sql.WriteString(sign)
		sql.WriteString(",")

	}

	sqltmp := sql.String()[0 : len(sql.String())-1]
	sql.Reset()
	sql.WriteString(sqltmp)

	//Formo los joins
	for el := joins.Front(); el != nil; el = el.Next() {
		j := el.Value
		jn := j.(*Join)
		sql.WriteString("\n")
		sql.WriteString(jn.Sql)
	}

	//Formo where
	where := GetWhereSql(entityRef, query, joins)

	if strings.TrimSpace(where) == "" {
		sql.WriteString("\n")
		sql.WriteString("WHERE ")
	} else {
		sql.WriteString(where)
		sql.WriteString("\n")
		sql.WriteString("AND ")
	}

	sql.WriteString(model)
	sql.WriteString(".")
	sql.WriteString("id")
	sql.WriteString(" = ")
	sql.WriteString(join.Alias)
	sql.WriteString(".")
	sql.WriteString("id")

	//fmt.Println("sql--->", sql.String())

	return sql.String(), nil

}

func GetValoresUpdate(entidadRef interface{}) *orderedmap.OrderedMap {

	//Obtengo los valores
	v := reflect.ValueOf(entidadRef).Elem()

	//Valores
	valores := orderedmap.NewOrderedMap()

	//Recorro los campos
	for i := 0; i < v.NumField(); i++ {

		//Recupero nombre de campo
		f := reflect.TypeOf(entidadRef).Elem().Field(i).Name

		//Recupero el campo
		ref, _ := corereflect.GetField(entidadRef, f)

		//Recupero el kinf
		k := reflect.TypeOf(ref).Kind()

		if k != reflect.Struct {

			//Si el campo es entidad diferente a un struct
			valores.Set(f, ref)

		} else {
			//Recupero valor de la referecia referencia Entidad
			ve := reflect.ValueOf(ref)
			if f == "Entidad" {
				for j := 0; j < ve.NumField(); j++ {
					f = reflect.TypeOf(ref).Field(j).Name
					if f != "Id" {
						ref, _ := corereflect.GetField(ref, f)
						valores.Set(f, ref)
					}
				}
			} else {

				//Recupero el tag gorm
				t, _ := corereflect.GetFieldTag(entidadRef, f, "gorm")

				//Si existe el tag foreignKey
				if strings.Contains(t, "foreignKey") {

					//Recupero el foreingkey
					fk := strings.Split(t, "foreignKey:")[1]

					//Recupero el campo del fk
					ffk, _ := corereflect.GetField(entidadRef, fk)

					//Recupero el campo Id de la entidad realacionada
					id, _ := corereflect.GetField(ref, "Id")

					//Set del valor
					str := fmt.Sprintf("%v", id)
					if reflect.ValueOf(ffk).Kind() == reflect.Ptr && (str == "0" || str == "") {
						valores.Set(fk, nil)
					} else {
						valores.Set(fk, id)
					}
				}
			}
		}
	}

	//Recupero el campo
	ref, _ := corereflect.GetField(entidadRef, "Id")
	valores.Set("Id", ref)

	return valores

}
