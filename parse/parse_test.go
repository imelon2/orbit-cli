package parse_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hokaccha/go-prettyjson"
	"github.com/imelon2/orbit-cli/common/tx"
	"github.com/imelon2/orbit-cli/parse"
)

var (
	CLIENT_URL string
	HASH       string
	ERROR_MSG  string
)

func init() {
	CLIENT_URL = "http://localhost:8547"
	HASH = "0x68c5c3bebd17776a3b39c9cc3194257062286ace2fdb0701af717400d6aa0c01"
	ERROR_MSG = "0xfadf238a000000000000000000000000000000000000000000000000000000000000004d0000000000000000000000000000000000000000000000000000000000000063"

}

func Test_ParseErrorByBytes(t *testing.T) {
	parse, err := parse.NewParse()
	if err != nil {
		log.Fatal(err)
	}

	errorJson, err := parse.ParseErrorByBytes(ERROR_MSG)
	if err != nil {
		log.Fatal(err)
	}

	formatter := prettyjson.NewFormatter()
	formatter.Indent = 2

	coloredJson, err := formatter.Marshal(errorJson)
	if err != nil {
		log.Fatalf("Failed to Marshal calldata: %v", err)
	}
	fmt.Println(string(coloredJson))
}

func Test_ParseError(t *testing.T) {
	client, err := ethclient.Dial(CLIENT_URL)
	if err != nil {
		log.Fatal(err)
	}
	hash := common.HexToHash(HASH)
	txLib := tx.NewTxLib(client, &hash)

	_, errorData, status, err := txLib.GetTransactionReturn()
	if err != nil {
		log.Fatal(err)
	}

	if *status {
		fmt.Printf("transaction hash %s is SUCCESS", hash.Hex())
		return
	}

	parse, err := parse.NewParse()
	if err != nil {
		log.Fatal(err)
	}

	errorJson, err := parse.ParseError(*errorData)
	if err != nil {
		log.Fatal(err)
	}

	formatter := prettyjson.NewFormatter()
	formatter.Indent = 2

	coloredJson, err := formatter.Marshal(errorJson)
	if err != nil {
		log.Fatalf("Failed to Marshal calldata: %v", err)
	}

	fmt.Println(string(coloredJson))
}
