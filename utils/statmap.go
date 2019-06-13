package utils

func IncStat(key string, stats map[string]int) {
	_, ok := stats[key]
	if !ok {
		stats[key] = 0
	}
	stats[key]++
}

func SetStat(key string, stats map[string]int, value int) {
	stats[key] = value
}

func StatValue(key string, stats map[string]int) int {
	v, ok := stats[key]
	if !ok {
		return 0
	}
	return v
}
