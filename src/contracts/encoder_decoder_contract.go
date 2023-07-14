package contracts

import (
	"bytes"
	"io"
)

// Contract to encode and decode response body.
type IEncodeDecode interface {
	Decode(encoding string, body io.ReadCloser) ([]byte, error)
	Encode(encoding string, body []byte) (bytes.Buffer, error)
}
