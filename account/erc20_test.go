package account

import (
	"testing"

	"log"

	"github.com/mehdi124/crypton/storage"
	"github.com/stretchr/testify/assert"
)

func createAccount(t *testing.T, name, password string) (*Erc20Account, error) {
	erc20Account, err := Create(name, password)
	return erc20Account, err
}

func TestImportAndCreateAccount(t *testing.T) {

	erc20Account, err := createAccount(t, "account2", "test")
	log.Println(err, "error")
	assert.Nil(t, err)

	//importErc20Account, err := Import(erc20Account.Export())
	//assert.Equal(t, erc20Account.Export(), importErc20Account.Export())
	//assert.Equal(t, erc20Account.Address(), importErc20Account.Address())

	pubInfos, err := storage.GetWalletsList()
	assert.Nil(t, err)
	log.Println(pubInfos)

	pvk, err := storage.GetPrivateInfo("account2", "test")
	assert.Nil(t, err)

	log.Println("pvk", pvk, erc20Account.PrivateKeyToHex())
	assert.Equal(t, erc20Account.PrivateKeyToHex(), pvk)
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
