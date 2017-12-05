package runtime

import "github.com/hetznercloud/hcloud-go/hcloud/runtime/conversion"

// Scheme defines methods for serializing and deserializing objects, a type
// registry for converting group, version, and kind information to and from Go
// schemas, and mappings between Go schemas of different versions.
type Scheme struct {
	// converter stores all registered conversion functions. It also has
	// default coverting behavior.
	converter conversion.Converter
}

// NewScheme creates a new Scheme. This scheme is pluggable by default.
func NewScheme() *Scheme {
	s := &Scheme{}
	s.converter = conversion.NewConverter()
	return s
}

// Convert will attempt to convert in into out. Both must be pointers. For easy
// testing of conversion functions. Returns an error if the conversion isn't
// possible. You can call this with types that haven't been registered (for example,
// a to test conversion of types that are nested within registered types). The
// context interface is passed to the convertor.
func (s *Scheme) Convert(in, out interface{}) error {
	return s.converter.Convert(in, out)
}

// Converter allows access to the converter for the scheme
func (s *Scheme) Converter() conversion.Converter {
	return s.converter
}

// AddConversionFuncs adds functions to the list of conversion functions. The given
// functions should know how to convert between two of your API objects, or their
// sub-objects.
func (s *Scheme) AddConversionFuncs(conversionFuncs ...interface{}) error {
	for _, f := range conversionFuncs {
		if err := s.converter.RegisterConversionFunc(f); err != nil {
			return err
		}
	}
	return nil
}
