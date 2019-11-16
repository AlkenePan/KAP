package crypto

import (
	"fmt"
	"testing"
)

func TestCrypto(t *testing.T) {
	//fmt.Println("test RSA")
	priv, pub := GenerateKeyPair(2048)
	fmt.Println("pub key: ", string(PublicKeyToBytes(pub)))
	fmt.Println("priv key: ", string(PrivateKeyToBytes(priv)))
	plaintext := []byte("gaba gaba gaba hey")
	plaintextBase64Encoded := BytesBase64Encode(plaintext)
	ciphertext := EncryptWithPublicKey(plaintextBase64Encoded, pub)
	ciphertextBase64Encoded := BytesBase64Encode(ciphertext)
	ciphertextBase64Decoded := BytesBase64Decode(ciphertextBase64Encoded)

	decryptedtext := DecryptWithPrivateKey(ciphertextBase64Decoded, priv)
	decryptedtextBase64Decoded := BytesBase64Decode(decryptedtext)

	fmt.Println(string(plaintext))
	fmt.Println(string(ciphertextBase64Encoded))
	fmt.Println(string(decryptedtextBase64Decoded))
	if string(plaintext) != string(decryptedtextBase64Decoded) {
		t.Error(
			"expected", plaintext,
			"got", decryptedtextBase64Decoded,
			)
	}
}
