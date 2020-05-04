package pool

import (
	"math"
	"syscall"
)

type EntropyPool struct {
	fd int
}

func OpenPool() *EntropyPool {
	fd, err := syscall.Open("/dev/random", syscall.O_RDWR, 666)
	if err != nil {
		panic(err)
	}
	return &EntropyPool{
		fd: fd,
	}
}

func (pool *EntropyPool) Cleardown() {
	syscall.Close(pool.fd)
}

// ComputeBitsNeeded reverses the kernel's asymptotic algorithm to determine how much entropy
// we need to add to the pool in order for entropy_avail to reach the target.
// https://elixir.bootlin.com/linux/v5.4.35/source/drivers/char/random.c#L727
func ComputeBitsNeeded(entropyAvailable int, entropyTarget int, poolCapacity int, maxBits int) int {
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
