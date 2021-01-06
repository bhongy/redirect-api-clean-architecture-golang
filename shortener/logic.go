package shortener

import (
	"errors"
	"fmt"
	"time"

	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	ErrRedirectInvalid  = errors.New("Redirect Invalid")
)

type redirectService struct {
	repo RedirectRepository
}

func NewRedirectService(repo RedirectRepository) RedirectService {
	return &redirectService{repo}
}

func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.repo.Find(code)
}

func (r *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return fmt.Errorf("service.Redirect.Store: %v", ErrRedirectInvalid)
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now()
	return r.repo.Store(redirect)
}
