package backend

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/iliyanmotovski/raytracer/backend/vector"
)

// ConfigRepository is an abstraction over some repository
type ConfigRepository interface {
	Get(context.Context) (*Config, error)
	Upsert(context.Context, *Config) (*Config, error)
}

// Config represents the scene configuration
type Config struct {
	Light    *vector.Vector
	Scene    *vector.Vector
	Polygons Polygons
}

// Configurator is an abstraction over some configurator - txt file, yaml etc...
type Configurator interface {
	Parse(ctx context.Context, configRepo ConfigRepository) (*Config, error)
}

// textFileConfigurator is a concrete implementation of txt file configurator
type textFileConfigurator struct {
	path string
}

// NewTextFileConfigurator creates new textFileConfigurator
func NewTextFileConfigurator(path string) Configurator {
	return &textFileConfigurator{path: path}
}

// Parse reads the txt file and populates the Config
func (t *textFileConfigurator) Parse(ctx context.Context, configRepo ConfigRepository) (*Config, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	file, err := os.Open(t.path)
	if err != nil {
		return &Config{}, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for scanner.Scan() {
		fileTextLines = append(fileTextLines, scanner.Text())
	}

	c := &Config{Light: &vector.Vector{}, Scene: &vector.Vector{}}

	j := 0
	polyCount := 0
	for i, line := range fileTextLines {
		if i == 0 || i == 1 {
			x, y, err := parseFirstAndSecond(line)
			if err != nil {
				return &Config{}, err
			}
			switch i {
			case 0:
				c.Scene.X = x
				c.Scene.Y = y
			case 1:
				c.Light.X = x
				c.Light.Y = y
			}
			continue
		}

		if i == 2 {
			polyCount, err = strconv.Atoi(line)
			if err != nil {
				return &Config{}, err
			}
			c.Polygons = make(Polygons, polyCount)
			continue
		}

		polyCoords := strings.Split(line, " ")

		verteciesCount, err := strconv.Atoi(polyCoords[0])
		if err != nil {
			return &Config{}, err
		}

		polyCoords = polyCoords[1:len(polyCoords)]

		c.Polygons[j] = &Polygon{VerticesCount: verteciesCount}
		c.Polygons[j].Loop = make(vector.Loop, verteciesCount)

		k := 0
		for g := range c.Polygons[j].Loop {
			x, err := strconv.ParseFloat(polyCoords[k], 64)
			if err != nil {
				return &Config{}, err
			}
			y, err := strconv.ParseFloat(polyCoords[k+1], 64)
			if err != nil {
				return &Config{}, err
			}

			c.Polygons[j].Loop[g] = &vector.Vector{X: x, Y: y}
			k += 2
		}
		j++
	}

	if err := scanner.Err(); err != nil {
		return &Config{}, err
	}

	persisted, err := configRepo.Upsert(ctx, c)
	if err != nil {
		return &Config{}, err
	}

	return persisted, nil
}

// parses the first and second lines of the config as they are identical
func parseFirstAndSecond(line string) (float64, float64, error) {
	x, y, err := 0.0, 0.0, error(nil)
	scene := strings.Split(line, " ")

	if x, err = strconv.ParseFloat(scene[0], 64); err != nil {
		return 0, 0, err
	}
	if y, err = strconv.ParseFloat(scene[1], 64); err != nil {
		return 0, 0, err
	}

	return x, y, nil
}
