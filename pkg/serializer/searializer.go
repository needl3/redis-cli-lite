package serializer

type Searializer struct {
	Parser  Parser
	Encoder Encoder
	Printer interface{} // TODO: Implement to pretty print
}

func New() *Searializer {
	return &Searializer{
		Parser:  ParserClient,
		Encoder: EncoderClient,
	}
}
