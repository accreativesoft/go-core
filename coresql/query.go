package coresql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/accreativesoft/go-core/corecons"
	"github.com/accreativesoft/go-core/coredto"
	"github.com/accreativesoft/go-core/coreerror"
	"github.com/accreativesoft/go-core/coremsg"
	"github.com/accreativesoft/go-core/corereflect"
	"github.com/elliotchance/orderedmap"
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Join struct {
	Alias      string
	Sql        string
	Campos     []string
	Referencia interface{}
}

func GetEntidad(trn *gorm.DB, entidadRef interface{}, query coredto.Query) error {

	//Recupero el nombre del dialector
	dialector := trn.Dialector.Name()

	//Tomo solo un resultado
	query.PrimerResultado = 0
	query.ResultadoMaximo = 1

	//Recupero la sentencia select
	sql, e := GetSql(dialector, entidadRef, query)
	if e != nil {
		return e
	}

	//Set de los parametros del sql
	values := make([]interface{}, 0)
	for _, filtro := range query.Filtros {
		values = append(values, filtro.Valor)
	}

	//Obtengo las filas
	rows, e := trn.Raw(sql, values...).Rows()
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}

	e = Map(entidadRef, query, rows)
	if e != nil {
		return e
	}

	return nil
}

func GetEntidadList(trn *gorm.DB, entidadRef interface{}, query coredto.Query, listaRef interface{}) error {

	//Recupero el nombre del dialector
	dialector := trn.Dialector.Name()

	//Recupero la sentencia select
	sql, e := GetSql(dialector, entidadRef, query)
	if e != nil {
		return e
	}

	//Set de los parametros del sql
	values := make([]interface{}, 0)
	for _, filtro := range query.Filtros {
		values = append(values, filtro.Valor)
	}

	//Obtengo las filas
	rows, e := trn.Raw(sql, values...).Rows()
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}

	e = MapLista(entidadRef, listaRef, query, rows)
	if e != nil {
		return e
	}

	return nil

}

func GetObjetoList(trn *gorm.DB, entidadRef interface{}, query coredto.Query, listaRef *[]interface{}) error {

	//Recupero el nombre del dialector
	dialector := trn.Dialector.Name()

	//Recupero la sentencia select
	sql, e := GetSql(dialector, entidadRef, query)
	if e != nil {
		return e
	}

	//Set de los parametros del sql
	values := make([]interface{}, 0)
	for _, filtro := range query.Filtros {
		values = append(values, filtro.Valor)
	}

	//Obtengo las filas
	rows, e := trn.Raw(sql, values...).Rows()
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}

	//Formo el listado de valores a mapear con el resultado de la consulta
	valores := make([]interface{}, 0)
	GetValores(entidadRef, query.Campos, &valores)

	//Agrego el objeto a la lista
	//listaValor := reflect.ValueOf(listaRef).Elem()

	for rows.Next() {

		//Creo lista de objetos
		objectRef := make([]interface{}, 0)
		listaValor := reflect.ValueOf(&objectRef).Elem()

		//Map de fila con los valores enviados
		e := rows.Scan(valores...)
		if e != nil {
			log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
			return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
		}

		// data is passed by value
		for _, v := range valores {
			valor := reflect.ValueOf(v).Elem()
			listaValor.Set(reflect.Append(listaValor, valor))
		}

		*listaRef = append(*listaRef, &objectRef)

	}

	return nil

}

func GetObjeto(trn *gorm.DB, entidadRef interface{}, query coredto.Query, listaRef *[]interface{}) error {

	//Recupero el nombre del dialector
	dialector := trn.Dialector.Name()

	//Tomo solo un resultado
	query.PrimerResultado = 0
	query.ResultadoMaximo = 1

	//Recupero la sentencia select
	sql, e := GetSql(dialector, entidadRef, query)
	if e != nil {
		return e
	}

	//Set de los parametros del sql
	values := make([]interface{}, 0)
	for _, filtro := range query.Filtros {
		values = append(values, filtro.Valor)
	}

	//Obtengo las filas
	rows, e := trn.Raw(sql, values...).Rows()
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}

	//Formo el listado de valores a mapear con el resultado de la consulta
	valores := make([]interface{}, 0)
	GetValores(entidadRef, query.Campos, &valores)

	//Agrego el objeto a la lista
	listaValor := reflect.ValueOf(listaRef).Elem()

	for rows.Next() {

		//Map de fila con los valores enviados
		e := rows.Scan(valores...)
		if e != nil {
			log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
			return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
		}

		//Recorro los valores
		for _, v := range valores {
			objectValor := reflect.ValueOf(v).Elem()
			listaValor.Set(reflect.Append(listaValor, objectValor))
		}

	}

	return nil

}

