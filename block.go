package ledger

import (
	"encoding/hex"
	"fmt"
	"github.com/cloudstruct/go-cardano-ledger/cbor"
	"golang.org/x/crypto/blake2b"
)

type Blake2b256 [32]byte

func (b Blake2b256) String() string {
	return hex.EncodeToString([]byte(b[:]))
}

type Blake2b224 [28]byte

func (b Blake2b224) String() string {
	return hex.EncodeToString([]byte(b[:]))
}

func NewBlockFromCbor(blockType uint, data []byte) (interface{}, error) {
	switch blockType {
	case BLOCK_TYPE_BYRON_EBB:
		return NewByronEpochBoundaryBlockFromCbor(data)
	case BLOCK_TYPE_BYRON_MAIN:
		return NewByronMainBlockFromCbor(data)
	case BLOCK_TYPE_SHELLEY:
		return NewShelleyBlockFromCbor(data)
	case BLOCK_TYPE_ALLEGRA:
		return NewAllegraBlockFromCbor(data)
	case BLOCK_TYPE_MARY:
		return NewMaryBlockFromCbor(data)
	case BLOCK_TYPE_ALONZO:
		return NewAlonzoBlockFromCbor(data)
	case BLOCK_TYPE_BABBAGE:
		return NewBabbageBlockFromCbor(data)
	}
	return nil, fmt.Errorf("unknown node-to-client block type: %d", blockType)
}

// XXX: should this take the block header type instead?
func NewBlockHeaderFromCbor(blockType uint, data []byte) (interface{}, error) {
	switch blockType {
	case BLOCK_TYPE_BYRON_EBB:
		return NewByronEpochBoundaryBlockHeaderFromCbor(data)
	case BLOCK_TYPE_BYRON_MAIN:
		return NewByronMainBlockHeaderFromCbor(data)
	case BLOCK_TYPE_SHELLEY, BLOCK_TYPE_ALLEGRA, BLOCK_TYPE_MARY, BLOCK_TYPE_ALONZO:
		return NewShelleyBlockHeaderFromCbor(data)
	case BLOCK_TYPE_BABBAGE:
		return NewBabbageBlockHeaderFromCbor(data)
	}
	return nil, fmt.Errorf("unknown node-to-node block type: %d", blockType)
}

func generateBlockHeaderHash(data []byte, prefix []byte) (string, error) {
	tmpHash, err := blake2b.New256(nil)
	if err != nil {
		return "", err
	}
	if prefix != nil {
		tmpHash.Write(prefix)
	}
	tmpHash.Write(data)
	return hex.EncodeToString(tmpHash.Sum(nil)), nil
}

func extractHeaderFromBlockCbor(data []byte) ([]byte, error) {
	// Parse outer list to get at header CBOR
	var rawBlock []cbor.RawMessage
	if _, err := cbor.Decode(data, &rawBlock); err != nil {
		return nil, err
	}
	return []byte(rawBlock[0]), nil
}
