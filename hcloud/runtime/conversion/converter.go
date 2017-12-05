package conversion

import (
	"fmt"
	"reflect"
)

type Converter interface {
	RegisterConversionFunc(conversionFunc interface{}) error
	Convert(src, dest interface{}) error
	DefaultConvert(src, dest interface{}) error
}

func NewConverter() Converter {
	return &converter{
		conversionFuncs: newConversionFuncs(),
	}
}

type converter struct {
	conversionFuncs conversionFuncs
}

func (c *converter) RegisterConversionFunc(conversionFunc interface{}) error {
	return c.conversionFuncs.Add(conversionFunc)
}

func (c *converter) Convert(src, dest interface{}) error {
	return c.doConvert(src, dest, c.convertion)
}

func (c *converter) DefaultConvert(src, dest interface{}) error {
	return c.doConvert(src, dest, c.defaultConvertion)
}

func (c *converter) convertion(sv, dv reflect.Value, scope *scope) error {
	// Convert sv to dv.
	dt, st := dv.Type(), sv.Type()
	pair := typePair{st, dt}
	if fv, ok := c.conversionFuncs.fns[pair]; ok {
		return c.callConverter(sv, dv, fv, scope)
	}
	return fmt.Errorf("No converter registered for conversion of %v to %v", st, dt)
}

func (c *converter) defaultConvertion(sv, dv reflect.Value, scope *scope) (err error) {
	if sv.Type().AssignableTo(dv.Type()) {
		dv.Set(sv)
		return
	}

	if !sv.IsValid() {
		return
	}
	dt, st := dv.Type(), sv.Type()
	if dt.Kind() != reflect.Struct || st.Kind() != reflect.Struct {
		return fmt.Errorf("No struct")
	}

	for _, sourceFieldType := range reflectFields(st) {
		if sf := sv.FieldByName(sourceFieldType.Name); sf.IsValid() {
			if df := dv.FieldByName(sourceFieldType.Name); df.IsValid() && df.CanSet() && !set(df, sf) {
				continue
			}
		}
	}

	return nil
}

type conversionFunc func(sv, dv reflect.Value, scope *scope) error

func (c *converter) doConvert(src, dest interface{}, f conversionFunc) error {
	dv, err := enforcePtr(dest)
	if err != nil {
		return err
	}
	if !dv.CanAddr() && !dv.CanSet() {
		return fmt.Errorf("can't write to dest")
	}
	sv, err := enforcePtr(src)
	if err != nil {
		return err
	}
	s := &scope{
		converter: c,
	}
	return f(sv, dv, s)
}

func (c *converter) callConverter(sv, dv, custom reflect.Value, scope *scope) error {
	if !sv.CanAddr() {
		sv2 := reflect.New(sv.Type())
		sv2.Elem().Set(sv)
		sv = sv2
	} else {
		sv = sv.Addr()
	}
	if !dv.CanAddr() {
		if !dv.CanSet() {
			return fmt.Errorf("can't addr or set dest")
		}
		dvOrig := dv
		dv := reflect.New(dvOrig.Type())
		defer func() { dvOrig.Set(dv) }()
	} else {
		dv = dv.Addr()
	}
	args := []reflect.Value{sv, dv, reflect.ValueOf(scope)}
	ret := custom.Call(args)[0].Interface()
	// This convolution is necessary because nil interfaces won't convert
	// to errors.
	if ret == nil {
		return nil
	}
	return ret.(error)
}

type typePair struct {
	source reflect.Type
	dest   reflect.Type
}

// newConversionFuncs creates a new ConversionFuncs mapping
func newConversionFuncs() conversionFuncs {
	return conversionFuncs{fns: make(map[typePair]reflect.Value)}
}

// ConversionFuncs holds the typePair to conversionFunc mapping
type conversionFuncs struct {
	fns map[typePair]reflect.Value
}

// Add adds the provided conversion functions to the lookup table - they must have the signature
// `func(type1, type2, Scope) error`. Functions are added in the order passed and will override
// previously registered pairs.
func (c conversionFuncs) Add(fns ...interface{}) error {
	for _, fn := range fns {
		fv := reflect.ValueOf(fn)
		ft := fv.Type()
		if err := verifyConversionFunctionSignature(ft); err != nil {
			return err
		}
		c.fns[typePair{ft.In(0).Elem(), ft.In(1).Elem()}] = fv
	}
	return nil
}

// enforcePtr ensures that obj is a pointer of some sort. Returns a reflect.Value
// of the dereferenced pointer, ensuring that it is settable/addressable.
// Returns an error if this is not possible.
func enforcePtr(obj interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		if v.Kind() == reflect.Invalid {
			return reflect.Value{}, fmt.Errorf("expected pointer, but got invalid kind")
		}
		return reflect.Value{}, fmt.Errorf("expected pointer, but got %v type", v.Type())
	}
	if v.IsNil() {
		return reflect.Value{}, fmt.Errorf("expected pointer, but got nil")
	}
	return v.Elem(), nil
}

// Scope is passed to conversion funcs to allow them to continue an ongoing conversion.
type Scope interface {
	// Call Convert to convert sub-objects.
	Convert(src, dest interface{}) error
	DefaultConvert(src, dest interface{}) error
}

// scope contains information about an ongoing conversion.
type scope struct {
	converter Converter
}

func (s *scope) Convert(src, dest interface{}) error {
	return s.converter.Convert(src, dest)
}

func (s *scope) DefaultConvert(src, dest interface{}) error {
	return s.converter.DefaultConvert(src, dest)
}
