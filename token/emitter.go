package token

type Emitter interface {
	Emit(token *Token)
}

type EmitterFunc func(token *Token)

func (fn EmitterFunc) Emit(token *Token) {
	fn(token)
}

type Iterator interface {
	Next() *Token
}

type IteratorFunc func() *Token

func (fn IteratorFunc) Next() *Token {
	return fn()
}

type EmitterIterator struct {
	tokens      chan *Token
	tokenMapper map[Type]string
}

func (e *EmitterIterator) Emit(token *Token) {
	if name, ok := e.tokenMapper[token.Type]; ok {
		token.Name = name
	}
	e.tokens <- token
}

func (e EmitterIterator) Next() *Token {
	value, ok := <-e.tokens
	if !ok {
		return &Token{
			Type: EOF,
		}
	}

	return value
}

func NewEmitterIterator(tokenMapper map[Type]string) *EmitterIterator {
	return &EmitterIterator{
		tokens:      make(chan *Token, 2),
		tokenMapper: tokenMapper,
	}
}
