package ledger

import (
	"fmt"
)

func NewTransactionFromCbor(txType uint, data []byte) (interface{}, error) {
	switch txType {
	case TX_TYPE_BYRON:
		return NewByronTransactionFromCbor(data)
	case TX_TYPE_SHELLEY:
		return NewShelleyTransactionFromCbor(data)
	case TX_TYPE_ALLEGRA:
		return NewAllegraTransactionFromCbor(data)
	case TX_TYPE_MARY:
		return NewMaryTransactionFromCbor(data)
	case TX_TYPE_ALONZO:
		return NewAlonzoTransactionFromCbor(data)
	case TX_TYPE_BABBAGE:
		return NewBabbageTransactionFromCbor(data)
	}
	return nil, fmt.Errorf("unknown transaction type: %d", txType)
}
