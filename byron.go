package ledger

import (
	"fmt"
	"github.com/fxamacker/cbor/v2"
)

const (
	BLOCK_TYPE_BYRON_EBB  = 0
	BLOCK_TYPE_BYRON_MAIN = 1

	BLOCK_HEADER_TYPE_BYRON = 0

	TX_TYPE_BYRON = 0
)

type ByronMainBlockHeader struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_             struct{} `cbor:",toarray"`
	id            string
	ProtocolMagic uint32
	PrevBlock     Blake2b256
	BodyProof     interface{}
	ConsensusData struct {
		// Tells the CBOR decoder to convert to/from a struct and a CBOR array
		_ struct{} `cbor:",toarray"`
		// [slotid, pubkey, difficulty, blocksig]
		SlotId struct {
			// Tells the CBOR decoder to convert to/from a struct and a CBOR array
			_     struct{} `cbor:",toarray"`
			Epoch uint64
			Slot  uint16
		}
		PubKey     []byte
		Difficulty struct {
			// Tells the CBOR decoder to convert to/from a struct and a CBOR array
			_       struct{} `cbor:",toarray"`
			Unknown uint64
		}
		BlockSig []interface{}
	}
	ExtraData struct {
		// Tells the CBOR decoder to convert to/from a struct and a CBOR array
		_            struct{} `cbor:",toarray"`
		BlockVersion struct {
			// Tells the CBOR decoder to convert to/from a struct and a CBOR array
			_       struct{} `cbor:",toarray"`
			Major   uint16
			Minor   uint16
			Unknown uint8
		}
		SoftwareVersion struct {
			// Tells the CBOR decoder to convert to/from a struct and a CBOR array
			_       struct{} `cbor:",toarray"`
			Name    string
			Unknown uint32
		}
		Attributes interface{}
		ExtraProof Blake2b256
	}
}

func (h *ByronMainBlockHeader) Id() string {
	return h.id
}

// TODO: flesh this out
type ByronTransaction interface{}

type ByronMainBlockBody struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_         struct{} `cbor:",toarray"`
	TxPayload []ByronTransaction
	// We keep this field as raw CBOR, since it contains a map with []byte
	// keys, which Go doesn't allow
	SscPayload cbor.RawMessage
	DlgPayload []interface{}
	UpdPayload []interface{}
}

type ByronEpochBoundaryBlockHeader struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_             struct{} `cbor:",toarray"`
	id            string
	ProtocolMagic uint32
	PrevBlock     Blake2b256
	BodyProof     interface{}
	ConsensusData struct {
		// Tells the CBOR decoder to convert to/from a struct and a CBOR array
		_          struct{} `cbor:",toarray"`
		Epoch      uint64
		Difficulty struct {
			// Tells the CBOR decoder to convert to/from a struct and a CBOR array
			_     struct{} `cbor:",toarray"`
			Value uint64
		}
	}
	ExtraData interface{}
}

func (h *ByronEpochBoundaryBlockHeader) Id() string {
	return h.id
}

type ByronMainBlock struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_      struct{} `cbor:",toarray"`
	Header ByronMainBlockHeader
	Body   ByronMainBlockBody
	Extra  []interface{}
}

func (b *ByronMainBlock) Id() string {
	return b.Header.Id()
}

type ByronEpochBoundaryBlock struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_      struct{} `cbor:",toarray"`
	Header ByronEpochBoundaryBlockHeader
	Body   []Blake2b224
	Extra  []interface{}
}

func (b *ByronEpochBoundaryBlock) Id() string {
	return b.Header.Id()
}

func NewByronEpochBoundaryBlockFromCbor(data []byte) (*ByronEpochBoundaryBlock, error) {
	var byronEbbBlock ByronEpochBoundaryBlock
	if err := cbor.Unmarshal(data, &byronEbbBlock); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	rawBlockHeader, err := extractHeaderFromBlockCbor(data)
	if err != nil {
		return nil, err
	}
	// Prepend bytes for CBOR list wrapper
	// The block hash is calculated with these extra bytes, so we have to add them to
	// get the correct value
	byronEbbBlock.Header.id, err = generateBlockHeaderHash(rawBlockHeader, []byte{0x82, BLOCK_TYPE_BYRON_EBB})
	return &byronEbbBlock, err
}

func NewByronEpochBoundaryBlockHeaderFromCbor(data []byte) (*ByronEpochBoundaryBlockHeader, error) {
	var err error
	var byronEbbBlockHeader ByronEpochBoundaryBlockHeader
	if err := cbor.Unmarshal(data, &byronEbbBlockHeader); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	// Prepend bytes for CBOR list wrapper
	// The block hash is calculated with these extra bytes, so we have to add them to
	// get the correct value
	byronEbbBlockHeader.id, err = generateBlockHeaderHash(data, []byte{0x82, BLOCK_TYPE_BYRON_EBB})
	return &byronEbbBlockHeader, err
}

func NewByronMainBlockFromCbor(data []byte) (*ByronMainBlock, error) {
	var byronMainBlock ByronMainBlock
	if err := cbor.Unmarshal(data, &byronMainBlock); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	rawBlockHeader, err := extractHeaderFromBlockCbor(data)
	if err != nil {
		return nil, err
	}
	// Prepend bytes for CBOR list wrapper
	// The block hash is calculated with these extra bytes, so we have to add them to
	// get the correct value
	byronMainBlock.Header.id, err = generateBlockHeaderHash(rawBlockHeader, []byte{0x82, BLOCK_TYPE_BYRON_MAIN})
	return &byronMainBlock, err
}

func NewByronMainBlockHeaderFromCbor(data []byte) (*ByronMainBlockHeader, error) {
	var err error
	var byronMainBlockHeader ByronMainBlockHeader
	if err := cbor.Unmarshal(data, &byronMainBlockHeader); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	// Prepend bytes for CBOR list wrapper
	// The block hash is calculated with these extra bytes, so we have to add them to
	// get the correct value
	byronMainBlockHeader.id, err = generateBlockHeaderHash(data, []byte{0x82, BLOCK_TYPE_BYRON_MAIN})
	return &byronMainBlockHeader, err
}

func NewByronTransactionFromCbor(data []byte) (*ByronTransaction, error) {
	var byronTx ByronTransaction
	if err := cbor.Unmarshal(data, &byronTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &byronTx, nil
}
