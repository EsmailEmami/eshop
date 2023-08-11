package helpers

type KeyValueResponse[K any, V any] struct {
	Key   K `json:"key"`
	Value V `json:"value"`
}

func NewKeyValueResponse[K any, V any](key K, value V) *KeyValueResponse[K, V] {
	return &KeyValueResponse[K, V]{
		Key:   key,
		Value: value,
	}
}
