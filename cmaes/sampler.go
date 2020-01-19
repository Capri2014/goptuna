package cmaes

import (
	"math/rand"
	"sync"

	"github.com/c-bata/goptuna"
)

var _ goptuna.RelativeSampler = &Sampler{}

// Sampler returns the next search points by using TPE.
type Sampler struct {
	rng *rand.Rand
	mu  sync.Mutex
}

func (s *Sampler) Sample(
	study *goptuna.Study,
	trial goptuna.FrozenTrial,
	searchSpace map[string]interface{},
) (map[string]float64, error) {
	panic("implement me")
}

// NewSampler returns the TPE sampler.
func NewSampler(opts ...SamplerOption) *Sampler {
	sampler := &Sampler{
		rng: rand.New(rand.NewSource(0)),
	}

	for _, opt := range opts {
		opt(sampler)
	}
	return sampler
}
