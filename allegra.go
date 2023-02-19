package ledger

import (
	"fmt"
	"github.com/cloudstruct/go-cardano-ledger/cbor"
)

const (
	ERA_ID_ALLEGRA = 2

	BLOCK_TYPE_ALLEGRA = 3

	BLOCK_HEADER_TYPE_ALLEGRA = 2

	TX_TYPE_ALLEGRA = 2
)

type AllegraBlock struct {
	cbor.StructAsArray
	Header                 *AllegraBlockHeader
	TransactionBodies      []AllegraTransactionBody
	TransactionWitnessSets []ShelleyTransactionWitnessSet
	TransactionMetadataSet map[uint]cbor.Value
}

func (b *AllegraBlock) Hash() string {
	return b.Header.Hash()
}

func (b *AllegraBlock) BlockNumber() uint64 {
	return b.Header.BlockNumber()
}

func (b *AllegraBlock) SlotNumber() uint64 {
	return b.Header.SlotNumber()
}

func (b *AllegraBlock) Era() Era {
	return eras[ERA_ID_ALLEGRA]
}

func (b *AllegraBlock) Transactions() []Transaction {
	// TODO
	return nil
}

type AllegraBlockHeader struct {
	ShelleyBlockHeader
}

func (h *AllegraBlockHeader) Era() Era {
	return eras[ERA_ID_ALLEGRA]
}

type AllegraTransactionBody struct {
	ShelleyTransactionBody
	ValidityIntervalStart uint64 `cbor:"8,keyasint,omitempty"`
}

type AllegraTransaction struct {
	cbor.StructAsArray
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
