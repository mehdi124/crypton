package storage

import (
	"os"

	"fmt"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"
const Dir = "./wallets/"

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

	//TODO check name and address exists or not ???
	db, err := New(Dir, nil)
	if err != nil {
		fmt.Println("Error", err)
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
