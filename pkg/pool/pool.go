package pool

import (
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

func (pool *EntropyPool) GetEntropyAvail() int {
	return readIntFromFd(pool.availFd)
}

func (pool *EntropyPool) GetWriteWakeupThreshold() int {
	return readIntFromFd(pool.threshFd)
}

func (pool *EntropyPool) GetPoolSize() int {
	return readIntFromFd(pool.sizeFd)
}

func (pool *EntropyPool) GetBitsNeeded(entropyTarget int, maxBits int) (int, int) {
	entropyAvailable := pool.GetEntropyAvail()
	poolCapacity := pool.GetPoolSize()
	bitsNeeded := computeBitsNeeded(entropyAvailable, entropyTarget, poolCapacity, maxBits)
	return entropyAvailable, bitsNeeded
}
