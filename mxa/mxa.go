// Agilent MXA Spectrum Analyzer Driver

package mxa

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	vi "github.com/jpoirier/visa"
)

type Driver struct {
	vi.Driver
}

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

// SetScreenTitle Sets screen title.
func (d *Driver) SetScreenTitle(title string) (status vi.Status) {
	b := fmt.Sprintf("DISP:ANN:TITL:DATA '%s'", title)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SaveScreenShot Saves name screenshot.
func (d *Driver) SaveScreenShot(name string) (status vi.Status) {
	b := fmt.Sprintf("MMEM:STOR:SCR '%s'", name)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// def getFile(self, strSrcFilename, strDstFilename):
//   """
//     Get file contents as raw binary data and save to file on local pc.
//     Usually used to fetch a screenshot or trace.
//   """
//   strSrcFilename = strSrcFilename.upper()  # Convert filename to all uppercase per GPIB standard

//   # Agilent Command to read data file
//   self.write("MMEM:DATA? '%s'" % strSrcFilename)
//   strData = self.Instrument.read_raw()

//   # data is stored as a packed string in strData. Entire file is stored here, limited only by PC RAM.
//   if strData[0] == '#':  # must start with the # sign
//     Length_of_LengthBytes = int(strData[1],10)  # Second byte tells how many bytes it takes to tell us the full file length
//     dataLen = int(strData[2:2+Length_of_LengthBytes],10)  # Read bytes to tell how long data is
//     #print 'File Length is: %d' % dataLen

//     outfile = open(strDstFilename,'wb')

//     for i in range(2+Length_of_LengthBytes,len(strData)-1):
//       outfile.write(strData[i])
//       outfile.flush()
//       outfile.close

// DeleteFile Deletes file name; name must include the full path.
func (d *Driver) DeleteFile(name string) (status vi.Status) {
	b := fmt.Sprintf("MMEM:DEL '%s'", name)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// CreateFolder Creates folder name; name must include the full path.
func (d *Driver) CreateFolder(name string) (status vi.Status) {
	b := fmt.Sprintf("MMEM:MDIR '%s'", name)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetTraceType Sets trace number to trace stype.
func (d *Driver) SetTraceType(number int, stype string) (status vi.Status) {
	stype = strings.ToUpper(stype)
	var b string
	switch stype {
	default:
		return -1
	case "CLEAR", "WRITE":
		b = fmt.Sprintf("TRAC%d:TYPE WRITE", number)
	case "AVERAGE":
		b = fmt.Sprintf("TRAC%d:TYPE AVERAGE", number)
	case "MAX", "MAXH", "MAXHOLD":
		b = fmt.Sprintf("TRAC%d:TYPE MAXHOLD", number)
	case "MIN", "MINH", "MINHOLD":
		b = fmt.Sprintf("TRAC%d:TYPE MINHOLD", number)
	}
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetTraceClearWrite Sets number trace to clear write.
func (d *Driver) SetTraceClearWrite(trace uint32) (status vi.Status) {
	b := fmt.Sprintf("TRAC%d:TYPE WRITE", trace)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// ClearTrace Clears/Resets number trace.
func (d *Driver) ClearTrace(trace uint32) (status vi.Status) {
	b := fmt.Sprintf("TRAC:CLEAR TRACE%d", trace)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// ClearAllTraces Clears/Resets all traces.
func (d *Driver) ClearAllTraces() (status vi.Status) {
	b := []byte("TRAC:CLEAR:ALL")
	_, status = d.Write(b, uint32(len(b)))
	return
}

// SetCenterFreqKHz Sets the center crequency mhz.
func (d *Driver) SetCenterFreqKHz(mhz float32) (status vi.Status) {
	b := fmt.Sprintf("FREQ:CENT %f KHZ", mhz)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetCenterFreqMHz Sets the center crequency mhz.
func (d *Driver) SetCenterFreqMHz(mhz float32) (status vi.Status) {
	b := fmt.Sprintf("FREQ:CENT %f MHZ", mhz)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetCenterFreqGHz Sets the center crequency ghz.
func (d *Driver) SetCenterFreqGHz(ghz float32) (status vi.Status) {
	b := fmt.Sprintf("FREQ:CENT %f GHZ", ghz)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// GetCenterFreqMHz returns the center frequency mhz).
func (d *Driver) GetCenterFreqMHz() (mhz float32, status vi.Status) {
	b := []byte("FREQ:CENT?")
	_, status = d.Write(b, uint32(len(b)))
	if status < vi.SUCCESS {
		return
	}
	buffer, retCount, status := d.Read(50)
	if status < vi.SUCCESS && retCount > 0 {
		return
	}
	t, err := strconv.ParseFloat(string(buffer), 32)
	if err != nil {
		return mhz, -1
	}
	mhz = float32(t / 1000.0 / 1000.0)
	return
}

// TBD - setFreqSpan seems to only work for Spectrum Analyzer.  How do I set for LTE ACP Measurement?

//   def setFreqSpan_MHz(self, span):
//     """
//       Sets the Spectrum Analzer Frequency Span.  units = MHz
//     """
//     self.write("FREQ:SPAN %f MHZ" % (span) )

//   def setFreqSpan_kHz(self, span):
//     """
//       Sets the Spectrum Analzer Frequency Span.  units = kHz
//     """
//     self.write("FREQ:SPAN %f KHZ" % (span) )

//   def setFreqSpan_GHz(self, span):
//     """
//       Sets the Spectrum Analzer Frequency Span.  units = GHz
//     """
//     self.write("FREQ:SPAN %f GHZ" % (span) )

//   def getFreqSpan_MHz(self):
//     """
//       Gets the Spectrum Analzer Frequency Span.  units = MHz
//     """
//     return float(self.read("FREQ:SPAN?")) / 1000.0 / 1000.0

//   def setACP_FreqSpan_MHz(self, span):
//     """
//       Sets the ACP Measurement Frequency Span.  units = MHz
//     """
//     self.write("ACP:FREQ:SPAN %f MHZ" % (span) )

//   def setMarkerMode(self, marker_num = 1, mode = "NORMAL", delta_marker = 1):
//     """
//       Sets the specified marker type for the specified marker number.
//     """
//     if mode.upper() == "NORMAL" or mode.upper() == "POS":
//         self.write("CALC:MARK%d:MODE POS" % (marker_num) )
//     elif mode.upper() == "DELTA":
//         self.write("CALC:MARK%d:MODE DELT" % (marker_num) )
//         self.write("CALC:MARK%d:REF %d" % (marker_num, delta_marker) )
//     elif mode.upper() == "FIXED":
//         self.write("CALC:MARK%d:MODE FIX" % (marker_num) )

// SetMarkerModeNorm Puts the specified marker in normal mode.
func (d *Driver) SetMarkerModeNorm(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:MODE POS", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerModeDelta Activates marker in Delta mode and sets relMarker
// property.  Note, if relative marker was OFF, it's enabled in Fixed mode
// at the same location as marker.
func (d *Driver) SetMarkerModeDelta(marker, relMarker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:MODE DELT", marker)

	_, status = d.Write([]byte(b), uint32(len(b)))
	if status < vi.SUCCESS {
		return
	}
	b = fmt.Sprintf("CALC:MARK%d:REF %d", relMarker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerModeFixed Puts the specified marker in fixed mode.
func (d *Driver) SetMarkerModeFixed(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:MODE FIX", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

//   def setMarkerFunction(self, marker_num = 1, fxn = "BAND_POWER"):
//     """
//       Sets the specified marker function for the specified marker number.
//     """

//     if fxn.upper() == "BAND_POWER" or fxn.upper() == "POWER" or fxn.upper() == "POW" or fxn.upper() == "BPOW":
//         self.write("CALC:MARK%d:FUNC BPOW" % (marker_num) )
//     elif fxn.upper() == "BAND_DENSITY" or fxn.upper() == "DENSITY" or fxn.upper() == "DEN" or fxn.upper() == "BDEN":
//         self.write("CALC:MARK%d:FUNC BDEN" % (marker_num) )
//     elif fxn.upper() == "NOISE":
//         self.write("CALC:MARK%d:FUNC NOISE" % (marker_num) )
//     elif fxn.upper() == "OFF":
//         self.write("CALC:MARK%d:FUNC OFF" % (marker_num) )

// SetMarkerFuncNoise Sets marker function to noise.
func (d *Driver) SetMarkerFuncNoise(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:FUNC NOISE", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerFuncBandPower Sets marker function to band/interval power.
func (d *Driver) SetMarkerFuncBandPower(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:FUNC BPOW", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerFuncBandDensity Sets marker band/interval density.
func (d *Driver) SetMarkerFuncBandDensity(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:FUNC BDEN", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerTraceNum Sets marker trace number.
func (d *Driver) SetMarkerTraceNum(marker, number uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:TRAC %d", marker, number)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerLinesOn Sets marker lines on.
func (d *Driver) SetMarkerLinesOn(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:LINES ON", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerLinesOff Sets marker lines off.
func (d *Driver) SetMarkerLinesOff(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:LINES OFF", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerFuncOff Sets marker function off.
func (d *Driver) SetMarkerFuncOff(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:FUNC OFF", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerFuncBandSpanMHz Sets marker to band adjust span mhz.
func (d *Driver) SetMarkerFuncBandSpanMHz(marker, mhz uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:FUNC:BAND:SPAN %f MHZ", marker, mhz)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerOff Sets marker off.
func (d *Driver) SetMarkerOff(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:MODE OFF", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetAllMarkersOff Sets all markers off.
func (d *Driver) SetAllMarkersOff() (status vi.Status) {
	b := []byte]("CALC:MARK:AOFF")
	_, status = d.Write(b, uint32(len(b)))
	return
}

// SetMarkerXValMHz Sets marker X-Axis value to mhz.
func (d *Driver) SetMarkerXValMHz(marker, mhz uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:X %f MHZ", marker, mhz)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

//   def getMarker_X_Value_MHz(self, marker_num):
//     """
//       Gets the X-Axis value for the specified marker.
//       Units are MHz.
//     """
//     return float(self.read("CALC:MARK%d:X?" % (marker_num) ) ) / 1000.0 / 1000.0

// SetMarkerYValDbm Sets marker Y-Axis value to dbm.
// Note, fixed type marker only.
func (d *Driver) SetMarkerYValDbm(marker uint32, dbm float32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:Y %f", marker, dbm)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

//   def getMarker_Y_Value(self, marker_num):
//     """
//       Gets the Y-Axis value for the specified marker.
//       Units are typically dBm.
//     """
//     return float(self.read("CALC:MARK%d:Y?" % (marker_num) ) )

// SetMarkerPeakSearch Sets marker peak (max) search level on.
func (d *Driver) SetMarkerPeakSearch(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:MAX", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerNextPeak Sets marker next peak on.
// Note, peak search performed if not done previously.
func (d *Driver) SetMarkerNextPeak(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:MAX:NEXT", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerNextPeakR Sets marker next right peak on.
// Note, peak search performed if not done previously.
func (d *Driver) SetMarkerNextPeakR(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:MAX:RIGHT", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerNextPeakL Sets marker next left peak on.
// Note, peak search performed if not done previously.
func (d *Driver) SetMarkerNextPeakL(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:MAX:LEFT", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerContPeak Sets marker continuous peak search on.
func (d *Driver) SetMarkerContPeak(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:CPSEARCH: ON", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerContPeakOff Sets marker continuous peak search off.
func (d *Driver) SetMarkerContPeak(marker uint32) (status vi.Status) {
	b := fmt.Sprintf("CALC:MARK%d:CPSEARCH: OFF", marker)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetMarkerTableOn Sets the market table display off.
func (d *Driver) SetMarkerTableOn() (status vi.Status) {
	b := []byte]("CALC:MARK:TABL ON")
	_, status = d.Write(b, uint32(len(b)))
	return
}

// SetMarkerTableOff Sets the market table display off.
func (d *Driver) SetMarkerTableOff() (status vi.Status) {
	b := []byte]("CALC:MARK:TABL OFF")
	_, status = d.Write(b, uint32(len(b)))
	return
}

// SetPeakTableOn Sets the peak table display on.
func (d *Driver) SetMarkerTableOn() (status vi.Status) {
	b := []byte]("CALC:MARK:PEAK:TABL:STATE ON")
	_, status = d.Write(b, uint32(len(b)))
	return
}

// SetPeakTableOff Sets the peak table display to off.
func (d *Driver) SetMarkerTableOff() (status vi.Status) {
	b := []byte]("CALC:MARK:PEAK:TABL:STATE OFF")
	_, status = d.Write(b, uint32(len(b)))
	return
}

// SaveMarkerTable Saves the marker table to filename.
func (d *Driver) SaveMarkerTable(filename string) (status vi.Status) {
	status = d.SetMarkerTableOn()
	if status < vi.SUCCESS {
		return
	}
	b := fmt.Sprintf("MMEM:STOR:RES:MTAB '%s'", filename)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SavePeakTable Saves the peak table to filename.
func (d *Driver) SavePeakTable(filename string) (status vi.Status) {
	status = d.SetPeakTableOn()
	if status < vi.SUCCESS {
		return
	}
	b := fmt.Sprintf("MMEM:STOR:RES:PTAB '%s'", filename)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SaveSpectogram Saves the spectogram to filename.
func (d *Driver) SavePeakTable(filename string) (status vi.Status) {
	b := fmt.Sprintf("MMEM:STOR:RES:SPEC '%s'", filename)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// ShowLTE_ACP Sets LTE mode and ACP measurement screen on.
func (d *Driver) ShowLTE_ACP() (status vi.Status) {
	b := []byte("INST LTE")
	_, status = d.Write(b, uint32(len(b)))
	if status < vi.SUCCESS {
		return
	}
	b = []byte("CONF:ACP")
	_, status = d.Write(b, uint32(len(b)))
	return
}

//   # TBD - Resize Marker Table
//   #       Does not seem possible via remote commands.

// ShowSpectrumAnalyzer Sets spectrum analyzer mode on.
func (d *Driver) ShowSpectrumAnalyzer() (status vi.Status) {
	b := []byte]("INST SA")
	_, status = d.Write(b, uint32(len(b)))
	return
}

// SetRefLevel Sets the reference level to dbm.
// Note, if too high then set to max allowed.
func (d *Driver) SetRefLevel(dbm float32) (status vi.Status) {
	b := fmt.Sprintf("DISP:WIND:TRAC:Y:RLEV %f DBM", dbm)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}

// SetRefLevelOffset Sets reference level offset to dbm.
func (d *Driver) SetRefLevelOffset(dbm float32) (status vi.Status) {
	b := fmt.Sprintf("DISP:WIND:TRAC:Y:RLEV:OFFSET %f", dbm)
	_, status = d.Write([]byte(b), uint32(len(b)))
	return
}
