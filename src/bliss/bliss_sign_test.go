package bliss

import (
	_ "fmt"
	_ "io/ioutil"
	"sampler"
	_ "strconv"
	_ "strings"
	"testing"
)

func TestSignVerify(t *testing.T) {
	for i := 0; i <= 4; i++ {
		seed := make([]uint8, sampler.SHA_512_DIGEST_LENGTH)
		for i := 0; i < len(seed); i++ {
			seed[i] = uint8(i % 8)
		}
		entropy, err := sampler.NewEntropy(seed)
		if err != nil {
			t.Errorf("Error in initializing entropy: %s", err.Error())
		}

		key, err := GeneratePrivateKey(i, entropy)
		if err != nil {
			t.Errorf("Error in generating private key: %s", err.Error())
		}

		pub := key.PublicKey()
		msg := []byte("Hello world")
		sig, err := key.Sign(msg, entropy)
		if err != nil {
			t.Errorf("Failed to generate signature for version %d: %s", i, err.Error())
		}
		_, err = pub.Verify(msg, sig)
		if err != nil {
			t.Errorf("Failed to verify signature for version %d: %s", i, err.Error())
		}
	}
}
