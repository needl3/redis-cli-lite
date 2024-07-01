// Parser is used to parse RESP specific response given back by the redis server
// Into client specific representation

package serializer

import "strings"

type Parser struct {
	RETURN_TYPE string
}

var ParserClient Parser

func init() {
	ParserClient = Parser{
		RETURN_TYPE: "+",
	}
}

func (prsr Parser) Parse(resp []byte) string {
	return strings.Trim(string(resp), "+")
}
