package account

import (
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
)

func createAccount(t *testing.T) *Erc20Account {
	erc20Account, err := Create("testaccount", "test")
	assert.Nil(t, err)
	log.Println("private key", erc20Account.Export())
	return erc20Account
}

func TestImportAndCreateAccount(t *testing.T) {

	erc20Account := createAccount(t)
	//importErc20Account, err := Import(erc20Account.Export())
	//assert.Equal(t, erc20Account.Export(), importErc20Account.Export())
	//assert.Equal(t, erc20Account.Address(), importErc20Account.Address())
	log.Println(erc20Account.Export(), erc20Account.Address())
}

/*func TestBalance(t *testing.T) {

	erc20Account := createAccount(t)
	balance, err := erc20Account.ETHBalance()
	assert.Nil(t, err)
	log.Println(balance, "balance")

	balance, err = erc20Account.TokenBalance("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	//assert.Nil(t, err)
	log.Println("token balance", balance, err)

}*/
