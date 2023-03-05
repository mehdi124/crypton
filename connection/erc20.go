package connection

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

var APIKey = "e926ac6aae5543f099859ad3a9293081"
var URL string

const (
	Mainnet string = "mainnet"
	Goerli         = "goerli"
	Sepolia        = "sepolia"
)

func GetInfuraUrl(network string) string {
	URL = "https://" + network + ".infura.io/v3/" + APIKey
	return URL
}

func ConnectHttp(network string) (*ethclient.Client, error) {

	GetInfuraUrl(network)
	URL = "http://localhost:7545"
	client, err := ethclient.Dial(URL)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Connect(network string) (*ethclient.Client, error) {

	GetInfuraUrl(network)
	URL = "http://localhost:7545"
	client, err := ethclient.Dial(URL)
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

/*func SetTokenInfo(tokenContract string) (*TokenContract, error) {

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

}*/
