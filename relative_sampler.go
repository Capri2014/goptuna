package goptuna

// RelativeSampler returns the next search points
type RelativeSampler interface {
	// InferRelativeSearchSpace infer the search space that will be used by relative sampling in the target trial.
	//
	// This method is called right before SampleRelative method, and the search space
	// returned by this method is pass to it. The parameters not  contained in the
	// search space will be sampled by using SampleIndependent method.
	InferRelativeSearchSpace(*Study, Trial) (map[string]interface{}, error)

	// SampleRelative samples multiple dimensional parameters in a given search space.
	//
	// This method is called once at the beginning of each trial, i.e., right before the
	// evaluation of the objective function. This method is suitable for sampling algorithms
	// that use relationship between parameters such as Gaussian Process and CMA-ES.
	SampleRelative(*Study, FrozenTrial, string, interface{}) (float64, error)
}

// IntersectionSearchSpace return return the intersection search space of the Study.
//
// Intersection search space contains the intersection of parameter distributions that have been
// suggested in the completed trials of the study so far.
// If there are multiple parameters that have the same name but different distributions,
// neither is included in the resulting search space
// (i.e., the parameters with dynamic value ranges are excluded).
func IntersectionSearchSpace(study *Study) (map[string]interface{}, error) {
	var searchSpace map[string]interface{}

	trials, err := study.GetTrials()
	if err != nil {
		return nil, err
	}

	for i := range trials {
		if trials[i].State == TrialStateComplete {
			continue
		}

		if searchSpace == nil {
			searchSpace = trials[i].Distributions
			continue
		}

		exists := func(name string) bool {
			for name2 := range trials[i].Distributions {
				if name == name2 {
					return true
				}
			}
			return false
		}

		deleteParams := make([]string, 0, len(searchSpace))
		for name := range searchSpace {
			if !exists(name) {
				deleteParams = append(deleteParams, name)
			} else if trials[i].Distributions[name] != searchSpace[name] {
				deleteParams = append(deleteParams, name)
			}
		}

		for j := range deleteParams {
			delete(searchSpace, deleteParams[j])
		}
	}
	return searchSpace, nil
}
