package api

// HTTP REST API сервера GoSearch.
// Прикладной интерфейс разработки для веб-приложения и других клиентов.

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/pprof"

	"gosearch/pkg/engine"

	"github.com/gorilla/mux"
)

// Service - служба API.
type Service struct {
	router *mux.Router
	engine *engine.Service
}

// ErrBadRequest - неверный запрос.
var ErrBadRequest = errors.New("неверный запрос")

// New - конструктор службы API.
func New(router *mux.Router, engine *engine.Service) *Service {
	s := Service{
		router: router,
		engine: engine,
	}
	s.endpoints()
	return &s
}

func (s *Service) endpoints() {
	// профилирование приложения
	s.router.HandleFunc("/debug/pprof/", pprof.Index)
	s.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	s.router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	s.router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	s.router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	s.router.Handle("/debug/pprof/block", pprof.Handler("block"))
	// поиск
	s.router.HandleFunc("/search/{query}", s.Search)
	// веб-приложение
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// Search ищет документы по запросу.
func (s *Service) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	result := s.engine.Search(mux.Vars(r)["query"])

	json.NewEncoder(w).Encode(result)
}
