package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/vukit/gophkeeper/internal/server/handlers"
	"github.com/vukit/gophkeeper/internal/server/logger"
	"github.com/vukit/gophkeeper/internal/server/model"
	"github.com/vukit/gophkeeper/internal/server/repositories/postgresql"
)

// handler структура обработчика HTTP запросов
type handler struct {
	tokenAuth *jwtauth.JWTAuth
	repoDB    handlers.RepoDB
	repoFile  handlers.RepoFile
	mLogger   *logger.Logger
}

var ErrNotFindUserID = errors.New("not find user id")

// NewHandler возвращает обработчик HTTP запросов
func NewHandler(tokenAuth *jwtauth.JWTAuth, repoDB handlers.RepoDB, repoFile handlers.RepoFile, mLogger *logger.Logger) handler {
	return handler{
		tokenAuth: tokenAuth,
		repoDB:    repoDB,
		repoFile:  repoFile,
		mLogger:   mLogger,
	}
}

// Index endpoint индексной страницы
func (h *handler) Index(w http.ResponseWriter, r *http.Request) {
	result := strings.Builder{}
	result.WriteString(`
	<!doctype html>
	<html lang="en">
	<head>
	  <meta charset="utf-8">
	  <title>GophKeeper Index Page</title>
	  <style>
	  	body {
			font-family: sans-serif;
			text-align: center;
		}
	  	.title {
			color: #3b4151;
			font-size: 24px;
			margin: 20px;
		}
		a {
			color: #4990e2;
			font-size: 14px;
			text-decoration: none;
			margin: 10px;
		}
	  </style>
	</head>
	<body>
	<h2 class="title">GophKeeper Index Page</h1>
	<a href="/swagger/index.html">Swagger API documetation</a>
	</body>
	</html>`)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintln(w, result.String())
}

// SignUp endpoint регистрации пользователя
//
// @Tags        User
// @Summary     Регистрация пользователя
// @Param       value body model.User true "данные пользователя"
// @Accept      json
// @Produce     json
// @Success     200 {object} object
// @Failure     400 {object} model.ErrorResponse
// @Failure     406 {object} model.ErrorResponse
// @Router /signup [post]
func (h *handler) SignUp(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		user, err := getUserFromBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		if err = user.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		userID, err := h.repoDB.SaveUser(ctx, user)
		if err != nil {
			switch {
			case errors.Is(err, postgresql.ErrDBUsernameIsAlreadyTaken):
				w.WriteHeader(http.StatusConflict)
			default:
				w.WriteHeader(http.StatusUnauthorized)
			}

			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		if err = h.setJWToken(w, userID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		fmt.Fprintf(w, "{}")
	}
}

// SignIn endpoint аутентификации пользователя
//
// @Tags        User
// @Summary     Аутентификация пользователя
// @Param       value body model.User true "данные пользователя"
// @Accept      json
// @Produce     json
// @Success     200 {object} object
// @Failure     400 {object} model.ErrorResponse
// @Failure     406 {object} model.ErrorResponse
// @Router /signin [post]
func (h *handler) SignIn(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		user, err := getUserFromBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		if err = user.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		userID, err := h.repoDB.FindUser(ctx, user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		if err = h.setJWToken(w, userID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		fmt.Fprintf(w, "{}")
	}
}

// SaveLogin endpoint сохранения данных логина пользователя
//
// @Tags        Logins
// @Summary     Cохраняет данные логина пользователя
// @Param       value body model.Login true "данные логина"
// @Accept      json
// @Produce     json
// @Success     200 {object} object
// @Failure     400 {object} model.ErrorResponse
// @Failure     406 {object} model.ErrorResponse
// @Router /logins [post]
func (h *handler) SaveLogin(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		login := model.Login{UserID: userID}

		decoder := json.NewDecoder(r.Body)

		err = decoder.Decode(&login)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		if err = login.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		err = h.repoDB.SaveLogin(ctx, &login)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		fmt.Fprintf(w, "{}")
	}
}

// DeleteLogin endpoint удаления данных логина пользователя
//
// @Tags        Logins
// @Summary     Удаляет данные логина пользователя
// @Param       id path integer true "id логина"
// @Accept      json
// @Produce     json
// @Success     200 {object} object
// @Failure     400 {object} model.ErrorResponse
// @Failure     406 {object} model.ErrorResponse
// @Router /logins/{id} [delete]
func (h *handler) DeleteLogin(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: "invalid login id = " + chi.URLParam(r, "id")})

			return
		}

		login := model.Login{ID: id, UserID: userID}

		err = h.repoDB.DeleteLogin(ctx, &login)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		fmt.Fprintf(w, "{}")
	}
}

