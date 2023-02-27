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
	cbor.DecodeStoreCbor
	hash          string
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

func (h *ByronMainBlockHeader) UnmarshalCBOR(cborData []byte) error {
	return h.UnmarshalCborGeneric(cborData, h)
}

func (h *ByronMainBlockHeader) Hash() string {
	if h.hash == "" {
		// Prepend bytes for CBOR list wrapper
		// The block hash is calculated with these extra bytes, so we have to add them to
		// get the correct value
		h.hash = generateBlockHeaderHash(h.Cbor(), []byte{0x82, BLOCK_TYPE_BYRON_EBB})
	}
	return h.hash
}

func (h *ByronMainBlockHeader) BlockNumber() uint64 {
	// Byron blocks don't store the block number in the block
	return 0
}

func (h *ByronMainBlockHeader) SlotNumber() uint64 {
	return uint64(h.ConsensusData.SlotId.Slot)
}

func (h *ByronMainBlockHeader) Era() Era {
	return eras[ERA_ID_BYRON]
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
	cbor.DecodeStoreCbor
	hash          string
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

func (h *ByronEpochBoundaryBlockHeader) UnmarshalCBOR(cborData []byte) error {
	return h.UnmarshalCborGeneric(cborData, h)
}

func (h *ByronEpochBoundaryBlockHeader) Hash() string {
	if h.hash == "" {
		// Prepend bytes for CBOR list wrapper
		// The block hash is calculated with these extra bytes, so we have to add them to
		// get the correct value
		h.hash = generateBlockHeaderHash(h.Cbor(), []byte{0x82, BLOCK_TYPE_BYRON_MAIN})
	}
	return h.hash
}

func (h *ByronEpochBoundaryBlockHeader) BlockNumber() uint64 {
	// Byron blocks don't store the block number in the block
	return 0
}

func (h *ByronEpochBoundaryBlockHeader) SlotNumber() uint64 {
	// There is no slot on boundary blocks
	return 0
}

func (h *ByronEpochBoundaryBlockHeader) Era() Era {
	return eras[ERA_ID_BYRON]
}

type ByronMainBlock struct {
	cbor.StructAsArray
	cbor.DecodeStoreCbor
	Header *ByronMainBlockHeader
	Body   ByronMainBlockBody
	Extra  []interface{}
}

func (b *ByronMainBlock) UnmarshalCBOR(cborData []byte) error {
	return b.UnmarshalCborGeneric(cborData, b)
}

func (b *ByronMainBlock) Hash() string {
	return b.Header.Hash()
}

func (b *ByronMainBlock) BlockNumber() uint64 {
	return b.Header.BlockNumber()
}

func (b *ByronMainBlock) SlotNumber() uint64 {
	return b.Header.SlotNumber()
}

func (b *ByronMainBlock) Era() Era {
	return b.Header.Era()
}

func (b *ByronMainBlock) Transactions() []TransactionBody {
	// TODO
	return nil
}

type ByronEpochBoundaryBlock struct {
	cbor.StructAsArray
	cbor.DecodeStoreCbor
	Header *ByronEpochBoundaryBlockHeader
	Body   []Blake2b224
	Extra  []interface{}
}

func (b *ByronEpochBoundaryBlock) UnmarshalCBOR(cborData []byte) error {
	return b.UnmarshalCborGeneric(cborData, b)
}

func (b *ByronEpochBoundaryBlock) Hash() string {
	return b.Header.Hash()
}

func (b *ByronEpochBoundaryBlock) BlockNumber() uint64 {
	return b.Header.BlockNumber()
}

func (b *ByronEpochBoundaryBlock) SlotNumber() uint64 {
	return b.Header.SlotNumber()
}

func (b *ByronEpochBoundaryBlock) Era() Era {
	return b.Header.Era()
}

func (b *ByronEpochBoundaryBlock) Transactions() []TransactionBody {
	// Boundary blocks don't have transactions
	return nil
}

func NewByronEpochBoundaryBlockFromCbor(data []byte) (*ByronEpochBoundaryBlock, error) {
	var byronEbbBlock ByronEpochBoundaryBlock
	if _, err := cbor.Decode(data, &byronEbbBlock); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &byronEbbBlock, nil
}

func NewByronEpochBoundaryBlockHeaderFromCbor(data []byte) (*ByronEpochBoundaryBlockHeader, error) {
	var byronEbbBlockHeader ByronEpochBoundaryBlockHeader
	if _, err := cbor.Decode(data, &byronEbbBlockHeader); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &byronEbbBlockHeader, nil
}

func NewByronMainBlockFromCbor(data []byte) (*ByronMainBlock, error) {
	var byronMainBlock ByronMainBlock
	if _, err := cbor.Decode(data, &byronMainBlock); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &byronMainBlock, nil
}

func NewByronMainBlockHeaderFromCbor(data []byte) (*ByronMainBlockHeader, error) {
	var byronMainBlockHeader ByronMainBlockHeader
	if _, err := cbor.Decode(data, &byronMainBlockHeader); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &byronMainBlockHeader, nil
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
