// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	entity "github.com/SETTER2000/prove/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockIndra is a mock of Indra interface.
type MockIndra struct {
	ctrl     *gomock.Controller
	recorder *MockIndraMockRecorder
}

// MockIndraMockRecorder is the mock recorder for MockIndra.
type MockIndraMockRecorder struct {
	mock *MockIndra
}

// NewMockIndra creates a new mock instance.
func NewMockIndra(ctrl *gomock.Controller) *MockIndra {
	mock := &MockIndra{ctrl: ctrl}
	mock.recorder = &MockIndraMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIndra) EXPECT() *MockIndraMockRecorder {
	return m.recorder
}

// CardListUserID mocks base method.
func (m *MockIndra) CardListUserID(ctx context.Context, u *entity.User) (*entity.CardList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CardListUserID", ctx, u)
	ret0, _ := ret[0].(*entity.CardList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CardListUserID indicates an expected call of CardListUserID.
func (mr *MockIndraMockRecorder) CardListUserID(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CardListUserID", reflect.TypeOf((*MockIndra)(nil).CardListUserID), ctx, u)
}

// ReadService mocks base method.
func (m *MockIndra) ReadService() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadService")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadService indicates an expected call of ReadService.
func (mr *MockIndraMockRecorder) ReadService() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadService", reflect.TypeOf((*MockIndra)(nil).ReadService))
}

// Register mocks base method.
func (m *MockIndra) Register(arg0 context.Context, arg1 *entity.Authentication) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockIndraMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockIndra)(nil).Register), arg0, arg1)
}

// SaveCard mocks base method.
func (m *MockIndra) SaveCard(arg0 context.Context, arg1 *entity.Card) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCard", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCard indicates an expected call of SaveCard.
func (mr *MockIndraMockRecorder) SaveCard(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCard", reflect.TypeOf((*MockIndra)(nil).SaveCard), arg0, arg1)
}

// SavePass mocks base method.
func (m *MockIndra) SavePass(arg0 context.Context, arg1 *entity.Pass) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SavePass", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SavePass indicates an expected call of SavePass.
func (mr *MockIndraMockRecorder) SavePass(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePass", reflect.TypeOf((*MockIndra)(nil).SavePass), arg0, arg1)
}

// SaveService mocks base method.
func (m *MockIndra) SaveService() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveService")
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveService indicates an expected call of SaveService.
func (mr *MockIndraMockRecorder) SaveService() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveService", reflect.TypeOf((*MockIndra)(nil).SaveService))
}

// SaveText mocks base method.
func (m *MockIndra) SaveText(arg0 context.Context, arg1 *entity.Text) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveText", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveText indicates an expected call of SaveText.
func (mr *MockIndraMockRecorder) SaveText(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveText", reflect.TypeOf((*MockIndra)(nil).SaveText), arg0, arg1)
}

// ShortLink mocks base method.
func (m *MockIndra) ShortLink(arg0 context.Context, arg1 *entity.Indra) (*entity.Indra, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShortLink", arg0, arg1)
	ret0, _ := ret[0].(*entity.Indra)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShortLink indicates an expected call of ShortLink.
func (mr *MockIndraMockRecorder) ShortLink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShortLink", reflect.TypeOf((*MockIndra)(nil).ShortLink), arg0, arg1)
}

// UserFindByID mocks base method.
func (m *MockIndra) UserFindByID(arg0 context.Context, arg1 string) (*entity.Authentication, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserFindByID", arg0, arg1)
	ret0, _ := ret[0].(*entity.Authentication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserFindByID indicates an expected call of UserFindByID.
func (mr *MockIndraMockRecorder) UserFindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserFindByID", reflect.TypeOf((*MockIndra)(nil).UserFindByID), arg0, arg1)
}

// UserFindByLogin mocks base method.
func (m *MockIndra) UserFindByLogin(arg0 context.Context, arg1 string) (*entity.Authentication, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserFindByLogin", arg0, arg1)
	ret0, _ := ret[0].(*entity.Authentication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserFindByLogin indicates an expected call of UserFindByLogin.
func (mr *MockIndraMockRecorder) UserFindByLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserFindByLogin", reflect.TypeOf((*MockIndra)(nil).UserFindByLogin), arg0, arg1)
}

// MockIndraRepo is a mock of IndraRepo interface.
type MockIndraRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIndraRepoMockRecorder
}

// MockIndraRepoMockRecorder is the mock recorder for MockIndraRepo.
type MockIndraRepoMockRecorder struct {
	mock *MockIndraRepo
}

// NewMockIndraRepo creates a new mock instance.
func NewMockIndraRepo(ctrl *gomock.Controller) *MockIndraRepo {
	mock := &MockIndraRepo{ctrl: ctrl}
	mock.recorder = &MockIndraRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIndraRepo) EXPECT() *MockIndraRepoMockRecorder {
	return m.recorder
}

// CardListGetUserID mocks base method.
func (m *MockIndraRepo) CardListGetUserID(arg0 context.Context, arg1 *entity.User) (*entity.CardList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CardListGetUserID", arg0, arg1)
	ret0, _ := ret[0].(*entity.CardList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CardListGetUserID indicates an expected call of CardListGetUserID.
func (mr *MockIndraRepoMockRecorder) CardListGetUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CardListGetUserID", reflect.TypeOf((*MockIndraRepo)(nil).CardListGetUserID), arg0, arg1)
}

// Delete mocks base method.
func (m *MockIndraRepo) Delete(arg0 context.Context, arg1 *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIndraRepoMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIndraRepo)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockIndraRepo) Get(arg0 context.Context, arg1 *entity.Indra) (*entity.Indra, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*entity.Indra)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIndraRepoMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIndraRepo)(nil).Get), arg0, arg1)
}

