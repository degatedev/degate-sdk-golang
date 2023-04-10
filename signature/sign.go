package signature

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/degatedev/degate-sdk-golang/signature/babyjub"
	"github.com/degatedev/degate-sdk-golang/signature/poseidon"
	"github.com/degatedev/degate-sdk-golang/signature/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethersphere/bee/pkg/crypto"
	"github.com/ethersphere/bee/pkg/crypto/eip712"
)

const (
	DomainTypeName    = "DeGate Protocol"
	DomainTypeVersion = "0.1.0"
)

func PrivateKeySign(privateKey, message string) (sig string, signature []byte, err error) {
	priKey, err := crypto2.HexToECDSA(privateKey)
	if err != nil {
		return
	}
	msg := fmt.Sprintf("\u0019Ethereum Signed Message:\n%d%s", len(message), message)
	hash3 := crypto2.Keccak256Hash([]byte(msg)).Bytes()

	signature, err = crypto2.Sign(hash3, priKey)
	if err != nil {
		return
	}

	signature[64] += 27
	sig = hexutil.Encode(signature)
	return
}

func EIP712PublicKey(publicKeyX, publicKeyY string) (publicKey string, err error) {
	bigX := new(big.Int)
	bigY := new(big.Int)
	if publicKeyX[0:2] == "0x" {
		publicKeyX = publicKeyX[2:]
		bigX.SetString(publicKeyX, 16)
	} else {
		bigX.SetString(publicKeyX, 10)
	}
	if publicKeyY[0:2] == "0x" {
		publicKeyY = publicKeyY[2:]
		bigY.SetString(publicKeyY, 16)
	} else {
		bigY.SetString(publicKeyY, 10)
	}
	x := utils.NewIntFromString(bigX.String())
	y := utils.NewIntFromString(bigY.String())
	p := &babyjub.Point{X: x, Y: y}
	pComp := p.Compress()
	is, y := babyjub.UnpackSignY(pComp)
	pComp2 := babyjub.PackSignY(is, y)
	emptyPointComp := [32]byte{}
	for i := 0; i < 32; i++ {
		emptyPointComp[31-i] = pComp2[i]
	}
	publicKey = hex.EncodeToString(emptyPointComp[:])
	pk := new(big.Int)
	pk.SetString(publicKey, 16)
	publicKey = pk.String()
	return
}

func GetSignDataPublicKey(signature []byte) (publicKeyX, publicKeyY, secretKey string) {
	var (
		pk *babyjub.Point
	)

	hash := sha256.Sum256(signature)

	res := new(big.Int)
	for i, b := range hash {
		n := new(big.Int)
		n.SetBytes([]byte{b})
		res.Add(res, n.Lsh(n, uint(i*8)))
	}

	p := babyjub.NewPrivKeyScalar(res).Public()

	res.Mod(res, babyjub.SubOrder)
	pk = babyjub.NewPoint().Mul(res, babyjub.B8)
	pk.Projective()
	return p.X.String(), p.Y.String(), res.String()
}

func PrivateKeyEIP712Sign(privateKey string, typeData *eip712.TypedData, isEndJudgment bool) (signature string, err error) {
	var (
		expected []byte
		priKey   *ecdsa.PrivateKey
		sig      []byte
	)

	if expected, err = hex.DecodeString(privateKey); err != nil {
		return
	}
	if priKey, err = crypto.DecodeSecp256k1PrivateKey(expected); err != nil {
		return
	}

	signer1 := NewDefaultSigner(priKey)
	if sig, err = signer1.SignTypedData(typeData); err != nil {
		return
	}

	if isEndJudgment {
		signature = hexutil.Encode(sig) + "02"
	} else {
		signature = hexutil.Encode(sig)
	}
	return
}

func EddsaSign(pri string, msg []*big.Int, t, nRoundsF, nRoundsP int) (signature string, err error) {
	hash, err := poseidon.EddsaHash(msg, t, nRoundsF, nRoundsP, true)
	if err != nil {
		// err = &dgerror.Error{Code: ecode.OrderBookEddsaSignatureError, Message: "error eddsa signature"}
		return
	}

	return EddsaSignHash(pri, hash)
}

