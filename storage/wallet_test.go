package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/mehdi124/crypton/account"
	"github.com/stretchr/testify/assert"
)

func TestStoreWallet(t *testing.T) {

	ercAccount, err := Create("test", "test")
	assert.Nil(t, err)

	dir := "./storage/wallets/"

	db, err := New(dir, nil)

	log.Println("private key", ercAccount.Export())

	if err != nil {
		fmt.Println("Error", err)
	}

	db.Write("wallets", ercAccount.Name, ercAccount)

	records, err := db.ReadAll("wallets")
	if err != nil {
		fmt.Println("Error read all", err)
	}

	fmt.Println(records)

	allusers := []account.Erc20Account{}
	for _, f := range records {

		employeeFund := account.Erc20Account{}
		if err := json.Unmarshal([]byte(f), &employeeFund); err != nil {
			fmt.Println("unmarshal error", err)
		}

		allusers = append(allusers, employeeFund)

	}

	fmt.Println(allusers)

}
