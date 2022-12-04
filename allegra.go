package ledger

import (
	"fmt"
	"github.com/cloudstruct/go-cardano-ledger/cbor"
)

const (
	BLOCK_TYPE_ALLEGRA = 3

	BLOCK_HEADER_TYPE_ALLEGRA = 2

	TX_TYPE_ALLEGRA = 2
)

type AllegraBlock struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_                      struct{} `cbor:",toarray"`
	Header                 ShelleyBlockHeader
	TransactionBodies      []AllegraTransactionBody
	TransactionWitnessSets []ShelleyTransactionWitnessSet
	TransactionMetadataSet map[uint]cbor.Value
}

func (b *AllegraBlock) Id() string {
	return b.Header.Id()
}

type AllegraTransactionBody struct {
	ShelleyTransactionBody
	ValidityIntervalStart uint64 `cbor:"8,keyasint,omitempty"`
}

type AllegraTransaction struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_          struct{} `cbor:",toarray"`
	Body       AllegraTransactionBody
	WitnessSet ShelleyTransactionWitnessSet
	Metadata   cbor.Value
}

func NewAllegraBlockFromCbor(data []byte) (*AllegraBlock, error) {
	var allegraBlock AllegraBlock
	if _, err := cbor.Decode(data, &allegraBlock); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	rawBlockHeader, err := extractHeaderFromBlockCbor(data)
	if err != nil {
		return nil, err
	}
	allegraBlock.Header.id, err = generateBlockHeaderHash(rawBlockHeader, nil)
	return &allegraBlock, err
}

func NewAllegraTransactionBodyFromCbor(data []byte) (*AllegraTransactionBody, error) {
	var allegraTx AllegraTransactionBody
	if _, err := cbor.Decode(data, &allegraTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &allegraTx, nil
}

func NewAllegraTransactionFromCbor(data []byte) (*AllegraTransaction, error) {
	var allegraTx AllegraTransaction
	if _, err := cbor.Decode(data, &allegraTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &allegraTx, nil
}
