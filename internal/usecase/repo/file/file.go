package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/internal/entity"
	"github.com/SETTER2000/prove/scripts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
)

// InFiles -.
type (
	producer struct {
		file    *os.File
		encoder *json.Encoder
	}

	consumer struct {
		file    *os.File
		cfg     *config.Config
		decoder *json.Decoder
	}

	InFiles struct {
		//lock sync.Mutex // <-- этот мьютекс защищает
		m   map[entity.UserID]entity.Proves
		cfg *config.Config
		r   *consumer
		w   *producer
	}
)

// New слой взаимодействия с файловым хранилищем
func New() *InFiles {
	return &InFiles{
		cfg: config.GetConfig(),
		m:   make(map[entity.UserID]entity.Proves),
		// создаём новый потребитель
		r: NewConsumer(),
		// создаём новый производитель
		w: NewProducer(),
	}
}

// NewProducer производитель
func NewProducer() *producer {
	file, _ := os.OpenFile(config.GetConfig().FileStorage, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}
}

// NewConsumer потребитель
func NewConsumer() *consumer {
	file, _ := os.OpenFile(config.GetConfig().FileStorage, os.O_RDONLY|os.O_CREATE, 0777)
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}
}

// Post сохранить данные в файл
func (i *InFiles) Post(ctx context.Context, ind *entity.Prove) error {
	//i.lock.Lock()
	//defer i.lock.Unlock()
	return i.post(ind)
}

func (i *InFiles) post(ind *entity.Prove) error {
	i.m[ind.UserID] = append(i.m[ind.UserID], *ind)
	return nil
}

// Put обновить данные в файле
func (i *InFiles) Put(ctx context.Context, ind *entity.Prove) error {
	ln := len(i.m[ind.UserID])
	if ln < 1 {
		i.Post(ctx, ind)
		return nil
	}
	for j := 0; j < ln; j++ {
		if i.m[entity.UserID(ind.UserID)][j].Slug == ind.Slug {
			i.m[entity.UserID(ind.UserID)][j].URL = ind.URL
		}
	}
	return i.Post(ctx, ind)
}

// Close - закрывает открытый файл
func (p *producer) Close() error {
	return p.file.Close()
}

// Get получить конкретный URL по идентификатору
func (i *InFiles) Get(ctx context.Context, ind *entity.Prove) (*entity.Prove, error) {
	return i.searchBySlug(ind)
}

func (i *InFiles) searchUID(ind *entity.Prove) (*entity.Prove, error) {
	for _, short := range i.m[entity.UserID(ind.UserID)] {
		if short.Slug == ind.Slug {
			ind.URL = short.URL
			ind.UserID = short.UserID
			ind.Del = short.Del
			fmt.Println("НАШЁЛ URL: ", ind.URL)
			break
		}
	}
	return ind, nil
}

// search by slug
func (i *InFiles) searchBySlug(ind *entity.Prove) (*entity.Prove, error) {
	shorts := entity.Proves{}
	for _, uid := range i.m {
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
func (i *InFiles) getAll() error {
	ind := &entity.Prove{}
	for {
		if err := i.r.decoder.Decode(&ind); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		i.m[entity.UserID(ind.UserID)] = append(i.m[entity.UserID(ind.UserID)], *ind)
	}
	return nil
}
func (i *InFiles) getAllUserID(u *entity.User) (*entity.User, error) {
	lst := entity.List{}
	shorts := i.m[entity.UserID(u.UserID)]
	for _, short := range shorts {
		if short.UserID == entity.UserID(u.UserID) {
			lst.URL = entity.URL(short.URL)
			//lst.Slug = short.Slug
			lst.ShortURL = entity.URL(scripts.GetHost(i.cfg.HTTP, string(short.Slug)))
			u.Urls = append(u.Urls, lst)
		}
	}
	return u, nil
}

// GetAllUrls получить все URL
func (i *InFiles) GetAllUrls() (entity.CountURLs, error) {
	var c int
	for _, item := range i.m {
		c += len(item)
	}

	return entity.CountURLs(c), nil
}

// GetAllUsers получить всех пользователей
func (i *InFiles) GetAllUsers() (entity.CountUsers, error) {
	return entity.CountUsers(len(i.m)), nil
}

// GetAll получить все URL пользователя по идентификатору
func (i *InFiles) GetAll(ctx context.Context, u *entity.User) (*entity.User, error) {
	//i.lock.Lock()
	//defer i.lock.Unlock()
	return i.getAllUserID(u)
}

// Delete - удаляет URLы переданный в запросе, только если есть права данного пользователя
func (i *InFiles) Delete(ctx context.Context, u *entity.User) error {
	//i.lock.Lock()
	//defer i.lock.Unlock()
	return i.delete(u)
}

func (i *InFiles) delete(u *entity.User) error {
	for j := 0; j < len(i.m[entity.UserID(u.UserID)]); j++ {
		for _, slug := range u.DelLink {
			if i.m[entity.UserID(u.UserID)][j].Slug == slug {
				// изменяет флаг del на true, в результате url становиться недоступным для пользователя
				i.m[entity.UserID(u.UserID)][j].Del = true
			}
		}
	}
	return nil
}

// Close закрывает соединение потребителя
func (c *consumer) Close() error {
	return c.file.Close()
}

// ---

func (i *InFiles) SavePass(context.Context, *entity.Pass) error {
	return nil
}
func (i *InFiles) SaveText(context.Context, *entity.Text) error {
	return nil
}
func (i *InFiles) SaveGroup(ctx context.Context, c *entity.Group) error {
	return status.Errorf(codes.Unimplemented, "method SaveGroup not implemented")
}
func (i *InFiles) GroupList(context.Context) (*entity.GroupList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GroupList not implemented")
}
func (i *InFiles) SaveTask(ctx context.Context, c *entity.Task) error {
	return status.Errorf(codes.Unimplemented, "method SaveGroup not implemented")
}
func (i *InFiles) TaskList(context.Context) (*entity.TaskList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GroupList not implemented")
}
func (i *InFiles) SaveCard(context.Context, *entity.Card) error {
	return nil
}
func (i *InFiles) CardListGetUserID(ctx context.Context, u *entity.User) (*entity.CardList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllUsers not implemented")
}
func (i *InFiles) Registry(ctx context.Context, a *entity.Authentication) error {
	if err := a.Validate(); err != nil {
		return err
	}
	if err := a.BeforeCreate(); err != nil {
		return err
	}

	return nil
}

func (i *InFiles) GetByLogin(ctx context.Context, l string) (*entity.Authentication, error) {
	var a entity.Authentication

	return &a, nil
}

func (i *InFiles) GetByID(context.Context, string) (*entity.Authentication, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllUsers not implemented")
}

// Read - читает данные из файла
func (i *InFiles) Read() error {
	for {
		if err := i.r.decoder.Decode(&i.m); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

// Save перезаписать файл с новыми данными
func (i *InFiles) Save() error {
	//переводит курсор в начало файла
	_, err := i.w.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	// очищает файл
	err = i.w.file.Truncate(0)
	if err != nil {
		return err
	}
	// пишем из памяти в файл
	return i.w.encoder.Encode(i.m)
}
