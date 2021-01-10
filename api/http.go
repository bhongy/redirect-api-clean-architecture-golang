package api

import (
	"io/ioutil"
	"log"
	"net/http"

	jsonSerializer "github.com/bhongy/tmp-clean-arch-golang/serializer/json"
	msgpackSerializer "github.com/bhongy/tmp-clean-arch-golang/serializer/msgpack"
	"github.com/bhongy/tmp-clean-arch-golang/shortener"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService shortener.RedirectService
}

func NewHandler(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{redirectService}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Is(err, shortener.ErrRedirectNotFound) {
			status := http.StatusNotFound
			http.Error(w, http.StatusText(status), status)
			return
		}
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	redirect, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Is(err, shortener.ErrRedirectNotFound) {
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	responseBody, err := h.serializer(contentType).Encode(redirect)
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) serializer(contentType string) shortener.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &msgpackSerializer.Redirect{}
	}
	return &jsonSerializer.Redirect{}
}
