package logs

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hokaccha/go-prettyjson"
)

func BoldString(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func BoldGreenString(str string) string {
	return "\x1b[32;1m" + str + "\033[0m"
}

func GrayString(str string) string {
	return "\033[38;5;242m" + str + "\033[0m"
}

func PrintFromatter(v interface{}) {
	formatter := prettyjson.NewFormatter()
	formatter.Indent = 2

	coloredJson, err := formatter.Marshal(v)
	if err != nil {
		log.Fatalf("Failed to Marshal calldata: %v", err)
	}
	fmt.Println(string(coloredJson))
}

func PrintReceiptFromatter(receipt *types.Receipt) {
	formatter := prettyjson.NewFormatter()
	formatter.Indent = 2

	txInfo := map[string]interface{}{
		"hash":   receipt.TxHash,
		"status": receipt.Status,
	}

	coloredJson, err := formatter.Marshal(txInfo)
	if err != nil {
		log.Fatalf("Failed to Marshal calldata: %v", err)
	}
	fmt.Println(string(coloredJson))
}

func PrintBlockScope(from int, limit int) {
	fmt.Printf("\rSearch event from %d block number | total searched "+BoldString("%d blocks")+" ...🚀", from, limit)
}
