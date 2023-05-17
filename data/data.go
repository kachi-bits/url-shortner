package data

type DataStore interface {
	SearchUrl(string) ([]map[string]string, error)
	Store(map[string]string) error
	Find(string) (map[string]string, error)
	Delete(string) (bool, error)
}
