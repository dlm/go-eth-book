package gethbk

import (
	"log"
)

func logFatal(v ...interface{}) {
	log.Fatal(v)
}

func logInfo(v ...interface{}) {
	log.Println(v)
}

func checkForError(err error) {
	if err != nil {
		logFatal(err)
		panic(err)
	}
}

func checkOk(ok bool, msg string) {
	if !ok {
		logFatal(msg)
		panic(msg)
	}
}

