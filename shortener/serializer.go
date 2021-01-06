package shortener

type RedirectSerializer interface {
	Decode(b []byte) (*Redirect, error)
	Encode(redirect *Redirect) ([]byte, error)
}
