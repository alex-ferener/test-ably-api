package main

func MinMaxAvg(values []int64) (int64, int64, float64) {
	var max = values[0]
	var min = values[0]
	var sum int64
	var avg float64

	for _, value := range values {
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
		sum += value
	}
	avg = float64(sum / int64(len(values)))

	return min, max, avg
}