func NumeroRegistros(trn *gorm.DB, entidadRef interface{}, filtros []coredto.Filtro) (int, error) {

	//Recupero la sentencia select
	sql, e := GetSqlCount(entidadRef, filtros)
	if e != nil {
		return 0, e
	}

	//Set de los parametros del sql
	values := make([]interface{}, 0)
	for _, filtro := range filtros {
		values = append(values, filtro.Valor)
	}

	//Obtengo las filas
	rows, e := trn.Raw(sql, values...).Rows()
	if e != nil {
		log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
		return 0, coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
	}

	//Valor
	var valor int

	for rows.Next() {
		//Map de fila con los valores enviados
		e := rows.Scan(&valor)
		if e != nil {
			log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
			return 0, coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
		}
	}

	return valor, nil
}

func GetSql(dialector string, entityRef interface{}, query coredto.Query) (string, error) {

	//Variable sql unir las sentencias select, joins, where , order
	var sql strings.Builder

	//Recupero los joins de la consulta
	joins, e := GetJoins(entityRef, query)
	if e != nil {
		return "", e
	}

	//Formo los selects
	selectSql, e := GetSelectSql(entityRef, query, joins)
	if e != nil {
		return "", e
	}
	sql.WriteString(selectSql)

	//Formo los joins
	for el := joins.Front(); el != nil; el = el.Next() {
		j := el.Value
		join := j.(*Join)
		sql.WriteString(join.Sql)
	}

	//Formo los where
	sql.WriteString(GetWhereSql(entityRef, query, joins))

	//Formo el orden
	sql.WriteString(GetOrderSql(entityRef, query, joins))

	sql.WriteString(GetLimite(dialector, query))

	//fmt.Println("sql-->", sql.String())

	return sql.String(), nil
}

func GetSqlCount(entityRef interface{}, filtros []coredto.Filtro) (string, error) {

	//Variable sql unir las sentencias select, joins, where , order
	var sql strings.Builder

	//Constrtuyo el query
	var query = coredto.Query{}
	query.Filtros = filtros

	//Recupero los joins de la consulta
	joins, e := GetJoins(entityRef, query)
	if e != nil {
		return "", e
	}

	//Formo los selects
	sql.WriteString("SELECT COUNT(1)")

	//Formo los joins
	for el := joins.Front(); el != nil; el = el.Next() {
		j := el.Value
		join := j.(*Join)
		sql.WriteString(join.Sql)
	}

	//Formo los where
	sql.WriteString(GetWhereSql(entityRef, query, joins))

	//fmt.Println("sql-->", sql.String())

	return sql.String(), nil
}

