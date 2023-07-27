package sql

import (
	"context"
	"database/sql"
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
	"strconv"
)

const (
	driverName = "pgx"
	maxCredit  = 1000.00
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

	// Balance initialization
	var credit string
	err = i.w.db.QueryRow(
		"INSERT INTO public.balance(user_id, uploaded_at) VALUES ($1,now()) RETURNING credit",
		a.ID,
	).Scan(&credit)

	if err != nil {
		log.Printf("error registration: %s", err)
		return er.ErrBadRequest
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

// GetAdmin вернёт nil, если пользователь является админом.
// Админом является тот кто первый зарегистрировался.
func (i *InSQL) GetAdmin(u *entity.User) bool {
	var id string
	rows, err := i.r.db.Query("SELECT user_id FROM public.\"user\" ORDER BY uploaded_at LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	if u.UserID == id {
		return true
	}
	return false
}

// GetSolution обновление баланса пользователя по его id
func (i *InSQL) GetSolution(ctx context.Context, s *entity.SolutionData) error {
	tx, err := i.w.db.Begin()
	if err != nil {
		log.Printf("%e", err)
	}
	ct := context.Background()
	defer tx.Rollback()

	credit, err := i.GetBalance(ctx, s)
	if err != nil {
		log.Printf("error GetBalance(ctx, s): %s", err.Error())
		return er.ErrBadRequest
	}
	u := entity.User{
		UserID: string(s.UserID),
	}
	t := entity.Task{
		TaskID: s.TaskID,
	}
	task, err := i.TaskKey(ctx, &u, &t)
	if err != nil {
		log.Printf("error, TaskKey(ctx, &u, &t): %s", err.Error())
		return er.ErrBadRequest
	}
	f, err := strconv.ParseFloat(task.Price, 64)
	if err != nil {
		log.Printf("error Atoi %s\n", err)
	}

	stmt, err := tx.PrepareContext(ct, "UPDATE balance SET credit=$1 WHERE user_id=$2")
	if err != nil {
		return err
	}
	defer stmt.Close()
	crd := credit + f
	if crd > maxCredit {
		return er.ErrExhaustedCredit
	}
	if _, err = stmt.Exec(&crd, s.UserID); err != nil {
		return err
	}

	return tx.Commit()
}

// GetBalance получить информацию о кредитоспособности,
// если вернёт err!=nil - кредит исчерпан
func (i *InSQL) GetBalance(ctx context.Context, u *entity.SolutionData) (float64, error) {
	var credit float64

	rows, err := i.w.db.Query("SELECT credit FROM balance WHERE user_id = $1 ", u.UserID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&credit)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	if credit <= maxCredit && credit > 0 {
		if credit == 0.01 {
			credit = 0
		}
		return credit, nil
	}

	return 0, fmt.Errorf("exceeded credit")
}

// Balance получение текущего баланса пользователя
func (i *InSQL) Balance(ctx context.Context) (*entity.Balance, error) {
	var credit sql.NullFloat64

	q := `SELECT credit FROM public."balance" WHERE user_id=$1`

	rows, err := i.w.db.Queryx(q, ctx.Value(i.cfg.Cookie.AccessTokenName).(string))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	b := entity.Balance{}
	for rows.Next() {
		err = rows.Scan(&credit)
		if err != nil {
			return nil, err
		}
		b.Current = float32(credit.Float64)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &b, nil
}

// BalanceAdd пополнить баланс
func (i *InSQL) BalanceAdd(ctx context.Context, s *entity.Balance) error {
	tx, err := i.w.db.Begin()
	if err != nil {
		log.Printf("%e", err)
	}
	ct := context.Background()
	defer tx.Rollback()
	sl := &entity.SolutionData{
		UserID: entity.UserID(s.UserID),
	}
	credit, err := i.GetBalance(ctx, sl)
	if err != nil {
		log.Printf("error GetBalance(ctx, s): %s", err.Error())
		return er.ErrBadRequest
	}

	stmt, err := tx.PrepareContext(ct, "UPDATE balance SET credit=$1 WHERE user_id=$2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	crd := credit - float64(s.Current)
	if crd > 0.01 {
		if _, err = stmt.Exec(&crd, s.UserID); err != nil {
			return err
		}
	}
	return tx.Commit()
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

// TaskKey вернуть одну задачу
func (i *InSQL) TaskKey(ctx context.Context, u *entity.User, t *entity.Task) (*entity.Task, error) {
	var taskID, taskName, description, price, uploadedAt string
	q := `SELECT task_id, task_name,description,  price, uploaded_at FROM task WHERE task_id=$1`
	rows, err := i.w.db.Queryx(q, t.TaskID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	s := &entity.Task{}
	for rows.Next() {
		err = rows.Scan(&taskID, &taskName, &description, &price, &uploadedAt)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	s.TaskID = t.TaskID
	s.TaskName = taskName
	s.Description = description
	s.Price = price
	s.UploadedAt = uploadedAt

	return s, nil
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
    price        		NUMERIC(9,2) NOT NULL DEFAULT 0,
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

CREATE TABLE IF NOT EXISTS public.balance
(
    balance_id      UUID NOT NULL DEFAULT uuid_generate_v1(),
    CONSTRAINT 	    balance_id_balance PRIMARY KEY (balance_id),
    credit           NUMERIC(9,2)   NOT NULL DEFAULT 0.01,
    uploaded_at 	TIMESTAMP(0) WITH TIME ZONE,
    user_id         uuid,
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