// FindLogins endpoint возвращает данные логинов пользователя
//
// @Tags        Logins
// @Summary     Возвращает данные логинов пользователя
// @Accept      json
// @Produce     json
// @Success     200 {array}  model.Login
// @Failure     500 {object} model.ErrorResponse
// @Router /logins [get]
func (h *handler) FindLogins(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		logins, err := h.repoDB.FindLogins(ctx, model.User{ID: userID})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		body := &bytes.Buffer{}
		encoder := json.NewEncoder(body)

		err = encoder.Encode(logins)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		_, err = w.Write(body.Bytes())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}
	}
}

// SaveCard endpoint сохраняет данные банковской карты пользователя
//
// @Tags        Cards
// @Summary     Cохраняет данные банковской карты пользователя
// @Param       value body model.Card true "данные банковской карты"
// @Accept      json
// @Produce     json
// @Success     200 {object} object
// @Failure     400 {object} model.ErrorResponse
// @Failure     406 {object} model.ErrorResponse
// @Router /cards [post]
func (h *handler) SaveCard(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		card := model.Card{UserID: userID}

		decoder := json.NewDecoder(r.Body)

		err = decoder.Decode(&card)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		if err = card.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		err = h.repoDB.SaveCard(ctx, &card)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		fmt.Fprintf(w, "{}")
	}
}

// DeleteCard endpoint удаляет данные банковской карты пользователя
//
// @Tags        Cards
// @Summary     Удаляет данные банковской карты пользователя
// @Param       id path integer true "id банковской карты"
// @Accept      json
// @Produce     json
// @Success     200 {object} object
// @Failure     400 {object} model.ErrorResponse
// @Failure     406 {object} model.ErrorResponse
// @Router /cards/{id} [delete]
func (h *handler) DeleteCard(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: "invalid login id = " + chi.URLParam(r, "id")})

			return
		}

		card := model.Card{ID: id, UserID: userID}

		err = h.repoDB.DeleteCard(ctx, &card)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		fmt.Fprintf(w, "{}")
	}
}

// FindCards endpoint возвращает данные банковских карт пользователя
//
// @Tags        Cards
// @Summary     Возвращает данные банковских карт пользователя
// @Accept      json
// @Produce     json
// @Success     200 {array}  model.Card
// @Failure     500 {object} model.ErrorResponse
// @Router /cards [get]
func (h *handler) FindCards(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		logins, err := h.repoDB.FindCards(ctx, model.User{ID: userID})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		body := &bytes.Buffer{}
		encoder := json.NewEncoder(body)

		err = encoder.Encode(logins)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		_, err = w.Write(body.Bytes())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}
	}
}

