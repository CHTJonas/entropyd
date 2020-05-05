package pool

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"syscall"
)

type EntropyPool struct {
	randFd   int
	availFd  int
	threshFd int
	sizeFd   int
}

func OpenPool() *EntropyPool {
	randFd := getFd("/dev/random", syscall.O_RDWR, 0666)
	availFd := getFd("/proc/sys/kernel/random/entropy_avail", syscall.O_RDONLY, 0444)
	threshFd := getFd("/proc/sys/kernel/random/write_wakeup_threshold", syscall.O_RDONLY, 0644)
	sizeFd := getFd("/proc/sys/kernel/random/poolsize", syscall.O_RDONLY, 0444)
	return &EntropyPool{
		randFd:   randFd,
		availFd:  availFd,
		threshFd: threshFd,
		sizeFd:   sizeFd,
	}
}

func (pool *EntropyPool) Cleardown() {
	syscall.Close(pool.randFd)
	syscall.Close(pool.availFd)
	syscall.Close(pool.threshFd)
	syscall.Close(pool.sizeFd)
}

func getFd(path string, mode int, perm uint32) int {
	fd, err := syscall.Open(path, mode, perm)
	if err != nil {
		panic(err)
	}
	fmt.Printf("FD for %s is %d\n", path, fd)
	return fd
}

func readIntFromFd(fd int) int {
	buffer := make([]byte, 10, 100)
	n, err := syscall.Read(fd, buffer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read %d bytes from %d: %s", n, fd, buffer)
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

func (pool *EntropyPool) GetEntropyAvail() int {
	return readIntFromFd(pool.availFd)
}

func (pool *EntropyPool) GetWriteWakeupThreshold() int {
	return readIntFromFd(pool.threshFd)
}

func (pool *EntropyPool) GetPoolSize() int {
	return readIntFromFd(pool.sizeFd)
}

func (pool *EntropyPool) GetBitsNeeded(maxBits int) int {
	entropyAvailable := pool.GetEntropyAvail()
	poolCapacity := pool.GetPoolSize()
	entropyTarget := poolCapacity * 6 / 7
	bitsNeeded := computeBitsNeeded(entropyAvailable, entropyTarget, poolCapacity, maxBits)
	fmt.Printf("Entropy available: %d. Entropy target: %d. Entropy needed: %d.\n", entropyAvailable, entropyTarget, bitsNeeded)
	return bitsNeeded
}

// ComputeBitsNeeded reverses the kernel's asymptotic algorithm to determine how much entropy
// we need to add to the pool in order for entropy_avail to reach the target.
// https://elixir.bootlin.com/linux/v5.4.35/source/drivers/char/random.c#L727
func computeBitsNeeded(entropyAvailable int, entropyTarget int, poolCapacity int, maxBits int) int {
	bitsWanted := 0
	for {
		bits := float64(4*(entropyTarget-entropyAvailable)*poolCapacity) / float64(3*(poolCapacity-entropyAvailable))
		if bitsWanted >= maxBits {
			return maxBits
		}
		if bits <= float64(poolCapacity)/2 {
			return int(math.Ceil(bits + float64(bitsWanted)))
		}
		// Pretend we added poolCapacity / 2 and go around again
		bitsWanted = poolCapacity / 2
		entropyAvailable += 3 * (poolCapacity - entropyAvailable) / 8
	}
}
