package main

type ZScoreNorm struct {
	mean float64
	std  float64
}

func NewZScoreNorm(mean, std float64) *ZScoreNorm {
	return &ZScoreNorm{mean: mean, std: std}
}

func (z *ZScoreNorm) Normalize(x float64) float64 {
	return (x - z.mean) / z.std
}

func (z *ZScoreNorm) Denormalize(x float64) float64 {
	return x*z.std + z.mean
}
