package goptuna_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/c-bata/goptuna"
)

func TestIntersectionSearchSpace(t *testing.T) {
	tests := []struct {
		name         string
		study        func() *goptuna.Study
		expectedKeys []string
		want         map[string]interface{}
		wantErr      bool
	}{
		{
			name: "No trial",
			study: func() *goptuna.Study {
				study, err := goptuna.CreateStudy("sampler_test")
				if err != nil {
					panic(err)
				}
				return study
			},
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "One trial",
			study: func() *goptuna.Study {
				study, err := goptuna.CreateStudy("sampler_test")
				if err != nil {
					panic(err)
				}

				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					x, _ := trial.SuggestInt("x", 0, 10)
					y, _ := trial.SuggestUniform("y", -3, 3)
					return float64(x) + y, nil
				}, 1); err != nil {
					panic(err)
				}
				return study
			},
			want: map[string]interface{}{
				"x": goptuna.IntUniformDistribution{
					High: 10,
					Low:  0,
				},
				"y": goptuna.UniformDistribution{
					High: 3,
					Low:  -3,
				},
			},
			wantErr: false,
		},
		{
			name: "Second trial (only 'y' parameter is suggested in this trial)",
			study: func() *goptuna.Study {
				study, err := goptuna.CreateStudy("sampler_test")
				if err != nil {
					panic(err)
				}

				// First Trial
				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					x, _ := trial.SuggestInt("x", 0, 10)
					y, _ := trial.SuggestUniform("y", -3, 3)
					return float64(x) + y, nil
				}, 1); err != nil {
					panic(err)
				}

				// Second Trial
				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					y, _ := trial.SuggestUniform("y", -3, 3)
					return y, nil
				}, 1); err != nil {
					panic(err)
				}
				return study
			},
			want: map[string]interface{}{
				"y": goptuna.UniformDistribution{
					High: 3,
					Low:  -3,
				},
			},
			wantErr: false,
		},
		{
			name: "Failed or pruned trials are not considered in the calculation of an intersection search space.",
			study: func() *goptuna.Study {
				study, err := goptuna.CreateStudy("sampler_test")
				if err != nil {
					panic(err)
				}

				// First Trial
				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					x, _ := trial.SuggestInt("x", 0, 10)
					y, _ := trial.SuggestUniform("y", -3, 3)
					return float64(x) + y, nil
				}, 1); err != nil {
					panic(err)
				}

				// Second Trial
				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					y, _ := trial.SuggestUniform("y", -3, 3)
					return y, nil
				}, 1); err != nil {
					panic(err)
				}

				// Failed trial (ignore error)
				_ = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					_, _ = trial.SuggestUniform("y", -3, 3)
					return 0.0, errors.New("something error")
				}, 1)
				// Pruned trial
				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					_, _ = trial.SuggestUniform("y", -3, 3)
					return 0.0, goptuna.ErrTrialPruned
				}, 1); err != nil {
					panic(err)
				}
				return study
			},
			want: map[string]interface{}{
				"y": goptuna.UniformDistribution{
					High: 3,
					Low:  -3,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := goptuna.IntersectionSearchSpace(tt.study())
			if (err != nil) != tt.wantErr {
				t.Errorf("IntersectionSearchSpace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("IntersectionSearchSpace() return %d items, but want %d", len(got), len(tt.want))
			}
			for key := range tt.want {
				if distribution, ok := got[key]; !ok {
					t.Errorf("IntersectionSearchSpace() should have %s key", key)
				} else if !reflect.DeepEqual(distribution, tt.want[key]) {
					t.Errorf("IntersectionSearchSpace()[%s] = %v, want %v", key, distribution, tt.want[key])
				}
			}
		})
	}
}
