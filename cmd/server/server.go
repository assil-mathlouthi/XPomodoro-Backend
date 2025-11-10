package server

import (
	"backend/middleware"
	"backend/services/pomodoros"
	"backend/services/ranking"
	"backend/services/stats"
	"backend/services/user"
	"database/sql"
	"log"
	"net/http"

	_ "backend/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	addr string
	db   *sql.DB
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	// Add Swagger (no auth required)
	subrouter.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/api/v1/swagger/doc.json"), // The url pointing to API definition
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
	// Create authenticated subrouter with JWT middleware
	authSubrouter := subrouter.PathPrefix("").Subrouter()
	authSubrouter.Use(middleware.JWTMiddleware)

	// Register user routes (some may need auth, some may not)
	userRepo := user.NewUserRepoImpl(s.db)
	userHandler := user.NewHandler(userRepo)
	userHandler.RegisterRoutes(subrouter, authSubrouter)

	// Register pomodoro routes (protected)
	pomodoroRepo := pomodoros.NewPomodoroRepoImpl(s.db)
	pomodoroHandler := pomodoros.NewHandler(pomodoroRepo)
	pomodoroHandler.RegisterRoutes(authSubrouter)

	// Register stats routes (protected)
	statsRepo := stats.NewStatsRepoImpl(s.db)
	statsHandler := stats.NewHandler(statsRepo)
	statsHandler.RegisterRoutes(authSubrouter)

	// Register ranking routes (protected)
	rankRepo := ranking.NewRankingRepoImpl(s.db)
	rankHandler := ranking.NewHandler(rankRepo)
	rankHandler.RegisterRoutes(authSubrouter)

	// --------------------------------------
	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, middleware.EnableCORS(router))
}
