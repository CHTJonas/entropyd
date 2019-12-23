package pool

import "syscall"

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
