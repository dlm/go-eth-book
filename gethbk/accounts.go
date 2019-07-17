package gethbk

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"math"
	"math/big"
	"regexp"
	"strings"
)

func gweiToWei(balance *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	return new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
}

func logBalance(
	client *ethclient.Client,
	address common.Address,
	blockNumber *big.Int) {
	balance, err := client.BalanceAt(context.Background(), address, blockNumber)
	checkForError(err)

	ethValue := gweiToWei(balance)
	logInfo(ethValue)
}

func AccountsBalances(client *ethclient.Client, accountHex string) {
	address := common.HexToAddress(accountHex)
	logBalance(client, address, nil)

	blockNumber := big.NewInt(5532993)
	logBalance(client, address, blockNumber)
}

type Wallet struct {
	private *ecdsa.PrivateKey
}

func (k Wallet) Private() *ecdsa.PrivateKey {
	return k.private
}

func (k Wallet) Public() *ecdsa.PublicKey {
	publicKey := k.Private().Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	checkOk(ok, "cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	return publicKeyECDSA
}

func (k Wallet) PrivateAsString() string {
	bytes := crypto.FromECDSA(k.Private())
	return hexutil.Encode(bytes)[2:]
}

func (k Wallet) PublicAsString() string {
	bytes := crypto.FromECDSAPub(k.Public())
	return hexutil.Encode(bytes)[4:]
}

func (k Wallet) PublicToAddress() string {
	pub := k.Public()
	return strings.ToLower(crypto.PubkeyToAddress(*pub).Hex())
}

func (k Wallet) PublicToAddressFromSha3() string {
	bytes := crypto.FromECDSAPub(k.Public())

	hash := sha3.NewLegacyKeccak256()
	hash.Write(bytes[1:])
	return hexutil.Encode(hash.Sum(nil)[12:])
}

func MakeWallet() Wallet {
	privateKey, err := crypto.GenerateKey()
	checkForError(err)
	return Wallet{privateKey}
}

func AccountsGeneratingNewWallets() {
	key := MakeWallet()
	logInfo(key.PrivateAsString())
	logInfo(key.PublicAsString())
	logInfo(key.PublicToAddress())
	logInfo(key.PublicToAddressFromSha3())
}

func newKeyStore(dir string) *keystore.KeyStore {
	return keystore.NewKeyStore(
		dir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)
}

func createKeyStore(dir string, password string) {
	ks := newKeyStore(dir)
	account, err := ks.NewAccount(password)
	checkForError(err)
	logInfo(account.Address.Hex())
}

func importKeyStore(dir string, password string, jsonKeystore []byte) {
	ks := newKeyStore(dir)
	account, err := ks.Import(jsonKeystore, password, password)
	checkForError(err)
	logInfo(account.Address.Hex())
}

func exampleKeyStore() []byte {
	return []byte(`{"address":"657e61650a4902042e736b288ffe71fa3610b02a","crypto":{"cipher":"aes-128-ctr","ciphertext":"3b98d349bafef6e16120fd9198d35ec036552d035580dd18f446532b607439ca","cipherparams":{"iv":"06dae92c156ed069d9b015153211508d"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"d4d050ace2d33a0423d5e1e5c6f07dbb39387161b22632d03b5ad011d553f8c2"},"mac":"4884a6192e26d16a6b0d32431be9d5ce56bd64d6392d29d090cba017b8c022e5"},"id":"efc09b9b-5d2b-46d3-aa7b-40570ea2104a","version":3}`)
}

func AccountsKeystores() {
	password := "secret"
	createKeyStore("./wallets", password)
	importKeyStore("./wallets", password, exampleKeyStore())

}

func isEthAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

func isContract(client *ethclient.Client, hex string) bool {
	address := common.HexToAddress(hex)
	bytecode, err := client.CodeAt(context.Background(), address, nil)
	checkForError(err)
	return len(bytecode) > 0
}

func AddressCheck(client *ethclient.Client) {
	valid := "0x323b5d4c32345ced77393b3530b1eed0f346429d"
	logInfo("is valid:", isEthAddress(valid))

	invalid := "0xZYXb5d4c32345ced77393b3530b1eed0f346429d"
	logInfo("is valid:", isEthAddress(invalid))

	contract := "0xe41d2489571d322189246dafa5ebde1f4699f498"
	logInfo("is contract:", isContract(client, contract))

	account := "0x8e215d06ea7ec1fdb4fc5fd21768f4b34ee92ef4"
	logInfo("is contract:", isContract(client, account))
}

func AddressPlay(client *ethclient.Client) string {
	password := "secret"
	dir := "./wallets"

	// create the keystore and grab the wallet
	ks := newKeyStore(dir)
	if len(ks.Wallets()) != 1 {
		panic("Unexpected number of wallets")
	}
	wallet := ks.Wallets()[0]
	status, _ := wallet.Status()
	logInfo(status, len(wallet.Accounts()))

	if len(wallet.Accounts()) != 1 {
		panic("Unexpected number of accounts")
	}
	account := wallet.Accounts()[0]
	expectedHex := "0x0eeEabE10a62D7D7C3e7f31Fb9B0eCD9455F279b"
	logInfo(account.Address.Hex())
	logInfo("  ===> expected", expectedHex)

	// unlock the account in the wallet (which unlocks the account)
	ks.Unlock(account, password)
	status, _ = wallet.Status()
	logInfo(status, len(wallet.Accounts()))

	// get the balance
	balance, err := client.BalanceAt(context.Background(), account.Address, nil)
	checkForError(err)
	ethValue := gweiToWei(balance)
	logInfo(ethValue)

	// get the private key
	keyJson, err :=  ks.Export(account, password, password)
	privateKey, err := keystore.DecryptKey(keyJson, password)
	checkForError(err)

	key := Wallet{privateKey.PrivateKey}
	return key.PrivateAsString()
}