func GetJoins(entityRef interface{}, query coredto.Query) (*orderedmap.OrderedMap, error) {

	//Formo arreglo para generar todos los joins en base a Campos, Filtros, Ordenamientos
	campos := make([]string, 0)

	//Lleno campos
	campos = append(campos, query.Campos...)

	//Lleno Filtros
	for _, filtro := range query.Filtros {
		campos = append(campos, filtro.Campo)
	}

	//Lleno ordenamiento
	for _, orden := range query.Ordenamientos {
		campos = append(campos, orden.Campo)
	}

	//Definicion de varaibles
	secuencia := 2
	joins := orderedmap.NewOrderedMap()
	//relaciones := make(map[string]bool)
	//joins := make(map[string]*Join)

	//Recupero el nombre de la referencia
	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])

	//Crea primer join de la entidad principal
	join := Join{}
	join.Alias = "e1"
	join.Sql = "\n" + "FROM " + strcase.ToSnake(model) + " e1"
	joins.Set(model, &join)

	//Indicador si es slibe
	isSlice := false

	for _, campo := range campos {

		//Asigno la referencia principal
		ref := entityRef

		//Split del campo
		propiedades := strings.Split(campo, ".")

		if len(propiedades) > 1 {

			//Variables claveJoins y alias
			claveJoins := model
			alias := "e1"

			//Obtengo solo las propiedades con relaciones usuario.direccion.ubicacion
			penultProp := propiedades[len(propiedades)-1]
			index := strings.Index(campo, penultProp)
			relacion := campo[0 : index-1]

			//if _, ok := relaciones[relacion]; !ok {

			//Split de las propiedades de la relacion
			propiedades := strings.Split(relacion, ".")

			//Recorro las propiedades de la relacion
			for _, propiedad := range propiedades {

				//Recupero el padre
				claveJoins = claveJoins + "." + propiedad
				if j, ok := joins.Get(claveJoins); ok {
					join := j.(*Join)
					ref = join.Referencia
					alias = join.Alias
				}

				if _, ok := joins.Get(claveJoins); !ok {

					//Recupero el tag  gorm
					t, e := corereflect.GetFieldTag(ref, strcase.ToCamel(propiedad), "gorm")
					if e != nil {
						log.Error().Err(e).Msg(propiedad)
						log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
						return nil, coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
					}
					//fmt.Println("l->", t)

					//obtengo el foreignKey
					f := strcase.ToSnake(strings.Split(t, "foreignKey:")[1])
					if e != nil {
						log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
						return nil, coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
					}
					//fmt.Println("f->", f)

					//Asigno el objeto de la propiedad
					u, _ := corereflect.GetField(ref, strcase.ToCamel(propiedad))

					if reflect.TypeOf(u).Kind().String() == "slice" {
						//Si es una referencia de un arrglo rrecupeta el elemnte para obtener el typr
						items := reflect.ValueOf(u)
						datatype := items.Index(0).Type()
						u = reflect.New(datatype).Interface()
						isSlice = true
					}

					if e != nil {
						log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
						return nil, coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
					}
					ref = u

					//Obtengo el tipo de dato de la porpiedad y obtengo l struct para obtener el nombre de la tabla
					dt := reflect.ValueOf(u).Type().String()
					table := strcase.ToSnake(strings.Split(dt, ".")[1])

					//Formo el join
					join := Join{}
					join.Alias = "e" + strconv.Itoa(secuencia)
					if isSlice {
						join.Sql = "\n" + "LEFT JOIN " + table + " " + join.Alias + " ON " + join.Alias + "." + f + " = " + alias + ".id"
						isSlice = false
					} else {
						join.Sql = "\n" + "LEFT JOIN " + table + " " + join.Alias + " ON " + join.Alias + ".id" + " = " + alias + "." + f
					}
					join.Referencia = ref
					joins.Set(claveJoins, &join)

					alias = join.Alias
					secuencia++

				}

			}

			//relaciones[relacion] = true

			//}

		}
	}

	return joins, nil
}

