package util

import (
	"math/big"
	"regexp"
	"strconv"
	"time"
)

var zeroHash = regexp.MustCompile("^0?x?0+$")

func IsZeroHash(s string) bool {
	return zeroHash.MatchString(s)
}

func ToHex(n int64) string {
	return "0x0" + strconv.FormatInt(n, 16)
}

func FormatBigToString(reward *big.Int) string {
	return reward.String()
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}

func String2Big(num string) *big.Int {
	n := new(big.Int)
	n.SetString(num, 0)
	return n
}
