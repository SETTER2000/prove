package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/internal/entity"
	"github.com/SETTER2000/prove/internal/usecase"
	"github.com/SETTER2000/prove/pkg/er"
	"github.com/SETTER2000/prove/pkg/log/logger"
	"github.com/SETTER2000/prove/pkg/middleware/encrypt"
	"github.com/SETTER2000/prove/pkg/middleware/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"io"
	"net/http"
	"os"
	"time"
)

type ServerHandler struct {
	s   usecase.Prove
	l   logger.Logger
	cfg *config.Config
}

func NewServerHandler(s usecase.Prove) *ServerHandler {
	return &ServerHandler{
		s:   s,
		l:   logger.GetLogger(),
		cfg: config.GetConfig(),
	}
}
func (sh *ServerHandler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	headerTypes := []string{
		"application/javascript",
		"application/x-gzip",
		"application/gzip",
		"application/json",
		"text/css",
		"text/html",
		"text/plain",
		"text/xml",
	}
	// AllowContentType применяет белый список запросов Content-Types,
	// в противном случае отвечает статусом 415 Unsupported Media Type.
	router.Use(middleware.AllowContentType(headerTypes...))
	router.Use(middleware.Compress(5, headerTypes...))
	router.Use(middleware.AllowContentEncoding("deflate", "gzip"))
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypePlainText))
	router.Use(encrypt.EncryptionKeyCookie)
	router.Use(gzip.DeCompressGzip)
	router.Mount("/debug", middleware.Profiler())

	// Swagger
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	h := router.Route("/api", func(r chi.Router) { r.Routes() })
	{
		h.Post("/group", sh.hGroup)
		h.Get("/group", sh.hGroups)
		h.Post("/task", sh.hTask)
		h.Get("/task", sh.hTasks)

		h.Post("/card", sh.hCard)
		h.Get("/card", sh.hCards)
		h.Post("/pass", sh.hPass)
		h.Post("/text", sh.hText)

		h.Route("/user", func(r chi.Router) {
			r.Post("/register", sh.handleUserCreate)
			r.Post("/login", sh.handleUserLogin)
		})
	}

	root := router.Route("/", func(r chi.Router) { r.Routes() })
	{
		root.Get("/ping", sh.connect)
	}

	return router
}

// @Summary     Return JSON empty
// @Description Save data task
// @ID          Создание задачи.
// @Tags  	    prove
// @Accept      json
// @Produce     json
// @Success     200 {object} response — данные успешно сохранены
// @Failure     400 {object} response — неверный формат запроса
// @Failure     409 {object} response — имя уже занято
// @Failure     500 {object} response — внутренняя ошибка сервера
// @Router      /task [post]
func (sh *ServerHandler) hTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c := &entity.Task{}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &c); err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	err = sh.s.SaveTask(ctx, c)
	if err != nil {
		if errors.Is(err, er.ErrAlreadyExists) {
			sh.error(w, r, http.StatusConflict, err)
			return
		}
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	sh.respond(w, r, http.StatusOK, c)
}

// @Summary     Return JSON empty
// @Description Save data group
// @ID          Создание группы.
// @Tags  	    prove
// @Accept      json
// @Produce     json
// @Success     200 {object} response — данные успешно сохранены
// @Failure     400 {object} response — неверный формат запроса
// @Failure     409 {object} response — имя уже занято
// @Failure     500 {object} response — внутренняя ошибка сервера
// @Router      /group [post]
func (sh *ServerHandler) hGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c := &entity.Group{}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &c); err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	err = sh.s.SaveGroup(ctx, c)
	if err != nil {
		if errors.Is(err, er.ErrAlreadyExists) {
			sh.error(w, r, http.StatusConflict, err)
			return
		}
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	sh.respond(w, r, http.StatusOK, c)
}

// @Summary     Return JSON empty
// @Description Save data card
// @ID          Сохранение данных банковских карт.
// @Tags  	    prove
// @Accept      json
// @Produce     json
// @Success     200 {object} response — данные карты успешно сохранены
// @Failure     400 {object} response — неверный формат запроса
// @Failure     409 {object} response — логин уже занят ??
// @Failure     500 {object} response — внутренняя ошибка сервера
// @Router      /card [post]
func (sh *ServerHandler) hCard(w http.ResponseWriter, r *http.Request) {
	userID, err := sh.IsAuthenticated(w, r)
	if err != nil {
		sh.respond(w, r, http.StatusUnauthorized, nil)
		return
	}
	// Для любых данных должна быть возможность хранения
	// произвольной текстовой метаинформации (принадлежность
	// данных к веб-сайту, личности или банку, списки одноразовых
	// кодов активации и прочее).
	ctx := r.Context()
	c := &entity.Card{}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &c); err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	c.UserID = userID
	err = sh.s.SaveCard(ctx, c)
	if err != nil {
		if errors.Is(err, er.ErrAlreadyExists) {
			sh.error(w, r, http.StatusConflict, err)
			return
		}
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	//sh.SessionCreated(w, r, a.ID)
	c.Sanitize()
	sh.respond(w, r, http.StatusOK, c)
}

