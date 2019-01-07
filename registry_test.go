package annotation

import (
	"reflect"
	"testing"
)

var InjectableAnnotationType = reflect.TypeOf((*InjectableAnnotation)(nil))

type InjectableAnnotation struct {
	Annotation `@name:"@Injectable"`
	Name       string `json:"name" @default:"" @required:"false"`
}

var InjectableServiceType = reflect.TypeOf((*InjectableService)(nil))

type InjectableService struct {
	TypeAnnotation `annotation:"@Injectable('name': 'InjectableService')"`
}

func TestNew(t *testing.T) {
	type args struct {
		parent Registry
	}
	tests := []struct {
		name string
		args args
		want Registry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.parent); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Register(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnnotationRegistry_Register(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "@Injectable",
			fields: fields{
				parent:          nil,
				definitions:     make(map[string]*Definition, 0),
				typeDefinitions: make(map[reflect.Type]*TypeDefinition, 0),
			},
			args: args{
				t: InjectableAnnotationType,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			if err := r.Register(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
			a, err := r.GetTypeAnnotation(InjectableServiceType, InjectableAnnotationType)
			if nil != err {
				t.Error(err)
			}
			t.Log(a)
		})
	}
}

func Test_findAnnotatedFields(t *testing.T) {
	type args struct {
		t    reflect.Type
		ft   reflect.Type
		deep bool
	}
	tests := []struct {
		name string
		args args
		want map[string]*reflect.StructField
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findAnnotatedFields(tt.args.t, tt.args.ft, tt.args.deep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAnnotatedFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotationRegistry_getAnnotation(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		f *reflect.StructField
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			got, err := r.getAnnotation(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.getAnnotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnnotationRegistry.getAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotationRegistry_getTypeDefinition(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "InjectableService",
			fields: fields{
				parent:          nil,
				definitions:     make(map[string]*Definition, 0),
				typeDefinitions: make(map[reflect.Type]*TypeDefinition, 0),
			},
			args: args{
				t: InjectableServiceType,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			r.Register(InjectableAnnotationType)

			got, err := r.getTypeDefinition(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.getTypeDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

func TestGetTypeAnnotation(t *testing.T) {
	type args struct {
		t  reflect.Type
		at reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTypeAnnotation(tt.args.t, tt.args.at)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTypeAnnotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTypeAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotationRegistry_GetTypeAnnotation(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		t  reflect.Type
		at reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			got, err := r.GetTypeAnnotation(tt.args.t, tt.args.at)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.GetTypeAnnotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnnotationRegistry.GetTypeAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTypeAnnotations(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTypeAnnotations(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTypeAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTypeAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotationRegistry_GetTypeAnnotations(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			got, err := r.GetTypeAnnotations(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.GetTypeAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnnotationRegistry.GetTypeAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFieldAnnotation(t *testing.T) {
	type args struct {
		t    reflect.Type
		name string
		at   reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFieldAnnotation(tt.args.t, tt.args.name, tt.args.at)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldAnnotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotationRegistry_GetFieldAnnotation(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		t    reflect.Type
		name string
		at   reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			got, err := r.GetFieldAnnotation(tt.args.t, tt.args.name, tt.args.at)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.GetFieldAnnotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnnotationRegistry.GetFieldAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFieldAnnotations(t *testing.T) {
	type args struct {
		t    reflect.Type
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFieldAnnotations(tt.args.t, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotationRegistry_GetFieldAnnotations(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		t    reflect.Type
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			got, err := r.GetFieldAnnotations(tt.args.t, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.GetFieldAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnnotationRegistry.GetFieldAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllFieldAnnotations(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllFieldAnnotations(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllFieldAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllFieldAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotationRegistry_GetAllFieldAnnotations(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			got, err := r.GetAllFieldAnnotations(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.GetAllFieldAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnnotationRegistry.GetAllFieldAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMethodAnnotation(t *testing.T) {
	type args struct {
		t    reflect.Type
		name string
		at   reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMethodAnnotation(tt.args.t, tt.args.name, tt.args.at)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMethodAnnotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMethodAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotationRegistry_GetMethodAnnotation(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		t    reflect.Type
		name string
		at   reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			got, err := r.GetMethodAnnotation(tt.args.t, tt.args.name, tt.args.at)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.GetMethodAnnotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnnotationRegistry.GetMethodAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMethodAnnotations(t *testing.T) {
	type args struct {
		t    reflect.Type
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMethodAnnotations(tt.args.t, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMethodAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMethodAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotationRegistry_GetMethodAnnotations(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		t    reflect.Type
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			got, err := r.GetMethodAnnotations(tt.args.t, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.GetMethodAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnnotationRegistry.GetMethodAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllMethodAnnotations(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllMethodAnnotations(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMethodAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllMethodAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnotationRegistry_GetAllMethodAnnotations(t *testing.T) {
	type fields struct {
		parent          Registry
		definitions     map[string]*Definition
		typeDefinitions map[reflect.Type]*TypeDefinition
	}
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]map[reflect.Type]Annotation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AnnotationRegistry{
				parent:          tt.fields.parent,
				definitions:     tt.fields.definitions,
				typeDefinitions: tt.fields.typeDefinitions,
			}
			got, err := r.GetAllMethodAnnotations(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnotationRegistry.GetAllMethodAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnnotationRegistry.GetAllMethodAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}
