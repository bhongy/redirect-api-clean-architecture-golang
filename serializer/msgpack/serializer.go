package msgpack

import (
	"fmt"

	"github.com/bhongy/tmp-clean-arch-golang/shortener"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

func (r *Redirect) Decode(b []byte) (*shortener.Redirect, error) {
	var redirect shortener.Redirect
	if err := msgpack.Unmarshal(b, &redirect); err != nil {
		return nil, fmt.Errorf("serializer.Redirect.Decode: %v", err)
	}
	return &redirect, nil
}

func (r *Redirect) Encode(redirect *shortener.Redirect) ([]byte, error) {
	b, err := msgpack.Marshal(redirect)
	if err != nil {
		return nil, fmt.Errorf("serializer.Redirect.Encode: %v", err)
	}
	return b, nil
}
