package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gourses/demo/storage"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/exp/slog"
)

type API struct {
	s   *storage.Storage
	srv *http.Server
}

func New(addr string, s *storage.Storage) *API {
	r := httprouter.New()
	a := &API{
		s: s,
		srv: &http.Server{
			Addr:              addr,
			Handler:           r,
			ReadHeaderTimeout: time.Second * 5,
		},
	}
	r.GET("/notes", a.GetNotes)
	return a
}

func (a *API) Run() error {
	err := a.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (a *API) Close() {
	a.srv.Close()
}

func (a *API) GetNotes(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	notes, err := a.s.GetNotes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		slog.Error(err.Error())
	}
}
