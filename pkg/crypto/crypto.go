package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
)

func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatalln(err)
	}
	return privKey, &privKey.PublicKey
}

func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		})
	return privBytes

}

func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Fatalln(err)
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	return pubBytes
}

func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatalln(err)
		}

	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Fatalln(err)
	}
	return key
}

func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Fatalln(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Fatalln("not ok")
	}

	return key
}

func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return ciphertext
}

func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return plaintext
}

func BytesBase64Encode(bytes []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
	base64.StdEncoding.Encode(dst, bytes)
	return dst
}
func BytesBase64Decode(bytes []byte) []byte {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(bytes)))
	n, err := base64.StdEncoding.Decode(dst, bytes)
	if err != nil {
		log.Fatalln(err)
	}
	return dst[:n]
}

