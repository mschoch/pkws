package pkws

// Storage is the interface for private storage within a room
type Storage interface {
	Get(key string) ([]byte, error)
	Put(key string, val []byte) error
	Delete(key string) (bool, error)
	DeleteAll() error
}

// EphemeralStorage is an in-memory only implementation
// of the Storage interface
type EphemeralStorage struct {
	data map[string][]byte
}

func NewEphemeralStorage() *EphemeralStorage {
	return &EphemeralStorage{
		data: make(map[string][]byte),
	}
}

func (e *EphemeralStorage) Get(key string) ([]byte, error) {
	if rv, ok := e.data[key]; ok {
		return rv, nil
	}
	return nil, nil
}

func (e *EphemeralStorage) Put(key string, val []byte) error {
	e.data[key] = val
	return nil
}

func (e *EphemeralStorage) Delete(key string) (bool, error) {
	_, existed := e.data[key]
	delete(e.data, key)
	return existed, nil
}

func (e *EphemeralStorage) DeleteAll() error {
	e.data = make(map[string][]byte)
	return nil
}
