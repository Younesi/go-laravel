package cache

import (
	"bytes"
	"encoding/gob"
)

type Cache interface {
	Has(string) (bool, error)
	Get(string) (interface{}, error)
	Set(string, interface{}, ...int) error
	Forget(string) error
	EmptyByMatch(string) error
	Empty() error
}

type Entry map[string]interface{}

func encode(item Entry) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	if err := e.Encode(item); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func decode(str string) (Entry, error) {
	item := Entry{}
	b := bytes.Buffer{}
	b.Write([]byte(str))
	d := gob.NewDecoder(&b)
	if err := d.Decode(&item); err != nil {
		return nil, err
	}

	return item, nil
}
