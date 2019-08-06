package gethbk

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func MakeClient(
	url string,
	dial func(string) (*ethclient.Client, error),
) *ethclient.Client {
	client, err := ethclient.Dial(url)
	checkForError(err)
	log.Println("Connection success")
	return client
}

func Client(url string) *ethclient.Client {
	return MakeClient(url, ethclient.Dial)
}
