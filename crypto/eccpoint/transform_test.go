package eccpoint

import (
	"bytes"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1"
)

func TestTranspose(t *testing.T) {
	transpose := []byte("Random Secret Value")

	clientPrivateKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		t.Fatalf("GeneratePrivateKey: %s", err)
	}
	clientPublicKey := clientPrivateKey.PubKey()

	pubTranspose1, _ := transposePublicKey(clientPublicKey, transpose)
	if clientPublicKey.X.Cmp(pubTranspose1.X) == 0 || clientPublicKey.Y.Cmp(pubTranspose1.Y) == 0 {
		t.Error("No transposition")
	}
	pubTranspose2, privTranspose2 := transposePrivateKey(clientPublicKey, clientPrivateKey, transpose)
	if pubTranspose2.X.Cmp(pubTranspose1.X) != 0 || pubTranspose2.Y.Cmp(pubTranspose1.Y) != 0 {
		t.Error("Transposition Error Public Key")
	}
	if bytes.Equal(privTranspose2.Serialize(), clientPrivateKey.Serialize()) {
		t.Error("Private key no transposition")
	}
	pubTranspose3 := privTranspose2.PubKey()
	if pubTranspose3.X.Cmp(pubTranspose1.X) != 0 || pubTranspose3.Y.Cmp(pubTranspose1.Y) != 0 {
		t.Error("Transposition Error Private Key")
	}
}
