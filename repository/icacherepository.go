package repository

type ICacheRepository interface {
	InsertTokens(tokens []string)
	StoreKeyUnionOfTokens(keyTop string, t string, keys []string) error
	GetCountOfTokensInSortedSet(key string) (int, error)
	GetTopValuesOfSortedSet(key string, n string) ([]string, error)
	ExpireKey(key string, t int) error
}
