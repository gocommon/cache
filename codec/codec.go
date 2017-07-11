package codec

// Codec Codec
type Codec interface {
	Encode(v interface{}) (string, error)
	Decode(data string, v interface{}) error
}
