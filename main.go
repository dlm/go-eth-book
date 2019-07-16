package main

import (
	"log"

	"github.com/dlm/go-eth-book/gethbk"
)

const CONFIG_LOCAL string = "CONFIG_LOCAL"
const CONFIG_MAINNET string = "CONFIG_MAINNET"
const CONFIG string = CONFIG_MAINNET

func url(config string) string {
	address := "https://mainnet.infura.io/v3/c540c4fe10614f27a05a22bd463f0996"
	if config == CONFIG_LOCAL {
		address = "http://localhost:8545"
	}
	return address
}

func accountHex(config string) string {
	hex := "0x71c7656ec7ab88b098defb751b7401b5f6d8976f"
	if config == CONFIG_LOCAL {
		hex = "0xe280029a7867ba5c9154434886c241775ea87e53"
	}
	return hex
}

func main() {
	log.Println("Running with configuration ", CONFIG)

	client := gethbk.Client(url(CONFIG))

	gethbk.AccountsBalances(client, accountHex(CONFIG))
	gethbk.AccountsGeneratingNewWallets()
	gethbk.AccountsKeystores()
	gethbk.AddressCheck(client)

	gethbk.TransactionsQueryingBlocks(client)
	gethbk.TransactionsQueryingTransactions(client)
}