// GetAll mocks base method.
func (m *MockIndraRepo) GetAll(arg0 context.Context, arg1 *entity.User) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0, arg1)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockIndraRepoMockRecorder) GetAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockIndraRepo)(nil).GetAll), arg0, arg1)
}

// GetAllUrls mocks base method.
func (m *MockIndraRepo) GetAllUrls() (entity.CountURLs, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUrls")
	ret0, _ := ret[0].(entity.CountURLs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUrls indicates an expected call of GetAllUrls.
func (mr *MockIndraRepoMockRecorder) GetAllUrls() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUrls", reflect.TypeOf((*MockIndraRepo)(nil).GetAllUrls))
}

// GetAllUsers mocks base method.
func (m *MockIndraRepo) GetAllUsers() (entity.CountUsers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers")
	ret0, _ := ret[0].(entity.CountUsers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockIndraRepoMockRecorder) GetAllUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockIndraRepo)(nil).GetAllUsers))
}

// GetByID mocks base method.
func (m *MockIndraRepo) GetByID(arg0 context.Context, arg1 string) (*entity.Authentication, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*entity.Authentication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIndraRepoMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIndraRepo)(nil).GetByID), arg0, arg1)
}

// GetByLogin mocks base method.
func (m *MockIndraRepo) GetByLogin(arg0 context.Context, arg1 string) (*entity.Authentication, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByLogin", arg0, arg1)
	ret0, _ := ret[0].(*entity.Authentication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByLogin indicates an expected call of GetByLogin.
func (mr *MockIndraRepoMockRecorder) GetByLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByLogin", reflect.TypeOf((*MockIndraRepo)(nil).GetByLogin), arg0, arg1)
}

// Post mocks base method.
func (m *MockIndraRepo) Post(arg0 context.Context, arg1 *entity.Indra) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Post indicates an expected call of Post.
func (mr *MockIndraRepoMockRecorder) Post(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockIndraRepo)(nil).Post), arg0, arg1)
}

// Put mocks base method.
func (m *MockIndraRepo) Put(arg0 context.Context, arg1 *entity.Indra) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockIndraRepoMockRecorder) Put(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockIndraRepo)(nil).Put), arg0, arg1)
}

// Read mocks base method.
func (m *MockIndraRepo) Read() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read")
	ret0, _ := ret[0].(error)
	return ret0
}

// Read indicates an expected call of Read.
func (mr *MockIndraRepoMockRecorder) Read() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockIndraRepo)(nil).Read))
}

// Registry mocks base method.
func (m *MockIndraRepo) Registry(arg0 context.Context, arg1 *entity.Authentication) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Registry", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Registry indicates an expected call of Registry.
func (mr *MockIndraRepoMockRecorder) Registry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Registry", reflect.TypeOf((*MockIndraRepo)(nil).Registry), arg0, arg1)
}

// Save mocks base method.
func (m *MockIndraRepo) Save() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save")
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIndraRepoMockRecorder) Save() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIndraRepo)(nil).Save))
}

// SaveCard mocks base method.
func (m *MockIndraRepo) SaveCard(arg0 context.Context, arg1 *entity.Card) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCard", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCard indicates an expected call of SaveCard.
func (mr *MockIndraRepoMockRecorder) SaveCard(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCard", reflect.TypeOf((*MockIndraRepo)(nil).SaveCard), arg0, arg1)
}

// SavePass mocks base method.
func (m *MockIndraRepo) SavePass(arg0 context.Context, arg1 *entity.Pass) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SavePass", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SavePass indicates an expected call of SavePass.
func (mr *MockIndraRepoMockRecorder) SavePass(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePass", reflect.TypeOf((*MockIndraRepo)(nil).SavePass), arg0, arg1)
}

// SaveText mocks base method.
func (m *MockIndraRepo) SaveText(arg0 context.Context, arg1 *entity.Text) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveText", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveText indicates an expected call of SaveText.
func (mr *MockIndraRepoMockRecorder) SaveText(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveText", reflect.TypeOf((*MockIndraRepo)(nil).SaveText), arg0, arg1)
}
