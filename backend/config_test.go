package backend_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iliyanmotovski/raytracer/backend"
	"github.com/iliyanmotovski/raytracer/backend/vector"
)

const data = "800 500\n" +
	"250 300\n" +
	"3\n" +
	"3 600 200 646 133 646 261\n" +
	"6 131 188 54 136 86 32 220 32 238 114 209 163\n" +
	"5 412 364 454 251 537 257 601 350 528 430"

func TestMain(m *testing.M) {
	f, _ := os.Create("test.txt")
	f.WriteString(data)
	f.Close()

	m.Run()

	os.Remove("test.txt")
}

func TestParseConfigFromTextFile(t *testing.T) {
	config := &backend.Config{
		Light: &vector.Vector{X: 250, Y: 300},
		Scene: &vector.Vector{X: 800, Y: 500},
		Polygons: backend.Polygons{
			{
				VerticesCount: 3,
				Loop: vector.Loop{
					{X: 600, Y: 200},
					{X: 646, Y: 133},
					{X: 646, Y: 261},
				},
			},
			{
				VerticesCount: 6,
				Loop: vector.Loop{
					{X: 131, Y: 188},
					{X: 54, Y: 136},
					{X: 86, Y: 32},
					{X: 220, Y: 32},
					{X: 238, Y: 114},
					{X: 209, Y: 163},
				},
			},
			{
				VerticesCount: 5,
				Loop: vector.Loop{
					{X: 412, Y: 364},
					{X: 454, Y: 251},
					{X: 537, Y: 257},
					{X: 601, Y: 350},
					{X: 528, Y: 430},
				},
			},
		},
	}

	configRepo := new(backend.FakeConfigRepository)
	configRepo.On("Upsert", config).Return(config, nil)

	c := backend.NewTextFileConfigurator("test.txt")
	got, err := c.Parse(context.Background(), configRepo)

	assert.Nil(t, err)
	assert.Equal(t, config, got)

	configRepo.AssertExpectations(t)
}

func TestParseConfigFromNonexistentFile(t *testing.T) {
	c := backend.NewTextFileConfigurator("non-existent.txt")
	_, err := c.Parse(context.Background(), nil)

	want := &os.PathError{
		Op:   "open",
		Path: "non-existent.txt",
		Err:  err.(*os.PathError).Err,
	}

	assert.Equal(t, want, err)
}

func TestParseConfigFromTextFileWithRepositoryFailure(t *testing.T) {
	config := &backend.Config{
		Light: &vector.Vector{X: 250, Y: 300},
		Scene: &vector.Vector{X: 800, Y: 500},
		Polygons: backend.Polygons{
			{
				VerticesCount: 3,
				Loop: vector.Loop{
					{X: 600, Y: 200},
					{X: 646, Y: 133},
					{X: 646, Y: 261},
				},
			},
			{
				VerticesCount: 6,
				Loop: vector.Loop{
					{X: 131, Y: 188},
					{X: 54, Y: 136},
					{X: 86, Y: 32},
					{X: 220, Y: 32},
					{X: 238, Y: 114},
					{X: 209, Y: 163},
				},
			},
			{
				VerticesCount: 5,
				Loop: vector.Loop{
					{X: 412, Y: 364},
					{X: 454, Y: 251},
					{X: 537, Y: 257},
					{X: 601, Y: 350},
					{X: 528, Y: 430},
				},
			},
		},
	}

	configRepo := new(backend.FakeConfigRepository)
	configRepo.On("Upsert", config).Return(&backend.Config{}, errors.New("error"))

	c := backend.NewTextFileConfigurator("test.txt")
	_, err := c.Parse(context.Background(), configRepo)

	assert.Equal(t, errors.New("error"), err)

	configRepo.AssertExpectations(t)
}
