// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"github.com/SETTER2000/prove/internal/entity"
	"github.com/SETTER2000/prove/internal/usecase/repo/file"
	"github.com/SETTER2000/prove/internal/usecase/repo/memory"
	"github.com/SETTER2000/prove/internal/usecase/repo/sql"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type (
	// Prove -.
	// Post -.
	Prove interface {
		SaveText(context.Context, *entity.Text) error
		SavePass(context.Context, *entity.Pass) error
		SaveCard(context.Context, *entity.Card) error
		CardListUserID(ctx context.Context, u *entity.User) (*entity.CardList, error)
		Register(context.Context, *entity.Authentication) error
		UserFindByLogin(context.Context, string) (*entity.Authentication, error)
		UserFindByID(context.Context, string) (*entity.Authentication, error)
		ShortLink(context.Context, *entity.Prove) (*entity.Prove, error)
		ReadService() error
		SaveService() error
	}

	// ProveRepo -.
	ProveRepo interface {
		SaveText(context.Context, *entity.Text) error
		SavePass(context.Context, *entity.Pass) error
		SaveCard(context.Context, *entity.Card) error
		CardListGetUserID(context.Context, *entity.User) (*entity.CardList, error)
		Registry(context.Context, *entity.Authentication) error
		GetByLogin(context.Context, string) (*entity.Authentication, error)
		GetByID(context.Context, string) (*entity.Authentication, error)
		// Post от shorturl
		Post(context.Context, *entity.Prove) error
		Put(context.Context, *entity.Prove) error
		Get(context.Context, *entity.Prove) (*entity.Prove, error)
		GetAll(context.Context, *entity.User) (*entity.User, error)
		GetAllUrls() (entity.CountURLs, error)
		GetAllUsers() (entity.CountUsers, error)
		Delete(context.Context, *entity.User) error
		Read() error
		Save() error
	}
)

// Строка ниже не несёт функциональной нагрузки
// её можно убрать без последствий для работы программы
// это отладочная строка, в этой строке приведением типов проверяем,
// реализует ли структура *batchPostProvider интерфейс BatchPostProvider —
// если нет или если методы прописаны неверно,
// то компилятор выдаст на этой строке ошибку типизации
var _ ProveRepo = (*memory.Memory)(nil)
var _ ProveRepo = (*file.InFiles)(nil)
var _ ProveRepo = (*sql.InSQL)(nil)
