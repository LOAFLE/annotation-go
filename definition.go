package annotation

import "reflect"

type Definition struct {
	t  reflect.Type
	rt reflect.Type
}

type TypeDefinition struct {
	t  reflect.Type
	rt reflect.Type

	typeAnnotation   map[reflect.Type]Annotation
	fieldAnnotation  map[string]map[reflect.Type]Annotation
	methodAnnotation map[string]map[reflect.Type]Annotation
}
