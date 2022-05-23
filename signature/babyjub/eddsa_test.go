package babyjub

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/degatedev/degate-sdk-golang/signature/constants"
	"github.com/degatedev/degate-sdk-golang/signature/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicKey(t *testing.T) {
	var k PrivateKey
	for i := 0; i < 32; i++ {
		k[i] = byte(i)
	}
	pk := k.Public()
	assert.True(t, pk.X.Cmp(constants.Q) == -1)
	assert.True(t, pk.Y.Cmp(constants.Q) == -1)
}

func TestSignVerifyMimc7(t *testing.T) {
	var k PrivateKey
	_, err := hex.Decode(k[:],
		[]byte("529ace6f01e6a8c67f3450a045230788f04e66dc2e5343ec8b4616e049a5d1"))
	if err != nil {
		panic(err)
	}
	msgBuf, err := hex.DecodeString("19546049c35a0fe619b77c0f76b89d95c9769845c06374edfe126ad48e8f9603")
	if err != nil {
		panic(err)
	}
	msg := utils.SetBigIntFromLEBytes(new(big.Int), msgBuf)

	pk := k.Public()

	sig := k.SignMimc7(msg)

	ok := pk.VerifyMimc7(msg, sig)

	sigBuf := sig.Compress()
	sig2, err := new(Signature).Decompress(sigBuf)

	fmt.Println("sig:", sig2.S)

	ok = pk.VerifyMimc7(msg, sig2)

	fmt.Println(ok)

	sx := new(big.Int)
	sx.SetString("1485da6fd05516111b6cf97099c4cd0cb1dd5a1bbd6be62bed4795d6001a9698", 16)
	sy := new(big.Int)
	sy.SetString("031ed326bc1e658f36854ac0710e93b9622add4ca73fbb286f48e83c98036bff", 16)
	s1 := new(big.Int)
	s1.SetString("051e0a6efc2cb152cc9a1e4ada7e0d73ae96ffe23c39067e2789bf19f29299b0", 16)
	sig3 := Signature{
		R8: &Point{X: sx, Y: sy},
		S:  s1,
	}
	px := new(big.Int)
	px.SetString("2535b90454fd5772ab81501364db841a4f81fce43f5acc643dce14660e984995", 16)
	py := new(big.Int)
	py.SetString("11b2b8253d6863d84be7fed761e8e1522984765862553f03d653734807e57342", 16)
	pk3 := PublicKey{X: px, Y: py}
	ok = pk3.VerifyMimc7(msg, &sig3)
	fmt.Println(ok)
}

func TestSignVerifyPoseidon(t *testing.T) {
	var k PrivateKey
	_, err := hex.Decode(k[:],
		[]byte("0001020304050607080900010203040506070809000102030405060708090001"))
	require.Nil(t, err)
	msgBuf, err := hex.DecodeString("00010203040506070809")
	if err != nil {
		panic(err)
	}
	msg := utils.SetBigIntFromLEBytes(new(big.Int), msgBuf)

	pk := k.Public()
	assert.Equal(t,
		"15872208232780880391323496162615626329490592476459343692724793783715106083082",
		pk.X.String())
	assert.Equal(t,
		"3297629380257478865105287016917085619944486593062417198110858086548618481395",
		pk.Y.String())

	sig := k.SignPoseidon(msg)
	assert.Equal(t,
		"19739533670544032605675820282046888963403853040461810014389227725175154446510",
		sig.R8.X.String())
	assert.Equal(t,
		"20733987715880802403713959152718145523293348408915441770046155699170249518419",
		sig.R8.Y.String())
	assert.Equal(t,
		"318044161283975110217224102989208461797606283138954992849019458277313526528",
		sig.S.String())

	ok := pk.VerifyPoseidon(msg, sig)
	assert.Equal(t, true, ok)

	sigBuf := sig.Compress()
	sig2, err := new(Signature).Decompress(sigBuf)
	assert.Equal(t, nil, err)

	// assert.Equal(t, ""+
	// 	"dfedb4315d3f2eb4de2d3c510d7a987dcab67089c8ace06308827bf5bcbe02a2"+
	// 	"9d043ece562a8f82bfc0adb640c0107a7d3a27c1c7c1a6179a0da73de5c1b203",
	// 	hex.EncodeToString(sigBuf[:]))

	ok = pk.VerifyPoseidon(msg, sig2)
	assert.Equal(t, true, ok)
}

func TestVerifyPoseidon(t *testing.T) {
	msg := utils.NewIntFromString("18907120458743615336946847248227397370763473802204269898187195559525130063203")
	key := utils.NewIntFromString("56869496543825")

	var k PrivateKey
	k = utils.BigIntLEBytes(key)
	pk := k.Public()
	assert.Equal(t,
		"9255092729144892245186624611131828247442112563544941131408300200214096116351",
		pk.X.String())
	assert.Equal(t,
		"8460370541846376796657659750509399834188652251932899797602116208684247832083",
		pk.Y.String())

	x, ok := big.NewInt(0).SetString("12752937249904285198676276090843566060401682639184875784873451302664399892304", 10)
	require.True(t, ok)
	y, ok := big.NewInt(0).SetString("13530361082613950739674235863189737173485045373827356210876301607961589355327", 10)
	require.True(t, ok)
	s, ok := big.NewInt(0).SetString("7616254846080660730932216519770737127037155777726245055053503272117180880572", 10)
	require.True(t, ok)

	sig := &Signature{
		&Point{
			X: x,
			Y: y,
		},
		s,
	}

	ok = pk.VerifyPoseidon(msg, sig)
	assert.Equal(t, true, ok)
}

