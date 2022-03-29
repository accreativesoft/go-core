package coresql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/accreativesoft/go-core/coreerror"
	"github.com/accreativesoft/go-core/coremsg"
	"github.com/accreativesoft/go-core/corereflect"
	"github.com/elliotchance/orderedmap"
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Insertar(trn *gorm.DB, entidadRef interface{}) error {

	//Recupero el nombre del dialector
	dialector := trn.Dialector.Name()

	//Recupero la sentencia insert
	sql := GetInsertSql(dialector, entidadRef)

	//Recupero los valores para insert
	vat := GetValoresInsert(entidadRef)

	//Lleno los valores
	values := make([]interface{}, 0)
	for el := vat.Front(); el != nil; el = el.Next() {
		values = append(values, el.Value)
	}

	//Obtengo el id del ultimo registro insertado
	var id uint

	switch dialector {
	case "postgres":
		//Ejecuto y tomo el ultimo id insertado
		rows, e := trn.Raw(sql, values...).Rows()
		if e != nil {
			log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
			return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
		}
		for rows.Next() {
			e := rows.Scan(&id)
			if e != nil {
				log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
				return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
			}
		}
	default:
		//Ejecuto el insert
		e := trn.Exec(sql, values...).Error
		if e != nil {
			log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
			return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
		}
		//Recupero el utimo insertado
		rows, er := trn.Raw("select LAST_INSERT_ID()").Rows()
		if er != nil {
			log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
			return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
		}
		for rows.Next() {
			e := rows.Scan(&id)
			if e != nil {
				log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
				return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
			}
		}
	}

	//Set del id del ultimo registro insertado
	object := reflect.ValueOf(entidadRef).Elem()
	valor := reflect.ValueOf(&id).Elem()
	object.FieldByName("Id").Set(valor)

	return nil
}

func GetInsertSql(dialector string, entityRef interface{}) string {

	sign := "?"

	var sql strings.Builder

	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])
	sql.WriteString("INSERT INTO ")
	sql.WriteString(model)
	sql.WriteString(" ( ")

	//Obtengo el valor
	v := reflect.ValueOf(entityRef).Elem()
	secuencia := 0

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
			sql.WriteString(f)
			sql.WriteString(", ")
			secuencia++
		} else {

			//Recupero valor de la referecia referencia Entidad
			ve := reflect.ValueOf(ref)

			//Si es una estructura
			if f == "Entidad" {

				for j := 0; j < ve.NumField(); j++ {
					f = reflect.TypeOf(ref).Field(j).Name
					f = strcase.ToSnake(f)
					if f != "id" {
						sql.WriteString(f)
						sql.WriteString(", ")
						secuencia++
					}
				}
			}
		}
	}

	sqltmp := sql.String()
	sql.Reset()
	sql.WriteString(sqltmp[0 : len(sqltmp)-2])

	//Escribo los valores
	sql.WriteString(" ) ")
	sql.WriteString("\n")
	sql.WriteString("VALUES (")
	for i := 0; i < secuencia; i++ {
		sql.WriteString(sign)
		if i < secuencia-1 {
			sql.WriteString(", ")
		}
	}
	sql.WriteString(")")

	switch dialector {
	case "postgres":
		sql.WriteString(" RETURNING id")
	default:

	}
	return sql.String()
}

func GetValoresInsert(entidadRef interface{}) *orderedmap.OrderedMap {

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

	return valores

}
