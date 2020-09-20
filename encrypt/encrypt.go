package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/jmattheis/shargo/proto"
	"gopkg.in/mgo.v2/bson"
)

func Sha256(s string) []byte {
	x := sha256.Sum256([]byte(s))
	return x[:]
}

func Encrypt(key []byte, p *proto.Packet) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	data, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func Decrypt(key, data []byte) (*proto.Packet, error) {

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create cipher: %s", err)
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("could not create gcm: %s", err)
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open with gcm: %s", err)
	}
	var p proto.Packet

	err = bson.Unmarshal(plaintext, &p)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal plaintext: %s", err)
	}

	return &p, nil
}
