package gethbk

import (
	"log"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Client(url string) *ethclient.Client {
	client, err := ethclient.Dial(url)
	checkForError(err)
	log.Println("Connection success")
	return client
}

