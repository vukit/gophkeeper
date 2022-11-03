package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strconv"

	"github.com/vukit/gophkeeper/internal/client/config"
	"github.com/vukit/gophkeeper/internal/client/logger"
	"github.com/vukit/gophkeeper/internal/client/model"
)

type httpService struct {
	client  *http.Client
	baseURL string
	mLogger *logger.Logger
	cs      *CryptoService
}

// NewHTTPService возваращает сервис с методам для обмена данными с сервером по HTTP протоколу
func NewHTTPService(cfg *config.Config, mLogger *logger.Logger) *httpService {
	service := &httpService{}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	service.client = &http.Client{Jar: jar}

	service.baseURL = fmt.Sprintf("%s://%s/api", cfg.ServerProtocol, cfg.ServerAddress)

	service.mLogger = mLogger

	return service
}

// SetCryptoService устанавливает сервис симметричного шифрования
func (s *httpService) SetCryptoService(cs *CryptoService) {
	s.cs = cs
}

// SignIn метод аутентификация пользователя
func (s *httpService) SignIn(ctx context.Context, user *model.User) (err error) {
	body := &bytes.Buffer{}
	encoder := json.NewEncoder(body)

	err = encoder.Encode(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/signin", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// SignUp метод регистрации пользователя
func (s *httpService) SignUp(ctx context.Context, user *model.User) (err error) {
	body := &bytes.Buffer{}
	encoder := json.NewEncoder(body)

	err = encoder.Encode(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/signup", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// SaveLogin метод сохранения данных логина пользователя
func (s *httpService) SaveLogin(ctx context.Context, login *model.Login) (err error) {
	login.Username = hex.EncodeToString((s.cs.Encrypt([]byte(login.Username))))
	login.Password = hex.EncodeToString((s.cs.Encrypt([]byte(login.Password))))
	login.MetaInfo = hex.EncodeToString((s.cs.Encrypt([]byte(login.MetaInfo))))

	body := &bytes.Buffer{}
	encoder := json.NewEncoder(body)

	err = encoder.Encode(login)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/logins", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// DeleteLogin метод удаления данных логина пользователя
func (s *httpService) DeleteLogin(ctx context.Context, login *model.Login) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s%s%d", s.baseURL, "/logins/", login.ID), &bytes.Buffer{})
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// GetLogins метод возвращает данные логинов пользователя
func (s *httpService) GetLogins(ctx context.Context) (logins []model.Login, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/logins", &bytes.Buffer{})
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	logins = make([]model.Login, 0)

	err = decoder.Decode(&logins)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(logins); i++ {
		data, err := hex.DecodeString(logins[i].Username)
		if err != nil {
			return nil, fmt.Errorf("error decrypted login with id = %d: %w", logins[i].ID, err)
		}

		data, err = s.cs.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("error decrypted login with id = %d: %w", logins[i].ID, err)
		}

		logins[i].Username = string(data)

		data, err = hex.DecodeString(logins[i].Password)
		if err != nil {
			return nil, fmt.Errorf("error decrypted login with id = %d: %w", logins[i].ID, err)
		}

		data, err = s.cs.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("error decrypted login with id = %d: %w", logins[i].ID, err)
		}

		logins[i].Password = string(data)

		data, err = hex.DecodeString(logins[i].MetaInfo)
		if err != nil {
			return nil, fmt.Errorf("error decrypted login with id = %d: %w", logins[i].ID, err)
		}

		data, err = s.cs.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("error decrypted login with id = %d: %w", logins[i].ID, err)
		}

		logins[i].MetaInfo = string(data)
	}

	return logins, nil
}

// SaveCard метод сохранения данных банковской карты пользователя
func (s *httpService) SaveCard(ctx context.Context, card *model.Card) (err error) {
	card.Bank = hex.EncodeToString((s.cs.Encrypt([]byte(card.Bank))))
	card.Number = hex.EncodeToString((s.cs.Encrypt([]byte(card.Number))))
	card.Date = hex.EncodeToString((s.cs.Encrypt([]byte(card.Date))))
	card.CVV = hex.EncodeToString((s.cs.Encrypt([]byte(card.CVV))))
	card.MetaInfo = hex.EncodeToString((s.cs.Encrypt([]byte(card.MetaInfo))))

	body := &bytes.Buffer{}
	encoder := json.NewEncoder(body)

	err = encoder.Encode(card)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/cards", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCard метод удаления данных банковской карты пользователя
func (s *httpService) DeleteCard(ctx context.Context, card *model.Card) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s%s%d", s.baseURL, "/cards/", card.ID), &bytes.Buffer{})
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// GetCards метод возвращает данные банковских карт пользователя
func (s *httpService) GetCards(ctx context.Context) (cards []model.Card, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/cards", &bytes.Buffer{})
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	cards = make([]model.Card, 0)

	err = decoder.Decode(&cards)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(cards); i++ {
		data, err := hex.DecodeString(cards[i].Bank)
		if err != nil {
			return nil, fmt.Errorf("error decrypted card with id = %d: %w", cards[i].ID, err)
		}

		data, err = s.cs.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("error decrypted card with id = %d: %w", cards[i].ID, err)
		}

		cards[i].Bank = string(data)

		data, err = hex.DecodeString(cards[i].Number)
		if err != nil {
			return nil, fmt.Errorf("error decrypted card with id = %d: %w", cards[i].ID, err)
		}

		data, err = s.cs.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("error decrypted card with id = %d: %w", cards[i].ID, err)
		}

		cards[i].Number = string(data)

		data, err = hex.DecodeString(cards[i].Date)
		if err != nil {
			return nil, fmt.Errorf("error decrypted card with id = %d: %w", cards[i].ID, err)
		}

		data, err = s.cs.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("error decrypted card with id = %d: %w", cards[i].ID, err)
		}

		cards[i].Date = string(data)

		data, err = hex.DecodeString(cards[i].CVV)
		if err != nil {
			return nil, fmt.Errorf("error decrypted card with id = %d: %w", cards[i].ID, err)
		}

		data, err = s.cs.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("error decrypted card with id = %d: %w", cards[i].ID, err)
		}

		cards[i].CVV = string(data)

		data, err = hex.DecodeString(cards[i].MetaInfo)
		if err != nil {
			return nil, fmt.Errorf("error decrypted card with id = %d: %w", cards[i].ID, err)
		}

		data, err = s.cs.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("error decrypted card with id = %d: %w", cards[i].ID, err)
		}

		cards[i].MetaInfo = string(data)
	}

	return cards, nil
}

