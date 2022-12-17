package ledger

import (
	"fmt"
	"github.com/cloudstruct/go-cardano-ledger/cbor"
)

const (
	ERA_ID_MARY = 3

	BLOCK_TYPE_MARY = 4

	BLOCK_HEADER_TYPE_MARY = 3

	TX_TYPE_MARY = 3
)

type MaryBlock struct {
	cbor.StructAsArray
	Header                 ShelleyBlockHeader
	TransactionBodies      []MaryTransactionBody
	TransactionWitnessSets []ShelleyTransactionWitnessSet
	TransactionMetadataSet map[uint]cbor.Value
}

func (b *MaryBlock) Id() string {
	return b.Header.Id()
}

type MaryTransactionBody struct {
	AllegraTransactionBody
	//Outputs []MaryTransactionOutput `cbor:"1,keyasint,omitempty"`
	Outputs []cbor.Value `cbor:"1,keyasint,omitempty"`
	// TODO: further parsing of this field
	Mint cbor.Value `cbor:"9,keyasint,omitempty"`
}

type MaryTransaction struct {
	cbor.StructAsArray
	Body       MaryTransactionBody
	WitnessSet ShelleyTransactionWitnessSet
	Metadata   cbor.Value
}

// TODO: support both forms
/*
transaction_output = [address, amount : value]
value = coin / [coin,multiasset<uint>]
*/
//type MaryTransactionOutput interface{}

type MaryTransactionOutput cbor.Value

func NewMaryBlockFromCbor(data []byte) (*MaryBlock, error) {
	var maryBlock MaryBlock
	if _, err := cbor.Decode(data, &maryBlock); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	rawBlockHeader, err := extractHeaderFromBlockCbor(data)
	if err != nil {
		return nil, err
	}
	maryBlock.Header.id, err = generateBlockHeaderHash(rawBlockHeader, nil)
	return &maryBlock, err
}

func NewMaryTransactionBodyFromCbor(data []byte) (*MaryTransactionBody, error) {
	var maryTx MaryTransactionBody
	if _, err := cbor.Decode(data, &maryTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &maryTx, nil
}

func NewMaryTransactionFromCbor(data []byte) (*MaryTransaction, error) {
	var maryTx MaryTransaction
	if _, err := cbor.Decode(data, &maryTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &maryTx, nil
}
