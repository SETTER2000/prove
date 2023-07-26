package memory

import (
	"context"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/internal/entity"
	"github.com/SETTER2000/prove/pkg/er"
	"github.com/SETTER2000/prove/scripts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type Memory struct {
	cfg  *config.Config
	m    map[entity.UserID]entity.Proves // <-- это поле под ним
	lock sync.Mutex
}

// New слой взаимодействия с хранилищем в памяти.
func New() *Memory {
	return &Memory{
		cfg: config.GetConfig(),
		//m:   make(map[entity.UserID]entity.Proves),
	}
}

func (s *Memory) SaveText(context.Context, *entity.Text) error {
	return nil
}
func (s *Memory) SavePass(context.Context, *entity.Pass) error {
	return nil
}
func (s *Memory) SaveCard(context.Context, *entity.Card) error {
	return nil
}
func (s *Memory) CardListGetUserID(ctx context.Context, u *entity.User) (*entity.CardList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllUsers not implemented")
}
func (s *Memory) Registry(context.Context, *entity.Authentication) error {
	return nil
}
func (s *Memory) GetByLogin(context.Context, string) (*entity.Authentication, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllUsers not implemented")
}
func (s *Memory) GetByID(context.Context, string) (*entity.Authentication, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllUsers not implemented")
}

// **********--------------------------

// Get получить конкретный URL по идентификатору этого URL и
// только если этот линк записал текущий пользователь.
func (s *Memory) Get(ctx context.Context, ind *entity.Prove) (*entity.Prove, error) {
	u, err := s.searchBySlug(ind)
	if err != nil {
		return nil, er.ErrNotFound
	}
	return u, nil
}

// search by slug
func (s *Memory) searchBySlug(ind *entity.Prove) (*entity.Prove, error) {
	shorts := entity.Proves{}
	for _, uid := range s.m {
		for j := 0; j < len(uid); j++ {
			shorts = append(shorts, uid[j])
		}
	}
	for _, short := range shorts {
		if short.Slug == ind.Slug {
			ind.URL = short.URL
			ind.UserID = short.UserID
			ind.Del = short.Del
			break
		}
	}
	return ind, nil
}

// GetAllUrls получить все URL
func (s *Memory) GetAllUrls() (entity.CountURLs, error) {
	var c int
	for _, item := range s.m {
		c += len(item)
	}

	return entity.CountURLs(c), nil
}

// GetAllUsers получить всех пользователей
func (s *Memory) GetAllUsers() (entity.CountUsers, error) {
	return entity.CountUsers(len(s.m)), nil
}

// GetAll получить все URL пользователя по идентификатору.
func (s *Memory) GetAll(ctx context.Context, u *entity.User) (*entity.User, error) {
	l := entity.List{}
	for _, i := range s.m[entity.UserID(u.UserID)] {

		l.URL = i.URL
		//l.Slug = i.Slug
		l.ShortURL = entity.URL(scripts.GetHost(s.cfg.HTTP, string(i.Slug)))
		u.Urls = append(u.Urls, l)
	}
	return u, nil
}

// Put обновить данные в память.
func (s *Memory) Put(ctx context.Context, ind *entity.Prove) error {
	ln := len(s.m[ind.UserID])
	if ln < 1 {
		s.Post(ctx, ind)
		return nil
	}
	for j := 0; j < ln; j++ {
		if s.m[ind.UserID][j].Slug == ind.Slug {
			s.m[ind.UserID][j].URL = ind.URL
			s.m[ind.UserID][j].Del = ind.Del
			return nil
		}
	}
	return s.Post(ctx, ind)
}

// Post сохранить данные в память.
//
//	{
//		UserID: ShortURL{
//			UserID: str1
//		}
//	}
func (s *Memory) Post(ctx context.Context, ind *entity.Prove) error {
	s.m[ind.UserID] = append(s.m[ind.UserID], *ind)
	return nil
}

// Delete - удаляет URLы переданный в запросе, только если есть права данного пользователя.
func (s *Memory) Delete(ctx context.Context, u *entity.User) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.delete(u)
}

func (s *Memory) delete(u *entity.User) error {
	for j := 0; j < len(s.m[entity.UserID(u.UserID)]); j++ {
		for _, slug := range u.DelLink {
			if s.m[entity.UserID(u.UserID)][j].Slug == slug {
				// изменяет флаг del на true, в результате url
				// становиться недоступным для пользователя
				s.m[entity.UserID(u.UserID)][j].Del = true
			}
		}
	}
	return nil
}

// Read - читает данные из памяти.
func (s *Memory) Read() error {
	return nil
}

// Save - сохраняет данные в памяти.
func (s *Memory) Save() error {
	return nil
}