// SaveFile метод сохранения данных файла пользователя
func (s *httpService) SaveFile(ctx context.Context, file *model.File) (err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	var (
		part          io.Writer
		encryptedFile io.Reader
	)

	if _, err = os.Stat(file.Path); err == nil {
		src, _ := os.Open(file.Path)
		defer src.Close()

		filename := hex.EncodeToString((s.cs.Encrypt([]byte(filepath.Base(src.Name())))))

		part, err = writer.CreateFormFile("file", filename)
		if err != nil {
			return err
		}

		encryptedFile, err = s.cs.EncryptFile(src)
		if err != nil {
			return err
		}

		_, err = io.Copy(part, encryptedFile)
		if err != nil {
			return err
		}
	}

	part, err = writer.CreateFormField("id")
	if err != nil {
		return err
	}

	_, err = part.Write([]byte(strconv.Itoa(file.ID)))
	if err != nil {
		return err
	}

	part, err = writer.CreateFormField("metainfo")
	if err != nil {
		return err
	}

	metaInfo := hex.EncodeToString((s.cs.Encrypt([]byte(file.MetaInfo))))

	_, err = part.Write([]byte(metaInfo))
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/files", body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile метод удаления данных файла пользователя
func (s *httpService) DeleteFile(ctx context.Context, file *model.File) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s%s%d", s.baseURL, "/files/", file.ID), &bytes.Buffer{})
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// GetFiles метод возвращает данные файлов пользователя
func (s *httpService) GetFiles(ctx context.Context) (files []model.File, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/files", &bytes.Buffer{})
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	files = make([]model.File, 0)

	err = decoder.Decode(&files)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(files); i++ {
		data, err := hex.DecodeString(files[i].MetaInfo)
		if err != nil {
			return nil, fmt.Errorf("error decrypted file with id = %d: %w", files[i].ID, err)
		}

		data, err = s.cs.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("error decrypted file with id = %d: %w", files[i].ID, err)
		}

		files[i].MetaInfo = string(data)

		data, err = hex.DecodeString(files[i].Name)
		if err != nil {
			return nil, fmt.Errorf("error decrypted file with id = %d: %w", files[i].ID, err)
		}

		data, err = s.cs.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("error decrypted file with id = %d: %w", files[i].ID, err)
		}

		files[i].Name = string(data)
	}

	return files, nil
}

// DownloadFile метод скачивания файла пользователя
func (s *httpService) DownloadFile(ctx context.Context, file *model.File, downloadFolder string) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s%d", s.baseURL, "/files/", file.ID), &bytes.Buffer{})
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp.StatusCode, resp.Body)
	if err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(downloadFolder, file.Name))
	if err != nil {
		return err
	}
	defer dst.Close()

	decryptedBody, err := s.cs.DecryptFile(resp.Body)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, decryptedBody)
	if err != nil {
		return err
	}

	return nil
}

func checkStatusCode(statusCode int, body io.ReadCloser) error {
	if statusCode == http.StatusUnauthorized {
		return errors.New(http.StatusText(statusCode))
	}

	if statusCode != http.StatusOK {
		var answer struct {
			Error string
		}

		decoder := json.NewDecoder(body)

		err := decoder.Decode(&answer)
		if err != nil {
			return err
		}

		return errors.New(answer.Error)
	}

	return nil
}