func GetSelectSql(entityRef interface{}, query coredto.Query, joins *orderedmap.OrderedMap) (string, error) {

	var sqlSelect strings.Builder
	sqlSelect.WriteString("SELECT")

	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])

	for _, campo := range query.Campos {

		//Asigno la referencia principal
		ref := entityRef

		//tipoDato
		tipoDato := ""

		//error
		var e error

		//Split del campo
		var join *Join
		propiedades := strings.Split(campo, ".")

		if len(propiedades) == 1 {
			//Formo mi clave
			claveJoins := model
			j, _ := joins.Get(claveJoins)
			join = j.(*Join)

			//Recupero
			propiedad := strcase.ToCamel(propiedades[0])
			tipoDato, e = corereflect.GetFieldType(ref, propiedad)
			if e != nil {
				log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
				return "", coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
			}

		} else {

			//Recorro las propiedades de la relacion
			for i, propiedad := range propiedades {

				propiedad = strcase.ToCamel(propiedad)

				if i == len(propiedades)-1 {
					tipoDato, e = corereflect.GetFieldType(ref, propiedad)
					if e != nil {
						log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
						return "", coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
					}
				} else {
					//Asigno el objeto de la propiedad
					ref, e = corereflect.GetField(ref, propiedad)
					if e != nil {
						log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
						return "", coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
					}
				}
			}

			//Obtengo solo las propiedades con relaciones usuario.direccion.ubicacion
			ultProp := propiedades[len(propiedades)-1]
			index := strings.Index(campo, ultProp)
			relacion := campo[0 : index-1]

			//Formo mi clave
			claveJoins := model + "." + relacion
			j, _ := joins.Get(claveJoins)
			join = j.(*Join)

		}

		//Formo mi Select
		propiedad := strcase.ToSnake(propiedades[len(propiedades)-1])
		sqlSelect.WriteString("\n")
		if strings.Contains(tipoDato, "int") || strings.Contains(tipoDato, "uint") || strings.Contains(tipoDato, "byte") || strings.Contains(tipoDato, "rune") || strings.Contains(tipoDato, "float") {
			sqlSelect.WriteString("COALESCE(")
			sqlSelect.WriteString(join.Alias)
			sqlSelect.WriteString(".")
			sqlSelect.WriteString(propiedad)
			sqlSelect.WriteString(",0)")
		} else if tipoDato == "bool" {
			sqlSelect.WriteString("COALESCE(")
			sqlSelect.WriteString(join.Alias)
			sqlSelect.WriteString(".")
			sqlSelect.WriteString(propiedad)
			sqlSelect.WriteString(",false)")
		} else if tipoDato == "*time.Time" {
			sqlSelect.WriteString(join.Alias)
			sqlSelect.WriteString(".")
			sqlSelect.WriteString(propiedad)
		} else {
			sqlSelect.WriteString("COALESCE(")
			sqlSelect.WriteString(join.Alias)
			sqlSelect.WriteString(".")
			sqlSelect.WriteString(propiedad)
			sqlSelect.WriteString(",'')")
		}
		sqlSelect.WriteString(", ")
	}

	return sqlSelect.String()[0 : len(sqlSelect.String())-2], nil

}

func GetWhereSql(entityRef interface{}, query coredto.Query, joins *orderedmap.OrderedMap) string {

	var sqlWhere strings.Builder

	//Si la longitud de los filtros es mayor a cero
	if len(query.Filtros) > 0 {
		sqlWhere.WriteString("\n")
		sqlWhere.WriteString("WHERE ")
	} else {
		sqlWhere.WriteString("    ")
	}

	grupos := orderedmap.NewOrderedMap()

	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])

	for i, filtro := range query.Filtros {

		campo := filtro.Campo

		//Split del campo
		var join *Join
		propiedades := strings.Split(campo, ".")

		if len(propiedades) == 1 {

			//Formo mi clave
			claveJoins := model
			j, _ := joins.Get(claveJoins)
			join = j.(*Join)

		} else {

			//Obtengo solo las propiedades con relaciones usuario.direccion.ubicacion
			ultProp := propiedades[len(propiedades)-1]
			index := strings.Index(campo, ultProp)
			relacion := campo[0 : index-1]

			//Formo mi clave
			claveJoins := model + "." + relacion
			j, _ := joins.Get(claveJoins)
			join = j.(*Join)

		}

		//Recupero el operador y el nombre de la columna
		col := join.Alias + "." + strcase.ToSnake(propiedades[len(propiedades)-1])
		cmp, opr, sign := GetOperator(col, filtro.Operador)

		if len(filtro.GrupoAndOr) == 0 {

			//Formo mi sentencia where sin Grupos
			if !strings.HasSuffix(sqlWhere.String(), "WHERE ") {
				sqlWhere.WriteString("AND ")
			}
			sqlWhere.WriteString(cmp)
			sqlWhere.WriteString(opr)
			sqlWhere.WriteString(sign)
			//Condicion para no aumentar salto de linea si es el ultmo filtro
			if i != len(query.Filtros)-1 {
				sqlWhere.WriteString("\n")
			}
		} else {

			//Formo mi sentencia where con Grupos
			grupo := strings.Split(filtro.GrupoAndOr, "~")
			claveGrupo := grupo[1]
			condicionGrupo := grupo[0]
			if g, ok := grupos.Get(claveGrupo); !ok {
				//Verifico si anteriormente existe un salto de linea
				salto := "\n"
				if strings.HasSuffix(sqlWhere.String(), "\n") {
					salto = ""
				}
				grupos.Set(claveGrupo, salto+"AND ("+cmp+opr+sign)
			} else {
				gr := g.(string)
				grupos.Set(claveGrupo, gr+" "+condicionGrupo+" "+cmp+opr+sign)
			}
		}
	}

	//Cierro parentensis de grupos
	i := 0
	for el := grupos.Front(); el != nil; el = el.Next() {
		g := el.Value
		gr := g.(string)
		//Comparacion si no viene ninguna condicion y solo viene grupos
		if i == 0 && strings.Compare(sqlWhere.String(), "\nWHERE ") == 0 {
			sqlWhere.Reset()
			grupo := strings.Replace(gr, "AND ", "WHERE ", 1) + ")"
			sqlWhere.WriteString(grupo)
		} else {
			grupo := gr + ")"
			sqlWhere.WriteString(grupo)
		}
		i++
	}

	return sqlWhere.String()[0:len(sqlWhere.String())]

}