// @Summary     Return JSON empty
// @Description Save data text
// @ID          Сохраняет произвольные текстовые данные.
// @Tags  	    prove
// @Accept      json
// @Produce     json
// @Success     200 {object} response — данные карты успешно сохранены
// @Failure     400 {object} response — неверный формат запроса
// @Failure     409 {object} response — повтор
// @Failure     500 {object} response — внутренняя ошибка сервера
// @Router      /pass [post]
func (sh *ServerHandler) hText(w http.ResponseWriter, r *http.Request) {
	userID, err := sh.IsAuthenticated(w, r)
	if err != nil {
		sh.respond(w, r, http.StatusUnauthorized, nil)
		return
	}
	// Для любых данных должна быть возможность хранения произвольной
	// текстовой метаинформации (принадлежность данных к веб-сайту,
	// личности или банку, списки одноразовых кодов активации и прочее).
	ctx := r.Context()
	c := &entity.Text{}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &c); err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	c.UserID = userID
	err = sh.s.SaveText(ctx, c)
	if err != nil {
		if errors.Is(err, er.ErrAlreadyExists) {
			sh.error(w, r, http.StatusConflict, err)
			return
		}
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	c.Sanitize()
	sh.respond(w, r, http.StatusOK, c)
}

// @Summary     Return JSON empty
// @Description Save data login/password
// @ID          Сохранение пары логин/пароль.
// @Tags  	    prove
// @Accept      json
// @Produce     json
// @Success     200 {object} response — данные карты успешно сохранены
// @Failure     400 {object} response — неверный формат запроса
// @Failure     409 {object} response — повтор
// @Failure     500 {object} response — внутренняя ошибка сервера
// @Router      /pass [post]
func (sh *ServerHandler) hPass(w http.ResponseWriter, r *http.Request) {
	userID, err := sh.IsAuthenticated(w, r)
	if err != nil {
		sh.respond(w, r, http.StatusUnauthorized, nil)
		return
	}
	// Для любых данных должна быть возможность хранения произвольной
	// текстовой метаинформации (принадлежность данных к веб-сайту,
	// личности или банку, списки одноразовых кодов активации и прочее).
	ctx := r.Context()
	c := &entity.Pass{}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &c); err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	c.UserID = userID
	err = sh.s.SavePass(ctx, c)
	if err != nil {
		if errors.Is(err, er.ErrAlreadyExists) {
			sh.error(w, r, http.StatusConflict, err)
			return
		}
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	c.Sanitize()
	sh.respond(w, r, http.StatusOK, c)
}

// GET /ping, который при запросе проверяет соединение с базой данных
// при успешной проверке хендлер должен вернуть HTTP-статус 200 OK
// при неуспешной — 500 Internal Server Error
func (sh *ServerHandler) connect(w http.ResponseWriter, r *http.Request) {
	dsn, ok := os.LookupEnv("DATABASE_URI")
	if !ok || dsn == "" {
		dsn = sh.cfg.ConnectDB
		if dsn == "" {
			sh.l.Info("connect DSN string is empty: %v\n", dsn)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		db, err := pgx.Connect(r.Context(), os.Getenv("DATABASE_URI"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer db.Close(context.Background())
		sh.respond(w, r, http.StatusOK, "connect... ")
	}
}

// @Summary     Return JSON empty
// @Description Redirect to log URL
// @ID          Регистрация пользователя
// @Tags  	    prove
// @Accept      json
// @Produce     json
// @Success     200 {object} response — пользователь успешно зарегистрирован и аутентифицирован
// @Failure     400 {object} response — неверный формат запроса
// @Failure     409 {object} response — логин уже занят
// @Failure     500 {object} response — внутренняя ошибка сервера
// @Router      /user/register [post]
func (sh *ServerHandler) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	// Регистрация производится по паре логин/пароль.
	// Каждый логин должен быть уникальным. После успешной регистрации
	// должна происходить автоматическая аутентификация пользователя.
	// Для передачи аутентификационных данных используйте
	// механизм cookies или HTTP-заголовок Authorization.
	ctx := r.Context()
	a := &entity.Authentication{Config: sh.cfg}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &a); err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	err = sh.s.Register(ctx, a)
	if err != nil {
		if errors.Is(err, er.ErrAlreadyExists) {
			sh.error(w, r, http.StatusConflict, err)
			return
		}
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	sh.SessionCreated(w, r, a.ID)
	a.Sanitize()
	sh.respond(w, r, http.StatusOK, a)
}

// @Summary     Return JSON empty
// @ID          Аутентификация пользователя
// @Tags  	    prove
// @Accept      json
// @Produce     json
// @Success     200 {object} response — пользователь успешно аутентифицирован
// @Failure     400 {object} response — неверный формат запроса
// @Failure     401 {object} response — неверная пара логин/пароль
// @Failure     500 {object} response — внутренняя ошибка сервера
// @Router      /user/login [post]
func (sh *ServerHandler) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	a := &entity.Authentication{Config: sh.cfg}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &a); err != nil {
		sh.error(w, r, http.StatusInternalServerError, err)
		return
	}
	if err := a.Validate(); err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	if err := a.BeforeCreate(); err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	u, err := sh.s.UserFindByLogin(r.Context(), a.Login) // will return the user by login
	if err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
		return
	}
	if !u.ComparePassword(a.Password) {
		sh.error(w, r, http.StatusUnauthorized, er.ErrIncorrectLoginOrPass)
		return
	}
	sh.SessionCreated(w, r, u.ID)
	sh.respond(w, r, http.StatusOK, nil)
}

