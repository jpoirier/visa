// MXA Spectrum Analzer driver

package mxa

import (
	"fmt"
	"strconv"
	"strings"

	vi "github.com/jpoirier/visa"
)

type Mxa struct {
	vi.Visa
	instr vi.Object
}

// SetScreenTitle
func (m *Mxa) SetScreenTitle(title string) (status vi.Status) {
	b := fmt.Sprintf("DISP:ANN:TITL:DATA '%s'", title)
	_, status = m.instr.Write([]byte(b), uint32(len(b)))
	return
}

func (m *Mxa) SaveScreenShot(name string) (status vi.Status) {
	b := fmt.Sprintf("MMEM:STOR:SCR '%s'", name)
	_, status = m.instr.Write([]byte(b), uint32(len(b)))
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
func (m *Mxa) DeleteFile(name string) (status vi.Status) {
	b := fmt.Sprintf("MMEM:DEL '%s'", name)
	_, status = m.instr.Write([]byte(b), uint32(len(b)))
	return
}

// CreateFolder Creates folder name; name must include the full path.
func (m *Mxa) CreateFolder(name string) (status vi.Status) {
	b := fmt.Sprintf("MMEM:MDIR '%s'", name)
	_, status = m.instr.Write([]byte(b), uint32(len(b)))
	return
}

// SetTraceType Sets trace number to trace type.
func (m *Mxa) SetTraceType(number int, stype string) (status vi.Status) {
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
	_, status = m.instr.Write([]byte(b), uint32(len(b)))
	return
}

// SetTraceClearWrite Sets trace number to Clear Write.
func (m *Mxa) SetTraceClearWrite(number int) (status vi.Status) {
	b := fmt.Sprintf("TRAC%d:TYPE WRITE", number)
	_, status = m.instr.Write([]byte(b), uint32(len(b)))
	return
}

// ClearTrace Clears/Resets trace number.
func (m *Mxa) ClearTrace(number int) (status vi.Status) {
	b := fmt.Sprintf("TRAC:CLEAR TRACE%d", number)
	_, status = m.instr.Write([]byte(b), uint32(len(b)))
	return
}

// ClearAllTraces Clears/Resets all traces.
func (m *Mxa) ClearAllTraces(number int) (status vi.Status) {
	b := []byte("TRAC:CLEAR:ALL")
	_, status = m.instr.Write(b, uint32(len(b)))
	return
}

// SetCenterFreqKHz Sets the center crequency freqKhz.
func (m *Mxa) SetCenterFreqKHz(freqKhz float32) (status vi.Status) {
	b := fmt.Sprintf("FREQ:CENT %f KHZ", freqKhz)
	_, status = m.instr.Write([]byte(b), uint32(len(b)))
	return
}

// SetCenterFreqMHz Sets the center crequency freqMhz.
func (m *Mxa) SetCenterFreqMHz(freqMhz float32) (status vi.Status) {
	b := fmt.Sprintf("FREQ:CENT %f MHZ", freqMhz)
	_, status = m.instr.Write([]byte(b), uint32(len(b)))
	return
}

// SetCenterFreqGHz Sets the center crequency freqGhz.
func (m *Mxa) SetCenterFreqGHz(freqGhz float32) (status vi.Status) {
	b := fmt.Sprintf("FREQ:CENT %f GHZ", freqGhz)
	_, status = m.instr.Write([]byte(b), uint32(len(b)))
	return
}

// GetCenterFreqMHz returns the center frequency (MHz).
func (m *Mxa) GetCenterFreqMHz() (freqMhz float32, status vi.Status) {
	b := []byte("FREQ:CENT?")
	_, status = m.instr.Write(b, uint32(len(b)))
	if status < vi.SUCCESS {
		return
	}
	buffer, retCount, status := m.instr.Read(50)
	if status < vi.SUCCESS && retCount > 0 {
		return
	}
	t, err := strconv.ParseFloat(string(buffer), 32)
	if err != nil {
		return freqMhz, -1
	}
	freqMhz = t / 1000.0 / 1000.0
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

//   def setMarkerMode_Normal(self, marker_num):
//     """
//       Activates the specified marker in Normal mode.
//     """
//     self.write("CALC:MARK%d:MODE POS" % (marker_num) )

//   def setMarkerMode_Delta(self, marker_num, relative_marker_num):
//     """
//       Activates the specified marker in Delta mode.
//       Sets the "Relative To" marker property.
//       Note: If relative marker is previously OFF, it will be enabled in Fixed mode at the same location as marker_num.
//     """
//     self.write("CALC:MARK%d:MODE DELT" % (marker_num) )
//     self.write("CALC:MARK%d:REF %d" % (marker_num, relative_marker_num) )

//   def setMarkerMode_Fixed(self, marker_num):
//     """
//       Activates the specified marker in Fixed mode.
//     """
//     self.write("CALC:MARK%d:MODE FIX" % (marker_num) )

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

//   def setMarkerFunction_Noise(self, marker_num):
//     """
//       Sets the Marker Function to Marker Noise for the specified marker.
//     """
//     self.write("CALC:MARK%d:FUNC NOISE" % (marker_num) )

//   def setMarkerFunction_BandPower(self, marker_num):
//     """
//       Sets the Marker Function to Band/Interval Power for the specified marker.
//     """
//     self.write("CALC:MARK%d:FUNC BPOW" % (marker_num) )

//   def setMarkerFunction_BandDensity(self, marker_num):
//     """
//       Sets the Marker Function to Band/Interval Density for the specified marker.
//     """
//     self.write("CALC:MARK%d:FUNC BDEN" % (marker_num) )

//   def setMarkerTraceNum(self, marker_num, trace_num):
//     """
//       Sets the Marker Trace Number for the specified marker.
//     """
//     self.write("CALC:MARK%d:TRAC %d" % (marker_num, trace_num) )

//   def setMarkerLines_On(self, marker_num):
//     """
//       Sets the Lines ON for the specified marker.
//     """
//     self.write("CALC:MARK%d:LINES ON" % (marker_num) )

//   def setMarkerLines_Off(self, marker_num):
//     """
//       Sets the Lines OFF for the specified marker.
//     """
//     self.write("CALC:MARK%d:LINES OFF" % (marker_num) )

//   def setMarkerFunction_Off(self, marker_num):
//     """
//       Sets the Marker Function to Off for the specified marker.
//     """
//     self.write("CALC:MARK%d:FUNC OFF" % (marker_num) )

//   def setMarkerFunction_BandSpan_MHz(self, marker_num, span):
//     """
//       Sets the Marker Function to Band Adjust Span for the specified marker.
//     """
//     self.write("CALC:MARK%d:FUNC:BAND:SPAN %f MHZ" % (marker_num, span) )

//   def setMarkerMode_Off(self, marker_num):
//     """
//       Sets specified marker OFF.
//     """
//     self.write("CALC:MARK%d:MODE OFF" % (marker_num) )

//   def setAllMarkers_Off(self):
//     """
//       Sets specified marker OFF.
//     """
//     self.write("CALC:MARK:AOFF")

//   def setMarker_X_Value_MHz(self, marker_num, x_value):
//     """
//       Sets the X-Axis value for the specified marker in MHz
//     """
//     self.write("CALC:MARK%d:X %f MHZ" % (marker_num, x_value) )

//   def getMarker_X_Value_MHz(self, marker_num):
//     """
//       Gets the X-Axis value for the specified marker.
//       Units are MHz.
//     """
//     return float(self.read("CALC:MARK%d:X?" % (marker_num) ) ) / 1000.0 / 1000.0

//   def setMarker_Y_Value(self, marker_num, y_value):
//     """
//       Sets the Y-Axis value for the specified marker.
//       Units are assumed to be dBm.
//       Note: Can only set Y Value for Fixed type marker.
//     """
//     self.write("CALC:MARK%d:Y %f" % (marker_num, y_value) )

//   def getMarker_Y_Value(self, marker_num):
//     """
//       Gets the Y-Axis value for the specified marker.
//       Units are typically dBm.
//     """
//     return float(self.read("CALC:MARK%d:Y?" % (marker_num) ) )

//   def setMarkerPeakSearch(self, marker_num):
//     """
//       Sets the specified marker to the Peak (max) level.
//     """
//     self.write("CALC:MARK%d:MAX" % (marker_num) )

//   def setMarkerNextPeak(self, marker_num):
//     """
//       Sets the specified marker to the Next Peak
//       If Peak Search has not been performed, this will be the same as Peak Search.
//     """
//     self.write("CALC:MARK%d:MAX:NEXT" % (marker_num) )

//   def setMarkerNextPeak_Right(self, marker_num):
//     """
//       Sets the specified marker to the Next Peak Right.
//       If Peak Search has not been performed, this will be the same as Peak Search.
//     """
//     self.write("CALC:MARK%d:MAX:RIGHT" % (marker_num) )

//   def setMarkerNextPeak_Left(self, marker_num):
//     """
//       Sets the specified marker to the Next Peak Left.
//       If Peak Search has not been performed, this will be the same as Peak Search.
//     """
//     self.write("CALC:MARK%d:MAX:LEFT" % (marker_num) )

//   def setMarkerContinuousPeak_On(self, marker_num):
//     """
//       Sets the specified marker to Continuous Peak Search.
//     """
//     self.write("CALC:MARK%d:CPSEARCH: ON" % (marker_num) )

//   def setMarkerContinuousPeak_Off(self, marker_num):
//     """
//       Disables Continuous Peak Search for specified marker.
//     """
//     self.write("CALC:MARK%d:CPSEARCH: OFF" % (marker_num) )

//   def setMarkerTable_On(self):
//     """
//       Turns the Marker Table Display ON.
//     """
//     self.write("CALC:MARK:TABL ON")

//   def setMarkerTable_Off(self):
//     """
//       Turns the Marker Table Display OFF.
//     """
//     self.write("CALC:MARK:TABL OFF")

//   def setPeakTable_On(self):
//     """
//       Turns the Peak Table Display ON.
//     """
//     self.write("CALC:MARK:PEAK:TABL:STATE ON")

//   def setPeakTable_Off(self):
//     """
//       Turns the Peak Table Display OFF.
//     """
//     self.write("CALC:MARK:PEAK:TABL:STATE OFF")

//   def saveMarkerTable(self, strFilename):
//     """
//       Saves Marker Table to MXA drive.
//       Always activates marker table view first.
//     """
//     self.setMarkerTable_On()
//     self.write("MMEM:STOR:RES:MTAB '%s'" % (strFilename) )

//   def savePeakTable(self, strFilename):
//     """
//       Saves Peak Table to MXA drive.
//       Always activates peak table view first.
//     """
//     self.setPeakTable_On()
//     self.write("MMEM:STOR:RES:PTAB '%s'" % (strFilename) )

//   def saveSpectogram(self, strFilename):
//     """
//       Saves Spectogram.
//     """
//     self.write("MMEM:STOR:RES:SPEC '%s'" % (strFilename) )

//   # TBD - Resize Marker Table
//   #       Does not seem possible via remote commands.

//   def showLTE_ACP(self):
//     """
//       Sets MXA to LTE mode, ACP Measurement screen.
//     """
//     self.write("INST LTE")
//     self.write("CONF:ACP")

//   def showSpectrumAnalyzer(self):
//     """
//       Sets MXA to Spectrum Analyzer mode.
//     """
//     self.write("INST SA")

//   def setRefLevel(self, refLevel):
//     """
//       Sets Reference Level in dBm.
//       Note: MXA will convert to max allowed if input is too high.
//     """
//     self.write("DISP:WIND:TRAC:Y:RLEV %f DBM" % (refLevel) )

//   def setRefLevelOffset(self, refLevelOffset):
//     """
//       Sets Reference Level Offset in dB.
//     """
//     self.write("DISP:WIND:TRAC:Y:RLEV:OFFSET %f" % (refLevelOffset) )

// if __name__ == "__main__":

//   TestMXA = MXA(GPIB_Address = 18)

//   print 'IDN: ' + TestMXA.GetID()
