package utils

import (
	"aad-auth-proxy/constants"
	"aad-auth-proxy/contracts"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"

	log "github.com/sirupsen/logrus"
)

type EncoderDecoder struct {
}

func NewEncoderDecoder() contracts.IEncoderDecoder {
	return &EncoderDecoder{}
}

func (encoderDecoder *EncoderDecoder) Decode(encoding string, body io.ReadCloser) ([]byte, error) {
	var reader io.ReadCloser
	var err error

	switch encoding {
	case constants.ENCODING_GZIP:
		reader, err = gzip.NewReader(body)
	case constants.ENCODING_DEFLATE_ZLIB:
		reader, err = zlib.NewReader(body)
	default:
		reader, err = body, nil
	}

	if err != nil {
		log.WithFields(log.Fields{
			"Encoding": encoding,
		}).Errorln("Failed to decode response body", err)
		return nil, err
	}

	defer reader.Close()

	var decodedBody []byte
	decodedBody, err = io.ReadAll(reader)
	if err != nil {
		log.WithFields(log.Fields{
			"Encoding": encoding,
		}).Errorln("Failed to decode response body", err)
		return nil, err
	}

	return decodedBody, nil
}

func (encoderDecoder *EncoderDecoder) Encode(encoding string, data []byte) (bytes.Buffer, error) {
	var buffer bytes.Buffer
	var writer io.Writer
	var err error

	switch encoding {
	case constants.ENCODING_GZIP:
		writer = gzip.NewWriter(&buffer)
	case constants.ENCODING_DEFLATE_ZLIB:
		writer = zlib.NewWriter(&buffer)
	default:
		writer = io.Writer(&buffer)
	}

	_, err = writer.Write(data)
	if err != nil {
		log.WithFields(log.Fields{
			"Encoding": encoding,
		}).Errorln("Failed to write response data", err)
		return buffer, err
	}

	if encoding == constants.ENCODING_GZIP || encoding == constants.ENCODING_DEFLATE_ZLIB {
		err = writer.(io.Closer).Close()
		if err != nil {
			log.WithFields(log.Fields{
				"Encoding": encoding,
			}).Errorln("Failed to encode response data", err)
			return buffer, err
		}
	}

	return buffer, nil
}