// @Summary     Return JSON empty
// @Description Получить список
// @ID          hTasks
// @Tags  	    prove-server
// @Accept      text
// @Produce     text
// @Success     200 {object} response — успешная обработка запроса
// @Success     204 {object} response — нет данных для ответ
// @Failure     500 {object} response — внутренняя ошибка сервера
// @Router      /task [get]
func (sh *ServerHandler) hTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ol, err := sh.s.TaskList(ctx)
	if err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
	}

	if len(*ol) < 1 {
		sh.respond(w, r, http.StatusNoContent, "нет данных для ответа")
	}

	sh.respond(w, r, http.StatusOK, ol)
}

// @Summary     Return JSON empty
// @Description Получить список
// @ID          hGroups
// @Tags  	    prove-server
// @Accept      text
// @Produce     text
// @Success     200 {object} response — успешная обработка запроса
// @Success     204 {object} response — нет данных для ответ
// @Failure     500 {object} response — внутренняя ошибка сервера
// @Router      /group [get]
func (sh *ServerHandler) hGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ol, err := sh.s.GroupList(ctx)
	if err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
	}

	if len(*ol) < 1 {
		sh.respond(w, r, http.StatusNoContent, "нет данных для ответа")
	}

	sh.respond(w, r, http.StatusOK, ol)
}

// @Summary     Return JSON empty
// @Description Получение списка загруженных номеров карт
// @ID          hCards
// @Tags  	    prove-server
// @Accept      text
// @Produce     text
// @Success     200 {object} response — успешная обработка запроса
// @Success     204 {object} response — нет данных для ответа
// @Failure     401 {object} response — пользователь не аутентифицирован
// @Failure     500 {object} response — внутренняя ошибка сервера
// @Router      /card [get]
func (sh *ServerHandler) hCards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := entity.User{}

	userID, err := sh.IsAuthenticated(w, r)
	if err != nil {
		sh.respond(w, r, http.StatusUnauthorized, nil)
		return
	}

	u.UserID = userID
	ol, err := sh.s.CardListUserID(ctx, &u)
	if err != nil {
		sh.error(w, r, http.StatusBadRequest, err)
	}

	if len(*ol) < 1 {
		sh.respond(w, r, http.StatusNoContent, "нет данных для ответа")
	}

	sh.respond(w, r, http.StatusOK, ol)
}

// ContentTypeCheck проверка соответствует ли content-type запроса endpoint
func (sh *ServerHandler) ContentTypeCheck(w http.ResponseWriter, r *http.Request, t string) bool {
	return r.Header.Get("Content-Type") == t
}

func (sh *ServerHandler) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	sh.respond(w, r, code, map[string]string{"error": err.Error()})
}
func (sh *ServerHandler) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (sh *ServerHandler) SessionCreated(w http.ResponseWriter, r *http.Request, data string) error {
	en := encrypt.Encrypt{}
	// ...создать подписанный секретным ключом токен,
	token, err := en.EncryptToken(sh.cfg.SecretKey, data)
	if err != nil {
		fmt.Printf("Encrypt error: %v\n", err)
		return err
	}
	// ...установить куку с именем access_token,
	// а в качестве значения установить зашифрованный,
	// подписанный токен
	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(time.Minute * 60),
	})

	_, err = en.DecryptToken(token, sh.cfg.SecretKey)
	if err != nil {
		fmt.Printf(" Decrypt error: %v\n", err)
		return err
	}

	return nil
}

// IsAuthenticated проверка авторизирован ли пользователь, по сути проверка токена на пригодность
// TODO возможно здесь нужно реализовать проверку времени жизни токена
func (sh *ServerHandler) IsAuthenticated(w http.ResponseWriter, r *http.Request) (string, error) {
	var e encrypt.Encrypt
	ctx := r.Context()
	at, err := r.Cookie("access_token")
	if err == http.ErrNoCookie {
		return "", err
	}
	// если кука обнаружена, то расшифровываем токен,
	// содержащийся в ней, и проверяем подпись
	dt, err := e.DecryptToken(at.Value, sh.cfg.SecretKey)
	if err != nil {
		fmt.Printf("error decrypt cookie: %e", err)
		return "", err
	}

	_, err = sh.s.UserFindByID(r.Context(), dt)
	if err != nil {
		return "", err
	}
	return ctx.Value(sh.cfg.AccessTokenName).(string), nil
}
