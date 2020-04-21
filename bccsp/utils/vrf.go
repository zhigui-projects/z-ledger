package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

func CalcEndorser(value []byte, candidate, threshold int64) (bool, *big.Int) {
	h := sha256.Sum256(value)
	i := new(big.Int)
	i.SetString(hex.EncodeToString(h[:]), 16)
	i.Mul(i, big.NewInt(candidate))
	i.Rsh(i, 256)
	return i.Cmp(big.NewInt(threshold)) >= 0, i
}
