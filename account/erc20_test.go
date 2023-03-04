package account

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ImportAndCreateAccountTest(t *testing.T) {

	erc20Account, err := account.Create()
	assert.Nil(t, err)
	fmt.Println("private key", erc20Account.Export())

}
