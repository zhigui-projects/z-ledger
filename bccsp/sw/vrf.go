package sw

import (
	"github.com/hyperledger/fabric/bccsp"
	"github.com/ontio/ontology-crypto/ec"
	"github.com/ontio/ontology-crypto/vrf"
)

type ecdsaPrivateKeyVrf struct{}

func (s *ecdsaPrivateKeyVrf) Vrf(k bccsp.Key, msg []byte) (rand, proof []byte, err error) {
	pri := k.(*ecdsaPrivateKey).privKey
	ecPri := &ec.PrivateKey{ec.ECDSA, pri}
	return vrf.Vrf(ecPri, msg)
}

type ecdsaPublicKeyKeyVrfVerifier struct{}

func (s *ecdsaPublicKeyKeyVrfVerifier) VrfVerify(k bccsp.Key, msg, rand, proof []byte) (bool, error) {
	pub := k.(*ecdsaPublicKey).pubKey
	ecPub := &ec.PublicKey{ec.ECDSA, pub}
	return vrf.Verify(ecPub, msg, rand, proof)
}
