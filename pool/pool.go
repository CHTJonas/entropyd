package pool

import (
	"log"
	"syscall"
)

type EntropyPool struct {
	randFd   int
	availFd  int
	threshFd int
	sizeFd   int
	logger   *log.Logger
}

func OpenPool() (*EntropyPool, error) {
	randFd, err := getFd("/dev/random", syscall.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	availFd, err := getFd("/proc/sys/kernel/random/entropy_avail", syscall.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	threshFd, err := getFd("/proc/sys/kernel/random/write_wakeup_threshold", syscall.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	sizeFd, err := getFd("/proc/sys/kernel/random/poolsize", syscall.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	return &EntropyPool{
		randFd:   randFd,
		availFd:  availFd,
		threshFd: threshFd,
		sizeFd:   sizeFd,
	}, nil
}

func (pool *EntropyPool) Cleardown() {
	syscall.Close(pool.randFd)
	syscall.Close(pool.availFd)
	syscall.Close(pool.threshFd)
	syscall.Close(pool.sizeFd)
}

func (pool *EntropyPool) GetEntropyAvail() (int, error) {
	return readIntFromFd(pool.availFd)
}

func (pool *EntropyPool) GetWriteWakeupThreshold() (int, error) {
	return readIntFromFd(pool.threshFd)
}

func (pool *EntropyPool) GetPoolSize() (int, error) {
	return readIntFromFd(pool.sizeFd)
}

func (pool *EntropyPool) GetBitsNeeded(entropyTarget int, maxBits int) (int, int, error) {
	entropyAvailable, err := pool.GetEntropyAvail()
	if err != nil {
		return 0, 0, err
	}
	poolCapacity, err := pool.GetPoolSize()
	if err != nil {
		return 0, 0, err
	}
	bitsNeeded, err := computeBitsNeeded(entropyAvailable, entropyTarget, poolCapacity, maxBits)
	if err != nil {
		return 0, 0, err
	}
	return entropyAvailable, bitsNeeded, nil
}
