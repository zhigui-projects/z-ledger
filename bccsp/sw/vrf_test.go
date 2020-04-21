package sw

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEcdsaPrivateKeyVrf(t *testing.T) {
	t.Parallel()

	vrfer := &ecdsaPrivateKeyVrf{}
	verifierPublicKey := &ecdsaPublicKeyKeyVrfVerifier{}

	// Generate a key
	lowLevelKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	k := &ecdsaPrivateKey{lowLevelKey}
	pk, err := k.PublicKey()
	assert.NoError(t, err)

	// Sign
	msg := []byte("Hello World")
	rand, proof, err := vrfer.Vrf(k, msg)
	assert.NoError(t, err)
	assert.NotNil(t, rand)
	assert.NotNil(t, proof)

	// Verify
	valid, err := verifierPublicKey.VrfVerify(pk, msg, rand, proof)
	assert.NoError(t, err)
	assert.True(t, valid)
}
