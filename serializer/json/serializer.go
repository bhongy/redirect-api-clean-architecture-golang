package json

import (
	"encoding/json"
	"fmt"

	"github.com/bhongy/rediret-api-clean-architecture-golang/shortener"
)

type Redirect struct{}

func (r *Redirect) Decode(b []byte) (*shortener.Redirect, error) {
	var redirect shortener.Redirect
	if err := json.Unmarshal(b, &redirect); err != nil {
		return nil, fmt.Errorf("serializer.Redirect.Decode: %v", err)
	}
	return &redirect, nil
}

func (r *Redirect) Encode(redirect *shortener.Redirect) ([]byte, error) {
	b, err := json.Marshal(redirect)
	if err != nil {
		return nil, fmt.Errorf("serializer.Redirect.Encode: %v", err)
	}
	return b, nil
}
