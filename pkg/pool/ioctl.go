package pool

import (
	"fmt"
	"syscall"
	"unsafe"
)

// RndAddEntropy is the ioctl request defined by the Linux kernel.
const RndAddEntropy = 0x40085203

// MaxDataBytes is the maximum size of the ioctl request payload.
const MaxDataBytes = 1016

type randPoolInfo struct {
	entropyCount int32
	bufSize      int32
	buf          []byte
}

func (pool *EntropyPool) AddEntropy(entropy *Entropy) error {
	if err := entropy.Validate(); err != nil {
		return fmt.Errorf("failed to validate entropy: %v", err)
	}
	payload := unsafe.Pointer(&randPoolInfo{
		entropyCount: int32(entropy.Count),
		bufSize:      int32(len(entropy.Data)),
		buf:          entropy.Data,
	})
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, uintptr(pool.randFd), uintptr(RndAddEntropy), uintptr(payload))
	if ep != 0 {
		err := syscall.Errno(ep)
		return err
	}
	return nil
}
