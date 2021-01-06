package shortener

type RedirectService interface {
	Find(code string) (*Redirect, error)
	Store(r *Redirect) error
}
