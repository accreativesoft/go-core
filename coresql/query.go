package coresql

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/accreativesoft/go-core/corecons"
	"github.com/accreativesoft/go-core/coredto"
	"github.com/accreativesoft/go-core/coreerror"
	"github.com/accreativesoft/go-core/coremsg"
	"github.com/accreativesoft/go-core/corereflect"
	"github.com/elliotchance/orderedmap"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

type Join struct {
	Alias  string
	Sql    string
	Campos []string
}

func Get(trn *gorm.DB, entidadRef interface{}, query coredto.Query) error {

	//Recupero el nombre del dialector
	dialector := trn.Dialector.Name()

	//Tomo solo un resultado
	query.PrimerResultado = 0
	query.ResultadoMaximo = 1

	//Recupero la sentencia select
	sql := GetSql(dialector, entidadRef, query)

	//Set de los parametros del sql
	values := make([]interface{}, 0)
	for _, filtro := range query.Filtros {
		values = append(values, filtro.Valor)
	}

	//Obtengo las filas
	rows, e := trn.Raw(sql, values...).Rows()
	if e != nil {
		return coreerror.NewError(coremsg.MSG_ERROR_SQL, e.Error())
	}

	//Formo el listado de valores a mapear con el resultado de la consulta
	valores := make([]interface{}, 0)
	GetValores(entidadRef, query, &valores)

	//Campos a mapear
	campos := query.Campos

	//Tipo de onjecto
	object := reflect.ValueOf(entidadRef).Elem()

	for rows.Next() {

		//Map de fila con los valores enviados
		rows.Scan(valores...)

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
	return nil
}

func GetLista(trn *gorm.DB, entidadRef interface{}, query coredto.Query, listaRef interface{}) error {

	//Recupero el nombre del dialector
	dialector := trn.Dialector.Name()

	//Recupero la sentencia select
	sql := GetSql(dialector, entidadRef, query)

	//Set de los parametros del sql
	values := make([]interface{}, 0)
	for _, filtro := range query.Filtros {
		values = append(values, filtro.Valor)
	}

	//Obtengo las filas
	rows, e := trn.Raw(sql, values...).Rows()
	if e != nil {
		return coreerror.NewError(coremsg.MSG_ERROR_SQL, e.Error(), "caca")
	}

	//Formo el listado de valores a mapear con el resultado de la consulta
	valores := make([]interface{}, 0)
	GetValores(entidadRef, query, &valores)

	//Campos a mapear
	campos := query.Campos

	//Tipo de onjecto
	typeObject := reflect.TypeOf(entidadRef).Elem()

	for rows.Next() {

		//Map de fila con los valores enviados
		rows.Scan(valores...)

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

func NumeroRegistros(trn *gorm.DB, entidadRef interface{}, filtros []coredto.Filtro) (int, error) {

	//Recupero la sentencia select
	sql := GetSqlCount(entidadRef, filtros)

	//Set de los parametros del sql
	values := make([]interface{}, 0)
	for _, filtro := range filtros {
		values = append(values, filtro.Valor)
	}

	//Obtengo las filas
	rows, e := trn.Raw(sql, values...).Rows()
	if e != nil {
		return 0, coreerror.NewError(coremsg.MSG_ERROR_SQL, e.Error())
	}

	//Valor
	var valor int

	for rows.Next() {
		//Map de fila con los valores enviados
		rows.Scan(&valor)
	}

	return valor, nil
}

func GetSql(dialector string, entityRef interface{}, query coredto.Query) string {

	//Variable sql unir las sentencias select, joins, where , order
	var sql strings.Builder

	//Recupero los joins de la consulta
	joins := GetJoins(entityRef, query)

	//Formo los selects
	sql.WriteString(GetSelectSql(entityRef, query, joins))

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

	return sql.String()
}

func GetSqlCount(entityRef interface{}, filtros []coredto.Filtro) string {

	//Variable sql unir las sentencias select, joins, where , order
	var sql strings.Builder

	//Constrtuyo el query
	var query = coredto.Query{}
	query.Filtros = filtros

	//Recupero los joins de la consulta
	joins := GetJoins(entityRef, query)

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

	return sql.String()
}

func GetJoins(entityRef interface{}, query coredto.Query) *orderedmap.OrderedMap {

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
	relaciones := make(map[string]bool)
	//joins := make(map[string]*Join)

	//Recupero el nombre de la referencia
	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])

	//Crea primer join de la entidad principal
	join := Join{}
	join.Alias = "e1"
	join.Sql = "\n" + "FROM " + model + " e1"
	//joins[model] = &join
	joins.Set(model, &join)

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

			if _, ok := relaciones[relacion]; !ok {

				//Split de las propiedades de la relacion
				propiedades := strings.Split(relacion, ".")

				//Recorro las propiedades de la relacion
				for _, propiedad := range propiedades {

					claveJoins = claveJoins + "." + propiedad

					if _, ok := joins.Get(claveJoins); !ok {

						//Recupero el tag  gorm
						t, _ := corereflect.GetFieldTag(ref, strcase.ToCamel(propiedad), "gorm")
						//fmt.Println("l->", t)

						//obtengo el foreignKey
						f := strcase.ToSnake(strings.Split(t, "foreignKey:")[1])
						//fmt.Println("f->", f)

						//Asigno el objeto de la propiedad
						u, _ := corereflect.GetField(ref, strcase.ToCamel(propiedad))
						ref = u

						//Formo el join
						join := Join{}
						join.Alias = "e" + strconv.Itoa(secuencia)
						join.Sql = "\n" + "LEFT JOIN " + strcase.ToSnake(propiedad) + " " + join.Alias + " ON " + join.Alias + ".id" + " = " + alias + "." + f
						//joins[claveJoins] = &join
						joins.Set(claveJoins, &join)

						alias = join.Alias
						secuencia++

					}

				}

				relaciones[relacion] = true

			}

		}
	}

	return joins
}

