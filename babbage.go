package ledger

import (
	"fmt"
	"github.com/fxamacker/cbor/v2"
)

const (
	BLOCK_TYPE_BABBAGE = 6

	BLOCK_HEADER_TYPE_BABBAGE = 5

	TX_TYPE_BABBAGE = 5
)

type BabbageBlock struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_                      struct{} `cbor:",toarray"`
	Header                 BabbageBlockHeader
	TransactionBodies      []BabbageTransactionBody
	TransactionWitnessSets []AlonzoTransactionWitnessSet
	// TODO: figure out how to parse properly
	// We use RawMessage here because the content is arbitrary and can contain data that
	// cannot easily be represented in Go (such as maps with bytestring keys)
	TransactionMetadataSet map[uint]cbor.RawMessage
	InvalidTransactions    []uint
}

func (b *BabbageBlock) Id() string {
	return b.Header.Id()
}

type BabbageBlockHeader struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_    struct{} `cbor:",toarray"`
	id   string
	Body struct {
		// Tells the CBOR decoder to convert to/from a struct and a CBOR array
		_             struct{} `cbor:",toarray"`
		BlockNumber   uint64
		Slot          uint64
		PrevHash      Blake2b256
		IssuerVkey    interface{}
		VrfKey        interface{}
		VrfResult     interface{}
		BlockBodySize uint32
		BlockBodyHash Blake2b256
		OpCert        struct {
			// Tells the CBOR decoder to convert to/from a struct and a CBOR array
			_              struct{} `cbor:",toarray"`
			HotVkey        interface{}
			SequenceNumber uint32
			KesPeriod      uint32
			Signature      interface{}
		}
		ProtoVersion struct {
			// Tells the CBOR decoder to convert to/from a struct and a CBOR array
			_     struct{} `cbor:",toarray"`
			Major uint64
			Minor uint64
		}
	}
	Signature interface{}
}

func (h *BabbageBlockHeader) Id() string {
	return h.id
}

type BabbageTransactionBody struct {
	AlonzoTransactionBody
	CollateralReturn ShelleyTransactionOutput  `cbor:"16,keyasint,omitempty"`
	TotalCollateral  uint64                    `cbor:"17,keyasint,omitempty"`
	ReferenceInputs  []ShelleyTransactionInput `cbor:"18,keyasint,omitempty"`
}

type BabbageTransaction struct {
	Body       BabbageTransactionBody
	WitnessSet AlonzoTransactionWitnessSet
	IsValid    bool
	Metadata   interface{}
}

func NewBabbageBlockFromCbor(data []byte) (*BabbageBlock, error) {
	var babbageBlock BabbageBlock
	if err := cbor.Unmarshal(data, &babbageBlock); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	rawBlockHeader, err := extractHeaderFromBlockCbor(data)
	if err != nil {
		return nil, err
	}
	babbageBlock.Header.id, err = generateBlockHeaderHash(rawBlockHeader, nil)
	return &babbageBlock, err
}

func NewBabbageBlockHeaderFromCbor(data []byte) (*BabbageBlockHeader, error) {
	var err error
	var babbageBlockHeader BabbageBlockHeader
	if err := cbor.Unmarshal(data, &babbageBlockHeader); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	babbageBlockHeader.id, err = generateBlockHeaderHash(data, nil)
	return &babbageBlockHeader, err
}

func NewBabbageTransactionBodyFromCbor(data []byte) (*BabbageTransactionBody, error) {
	var babbageTx BabbageTransactionBody
	if err := cbor.Unmarshal(data, &babbageTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &babbageTx, nil
}

func NewBabbageTransactionFromCbor(data []byte) (*BabbageTransaction, error) {
	var babbageTx BabbageTransaction
	if err := cbor.Unmarshal(data, &babbageTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &babbageTx, nil
}
