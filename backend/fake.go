package backend

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// FakeConfigRepository is a fake implementation of ConfigRepository
// using the testify mock package to generate mocks
type FakeConfigRepository struct {
	mock.Mock
}

// Get is a no-op
func (f FakeConfigRepository) Get(context.Context) (*Config, error) {
	args := f.Called()
	return args.Get(0).(*Config), args.Error(1)
}

// Upsert is a no-op
func (f FakeConfigRepository) Upsert(ctx context.Context, cfg *Config) (*Config, error) {
	args := f.Called(cfg)
	return args.Get(0).(*Config), args.Error(1)
}

// FakeSceneRepository is a fake implementation of SceneRepository
// using the testify mock package to generate mocks
type FakeSceneRepository struct {
	mock.Mock
}

// Get is a no-op
func (f FakeSceneRepository) Get(context.Context) (*Scene, error) {
	args := f.Called()
	return args.Get(0).(*Scene), args.Error(1)
}

// Upsert is a no-op
func (f FakeSceneRepository) Upsert(ctx context.Context, cfg *Scene) (*Scene, error) {
	args := f.Called(cfg)
	return args.Get(0).(*Scene), args.Error(1)
}
