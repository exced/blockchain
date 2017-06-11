package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
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

func OpenRsaFile(path string) (*rsa.PrivateKey, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(f)
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
