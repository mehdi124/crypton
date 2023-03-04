package account

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mehdi124/crypton/connection"
)

type Erc20Account struct {
	privateKey *ecdsa.PrivateKey
	publickKey *ecdsa.PublicKey
}

func Import(hexPrivateKey string) (*Erc20Account, error) {

	privatekey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}

	pbKey := privateKey.Public()
	publicKey, ok := pbKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	return &Erc20Account{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

func Create() (*Erc20Account, error) {

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	pbKey := privateKey.Public()
	publicKey, ok := pbKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	return &Erc20Account{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (account *Erc20Account) Export() string {

	privateKeyBytes := crypto.FromECDSA(account.privateKey)
	return hexutil.Encode(privateKeyBytes)[2:]

}

func (account *Erc20Account) Address() string {
	return crypto.PubkeyToAddress(account.publickKey).Hex()
}

func (account *Erc20Account) ETHBalance() (big.Float, error) {

	client, err := connection.Connect("localhost:4545")
	if err != nil {
		return new(big.Float), err
	}

	hexAddress := account.Address()
	acc := common.HexToAddress(hexAddress)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return new(big.Float), err
	}

	return convertEthValue(balance), nil
}

func convertEthValue(balanceAt string) big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(balanceAt.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return ethValue
}

func (account *Erc20Account) TokenBalance(contract string) (big.Float, err) {

	instance, err := connection.GetTokenInstance(contract)
	if err != nil {
		return new(big.Float), err
	}

	address := common.HexToAddress(account.Address())
	balance, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return new(big.Float), err
	}

	return convertTokenValue(balance, decimals), nil
}

func convertTokenValue(balance string, decimals int) big.Float {
	fbal := new(big.Float)
	fbal.SetString(balance.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	return value
}
