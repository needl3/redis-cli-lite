// Encoder is used to encode commands to RESP specification for send to finally send to redis server

package serializer

import (
	"fmt"
	"strconv"
	"strings"
)

type Encoder struct {
	CLRF         string
	I_ARR        string
	I_BULKSTRING string
	ARRAY_DECLR  func(nArgs int) string
}

var EncoderClient Encoder

func init() {
	CLRF := "\r\n"
	I_ARR := "*"
	I_BULKSTRING := "$"
	ARRAY_DECLR := func(nArgs int) string {
		return fmt.Sprintf("%s%s%s", I_ARR, strconv.Itoa(nArgs), CLRF)
	}
	EncoderClient = Encoder{
		CLRF,
		I_ARR,
		I_BULKSTRING,
		ARRAY_DECLR,
	}
}

func (encdr Encoder) Encode(cmd string) []byte {
	args := strings.Split(cmd, " ")
	nArgs := len(args)

	encoded := encdr.ARRAY_DECLR(nArgs)
	for _, c := range args {
		nStr := len(c)
		encoded += encdr.I_BULKSTRING + strconv.Itoa(nStr) + encdr.CLRF
		encoded += c + encdr.CLRF
	}
	return []byte(encoded)
}
