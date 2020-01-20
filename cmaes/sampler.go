package cmaes

import (
	"math/rand"

	"gonum.org/v1/gonum/mat"

	"github.com/c-bata/goptuna"
)

var _ goptuna.RelativeSampler = &Sampler{}

// Sampler returns the next search points by using CMA-ES.
type Sampler struct {
	rng            *rand.Rand
	nStartUpTrials int
	mu             mat.Dense
	sigma          mat.Dense
}

func (s *Sampler) Sample(
	study *goptuna.Study,
	trial goptuna.FrozenTrial,
	searchSpace map[string]interface{},
) (map[string]float64, error) {
	if searchSpace == nil || len(searchSpace) == 0 {
		return nil, nil
	}

	if len(searchSpace) == 1 {
		// TODO(c-bata): Add warn log "CMA-ES does not support optimization of 1-D search space."
		return nil, nil
	}

	trials, err := study.GetTrials()
	if err != nil {
		return nil, err
	}
	completed := make([]goptuna.FrozenTrial, 0, len(trials))
	for i := range trials {
		if trials[i].State == goptuna.TrialStateComplete {
			completed = append(completed, trials[i])
		}
	}

	if len(completed) < s.nStartUpTrials {
		return nil, err
	}

	// TODO: sample parameters.

	params := make(map[string]float64, len(searchSpace))
	return params, nil
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
