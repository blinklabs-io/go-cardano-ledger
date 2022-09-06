package ledger

import (
	"fmt"
	"github.com/fxamacker/cbor/v2"
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
	// TODO: figure out how to parse properly
	// We use RawMessage here because the content is arbitrary and can contain data that
	// cannot easily be represented in Go (such as maps with bytestring keys)
	TransactionMetadataSet map[uint]cbor.RawMessage
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
	// TODO: figure out how to parse properly
	// We use RawMessage here because the content is arbitrary and can contain data that
	// cannot easily be represented in Go (such as maps with bytestring keys)
	Metadata cbor.RawMessage
}

func NewAllegraBlockFromCbor(data []byte) (*AllegraBlock, error) {
	var allegraBlock AllegraBlock
	if err := cbor.Unmarshal(data, &allegraBlock); err != nil {
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
	if err := cbor.Unmarshal(data, &allegraTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &allegraTx, nil
}

func NewAllegraTransactionFromCbor(data []byte) (*AllegraTransaction, error) {
	var allegraTx AllegraTransaction
	if err := cbor.Unmarshal(data, &allegraTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &allegraTx, nil
}
