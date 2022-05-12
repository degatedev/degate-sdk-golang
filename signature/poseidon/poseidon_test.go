package poseidon

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/degatedev/degatesdk/signature/utils"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/blake2b"
)

func TestBlake2bVersion(t *testing.T) {
	h := blake2b.Sum256([]byte("poseidon_constants"))
	assert.Equal(t,
		"e57ba154fb2c47811dc1a2369b27e25a44915b4e4ece4eb8ec74850cb78e01b1",
		hex.EncodeToString(h[:]))
}

func TestPoseidonHash(t *testing.T) {
	b0 := big.NewInt(0)
	b1 := big.NewInt(1)
	b2 := big.NewInt(2)

	h, err := Hash([]*big.Int{b1})
	assert.Nil(t, err)
	assert.Equal(t,
		"11316722965829087614032985243432266723826890185209218714357779037968059437034",
		h.String())

	h, err = Hash([]*big.Int{b1, b2})
	assert.Nil(t, err)
	assert.Equal(t,
		"18034868597434240293665220970421168445584131937984445797953356852217236273181",
		h.String())

	h, err = Hash([]*big.Int{b0, b1, b2, b0, b1, b2, b0, b1, b2})
	assert.Nil(t, err)
	assert.Equal(t,
		"901171399164243716242986386548401611619556238528027833916499064730590161382",
		h.String())

	h, err = Hash([]*big.Int{b0, b1, b2, b0, b1, b2})
	assert.Nil(t, err)
	assert.Equal(t,
		"13102070988478037395154308865607405548746274688434317574093002894058697028363",
		h.String())
}

func TestBigIntPoseidonHash(t *testing.T) {
	b0, ok := big.NewInt(0).SetString("69588426711107115100232500042334179657931174539151555867956034570704220523596", 10)
	assert.True(t, ok)

	h, err := Hash([]*big.Int{b0})
	assert.Nil(t, err)
	assert.Equal(t,
		"17301542653460600976115435789627461515455895446166776549412913422670972634442",
		h.String())
}

func TestErrorInputs(t *testing.T) {
	b0 := big.NewInt(0)
	b1 := big.NewInt(1)
	b2 := big.NewInt(2)

	_, err := Hash([]*big.Int{b1, b2, b0, b0, b0, b0})
	assert.Nil(t, err)

	_, err = Hash([]*big.Int{b1, b2, b0, b0, b0, b0, b0})
	assert.NotNil(t, err)
	assert.Equal(t, "invalid inputs length 7, max 7", err.Error())

	_, err = Hash([]*big.Int{b1, b2, b0, b0, b0, b0, b0, b0})
	assert.NotNil(t, err)
	assert.Equal(t, "invalid inputs length 8, max 7", err.Error())
}

func BenchmarkPoseidonHash(b *testing.B) {
	b0 := big.NewInt(0)
	b1 := utils.NewIntFromString("12242166908188651009877250812424843524687801523336557272219921456462821518061") //nolint:lll
	b2 := utils.NewIntFromString("12242166908188651009877250812424843524687801523336557272219921456462821518061") //nolint:lll

	bigArray4 := []*big.Int{b1, b2, b0, b0, b0, b0}

	for i := 0; i < b.N; i++ {
		Hash(bigArray4) //nolint:errcheck,gosec
	}
}
