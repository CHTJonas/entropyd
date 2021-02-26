package malc

import "math"

func roundToOctet(bits int) int {
	multiple := float64(bits) / 8.0
	return int(math.Ceil(multiple)) * 8
}

func clamp(value, min, max int) int {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}
