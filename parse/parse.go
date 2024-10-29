package parse

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/imelon2/orbit-cli/common/path"
	"github.com/imelon2/orbit-cli/common/utils"
)

// Revert reason ID for "Error(string)"
var revertReasonID = "08c379a0"

type Parse struct {
	Abi *abi.ABI
}

type CalldataLog struct {
	Function string      `json:"function"`
	Params   interface{} `json:"params"`
}

type EventLog struct {
	Name      string      `json:"name"`
	Signature string      `json:"signature"`
	Topic     string      `json:"topic"`
	Params    interface{} `json:"params"`
}

type ErrorDataLog struct {
	System  string        `json:"system,omitempty"`
	Message string        `json:"message,omitempty"` // 빈 값 허용
	Custom  CustomMessage `json:"custom,omitempty"`  // nil 값 허용
	Hex     string        `json:"hex"`
}

type CustomMessage struct {
	Name   string      `json:"name,omitempty"`
	Params interface{} `json:"params,omitempty"`
}

func NewParse() (*Parse, error) {
	abiPath := path.GetAbiPath()
	_abi, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, fmt.Errorf("failed read abi file: %v", err)
	}
	parsedABI, err := abi.JSON(strings.NewReader(string(_abi)))
	if err != nil {
		return nil, fmt.Errorf("failed get abi json: %v", err)
	}
	return &Parse{
		Abi: &parsedABI,
	}, nil
}

func (parse *Parse) ParseCalldata(calldata []byte) (*CalldataLog, error) {
	method, err := parse.Abi.MethodById(calldata[:4]) // function selector
	if err != nil {
		return nil, fmt.Errorf("failed to get method from calldata: %v", err)
	}
	hex, err := method.Inputs.Unpack(calldata[4:]) // data
	if err != nil {
		return nil, fmt.Errorf("failed to unpack calldata: %v", err)
	}

	jsonCalldata := make(map[string]interface{})
	for i, data := range hex {
		jsonCalldata[method.Inputs[i].Name] = utils.ConvertBytesToHex(data)
	}

	return &CalldataLog{
		Function: method.RawName,
		Params:   jsonCalldata,
	}, nil
}

func (parse *Parse) ParseEvent(logs []*types.Log) (*[]EventLog, error) {
	result := make([]EventLog, 0)
	for _, log := range logs {
		signature := log.Topics[0]
		event, err := parse.Abi.EventByID(signature)
		if err != nil {
			return nil, fmt.Errorf("failed get event id: %v", err)
		}
		log.Topics = log.Topics[1:] // remove event signature

		logDatas, err := event.Inputs.Unpack(log.Data)
		if err != nil {
			return nil, fmt.Errorf("failed decode unpack event: %v", err)
		}

		// parse event data
		eventData := make(map[string]interface{})
		dataIndex := 0
		topicIndex := 0
		for _, eventInput := range event.Inputs {
			if eventInput.Indexed {
				topic := log.Topics[topicIndex]
				data, err := utils.ConvertIndexedInput(eventInput, topic.Bytes())
				if err != nil {
					return nil, err
				}
				eventData[eventInput.Name] = data
				topicIndex++
			} else {
				eventData[eventInput.Name] = logDatas[dataIndex]
				dataIndex++
			}
		}
		result = append(result, EventLog{
			Name:      event.RawName,
			Signature: event.Sig,
			Topic:     signature.Hex(),
			Params:    utils.ConvertBytesToHex(eventData),
		})
	}

	return &result, nil
}

func (parse *Parse) ParseError(rpcErr rpc.DataError) (*ErrorDataLog, error) {
	errRuslt := ErrorDataLog{}

	system := rpcErr.Error() // system error
	sError, ok := rpcErr.ErrorData().(string)
	if !ok {
		return nil, fmt.Errorf("failed convert error data to string")
	}
	sError = utils.Unhexlify(sError)

	if system == "" /* Unexpected Error */ {
		errRuslt.Hex = "0x"
		errRuslt.Message = "NULL"
	} else if sError[:8] == revertReasonID /* String Error */ {
		bError, err := hex.DecodeString(sError[8:]) // remove error signature
		if err != nil {
			return nil, fmt.Errorf("failed hex decode string: %v", err)
		}
		stringType, err := abi.NewType("string", "", nil)
		if err != nil {
			return nil, fmt.Errorf("failed new type on abi: %v", err)
		}
		args := abi.Arguments{
			{Type: stringType},
		}
		decodeError, err := args.Unpack(bError)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack error: %v", err)
		}
		for _, s := range decodeError {
			_s, _ := s.(string)
			errRuslt.Message = _s
			break
		}
		errRuslt.System = system
		errRuslt.Hex = sError
	} else /* Custom Error */ {
		bError, err := hex.DecodeString(sError)
		if err != nil {
			return nil, fmt.Errorf("failed hex decode string: %v", err)
		}

		/* generate signature */
		var sigdata [4]byte
		copy(sigdata[:], bError[:4])

		errorAbi, err := parse.Abi.ErrorByID(sigdata)
		if err != nil {
			return nil, fmt.Errorf("failed get error id: %v", err)
		}

		decodedError, err := errorAbi.Inputs.Unpack(bError[4:])
		if err != nil {
			return nil, fmt.Errorf("failed unpack error: %v", err)
		}
		jsonError := make(map[string]interface{})
		for i, data := range decodedError {
			jsonError[errorAbi.Inputs[i].Name] = utils.ConvertBytesToHex(data)
		}

		errRuslt.Custom = CustomMessage{
			Name:   errorAbi.Sig,
			Params: jsonError,
		}
		errRuslt.Hex = sError
	}

	return &errRuslt, nil
}

func (parse *Parse) ParseErrorByBytes(err string) (*ErrorDataLog, error) {
	sError := utils.Unhexlify(err)

	errRuslt := ErrorDataLog{}
	if sError[:8] == revertReasonID /* String Error */ {
		bError, err := hex.DecodeString(sError[8:]) // remove error signature
		if err != nil {
			return nil, fmt.Errorf("failed hex decode string: %v", err)
		}
		stringType, err := abi.NewType("string", "", nil)
		if err != nil {
			return nil, fmt.Errorf("failed new type on abi: %v", err)
		}
		args := abi.Arguments{
			{Type: stringType},
		}
		decodeError, err := args.Unpack(bError)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack error: %v", err)
		}
		for _, s := range decodeError {
			_s, _ := s.(string)
			errRuslt.Message = _s
			break
		}
	} else /* Custom Error */ {
		bError, err := hex.DecodeString(sError)
		if err != nil {
			return nil, fmt.Errorf("failed hex decode string: %v", err)
		}

		/* generate signature */
		var sigdata [4]byte
		copy(sigdata[:], bError[:4])

		errorAbi, err := parse.Abi.ErrorByID(sigdata)
		if err != nil {
			return nil, fmt.Errorf("failed get error id: %v", err)
		}

		decodedError, err := errorAbi.Inputs.Unpack(bError[4:])
		if err != nil {
			return nil, fmt.Errorf("failed unpack error: %v", err)
		}
		jsonError := make(map[string]interface{})
		for i, data := range decodedError {
			jsonError[errorAbi.Inputs[i].Name] = utils.ConvertBytesToHex(data)
		}

		errRuslt.Custom = CustomMessage{
			Name:   errorAbi.Sig,
			Params: jsonError,
		}
	}

	errRuslt.Hex = "0x" + sError
	return &errRuslt, nil
}
