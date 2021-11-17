package types

import (
	"bytes"
	"strconv"
)

//
// byte parsing methods for TradeMsg json bytes
//

const (
	msgTypeSkipSep   = 3
	msgVolumeSkipSep = 21
	msgPriceSkipSep  = 25
	msgProductSkipSep  = 29
	tokenSep         = '"'
)

type ParseToken struct {
	Token      string
	SepSkipCnt int
}

func ParseType(msg []byte) (string, int) {
	return ParseString(tokenSep, msgTypeSkipSep, msg)
}

func ParseProductID(msg []byte) (string, int) {
	return ParseString(tokenSep, msgProductSkipSep, msg)
}

func ParsePrice(msg []byte) (float64, int) {
	return ParseF64(tokenSep, msgPriceSkipSep, msg)
}

func ParseVolume(msg []byte) (float64, int) {
	return ParseF64(tokenSep, msgVolumeSkipSep, msg)
}

func ParseString(tokenSep byte, skipCnt int, msg []byte) (string, int){
	val, startIdx := parseVal(tokenSep, skipCnt, msg)
	if startIdx == -1 {
		return "", -1
	}

	return string(val), startIdx
}

func ParseF64(tokenSep  byte, skipCnt int, msg []byte) (float64, int){
	val, startIdx := ParseString(tokenSep, skipCnt, msg)
	if startIdx == -1 {
		return -1, -1
	}

	f64Val, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return -1, -1
	}

	return f64Val, startIdx
}

func parseVal(tokenSep byte, skipCnt int, msg []byte) ([]byte, int) {
	if len(msg) == 0 {
		return nil, -1
	}

	var (
		accIdx, beginIdx, endIdx int
	)
	for i := 0; i < skipCnt; i++ {
		beginIdx = bytes.IndexByte(msg[accIdx:], tokenSep)
		if beginIdx == -1 {
			return nil, -1
		}
		accIdx += beginIdx + 1
	}
	endIdx = bytes.IndexByte(msg[accIdx+1:], tokenSep)
	if endIdx == -1 {
		return nil, -1
	}
	return msg[accIdx:accIdx+endIdx+1], accIdx
}