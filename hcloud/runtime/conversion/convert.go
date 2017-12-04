package conversion

// DefaultConversions to convert basic values
var DefaultConversions = []interface{}{
	ConvertSliceByteToSliceByte,
}

// ConvertSliceByteToSliceByte prevents recursing into every byte
func ConvertSliceByteToSliceByte(in *[]byte, out *[]byte, s Scope) error {
	if *in == nil {
		*out = nil
		return nil
	}
	*out = make([]byte, len(*in))
	copy(*out, *in)
	return nil
}
