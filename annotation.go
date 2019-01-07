package annotation

import (
	"regexp"
)

const (
	AnnotationTag = "annotation"

	NameTag     = "@name"
	DefaultTag  = "@default"
	RequiredTag = "@required"

	MethodAnnotationPrefix = "_"
)

var (
	AnnotationRGX     = regexp.MustCompile(`@((?s).*?)\(((?s).*?)\)`)
	AnnotationBodyRGX = regexp.MustCompile(`'([^\\\\']|\\')*'`)
)
