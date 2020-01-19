package cmaes

import (
	"errors"
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

func (s *Sampler) InferRelativeSearchSpace(study *goptuna.Study, trial goptuna.Trial) (map[string]interface{}, error) {
	intersection, err := goptuna.IntersectionSearchSpace(study)
	if err != nil {
		return nil, err
	}
	if intersection == nil {
		return nil, nil
	}
	searchSpace := make(map[string]interface{}, len(intersection))
	for paramName := range intersection {
		distribution, ok := intersection[paramName].(goptuna.Distribution)
		if !ok {
			return nil, errors.New("failed to cast distribution")
		}
		if distribution.Single() {
			continue
		}
		searchSpace[paramName] = distribution
	}
	return searchSpace, nil
}

func (s *Sampler) SampleRelative(*goptuna.Study, goptuna.FrozenTrial, string, interface{}) (float64, error) {
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
