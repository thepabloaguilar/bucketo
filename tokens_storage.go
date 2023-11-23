package bucketo

type TokensStorage interface {
	GetTokens() (int64, error)
	SetTokens(tokens int64) error
}

type InMemoryTokenStorage struct {
	tokens int64
}

func NewInMemoryTokenStorage(initialTokens int64) *InMemoryTokenStorage {
	return &InMemoryTokenStorage{
		tokens: initialTokens,
	}
}

func (s *InMemoryTokenStorage) GetTokens() (int64, error) {
	return s.tokens, nil
}

func (s *InMemoryTokenStorage) SetTokens(tokens int64) error {
	s.tokens = tokens

	return nil
}