// SaveFile endpoint сохраняет данные файла пользователя
//
// @Tags        Files
// @Summary     Cохраняет данные файла пользователя
// @Param   	id formData integer true  "id файла"
// @Param   	metainfo formData string true  "metainfo файла"
// @Param   	file formData file true  "содержимое файла"
// @Accept      multipart/form-data
// @Produce     json
// @Success     200 {object} object
// @Failure     400 {object} model.ErrorResponse
// @Failure     406 {object} model.ErrorResponse
// @Router /files [post]
func (h *handler) SaveFile(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		fileID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: "invalid file id = " + r.FormValue("id")})

			return
		}

		file := &model.File{ID: fileID, UserID: userID, MetaInfo: r.FormValue("metainfo")}

		oldFile, err := h.repoDB.FindFile(ctx, file.ID, file.UserID)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		if oldFile.Path == "" && fileID != 0 {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: "file not exist"})

			return
		}

		file.Path = oldFile.Path
		file.Name = oldFile.Name

		mpFile, mpFH, err := r.FormFile("file")
		if err == nil {
			defer mpFile.Close()

			errUpdateFile := h.repoFile.DeleteFile(ctx, oldFile.Path)
			if errUpdateFile != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				fmt.Fprint(w, model.ErrorResponse{Error: errUpdateFile.Error()})

				return
			}

			newFilePath, errUpdateFile := h.repoFile.SaveFile(ctx, mpFile)
			if errUpdateFile != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				fmt.Fprint(w, model.ErrorResponse{Error: errUpdateFile.Error()})

				return
			}

			file.Path = newFilePath
			file.Name = mpFH.Filename
		}

		if err = file.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		err = h.repoDB.SaveFile(ctx, file)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		fmt.Fprintf(w, "{}")
	}
}

// DeleteFile endpoint удаления данных файла пользователя
//
// @Tags        Files
// @Summary     Удаляет данные файла пользователя
// @Param       id path integer true "id файла"
// @Accept      json
// @Produce     json
// @Success     200 {object} object
// @Failure     400 {object} model.ErrorResponse
// @Failure     406 {object} model.ErrorResponse
// @Router /files/{id} [delete]
func (h *handler) DeleteFile(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: "invalid file id = " + chi.URLParam(r, "id")})

			return
		}

		file, err := h.repoDB.FindFile(ctx, id, userID)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		err = h.repoFile.DeleteFile(ctx, file.Path)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		err = h.repoDB.DeleteFile(ctx, file)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		fmt.Fprintf(w, "{}")
	}
}

// FindFiles endpoint возвращает данные файлов пользователя
//
// @Tags        Files
// @Summary     Возвращает данные файлов пользователя
// @Accept      json
// @Produce     json
// @Success     200 {array}  model.File
// @Failure     500 {object} model.ErrorResponse
// @Router /files [get]
func (h *handler) FindFiles(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		files, err := h.repoDB.FindFiles(ctx, model.User{ID: userID})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		body := &bytes.Buffer{}
		encoder := json.NewEncoder(body)

		err = encoder.Encode(files)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		_, err = w.Write(body.Bytes())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}
	}
}

// FindFile endpoint выгружает файл пользователю
//
// @Tags        Files
// @Summary     Выгружает файл пользователю
// @Param       id path integer true "id файла"
// @Accept      json
// @Produce     octet-stream
// @Produce     json
// @Success     200 {file} schema
// @Failure     400 {object} model.ErrorResponse
// @Failure     406 {object} model.ErrorResponse
// @Router /files/{id} [get]
func (h *handler) FindFile(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userID, err := getUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, model.ErrorResponse{Error: "invalid file id = " + chi.URLParam(r, "id")})

			return
		}

		file, err := h.repoDB.FindFile(ctx, id, userID)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}

		src, err := h.repoFile.GetFile(ctx, file.Path)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, model.ErrorResponse{Error: err.Error()})

			return
		}
		defer src.Close()

		w.Header().Set("Content-Type", "application/octet-stream")

		_, err = io.Copy(w, src)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)

			return
		}
	}
}

func getUserFromBody(r *http.Request) (user model.User, err error) {
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&user)

	return
}

func (h *handler) setJWToken(w http.ResponseWriter, userID int) (err error) {
	claims := map[string]interface{}{"user_id": strconv.Itoa(userID)}
	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour))

	_, tokenString, err := h.tokenAuth.Encode(claims)
	if err == nil {
		http.SetCookie(w, &http.Cookie{Name: "jwt", Value: tokenString, Path: "/"})
	}

	return
}

func getUserID(r *http.Request) (id int, err error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return 0, err
	}

	userID, ok := claims["user_id"]
	if !ok {
		return 0, ErrNotFindUserID
	}

	id, err = strconv.Atoi(userID.(string))
	if err != nil {
		return 0, ErrNotFindUserID
	}

	return
}
