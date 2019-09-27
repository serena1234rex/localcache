package localcache

import (
	"bytes"
	"encoding/gob"
)

func DefaultSerializeFunc(value interface{}) (interface{}, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(value)
	return buf.Bytes(), err
}

func DefaultDeserializeFunc(value interface{}) (interface{}, error) {
	dec := gob.NewDecoder(bytes.NewBuffer(value.([]byte)))
	var ret string
	err := dec.Decode(&ret)
	return ret, err
}