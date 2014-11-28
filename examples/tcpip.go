//
// This example demonstrates opening a simple TCPIP connection and
// does a read and checks a few properties.
//
// The general flow of the code is
//      Open Resource Manager
//      Open a session to the TCP/IP site at NI
//      Perform a read, and check properties
//      Close all VISA Sessions
//

package main

import (
	"fmt"
	"unsafe"

	vi "github.com/jpoirier/visa"
)

func main() {
	rm, status := vi.OpenDefaultRM()
	if status < vi.SUCCESS {
		fmt.Println("Could not open a session to the VISA Resource Manager!")
		return
	}
	defer rm.Close()

	instr, status := rm.Open("TCPIP0::ftp.ni.com::21::SOCKET", vi.NULL, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Println("An error occurred opening the session to TCPIP0::ftp.ni.com::21::SOCKET")
		return
	}
	defer instr.Close()

	status = instr.SetAttribute(vi.ATTR_TCPIP_NODELAY, vi.TRUE)
	if status < vi.SUCCESS {
		fmt.Println("An error occurred setting the attributes...")
		return
	}
	b, _, status := instr.Read(25)
	if status < vi.SUCCESS {
		fmt.Printf("Read failed with error code %x \n", status)
		return
	}
	fmt.Printf("The server response is:\n %s\n\n", string(b))

	buf := make([]byte, 100)
	instr.GetAttribute(vi.ATTR_TCPIP_ADDR, unsafe.Pointer(&buf[0]))
	fmt.Printf(" Address:  %s\n", string(buf))

	buf = nil
	instr.GetAttribute(vi.ATTR_TCPIP_HOSTNAME, unsafe.Pointer(&buf[0]))
	fmt.Printf(" Host Name:  %s\n", string(buf))

	buf = nil
	instr.GetAttribute(vi.ATTR_TCPIP_PORT, unsafe.Pointer(&buf[0]))
	fmt.Printf(" Port:  %s\n", string(buf))

	buf = nil
	instr.GetAttribute(vi.ATTR_RSRC_CLASS, unsafe.Pointer(&buf[0]))
	fmt.Printf(" Resource Class:  %s\n", string(buf))
}
