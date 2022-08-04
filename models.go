package main

import "github.com/dmitryikh/leaves"

type LGBModel struct {
	model *leaves.Ensemble
}

func NewLGBModel(path string) (*LGBModel, error) {
	model, err := leaves.LGEnsembleFromFile(path, true)
	if err != nil {
		return nil, err
	}
	return &LGBModel{model: model}, nil
}
