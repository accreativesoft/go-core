package coreutil

import "reflect"

//Ejecuta metodo por refleccion
func InvokeFunc(objRef interface{}, nfunc string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(objRef).MethodByName(nfunc).Call(inputs)
}
