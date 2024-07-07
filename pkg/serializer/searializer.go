package serializer

type Searializer struct {
	Parser  Parser
	Encoder Encoder
	Printer interface{}
}

func New() *Searializer {
	return &Searializer{
		Parser:  ParserClient,
		Encoder: EncoderClient,
	}
}
