// Package eccpoint implements transposition on points and skalars.
package eccpoint

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/decred/dcrd/dcrec/secp256k1"
)

func hash(d []byte) []byte {
	r := sha256.Sum256(d)
	return r[:]
}

func calcTranspose(transpose []byte, try []byte) *secp256k1.ModNScalar {
	mac := hmac.New(sha256.New, transpose)
	mac.Write(try)
	r := new(secp256k1.ModNScalar)
	r.SetByteSlice(mac.Sum(nil))
	return r
}

func transposePublicKey(pubKey *secp256k1.PublicKey, transpose []byte) (*secp256k1.PublicKey, *secp256k1.ModNScalar) {
	curve := secp256k1.S256()
	pubKeyS := pubKey.SerializeCompressed()
	try := make([]byte, len(pubKeyS))
	copy(try, pubKeyS)
	for c := 0; c < 1000; c++ {
		r := new(secp256k1.PublicKey)
		transpose := calcTranspose(transpose, try)
		transposeB := transpose.Bytes()
		r.X, r.Y = curve.ScalarMult(pubKey.X, pubKey.Y, transposeB[:])
		if curve.IsOnCurve(r.X, r.Y) {
			return r, transpose
		}
		try = hash(try)
	}
	panic("secp256k1: TransposePublicKey could not find point")
	return nil, nil
}

func transposePrivateKey(pubKey *secp256k1.PublicKey, privateKey *secp256k1.PrivateKey, transpose []byte) (*secp256k1.PublicKey, *secp256k1.PrivateKey) {
	pub, transposeC := transposePublicKey(pubKey, transpose)
	priv := new(secp256k1.ModNScalar)
	priv.SetByteSlice(privateKey.Serialize())
	return pub, secp256k1.NewPrivateKey(transposeC.Mul(priv))
}
