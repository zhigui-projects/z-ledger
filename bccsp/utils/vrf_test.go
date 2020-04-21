package utils

import (
	"testing"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/vrf"
)

func computeVrf(t *testing.T, kt keypair.KeyType, curve byte) {
	pri, _, err := keypair.GenerateKeyPair(kt, curve)
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("test")
	v, _, err := vrf.Vrf(pri, msg)
	if err != nil {
		t.Fatalf("compute vrf: %v", err)
	}
	ret := CalcEndorser(v, 5, 2)
	t.Logf("CalcEndorser ret: %v", ret)
}

func TestCalcEndorser(t *testing.T) {
	computeVrf(t, keypair.PK_ECDSA, keypair.SECP256K1)
	computeVrf(t, keypair.PK_ECDSA, keypair.P224)
	computeVrf(t, keypair.PK_ECDSA, keypair.P256)
	computeVrf(t, keypair.PK_ECDSA, keypair.P384)
	computeVrf(t, keypair.PK_SM2, keypair.SM2P256V1)
}
