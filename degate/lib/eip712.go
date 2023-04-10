package lib

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/request"
	"github.com/degatedev/degate-sdk-golang/signature"
	goSha3 "github.com/degatedev/degate-sdk-golang/signature/go-sha3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethersphere/bee/pkg/crypto/eip712"
)

func encodeAndhash(typedData *eip712.TypedData) ([]byte, error) {
	rawData, err := signature.EncodeForSigning(typedData)
	if err != nil {
		return nil, err
	}

	sighash, err := signature.LegacyKeccak256(rawData)
	if err != nil {
		return nil, err
	}

	return sighash, nil
}

// GetEIP712Hash encode EIP712 typed data hash for sign
func GetWithdrawEIP712Hash(owner, accountID, tokenID, amount, feeTokenID, maxFee, minGas, validUntil, storageID, exchange, to string,
	chainID int64) ([]byte, error) {
	var EIP712DomainType = []eip712.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	}

	var AccountUpdate = []eip712.Type{
		{Name: "owner", Type: "address"},
		{Name: "accountID", Type: "uint32"},
		{Name: "tokenID", Type: "uint32"},
		{Name: "amount", Type: "uint248"},
		{Name: "feeTokenID", Type: "uint32"},
		{Name: "maxFee", Type: "uint96"},
		{Name: "to", Type: "address"},
		{Name: "minGas", Type: "uint248"},
		{Name: "validUntil", Type: "uint32"},
		{Name: "storageID", Type: "uint32"},
	}

	var types = map[string][]eip712.Type{
		"EIP712Domain": EIP712DomainType,
		"Withdrawal":   AccountUpdate,
	}

	var messageType = map[string]interface{}{
		"owner":      common.HexToAddress(owner).Hex(),
		"accountID":  accountID,
		"tokenID":    tokenID,
		"amount":     amount,
		"feeTokenID": feeTokenID,
		"maxFee":     maxFee,
		"to":         common.HexToAddress(to).Hex(),
		"minGas":     minGas,
		"validUntil": validUntil,
		"storageID":  storageID,
	}

	var DomainType = eip712.TypedDataDomain{
		Name:              signature.DomainTypeName,
		Version:           signature.DomainTypeVersion,
		ChainId:           math.NewHexOrDecimal256(chainID),
		VerifyingContract: common.HexToAddress(exchange).Hex(),
	}

	return encodeAndhash(&eip712.TypedData{
		Types:       types,
		PrimaryType: "Withdrawal",
		Message:     messageType,
		Domain:      DomainType,
	})
}

// GetTransferEIP712Hash get transfer EIP712 typed data hash for sign
func GetTransferEIP712Hash(owner, accountID, tokenID, amount, feeTokenID, maxFee, validUntil, storageID, exchange, to string,
	chainID int64) ([]byte, error) {
	var EIP712DomainType = []eip712.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	}

	var AccountUpdate = []eip712.Type{
		{Name: "owner", Type: "address"},
		{Name: "accountID", Type: "uint32"},
		{Name: "tokenID", Type: "uint32"},
		{Name: "amount", Type: "uint248"},
		{Name: "feeTokenID", Type: "uint32"},
		{Name: "maxFee", Type: "uint96"},
		{Name: "to", Type: "address"},
		{Name: "validUntil", Type: "uint32"},
		{Name: "storageID", Type: "uint32"},
	}

	var types = map[string][]eip712.Type{
		"EIP712Domain": EIP712DomainType,
		"Transfer":     AccountUpdate,
	}

	var messageType = map[string]interface{}{
		"owner":      common.HexToAddress(owner).Hex(),
		"accountID":  accountID,
		"tokenID":    tokenID,
		"amount":     amount,
		"feeTokenID": feeTokenID,
		"maxFee":     maxFee,
		"to":         common.HexToAddress(to).Hex(),
		"validUntil": validUntil,
		"storageID":  storageID,
	}

	var DomainType = eip712.TypedDataDomain{
		Name:              signature.DomainTypeName,
		Version:           signature.DomainTypeVersion,
		ChainId:           math.NewHexOrDecimal256(chainID),
		VerifyingContract: common.HexToAddress(exchange).Hex(),
	}

	return encodeAndhash(&eip712.TypedData{
		Types:       types,
		PrimaryType: "Transfer",
		Message:     messageType,
		Domain:      DomainType,
	})
}

func EncodeWithdrawRequest(exchange string, r *request.WithdrawRequest) (ss []string, err error) {
	if !model.IsETHAddress(exchange) {
		err = errors.New("illegal exchange address")
		return
	}
	var (
		is bool
	)
	var param []interface{}
	minGas := new(big.Int)
	minGas, is = minGas.SetString(r.MinGas, 10)
	if !is {
		err = errors.New("error min gas")
		return
	}
	amount, is := new(big.Int).SetString(r.Token.Volume, 10)
	if !is {
		err = errors.New("error quantity")
		return
	}
	param = append(param, minGas)
	param = append(param, r.To)
	param = append(param, amount)
	var onChainDataHashBig = new(big.Int)
	onChainDataHash := goSha3.GoSHA3(
		[]string{"uint248", "address", "uint248"},
		param,
	)
	if len(onChainDataHash) > 20 {
		onChainDataHash = onChainDataHash[0:20]
	}
	onChainDataHashString := hex.EncodeToString(onChainDataHash)
	onChainDataHashBig.SetString(onChainDataHashString, 16)

	var exchangeBig = new(big.Int)
	exchangeBig.SetString(exchange[2:], 16)

	maxFee, is := new(big.Int).SetString(r.MaxFee.Volume, 10)
	if !is {
		err = errors.New("error fee quantity")
		return
	}
	userAppKey := new(big.Int).SetUint64(0)

	var inpBI []*big.Int
	inpBI = append(inpBI, exchangeBig)
	inpBI = append(inpBI, new(big.Int).SetUint64(uint64(r.AccountId)))
	inpBI = append(inpBI, new(big.Int).SetUint64(uint64(r.Token.TokenId)))
	inpBI = append(inpBI, amount)
	inpBI = append(inpBI, new(big.Int).SetUint64(uint64(r.MaxFee.TokenId)))
	inpBI = append(inpBI, maxFee)
	inpBI = append(inpBI, onChainDataHashBig)
	inpBI = append(inpBI, new(big.Int).SetInt64(r.ValidUntil))
	inpBI = append(inpBI, new(big.Int).SetUint64(r.StorageId))
	inpBI = append(inpBI, userAppKey)

	for _, item := range inpBI {
		ss = append(ss, item.String())
	}
	return
}
