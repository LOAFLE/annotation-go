package annotation

import "reflect"

var AnnotationType = reflect.TypeOf((*Annotation)(nil)).Elem()

const AnnotationName = "Annotation"

type Annotation interface {
}

var TypeAnnotationType = reflect.TypeOf((*TypeAnnotation)(nil)).Elem()

const TypeAnnotationName = "TypeAnnotation"

type TypeAnnotation interface {
	Annotation
}

var MethodAnnotationType = reflect.TypeOf((*MethodAnnotation)(nil)).Elem()

const MethodAnnotationName = "MethodAnnotation"

type MethodAnnotation interface {
	Annotation
}
