package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func CalcEndorser(value []byte, candidate, threshold int64) bool {
	h := sha256.Sum256(value)
	i := new(big.Int)
	i.SetString(hex.EncodeToString(h[:]), 16)
	i.Mul(i, big.NewInt(candidate))
	i.Rsh(i, 256)
	fmt.Println("CalcEndorser ret:", i)
	return i.Cmp(big.NewInt(threshold)) >= 0
}
