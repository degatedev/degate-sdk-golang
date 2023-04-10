package model

import (
	"encoding/json"
	"regexp"
)

var (
	reETHAddr = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	reVolume  = regexp.MustCompile("^[+]{0,1}(\\d+)$")
)

func Copy(to, from interface{}) (err error) {
	b, err := json.Marshal(from)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, to)
	if err != nil {
		return
	}
	return nil
}

func IsETHAddress(address string) bool {
	return reETHAddr.MatchString(address)
}

func IsVolumeLegal(number string) bool {
	return reVolume.MatchString(number)
}
