package usecase

import (
	"context"
	"errors"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/internal/entity"
	"github.com/SETTER2000/prove/pkg/er"
	"github.com/SETTER2000/prove/scripts"
	"log"
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

func (uc *ProveUseCase) GetAdmin(u *entity.User) bool {
	if ok := uc.repo.GetAdmin(u); !ok {
		return false
	}

	return true
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
func (uc *ProveUseCase) SaveGroup(ctx context.Context, c *entity.Group) error {
	err := c.Validate()
	if err != nil {
		return err
	}

	err = uc.repo.SaveGroup(ctx, c)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ProveUseCase) SaveSolution(ctx context.Context, c *entity.Solution) error {
	err := c.Validate()
	if err != nil {
		return err
	}

	err = uc.repo.SaveSolution(ctx, c)
	if err != nil {
		return err
	}
	return nil
}

// SaveTask создать задачу
func (uc *ProveUseCase) SaveTask(ctx context.Context, c *entity.Task) error {
	err := c.Validate()
	if err != nil {
		return err
	}

	err = uc.repo.SaveTask(ctx, c)
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
	if !scripts.IsValidUUID(auth.GroupID) {
		return er.ErrUUID
	}
	err := uc.repo.Registry(ctx, auth)
	if err != nil {
		return err
	}
	return nil
}

// GetSolution вернёт решение задачи
func (uc *ProveUseCase) GetSolution(ctx context.Context, s *entity.SolutionData) error {

	if !scripts.IsValidUUID(s.TaskID) {
		return er.ErrUUID
	}

	_, err := uc.repo.GetBalance(ctx, s)
	if err != nil {
		log.Printf("error, user credit exhausted: %s", err.Error())
		return er.ErrBadRequest
	}

	uc.repo.GetSolution(ctx, s)
	solution, err := scripts.FindAllMissingNumbers(s.Data)
	if err != nil {
		return er.ErrBadRequest
	} else {
		s.Solution = solution
	}

	return nil
}

// FindBalance получение текущего баланса пользователя
func (uc *ProveUseCase) FindBalance(ctx context.Context) (*entity.Balance, error) {
	b, err := uc.repo.Balance(ctx)
	if err == nil {
		return b, nil
	}
	return nil, ErrBadRequest
}

// GetBalance получить информацию о кредитоспособности,
// если вернёт err!=nil - кредит исчерпан
func (uc *ProveUseCase) GetBalance(ctx context.Context, s *entity.SolutionData) (float64, error) {
	r, err := uc.repo.GetBalance(ctx, s)
	if err != nil {
		log.Printf("error, user credit exhausted: %s", err.Error())
		return 0, er.ErrBadRequest
	}

	return r, nil
}

// BalanceAdd пополнить баланс
func (uc *ProveUseCase) BalanceAdd(ctx context.Context, s *entity.Balance) error {
	err := uc.repo.BalanceAdd(ctx, s)
	if err != nil {
		log.Printf("error, user credit exhausted: %s", err.Error())
		return er.ErrBadRequest
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

// TaskKey возвращает ответ по задаче
func (uc *ProveUseCase) TaskKey(ctx context.Context, u *entity.User, t *entity.Task) (*entity.Task, error) {

	ol, err := uc.repo.TaskKey(ctx, u, t)
	if err == nil {
		return ol, nil
	}
	return nil, ErrBadRequest
}

// GroupList возвращает список групп
func (uc *ProveUseCase) GroupList(ctx context.Context) (*entity.GroupList, error) {
	ol, err := uc.repo.GroupList(ctx)
	if err == nil {
		return ol, nil
	}
	return nil, ErrBadRequest
}

// TaskList возвращает список задач
func (uc *ProveUseCase) TaskList(ctx context.Context) (*entity.TaskList, error) {
	ol, err := uc.repo.TaskList(ctx)
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
