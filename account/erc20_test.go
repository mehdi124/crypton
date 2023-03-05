package account

import (
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
)

func TestImportAndCreateAccount(t *testing.T) {

	erc20Account, err := Create()
	assert.Nil(t, err)
	log.Println("private key", erc20Account.Export())

}