func GetSelectSql(entityRef interface{}, query coredto.Query, joins *orderedmap.OrderedMap) string {

	var sqlSelect strings.Builder
	sqlSelect.WriteString("SELECT")

	rType := fmt.Sprint(reflect.TypeOf(entityRef))
	model := strcase.ToLowerCamel(strings.Split(rType, ".")[1])

	for _, campo := range query.Campos {

		//Asigno la referencia principal
		ref := entityRef

		//tipoDato
		tipoDato := ""

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
			tipoDato, _ = corereflect.GetFieldType(ref, propiedad)

		} else {

			//Recorro las propiedades de la relacion
			for i, propiedad := range propiedades {

				propiedad = strcase.ToCamel(propiedad)

				if i == len(propiedades)-1 {
					tipoDato, _ = corereflect.GetFieldType(ref, propiedad)
				} else {
					//Asigno el objeto de la propiedad
					ref, _ = corereflect.GetField(ref, propiedad)
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
		if tipoDato == "int" || tipoDato == "int8" || tipoDato == "int16" || tipoDato == "int32" || tipoDato == "int64" || tipoDato == "uint" || tipoDato == "uint8" || tipoDato == "uint16" || tipoDato == "uint32" || tipoDato == "uint64" || tipoDato == "byte" || tipoDato == "rune" || tipoDato == "float32" || tipoDato == "float64" {
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
		} else {
			sqlSelect.WriteString("COALESCE(")
			sqlSelect.WriteString(join.Alias)
			sqlSelect.WriteString(".")
			sqlSelect.WriteString(propiedad)
			sqlSelect.WriteString(",'')")
		}
		sqlSelect.WriteString(", ")
	}

	return sqlSelect.String()[0 : len(sqlSelect.String())-2]

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

	grupos := make(map[string]string)

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
			if _, ok := grupos[claveGrupo]; !ok {
				//Verifico si anteriormente existe un salto de linea
				salto := "\n"
				if strings.HasSuffix(sqlWhere.String(), "\n") {
					salto = ""
				}
				grupos[claveGrupo] = salto + "AND (" + cmp + opr + sign
			} else {
				grupos[claveGrupo] = grupos[claveGrupo] + " " + condicionGrupo + " " + cmp + opr + sign
			}
		}
	}

	//Cierro parentensis de grupos
	for _, grupo := range grupos {
		grupo = grupo + ")"
		sqlWhere.WriteString(grupo)
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

func GetTipoDatos(entityRef interface{}, query coredto.Query) []string {

	//Formo arreglo para generar todos los joins en base a Campos, Filtros, Ordenamientos
	tipoDatos := make([]string, 0)

	for _, campo := range query.Campos {

		//Asigno la referencia principal
		ref := entityRef

		//Split del campo
		propiedades := strings.Split(campo, ".")

		if len(propiedades) == 1 {
			propiedad := strcase.ToCamel(propiedades[0])
			t, _ := corereflect.GetFieldType(ref, propiedad)
			tipoDatos = append(tipoDatos, t)

		} else if len(propiedades) > 1 {

			//Recorro las propiedades de la relacion
			for i, propiedad := range propiedades {

				propiedad = strcase.ToCamel(propiedad)

				if i == len(propiedades)-1 {
					t, _ := corereflect.GetFieldType(ref, propiedad)
					tipoDatos = append(tipoDatos, t)
				} else {
					//Asigno el objeto de la propiedad
					ref, _ = corereflect.GetField(ref, propiedad)
				}
			}
		}
	}

	return tipoDatos
}

func GetValores(entidadRef interface{}, query coredto.Query, valores interface{}) {

	tipoDatos := GetTipoDatos(entidadRef, query)

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
		default:
			var v string
			objectValor := reflect.ValueOf(&v)
			listaValor.Set(reflect.Append(listaValor, objectValor))
		}

	}
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
