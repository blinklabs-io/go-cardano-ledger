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
	cbor.DecodeStoreCbor
	Header                 *MaryBlockHeader
	TransactionBodies      []MaryTransactionBody
	TransactionWitnessSets []ShelleyTransactionWitnessSet
	TransactionMetadataSet map[uint]cbor.Value
}

func (b *MaryBlock) UnmarshalCBOR(cborData []byte) error {
	return b.UnmarshalCborGeneric(cborData, b)
}

func (b *MaryBlock) Hash() string {
	return b.Header.Hash()
}

func (b *MaryBlock) BlockNumber() uint64 {
	return b.Header.BlockNumber()
}

func (b *MaryBlock) SlotNumber() uint64 {
	return b.Header.SlotNumber()
}

func (b *MaryBlock) Era() Era {
	return eras[ERA_ID_MARY]
}

func (b *MaryBlock) Transactions() []Transaction {
	// TODO
	return nil
}

type MaryBlockHeader struct {
	ShelleyBlockHeader
}

func (h *MaryBlockHeader) Era() Era {
	return eras[ERA_ID_MARY]
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
	return &maryBlock, nil
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