func TestVerifyPoseidon2(t *testing.T) {
	msg := utils.NewIntFromString("69588426711107115100232500042334179657931174539151555867956034570704220523596")
	msg = msg.Mod(msg, constants.Q)
	key := utils.NewIntFromString("56869496543825")

	var k PrivateKey
	k = utils.BigIntLEBytes(key)
	pk := k.Public()
	assert.Equal(t,
		"9255092729144892245186624611131828247442112563544941131408300200214096116351",
		pk.X.String())
	assert.Equal(t,
		"8460370541846376796657659750509399834188652251932899797602116208684247832083",
		pk.Y.String())

	x, ok := big.NewInt(0).SetString("15162295769440257382486195264681544386788758457719201693385196316384812064800", 10)
	require.True(t, ok)
	y, ok := big.NewInt(0).SetString("2782493627416942909007936076956568507304418921277473381438986134099538816121", 10)
	require.True(t, ok)
	s, ok := big.NewInt(0).SetString("16835165705656063478925976830596286105859651486320548752684160221106715530538", 10)
	require.True(t, ok)

	sig := &Signature{
		&Point{
			X: x,
			Y: y,
		},
		s,
	}

	ok = pk.VerifyPoseidon(msg, sig)
	assert.Equal(t, true, ok)
}

func TestCompressDecompress(t *testing.T) {
	var k PrivateKey
	_, err := hex.Decode(k[:],
		[]byte("0001020304050607080900010203040506070809000102030405060708090001"))
	require.Nil(t, err)
	pk := k.Public()
	for i := 0; i < 64; i++ {
		msgBuf, err := hex.DecodeString(fmt.Sprintf("000102030405060708%02d", i))
		if err != nil {
			panic(err)
		}
		msg := utils.SetBigIntFromLEBytes(new(big.Int), msgBuf)
		sig := k.SignMimc7(msg)
		sigBuf := sig.Compress()
		sig2, err := new(Signature).Decompress(sigBuf)
		assert.Equal(t, nil, err)
		ok := pk.VerifyMimc7(msg, sig2)
		assert.Equal(t, true, ok)
	}
}

func TestSignatureCompScannerValuer(t *testing.T) {
	privK := NewRandPrivKey()
	var value driver.Valuer //nolint:gosimple this is done to ensure interface compatibility
	value = privK.SignPoseidon(big.NewInt(674238462)).Compress()
	scan := privK.SignPoseidon(big.NewInt(1)).Compress()
	fromDB, err := value.Value()
	assert.Nil(t, err)
	assert.Nil(t, scan.Scan(fromDB))
	assert.Equal(t, value, scan)
}

func TestSignatureScannerValuer(t *testing.T) {
	privK := NewRandPrivKey()
	var value driver.Valuer
	var scan sql.Scanner
	value = privK.SignPoseidon(big.NewInt(674238462))
	scan = privK.SignPoseidon(big.NewInt(1))
	fromDB, err := value.Value()
	assert.Nil(t, err)
	assert.Nil(t, scan.Scan(fromDB))
	assert.Equal(t, value, scan)
}

func TestPublicKeyScannerValuer(t *testing.T) {
	privKValue := NewRandPrivKey()
	pubKValue := privKValue.Public()
	privKScan := NewRandPrivKey()
	pubKScan := privKScan.Public()
	var value driver.Valuer
	var scan sql.Scanner
	value = pubKValue
	scan = pubKScan
	fromDB, err := value.Value()
	assert.Nil(t, err)
	assert.Nil(t, scan.Scan(fromDB))
	assert.Equal(t, value, scan)
}

func TestPublicKeyCompScannerValuer(t *testing.T) {
	privKValue := NewRandPrivKey()
	pubKCompValue := privKValue.Public().Compress()
	privKScan := NewRandPrivKey()
	pubKCompScan := privKScan.Public().Compress()
	var value driver.Valuer
	var scan sql.Scanner
	value = &pubKCompValue
	scan = &pubKCompScan
	fromDB, err := value.Value()
	assert.Nil(t, err)
	assert.Nil(t, scan.Scan(fromDB))
	assert.Equal(t, value, scan)
}

func BenchmarkBabyjubEddsa(b *testing.B) {
	var k PrivateKey
	_, err := hex.Decode(k[:],
		[]byte("0001020304050607080900010203040506070809000102030405060708090001"))
	require.Nil(b, err)
	pk := k.Public()

	const n = 256

	msgBuf, err := hex.DecodeString("00010203040506070809")
	if err != nil {
		panic(err)
	}
	msg := utils.SetBigIntFromLEBytes(new(big.Int), msgBuf)
	var msgs [n]*big.Int
	for i := 0; i < n; i++ {
		msgs[i] = new(big.Int).Add(msg, big.NewInt(int64(i)))
	}
	var sigs [n]*Signature

	b.Run("SignMimc7", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			k.SignMimc7(msgs[i%n])
		}
	})

	for i := 0; i < n; i++ {
		sigs[i%n] = k.SignMimc7(msgs[i%n])
	}

	b.Run("VerifyMimc7", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pk.VerifyMimc7(msgs[i%n], sigs[i%n])
		}
	})

	b.Run("SignPoseidon", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			k.SignPoseidon(msgs[i%n])
		}
	})

	for i := 0; i < n; i++ {
		sigs[i%n] = k.SignPoseidon(msgs[i%n])
	}

	b.Run("VerifyPoseidon", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// 校验
			pk.VerifyPoseidon(msgs[i%n], sigs[i%n])
		}
	})
}
