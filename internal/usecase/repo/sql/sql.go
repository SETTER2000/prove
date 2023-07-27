package sql

import (
	"context"
	"fmt"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/internal/entity"
	"github.com/SETTER2000/prove/pkg/er"
	"github.com/google/uuid"
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
		"INSERT INTO public.user(login, encrypted_passwd,email, surname, name,patronymic, group_id, uploaded_at) VALUES ($1,$2,$3,$4,$5,$6,$7,now()) RETURNING user_id",
		a.Login,
		a.EncryptPassword,
		a.Email,
		a.Surname,
		a.Name,
		a.Patronymic,
		a.GroupID,
	).Scan(&a.ID)
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return er.NewConflictError("", "", er.ErrAlreadyExists)
		}

		if err != nil {
			log.Printf("error registration: %s", err.Message)
			return er.ErrBadRequest
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

func (i *InSQL) SaveSolution(ctx context.Context, c *entity.Solution) error {
	uid, err := uuid.Parse(c.TaskID)
	if err != nil {
		log.Printf("error uuid %s: %s", c.TaskID, err.Error())
	}
	err = i.w.db.QueryRow(
		"INSERT INTO public.solution(description, task_id, encrypted_solution, uploaded_at) VALUES ($1,$2,$3,now()) RETURNING solution_id",
		c.Description,
		uid,
		c.Solution,
	).Scan(&c.SolutionID)
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return er.NewConflictError("", "", er.ErrAlreadyExists)
		}
		if err != nil {
			log.Printf("error save solution: %s", err.Message)
			return er.ErrBadRequest
		}
	}

	return nil
}

func (i *InSQL) SaveGroup(ctx context.Context, c *entity.Group) error {
	err := i.w.db.QueryRow(
		"INSERT INTO public.group(group_name, uploaded_at) VALUES ($1,now()) RETURNING group_name",
		c.GroupName,
	).Scan(&c.GroupName)
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return er.NewConflictError("", "", er.ErrAlreadyExists)
		}
	}

	return nil
}

func (i *InSQL) SaveTask(ctx context.Context, c *entity.Task) error {
	err := i.w.db.QueryRow(
		"INSERT INTO public.task(task_name, description,price, uploaded_at) VALUES ($1,$2,$3, now()) RETURNING task_name",
		c.TaskName,
		c.Description,
		c.Price,
	).Scan(&c.TaskName)
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return er.NewConflictError("", "", er.ErrAlreadyExists)
		}
	}

	return nil
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

