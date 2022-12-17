package ledger

import (
	"fmt"
	"github.com/cloudstruct/go-cardano-ledger/cbor"
)

const (
	ERA_ID_BABBAGE = 5

	BLOCK_TYPE_BABBAGE = 6

	BLOCK_HEADER_TYPE_BABBAGE = 5

	TX_TYPE_BABBAGE = 5
)

type BabbageBlock struct {
	cbor.StructAsArray
	Header                 BabbageBlockHeader
	TransactionBodies      []BabbageTransactionBody
	TransactionWitnessSets []AlonzoTransactionWitnessSet
	TransactionMetadataSet map[uint]cbor.Value
	InvalidTransactions    []uint
}

func (b *BabbageBlock) Id() string {
	return b.Header.Id()
}

type BabbageBlockHeader struct {
	cbor.StructAsArray
	id   string
	Body struct {
		cbor.StructAsArray
		BlockNumber   uint64
		Slot          uint64
		PrevHash      Blake2b256
		IssuerVkey    interface{}
		VrfKey        interface{}
		VrfResult     interface{}
		BlockBodySize uint32
		BlockBodyHash Blake2b256
		OpCert        struct {
			cbor.StructAsArray
			HotVkey        interface{}
			SequenceNumber uint32
			KesPeriod      uint32
			Signature      interface{}
		}
		ProtoVersion struct {
			cbor.StructAsArray
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
	cbor.StructAsArray
	Body       BabbageTransactionBody
	WitnessSet AlonzoTransactionWitnessSet
	IsValid    bool
	Metadata   cbor.Value
}

func NewBabbageBlockFromCbor(data []byte) (*BabbageBlock, error) {
	var babbageBlock BabbageBlock
	if _, err := cbor.Decode(data, &babbageBlock); err != nil {
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
	if _, err := cbor.Decode(data, &babbageBlockHeader); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	babbageBlockHeader.id, err = generateBlockHeaderHash(data, nil)
	return &babbageBlockHeader, err
}

func NewBabbageTransactionBodyFromCbor(data []byte) (*BabbageTransactionBody, error) {
	var babbageTx BabbageTransactionBody
	if _, err := cbor.Decode(data, &babbageTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &babbageTx, nil
}

func NewBabbageTransactionFromCbor(data []byte) (*BabbageTransaction, error) {
	var babbageTx BabbageTransaction
	if _, err := cbor.Decode(data, &babbageTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &babbageTx, nil
}
