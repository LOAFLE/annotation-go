package annotation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	lur "github.com/LOAFLE/util-go/reflect"
	yaml "gopkg.in/yaml.v2"
)

type Registry interface {
	Register(t reflect.Type) error
	GetTypeAnnotation(t reflect.Type, at reflect.Type) (Annotation, error)
	GetTypeAnnotations(t reflect.Type) (map[reflect.Type]Annotation, error)
	GetFieldAnnotation(t reflect.Type, name string, at reflect.Type) (Annotation, error)
	GetFieldAnnotations(t reflect.Type, name string) (map[reflect.Type]Annotation, error)
	GetFieldAnnotationsByType(t reflect.Type, at reflect.Type) (map[string]Annotation, error)
	GetAllFieldAnnotations(t reflect.Type) (map[string]map[reflect.Type]Annotation, error)
	GetMethodAnnotation(t reflect.Type, name string, at reflect.Type) (Annotation, error)
	GetMethodAnnotations(t reflect.Type, name string) (map[reflect.Type]Annotation, error)
	GetMethodAnnotationsByType(t reflect.Type, at reflect.Type) (map[string]Annotation, error)
	GetAllMethodAnnotations(t reflect.Type) (map[string]map[reflect.Type]Annotation, error)
}

func New(parent Registry) Registry {
	r := &AnnotationRegistry{
		parent:          parent,
		definitions:     make(map[string]*Definition, 0),
		typeDefinitions: make(map[reflect.Type]*TypeDefinition, 0),
	}
	if nil == r.parent {
		r.parent = SystemRegistry
	}
	return r
}

var SystemRegistry = &AnnotationRegistry{
	parent:          nil,
	definitions:     make(map[string]*Definition, 0),
	typeDefinitions: make(map[reflect.Type]*TypeDefinition, 0),
}

type AnnotationRegistry struct {
	parent          Registry
	definitions     map[string]*Definition
	typeDefinitions map[reflect.Type]*TypeDefinition
}

func Register(t reflect.Type) error {
	return SystemRegistry.Register(t)
}
func (r *AnnotationRegistry) Register(t reflect.Type) error {
	rt, _, _ := lur.GetTypeInfo(t)

	fields := findAnnotatedFields(t, AnnotationType, false)
	switch len(fields) {
	case 0:
		return fmt.Errorf("type[%s] is not Annotation", rt.Name())
	case 1:
	default:
		return fmt.Errorf("type[%s] have only one Annotation", rt.Name())
	}

	f := fields[AnnotationName]
	name := strings.TrimSpace(f.Tag.Get(NameTag))
	if "" == name {
		return fmt.Errorf("annotation name of type[%s] is not valid", rt.Name())
	}

	if _, ok := r.definitions[name]; ok {
		return fmt.Errorf("name[%s] of annotation exist already", name)
	}

	r.definitions[name] = &Definition{
		t:  t,
		rt: rt,
	}

	return nil
}

func findAnnotatedFields(t reflect.Type, ft reflect.Type, deep bool) map[string]*reflect.StructField {
	fields := make(map[string]*reflect.StructField, 0)

	rt, _, _ := lur.GetTypeInfo(t)
	if reflect.Struct != rt.Kind() {
		return fields
	}

LOOP:
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		if f.Anonymous {
			if f.Type == ft {
				fields[f.Name] = &f
				continue LOOP
			}
			if deep {
				_fields := findAnnotatedFields(f.Type, ft, deep)
				for _n, _f := range _fields {
					fields[_n] = _f
				}
			}
		}
	}
	return fields
}

