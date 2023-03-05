package storage

import (
	"os"

	"github.com/mehdi124/crypton/account"

	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

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

func main() {

	publicDir := "./storage/wallets/public/"
	privateDir := "./storage/wallets/private/"
	password := "test"

	dbPub, err := New(publicDir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	dbPk, err := New(privateDir, nil)
	if err != nil {
		fmt.Println("error 2 ", err)
	}

	pubInfo := NewPublicInfo()

	dbPub.Write("users", account.Erc20Account.Name, value)

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
