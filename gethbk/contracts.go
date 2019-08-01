package gethbk

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/dlm/go-eth-book/store"
)

func SmartContractsDeploy(client *ethclient.Client, privateKeyHex string) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	checkForError(err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	checkOk(ok, "fail asserting type")

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	checkForError(err)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	checkForError(err)

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address, tx, instance, err := store.DeployStore(auth, client)

	fmt.Println(address.Hex())
	fmt.Println(tx.Hash().Hex())

	_ = instance
}

func SmartContractsLoad(client *ethclient.Client) *store.Store {
	// output from previous section
	addressHex := "0xD2Bb7fF4Aa4ce4EA42C0278eb993107E12f2c391"
	txHashHex := "0xe1ac7943dbc70316081d51dc278242f77bfef4a2c2ca139e09909c1cd8125205"
	_ = txHashHex

	address := common.HexToAddress(addressHex)
	instance, err := store.NewStore(address, client)
	checkForError(err)

	return instance

}

func SmartContractsQuery(client *ethclient.Client, store *store.Store) {
	version, err := store.Version(nil)
	checkForError(err)
	fmt.Println(version)
}
