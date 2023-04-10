package lib

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
	"strings"

	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/request"
	"github.com/degatedev/degate-sdk-golang/signature"
	goSha3 "github.com/degatedev/degate-sdk-golang/signature/go-sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethersphere/bee/pkg/crypto/eip712"
)

const (
	message = "This signature creates a temporary Asset Key Pair for DeGate %s transactions, with no fees or blockchain transactions required.\\n\\nKey Nonce: \\n%d"
)

func CreateAppKey(privateKey string, baseUrl string, keyNonce uint) (secretKey string, publicKeyX string, publicKeyY string, err error) {
	var env string
	baseUrl = strings.ToLower(baseUrl)
	if strings.Contains(baseUrl, "dev") {
		env = "Dev"
	} else if strings.Contains(baseUrl, "testa") {
		env = "TestA"
	} else if strings.Contains(baseUrl, "testb") {
		env = "TestB"
	} else if strings.Contains(baseUrl, "testc") {
		env = "TestC"
	} else if strings.Contains(baseUrl, "testd") {
		env = "TestD"
	} else if strings.Contains(baseUrl, "teste") {
		env = "TestE"
	} else if strings.Contains(baseUrl, "testf") {
		env = "TestF"
	} else if strings.Contains(baseUrl, "testg") {
		env = "TestG"
	} else if strings.Contains(baseUrl, "testnet") {
		env = "Testnet"
	} else {
		env = "Mainnet"
	}
	_, sign, err := signature.PrivateKeySign(privateKey, fmt.Sprintf(message, env, keyNonce))
	if err != nil {
		return
	}
	x, y, secretKey := signature.GetSignDataPublicKey(sign)
	if px, is := new(big.Int).SetString(x, 10); !is {
		err = errors.New("error publicX")
		return
	} else {
		publicKeyX = hexutil.Encode(px.Bytes())
	}
	if py, is := new(big.Int).SetString(y, 10); !is {
		err = errors.New("error publicY")
		return
	} else {
		publicKeyY = hexutil.Encode(py.Bytes())
	}
	return
}

func SignUpdateAccount(privateKey string, owner, accountID string, feeTokenID uint64, maxFee, validUntil, nonce, publicKeyX, publicY, exchange string, referID uint64, chainID int64) (s string, err error) {
	if !model.IsETHAddress(exchange) {
		err = errors.New("illegal exchange address")
		return
	}
	var EIP712DomainType = []eip712.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	}

	var AccountUpdate = []eip712.Type{
		{Name: "owner", Type: "address"},
		{Name: "accountID", Type: "uint32"},
		{Name: "feeTokenID", Type: "uint32"},
		{Name: "maxFee", Type: "uint96"},
		{Name: "publicKey", Type: "uint256"},
		{Name: "validUntil", Type: "uint32"},
		{Name: "nonce", Type: "uint32"},
	}

	var types = map[string][]eip712.Type{
		"EIP712Domain":  EIP712DomainType,
		"AccountUpdate": AccountUpdate,
	}

	var messageType = map[string]interface{}{
		"owner":      common.HexToAddress(owner).Hex(),
		"accountID":  accountID,
		"feeTokenID": strconv.FormatUint(feeTokenID, 10),
		"maxFee":     maxFee,
		"validUntil": validUntil,
		"nonce":      nonce,
	}

	var DomainType = eip712.TypedDataDomain{
		Name:              signature.DomainTypeName,
		Version:           signature.DomainTypeVersion,
		ChainId:           math.NewHexOrDecimal256(chainID),
		VerifyingContract: common.HexToAddress(exchange).Hex(),
	}

	var typedData = eip712.TypedData{
		Types:       types,
		PrimaryType: "AccountUpdate",
		Message:     messageType,
		Domain:      DomainType,
	}

	pubKey, err := signature.EIP712PublicKey(publicKeyX, publicY)
	if err != nil {
		return
	}
	typedData.Message["publicKey"] = pubKey

	s, err = signature.PrivateKeyEIP712Sign(privateKey, &typedData, true)
	return
}

func SignHeader(privateKey string, owner string, time int64) (s string, err error) {
	if len(privateKey) == 0 {
		err = errors.New("AssetPrivateKey is empty")
		return
	}
	if !model.IsETHAddress(owner) {
		err = errors.New("illegal address")
		return
	}
	var inpBI []*big.Int
	ownerBig, is := new(big.Int).SetString(owner[2:], 16)
	if !is {
		err = errors.New("error owner")
		return
	}
	inpBI = append(inpBI, ownerBig)
	inpBI = append(inpBI, new(big.Int).SetInt64(time))

	s, err = signature.EddsaSign(privateKey, inpBI, len(inpBI)+1, 6, 53)
	if err != nil {
		return
	}
	return
}

