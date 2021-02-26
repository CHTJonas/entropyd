package pool

import (
	"syscall"
	"unsafe"
)

type randPoolInfo struct {
	entropyCount int
	bufSize      int
	buf          []byte
}

const rndAddEntropy = 0x40085203

func (pool *EntropyPool) AddEntropy(entropy *Entropy) {
	if err := entropy.Validate(); err != nil {
		panic(err)
	}
	arg := unsafe.Pointer(&randPoolInfo{
		entropyCount: entropy.Count,
		bufSize:      len(entropy.Data),
		buf:          entropy.Data,
	})
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, uintptr(pool.randFd), uintptr(rndAddEntropy), uintptr(arg))
	if ep != 0 {
		err := syscall.Errno(ep)
		panic(err)
	}
}
