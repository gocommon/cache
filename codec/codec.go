package codec

// Codec Codec
type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
	NewWithConf(jsonconf string) error
}