func (r *AnnotationRegistry) getAnnotation(f *reflect.StructField) (map[reflect.Type]Annotation, error) {
	annotations := make(map[reflect.Type]Annotation, 0)

	tag := strings.TrimSpace(f.Tag.Get(AnnotationTag))
	if "" == tag {
		return annotations, nil
	}

	if !AnnotationRGX.MatchString(tag) {
		return nil, fmt.Errorf("Tag of annotation[%s] is not match", tag)
	}

	rss := AnnotationRGX.FindAllStringSubmatch(tag, -1)
	if nil == rss || 0 == len(rss) {
		return annotations, nil
	}

	for _, rs := range rss {
		if 3 != len(rs) {
			return nil, fmt.Errorf("Tag of annotation[%s] is not valid", rs[0])
		}

		name := fmt.Sprintf("@%s", strings.TrimSpace(rs[1]))
		body := strings.TrimSpace(rs[2])

		def, ok := r.definitions[name]
		if !ok {
			return nil, fmt.Errorf("annotation[%s] is not exist", name)
		}

		v := reflect.New(def.rt)
		i := v.Interface()

		if "" != body {
			if !AnnotationBodyRGX.MatchString(body) {
				return nil, fmt.Errorf("Body[%s] of annotation[%s] is not valid", body, name)
			}

			body = AnnotationBodyRGX.ReplaceAllStringFunc(body, func(token string) string {
				switch len(token) {
				case 0, 1, 2:
					return "\"\""
				default:
					return strings.Replace(fmt.Sprintf("\"%s\"", token[1:len(token)-1]), "\\'", "'", -1)
				}
			})
			body = fmt.Sprintf("{%s}", strings.TrimSpace(body))

			err := json.Unmarshal([]byte(body), i)
			if nil != err {
				return nil, fmt.Errorf("Unmarshal failed %v", err)
			}
		}

		_v := reflect.Indirect(reflect.ValueOf(i))
		for index := 0; index < _v.NumField(); index++ {
			_f := _v.Field(index)
			_fs := def.rt.Field(index)

			if !_f.CanAddr() || !_f.CanInterface() {
				continue
			}

			if isBlank := reflect.DeepEqual(_f.Interface(), reflect.Zero(_f.Type()).Interface()); isBlank {
				// Set default configuration if blank
				if value := _fs.Tag.Get(DefaultTag); value != "" {
					if err := yaml.Unmarshal([]byte(value), _f.Addr().Interface()); err != nil {
						return nil, fmt.Errorf("Unmarshal failed %v", err)
					}
				} else if _fs.Tag.Get(RequiredTag) == "true" {
					// return error if it is required but blank
					return nil, fmt.Errorf("%s is required, but blank", _fs.Name)
				}
			}
		}

		annotations[def.t] = i
	}

	return annotations, nil
}

func (r *AnnotationRegistry) getTypeDefinition(t reflect.Type) (*TypeDefinition, error) {
	rt, _, _ := lur.GetTypeInfo(t)
	if reflect.Struct != rt.Kind() {
		return nil, fmt.Errorf("type[%s] is not struct", rt.Name())
	}

	td, ok := r.typeDefinitions[t]
	if ok {
		return td, nil
	}

	td = &TypeDefinition{
		t:  t,
		rt: rt,
	}

LOOP:
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		if f.Anonymous {
			if f.Type == TypeAnnotationType {
				if 0 < len(td.typeAnnotation) {
					return nil, fmt.Errorf("TypeAnnotation of type[%s] be defined one more time", rt.Name())
				}
				as, err := r.getAnnotation(&f)
				if nil != err {
					return nil, err
				}
				if 0 == len(as) {
					continue LOOP
				}
				td.typeAnnotation = as
			}
		} else {
			if f.Type == MethodAnnotationType {
				as, err := r.getAnnotation(&f)
				if nil != err {
					return nil, err
				}
				if 0 == len(as) {
					continue LOOP
				}
				if 0 == len(td.methodAnnotation) {
					td.methodAnnotation = make(map[string]map[reflect.Type]Annotation, 0)
				}
				_name := strings.TrimLeft(f.Name, MethodAnnotationPrefix)
				if "" == _name {
					return nil, fmt.Errorf("name[%s] of method annotation is not valid", _name)
				}
				td.methodAnnotation[_name] = as
				continue LOOP
			} else {
				as, err := r.getAnnotation(&f)
				if nil != err {
					return nil, err
				}
				if 0 == len(as) {
					continue LOOP
				}
				if 0 == len(td.fieldAnnotation) {
					td.fieldAnnotation = make(map[string]map[reflect.Type]Annotation, 0)
				}
				td.fieldAnnotation[f.Name] = as
				continue LOOP
			}
		}
	}

	r.typeDefinitions[t] = td

	return td, nil
}

func GetTypeAnnotation(t reflect.Type, at reflect.Type) (Annotation, error) {
	return SystemRegistry.GetTypeAnnotation(t, at)
}
func (r *AnnotationRegistry) GetTypeAnnotation(t reflect.Type, at reflect.Type) (Annotation, error) {
	td, err := r.getTypeDefinition(t)
	if nil != err {
		return nil, err
	}

	if 0 == len(td.typeAnnotation) {
		return nil, nil
	}

	a, ok := td.typeAnnotation[at]
	if !ok {
		return nil, nil
	}

	return a, nil
}

func GetTypeAnnotations(t reflect.Type) (map[reflect.Type]Annotation, error) {
	return SystemRegistry.GetTypeAnnotations(t)
}
func (r *AnnotationRegistry) GetTypeAnnotations(t reflect.Type) (map[reflect.Type]Annotation, error) {
	td, err := r.getTypeDefinition(t)
	if nil != err {
		return nil, err
	}

	if 0 == len(td.typeAnnotation) {
		return nil, nil
	}

	return td.typeAnnotation, nil
}

