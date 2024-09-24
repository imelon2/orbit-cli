package utils

import (
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

var Callopts = &bind.CallOpts{
	Pending: false, // 트랜잭션이 확정된 상태를 조회
	Context: nil,   // 컨텍스트가 필요한 경우 (예: 시간 초과)
}

func GetClient(url string) *ethclient.Client {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
