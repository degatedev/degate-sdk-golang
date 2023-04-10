package go_sha3

import (
	"crypto/sha256"
	"fmt"
)

// SoliditySHA3 solidity sha3
func GoSHA3(data ...interface{}) []byte {
	types, ok := data[0].([]string)
	if len(data) > 1 && ok {
		rest := data[1:]
		if len(rest) == len(types) {
			return gosha3(types, data[1:]...)
		}
		iface, ok := data[1].([]interface{})
		if ok {
			return gosha3(types, iface...)
		}
	}

	var v [][]byte
	for _, item := range data {
		b := parseBytes(item, -1)
		v = append(v, b)
	}
	return gosha3Legacy(v...)
}

// solsha3Legacy solidity sha3
func gosha3Legacy(data ...[]byte) []byte {
	//hash := sha3.NewLegacyKeccak256()
	hash := sha256.New()
	bs := concatByteSlices(data...)
	hash.Write(bs)
	return hash.Sum(nil)
}

// SoliditySHA3WithPrefix solidity sha3 with prefix
func SoliditySHA3WithPrefix(data []byte) []byte {
	result := GoSHA3(
		concatByteSlices(
			[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%v", len(data))),
			data,
		),
	)

	return result
}

// ConcatByteSlices concat byte slices
func ConcatByteSlices(arrays ...[]byte) []byte {
	return concatByteSlices(arrays...)
}

// solsha3 solidity sha3
func gosha3(types []string, values ...interface{}) []byte {
	var b [][]byte
	for i, typ := range types {
		b = append(b, pack(typ, values[i], false))
	}

	//hash := sha3.NewLegacyKeccak256()
	hash := sha256.New()
	bs := concatByteSlices(b...)
	hash.Write(bs)
	return hash.Sum(nil)
}

// Pack ...
func Pack(types []string, values []interface{}) []byte {
	if len(types) != len(values) {
		panic("type/value count mismatch")
	}

	var tight [][]byte
	for i, typ := range types {
		tight = append(tight, pack(typ, values[i], false))
	}

	return concatByteSlices(tight...)
}
