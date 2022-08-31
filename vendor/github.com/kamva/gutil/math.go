package gutil

// PercentOf calculate what percent [part] is of [total].
// ex. 300 is 12.5% of 2400
func PercentOf(total int64, part int64) float64 {
	if total == 0 {
		return 0
	}
	return (float64(part) * float64(100)) / float64(total)
}

func Percent(total int64, percent float64) float64 {
	return PercentFloat(float64(total), percent)
}

func PercentInt(total int64, percent int) float64 {
	return PercentFloat(float64(total), float64(percent))
}

func PercentFloat(total float64, percent float64) float64 {
	return total * percent / 100
}
