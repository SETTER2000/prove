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

// ProveUseCase -.
type ProveUseCase struct {
	repo ProveRepo
}

// New -.
func New(r ProveRepo) *ProveUseCase {
	return &ProveUseCase{
		repo: r,
	}
}

func (uc *ProveUseCase) UserFindByLogin(ctx context.Context, s string) (*entity.Authentication, error) {
	//auth.Login = ctx.Value(auth.Cookie.AccessTokenName).(string)
	a, err := uc.repo.GetByLogin(ctx, s)
	if err != nil {
		return nil, err
	}
	return a, nil
}
func (uc *ProveUseCase) UserFindByID(ctx context.Context, s string) (*entity.Authentication, error) {
	//auth.Login = ctx.Value(auth.Cookie.AccessTokenName).(string)
	a, err := uc.repo.GetByID(ctx, s)
	if err != nil {
		return nil, err
	}
	return a, nil
}
func (uc *ProveUseCase) SavePass(ctx context.Context, c *entity.Pass) error {
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

func (uc *ProveUseCase) SaveText(ctx context.Context, c *entity.Text) error {
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
func (uc *ProveUseCase) SaveCard(ctx context.Context, c *entity.Card) error {
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
func (uc *ProveUseCase) Register(ctx context.Context, auth *entity.Authentication) error {
	//auth.Login = ctx.Value(auth.Cookie.AccessTokenName).(string)
	err := uc.repo.Registry(ctx, auth)
	if err != nil {
		return err
	}
	return nil
}

// ShortLink принимает короткий URL и возвращает длинный (GET /api/{key})
func (uc *ProveUseCase) ShortLink(ctx context.Context, ind *entity.Prove) (*entity.Prove, error) {
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
func (uc *ProveUseCase) CardListUserID(ctx context.Context, u *entity.User) (*entity.CardList, error) {
	ol, err := uc.repo.CardListGetUserID(ctx, u)
	if err == nil {
		return ol, nil
	}
	return nil, ErrBadRequest
}

func (uc *ProveUseCase) SaveService() error {
	err := uc.repo.Save()
	if err != nil {
		return err
	}
	return nil
}
func (uc *ProveUseCase) ReadService() error {
	err := uc.repo.Read()
	if err != nil {
		return err
	}
	return nil
}