package serializer

type Searializer struct {
	Parser  Parser
	Encoder Encoder
}

func New() *Searializer {
	return &Searializer{
		Parser:  ParserClient,
		Encoder: EncoderClient,
	}
}
