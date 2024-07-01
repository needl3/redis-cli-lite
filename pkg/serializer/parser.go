// Parser is used to parse RESP specific response given back by the redis server
// Into client specific representation

package serializer

import "strings"

type Parser struct {
	OK_TYPE      string
	ERROR_TYPE   string
	CLRF         string
	I_ARR        string
	I_BULKSTRING string
}

var ParserClient Parser

func init() {
	ParserClient = Parser{
		OK_TYPE:      "+",
		ERROR_TYPE:   "-",
		CLRF:         "\r\n",
		I_ARR:        "*",
		I_BULKSTRING: "$",
	}
}

func (prsr Parser) Parse(resp []byte) string {
	// Time for some state machine
	return strings.Trim(string(resp), "+")
}