func GetFieldAnnotation(t reflect.Type, name string, at reflect.Type) (Annotation, error) {
	return SystemRegistry.GetFieldAnnotation(t, name, at)
}
func (r *AnnotationRegistry) GetFieldAnnotation(t reflect.Type, name string, at reflect.Type) (Annotation, error) {
	td, err := r.getTypeDefinition(t)
	if nil != err {
		return nil, err
	}

	if 0 == len(td.fieldAnnotation) {
		return nil, nil
	}

	as, ok := td.fieldAnnotation[name]
	if !ok || 0 == len(as) {
		return nil, nil
	}

	a, ok := as[at]
	if !ok {
		return nil, nil
	}

	return a, nil
}

func GetFieldAnnotations(t reflect.Type, name string) (map[reflect.Type]Annotation, error) {
	return SystemRegistry.GetFieldAnnotations(t, name)
}
func (r *AnnotationRegistry) GetFieldAnnotations(t reflect.Type, name string) (map[reflect.Type]Annotation, error) {
	td, err := r.getTypeDefinition(t)
	if nil != err {
		return nil, err
	}

	if 0 == len(td.fieldAnnotation) {
		return nil, nil
	}

	as, ok := td.fieldAnnotation[name]
	if !ok || 0 == len(as) {
		return nil, nil
	}

	return as, nil
}

func GetFieldAnnotationsByType(t reflect.Type, at reflect.Type) (map[string]Annotation, error) {
	return SystemRegistry.GetFieldAnnotationsByType(t, at)
}
func (r *AnnotationRegistry) GetFieldAnnotationsByType(t reflect.Type, at reflect.Type) (map[string]Annotation, error) {
	td, err := r.getTypeDefinition(t)
	if nil != err {
		return nil, err
	}

	if 0 == len(td.fieldAnnotation) {
		return nil, nil
	}

	as := make(map[string]Annotation, 0)

	for k, v := range td.fieldAnnotation {
		if a, ok := v[at]; ok {
			as[k] = a
		}
	}

	return as, nil
}

func GetAllFieldAnnotations(t reflect.Type) (map[string]map[reflect.Type]Annotation, error) {
	return SystemRegistry.GetAllFieldAnnotations(t)
}
func (r *AnnotationRegistry) GetAllFieldAnnotations(t reflect.Type) (map[string]map[reflect.Type]Annotation, error) {
	td, err := r.getTypeDefinition(t)
	if nil != err {
		return nil, err
	}

	if 0 == len(td.fieldAnnotation) {
		return nil, nil
	}

	return td.fieldAnnotation, nil
}

func GetMethodAnnotation(t reflect.Type, name string, at reflect.Type) (Annotation, error) {
	return SystemRegistry.GetMethodAnnotation(t, name, at)
}
func (r *AnnotationRegistry) GetMethodAnnotation(t reflect.Type, name string, at reflect.Type) (Annotation, error) {
	td, err := r.getTypeDefinition(t)
	if nil != err {
		return nil, err
	}

	if 0 == len(td.methodAnnotation) {
		return nil, nil
	}

	as, ok := td.methodAnnotation[name]
	if !ok || 0 == len(as) {
		return nil, nil
	}

	a, ok := as[at]
	if !ok {
		return nil, nil
	}

	return a, nil
}

func GetMethodAnnotations(t reflect.Type, name string) (map[reflect.Type]Annotation, error) {
	return SystemRegistry.GetMethodAnnotations(t, name)
}
func (r *AnnotationRegistry) GetMethodAnnotations(t reflect.Type, name string) (map[reflect.Type]Annotation, error) {
	td, err := r.getTypeDefinition(t)
	if nil != err {
		return nil, err
	}

	if 0 == len(td.methodAnnotation) {
		return nil, nil
	}

	as, ok := td.methodAnnotation[name]
	if !ok || 0 == len(as) {
		return nil, nil
	}

	return as, nil
}

func GetMethodAnnotationsByType(t reflect.Type, at reflect.Type) (map[string]Annotation, error) {
	return SystemRegistry.GetMethodAnnotationsByType(t, at)
}
func (r *AnnotationRegistry) GetMethodAnnotationsByType(t reflect.Type, at reflect.Type) (map[string]Annotation, error) {
	td, err := r.getTypeDefinition(t)
	if nil != err {
		return nil, err
	}

	if 0 == len(td.methodAnnotation) {
		return nil, nil
	}

	as := make(map[string]Annotation, 0)

	for k, v := range td.methodAnnotation {
		if a, ok := v[at]; ok {
			as[k] = a
		}
	}

	return as, nil
}

func GetAllMethodAnnotations(t reflect.Type) (map[string]map[reflect.Type]Annotation, error) {
	return SystemRegistry.GetAllMethodAnnotations(t)
}
func (r *AnnotationRegistry) GetAllMethodAnnotations(t reflect.Type) (map[string]map[reflect.Type]Annotation, error) {
	td, err := r.getTypeDefinition(t)
	if nil != err {
		return nil, err
	}

	if 0 == len(td.methodAnnotation) {
		return nil, nil
	}

	return td.methodAnnotation, nil
}
