package gethbk

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func getHeader(client *ethclient.Client, blockNumber *big.Int) *types.Header {
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	checkForError(err)
	return header
}

func getMostRecentHeader(client *ethclient.Client) *types.Header {
	return getHeader(client, nil)
}

func logHeader(header *types.Header) {
	logInfo(header.Number.String())
}

func getBlock(client *ethclient.Client, blockNumber *big.Int) *types.Block {
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	checkForError(err)
	return block
}

func logBlock(block *types.Block) {
	logInfo("Number:", block.Number().Uint64())
	logInfo("Time:", block.Time())
	logInfo("Difficulty: ", block.Difficulty().Uint64())
	logInfo("Hash:", block.Hash().Hex())
	logInfo("Num transactions:", len(block.Transactions()))
}

func TransactionsQueryingBlocks(client *ethclient.Client) {
	blockNumber := big.NewInt(5671744)

	header := getMostRecentHeader(client)
	logHeader(header)

	header = getHeader(client, blockNumber)
	logHeader(header)

	block := getBlock(client, blockNumber)
	logBlock(block)
}

func logTransaction(tx *types.Transaction) {
	logInfo("Hash:", tx.Hash().Hex())
	logInfo("Value:", tx.Value().String())
	logInfo("Gas:", tx.Gas())
	logInfo("Gas Price:", tx.GasPrice().Uint64())
	logInfo("Nonce:", tx.Nonce())
	logInfo("Data:", tx.Data())
	logInfo("To:", tx.To().Hex())
}

func getMessage(client *ethclient.Client, tx *types.Transaction) *types.Message {
	chainID, err := client.NetworkID(context.Background())
	checkForError(err)

	msg, err2 := tx.AsMessage(types.NewEIP155Signer(chainID))
	checkForError(err2)

	return &msg
}

func getReceipt(client *ethclient.Client, tx *types.Transaction) *types.Receipt {
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	checkForError(err)
	return receipt
}

func logReceipt(receipt *types.Receipt) {
	logInfo("Status:", receipt.Status)
	logInfo("Logs:", receipt.Logs)
}

func queryTransactionsOfBlock(client *ethclient.Client, blockNumber *big.Int) {
	block := getBlock(client, blockNumber)
	for _, tx := range block.Transactions() {
		logInfo("------------------------------------------")
		logTransaction(tx)

		msg := getMessage(client, tx)
		logInfo(msg)
		receipt := getReceipt(client, tx)
		logReceipt(receipt)
	}
}

func iterateOverTransactions(client *ethclient.Client, hex string) {
	blockHash := common.HexToHash(hex)
	count, err := client.TransactionCount(context.Background(), blockHash)
	checkForError(err)
	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		checkForError(err)
		logInfo(idx, tx.Hash().Hex())
	}
}

func queryTransactionByHex(client *ethclient.Client, hex string) {
	txHash := common.HexToHash(hex)
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	checkForError(err)

	logInfo("Hex", tx.Hash().Hex())
	logInfo("IsPending:", isPending)
}

func TransactionsQueryingTransactions(client *ethclient.Client) {
	blockNumber := big.NewInt(5671744)
	queryTransactionsOfBlock(client, blockNumber)

	blockHex := "0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9"
	iterateOverTransactions(client, blockHex)

	txHex := "0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2"
	queryTransactionByHex(client, txHex)

}

func TransactionsTransferringETH(client *ethclient.Client, privateKeyHex string) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	checkForError(err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	checkOk(ok, "cannot assert type: publicKey is not of type *ecdsa.PublicKey")

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	checkForError(err)

	value := big.NewInt(1000000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	checkForError(err)

	toHex := "0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d"
	toAddress := common.HexToAddress(toHex)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	checkForError(err)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	checkForError(err)

	err = client.SendTransaction(context.Background(), signedTx)
	checkForError(err)

	logInfo("tx sent: %s", signedTx.Hash().Hex())
}

func TransactionsSubscribingToNewBlocks() {
	client, err := ethclient.Dial("wss://ropsten.infura.io/ws")
	checkForError(err)

	headers := make(chan *types.Header)
	sub, err :=client.SubscribeNewHead(context.Background(), headers)
	checkForError(err)
	for {
		select {
		case err := <-sub.Err():
			logFatal(err)
		case header := <-headers:
			logInfo("-----------------------New Block--------------------")
			logInfo(header.Hash().Hex())
			block, err := client.BlockByHash(context.Background(), header.Hash())
			checkForError(err)
			logBlock(block)
		}
	}
}
