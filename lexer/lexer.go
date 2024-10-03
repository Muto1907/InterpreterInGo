package lexer

type Lexer struct {
	input            string
	currCharPosition int
	readPosition     int // 1 char after currCharPosition
	char             byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	return l
}
