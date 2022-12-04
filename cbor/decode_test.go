package cbor_test

import (
	"encoding/hex"
	"github.com/cloudstruct/go-cardano-ledger/cbor"
	"reflect"
	"testing"
)

type decodeTestDefinition struct {
	CborHex   string
	Object    interface{}
	BytesRead int
}

var decodeTests = []decodeTestDefinition{
	// Simple list of numbers
	{
		CborHex: "83010203",
		Object:  []interface{}{uint64(1), uint64(2), uint64(3)},
	},
	// Multiple CBOR objects
	{
		CborHex:   "81018102",
		Object:    []interface{}{uint64(1)},
		BytesRead: 2,
	},
}

func TestDecode(t *testing.T) {
	for _, test := range decodeTests {
		cborData, err := hex.DecodeString(test.CborHex)
		if err != nil {
			t.Fatalf("failed to decode CBOR hex: %s", err)
		}
		var dest interface{}
		bytesRead, err := cbor.Decode(cborData, &dest)
		if err != nil {
			t.Fatalf("failed to decode CBOR: %s", err)
		}
		if test.BytesRead > 0 {
			if bytesRead != test.BytesRead {
				t.Fatalf("expected to read %d bytes, read %d instead", test.BytesRead, bytesRead)
			}
		}
		if !reflect.DeepEqual(dest, test.Object) {
			t.Fatalf("CBOR did not decode to expected object\n  got: %#v\n  wanted: %#v", dest, test.Object)
		}
	}
}
