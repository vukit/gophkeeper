package postgresql

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/vukit/gophkeeper/internal/server/model"

	// Register pgx stdlib
	_ "github.com/jackc/pgx/v4/stdlib"
)

// RepoPostgreSQL структура PostgreSQL репозитория
type RepoPostgreSQL struct {
	db *sql.DB
}

// NewRepo возвращает PostgreSQL репозиторий
func NewRepo(dsn string) (repo RepoPostgreSQL, err error) {
	db, err := sql.Open("pgx", dsn)

	repo = RepoPostgreSQL{db: db}

	if err != nil {
		return repo, err
	}

	return repo, err
}

// SaveUser используется при регистрации пользователя
func (repo RepoPostgreSQL) SaveUser(ctx context.Context, user model.User) (userID int, err error) {
	if repo.db == nil {
		return 0, ErrDBNoDBConn
	}

	passwordHash := sha256.Sum256([]byte(user.Password))

	err = repo.db.QueryRowContext(ctx,
		`INSERT INTO users (username, password) VALUES($1, $2) RETURNING user_id`,
		user.Username,
		hex.EncodeToString(passwordHash[:])).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, ErrDBUsernameIsAlreadyTaken
		}

		return 0, err
	}

	return userID, err
}

// FindUser используется при аутентификации пользователя
func (repo RepoPostgreSQL) FindUser(ctx context.Context, user model.User) (userID int, err error) {
	if repo.db == nil {
		return 0, ErrDBNoDBConn
	}

	passwordHash := sha256.Sum256([]byte(user.Password))

	err = repo.db.QueryRowContext(ctx,
		`SELECT user_id FROM users WHERE username = $1 and password = $2`,
		user.Username,
		hex.EncodeToString(passwordHash[:])).Scan(&userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrDBInvalidUsernamePasswordPair
		default:
			return 0, err
		}
	}

	return userID, err
}

// SaveLogin используется при сохранении данных логина пользователя
func (repo RepoPostgreSQL) SaveLogin(ctx context.Context, login *model.Login) (err error) {
	if repo.db == nil {
		return ErrDBNoDBConn
	}

	if login.ID == 0 {
		_, err = repo.db.ExecContext(ctx,
			`INSERT INTO logins (user_id, username, password, metainfo) VALUES($1, $2, $3, $4)`,
			login.UserID,
			login.Username,
			login.Password,
			login.MetaInfo)
	} else {
		_, err = repo.db.ExecContext(ctx,
			`UPDATE logins SET username = $1, password = $2, metainfo = $3 WHERE login_id = $4 and user_id = $5`,
			login.Username,
			login.Password,
			login.MetaInfo,
			login.ID,
			login.UserID)
	}

	return err
}

// DeleteLogin используется при удалении данных логина пользователя
func (repo RepoPostgreSQL) DeleteLogin(ctx context.Context, login *model.Login) (err error) {
	if repo.db == nil {
		return ErrDBNoDBConn
	}

	_, err = repo.db.ExecContext(ctx,
		`DELETE FROM logins WHERE login_id = $1 and user_id = $2`,
		login.ID,
		login.UserID)

	return err
}