func SignOrder(privateKey, exchange string, sellVol, buyVol, feeVol string, fillAmountBOrs bool, storageId, accountId, validUntil, sellTokenId, buyTokenId, feeTokenId uint64) (s string, err error) {
	if len(privateKey) == 0 {
		err = errors.New("AssetPrivateKey is empty")
		return
	}
	if !model.IsETHAddress(exchange) {
		err = errors.New("illegal exchange address")
		return
	}
	exchangeBig, isBig := new(big.Int).SetString(exchange[2:], 16)
	if !isBig {
		err = errors.New("error exchange")
		return
	}
	storageIDBig := new(big.Int).SetUint64(storageId)
	accountIDBig := new(big.Int).SetUint64(uint64(accountId))
	tokenSIdBig := new(big.Int).SetUint64(uint64(sellTokenId))
	tokenBIdBig := new(big.Int).SetUint64(uint64(buyTokenId))
	amountSInBNBig, isBig := new(big.Int).SetString(sellVol, 10)
	if !isBig {
		err = errors.New("error sell volume")
		return
	}
	amountBInBNBig, isBig := new(big.Int).SetString(buyVol, 10)
	if !isBig {
		err = errors.New("error buy volume")
		return
	}
	validUntilBig := new(big.Int).SetInt64(int64(validUntil))
	fillAmountBOrSBig := new(big.Int)
	if fillAmountBOrs {
		fillAmountBOrSBig.SetInt64(1)
	} else {
		fillAmountBOrSBig.SetInt64(0)
	}
	feeTokenIdBig := new(big.Int).SetUint64(uint64(feeTokenId))
	maxFeeBipsBig, isBig := new(big.Int).SetString(feeVol, 10)
	if !isBig {
		err = errors.New("error fee volume")
		return
	}
	typeBig := new(big.Int).SetUint64(0)
	gridOffsetBig, _ := new(big.Int).SetString("0", 10)
	orderOffsetBig, _ := new(big.Int).SetString("0", 10)
	maxLevelBig := new(big.Int).SetUint64(0)
	userAppKey := new(big.Int).SetUint64(0)

	var inpBI []*big.Int
	inpBI = append(inpBI, exchangeBig)
	inpBI = append(inpBI, storageIDBig)
	inpBI = append(inpBI, accountIDBig)
	inpBI = append(inpBI, tokenSIdBig)
	inpBI = append(inpBI, tokenBIdBig)
	inpBI = append(inpBI, amountSInBNBig)
	inpBI = append(inpBI, amountBInBNBig)
	inpBI = append(inpBI, validUntilBig)
	inpBI = append(inpBI, fillAmountBOrSBig)
	inpBI = append(inpBI, big.NewInt(0))
	inpBI = append(inpBI, feeTokenIdBig)
	inpBI = append(inpBI, maxFeeBipsBig)
	inpBI = append(inpBI, typeBig)
	inpBI = append(inpBI, gridOffsetBig)
	inpBI = append(inpBI, orderOffsetBig)
	inpBI = append(inpBI, maxLevelBig)
	inpBI = append(inpBI, userAppKey)

	s, err = signature.EddsaSign(privateKey, inpBI, len(inpBI)+1, 6, 53)
	if err != nil {
		return
	}
	return
}

