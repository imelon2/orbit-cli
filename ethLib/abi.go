package ethlib

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type Calldata struct {
	Abi    *abi.ABI
	Method *abi.Method
	Data   []byte
}

func NewCalldata(abi *abi.ABI, calldata []byte) *Calldata {
	newCalldata := new(Calldata)
	newCalldata.Abi = abi
	newCalldata.Data = calldata
	return newCalldata
}

func GetAbi(abiReader *strings.Reader) (abi.ABI, error) {
	parsedABI, err := abi.JSON(abiReader)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to get ABI: %v", err)
	}
	return parsedABI, nil
}

func (c *Calldata) GetMethodById() (*Calldata, error) {
	method, err := c.Abi.MethodById(c.Data[:4])
	if err != nil {
		return nil, fmt.Errorf("failed to get method from calldata: %v", err)
	}
	c.Method = method
	return c, nil
}

func (c *Calldata) GetUnpackedHexdata() ([]interface{}, error) {
	if c.Method == nil {
		c.GetMethodById()
	}

	hex, err := c.Method.Inputs.Unpack(c.Data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack calldata: %v", err)
	}
	return hex, nil
}
