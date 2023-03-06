package account

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"strings"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mehdi124/crypton/connection"
	"github.com/mehdi124/crypton/storage"
)

type Erc20Account struct {
	Name       string
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

type ethHandlerResult struct {
	Result string `json:"result"`
	Error  struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
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

func Create(name, password string) (*Erc20Account, error) {

	name = handleName(name)

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	pbKey := privateKey.Public()
	publicKey, ok := pbKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	erc20Account := &Erc20Account{
		Name:       name,
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	err = storage.Store(erc20Account.Name, erc20Account.Address(), erc20Account.PrivateKeyToHex(), password)
	if err != nil {
		return nil, err
	}

	return erc20Account, nil
}

func (account *Erc20Account) PrivateKeyToHex() string {
	privateKeyBytes := crypto.FromECDSA(account.privateKey)
	return hexutil.Encode(privateKeyBytes)[2:]
}

func (account *Erc20Account) Export(name, password string) (string, error) {

	pkv, err := storage.GetPrivateInfo(name, password)
	if err != nil {
		return "", err
	}

	return pkv, nil
}

func (account *Erc20Account) Address() string {
	return crypto.PubkeyToAddress(*account.publicKey).Hex()
}

func (account *Erc20Account) ETHBalance() (*big.Float, error) {

	client, err := connection.Connect("mainnet")
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

	//address := account.Address()
	address := "0xc09cdd3874d2dabc80696ad21d2bf7441fda7832"

	data := crypto.Keccak256Hash([]byte("balanceOf(address)")).String()[0:10] + "000000000000000000000000" + address[2:]
	postBody, _ := json.Marshal(map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "eth_call",
		"params": []interface{}{
			map[string]string{
				"to":   contractAddr,
				"data": data,
			},
			"latest",
		},
	})

	requestBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(connection.GetInfuraUrl("mainnet"), "application/json", requestBody)
	if err != nil {
		return nil, err
	}

	ethResult := new(big.Int)

	if err := json.NewDecoder(resp.Body).Decode(&ethResult); err != nil {
		return nil, err
	}

	//ethResult.SetString(ethResult.Result[2:], 16)
	fmt.Println(ethResult, "ss", resp.Body)
	return new(big.Float), nil

	//return convertTokenValue(balance, 6), nil
}

func convertTokenValue(balance *big.Int, decimals int) *big.Float {
	fbal := new(big.Float)
	fbal.SetString(balance.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	return value
}

func (account *Erc20Account) ETHTransfer(to string, amount *big.Int) (string, error) {

	client, err := connection.Connect("rinkeby")
	if err != nil {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*account.publicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	toAddress := common.HexToAddress(to)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), account.privateKey)
	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}

func (account *Erc20Account) TokenTransfer(contract string, to string, amount *big.Int) (string, error) {

	client, err := connection.Connect("rinkeby")
	if err != nil {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*account.publicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	toAddress := common.HexToAddress(to)
	tokenAddress := common.HexToAddress(contract)

	transferFnSignature := []byte("transfer(address,uint256)")

	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d
	amount = new(big.Int)
	amount.SetString("1000000000000000000000", 10) // 1000 tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gasLimit) // 23256

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), account.privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}

func handleName(name string) string {

	name = strings.ReplaceAll(name, " ", "")
	name = strings.TrimSpace(name)
	return name
}