func SignOrderRequest(privateKey string, exchange string, r *model.DGOrderParam) (s string, err error) {
	if len(privateKey) == 0 {
		err = errors.New("AssetPrivateKey is empty")
		return
	}
	if !model.IsETHAddress(exchange) {
		err = errors.New("illegal exchange address")
		return
	}
	exchangeBig, isBig := new(big.Int).SetString(exchange[2:], 16)
	if !isBig {
		err = errors.New("error exchange")
		return
	}
	storageIDBig := new(big.Int).SetUint64(r.StorageId)
	accountIDBig := new(big.Int).SetUint64(r.AccountId)
	tokenSIdBig := new(big.Int).SetUint64(r.SellToken.TokenId)
	tokenBIdBig := new(big.Int).SetUint64(r.BuyToken.TokenId)
	amountSInBNBig, isBig := new(big.Int).SetString(r.SellToken.Volume, 10)
	if !isBig {
		err = errors.New("error sell volume")
		return
	}
	amountBInBNBig, isBig := new(big.Int).SetString(r.BuyToken.Volume, 10)
	if !isBig {
		err = errors.New("error buy volume")
		return
	}
	validUntilBig := new(big.Int).SetInt64(r.ValidUntil)
	fillAmountBOrSBig := new(big.Int)
	if r.FillAmountBOrs {
		fillAmountBOrSBig.SetInt64(1)
	} else {
		fillAmountBOrSBig.SetInt64(0)
	}
	feeTokenIdBig := new(big.Int).SetUint64(r.FeeToken.TokenId)
	feeVolumeBig, isBig := new(big.Int).SetString(r.FeeToken.Volume, 10)
	if !isBig {
		err = errors.New("error fee volume")
		return
	}
	feeBipsDec, err := decimal.NewFromString(r.FeeBips)
	if err != nil {
		return
	}
	feeBipsBig := feeBipsDec.BigInt()

	typeBig := new(big.Int).SetUint64(0)
	gridOffsetBig, _ := new(big.Int).SetString("0", 10)
	orderOffsetBig, _ := new(big.Int).SetString("0", 10)
	maxLevelBig := new(big.Int).SetUint64(0)
	userAppKey := new(big.Int).SetUint64(0)

	var inpBI []*big.Int
	inpBI = append(inpBI, exchangeBig)
	inpBI = append(inpBI, storageIDBig)
	inpBI = append(inpBI, accountIDBig)
	inpBI = append(inpBI, tokenSIdBig)
	inpBI = append(inpBI, tokenBIdBig)
	inpBI = append(inpBI, amountSInBNBig)
	inpBI = append(inpBI, amountBInBNBig)
	inpBI = append(inpBI, validUntilBig)
	inpBI = append(inpBI, fillAmountBOrSBig)
	inpBI = append(inpBI, big.NewInt(0))
	inpBI = append(inpBI, feeTokenIdBig)
	inpBI = append(inpBI, feeVolumeBig)
	inpBI = append(inpBI, feeBipsBig)
	inpBI = append(inpBI, typeBig)
	inpBI = append(inpBI, gridOffsetBig)
	inpBI = append(inpBI, orderOffsetBig)
	inpBI = append(inpBI, maxLevelBig)
	inpBI = append(inpBI, userAppKey)

	s, err = signature.EddsaSign(privateKey, inpBI, len(inpBI)+1, 6, 53)
	if err != nil {
		return
	}
	return
}

func SignCancelOrderNew(privateKey, exchange string, accountID, storageID uint64, gasMaxFee string, gasFeeTokenID uint64) (s string, err error) {
	if len(privateKey) == 0 {
		err = errors.New("no AssetPrivateKey")
		return
	}
	var (
		exchangeBig      = new(big.Int)
		accountIDBig     = new(big.Int)
		storageIDBig     = new(big.Int)
		gasMaxFeeBig     = new(big.Int)
		gasFeeTokenIDBig = new(big.Int)
		useAppKeyBig     = new(big.Int).SetUint64(0)
		isBig            bool
	)

	_, isBig = exchangeBig.SetString(exchange[2:], 16)
	if !isBig {
		err = errors.New("error exchange")
		return
	}
	accountIDBig.SetUint64(accountID)
	storageIDBig.SetUint64(storageID)
	_, isBig = gasMaxFeeBig.SetString(gasMaxFee, 10)
	if !isBig {
		err = errors.New("error gas fee")
		return
	}
	gasFeeTokenIDBig.SetUint64(gasFeeTokenID)

	var inpBI []*big.Int
	inpBI = append(inpBI, exchangeBig)
	inpBI = append(inpBI, accountIDBig)
	inpBI = append(inpBI, storageIDBig)
	inpBI = append(inpBI, gasMaxFeeBig)
	inpBI = append(inpBI, gasFeeTokenIDBig)
	inpBI = append(inpBI, useAppKeyBig)

	s, err = signature.EddsaSign(privateKey, inpBI, len(inpBI)+1, 6, 53)
	if err != nil {
		return
	}
	return
}

