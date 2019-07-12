package main

import (
	"log"

	"github.com/dlm/go-eth-book/gethbk"
)

const CONFIG_LOCAL string = "CONFIG_LOCAL"
const CONFIG_MAINNET string = "CONFIG_MAINNET"
const CONFIG string = CONFIG_MAINNET

func main() {
	log.Println("Running with configuration ", CONFIG)

	url := "https://mainnet.infura.io/v3/c540c4fe10614f27a05a22bd463f0996"
	if CONFIG == CONFIG_LOCAL {
		url = "http://localhost:8545"
	}
	_ = url
	client := gethbk.Client(url)

	accountHex := "0x71c7656ec7ab88b098defb751b7401b5f6d8976f"
	if CONFIG == CONFIG_LOCAL {
		accountHex = "0xe280029a7867ba5c9154434886c241775ea87e53"
	}
	_ = accountHex

	// gethbk.AccountsBalances(client, accountHex)
	// gethbk.AccountsGeneratingNewWallets()
	// gethbk.AccountsKeystores()
	gethbk.AddressCheck(client)
}
