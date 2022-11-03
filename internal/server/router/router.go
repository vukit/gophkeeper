package router

import (
	"context"
	"crypto/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/jwtauth"
	"github.com/vukit/gophkeeper/internal/server/handlers"
	handler "github.com/vukit/gophkeeper/internal/server/handlers/http"
	"github.com/vukit/gophkeeper/internal/server/logger"

	// Подключение Swagger
	_ "github.com/vukit/gophkeeper/internal/server/swagger"
)

// NewRouter возвращает маршрутизатор для HTTP запросов
func NewRouter(ctx context.Context, repoDB handlers.RepoDB, repoFile handlers.RepoFile, mLogger *logger.Logger) (r chi.Router, err error) {
	secret, err := generateSecret(32)
	if err != nil {
		return nil, err
	}

	tokenAuth := jwtauth.New("HS256", secret, nil)

	r = chi.NewRouter()

	r.Use(middleware.Compress(5))
	r.Mount("/swagger", httpSwagger.WrapHandler)

	h := handler.NewHandler(tokenAuth, repoDB, repoFile, mLogger)

	r.Handle("/", http.FileServer(http.Dir("./internal/server/static")))

	r.Post("/api/signup", h.SignUp(ctx))

	r.Post("/api/signin", h.SignIn(ctx))

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/api/logins", h.SaveLogin(ctx))
		r.Delete("/api/logins/{id}", h.DeleteLogin(ctx))
		r.Get("/api/logins", h.FindLogins(ctx))
		r.Post("/api/cards", h.SaveCard(ctx))
		r.Delete("/api/cards/{id}", h.DeleteCard(ctx))
		r.Get("/api/cards", h.FindCards(ctx))
		r.Post("/api/files", h.SaveFile(ctx))
		r.Delete("/api/files/{id}", h.DeleteFile(ctx))
		r.Get("/api/files", h.FindFiles(ctx))
		r.Get("/api/files/{id}", h.FindFile(ctx))
	})

	return r, nil
}

func generateSecret(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
