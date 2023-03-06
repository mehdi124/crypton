package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type PrivateInfo struct {
	PrivateKey string
}

func NewPrivateInfo(privateKey string) *PrivateInfo {
	return &PrivateInfo{
		PrivateKey: privateKey,
	}
}

func (d *Driver) WritePrivateInfo(collection, resource, password string, pk *PrivateInfo) error {

	if collection == "" {
		return fmt.Errorf("missing collection - no place to save records")
	}

	if resource == "" {
		return fmt.Errorf("missing resource - unable to save records (no name)")
	}

	mutex := d.getOrCreateMutext(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".txt")
	tempPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0775); err != nil {
		return err
	}

	ps := hashPassword(password)

	pkBytes := pk.Bytes()

	b, err := encrypt(ps, pkBytes)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(tempPath, b, 0644); err != nil {
		return nil
	}

	return os.Rename(tempPath, fnlPath)
}

func (d *Driver) ReadPrivateInfo(collection, resource, password string) (*PrivateInfo, error) {

	if collection == "" {
		return nil, fmt.Errorf("missing collection - unable to read")
	}

	if resource == "" {
		return nil, fmt.Errorf("missing resource - unable to read records (no name)")
	}

	record := filepath.Join(d.dir, collection, resource)

	if _, err := privateStat(record); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(record + ".txt")
	if err != nil {
		return nil, err
	}

	ps := hashPassword(password)

	decrypted, err := decrypt(ps, b)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(decrypted)
	dec := gob.NewDecoder(buf)

	pk := PrivateInfo{}

	if err := dec.Decode(&pk); err != nil {
		log.Fatal(err, pk)
	}

	//return json.Unmarshal(b, &pkInfo)
	return &pk, nil
}

func (d *Driver) DeletePrivateInfo(collection, resource string) error {

	path := filepath.Join(collection, resource)

	mutex := d.getOrCreateMutext(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := privateStat(dir); {

	case fi == nil, err != nil:
		return fmt.Errorf("unable to find file or directory name %v\n", path)

	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}

	return nil
}

func (d *Driver) getOrCreateMutext(collection string) *sync.Mutex {

	d.mutex.Lock()
	defer d.mutex.Unlock()

	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func stat(path string) (fi os.FileInfo, err error) {

	if fi, err = os.Stat(path); os.IsNotExist(err) {

		fi, err = os.Stat(path + ".json")

	}

	return fi, err

}

func privateStat(path string) (fi os.FileInfo, err error) {

	if fi, err = os.Stat(path); os.IsNotExist(err) {

		fi, err = os.Stat(path + ".txt")

	}

	return fi, err

}

func encrypt(key, data []byte) ([]byte, error) {
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

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func decrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (pk *PrivateInfo) Bytes() []byte {

	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(pk)
	return buf.Bytes()
}

func hashPassword(password string) []byte {
	ps := sha256.Sum256([]byte(password))
	log.Println(ps)
	return ps[:32]
}