func GetOperator(campo string, operador string) (string, string, string) {
	sign := "?"
	switch operador {
	case corecons.NOT_EQUALS:
		return campo, " != ", sign
	case corecons.GREATER_THAN_OR_EQUAL:
		return campo, " >= ", sign
	case corecons.LESS_THAN_OR_EQUAL:
		return campo, " <= ", sign
	case corecons.GREATER_THAN:
		return campo, " > ", sign
	case corecons.LESS_THAN:
		return campo, " < ", sign
	case corecons.EQUALS:
		return campo, " = ", sign
	case corecons.IS_NOT_NULL:
		return campo, " IS NOT NULL ", ""
	case corecons.IS_NULL:
		return campo, " IS NULL ", ""
	case corecons.IN:
		return campo, " IN ", "(" + sign + ")"
	case corecons.NOT_IN:
		return campo, " NOT IN ", "(" + sign + ")"
	case corecons.NOT_LIKE:
		return " LOWER(" + campo + ")", " NOT LIKE ", "LOWER(" + sign + ")"
	case corecons.LIKE:
		return " LOWER(" + campo + ")", " LIKE ", "LOWER(" + sign + ")"
	case corecons.STARTS_WITH:
		return " LOWER(" + campo + ")", " LIKE ", "LOWER(" + sign + ")"
	case corecons.ENDS_WITH:
		return " LOWER(" + campo + ")", " LIKE ", "LOWER(" + sign + ")"
	default:
		return campo, " IS NOT NULL", ""
	}
}

func GetOrderSql(entityRef interface{}, query coredto.Query, joins *orderedmap.OrderedMap) string {

	var sqlOrder strings.Builder

	//Si la longitud del ordenamiento es mayor a cero
	if len(query.Ordenamientos) > 0 {
		sqlOrder.WriteString("\n")
		sqlOrder.WriteString("ORDER BY")
	} else {
		sqlOrder.WriteString("  ")
	}

	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])

	for _, orden := range query.Ordenamientos {

		campo := orden.Campo

		//Split del campo
		var join *Join
		propiedades := strings.Split(campo, ".")
		order := "ASC"
		if orden.Orden != 1 {
			order = "DESC"
		}

		if len(propiedades) == 1 {
			//Formo mi clave
			claveJoins := model
			j, _ := joins.Get(claveJoins)
			join = j.(*Join)

		} else {

			//Obtengo solo las propiedades con relaciones usuario.direccion.ubicacion
			ultProp := propiedades[len(propiedades)-1]
			index := strings.Index(campo, ultProp)
			relacion := campo[0 : index-1]

			//Formo mi clave
			claveJoins := model + "." + relacion
			j, _ := joins.Get(claveJoins)
			join = j.(*Join)

		}

		//Formo mi Order by
		sqlOrder.WriteString("\n")
		sqlOrder.WriteString(join.Alias)
		sqlOrder.WriteString(".")
		sqlOrder.WriteString(strcase.ToSnake(propiedades[len(propiedades)-1]))
		sqlOrder.WriteString(" ")
		sqlOrder.WriteString(order)
		sqlOrder.WriteString(", ")
	}

	return sqlOrder.String()[0 : len(sqlOrder.String())-2]

}

