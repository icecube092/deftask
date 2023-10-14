package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"deftask/internal/service"

	"github.com/go-chi/chi/v5"
)

type Server interface {
	Run()
	Shutdown(ctx context.Context)
}

type server struct {
	svc        service.Service
	httpServer *http.Server
}

func New(address string, svc service.Service) Server {
	s := &server{
		svc:        svc,
		httpServer: nil,
	}

	router := chi.NewRouter()
	router.With(s.timeMeasure).
		Get("/{userID_1}/{userID_2}", s.handle)

	httpServer := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: time.Second,
	}

	s.httpServer = httpServer

	return s
}

func (s *server) Run() {
	fmt.Println("server start at", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {

	}
}

func (s *server) Shutdown(ctx context.Context) {
	fmt.Println("server stopping")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Println("shutdown server error", err)
	}
}

func (s *server) timeMeasure(h http.Handler) http.Handler {
	return http.Handler(
		http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				now := time.Now()

				h.ServeHTTP(writer, request)

				fmt.Println("request time:", time.Since(now))
			},
		),
	)
}

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	userID1Str := chi.URLParam(r, "userID_1")
	userID2Str := chi.URLParam(r, "userID_2")

	userID1, err := strconv.ParseInt(userID1Str, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID2, err := strconv.ParseInt(userID2Str, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dupl, err := s.svc.IsUserDuplicate(r.Context(), userID1, userID2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"dupes": dupl,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
