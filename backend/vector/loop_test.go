package vector

import (
	"log"
	"testing"
)

func TestIsPointInLoop(t *testing.T) {
	loop := Loop{
		{131, 188},
		{54, 136},
		{86, 32},
		{220, 32},
		{238, 114},
		{209, 163},
	}

	log.Println(loop.IsPointContainedInLoop(&Vector{800, 600}, true))
}
