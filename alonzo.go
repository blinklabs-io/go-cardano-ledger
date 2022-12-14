package ledger

import (
	"fmt"
	"github.com/cloudstruct/go-cardano-ledger/cbor"
)

const (
	ERA_ID_ALONZO = 4

	BLOCK_TYPE_ALONZO = 5

	BLOCK_HEADER_TYPE_ALONZO = 4

	TX_TYPE_ALONZO = 4
)

type AlonzoBlock struct {
	cbor.StructAsArray
	Header                 ShelleyBlockHeader
	TransactionBodies      []AlonzoTransactionBody
	TransactionWitnessSets []AlonzoTransactionWitnessSet
	TransactionMetadataSet map[uint]cbor.Value
	InvalidTransactions    []uint
}

func (b *AlonzoBlock) Id() string {
	return b.Header.Id()
}

type AlonzoTransactionBody struct {
	MaryTransactionBody
	ScriptDataHash  Blake2b256                `cbor:"11,keyasint,omitempty"`
	Collateral      []ShelleyTransactionInput `cbor:"13,keyasint,omitempty"`
	RequiredSigners []Blake2b224              `cbor:"14,keyasint,omitempty"`
	NetworkId       uint8                     `cbor:"15,keyasint,omitempty"`
}

type AlonzoTransactionWitnessSet struct {
	ShelleyTransactionWitnessSet
	PlutusScripts interface{}  `cbor:"3,keyasint,omitempty"`
	PlutusData    []cbor.Value `cbor:"4,keyasint,omitempty"`
	Redeemers     []cbor.Value `cbor:"5,keyasint,omitempty"`
}

type AlonzoTransaction struct {
	cbor.StructAsArray
	Body       AlonzoTransactionBody
	WitnessSet AlonzoTransactionWitnessSet
	IsValid    bool
	Metadata   cbor.Value
}

func NewAlonzoBlockFromCbor(data []byte) (*AlonzoBlock, error) {
	var alonzoBlock AlonzoBlock
	if _, err := cbor.Decode(data, &alonzoBlock); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	rawBlockHeader, err := extractHeaderFromBlockCbor(data)
	if err != nil {
		return nil, err
	}
	alonzoBlock.Header.id, err = generateBlockHeaderHash(rawBlockHeader, nil)
	return &alonzoBlock, err
}

func NewAlonzoTransactionBodyFromCbor(data []byte) (*AlonzoTransactionBody, error) {
	var alonzoTx AlonzoTransactionBody
	if _, err := cbor.Decode(data, &alonzoTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &alonzoTx, nil
}

func NewAlonzoTransactionFromCbor(data []byte) (*AlonzoTransaction, error) {
	var alonzoTx AlonzoTransaction
	if _, err := cbor.Decode(data, &alonzoTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &alonzoTx, nil
}
