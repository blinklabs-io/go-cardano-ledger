package cbor

import (
	_cbor "github.com/fxamacker/cbor/v2"
)

const (
	CBOR_TYPE_BYTE_STRING uint8 = 0x40
	CBOR_TYPE_TEXT_STRING uint8 = 0x60
	CBOR_TYPE_ARRAY       uint8 = 0x80
	CBOR_TYPE_MAP         uint8 = 0xa0

	// Only the top 3 bytes are used to specify the type
	CBOR_TYPE_MASK uint8 = 0xe0

	// Max value able to be stored in a single byte without type prefix
	CBOR_MAX_UINT_SIMPLE uint8 = 0x17
)

// Create an alias for RawMessage for convenience
type RawMessage = _cbor.RawMessage

// Useful for embedding and easier to remember
type StructAsArray struct {
	_ struct{} `cbor:",toarray"`
}
