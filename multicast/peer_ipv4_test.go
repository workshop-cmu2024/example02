package multicast

import (
	"fmt"
	"github.com/talostrading/sonic"
	"github.com/talostrading/sonic/net/ipv4"
	"log"
	"net"
	"net/netip"
	"sync"
	"testing"
	"time"
)

// Listing multicast group memberships: netstat -gsv

func TestUDPPeerIPv4_Addresses(t *testing.T) {
	if len(testInterfacesIPv4) == 0 {
		return
	}

	ioc := sonic.MustIO()
	defer ioc.Close()

	{
		_, err := NewUDPPeer(ioc, "udp", net.IPv4zero.String())
		if err == nil {
			t.Fatal("should have received an error as the address is missing the port")
		}
	}
	{
		_, err := NewUDPPeer(ioc, "udp4", net.IPv4zero.String())
		if err == nil {
			t.Fatal("should have received an error as the address is missing the port")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp", "")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), net.IPv4zero.String(); given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp4", "")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), net.IPv4zero.String(); given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp", ":0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), net.IPv4zero.String(); given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp4", ":0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), net.IPv4zero.String(); given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), "127.0.0.1"; given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp4", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), "127.0.0.1"; given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp", "localhost:0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), "127.0.0.1"; given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp4", "localhost:0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), "127.0.0.1"; given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}

	log.Println("ran")
}

func TestUDPPeerIPv4_JoinInvalidGroup(t *testing.T) {
	if len(testInterfacesIPv4) == 0 {
		return
	}

	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	if err := peer.Join("0.0.0.0:4555"); err == nil {
		t.Fatal("should not have joined")
	}

	log.Println("ran")
}

func TestUDPPeerIPv4_Join(t *testing.T) {
	if len(testInterfacesIPv4) == 0 {
		return
	}

	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	if err := peer.Join("224.0.0.0"); err != nil {
		t.Fatal(err)
	}

	addr, err := ipv4.GetMulticastInterfaceAddr(peer.socket)
	if err != nil {
		t.Fatal(err)
	}
	if !addr.IsUnspecified() {
		t.Fatal("multicast address should be unspecified")
	}

	log.Println("ran")
}

func TestUDPPeerIPv4_SetLoop1(t *testing.T) {
	if len(testInterfacesIPv4) == 0 {
		return
	}

	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	if peer.Loop() {
		t.Fatal("peer should not loop packets by default")
	}

	if err := peer.SetLoop(false); err != nil {
		t.Fatal(err)
	}
	if peer.Loop() {
		t.Fatal("peer should not loop packets")
	}

	if err := peer.SetLoop(true); err != nil {
		t.Fatal(err)
	}
	if !peer.Loop() {
		t.Fatal("peer should loop packets")
	}

	if err := peer.SetLoop(false); err != nil {
		t.Fatal(err)
	}
	if peer.Loop() {
		t.Fatal("peer should not loop packets")
	}

	log.Println("ran")
}

func TestUDPPeerIPv4_DefaultOutboundInterface(t *testing.T) {
	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	fmt.Println(peer.Outbound())
}

