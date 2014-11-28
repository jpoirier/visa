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
	"os"

	vi "github.com/jpoirier/visa"
)

func main() {
	rm, status := vi.OpenDefaultRM()
	if status < vi.SUCCESS {
		fmt.Println("Could not open a session to the VISA Resource Manager!")
		os.Exit(0)
	}

	instr, status := rm.Open("TCPIP0::ftp.ni.com::21::SOCKET", vi.NULL, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Println("An error occurred opening the session to TCPIP0::ftp.ni.com::21::SOCKET")
		rm.Close()
		os.Exit(0)
	}

	status = instr.SetAttribute(vi.ATTR_TCPIP_NODELAY, vi.TRUE)
	if status < vi.SUCCESS {
		fmt.Println("An error occurred setting the attributes...")
		instr.Close()
		rm.Close()
		os.Exit(0)
	}
	buf, _, status := instr.Read(25)
	if status < vi.SUCCESS {
		fmt.Printf("Read failed with error code %x \n", status)
		rm.Close()
		os.Exit(0)
	}
	fmt.Printf("The server response is:\n %s\n\n", string(buf))

	buf, _ = instr.GetAttribute(vi.ATTR_TCPIP_ADDR)
	fmt.Printf(" Address:  %s\n", string(buf))

	buf, _ = instr.GetAttribute(vi.ATTR_TCPIP_HOSTNAME)
	fmt.Printf(" Host Name:  %s\n", string(buf))

	buf, _ = instr.GetAttribute(vi.ATTR_TCPIP_PORT)
	fmt.Printf(" Port:  %s\n", string(buf))

	buf, _ = instr.GetAttribute(vi.ATTR_RSRC_CLASS)
	fmt.Printf(" Resource Class:  %s\n", string(buf))

	instr.Close()
	rm.Close()
}
