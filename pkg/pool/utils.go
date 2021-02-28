package pool

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"syscall"
)

var ErrTargetTooLarge = errors.New("Target pool size must be less than full")

// ComputeBitsNeeded reverses the kernel's asymptotic algorithm to determine
// how much entropy we need to add to the pool in order to reach the target.
// https://elixir.bootlin.com/linux/v5.4.35/source/drivers/char/random.c#L727
func computeBitsNeeded(entropyAvailable int, entropyTarget int, poolCapacity int, maxBits int) int {
	if entropyTarget >= poolCapacity {
		panic(ErrTargetTooLarge)
	}

	ea := float64(entropyAvailable)
	et := float64(entropyTarget)
	pc := float64(poolCapacity)
	mb := float64(maxBits)
	halfCapacity := pc / 2
	bitsWanted := float64(0)
	for {
		numerator := 4 * (et - ea) * pc
		denominator := 3 * (pc - ea)
		bitsThisRound := numerator / denominator
		// Linux only adds half a pool at once in order to compensate for wonky simplified formulae.
		if bitsThisRound <= halfCapacity {
			return int(math.Min(math.Ceil(bitsThisRound+bitsWanted), mb))
		}
		// Pretend we added half a pool and go around again.
		bitsWanted += halfCapacity
		if bitsWanted >= mb {
			return maxBits
		}
		ea += 3 * (pc - ea) / 8
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