func TestUDPPeerIPv4_SetOutboundInterfaceOnUnspecifiedIPAndPort(t *testing.T) {
	if len(testInterfacesIPv4) == 0 {
		return
	}

	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	for _, iff := range testInterfacesIPv4 {
		fmt.Println("setting", iff.Name, "as outbound")

		if err := peer.SetOutboundIPv4(iff.Name); err != nil {
			t.Fatal(err)
		}

		outboundInterface, outboundIP := peer.Outbound()
		fmt.Println("outbound for", iff.Name, outboundInterface, outboundIP)

		{
			addr, err := ipv4.GetMulticastInterfaceAddr(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%s GetMulticastInterface_Inet4Addr addr=%s\n", iff.Name, addr.String())
		}

		{
			interfaceAddr, multicastAddr, err := ipv4.GetMulticastInterfaceAddrAndGroup(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf(
				"%s GetMulticastInterface_IPMreq4Addr interface_addr=%s multicast_addr=%s\n",
				iff.Name, interfaceAddr.String(), multicastAddr.String())
		}

		{
			interfaceIndex, err := ipv4.GetMulticastInterfaceIndex(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%s GetMulticastInterface_Index interface_index=%d\n", iff.Name, interfaceIndex)
		}
	}

	log.Println("ran")
}

func TestUDPPeerIPv4_SetOutboundInterfaceOnUnspecifiedPort(t *testing.T) {
	if len(testInterfacesIPv4) == 0 {
		return
	}

	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	for _, iff := range testInterfacesIPv4 {
		fmt.Println("setting", iff.Name, "as outbound")

		if err := peer.SetOutboundIPv4(iff.Name); err != nil {
			t.Fatal(err)
		}

		outboundInterface, outboundIP := peer.Outbound()
		fmt.Println("outbound for", iff.Name, outboundInterface, outboundIP)

		{
			addr, err := ipv4.GetMulticastInterfaceAddr(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%s GetMulticastInterface_Inet4Addr addr=%s\n", iff.Name, addr.String())
		}

		{
			interfaceAddr, multicastAddr, err := ipv4.GetMulticastInterfaceAddrAndGroup(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf(
				"%s GetMulticastInterface_IPMreq4Addr interface_addr=%s multicast_addr=%s\n",
				iff.Name, interfaceAddr.String(), multicastAddr.String())
		}

		{
			interfaceIndex, err := ipv4.GetMulticastInterfaceIndex(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%s GetMulticastInterface_Index interface_index=%d\n", iff.Name, interfaceIndex)
		}
	}

	log.Println("ran")
}

func TestUDPPeerIPv4_TTL(t *testing.T) {
	if len(testInterfacesIPv4) == 0 {
		return
	}

	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "")
	if err != nil {
		t.Fatal(err)
	}

	actualTTL, err := ipv4.GetMulticastTTL(peer.NextLayer())
	if err != nil {
		t.Fatal(err)
	}

	if actualTTL != peer.TTL() {
		t.Fatalf("wrong TTL expected=%d given=%d", actualTTL, peer.TTL())
	}

	if peer.TTL() != 1 {
		t.Fatalf("peer TTL should be 1 by default")
	}

	setAndCheck := func(ttl uint8) {
		if err := peer.SetTTL(ttl); err != nil {
			t.Fatal(err)
		}

		if peer.TTL() != ttl {
			t.Fatalf("peer TTL should be %d", ttl)
		}
	}

	for i := 0; i <= 255; i++ {
		setAndCheck(uint8(i))
	}

	log.Println("ran")
}

func TestUDPPeerIPv4_Reader1(t *testing.T) {
	// 1 reader, 1 writer

	r := newTestRW(t, "udp", "")
	defer r.Close()

	multicastIP := "224.0.1.0"
	multicastPort := r.peer.LocalAddr().Port
	multicastAddr := fmt.Sprintf("%s:%d", multicastIP, multicastPort)
	if err := r.peer.Join(multicastIP); err != nil {
		t.Fatal(err)
	}

	w := newTestRW(t, "udp", "")
	defer w.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	start := time.Now()
	go func() {
		defer wg.Done()

		r.ReadLoop(func(err error, seq uint64, from netip.AddrPort) {
			if err != nil {
				t.Fatal(err)
			} else {
				if seq == 10 || time.Now().Sub(start).Seconds() > 1 /* just to not have it hang */ {
					r.Close()
				}
			}
		})
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			if err := w.WriteNext(multicastAddr); err != nil {
				t.Fatal(err)
			}
			time.Sleep(time.Millisecond)
		}
	}()

	wg.Wait()

	if len(r.ReceivedFrom()) != 1 {
		t.Fatal("should have received from exactly one source")
	}

	fmt.Println(r.received)
}

//func TestUDPPeerIPv4_Reader2(t *testing.T) {
//	// 1 reader, 2 writers
//
//	r := newTestRW(t, "udp", "")
//	defer r.Close()
//
//	multicastIP := "224.0.1.0"
//	multicastPort := r.peer.LocalAddr().Port
//	multicastAddr := fmt.Sprintf("%s:%d", multicastIP, multicastPort)
//	if err := r.peer.Join(multicastIP); err != nil {
//		t.Fatal(err)
//	}
//
//	w1 := newTestRW(t, "udp", "")
//	defer w1.Close()
//
//	var wg sync.WaitGroup
//	wg.Add(2)
//
//	start := time.Now()
//	go func() {
//		defer wg.Done()
//
//		r.ReadLoop(func(err error, seq uint64, from netip.AddrPort) {
//			if err != nil {
//				t.Fatal(err)
//			} else {
//				if seq == 10 || time.Now().Sub(start).Seconds() > 1 /* just to not have it hang */ {
//					r.Close()
//				}
//			}
//		})
//	}()
//
//	go func() {
//		defer wg.Done()
//		for i := 0; i < 10; i++ {
//			if err := w.WriteNext(multicastAddr); err != nil {
//				t.Fatal(err)
//			}
//			time.Sleep(time.Millisecond)
//		}
//	}()
//
//	wg.Wait()
//
//	if len(r.ReceivedFrom()) != 1 {
//		t.Fatal("should have received from exactly one source")
//	}
//}