func SignCancelOrder(privateKey, exchange, owner string, accountID, tokenID, storageID uint64, gasMaxFee string, gasFeeTokenID uint64) (s string, err error) {
	if len(privateKey) == 0 {
		err = errors.New("no AssetPrivateKey")
		return
	}
	var (
		exchangeBig      = new(big.Int)
		ownerBig         = new(big.Int)
		accountIDBig     = new(big.Int)
		tokenIDBig       = new(big.Int)
		storageIDBig     = new(big.Int)
		gasMaxFeeBig     = new(big.Int)
		gasFeeTokenIDBig = new(big.Int)
		isBig            bool
	)

	_, isBig = exchangeBig.SetString(exchange[2:], 16)
	if !isBig {
		err = errors.New("error exchange")
		return
	}
	_, isBig = ownerBig.SetString(owner[2:], 16)
	if !isBig {
		err = errors.New("error owner")
		return
	}
	accountIDBig.SetUint64(accountID)
	tokenIDBig.SetUint64(tokenID)
	storageIDBig.SetUint64(storageID)
	_, isBig = gasMaxFeeBig.SetString(gasMaxFee, 10)
	if !isBig {
		err = errors.New("error gas fee")
		return
	}
	gasFeeTokenIDBig.SetUint64(0) // here just set to 0

	var inpBI []*big.Int
	inpBI = append(inpBI, exchangeBig)
	inpBI = append(inpBI, ownerBig)
	inpBI = append(inpBI, accountIDBig)
	inpBI = append(inpBI, tokenIDBig)
	inpBI = append(inpBI, storageIDBig)
	inpBI = append(inpBI, gasMaxFeeBig)
	inpBI = append(inpBI, gasFeeTokenIDBig)

	s, err = signature.EddsaSign(privateKey, inpBI, len(inpBI)+1, 6, 53)
	if err != nil {
		return
	}
	return
}

func SignTransfer(privateKey string, exchange string, r *request.TransferRequest) (s string, err error) {
	if len(privateKey) == 0 {
		err = errors.New("no AssetPrivateKey")
		return
	}
	if !model.IsETHAddress(exchange) {
		err = errors.New("illegal exchange address")
		return
	}
	var toBig = new(big.Int)
	toBig.SetString(r.To[2:], 16)

	var exchangeBig = new(big.Int)
	exchangeBig.SetString(exchange[2:], 16)

	var inpBI []*big.Int
	inpBI = append(inpBI, exchangeBig)
	inpBI = append(inpBI, new(big.Int).SetUint64(uint64(r.AccountId)))
	inpBI = append(inpBI, new(big.Int).SetUint64(uint64(r.ToAccountId)))
	inpBI = append(inpBI, new(big.Int).SetUint64(uint64(r.Token.TokenId)))
	tokenAmount, is := new(big.Int).SetString(r.Token.Volume, 10)
	if !is {
		err = errors.New("error token volume")
		return
	}
	inpBI = append(inpBI, tokenAmount)
	inpBI = append(inpBI, new(big.Int).SetUint64(uint64(r.MaxFee.TokenId)))
	feeAmount, is := new(big.Int).SetString(r.MaxFee.Volume, 10)
	if !is {
		err = errors.New("error fee volume")
		return
	}
	userAppKey := new(big.Int).SetUint64(0)

	inpBI = append(inpBI, feeAmount)
	inpBI = append(inpBI, toBig)
	inpBI = append(inpBI, new(big.Int))
	inpBI = append(inpBI, new(big.Int))
	inpBI = append(inpBI, new(big.Int).SetInt64(r.ValidUntil))
	inpBI = append(inpBI, new(big.Int).SetUint64(r.StorageId))
	inpBI = append(inpBI, userAppKey)

	s, err = signature.EddsaSign(privateKey, inpBI, len(inpBI)+1, 6, 53)
	if err != nil {
		return
	}
	return
}

func SignTransferEcdsa(privateKey, owner, accountID, tokenID, amount, feeTokenID, maxFee, validUntil, storageID, exchange, to string, chainID int64) (s string, err error) {
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

	var typedData = eip712.TypedData{
		Types:       types,
		PrimaryType: "Transfer",
		Message:     messageType,
		Domain:      DomainType,
	}
	s, err = signature.PrivateKeyEIP712Sign(privateKey, &typedData, true)
	return
}

func SignWithdrawEddsa(privateKey string, exchange string, r *request.WithdrawRequest) (s string, err error) {
	if len(privateKey) == 0 {
		err = errors.New("no AssetPrivateKey")
		return
	}
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

	s, err = signature.EddsaSign(privateKey, inpBI, len(inpBI)+1, 6, 53)
	if err != nil {
		return
	}
	return
}

func SignWithdrawEcdsa(privateKey string, owner, accountID, tokenID, amount, feeTokenID, maxFee, minGas, validUntil, storageID, exchange, to string, chainID int64) (s string, err error) {
	if len(privateKey) == 0 {
		err = errors.New("no PrivateKey")
		return
	}
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

	var typedData = eip712.TypedData{
		Types:       types,
		PrimaryType: "Withdrawal",
		Message:     messageType,
		Domain:      DomainType,
	}

	s, err = signature.PrivateKeyEIP712Sign(privateKey, &typedData, true)
	return
}
