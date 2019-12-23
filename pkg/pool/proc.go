package pool

import (
	"bufio"
	"os"
	"strconv"
)

func GetEntropyAvail() int {
	path := "/proc/sys/kernel/random/entropy_avail"
	return readIntFromFile(path)
}

func GetWriteWakeupThreshold() int {
	path := "/proc/sys/kernel/random/write_wakeup_threshold"
	return readIntFromFile(path)
}

func GetPoolsize() int {
	path := "/proc/sys/kernel/random/poolsize"
	return readIntFromFile(path)
}

func readIntFromFile(path string) int {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	i, err := strconv.Atoi(scanner.Text())

	if err != nil {
		panic(err)
	}
	return i
}
