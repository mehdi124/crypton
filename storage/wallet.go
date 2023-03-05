package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"os"

	"github.com/mehdi124/crypton/account"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

func Encrypt(key, data []byte) ([]byte, error) {
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

func Decrypt(key, data []byte) ([]byte, error) {
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

const Version = "1.0.0"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {

	filepath.Clean(dir)
	opts := Options{}
	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	if _, err := os.Stat(dir); err != nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating db at '%s' ...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)

}

func (d *Driver) Write(collection, resource string, v []byte) error {

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

	if err := ioutil.WriteFile(tempPath, v, 0644); err != nil {
		return nil
	}

	return os.Rename(tempPath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {

	if collection == "" {
		return fmt.Errorf("missing collection - unable to read")
	}

	if resource == "" {
		return fmt.Errorf("missing resource - unable to read records (no name)")
	}

	record := filepath.Join(d.dir, collection, resource)

	if _, err := stat(record); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(record + ".txt")
	if err != nil {
		return err
	}

	return hex.EncodeToString(b)

}

func (d *Driver) ReadAll(collection string) ([]string, error) {

	if collection == "" {
		return nil, fmt.Errorf("missing collection name - unable to read")
	}

	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := ioutil.ReadDir(dir)
	var records []string

	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))

		if err != nil {
			return nil, err
		}

		records = append(records, string(b))

	}

	return records, nil
}

func (d *Driver) Delete(collection, resource string) error {

	path := filepath.Join(collection, resource)

	mutex := d.getOrCreateMutext(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {

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

func main() {

	dir := "./storage/wallets/"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println(account.Erc20Account)
	db.Write("users", account.Erc20Account.Name, value)

	records, err := db.ReadAll("wallets")
	if err != nil {
		fmt.Println("Error read all", err)
	}

	fmt.Println(records)

	allusers := []account.Erc20Account{}
	for _, f := range records {

		employeeFund := User{}
		if err := json.Unmarshal([]byte(f), &employeeFund); err != nil {
			fmt.Println("unmarshal error", err)
		}

		allusers = append(allusers, employeeFund)

	}

	fmt.Println(allusers)

	//if err := db.Delete("users", ""); err != nil {
	//fmt.Println("delete error", err)
	//}

}