func GetLimite(dialector string, query coredto.Query) string {
	var sqlLimite strings.Builder
	sqlLimite.WriteString("\n")
	if query.ResultadoMaximo > 0 {
		switch dialector {
		case "postgres":
			sqlLimite.WriteString("OFFSET ")
			sqlLimite.WriteString(strconv.FormatUint(uint64(query.PrimerResultado), 10))
			sqlLimite.WriteString(" LIMIT ")
			sqlLimite.WriteString(strconv.FormatUint(uint64(query.ResultadoMaximo), 10))
		default:
			sqlLimite.WriteString("LIMIT ")
			sqlLimite.WriteString(strconv.FormatUint(uint64(query.PrimerResultado), 10))
			sqlLimite.WriteString(",")
			sqlLimite.WriteString(strconv.FormatUint(uint64(query.ResultadoMaximo), 10))
		}

	}
	return sqlLimite.String()
}

func GetTipoDatos(entityRef interface{}, campos []string) ([]string, error) {

	//Formo arreglo para generar todos los joins en base a Campos, Filtros, Ordenamientos
	tipoDatos := make([]string, 0)
	var e error

	for _, campo := range campos {

		//Asigno la referencia principal
		ref := entityRef

		//Split del campo
		propiedades := strings.Split(campo, ".")

		if len(propiedades) == 1 {
			propiedad := strcase.ToCamel(propiedades[0])
			t, e := corereflect.GetFieldType(ref, propiedad)
			if e != nil {
				log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
				return nil, coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
			}
			tipoDatos = append(tipoDatos, t)

		} else if len(propiedades) > 1 {

			//Recorro las propiedades de la relacion
			for i, propiedad := range propiedades {

				propiedad = strcase.ToCamel(propiedad)

				if i == len(propiedades)-1 {
					t, e := corereflect.GetFieldType(ref, propiedad)
					if e != nil {
						log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
						return nil, coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
					}
					tipoDatos = append(tipoDatos, t)
				} else {
					//Asigno el objeto de la propiedad
					ref, e = corereflect.GetField(ref, propiedad)
					if e != nil {
						log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
						return nil, coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
					}
				}
			}
		}
	}

	return tipoDatos, nil
}

func GetValores(entidadRef interface{}, campos []string, valores interface{}) error {

	tipoDatos, e := GetTipoDatos(entidadRef, campos)
	if e != nil {
		return e
	}

	//Agrego el objeto a la lista
	listaValor := reflect.ValueOf(valores).Elem()

	for _, tipoDato := range tipoDatos {
		switch tipoDato {
		case "int":
			var v int
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "int8":
			var v int8
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "int16":
			var v int16
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "int32":
			var v int32
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "int64":
			var v int64
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "uint":
			var v uint
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "uint8":
			var v uint8
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "uint16":
			var v uint16
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "uint32":
			var v uint32
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "float32":
			var v float32
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "float64":
			var v float64
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "complex64":
			var v complex64
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "complex128":
			var v complex128
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "byte":
			var v byte
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "rune":
			var v rune
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "string":
			var v string
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "bool":
			var v bool
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "time.Time":
			var v time.Time
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*int":
			var v *int
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*int8":
			var v *int8
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*int16":
			var v *int16
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*int32":
			var v *int32
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*int64":
			var v *int64
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*uint":
			var v *uint
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*uint8":
			var v *uint8
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*uint16":
			var v *uint16
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*uint32":
			var v *uint32
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*float32":
			var v *float32
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*float64":
			var v *float64
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*complex64":
			var v *complex64
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*complex128":
			var v *complex128
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*byte":
			var v *byte
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*rune":
			var v *rune
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*string":
			var v *string
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*bool":
			var v *bool
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		case "*time.Time":
			var v *time.Time
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		default:
			var v string
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		}

	}

	return nil
}

func Sign(dialector string, secuencia *int) string {
	sign := "?"
	switch dialector {
	case "postgres":
		v := reflect.ValueOf(secuencia).Elem()
		vi := int(v.Int())
		vi++
		v.SetInt(int64(vi))
		sign = "$" + strconv.Itoa(vi)
		return sign
	default:
		return sign
	}
}