func (i *InSQL) Post(ctx context.Context, sh *entity.Prove) error {
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

func (i *InSQL) Put(ctx context.Context, sh *entity.Prove) error {
	return i.Post(ctx, sh)
}

// NewSQLConsumer потребитель
func NewSQLConsumer(db *sqlx.DB) *consumerSQL {
	return &consumerSQL{
		db: db,
	}
}

func (i *InSQL) Get(ctx context.Context, sh *entity.Prove) (*entity.Prove, error) {
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

// GroupList вернёт список групп
func (i *InSQL) GroupList(ctx context.Context) (*entity.GroupList, error) {
	var grupID, groupName, uploadedAt string

	q := `SELECT group_id, group_name, uploaded_at FROM "group" ORDER BY group_name`
	rows, err := i.w.db.Queryx(q)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	or := entity.GroupList{}
	ol := entity.Group{}
	for rows.Next() {
		err = rows.Scan(&grupID, &groupName, &uploadedAt)
		if err != nil {
			return nil, err
		}

		ol.GroupName = groupName
		ol.GroupID = grupID
		ol.UploadedAt = uploadedAt
		or = append(or, ol)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &or, nil
}

// TaskList вернёт список задач
func (i *InSQL) TaskList(ctx context.Context) (*entity.TaskList, error) {
	var taskID, taskName, description, price, uploadedAt string

	q := `SELECT task_id, task_name, description,price, uploaded_at FROM "task" ORDER BY task_name`
	rows, err := i.w.db.Queryx(q)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	or := entity.TaskList{}
	ol := entity.Task{}
	for rows.Next() {
		err = rows.Scan(&taskID, &taskName, &description, &price, &uploadedAt)
		if err != nil {
			return nil, err
		}

		ol.TaskName = taskName
		ol.TaskID = taskID
		ol.Description = description
		ol.Price = price
		ol.UploadedAt = uploadedAt
		or = append(or, ol)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &or, nil
}

// TaskKey вернуть ответ по задаче
func (i *InSQL) TaskKey(ctx context.Context, u *entity.User, t *entity.Task) (*entity.SolutionList, error) {
	var solution, taskID, solutionID string

	q := `SELECT encrypted_solution, task_id, solution_id FROM "solution" WHERE task_id=$1 ORDER BY LIMIT 1`
	rows, err := i.w.db.Queryx(q, t.TaskID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	sl := entity.SolutionList{}
	s := entity.Solution{}
	for rows.Next() {
		err = rows.Scan(&solution, &taskID, &solutionID)
		if err != nil {
			return nil, err
		}
		//
		s.Solution = solution
		s.TaskID = taskID
		s.SolutionID = solutionID
		sl = append(sl, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &sl, nil
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
CREATE TABLE IF NOT EXISTS public.group
(
    group_id    UUID NOT NULL DEFAULT uuid_generate_v1(),
    CONSTRAINT 	group_id_group PRIMARY KEY (group_id),
    group_name  VARCHAR(100) NOT NULL UNIQUE,
    uploaded_at TIMESTAMP(0) WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS public.user
(
    user_id     UUID NOT NULL DEFAULT uuid_generate_v1(),
    CONSTRAINT 	user_id_user PRIMARY KEY (user_id),
    login       	 VARCHAR(100) NOT NULL UNIQUE,
    encrypted_passwd VARCHAR(100) NOT NULL,
    email 			 VARCHAR (100) NOT NULL,
    surname			 VARCHAR(100) NOT NULL,
    name 			 VARCHAR(100) NOT NULL,
    patronymic 		 VARCHAR(100) NOT NULL,
    group_id uuid,
    foreign key (group_id) references public."group" (group_id)
        match simple on update no action on delete no action,
    uploaded_at      TIMESTAMP(0) WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS public.task
(
    task_id     UUID   NOT NULL DEFAULT uuid_generate_v1(),
    CONSTRAINT task_id_task PRIMARY KEY (task_id),
    task_name   		VARCHAR(100) NOT NULL UNIQUE,
    description TEXT 	NOT NULL,
    price        		NUMERIC NOT NULL DEFAULT 0,
    uploaded_at 		TIMESTAMP(0) WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS public.solution
(
    solution_id     UUID NOT NULL DEFAULT uuid_generate_v1(),
    CONSTRAINT 	solution_id_solution PRIMARY KEY (solution_id),
    encrypted_solution VARCHAR(100) NOT NULL,
    description TEXT 	NOT NULL,
    task_id uuid,
    foreign key (task_id) references public."task" (task_id)
        match simple on update no action on delete no action,
    uploaded_at      TIMESTAMP(0) WITH TIME ZONE
);


CREATE TABLE IF NOT EXISTS public.cash
(
    cash_id     UUID NOT NULL DEFAULT uuid_generate_v1(),
    CONSTRAINT 	cash_id_cash PRIMARY KEY (cash_id),
    price       MONEY 	NOT NULL,
    uploaded_at 		TIMESTAMP(0) WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS cash_task (
  	cash_id     UUID REFERENCES cash (cash_id) ON UPDATE CASCADE ON DELETE CASCADE,
	task_id 	UUID REFERENCES task (task_id) ON UPDATE CASCADE,
	amount      numeric NOT NULL DEFAULT 1,
	CONSTRAINT cash_task_pkey PRIMARY KEY (cash_id, task_id)
);


CREATE TABLE IF NOT EXISTS public.ball
(
    ball_id         UUID NOT NULL DEFAULT uuid_generate_v1(),
    CONSTRAINT 	ball_id_ball PRIMARY KEY (ball_id),
    price           MONEY   NOT NULL,
    uploaded_at 	TIMESTAMP(0) WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS task_user (
  	user_id     UUID REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE,
	task_id 	UUID REFERENCES task (task_id) ON UPDATE CASCADE,
	amount      numeric NOT NULL DEFAULT 1,
	CONSTRAINT user_task_pkey PRIMARY KEY (user_id, task_id)
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
