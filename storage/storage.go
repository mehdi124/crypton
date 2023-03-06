package storage

import (
	"encoding/json"
	"os"

	"fmt"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"
const Dir = "/home/mehdi/go/src/crypton/wallets/"

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

func Store(name, address, privateKey, password string) error {

	err := WalletExist(name, address)
	if err != nil {
		return err
	}

	db, err := New(Dir, nil)
	if err != nil {
		return err
	}

	pubInfo := NewPublicInfo(name, address)
	privInfo := NewPrivateInfo(privateKey)

	err = db.WritePublicInfo("public", name, pubInfo)
	if err != nil {
		return err
	}

	err = db.WritePrivateInfo("private", name, password, privInfo)
	if err != nil {
		return err
	}

	return nil
}

func WalletExist(name, address string) error {

	db, err := New(Dir, nil)
	if err != nil {
		return err
	}

	records, err := db.ReadAllPublicInfo("public")
	if err != nil {
		return nil
	}

	for _, record := range records {

		pubInfo := PublicInfo{}
		if err := json.Unmarshal([]byte(record), &pubInfo); err != nil {
			fmt.Errorf("unmarshal error", err)
		}

		if pubInfo.Name == name {
			return fmt.Errorf("%s wallet name already exist", name)
		}

		if pubInfo.Address == address {
			return fmt.Errorf("%s address already exist", address)
		}
	}

	return nil

}

func GetWalletsList() ([]PublicInfo, error) {

	db, err := New(Dir, nil)
	if err != nil {
		fmt.Println("error", err)
	}

	records, err := db.ReadAllPublicInfo("public")
	if err != nil {
		fmt.Println("er", err)
	}

	wallets := []PublicInfo{}

	for _, record := range records {

		pubInfo := PublicInfo{}
		if err := json.Unmarshal([]byte(record), &pubInfo); err != nil {
			fmt.Errorf("unmarshal error", err)
		}

		wallets = append(wallets, pubInfo)
	}

	return wallets, nil

}

func GetPrivateInfo(name, password string) (string, error) {

	db, err := New(Dir, nil)
	if err != nil {
		return "", err
	}

	privateInfo, err := db.ReadPrivateInfo("private", name, password)
	if err != nil {
		return "", err
	}
	return privateInfo.PrivateKey, nil
}
