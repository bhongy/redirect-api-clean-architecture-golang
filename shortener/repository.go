package shortener

// Repository is an interface to connect business logic to repository

type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
