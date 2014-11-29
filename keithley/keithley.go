// Keithley RF Switch driver

package keithley

import (
	"fmt"
	"os"

	vi "github.com/jpoirier/visa"
)

type Driver struct {
	vi.Driver
}

// keithley - works with Keithley S46 RF Switch.
// Eight SPDT unterminated coaxial relays (2-pole) and four multi-pole
// unterminated coaxial relays.
// 	Relay A, multipole = Chan  1...6
// 	Relay B, multipole = Chan  7..12
// 	Relay C, multipole = Chan 13..18
// 	Relay D, multipole = Chan 19..24
//	Relay 1..8, 2-pole = Chan 25..32
//
// Caution: Do not close more than one RF path per multiport switch.

// OpenGpib Opens a session to the specified resource.
func OpenGpib(rm vi.Session, ctrl, addr, mode, timeout uint32) (*Driver, vi.Status) {
	name := fmt.Sprintf("GPIB%d::%d", ctrl, addr)
	instr, status := rm.Open(name, mode, timeout)
	if status < vi.SUCCESS {
		fmt.Println("Error, OpenGpib failed with error: ", status)
		os.Exit(0)
	}
	return &Driver{instr}, status
}

// OpenTcp Opens a session to the specified resource.
func OpenTcp(rm vi.Session, ip string, mode, timeout uint32) (*Driver, vi.Status) {
	if len(ip) == 0 {
		fmt.Println("Error, empty ip address string.")
		os.Exit(0)
	}
	name := fmt.Sprintf("TCPIP::%s::INSTR", ip)
	instr, status := rm.Open(name, mode, timeout)
	if status < vi.SUCCESS {
		fmt.Println("Error, OpenGpib failed with error: ", status)
		os.Exit(0)
	}
	return &Driver{instr}, status
}

// Reset Resets the switch unit.
func (d *Driver) Reset() (status vi.Status) {
	b := []byte("*RST")
	_, status = d.Write(b, uint32(len(b)))
	return
}

// OpenChan Opens the specified channel. Where an open channel does not
// allow a signal to pass through.
func (d *Driver) OpenChan(ch uint32) (status vi.Status) {
	b := fmt.Sprintf("OPEN (@%d)", ch)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// OpenAllChans Opens all channels.
func (d *Driver) OpenAllChans() (status vi.Status) {
	b := []byte("OPEN:ALL")
	_, status = d.Write(b, uint32(len(b)))
	return
}

// CloseChan Closes the specified channel. Note, All other channels on relay
// are opened first to prevent multiple closed relays leading to damage.
func (d *Driver) CloseChan(ch uint32) (status vi.Status) {
	// Determine if ch is part of 2-port relay or multi-port relay (A..D)
	// Multi-Port Relay, A..D
	if ch > 0 && ch < 25 {
		// Open all ports on this relay
		for i := 1; i < 7; i++ {
			c := int((ch-1)/6)*6 + i
			d.OpenChan(uint32(c))
		}
	}
	b := fmt.Sprintf("CLOSE (@%d)", ch)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// ClosedChanList Returns a list of closed channels.
func (d *Driver) ClosedChanList(ch uint32) (list string, status vi.Status) {
	// Returns list of closed channels.
	// RF Switch returns format '(@1,2,3)'.
	// If no channels closed, switch returns '(@)'.

	b := []byte("CLOSE?")
	_, status = d.Write(b, uint32(len(b)))
	if status < vi.SUCCESS {
		return
	}
	buffer, _, status := d.Read(100)
	if status < vi.SUCCESS {
		return
	} else {
		return string(buffer), status
	}
}

// if len(strClosedChans) > 3:  # Len always greater than 3 if channel in the list.
//     if strClosedChans.find(",") >= 0:
//         closedChans = strClosedChans[2:len(strClosedChans)-1].split(",")
//     else:
//         closedChans.append(strClosedChans[2:len(strClosedChans)-1])
