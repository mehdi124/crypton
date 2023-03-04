package connection

import (
	token "github.com/crypton/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Connect(url) (ethclient.Client, error) {

	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	return client, nil
}

type TokenContract struct {
	Name     string
	Symbol   string
	Decimals int
}

func GetTokenInstance(tokenContract string) (*token.Token, error) {

	client, err := connection.Connect("localhost:4545")
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(tokenContract)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func SetTokenInfo(tokenContract string) (*TokenContract, error) {

	instance, err := GetTokenInstance(tokenContract)
	if err != nil {
		return nil, err
	}

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	return &TokenContract{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}, nil

}
