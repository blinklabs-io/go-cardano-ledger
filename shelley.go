package ledger

import (
	"fmt"
	"github.com/cloudstruct/go-cardano-ledger/cbor"
)

const (
	BLOCK_TYPE_SHELLEY = 2

	BLOCK_HEADER_TYPE_SHELLEY = 1

	TX_TYPE_SHELLEY = 1
)

type ShelleyBlock struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_                      struct{} `cbor:",toarray"`
	Header                 ShelleyBlockHeader
	TransactionBodies      []ShelleyTransactionBody
	TransactionWitnessSets []ShelleyTransactionWitnessSet
	TransactionMetadataSet map[uint]cbor.Value
}

func (b *ShelleyBlock) Id() string {
	return b.Header.Id()
}

type ShelleyBlockHeader struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_    struct{} `cbor:",toarray"`
	id   string
	Body struct {
		// Tells the CBOR decoder to convert to/from a struct and a CBOR array
		_                    struct{} `cbor:",toarray"`
		BlockNumber          uint64
		Slot                 uint64
		PrevHash             Blake2b256
		IssuerVkey           interface{}
		VrfKey               interface{}
		NonceVrf             interface{}
		LeaderVrf            interface{}
		BlockBodySize        uint32
		BlockBodyHash        Blake2b256
		OpCertHotVkey        interface{}
		OpCertSequenceNumber uint32
		OpCertKesPeriod      uint32
		OpCertSignature      interface{}
		ProtoMajorVersion    uint64
		ProtoMinorVersion    uint64
	}
	Signature interface{}
}

func (h *ShelleyBlockHeader) Id() string {
	return h.id
}

type ShelleyTransactionBody struct {
	Inputs  []ShelleyTransactionInput  `cbor:"0,keyasint,omitempty"`
	Outputs []ShelleyTransactionOutput `cbor:"1,keyasint,omitempty"`
	Fee     uint64                     `cbor:"2,keyasint,omitempty"`
	Ttl     uint64                     `cbor:"3,keyasint,omitempty"`
	// TODO: figure out how to parse properly
	Certificates []cbor.Value `cbor:"4,keyasint,omitempty"`
	// TODO: figure out how to parse this correctly
	// We keep the raw CBOR because it can contain a map with []byte keys, which
	// Go does not allow
	Withdrawals cbor.Value `cbor:"5,keyasint,omitempty"`
	Update      struct {
		// Tells the CBOR decoder to convert to/from a struct and a CBOR array
		_                    struct{} `cbor:",toarray"`
		ProtocolParamUpdates cbor.Value
		Epoch                uint64
	} `cbor:"6,keyasint,omitempty"`
	MetadataHash Blake2b256 `cbor:"7,keyasint,omitempty"`
}

type ShelleyTransactionInput struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_     struct{} `cbor:",toarray"`
	Id    Blake2b256
	Index uint32
}

type ShelleyTransactionOutput struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_       struct{} `cbor:",toarray"`
	Address Blake2b256
	Amount  uint64
}

type ShelleyTransactionWitnessSet struct {
	VkeyWitnesses      []interface{} `cbor:"0,keyasint,omitempty"`
	MultisigScripts    []interface{} `cbor:"1,keyasint,omitempty"`
	BootstrapWitnesses []interface{} `cbor:"2,keyasint,omitempty"`
}

type ShelleyTransaction struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_          struct{} `cbor:",toarray"`
	Body       ShelleyTransactionBody
	WitnessSet ShelleyTransactionWitnessSet
	Metadata   cbor.Value
}

func NewShelleyBlockFromCbor(data []byte) (*ShelleyBlock, error) {
	var shelleyBlock ShelleyBlock
	if _, err := cbor.Decode(data, &shelleyBlock); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	rawBlockHeader, err := extractHeaderFromBlockCbor(data)
	if err != nil {
		return nil, err
	}
	shelleyBlock.Header.id, err = generateBlockHeaderHash(rawBlockHeader, nil)
	return &shelleyBlock, err
}

func NewShelleyBlockHeaderFromCbor(data []byte) (*ShelleyBlockHeader, error) {
	var err error
	var shelleyBlockHeader ShelleyBlockHeader
	if _, err := cbor.Decode(data, &shelleyBlockHeader); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	shelleyBlockHeader.id, err = generateBlockHeaderHash(data, nil)
	return &shelleyBlockHeader, err
}

func NewShelleyTransactionBodyFromCbor(data []byte) (*ShelleyTransactionBody, error) {
	var shelleyTx ShelleyTransactionBody
	if _, err := cbor.Decode(data, &shelleyTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &shelleyTx, nil
}

func NewShelleyTransactionFromCbor(data []byte) (*ShelleyTransaction, error) {
	var shelleyTx ShelleyTransaction
	if _, err := cbor.Decode(data, &shelleyTx); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}
	return &shelleyTx, nil
}
