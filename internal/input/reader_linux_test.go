package input

import (
	"os"
	"strconv"
	"syscall"
	"testing"
	"unsafe"
)

// openTerminal returns the follower side of a fresh pseudo-terminal holding
// one line of input followed by end-of-file, so reading it never blocks.
func openTerminal(t *testing.T) *os.File {
	t.Helper()

	leader, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		t.Skipf("pseudo-terminal unavailable: %v", err)
	}
	t.Cleanup(func() { leader.Close() })

	var unlock int32
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, leader.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&unlock))); errno != 0 {
		t.Skipf("unlock pseudo-terminal: %v", errno)
	}
	var number uint32
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, leader.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&number))); errno != 0 {
		t.Skipf("pseudo-terminal number: %v", errno)
	}

	follower, err := os.OpenFile("/dev/pts/"+strconv.FormatUint(uint64(number), 10), os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		t.Skipf("open pseudo-terminal: %v", err)
	}
	t.Cleanup(func() { follower.Close() })

	if _, err := leader.WriteString("typed\n\x04"); err != nil {
		t.Fatalf("write to pseudo-terminal: %v", err)
	}

	return follower
}

func TestReadBytesFromTerminalReturnsNothing(t *testing.T) {
	t.Parallel()

	got, err := readBytes(openTerminal(t))
	if err != nil {
		t.Fatalf("readBytes(terminal) returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("readBytes(terminal) = %q, want nil", got)
	}
}
