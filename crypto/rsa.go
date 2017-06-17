package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
)

func Sign(hash []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
}

func Verify(sig, hash []byte, publicKey *rsa.PublicKey) error {
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, sig)
}

// GenRsaFile generate RSA Private Key and store it in a file
func GenRsaFile(path string) error {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey),
		})
	return ioutil.WriteFile(path, pemdata, 0644)
}

// OpenRsaFile open and retrieve RSA private key stored at given path
func OpenRsaFile(path string) (*rsa.PrivateKey, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(f)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// RsaID open rsa file at given path and retrieves owner ID
func RsaID(path string) (string, error) {
	// rsa key
	rsaPrivateKey, err := OpenRsaFile(path)
	if err != nil {
		return "", err
	}

	// hash private key to get id
	hash := sha256.New()
	io.WriteString(hash, fmt.Sprintf("%v", rsaPrivateKey))
	localID := fmt.Sprintf("%x", hash.Sum(nil))
	return localID, nil
}
