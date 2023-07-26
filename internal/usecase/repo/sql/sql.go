package sql

import (
	"context"
	"fmt"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/internal/entity"
	"github.com/SETTER2000/prove/pkg/er"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

const (
	driverName = "pgx"
)

type (
	producerSQL struct {
		db *sqlx.DB
	}

	consumerSQL struct {
		db *sqlx.DB
	}

	InSQL struct {
		cfg *config.Config
		r   *consumerSQL
		w   *producerSQL
	}
)

// New слой взаимодействия с db в данном случаи с postgresql
func New(db *sqlx.DB) *InSQL {
	return &InSQL{
		cfg: config.GetConfig(),
		// создаём новый потребитель
		r: NewSQLConsumer(db),
		// создаём новый производитель
		w: NewSQLProducer(db),
	}
}

// NewSQLProducer производитель
func NewSQLProducer(db *sqlx.DB) *producerSQL {
	return &producerSQL{
		db: db,
	}
}

func (i *InSQL) Registry(ctx context.Context, a *entity.Authentication) error {
	if err := a.Validate(); err != nil {
		return err
	}
	if err := a.BeforeCreate(); err != nil {
		return err
	}
	err := i.w.db.QueryRow(
		"INSERT INTO public.user(login, encrypted_passwd) VALUES ($1,$2) RETURNING user_id",
		a.Login,
		a.EncryptPassword,
	).Scan(&a.ID)
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return er.NewConflictError("", "", er.ErrAlreadyExists)
		}
	}
	return nil
}

func (i *InSQL) SaveText(ctx context.Context, c *entity.Text) error {
	tx, err := i.w.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO public.text (text_data, user_id, meta_data, uploaded_at) VALUES ($1,$2,$3,now())")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, c.Text, c.UserID, c.MetaData)
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return er.NewConflictError("old url", "http://testiki", er.ErrAlreadyExists)
		}
		return err
	}

	return tx.Commit()

}

func (i *InSQL) SavePass(ctx context.Context, c *entity.Pass) error {
	tx, err := i.w.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO public.pass (login , password, user_id, meta_data, uploaded_at) VALUES ($1,$2,$3,$4,now())")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, c.Login, c.Password, c.UserID, c.MetaData)
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return er.NewConflictError("old url", "http://testiki", er.ErrAlreadyExists)
		}
		return err
	}

	return tx.Commit()

}

func (i *InSQL) SaveCard(ctx context.Context, c *entity.Card) error {
	tx, err := i.w.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO public.card (number, user_id, meta_data, uploaded_at) VALUES ($1,$2,$3,now())")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, c.Number, c.UserID, c.MetaData)
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return er.NewConflictError("old url", "http://testiki", er.ErrAlreadyExists)
		}
		return err
	}

	return tx.Commit()

}

func (i *InSQL) Post(ctx context.Context, sh *entity.Indra) error {
	stmt, err := i.w.db.Prepare("INSERT INTO public.prove (slug, url, user_id) VALUES ($1,$2,$3)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(sh.Slug, sh.URL, sh.UserID)
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return er.NewConflictError("old url", "http://testiki", er.ErrAlreadyExists)
		}
	}
	return nil
}

func (i *InSQL) Put(ctx context.Context, sh *entity.Indra) error {
	return i.Post(ctx, sh)
}

// NewSQLConsumer потребитель
func NewSQLConsumer(db *sqlx.DB) *consumerSQL {
	return &consumerSQL{
		db: db,
	}
}

func (i *InSQL) Get(ctx context.Context, sh *entity.Indra) (*entity.Indra, error) {
	var slug, url, id string
	var del bool
	rows, err := i.w.db.Query("SELECT slug, url, user_id, del FROM prove WHERE slug = $1 OR url = $2 ", sh.Slug, sh.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&slug, &url, &id, &del)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	sh.Slug = entity.Slug(slug)
	sh.URL = entity.URL(url)
	sh.UserID = entity.UserID(id)
	sh.Del = del

	return sh, nil
}

