package cbor

import (
	"bytes"
	_cbor "github.com/fxamacker/cbor/v2"
)

func Decode(dataBytes []byte, dest interface{}) (int, error) {
	data := bytes.NewReader(dataBytes)
	dec := _cbor.NewDecoder(data)
	err := dec.Decode(dest)
	return dec.NumBytesRead(), err
}
