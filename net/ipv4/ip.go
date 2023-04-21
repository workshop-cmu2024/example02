package ipv4

import (
	"fmt"
	"github.com/talostrading/sonic"
	"net"
	"net/netip"
	"syscall"
	"unsafe"
)

func GetMulticastInterface(socket *sonic.Socket) (netip.Addr, error) {
	addr, err := syscall.GetsockoptInet4Addr(socket.RawFd(), syscall.IPPROTO_IP, syscall.IP_MULTICAST_IF)
	if err != nil {
		return netip.Addr{}, err
	} else {
		return netip.AddrFrom4(addr), nil
	}
}

func SetMulticastInterface(socket *sonic.Socket, iff *net.Interface) (netip.Addr, error) {
	if iff.Flags&net.FlagMulticast == 0 {
		return netip.Addr{}, fmt.Errorf("interface=%s does not support multicast", iff.Name)
	}

	addrs, err := iff.Addrs()
	if err != nil {
		return netip.Addr{}, err
	}

	var (
		interfaceAddr net.IP
		found         = false
	)
	for _, addr := range addrs {
		switch a := addr.(type) {
		case *net.IPAddr:
			interfaceAddr = a.IP
		case *net.IPNet:
			interfaceAddr = a.IP
		}
		if interfaceAddr.To4() != nil {
			found = true
			break
		}
	}

	if found {
		var addr [4]byte
		copy(addr[:], interfaceAddr)

		if err := syscall.SetsockoptInet4Addr(
			socket.RawFd(),
			syscall.IPPROTO_IP,
			syscall.IP_MULTICAST_IF,
			addr,
		); err != nil {
			return netip.Addr{}, err
		}
		return netip.AddrFrom4(addr), nil
	} else {
		return netip.Addr{}, fmt.Errorf("interface has no IPv4 address assigned")
	}
}

// AddMembership makes the given socket a member of the specified multicast IP.
func AddMembership(socket *sonic.Socket, ip netip.Addr) error {
	if !ip.Is4() && !ip.Is4In6() {
		return fmt.Errorf("expected an IPv4 address=%s", ip)
	}

	if !ip.IsMulticast() {
		return fmt.Errorf("expected a multicast address=%s", ip)
	}

	req := syscall.IPMreq{}
	copy(req.Multiaddr[:], ip.AsSlice())

	_, _, errno := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(socket.RawFd()),
		uintptr(syscall.IPPROTO_IP),
		uintptr(syscall.IP_ADD_MEMBERSHIP),
		uintptr(unsafe.Pointer(&req)),
		syscall.SizeofIPMreq,
		0)

	var err error
	if errno != 0 {
		err = errno
	}

	return err
}