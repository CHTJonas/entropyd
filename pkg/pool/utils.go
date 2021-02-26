package pool

import (
	"math"
	"strconv"
	"strings"
	"syscall"
)

// ComputeBitsNeeded reverses the kernel's asymptotic algorithm to determine how much entropy
// we need to add to the pool in order for entropy_avail to reach the target.
// https://elixir.bootlin.com/linux/v5.4.35/source/drivers/char/random.c#L727
func computeBitsNeeded(entropyAvailable int, entropyTarget int, poolCapacity int, maxBits int) int {
	bitsWanted := 0
	for {
		bitsThisRound := float64(4*(entropyTarget-entropyAvailable)*poolCapacity) / float64(3*(poolCapacity-entropyAvailable))
		// Linux only adds half a pool at once, to compensate for wonky simplified formulae.
		if bitsThisRound <= float64(poolCapacity)/2 {
			return int(math.Min(math.Ceil(bitsThisRound+float64(bitsWanted)), float64(maxBits)))
		}
		// Pretend we added poolCapacity / 2 and go around again
		bitsWanted += poolCapacity / 2
		if bitsWanted >= maxBits {
			return maxBits
		}
		entropyAvailable += 3 * (poolCapacity - entropyAvailable) / 8
	}
}

func getFd(path string, mode int, perm uint32) int {
	fd, err := syscall.Open(path, mode, perm)
	if err != nil {
		panic(err)
	}
	return fd
}

func readIntFromFd(fd int) int {
	buffer := make([]byte, 10, 100)
	n, err := syscall.Read(fd, buffer)
	if err != nil {
		panic(err)
	}
	str := strings.ReplaceAll(string(buffer[:n]), "\n", "")
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	_, err = syscall.Seek(fd, 0, 0)
	if err != nil {
		panic(err)
	}
	return i
}
