// Code generated by MockGen. DO NOT EDIT.
// Source: domain/inventory/pokemon/pokemon_repository.go
//
// Generated by this command:
//
//	mockgen -package pokemon -source domain/inventory/pokemon/pokemon_repository.go -destination domain/inventory/pokemon/mock_pokemon_repository.go
//

// Package pokemon is a generated GoMock package.
package pokemon

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockPokemonRepository is a mock of PokemonRepository interface.
type MockPokemonRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPokemonRepositoryMockRecorder
	isgomock struct{}
}

// MockPokemonRepositoryMockRecorder is the mock recorder for MockPokemonRepository.
type MockPokemonRepositoryMockRecorder struct {
	mock *MockPokemonRepository
}

// NewMockPokemonRepository creates a new mock instance.
func NewMockPokemonRepository(ctrl *gomock.Controller) *MockPokemonRepository {
	mock := &MockPokemonRepository{ctrl: ctrl}
	mock.recorder = &MockPokemonRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPokemonRepository) EXPECT() *MockPokemonRepositoryMockRecorder {
	return m.recorder
}

// FindById mocks base method.
func (m *MockPokemonRepository) FindById(ctx context.Context, pokemonId int) (*Pokemon, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, pokemonId)
	ret0, _ := ret[0].(*Pokemon)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockPokemonRepositoryMockRecorder) FindById(ctx, pokemonId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockPokemonRepository)(nil).FindById), ctx, pokemonId)
}

// Save mocks base method.
func (m *MockPokemonRepository) Save(ctx context.Context, pokemon *Pokemon, userId string, quantity int, now time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, pokemon, userId, quantity, now)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockPokemonRepositoryMockRecorder) Save(ctx, pokemon, userId, quantity, now any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockPokemonRepository)(nil).Save), ctx, pokemon, userId, quantity, now)
}