func Map(entidadRef interface{}, query coredto.Query, rows *sql.Rows) error {

	//Formo el listado de valores a mapear con el resultado de la consulta
	valores := make([]interface{}, 0)
	e := GetValores(entidadRef, query.Campos, &valores)
	if e != nil {
		return e
	}

	//Campos a mapear
	campos := query.Campos

	//Tipo de onjecto
	object := reflect.ValueOf(entidadRef).Elem()

	i := 0

	for rows.Next() {

		i++

		//Map de fila con los valores enviados
		e := rows.Scan(valores...)
		if e != nil {
			log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
			return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
		}

		//Recorro los valores
		for i, v := range valores {

			// Recupero el campo y el valor
			campo := campos[i]
			valor := reflect.ValueOf(v).Elem()

			// Realizo el split del campo (Ej. tipoBodega.ubicacion.idUbicacion)
			propiedades := strings.Split(campo, ".")

			//Recupero la primera propiedad
			propiedad := strcase.ToCamel(propiedades[0])

			if len(propiedades) == 1 {
				//Set del valor de la propiedad ej. nombre
				object.FieldByName(propiedad).Set(valor)
			} else {

				//Recupero el campo de la propiedad relacion
				ref := object.FieldByName(propiedad)

				//Recorro las propiidades de tipo relacion
				for j := 1; j < len(propiedades); j++ {

					//Obtengo la propiedad
					propiedad = strcase.ToCamel(propiedades[j])

					if j == len(propiedades)-1 {
						//Set del valor para la ultima  propiedad
						reflect.Indirect(ref).FieldByName(propiedad).Set(valor)
					} else {
						//Recupero la propiedad relacion
						ref = reflect.Indirect(ref).FieldByName(propiedad)
						//fmt.Println("ref--->", ref)
					}
				}
			}
		}
	}

	if i == 0 {
		v := reflect.ValueOf(entidadRef).Elem()
		v.Set(reflect.Zero(v.Type()))
	}

	return nil
}

func MapLista(entidadRef interface{}, listaRef interface{}, query coredto.Query, rows *sql.Rows) error {

	//Formo el listado de valores a mapear con el resultado de la consulta
	valores := make([]interface{}, 0)
	GetValores(entidadRef, query.Campos, &valores)

	//Campos a mapear
	campos := query.Campos

	//Tipo de onjecto
	typeObject := reflect.TypeOf(entidadRef).Elem()

	for rows.Next() {

		//Map de fila con los valores enviados
		e := rows.Scan(valores...)
		if e != nil {
			log.Error().Err(e).Msg(coremsg.MSG_FALLA_INFRAESTRUCTURA)
			return coreerror.NewError(coremsg.MSG_FALLA_INFRAESTRUCTURA, "")
		}

		//Creo objeto principal para llenar listado
		objectRef := reflect.New(typeObject).Interface()

		//Obtengo valor de la referecia objectRef
		object := reflect.ValueOf(objectRef).Elem()

		//Recorro los valores
		for i, v := range valores {

			// Recupero el campo y el valor
			campo := campos[i]
			valor := reflect.ValueOf(v).Elem()

			// Realizo el split del campo (Ej. tipoBodega.ubicacion.idUbicacion)
			propiedades := strings.Split(campo, ".")

			//Recupero la primera propiedad
			propiedad := strcase.ToCamel(propiedades[0])

			if len(propiedades) == 1 {
				//Set del valor de la propiedad ej. nombre
				object.FieldByName(propiedad).Set(valor)
			} else {

				//Recupero el campo de la propiedad relacion
				ref := object.FieldByName(propiedad)

				//Recorro las propiidades de tipo relacion
				for j := 1; j < len(propiedades); j++ {

					//Obtengo la propiedad
					propiedad = strcase.ToCamel(propiedades[j])

					if j == len(propiedades)-1 {
						//Set del valor para la ultima  propiedad
						reflect.Indirect(ref).FieldByName(propiedad).Set(valor)
					} else {
						//Recupero la propiedad relacion
						ref = reflect.Indirect(ref).FieldByName(propiedad)
						//fmt.Println("ref--->", ref)
					}
				}
			}
		}

		//Agrego el objeto a la lista
		listaValor := reflect.ValueOf(listaRef).Elem()
		objectValor := reflect.ValueOf(objectRef).Elem()
		listaValor.Set(reflect.Append(listaValor, objectValor))

	}

	return nil

}
