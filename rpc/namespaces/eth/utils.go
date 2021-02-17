package eth

import (
	"encoding/json"
	"strings"

	"github.com/cosmos/ethermint/x/evm/types"

	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	DefaultEVMErrorCode          = -32000
	VMExecuteException           = -32015
	VMExecuteExceptionInEstimate = 3

	RPCEthCall           = "eth_call"
	RPCEthEstimateGas    = "eth_estimateGas"
	RPCEthGetBlockByHash = "eth_getBlockByHash"

	RPCUnknowErr = "unknow"
	RPCNullData  = "null"
)

type cosmosError struct {
	Code      int    `json:"code"`
	Log       string `json:"log"`
	Codespace string `json:"codespace"`
}

func (c cosmosError) Error() string {
	return c.Log
}

type wrappedEthError struct {
	Wrap ethDataError `json:"0x00000000000000000000000000000000"`
}

type ethDataError struct {
	Error          string `json:"error"`
	ProgramCounter int    `json:"program_counter"`
	Reason         string `json:"reason"`
	Ret            string `json:"return"`
}

type DataError struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (d DataError) Error() string {
	return d.Msg
}

func (d DataError) ErrorData() interface{} {
	return d.Data
}

func (d DataError) ErrorCode() int {
	return d.Code
}

func newDataError(revert string, data string) *wrappedEthError {
	return &wrappedEthError{
		Wrap: ethDataError{
			Error:          "revert",
			ProgramCounter: 0,
			Reason:         revert,
			Ret:            data,
		}}
}

func TransformDataError(err error, method string) error {
	msg := err.Error()
	var realErr cosmosError
	if len(msg) > 0 {
		e := json.Unmarshal([]byte(msg), &realErr)
		if e != nil {
			return DataError{
				Code: DefaultEVMErrorCode,
				Msg:  err.Error(),
				Data: RPCNullData,
			}
		}
		if method == RPCEthGetBlockByHash {
			return DataError{
				Code: DefaultEVMErrorCode,
				Msg:  realErr.Error(),
				Data: RPCNullData,
			}
		}
		retErr, m := preProcessError(realErr, err.Error())
		if retErr != nil {
			return realErr
		}
		//if there have multi error type of EVM, this need a reactor mode to process error
		revert, f := m[vm.ErrExecutionReverted.Error()]
		if !f {
			revert = RPCUnknowErr
		}
		data, f := m[types.ErrorHexData]
		if !f {
			data = RPCNullData
		}
		switch method {
		case RPCEthEstimateGas:
			return DataError{
				Code: VMExecuteExceptionInEstimate,
				Msg:  revert,
				Data: data,
			}
		case RPCEthCall:
			return DataError{
				Code: VMExecuteException,
				Msg:  revert,
				Data: newDataError(revert, data),
			}
		default:
			return DataError{
				Code: DefaultEVMErrorCode,
				Msg:  revert,
				Data: newDataError(revert, data),
			}
		}

	}
	return DataError{
		Code: DefaultEVMErrorCode,
		Msg:  err.Error(),
		Data: RPCNullData,
	}
}

//Preprocess error string, the string of realErr.Log is most like:
//`["execution reverted","message","HexData","0x00000000000"];some failed information`
//we need marshalled json slice from realErr.Log and using segment tag `[` and `]` to cut it
func preProcessError(realErr cosmosError, origErrorMsg string) (error, map[string]string) {
	var logs []string
	lastSeg := strings.LastIndexAny(realErr.Log, "]")
	if lastSeg < 0 {
		return DataError{
			Code: DefaultEVMErrorCode,
			Msg:  origErrorMsg,
			Data: RPCNullData,
		}, nil
	}
	marshaler := realErr.Log[0 : lastSeg+1]
	e := json.Unmarshal([]byte(marshaler), &logs)
	if e != nil {
		return DataError{
			Code: DefaultEVMErrorCode,
			Msg:  origErrorMsg,
			Data: RPCNullData,
		}, nil
	}
	m := genericStringMap(logs)
	if m == nil {
		return DataError{
			Code: DefaultEVMErrorCode,
			Msg:  origErrorMsg,
			Data: RPCNullData,
		}, nil
	}
	return nil, m
}

func genericStringMap(s []string) map[string]string {
	var ret = make(map[string]string)
	if len(s)%2 != 0 {
		return nil
	}
	for i := 0; i < len(s); i += 2 {
		ret[s[i]] = s[i+1]
	}
	return ret
}