// FindLogins возвращает данные логинов пользователя
func (repo RepoPostgreSQL) FindLogins(ctx context.Context, user model.User) (logins []model.Login, err error) {
	if repo.db == nil {
		return nil, ErrDBNoDBConn
	}

	rows, err := repo.db.QueryContext(ctx,
		"SELECT login_id, username, password, metainfo FROM logins WHERE user_id = $1 ORDER BY login_id DESC",
		user.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	logins = make([]model.Login, 0)

	for rows.Next() {
		login := model.Login{}

		err = rows.Scan(&login.ID, &login.Username, &login.Password, &login.MetaInfo)
		if err != nil {
			return nil, err
		}

		logins = append(logins, login)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return logins, err
}

// SaveCard используется при сохранении данных банковской карты пользователя
func (repo RepoPostgreSQL) SaveCard(ctx context.Context, card *model.Card) (err error) {
	if repo.db == nil {
		return ErrDBNoDBConn
	}

	if card.ID == 0 {
		_, err = repo.db.ExecContext(ctx,
			`INSERT INTO cards (user_id, bank, number, date, cvv, metainfo) VALUES($1, $2, $3, $4, $5, $6)`,
			card.UserID,
			card.Bank,
			card.Number,
			card.Date,
			card.CVV,
			card.MetaInfo)
	} else {
		_, err = repo.db.ExecContext(ctx,
			`UPDATE cards SET bank = $1, number = $2, date = $3, cvv = $4, metainfo = $5 WHERE card_id = $6 and user_id = $7`,
			card.Bank,
			card.Number,
			card.Date,
			card.CVV,
			card.MetaInfo,
			card.ID,
			card.UserID)
	}

	return err
}

// DeleteCard используется при удалении данных банковской карты пользователя
func (repo RepoPostgreSQL) DeleteCard(ctx context.Context, card *model.Card) (err error) {
	if repo.db == nil {
		return ErrDBNoDBConn
	}

	_, err = repo.db.ExecContext(ctx,
		`DELETE FROM cards WHERE card_id = $1 and user_id = $2`,
		card.ID,
		card.UserID)

	return err
}

// FindCards возвращает данные банковских карт пользователя
func (repo RepoPostgreSQL) FindCards(ctx context.Context, user model.User) (cards []model.Card, err error) {
	if repo.db == nil {
		return nil, ErrDBNoDBConn
	}

	rows, err := repo.db.QueryContext(ctx,
		"SELECT card_id, bank, number, date, cvv, metainfo FROM cards WHERE user_id = $1 ORDER BY card_id DESC",
		user.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cards = make([]model.Card, 0)

	for rows.Next() {
		card := model.Card{}

		err = rows.Scan(&card.ID, &card.Bank, &card.Number, &card.Date, &card.CVV, &card.MetaInfo)
		if err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cards, err
}

// SaveFile используется при сохранении данных файла пользователя
func (repo RepoPostgreSQL) SaveFile(ctx context.Context, file *model.File) (err error) {
	if repo.db == nil {
		return ErrDBNoDBConn
	}

	if file.ID == 0 {
		_, err = repo.db.ExecContext(ctx,
			`INSERT INTO files (user_id, path, name, metainfo) VALUES($1, $2, $3, $4)`,
			file.UserID,
			file.Path,
			file.Name,
			file.MetaInfo)
	} else {
		_, err = repo.db.ExecContext(ctx,
			`UPDATE files SET path = $1, name = $2, metainfo = $3 WHERE file_id = $4 and user_id = $5`,
			file.Path,
			file.Name,
			file.MetaInfo,
			file.ID,
			file.UserID)
	}

	return err
}

// FindFile получает данные файла по Id пользователя и Id файла
func (repo RepoPostgreSQL) FindFile(ctx context.Context, fileID, userID int) (file *model.File, err error) {
	if repo.db == nil {
		return nil, ErrDBNoDBConn
	}

	file = &model.File{}

	err = repo.db.QueryRowContext(ctx,
		"SELECT file_id, user_id, path, name, metainfo FROM files WHERE file_id = $1 and user_id = $2",
		fileID, userID).Scan(&file.ID, &file.UserID, &file.Path, &file.Name, &file.MetaInfo)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return file, nil
}

// DeleteFile используется при удалении данных файла пользователя
func (repo RepoPostgreSQL) DeleteFile(ctx context.Context, file *model.File) (err error) {
	if repo.db == nil {
		return ErrDBNoDBConn
	}

	_, err = repo.db.ExecContext(ctx,
		`DELETE FROM files WHERE file_id = $1 and user_id = $2`,
		file.ID,
		file.UserID)

	return err
}

// FindFiles возвращает данные файлов пользователя
func (repo RepoPostgreSQL) FindFiles(ctx context.Context, user model.User) (files []model.File, err error) {
	if repo.db == nil {
		return nil, ErrDBNoDBConn
	}

	rows, err := repo.db.QueryContext(ctx,
		"SELECT file_id, name, metainfo FROM files WHERE user_id = $1 ORDER BY file_id DESC",
		user.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files = make([]model.File, 0)

	for rows.Next() {
		file := model.File{}

		err = rows.Scan(&file.ID, &file.Name, &file.MetaInfo)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return files, err
}

// Close закрывает соединение с БД
func (repo RepoPostgreSQL) Close() error {
	if repo.db == nil {
		return ErrDBNoDBConn
	}

	return repo.db.Close()
}
