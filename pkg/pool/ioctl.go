package pool

import (
	"syscall"
	"unsafe"

	"github.com/chtjonas/entropyd/pkg/entropy"
)

type randPoolInfo struct {
	entropyCount int
	bufSize      int
	buf          []byte
}

const rndAddEntropy = 0x40085203

func (pool *EntropyPool) AddEntropy(sample *entropy.Sample) {
	arg := unsafe.Pointer(&randPoolInfo{
		entropyCount: sample.GetBits(),
		bufSize:      sample.GetSize(),
		buf:          sample.GetData(),
	})
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, uintptr(pool.randFd), uintptr(rndAddEntropy), uintptr(arg))
	if ep != 0 {
		err := syscall.Errno(ep)
		panic(err)
	}
}
