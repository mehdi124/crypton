package account

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mehdi124/crypton/connection"
)

type Erc20Account struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func Import(hexPrivateKey string) (*Erc20Account, error) {

	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
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
	}, nil
}

type request struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

func (account *Erc20Account) Export() string {

	privateKeyBytes := crypto.FromECDSA(account.privateKey)
	return hexutil.Encode(privateKeyBytes)[2:]

}

func (account *Erc20Account) Address() string {
	return crypto.PubkeyToAddress(*account.publicKey).Hex()
}

func (account *Erc20Account) ETHBalance() (*big.Float, error) {

	client, err := connection.Connect("localhost:4545")
	if err != nil {
		return new(big.Float), err
	}

	hexAddress := account.Address()
	acc := common.HexToAddress(hexAddress)
	balance, err := client.BalanceAt(context.Background(), acc, nil)
	if err != nil {
		return new(big.Float), err
	}

	return convertEthValue(balance), nil
}

func convertEthValue(balanceAt *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(balanceAt.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return ethValue
}

func (account *Erc20Account) TokenBalance(contractAddr string) (*big.Float, error) {

	client, err := connection.ConnectHttp("mainnet")
	if err != nil {
		return new(big.Float), err
	}

	defer client.Close()

	address := account.Address()
	data := "0x70a08231" + fmt.Sprintf("%064s", address[2:])

	req := request{contractAddr, data}
	var resp string
	if err := client.Call(&resp, "eth_call", req, "latest"); err != nil {
		return new(big.Float), err
	}

	balance, err := hexutil.DecodeBig("0x" + strings.TrimLeft(resp[2:], "0")) // %064s means that the string is padded with 0 to 64 bytes
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance)

	return convertTokenValue(balance, 6), nil
}

func convertTokenValue(balance *big.Int, decimals int) *big.Float {
	fbal := new(big.Float)
	fbal.SetString(balance.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	return value
}
