package algo

func CalculateThroughput(segSize, time int) int {
	return int(float64(segSize) / (float64(time) / 1000))
}