func EddsaSignHash(pri string, hash *big.Int) (signature string, err error) {
	var (
		n              string
		k              babyjub.PrivateKey
		signatureRxHex string
		signatureRyHex string
		signatureSHex  string
	)
	if pri[0:2] == "0x" {
		n = pri[2:]
	} else {
		n = pri
	}

	k, err = LeInt2Buff32(n)
	if err != nil {
		// err = &dgerror.Error{Code: ecode.OrderBookEddsaSignatureError, Message: "error eddsa signature"}
		return
	}

	sig := k.SignPoseidon(hash)

	signatureRxHex = common.BigToHash(sig.R8.X).String()
	if len(signatureRxHex) != 66 {
		// err = &dgerror.Error{Code: ecode.OrderBookEddsaSignatureError, Message: "error eddsa signature"}
		return
	}
	signatureRxHex = signatureRxHex[2:]

	signatureRyHex = common.BigToHash(sig.R8.Y).String()
	if len(signatureRyHex) != 66 {
		// err = &dgerror.Error{Code: ecode.OrderBookEddsaSignatureError, Message: "error eddsa signature"}
		return
	}
	signatureRyHex = signatureRyHex[2:]

	signatureSHex = common.BigToHash(sig.S).String()
	if len(signatureSHex) != 66 {
		err = fmt.Errorf("error eddsa signature")
		return
	}
	signatureSHex = signatureSHex[2:]

	signature = "0x" + signatureRxHex + signatureRyHex + signatureSHex
	return
}

func EddsaSignVerify(signature string, publicKeyX, publicKeyY string, msg []*big.Int, t, nRoundsF, nRoundsP int) (is bool, err error) {
	hash, err := poseidon.EddsaHash(msg, t, nRoundsF, nRoundsP, true)
	if err != nil {
		return
	}
	return EddsaSignVerifyHash(signature, publicKeyX, publicKeyY, hash)
}

func EddsaSignVerifyHash(signature string, publicKeyX, publicKeyY string, hash *big.Int) (is bool, err error) {
	var (
		s    string
		sig  *babyjub.Signature
		sigX = new(big.Int)
		sigY = new(big.Int)
		sigS = new(big.Int)
		pk   *babyjub.PublicKey
		pkx  = new(big.Int)
		pky  = new(big.Int)
	)
	if len(signature) == 0 {
		err = fmt.Errorf("error eddsa signature null")
		return
	}
	if signature[0:2] == "0x" {
		s = signature[2:]
	} else {
		s = signature
	}
	if len(s) != 64*3 {
		err = fmt.Errorf("error eddsa signature")
		return
	}

	sigX.SetString(s[0:64], 16)
	sigY.SetString(s[64:128], 16)
	sigS.SetString(s[128:], 16)

	sig = &babyjub.Signature{
		R8: &babyjub.Point{
			X: sigX,
			Y: sigY,
		},
		S: sigS,
	}

	if len(publicKeyX) == 0 {
		err = fmt.Errorf("error eddsa publicKeyX null")
		return
	}
	if len(publicKeyY) == 0 {
		err = fmt.Errorf("error eddsa publicKeyY null")
		return
	}
	if publicKeyX[0:2] == "0x" {
		publicKeyX = publicKeyX[2:]
		pkx.SetString(publicKeyX, 16)
	} else {
		pkx.SetString(publicKeyX, 10)
	}
	if publicKeyY[0:2] == "0x" {
		publicKeyY = publicKeyY[2:]
		pky.SetString(publicKeyY, 16)
	} else {
		pky.SetString(publicKeyY, 10)
	}
	pk = &babyjub.PublicKey{
		X: pkx,
		Y: pky,
	}

	is = pk.VerifyPoseidon(hash, sig)
	return
}

func LeInt2Buff32(n string) (buff [32]byte, err error) {
	var (
		r = new(big.Int)
		o = 0
	)
	r.SetString(n, 10)
	buff = [32]byte{}
	for r.Cmp(big.NewInt(0)) > 0 && o < len(buff) {
		ss := big.NewInt(r.Int64())
		c := ss.And(ss, big.NewInt(255))
		buff[o] = uint8(c.Int64())
		o++
		r = r.Rsh(r, 8)
	}
	return
}