func (i *InSQL) GetByLogin(ctx context.Context, l string) (*entity.Authentication, error) {
	var a entity.Authentication
	var userID, login, encrypt string

	q := `SELECT user_id, login, encrypted_passwd FROM "user" WHERE login=$1`
	rows, err := i.w.db.Queryx(q, l)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&userID, &login, &encrypt)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	a.ID = userID
	a.Login = login
	a.EncryptPassword = encrypt

	return &a, nil
}

func (i *InSQL) GetByID(ctx context.Context, l string) (*entity.Authentication, error) {
	var a entity.Authentication
	var userID, login, encrypt string

	q := `SELECT user_id, login, encrypted_passwd FROM "user" WHERE user_id=$1`
	rows, err := i.w.db.Queryx(q, l)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&userID, &login, &encrypt)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	if userID == "" {
		return nil, er.ErrNotFound
	}

	a.ID = userID
	a.Login = login
	a.EncryptPassword = encrypt

	return &a, nil
}

// CardListGetUserID  вернуть список карт по UserID
func (i *InSQL) CardListGetUserID(ctx context.Context, u *entity.User) (*entity.CardList, error) {
	var number, userID, uploadedAt, metaData string
	//var accrual float32

	// 2020-12-10T15:15:45+03:00
	q := `SELECT number, user_id,  uploaded_at, meta_data FROM "card" WHERE user_id=$1 ORDER BY uploaded_at`

	rows, err := i.w.db.Queryx(q, u.UserID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	or := entity.CardResponse{}
	ol := entity.CardList{}
	for rows.Next() {
		err = rows.Scan(&number, &userID, &uploadedAt, &metaData)
		if err != nil {
			return nil, err
		}

		or.Number = number
		or.MetaData = metaData
		//or.Status = status
		//or.Accrual = accrual
		or.UploadedAt = uploadedAt
		ol = append(ol, or)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &ol, nil
}

func (i *InSQL) GetAll(context.Context, *entity.User) (*entity.User, error) {
	return nil, nil
}
func (i *InSQL) GetAllUrls() (entity.CountURLs, error) {

	return 0, nil
}
func (i *InSQL) GetAllUsers() (entity.CountUsers, error) {
	return 0, nil
}
func (i *InSQL) Delete(context.Context, *entity.User) error {
	return nil
}
func NewDB() (*sqlx.DB, error) {
	db, _ := sqlx.Open(driverName, config.GetConfig().ConnectDB)
	err := db.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	n := 100
	db.SetMaxIdleConns(n)
	db.SetMaxOpenConns(n)
	schema := `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS public.user
(
	user_id UUID NOT NULL DEFAULT uuid_generate_v1(),
  CONSTRAINT user_id_user PRIMARY KEY (user_id),
    login VARCHAR(100) NOT NULL UNIQUE,
    encrypted_passwd VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS public.card
(
    number      NUMERIC PRIMARY KEY,
    user_id     uuid,
    uploaded_at TIMESTAMP(0) WITH TIME ZONE,
    meta_data   VARCHAR,
    FOREIGN KEY (user_id) REFERENCES public."user" (user_id)
        MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE IF NOT EXISTS public.pass
(
    login       VARCHAR(100) NOT NULL,
    password    VARCHAR(100) NOT NULL,
    user_id     uuid,
    uploaded_at TIMESTAMP(0) WITH TIME ZONE,
    meta_data   VARCHAR,
    CONSTRAINT key_primary_unique
        PRIMARY KEY (login, password, user_id),
    FOREIGN KEY (user_id) REFERENCES public."user" (user_id)
        MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION

);
CREATE TABLE IF NOT EXISTS public.text
(
    id       SERIAL PRIMARY KEY,
    text_data    VARCHAR(100) NOT NULL,
    user_id     uuid NOT NULL,
    uploaded_at TIMESTAMP(0) WITH TIME ZONE,
    meta_data   VARCHAR,
    FOREIGN KEY (user_id) REFERENCES public."user" (user_id)
        MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION

);
`
	db.MustExec(schema)
	if err != nil {
		panic(err)
	}
	return db, nil
}

func (i *InSQL) Read() error {
	return nil
}
func (i *InSQL) Save() error {
	return nil
}
