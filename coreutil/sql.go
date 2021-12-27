package coreutil

import (
	"fmt"
	"go-core/cons"
	"go-core/dto"
	"reflect"
	"strconv"
	"strings"
)

func GetEntityListBase(entity interface{}, query dto.Query) {

	//Creacion de mapas a utilizar
	hsCampos := make(map[string]bool)
	hsFilters := make(map[string]string)

	//Lleno campos
	if query.Campos != nil {
		for _, campo := range query.Campos {
			hsCampos["campo¬"+campo] = true
		}
	}

	//Lleno Filtros
	if query.Filtros != nil {
		for _, filtro := range query.Filtros {
			if cons.NUMERO_REGISTROS == filtro.Campo {
				operador := cons.EQUALS
				if filtro.Operador != "" {
					operador = filtro.Operador
				}
				hsCampos["filtro¬"+filtro.Campo+"¬"+filtro.GrupoAndOr+operador+"¬"+filtro.Valor] = true
				hsFilters[filtro.GrupoAndOr+operador+filtro.Campo] = filtro.Valor
			}
		}
	}

	//Lleno ordenamiento
	for _, orden := range query.Ordenamientos {
		hsCampos["orden¬"+orden.Campo+"¬"+strconv.Itoa(orden.Orden)] = true
	}

	// Variables para formar JPA dinamico
	var sql strings.Builder
	var sqlSelect strings.Builder

	var sqlOrder strings.Builder
	/*


		var sqlJoins strings.Builder

		var sqlWhere strings.Builder
		hmJoins := make(map[string]string)
		hmParameters := make(map[string]string)
		hmFiltrosOrAnd := make(map[string]string)

		secuenciaAlias := 1
		secuenciaParametro := 1*/

	for campo := range hsCampos {

		claveJoins := ""
		campoArray := strings.Split(campo, "¬")

		tipoCampo := campoArray[0]
		properties := strings.Split(campoArray[1], ".")
		orden := " ASC"
		matchMode := ""
		filterValue := ""

		if tipoCampo == "orden" && campoArray[2] == "-1" {
			orden = " DESC"
		}

		if tipoCampo == "filtro" {
			orden = campoArray[2]
		}

		if tipoCampo == "filtro" {
			filterValue = hsFilters[matchMode+campoArray[1]]
		}

		fmt.Println(claveJoins)
		fmt.Println(campoArray)
		fmt.Println(tipoCampo)
		fmt.Println(properties)
		fmt.Println(orden)
		fmt.Println(filterValue)

		if len(properties) == 1 {

			fileld := properties[0]
			sentencia := "entidad." + fileld
			if tipoCampo == "campo" {
				sqlSelect.WriteString(sentencia)
				sqlSelect.WriteString(", ")
			} else if tipoCampo == "orden" {
				sqlOrder.WriteString(sentencia)
				sqlOrder.WriteString(orden)
				sqlOrder.WriteString(", ")
			} else if tipoCampo == "filtro" {

			}
		} else {

		}

	}

	rType := fmt.Sprint(reflect.TypeOf(entity))
	model := strings.Split(rType, ".")[1]

	sql.WriteString("SELECT ")
	sql.WriteString("( ")
	sql.WriteString(") ")
	sql.WriteString("FROM ")
	sql.WriteString(model)

	fmt.Println(sql.String())

}

func AgregarFiltro(campo string, operador string, valor string, tipoDato interface{}, sqlWhere strings.Builder, parametro1 string,
	parametro2 string, hmParametros map[string]string, hmFiltrosOrAnd map[string]string) {

	sign := "?"

	//value1 := ""
	//value2 := ""
	//SimpleDateFormat df = new SimpleDateFormat(ConstantesBase.FORMATO_FECHA);
	//DateTimeFormatter dtf = DateTimeFormatter.ofPattern(ConstantesBase.FORMATO_FECHA);

	if operador == cons.NOT_EQUALS {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" != ")
		sqlWhere.WriteString(sign)
		sqlWhere.WriteString(" AND ")
	} else if operador == cons.GREATER_THAN_OR_EQUAL {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" >= ")
		sqlWhere.WriteString(sign)
		sqlWhere.WriteString(" AND ")
	} else if operador == cons.LESS_THAN_OR_EQUAL {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" <= ")
		sqlWhere.WriteString(sign)
		sqlWhere.WriteString(" AND ")
	} else if operador == cons.GREATER_THAN {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" > ")
		sqlWhere.WriteString(sign)
		sqlWhere.WriteString(" AND ")
	} else if operador == cons.LESS_THAN {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" < ")
		sqlWhere.WriteString(sign)
		sqlWhere.WriteString(" AND ")
	} else if operador == cons.EQUALS {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" = ")
		sqlWhere.WriteString(sign)
		sqlWhere.WriteString(" AND ")
	} else if operador == cons.IS_NOT_NULL {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" IS NOT NULL  ")
		sqlWhere.WriteString(" AND ")
	} else if operador == cons.IS_NULL {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" IS NULL ")
		sqlWhere.WriteString(" AND ")
	} else if operador == cons.IN {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" IN ( ")
		sqlWhere.WriteString(sign)
		sqlWhere.WriteString(" ) ")
		sqlWhere.WriteString(" AND ")
	} else if operador == cons.NOT_IN {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" NOT IN ( ")
		sqlWhere.WriteString(sign)
		sqlWhere.WriteString(" ) ")
		sqlWhere.WriteString(" AND ")
	} else if operador == cons.EQUALS || operador == cons.STARTS_WITH || operador == cons.ENDS_WITH {
		if tipoDato == "string" {
			sqlWhere.WriteString("LOWER(")
			sqlWhere.WriteString(campo)
			sqlWhere.WriteString(") LIKE ")
			sqlWhere.WriteString(sign)
			sqlWhere.WriteString(" AND ")
		} else {
			sqlWhere.WriteString(campo)
			sqlWhere.WriteString(" = ")
			sqlWhere.WriteString(sign)
			sqlWhere.WriteString(" AND ")
		}
	} else if operador == cons.NOT_LIKE {
		if tipoDato == "string" {
			sqlWhere.WriteString("LOWER(")
			sqlWhere.WriteString(campo)
			sqlWhere.WriteString(") NOT LIKE ")
			sqlWhere.WriteString(sign)
			sqlWhere.WriteString(" AND ")
		} else {
			sqlWhere.WriteString(campo)
			sqlWhere.WriteString(" != :")
			sqlWhere.WriteString(sign)
			sqlWhere.WriteString(" AND ")
		}
	} else if operador == cons.BETWEEN {
		sqlWhere.WriteString(campo)
		sqlWhere.WriteString(" BETWEEN ")
		sqlWhere.WriteString(sign)
		sqlWhere.WriteString(" AND ")
		sqlWhere.WriteString(sign)
		sqlWhere.WriteString(" AND ")
	} else if strings.HasPrefix(operador, "OR") || strings.HasPrefix(operador, "AND") {

	}

}
