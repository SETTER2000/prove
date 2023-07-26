package usecase

import (
	"context"
	"errors"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/internal/entity"
	"github.com/SETTER2000/prove/pkg/er"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrBadRequest    = errors.New("bad request")
)

// IndraUseCase -.
type IndraUseCase struct {
	repo IndraRepo
}

// New -.
func New(r IndraRepo) *IndraUseCase {
	return &IndraUseCase{
		repo: r,
	}
}

func (uc *IndraUseCase) UserFindByLogin(ctx context.Context, s string) (*entity.Authentication, error) {
	//auth.Login = ctx.Value(auth.Cookie.AccessTokenName).(string)
	a, err := uc.repo.GetByLogin(ctx, s)
	if err != nil {
		return nil, err
	}
	return a, nil
}
func (uc *IndraUseCase) UserFindByID(ctx context.Context, s string) (*entity.Authentication, error) {
	//auth.Login = ctx.Value(auth.Cookie.AccessTokenName).(string)
	a, err := uc.repo.GetByID(ctx, s)
	if err != nil {
		return nil, err
	}
	return a, nil
}
func (uc *IndraUseCase) SavePass(ctx context.Context, c *entity.Pass) error {
	err := c.Validate()
	if err != nil {
		return err
	}

	err = uc.repo.SavePass(ctx, c)
	if err != nil {
		return err
	}
	return nil
}

func (uc *IndraUseCase) SaveText(ctx context.Context, c *entity.Text) error {
	err := c.Validate()
	if err != nil {
		return err
	}

	err = uc.repo.SaveText(ctx, c)
	if err != nil {
		return err
	}
	return nil
}
func (uc *IndraUseCase) SaveCard(ctx context.Context, c *entity.Card) error {
	err := c.Validate()
	if err != nil {
		return err
	}
	err = c.Luna()
	if err != nil {
		return err
	}

	err = uc.repo.SaveCard(ctx, c)
	if err != nil {
		return err
	}
	return nil
}
func (uc *IndraUseCase) Register(ctx context.Context, auth *entity.Authentication) error {
	//auth.Login = ctx.Value(auth.Cookie.AccessTokenName).(string)
	err := uc.repo.Registry(ctx, auth)
	if err != nil {
		return err
	}
	return nil
}

// ShortLink принимает короткий URL и возвращает длинный (GET /api/{key})
func (uc *IndraUseCase) ShortLink(ctx context.Context, ind *entity.Indra) (*entity.Indra, error) {
	ind.Config = config.GetConfig()
	sh, err := uc.repo.Get(ctx, ind)
	if err != nil {
		return nil, er.ErrBadRequest
	}

	//При запросе удалённого URL с помощью хендлера GET /{id} нужно вернуть статус 410 Gone
	if sh.Del {
		return nil, er.ErrStatusGone
	}
	return sh, nil
}

// CardListUserID возвращает все сохранённые карты пользователя
func (uc *IndraUseCase) CardListUserID(ctx context.Context, u *entity.User) (*entity.CardList, error) {
	ol, err := uc.repo.CardListGetUserID(ctx, u)
	if err == nil {
		return ol, nil
	}
	return nil, ErrBadRequest
}

func (uc *IndraUseCase) SaveService() error {
	err := uc.repo.Save()
	if err != nil {
		return err
	}
	return nil
}
func (uc *IndraUseCase) ReadService() error {
	err := uc.repo.Read()
	if err != nil {
		return err
	}
	return nil
}
