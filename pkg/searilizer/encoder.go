// Encoder is used to encode commands to RESP specification for send to finally send to redis server

package serializer

import "fmt"

type Encoder struct {
	CLRF string
}

var EncoderClient Encoder

func init() {
	EncoderClient = Encoder{
		CLRF: "\r\n",
	}
}

func (encdr Encoder) Encode(cmd string) []byte {
	encoded := fmt.Sprintf("*1%s$4%s%s%s", encdr.CLRF, encdr.CLRF, cmd, encdr.CLRF)
	return []byte(encoded)
}
