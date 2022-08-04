package repo

type RepositoryStorage interface {
	Set(key, value string) error
	SetWithTTL(key, value string, seconds int64) error
	Get(key string) (interface{}, error)
	Delete(key string) (interface{}, error)
	Search(key string) (interface{}, error)
}
