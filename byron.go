package ledger

import (
	"fmt"
	"github.com/cloudstruct/go-cardano-ledger/cbor"
)

const (
	ERA_ID_BYRON = 0

	BLOCK_TYPE_BYRON_EBB  = 0
	BLOCK_TYPE_BYRON_MAIN = 1

	BLOCK_HEADER_TYPE_BYRON = 0

	TX_TYPE_BYRON = 0
)

type ByronMainBlockHeader struct {
	cbor.StructAsArray
	id            string
	ProtocolMagic uint32
	PrevBlock     Blake2b256
	BodyProof     interface{}
	ConsensusData struct {
		cbor.StructAsArray
		// [slotid, pubkey, difficulty, blocksig]
		SlotId struct {
			cbor.StructAsArray
			Epoch uint64
			Slot  uint16
		}
		PubKey     []byte
		Difficulty struct {
			cbor.StructAsArray
			Unknown uint64
		}
		BlockSig []interface{}
	}
	ExtraData struct {
		cbor.StructAsArray
		BlockVersion struct {
			cbor.StructAsArray
			Major   uint16
			Minor   uint16
			Unknown uint8
		}
		SoftwareVersion struct {
			cbor.StructAsArray
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
type ByronTransactionBody interface{}

// TODO: flesh this out
type ByronTransaction interface{}

type ByronMainBlockBody struct {
	cbor.StructAsArray
	TxPayload  []ByronTransactionBody
	SscPayload cbor.Value
	DlgPayload []interface{}
	UpdPayload []interface{}
}

type ByronEpochBoundaryBlockHeader struct {
	cbor.StructAsArray
	id            string
	ProtocolMagic uint32
	PrevBlock     Blake2b256
	BodyProof     interface{}
	ConsensusData struct {
		cbor.StructAsArray
		Epoch      uint64
		Difficulty struct {
			cbor.StructAsArray
			Value uint64
		}
	}
	ExtraData interface{}
}

func (h *ByronEpochBoundaryBlockHeader) Id() string {
	return h.id
}

type ByronMainBlock struct {
	cbor.StructAsArray
	Header ByronMainBlockHeader
	Body   ByronMainBlockBody
	Extra  []interface{}
}

func (b *ByronMainBlock) Id() string {
	return b.Header.Id()
}

type ByronEpochBoundaryBlock struct {
	cbor.StructAsArray
	Header ByronEpochBoundaryBlockHeader
	Body   []Blake2b224
	Extra  []interface{}
}

func (b *ByronEpochBoundaryBlock) Id() string {
	return b.Header.Id()
}

func NewByronEpochBoundaryBlockFromCbor(data []byte) (*ByronEpochBoundaryBlock, error) {
	var byronEbbBlock ByronEpochBoundaryBlock
	if _, err := cbor.Decode(data, &byronEbbBlock); err != nil {
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
	if _, err := cbor.Decode(data, &byronEbbBlockHeader); err != nil {
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
	if _, err := cbor.Decode(data, &byronMainBlock); err != nil {
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
	if _, err := cbor.Decode(data, &byronMainBlockHeader); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	// Prepend bytes for CBOR list wrapper
	// The block hash is calculated with these extra bytes, so we have to add them to
	// get the correct value
	byronMainBlockHeader.id, err = generateBlockHeaderHash(data, []byte{0x82, BLOCK_TYPE_BYRON_MAIN})
	return &byronMainBlockHeader, err
}

func NewByronTransactionBodyFromCbor(data []byte) (*ByronTransactionBody, error) {
	var byronTx ByronTransactionBody
	if _, err := cbor.Decode(data, &byronTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &byronTx, nil
}

func NewByronTransactionFromCbor(data []byte) (*ByronTransaction, error) {
	var byronTx ByronTransaction
	if _, err := cbor.Decode(data, &byronTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &byronTx, nil
}
