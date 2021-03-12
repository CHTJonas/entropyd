package pool

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"syscall"
)

// ComputeBitsNeeded reverses the kernel's asymptotic algorithm to determine
// how much entropy we need to add to the pool in order to reach the target.
// https://elixir.bootlin.com/linux/v5.4.35/source/drivers/char/random.c#L727
func computeBitsNeeded(entropyAvailable int, entropyTarget int, poolCapacity int, maxBits int) (int, error) {
	if entropyTarget >= poolCapacity {
		return 0, fmt.Errorf("target %d must be less than the pool capacity %d", entropyTarget, poolCapacity)
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
			return int(math.Min(math.Ceil(bitsThisRound+bitsWanted), mb)), nil
		}
		// Pretend we added half a pool and go around again.
		bitsWanted += halfCapacity
		if bitsWanted >= mb {
			return maxBits, nil
		}
		ea += 3 * (pc - ea) / 8
	}
}

func getFd(path string, mode int, perm uint32) (int, error) {
	fd, err := syscall.Open(path, mode, perm)
	if err != nil {
		return 0, fmt.Errorf("error opening %s: %v", path, err)
	}
	return fd, nil
}

func readIntFromFd(fd int) (int, error) {
	buffer := make([]byte, 10, 100)
	n, err := syscall.Read(fd, buffer)
	if err != nil {
		fmt.Errorf("error reading from descriptor %v: %v", fd, err)
	}
	str := strings.ReplaceAll(string(buffer[:n]), "\n", "")
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Errorf("error converting %s to int: %v", str, err)
	}
	_, err = syscall.Seek(fd, 0, 0)
	if err != nil {
		fmt.Errorf("error resetting descriptor %v: %v", fd, err)
	}
	return i, nil
}
