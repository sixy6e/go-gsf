package gsf

import (
	"bytes"
	"encoding/binary"
	"time"
)

// TODO; check signed/unsigned type reads with c code (stemp, ltemp)
// Update: it will be difficult to determine the correct read type in cases
// where the spec differs to the c-code.
// In many instances the spec says int16, but the code reads it as a uint16
// and then converts to an int32.
// If the spec says int16, then we should be safe in reading as int16.
// However, I've come across many instances where the code does something
// very different to the spec.
// Best attempts will be made to infer the "more correct" type if something
// doesn't look right:
// (differs wildly in the spec vs code, as well as what the data represents).

type Seabeam struct {
	EclipseTime []uint16
}

func DecodeSeabeamSpecific(reader *bytes.Reader) (sensor_data Seabeam) {
	var buffer struct {
		EclipseTime uint16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)
	sensor_data.EclipseTime = []uint16{buffer.EclipseTime}

	return sensor_data
}

type Em12 struct {
	PingNumber    []int16
	Resolution    []int8
	PingQuality   []int8
	SoundVelocity []float32
	Mode          []int8
}

func DecodeEm12Specific(reader *bytes.Reader) (sensor_data Em12) {
	var buffer struct {
		PingNumber    int16
		Resolution    int8
		PingQuality   int8
		SoundVelocity int16
		Mode          int8
		Spare         [4]int32
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.Resolution = []int8{buffer.Resolution}
	sensor_data.PingQuality = []int8{buffer.PingQuality}
	sensor_data.SoundVelocity = []float32{float32(buffer.SoundVelocity) / 10.0}
	sensor_data.Mode = []int8{buffer.Mode}

	return sensor_data
}

type Em100 struct {
	ShipPitch       []float32
	TransducerPitch []float32
	Mode            []int8
	Power           []int8
	Attenuation     []int8
	Tvg             []int8
	PulseLength     []int8
	Counter         []int16
}

func DecodeEm100Specific(reader *bytes.Reader) (sensor_data Em100) {
	var buffer struct {
		ShipPitch       int16
		TransducerPitch int16
		Mode            int8
		Power           int8
		Attenuation     int8
		Tvg             int8
		PulseLength     int8
		Counter         int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.ShipPitch = []float32{float32(buffer.ShipPitch) / SCALE2}
	sensor_data.TransducerPitch = []float32{float32(buffer.TransducerPitch) / SCALE2}
	sensor_data.Mode = []int8{buffer.Mode}
	sensor_data.Power = []int8{buffer.Power}
	sensor_data.Attenuation = []int8{buffer.Attenuation}
	sensor_data.Tvg = []int8{buffer.Tvg}
	sensor_data.PulseLength = []int8{buffer.PulseLength}
	sensor_data.Counter = []int16{buffer.Counter}

	return sensor_data
}

type Em950 struct {
	PingNumber           []uint16
	Mode                 []uint8
	Quality              []uint8
	ShipPitch            []float32
	TransducerPitch      []float32
	SurfaceSoundVelocity []float32
}

func DecodeEm950Specific(reader *bytes.Reader) (sensor_data Em950) {
	var buffer struct {
		PingNumber           uint16
		Mode                 uint8
		Quality              uint8
		ShipPitch            int16
		TransducerPitch      int16
		SurfaceSoundVelocity uint16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.Quality = []uint8{buffer.Quality}
	sensor_data.ShipPitch = []float32{float32(buffer.ShipPitch) / SCALE2}
	sensor_data.TransducerPitch = []float32{float32(buffer.TransducerPitch) / SCALE2}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / 10.0}

	return sensor_data
}

type Em121A struct {
	PingNumber           []int16
	Mode                 []int8
	ValidBeams           []int8
	PulseLength          []int8
	BeamWidth            []int8
	TransmitPower        []int8
	TransmitStatus       []int8
	ReceiveStatus        []int8
	SurfaceSoundVelocity []float32
}

func DecodeEm121ASpecific(reader *bytes.Reader) (sensor_data Em121A) {
	var buffer struct {
		PingNumber           int16
		Mode                 int8
		ValidBeams           int8
		PulseLength          int8
		BeamWidth            int8
		TransmitPower        int8
		TransmitStatus       int8
		ReceiveStatus        int8
		SurfaceSoundVelocity int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.Mode = []int8{buffer.Mode}
	sensor_data.ValidBeams = []int8{buffer.ValidBeams}
	sensor_data.PulseLength = []int8{buffer.PulseLength}
	sensor_data.BeamWidth = []int8{buffer.BeamWidth}
	sensor_data.TransmitPower = []int8{buffer.TransmitPower}
	sensor_data.TransmitStatus = []int8{buffer.TransmitStatus}
	sensor_data.ReceiveStatus = []int8{buffer.ReceiveStatus}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.ReceiveStatus) / 10.0}

	return sensor_data
}

type Em121 struct {
	PingNumber           []int16
	Mode                 []int8
	ValidBeams           []int8
	PulseLength          []int8
	BeamWidth            []int8
	TransmitPower        []int8
	TransmitStatus       []int8
	ReceiveStatus        []int8
	SurfaceSoundVelocity []float32
}

func DecodeEm121Specific(reader *bytes.Reader) (sensor_data Em121) {
	var buffer struct {
		PingNumber           int16
		Mode                 int8
		ValidBeams           int8
		PulseLength          int8
		BeamWidth            int8
		TransmitPower        int8
		TransmitStatus       int8
		ReceiveStatus        int8
		SurfaceSoundVelocity int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.Mode = []int8{buffer.Mode}
	sensor_data.ValidBeams = []int8{buffer.ValidBeams}
	sensor_data.PulseLength = []int8{buffer.PulseLength}
	sensor_data.BeamWidth = []int8{buffer.BeamWidth}
	sensor_data.TransmitPower = []int8{buffer.TransmitPower}
	sensor_data.TransmitStatus = []int8{buffer.TransmitStatus}
	sensor_data.ReceiveStatus = []int8{buffer.ReceiveStatus}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.ReceiveStatus) / 10.0}

	return sensor_data
}

type Sass struct {
	LeftMostBeam       []int16
	RightMostBeam      []int16
	TotalNumverOfBeams []int16
	NavigationMode     []int16
	PingNumber         []int16
	MissionNumber      []int16
}

func DecodeSassSpecfic(reader *bytes.Reader) (sensor_data Sass) {
	var buffer struct {
		LeftMostBeam       int16
		RightMostBeam      int16
		TotalNumverOfBeams int16
		NavigationMode     int16
		PingNumber         int16
		MissionNumber      int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.LeftMostBeam = []int16{buffer.LeftMostBeam}
	sensor_data.RightMostBeam = []int16{buffer.RightMostBeam}
	sensor_data.TotalNumverOfBeams = []int16{buffer.TotalNumverOfBeams}
	sensor_data.NavigationMode = []int16{buffer.NavigationMode}
	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.MissionNumber = []int16{buffer.MissionNumber}

	return sensor_data
}

type Seamap struct {
	PortTransmit1        []float32
	PortTransmit2        []float32
	StarboardTransmit1   []float32
	StarboardTransmit2   []float32
	PortGain             []float32
	StarboardGain        []float32
	PortPulseLength      []float32
	StarboardPulseLength []float32
	PressureDepth        []float32 // only present in GSF >= 2.08
	Altitude             []float32
	Temperature          []float32
}

func DecodeSeamapSpecific(reader *bytes.Reader, gsfd GsfDetails) (sensor_data Seamap) {
	var (
		buffer1 struct {
			PortTransmit1        int16
			PortTransmit2        int16
			StarboardTransmit1   int16
			StarboardTransmit2   int16
			PortGain             int16
			StarboardGain        int16
			PortPulseLength      int16
			StarboardPulseLength int16
		}
		pressure_depth int16 // only present in GSF >= 2.08
		buffer2        struct {
			Altitude    int16
			Temperature int16
		}
	)
	_ = binary.Read(reader, binary.BigEndian, &buffer1)

	major, minor := gsfd.MajorMinor()

	if major > 2 || (major == 2 && minor > 7) {
		_ = binary.Read(reader, binary.BigEndian, &pressure_depth)
	} else {
		pressure_depth = 0 // treating as null
	}

	_ = binary.Read(reader, binary.BigEndian, &buffer2)

	sensor_data.PortTransmit1 = []float32{float32(buffer1.PortTransmit1) / 10.0}
	sensor_data.PortTransmit2 = []float32{float32(buffer1.PortTransmit2) / 10.0}
	sensor_data.StarboardTransmit1 = []float32{float32(buffer1.StarboardTransmit1) / 10.0}
	sensor_data.StarboardTransmit2 = []float32{float32(buffer1.StarboardTransmit2) / 10.0}
	sensor_data.PortGain = []float32{float32(buffer1.PortGain) / 10.0}
	sensor_data.StarboardGain = []float32{float32(buffer1.StarboardGain) / 10.0}
	sensor_data.PortPulseLength = []float32{float32(buffer1.PortPulseLength) / 10.0}
	sensor_data.StarboardPulseLength = []float32{float32(buffer1.StarboardPulseLength) / 10.0}
	sensor_data.PressureDepth = []float32{float32(pressure_depth) / 10.0}
	sensor_data.Altitude = []float32{float32(buffer2.Altitude) / 10.0}
	sensor_data.Temperature = []float32{float32(buffer2.Temperature) / 10.0}

	return sensor_data
}

type Seabat struct {
	PingNumber           []int16
	SurfaceSoundVelocity []float32
	Mode                 []int8
	Range                []int8
	TransmitPower        []int8
	ReceiveGain          []int8
}

func DecodeSeabatSpecific(reader *bytes.Reader, gsfd GsfDetails) (sensor_data Seabat) {
	var buffer struct {
		PingNumber           int16
		SurfaceSoundVelocity int16
		Mode                 int8
		Range                int8
		TransmitPower        int8
		ReceiveGain          int8
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / 10.0}
	sensor_data.Mode = []int8{buffer.Mode}
	sensor_data.Range = []int8{buffer.Range}
	sensor_data.TransmitPower = []int8{buffer.TransmitPower}
	sensor_data.ReceiveGain = []int8{buffer.ReceiveGain}

	return sensor_data
}

type Em1000 struct {
	PingNumber           []int16
	Mode                 []int8
	Quality              []int8
	ShipPitch            []float32
	TransducerPitch      []float32
	SurfaceSoundVelocity []float32
}

func DecodeEm1000Specific(reader *bytes.Reader) (sensor_data Em1000) {
	var buffer struct {
		PingNumber           int16
		Mode                 int8
		Quality              int8
		ShipPitch            int16
		TransducerPitch      int16
		SurfaceSoundVelocity int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.Mode = []int8{buffer.Mode}
	sensor_data.Quality = []int8{buffer.Quality}
	sensor_data.ShipPitch = []float32{float32(buffer.ShipPitch) / SCALE2}
	sensor_data.TransducerPitch = []float32{float32(buffer.TransducerPitch) / SCALE2}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / 10.0}

	return sensor_data
}

type TypeIIISeabeam struct {
	LeftMostBeam       []int16
	RightMostBeam      []int16
	TotalNumverOfBeams []int16
	NavigationMode     []int16
	PingNumber         []int16
	MissionNumber      []int16
}

func DecodeTypeIIISeabeamSpecific(reader *bytes.Reader) (sensor_data TypeIIISeabeam) {
	var buffer struct {
		LeftMostBeam       int16
		RightMostBeam      int16
		TotalNumverOfBeams int16
		NavigationMode     int16
		PingNumber         int16
		MissionNumber      int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.LeftMostBeam = []int16{buffer.LeftMostBeam}
	sensor_data.RightMostBeam = []int16{buffer.RightMostBeam}
	sensor_data.TotalNumverOfBeams = []int16{buffer.TotalNumverOfBeams}
	sensor_data.NavigationMode = []int16{buffer.NavigationMode}
	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.MissionNumber = []int16{buffer.MissionNumber}

	return sensor_data
}

type SbAmp struct {
	Hour         []int8
	Minute       []int8
	Second       []int8
	Hundredths   []int8
	BlockNumber  []int32
	AvgGateDepth []int16
}

func DecodeSbAmpSeabeamSpecific(reader *bytes.Reader) (sensor_data SbAmp) {
	var buffer struct {
		Hour         int8
		Minute       int8
		Second       int8
		Hundredths   int8
		BlockNumber  int32
		AvgGateDepth int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.Hour = []int8{buffer.Hour}
	sensor_data.Minute = []int8{buffer.Minute}
	sensor_data.Second = []int8{buffer.Second}
	sensor_data.Hundredths = []int8{buffer.Hundredths}
	sensor_data.BlockNumber = []int32{buffer.BlockNumber}
	sensor_data.AvgGateDepth = []int16{buffer.AvgGateDepth}

	return sensor_data
}

type SeabatII struct {
	PingNumber           []int16
	SurfaceSoundVelocity []float32
	Mode                 []int16
	SonarRange           []int16
	TransmitPower        []int16
	ReceiveGain          []int16
	ForeAftBandwidth     []float32
	AthwartBandwidth     []float32
}

func DecodeSeabatIISpecific(reader *bytes.Reader) (sensor_data SeabatII) {
	var buffer struct {
		PingNumber           int16
		SurfaceSoundVelocity int16
		Mode                 int16
		SonarRange           int16
		TransmitPower        int16
		ReceiveGain          int16
		ForeAftBandwidth     int8
		AthwartBandwidth     int8
		Spare                int32
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / 10.0}
	sensor_data.Mode = []int16{buffer.Mode}
	sensor_data.SonarRange = []int16{buffer.SonarRange}
	sensor_data.TransmitPower = []int16{buffer.TransmitPower}
	sensor_data.ReceiveGain = []int16{buffer.ReceiveGain}
	sensor_data.ForeAftBandwidth = []float32{float32(buffer.ForeAftBandwidth) / 10.0}
	sensor_data.AthwartBandwidth = []float32{float32(buffer.AthwartBandwidth) / 10.0}

	return sensor_data
}

type Seabat8101 struct {
	PingNumber           []int16
	SurfaceSoundVelocity []float32
	Mode                 []int16
	Range                []int16
	TransmitPower        []int16
	RecieveGain          []int16
	PulseWidth           []int16
	TvgSpreading         []int8
	TvgAbsorption        []int8
	ForeAftBandwidth     []float32
	AthwartBandwidth     []float32
	RangeFilterMin       []float32
	RangeFilterMax       []float32
	DepthFilterMin       []float32
	DepthFilterMax       []float32
	ProjectorType        []int8
}

func DecodeSeabat8101Specific(reader *bytes.Reader) (sensor_data Seabat8101) {
	var buffer struct {
		PingNumber           int16
		SurfaceSoundVelocity int16
		Mode                 int16
		Range                int16
		TransmitPower        int16
		RecieveGain          int16
		PulseWidth           int16
		TvgSpreading         int8
		TvgAbsorption        int8
		ForeAftBandwidth     int8
		AthwartBandwidth     int8
		RangeFilterMin       int16
		RangeFilterMax       int16
		DepthFilterMin       int16
		DepthFilterMax       int16
		ProjectorType        int8
		Spare                int32
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / 10.0}
	sensor_data.Mode = []int16{buffer.Mode}
	sensor_data.Range = []int16{buffer.Range}
	sensor_data.TransmitPower = []int16{buffer.TransmitPower}
	sensor_data.RecieveGain = []int16{buffer.RecieveGain}
	sensor_data.PulseWidth = []int16{buffer.PulseWidth}
	sensor_data.TvgSpreading = []int8{buffer.TvgSpreading}
	sensor_data.TvgAbsorption = []int8{buffer.TvgAbsorption}
	sensor_data.ForeAftBandwidth = []float32{float32(buffer.ForeAftBandwidth) / 10.0}
	sensor_data.AthwartBandwidth = []float32{float32(buffer.AthwartBandwidth) / 10.0}
	sensor_data.RangeFilterMin = []float32{float32(buffer.RangeFilterMin)}
	sensor_data.RangeFilterMax = []float32{float32(buffer.RangeFilterMax)}
	sensor_data.DepthFilterMin = []float32{float32(buffer.DepthFilterMin)}
	sensor_data.DepthFilterMax = []float32{float32(buffer.DepthFilterMax)}
	sensor_data.ProjectorType = []int8{buffer.ProjectorType}

	return sensor_data
}

type Seabeam2112 struct {
	Mode                   []int8
	SurfaceSoundVelocity   []float32
	SsvSource              []int8
	PingGain               []int8
	PulseWidth             []int8
	TransmitterAttenuation []int8
	NumberAlgorithms       []int8
	AlgorithmOrder         []string
}

func DecodeSeabeam2112Specific(reader *bytes.Reader) (sensor_data Seabeam2112) {
	var buffer struct {
		Mode                   int8
		SurfaceSoundVelocity   float32
		SsvSource              int8
		PingGain               int8
		PulseWidth             int8
		TransmitterAttenuation int8
		NumberAlgorithms       int8
		AlgorithmOrder         [5]byte
		Spare                  int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.Mode = []int8{buffer.Mode}
	sensor_data.SurfaceSoundVelocity = []float32{(float32(buffer.SurfaceSoundVelocity) + 130000.0) / 100.0}
	sensor_data.SsvSource = []int8{buffer.SsvSource}
	sensor_data.PingGain = []int8{buffer.PingGain}
	sensor_data.PulseWidth = []int8{buffer.PulseWidth}
	sensor_data.TransmitterAttenuation = []int8{buffer.TransmitterAttenuation}
	sensor_data.NumberAlgorithms = []int8{buffer.NumberAlgorithms}
	sensor_data.AlgorithmOrder = []string{string(buffer.AlgorithmOrder[:])}

	return sensor_data
}

type ElacMkII struct {
	Mode                  []int8
	PingNumber            []int16
	SurfaceSoundVelocity  []int16
	PulseLength           []int16
	ReceiverGainStarboard []int8
	ReceiverGainPort      []int8
}

func DecodeElacMkIISpecific(reader *bytes.Reader) (sensor_data ElacMkII) {
	var buffer struct {
		Mode                  int8
		PingNumber            int16
		SurfaceSoundVelocity  int16
		PulseLength           int16
		ReceiverGainStarboard int8
		ReceiverGainPort      int8
		Spare                 int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.Mode = []int8{buffer.Mode}
	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []int16{buffer.SurfaceSoundVelocity}
	sensor_data.PulseLength = []int16{buffer.PulseLength}
	sensor_data.ReceiverGainStarboard = []int8{buffer.ReceiverGainStarboard}
	sensor_data.ReceiverGainPort = []int8{buffer.ReceiverGainPort}

	return sensor_data
}

type CmpSass struct {
	Lfreq  []float32
	Lntens []float32
}

func DecodeCmpSass(reader *bytes.Reader) (sensor_data CmpSass) {
	var buffer struct {
		Lfreq  int16
		Lntens int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.Lfreq = []float32{float32(buffer.Lfreq) / 10.0}
	sensor_data.Lntens = []float32{float32(buffer.Lntens) / 10.0}

	return sensor_data
}

type Reson8100 struct {
	Latency              []int16
	PingNumber           []int32
	SonarID              []int32
	SonarModel           []int16
	Frequency            []int16
	SurfaceSoundVelocity []float32
	SampleRate           []int16
	PingRate             []int16
	Mode                 []int16
	Range                []int16
	TransmitPower        []int16
	ReceiveGain          []int16
	PulseWidth           []int16
	TvgSpreading         []int8
	TvgAbsorption        []int8
	ForeAftBandwidth     []float32
	AthwartBandwidth     []float32
	ProjectorType        []int8
	ProjectorAngle       []int16
	RangeFilterMin       []float32
	RangeFilterMax       []float32
	DepthFilterMin       []float32
	DepthFilterMax       []float32
	FiltersActive        []int8
	Temperature          []int16
	BeamSpacing          []float32
}

func DecodeReson8100(reader *bytes.Reader) (sensor_data Reson8100) {
	var buffer struct {
		Latency              int16
		PingNumber           int32
		SonarID              int32
		SonarModel           int16
		Frequency            int16
		SurfaceSoundVelocity int16
		SampleRate           int16
		PingRate             int16
		Mode                 int16
		Range                int16
		TransmitPower        int16
		ReceiveGain          int16
		PulseWidth           int16
		TvgSpreading         int8
		TvgAbsorption        int8
		ForeAftBandwidth     int8
		AthwartBandwidth     int8
		ProjectorType        int8
		ProjectorAngle       int16
		RangeFilterMin       int16
		RangeFilterMax       int16
		DepthFilterMin       int16
		DepthFilterMax       int16
		FiltersActive        int8
		Temperature          int16
		BeamSpacing          int16
		Spare                int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.Latency = []int16{buffer.Latency}
	sensor_data.PingNumber = []int32{buffer.PingNumber}
	sensor_data.SonarID = []int32{buffer.SonarID}
	sensor_data.SonarModel = []int16{buffer.SonarModel}
	sensor_data.Frequency = []int16{buffer.Frequency}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / 10.0}
	sensor_data.SampleRate = []int16{buffer.SampleRate}
	sensor_data.PingRate = []int16{buffer.PingRate}
	sensor_data.Mode = []int16{buffer.Mode}
	sensor_data.Range = []int16{buffer.Range}
	sensor_data.TransmitPower = []int16{buffer.TransmitPower}
	sensor_data.ReceiveGain = []int16{buffer.ReceiveGain}
	sensor_data.PulseWidth = []int16{buffer.PulseWidth}
	sensor_data.TvgSpreading = []int8{buffer.TvgSpreading}
	sensor_data.TvgAbsorption = []int8{buffer.TvgAbsorption}
	sensor_data.ForeAftBandwidth = []float32{float32(buffer.ForeAftBandwidth) / 10.0}
	sensor_data.AthwartBandwidth = []float32{float32(buffer.AthwartBandwidth) / 10.0}
	sensor_data.ProjectorType = []int8{buffer.ProjectorType}
	sensor_data.ProjectorAngle = []int16{buffer.ProjectorAngle}
	sensor_data.RangeFilterMin = []float32{float32(buffer.RangeFilterMin)}
	sensor_data.RangeFilterMax = []float32{float32(buffer.RangeFilterMax)}
	sensor_data.DepthFilterMin = []float32{float32(buffer.DepthFilterMin)}
	sensor_data.DepthFilterMax = []float32{float32(buffer.DepthFilterMax)}
	sensor_data.FiltersActive = []int8{buffer.FiltersActive}
	sensor_data.Temperature = []int16{buffer.Temperature}
	sensor_data.BeamSpacing = []float32{float32(buffer.BeamSpacing) / 10000.0}

	return sensor_data
}

type Em3 struct {
	ModelNumber          []int16
	PingNumber           []int16
	SerialNumber         []int16
	SurfaceSoundVelocity []float32
	TransducerDepth      []float32
	ValidBeams           []int16
	SampleRate           []int16
	DepthDifference      []float32
	OffsetMultiplier     []int8
	// RunTimeID                      []uint32 // not stored
	RunTimeModelNumber             [][]int16
	RunTimeDgTime                  [][]time.Time
	RunTimePingNumber              [][]int16
	RunTimeSerialNumber            [][]int16
	RunTimeSystemStatus            [][]uint32
	RunTimeMode                    [][]int8
	RunTimeFilterID                [][]int8
	RunTimeMinDepth                [][]float32
	RunTimeMaxDepth                [][]float32
	RunTimeAbsorption              [][]float32
	RunTimeTransmitPulseLength     [][]float32
	RunTimeTransmitBeamWidth       [][]float32
	RunTimePowerReduction          [][]int8
	RunTimeReceiveBeamWidth        [][]float32
	RunTimeReceiveBandwidth        [][]int16
	RunTimeReceiveGain             [][]int8
	RunTimeCrossOverAngle          [][]int8
	RunTimeSsvSource               [][]int8
	RunTimePortSwathWidth          [][]int16
	RunTimeBeamSpacing             [][]int8
	RunTimePortCoverageSector      [][]int8
	RunTimeStabilization           [][]int8
	RunTimeStarboardCoverageSector [][]int8
	RunTimeStarboardSwathWidth     [][]int16
	RunTimeHiloFreqAbsorpRatio     [][]int8
	RunTimeSwathWidth              [][]int16
	RunTimeCoverageSector          [][]int16
}

func DecodeEm3(reader *bytes.Reader) (sensor_data Em3) {
	var (
		buffer struct {
			ModelNumber          int16
			PingNumber           int16
			SerialNumber         int16
			SurfaceSoundVelocity int16
			TransducerDepth      int16
			ValidBeams           int16
			SampleRate           int16
			DepthDifference      int16
			OffsetMultiplier     int8
			RunTimeID            uint32
		}
		rt1 struct {
			ModelNumber             int16
			TvSec                   uint32
			TvNSec                  uint32
			PingNumber              int16
			SerialNumber            int16
			SystemStatus            uint32
			Mode                    int8
			FilterID                int8
			MinDepth                int16
			MaxDepth                int16
			Absoprtion              int16
			TransmitPulseLength     int16
			TransmitBeamWidth       int16
			PowerReduction          int8
			ReceiveBeamWidth        int8
			ReceiveBandwidth        int8
			ReceiveGain             int8
			CrossOverAnlge          int8
			SsvSource               int8
			PortSwathWidth          int16
			BeamSpacing             int8
			PortCoverageSector      int8
			Stabilization           int8
			StarboardCoverageSector int8
			StarboardSwathWidth     int16
			HiloFreqAbsorpRatio     int8
			Spare                   int32
		}
		rt2 struct {
			ModelNumber             int16
			TvSec                   uint32
			TvNSec                  uint32
			PingNumber              int16
			SerialNumber            int16
			SystemStatus            uint32
			Mode                    int8
			FilterID                int8
			MinDepth                int16
			MaxDepth                int16
			Absoprtion              int16
			TransmitPulseLength     int16
			TransmitBeamWidth       int16
			PowerReduction          int8
			ReceiveBeamWidth        int8
			ReceiveBandwidth        int8
			ReceiveGain             int8
			CrossOverAnlge          int8
			SsvSource               int8
			PortSwathWidth          int16
			BeamSpacing             int8
			PortCoverageSector      int8
			Stabilization           int8
			StarboardCoverageSector int8
			StarboardSwathWidth     int16
			HiloFreqAbsorpRatio     int8
			Spare                   int32
		}
	)
	model_number := make([]int16, 0, 2)
	dg_time := make([]time.Time, 0, 2)
	ping_number := make([]int16, 0, 2)
	serial_number := make([]int16, 0, 2)
	system_status := make([]uint32, 0, 2)
	mode := make([]int8, 0, 2)
	filter_id := make([]int8, 0, 2)
	min_depth := make([]float32, 0, 2)
	max_depth := make([]float32, 0, 2)
	absorption := make([]float32, 0, 2)
	transmit_pulse_length := make([]float32, 0, 2)
	transmit_beam_width := make([]float32, 0, 2)
	power_reduction := make([]int8, 0, 2)
	receive_beamwidth := make([]float32, 0, 2)
	receive_bandwidth := make([]int16, 0, 2)
	receive_gain := make([]int8, 0, 2)
	cross_over_angle := make([]int8, 0, 2)
	ssv_source := make([]int8, 0, 2)
	port_swath_width := make([]int16, 0, 2)
	beam_spacing := make([]int8, 0, 2)
	port_coverage_sector := make([]int8, 0, 2)
	stabilization := make([]int8, 0, 2)
	starboard_coverage_sector := make([]int8, 0, 2)
	starboard_swath_width := make([]int16, 0, 2)
	hilo_freq_absorp_ratio := make([]int8, 0, 2)
	swath_width := make([]int16, 0, 2)
	coverage_sector := make([]int16, 0, 2)

	_ = binary.Read(reader, binary.BigEndian, &buffer)

	// base set of values
	sensor_data.ModelNumber = []int16{buffer.ModelNumber}
	sensor_data.PingNumber = []int16{buffer.PingNumber}
	sensor_data.SerialNumber = []int16{buffer.SerialNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SerialNumber) / 10.0}
	sensor_data.TransducerDepth = []float32{float32(buffer.TransducerDepth) / SCALE2}
	sensor_data.ValidBeams = []int16{buffer.ValidBeams}
	sensor_data.SampleRate = []int16{buffer.SampleRate}
	sensor_data.DepthDifference = []float32{float32(buffer.DepthDifference) / SCALE2}
	sensor_data.OffsetMultiplier = []int8{buffer.OffsetMultiplier}

	// runtime values
	if (buffer.RunTimeID & 0x00000001) != 0 {
		_ = binary.Read(reader, binary.BigEndian, &rt1)
		model_number = append(model_number, rt1.ModelNumber)
		dg_time = append(dg_time, time.Unix(int64(rt1.TvSec), int64(rt1.TvNSec)).UTC())
		ping_number = append(ping_number, rt1.PingNumber)
		serial_number = append(serial_number, rt1.SerialNumber)
		system_status = append(system_status, rt1.SystemStatus)
		mode = append(mode, rt1.Mode)
		filter_id = append(filter_id, rt1.FilterID)
		min_depth = append(min_depth, float32(rt1.MinDepth))
		max_depth = append(max_depth, float32(rt1.MaxDepth))
		absorption = append(absorption, float32(rt2.Absoprtion)/SCALE2)
		transmit_pulse_length = append(transmit_pulse_length, float32(rt1.TransmitPulseLength))
		transmit_beam_width = append(transmit_beam_width, float32(rt1.TransmitBeamWidth)/10.0)
		power_reduction = append(power_reduction, rt1.PowerReduction)
		receive_beamwidth = append(receive_beamwidth, float32(rt1.ReceiveBeamWidth)/10.0)
		receive_bandwidth = append(receive_bandwidth, int16(rt1.ReceiveBandwidth)*50)
		receive_gain = append(receive_gain, rt1.ReceiveGain)
		cross_over_angle = append(cross_over_angle, rt1.CrossOverAnlge)
		ssv_source = append(ssv_source, rt1.SsvSource)
		// port_swath_width = append(port_swath_width, rt1.PortSwathWidth)
		beam_spacing = append(beam_spacing, rt1.BeamSpacing)
		//port_coverage_sector = append(port_coverage_sector, rt1.PortCoverageSector)
		stabilization = append(stabilization, rt1.Stabilization)
		// starboard_coverage_sector = append(starboard_coverage_sector, rt1.StarboardCoverageSector)
		// starboard_swath_width = append(starboard_swath_width, rt1.StarboardSwathWidth)
		hilo_freq_absorp_ratio = append(hilo_freq_absorp_ratio, rt1.HiloFreqAbsorpRatio)

		if rt1.StarboardSwathWidth != 0 {
			swath_width = append(swath_width, rt1.PortSwathWidth+rt1.StarboardSwathWidth)
			port_swath_width = append(port_swath_width, rt1.PortSwathWidth)
			starboard_swath_width = append(starboard_swath_width, rt1.StarboardSwathWidth)
		} else {
			swath_width = append(swath_width, rt1.PortSwathWidth)
			port_swath_width = append(port_swath_width, rt1.PortSwathWidth/2)
			starboard_swath_width = append(starboard_swath_width, rt1.PortSwathWidth/2)
		}

		if rt1.StarboardCoverageSector != 0 {
			coverage_sector = append(coverage_sector, int16(rt1.PortCoverageSector)+int16(rt1.StarboardCoverageSector))
			port_coverage_sector = append(port_coverage_sector, rt1.PortCoverageSector)
			starboard_coverage_sector = append(starboard_coverage_sector, rt1.StarboardCoverageSector)
		} else {
			coverage_sector = append(coverage_sector, int16(rt1.PortCoverageSector))
			port_coverage_sector = append(port_coverage_sector, rt1.PortCoverageSector/2)
			starboard_coverage_sector = append(starboard_coverage_sector, rt1.PortCoverageSector/2)
		}

		if (buffer.RunTimeID & 0x00000002) != 0 {
			_ = binary.Read(reader, binary.BigEndian, &rt2)
			model_number = append(model_number, rt2.ModelNumber)
			dg_time = append(dg_time, time.Unix(int64(rt2.TvSec), int64(rt2.TvNSec)).UTC())
			ping_number = append(ping_number, rt2.PingNumber)
			serial_number = append(serial_number, rt2.SerialNumber)
			system_status = append(system_status, rt2.SystemStatus)
			mode = append(mode, rt2.Mode)
			filter_id = append(filter_id, rt2.FilterID)
			min_depth = append(min_depth, float32(rt2.MinDepth))
			max_depth = append(max_depth, float32(rt2.MaxDepth))
			absorption = append(absorption, float32(rt2.Absoprtion)/SCALE2)
			transmit_pulse_length = append(transmit_pulse_length, float32(rt2.TransmitPulseLength))
			transmit_beam_width = append(transmit_beam_width, float32(rt2.TransmitBeamWidth)/10.0)
			power_reduction = append(power_reduction, rt2.PowerReduction)
			receive_beamwidth = append(receive_beamwidth, float32(rt2.ReceiveBeamWidth)/10.0)
			receive_bandwidth = append(receive_bandwidth, int16(rt2.ReceiveBandwidth)*50)
			receive_gain = append(receive_gain, rt2.ReceiveGain)
			cross_over_angle = append(cross_over_angle, rt2.CrossOverAnlge)
			ssv_source = append(ssv_source, rt2.SsvSource)
			// port_swath_width = append(port_swath_width, rt2.PortSwathWidth)
			beam_spacing = append(beam_spacing, rt2.BeamSpacing)
			port_coverage_sector = append(port_coverage_sector, rt2.PortCoverageSector)
			stabilization = append(stabilization, rt2.Stabilization)
			starboard_coverage_sector = append(starboard_coverage_sector, rt2.StarboardCoverageSector)
			// starboard_swath_width = append(starboard_swath_width, rt2.StarboardSwathWidth)
			hilo_freq_absorp_ratio = append(hilo_freq_absorp_ratio, rt2.HiloFreqAbsorpRatio)
		}

		if rt2.StarboardSwathWidth != 0 {
			swath_width = append(swath_width, rt2.PortSwathWidth+rt2.StarboardSwathWidth)
			port_swath_width = append(port_swath_width, rt2.PortSwathWidth)
			starboard_swath_width = append(starboard_swath_width, rt2.StarboardSwathWidth)
		} else {
			swath_width = append(swath_width, rt2.PortSwathWidth)
			port_swath_width = append(port_swath_width, rt2.PortSwathWidth/2)
			starboard_swath_width = append(starboard_swath_width, rt2.PortSwathWidth/2)
		}

		if rt2.StarboardCoverageSector != 0 {
			coverage_sector = append(coverage_sector, int16(rt2.PortCoverageSector)+int16(rt2.StarboardCoverageSector))
			port_coverage_sector = append(port_coverage_sector, rt2.PortCoverageSector)
			starboard_coverage_sector = append(starboard_coverage_sector, rt2.StarboardCoverageSector)
		} else {
			coverage_sector = append(coverage_sector, int16(rt2.PortCoverageSector))
			port_coverage_sector = append(port_coverage_sector, rt2.PortCoverageSector/2)
			starboard_coverage_sector = append(starboard_coverage_sector, rt2.PortCoverageSector/2)
		}
	}

	// insert runtime data
	sensor_data.RunTimeModelNumber = [][]int16{model_number}
	sensor_data.RunTimeDgTime = [][]time.Time{dg_time}
	sensor_data.RunTimePingNumber = [][]int16{ping_number}
	sensor_data.RunTimeSerialNumber = [][]int16{serial_number}
	sensor_data.RunTimeMode = [][]int8{mode}
	sensor_data.RunTimeFilterID = [][]int8{filter_id}
	sensor_data.RunTimeMinDepth = [][]float32{min_depth}
	sensor_data.RunTimeMaxDepth = [][]float32{max_depth}
	sensor_data.RunTimeAbsorption = [][]float32{absorption}
	sensor_data.RunTimeTransmitPulseLength = [][]float32{transmit_pulse_length}
	sensor_data.RunTimeTransmitBeamWidth = [][]float32{transmit_beam_width}
	sensor_data.RunTimePowerReduction = [][]int8{power_reduction}
	sensor_data.RunTimeReceiveBeamWidth = [][]float32{receive_beamwidth}
	sensor_data.RunTimeReceiveBandwidth = [][]int16{receive_bandwidth}
	sensor_data.RunTimeReceiveGain = [][]int8{receive_gain}
	sensor_data.RunTimeCrossOverAngle = [][]int8{cross_over_angle}
	sensor_data.RunTimeSsvSource = [][]int8{ssv_source}
	sensor_data.RunTimePortSwathWidth = [][]int16{port_swath_width}
	sensor_data.RunTimeBeamSpacing = [][]int8{beam_spacing}
	sensor_data.RunTimePortCoverageSector = [][]int8{port_coverage_sector}
	sensor_data.RunTimeStabilization = [][]int8{stabilization}
	sensor_data.RunTimeStarboardCoverageSector = [][]int8{starboard_coverage_sector}
	sensor_data.RunTimeStarboardSwathWidth = [][]int16{starboard_swath_width}
	sensor_data.RunTimeHiloFreqAbsorpRatio = [][]int8{hilo_freq_absorp_ratio}
	sensor_data.RunTimeSwathWidth = [][]int16{swath_width}
	sensor_data.RunTimeCoverageSector = [][]int16{coverage_sector}

	return sensor_data
}

// float64 may be overkill
// where scale factors are applied, float32 is used
// where it's confident float32 is enough to represent the value
// TODO; align into 64bit chunks
// the spec says binary integers are stored as either 1-byte unsigned, 2-byte signed or unsigned, or 4-byte signed
type Em4 struct {
	ModelNumber                       []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	PingCounter                       []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	SerialNumber                      []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceVelocity                   []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransducerDepth                   []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ValidDetections                   []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	SamplingFrequency                 []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	DopplerCorrectionScale            []int32     `tiledb:"dtype=int32,ftype=attr" filters:"zstd(level=16)"`
	VehicleDepth                      []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransmitSectors                   []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	TiltAngle                         [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	FocusRange                        [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	SignalLength                      [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	TransmitDelay                     [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	CenterFrequency                   [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	MeanAbsorption                    [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	WaveformId                        [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	SectorNumber                      [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	SignalBandwith                    [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeModelNumber                []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeDatagramTime               []time.Time `tiledb:"dtype=datetime_ns,ftype=attr" filters:"zstd(level=16)"`
	RunTimePingCounter                []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeSerialNumber               []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeOperatorStationStatus      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeProcessingUnitStatus       []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeBspStatus                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeHeadTransceiverStatus      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMode                       []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeFilterId                   []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMinDepth                   []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMaxDepth                   []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeAbsorption                 []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTransmitPulseLength        []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTransmitBeamWidth          []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTransmitPowerReduction     []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeReceiveBeamWidth           []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeReceiveBandwidth           []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeReceiveFixedGain           []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTvgCrossOverAngle          []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeSsvSource                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMaxPortSwathWidth          []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeBeamSpacing                []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMaxPortCoverage            []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeStabilization              []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMaxStbdCoverage            []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMaxStdbSwathWidth          []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTransmitAlongTilt          []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeFilterId2                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ProcessorUnitCpuLoad              []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ProcessorUnitSensorStatus         []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ProcessorUnitAchievedPortCoverage []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ProcessorUnitAchievedStbdCoverage []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ProcessorUnitYawStabilization     []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm4Specific(reader *bytes.Reader) (sensor_data EM4) {

	var (
		buffer struct {
			ModelNumber            int16
			PingCounter            int16
			SerialNumber           int16
			SurfaceVelocity        int16
			TransducerDepth        int32
			ValidDetections        int16
			SamplingFrequency1     int32
			SamplingFrequency2     int32
			DopplerCorrectionScale int32
			VehicleDepth           int32
			Spare                  [4]int32
			TransmitSectors        int16
		} // 48 bytes
		sector_buffer struct {
			TiltAngle       []float32
			FocusRange      []float32
			SignalLength    []float64
			TransmitDelay   []float64
			CenterFrequency []float64
			MeanAbsorption  []float32
			WaveformId      []uint8
			SectorNumber    []uint8
			SignalBandwith  []float64
			Spare           [4]int32
		}
		sector_buffer_base struct {
			TiltAngle       int16
			FocusRange      int16
			SignalLength    int32
			TransmitDelay   int32
			CenterFrequency int32
			MeanAbsorption  int16
			WaveformId      uint8
			SectorNumber    uint8
			SignalBandwith  int32
			Spare           [4]int32
		} // 40 bytes
		spare_buffer struct {
			Spare [4]int32
		} // 16 bytes
		runtime_buffer struct {
			RunTimeModelNumber            int16
			RunTimeDatagramTime_sec       int32
			RunTimeDatagramTime_nsec      int32
			RunTimePingCounter            int16
			RunTimeSerialNumber           int16
			RunTimeOperatorStationStatus  uint8
			RunTimeProcessingUnitStatus   uint8
			RunTimeBspStatus              uint8
			RunTimeHeadTransceiverStatus  uint8
			RunTimeMode                   uint8
			RunTimeFilterId               uint8
			RunTimeMinDepth               int16
			RunTimeMaxDepth               int16
			RunTimeAbsorption             int16
			RunTimeTransmitPulseLength    int16
			RunTimeTransmitBeamWidth      int16
			RunTimeTransmitPowerReduction uint8
			RunTimeReceiveBeamWidth       uint8
			RunTimeReceiveBandwidth       uint8
			RunTimeReceiveFixedGain       uint8
			RunTimeTvgCrossOverAngle      uint8
			RunTimeSsvSource              uint8
			RunTimeMaxPortSwathWidth      int16
			RunTimeBeamSpacing            uint8
			RunTimeMaxPortCoverage        uint8
			RunTimeStabilization          uint8
			RunTimeMaxStbdCoverage        uint8
			RunTimeMaxStdbSwathWidth      int16
			RunTimeTransmitAlongTilt      int16
			RunTimeFilterId2              uint8
			Spare                         [4]int32
		} // 63 bytes
		proc_buffer struct {
			ProcessorUnitCpuLoad              uint8
			ProcessorUnitSensorStatus         uint16
			ProcessorUnitAchievedPortCoverage uint8
			ProcessorUnitAchievedStbdCoverage uint8
			ProcessorUnitYawStabilization     int16
			Spare                             [4]int32
		} // 23 bytes
	)

	// first 46 bytes
	_ = binary.Read(reader, binary.BigEndian, &buffer)
	// n_bytes += 46

	// sector arrays
	// what if TransmitSectors == 0???
	for i := int16(0); i < buffer.TransmitSectors; i++ {
		_ = binary.Read(reader, binary.BigEndian, &sector_buffer_base)
		sector_buffer.TiltAngle = append(
			sector_buffer.TiltAngle,
			float32(sector_buffer_base.TiltAngle)/float32(100),
		)
		sector_buffer.FocusRange = append(
			sector_buffer.FocusRange,
			float32(sector_buffer_base.FocusRange)/float32(10),
		)
		sector_buffer.SignalLength = append(
			sector_buffer.SignalLength,
			float64(sector_buffer_base.SignalLength)/float64(1_000_000),
		)
		sector_buffer.TransmitDelay = append(
			sector_buffer.TransmitDelay,
			float64(sector_buffer_base.TransmitDelay)/float64(1_000_000),
		)
		sector_buffer.CenterFrequency = append(
			sector_buffer.CenterFrequency,
			float64(sector_buffer_base.CenterFrequency)/float64(1000),
		)
		sector_buffer.MeanAbsorption = append(
			sector_buffer.MeanAbsorption,
			float32(sector_buffer_base.MeanAbsorption)/float32(100),
		)
		sector_buffer.WaveformId = append(
			sector_buffer.WaveformId,
			sector_buffer_base.WaveformId,
		)
		sector_buffer.SectorNumber = append(
			sector_buffer.SectorNumber,
			sector_buffer_base.SectorNumber,
		)
		sector_buffer.SignalBandwith = append(
			sector_buffer.SignalBandwith,
			float64(sector_buffer_base.SignalBandwith)/float64(1000),
		)
	}

	// spare 16 bytes
	_ = binary.Read(reader, binary.BigEndian, &spare_buffer)

	// next 63 bytes for the RunTime info
	_ = binary.Read(reader, binary.BigEndian, &runtime_buffer)

	// next 23 bytes for the processing unit info
	_ = binary.Read(reader, binary.BigEndian, &proc_buffer)

	// populate generic
	sensor_data.ModelNumber = []int16{buffer.ModelNumber}
	sensor_data.PingCounter = []int16{buffer.PingCounter}
	sensor_data.SerialNumber = []int16{buffer.SerialNumber}
	sensor_data.SurfaceVelocity = []float32{float32(buffer.SurfaceVelocity) / float32(10)}
	sensor_data.TransducerDepth = []float64{float64(buffer.TransducerDepth) / float64(20000)}
	sensor_data.ValidDetections = []int16{buffer.ValidDetections}
	sensor_data.SamplingFrequency = []float64{float64(buffer.SamplingFrequency1) + float64(buffer.SamplingFrequency2)/float64(4_000_000_000)}
	sensor_data.DopplerCorrectionScale = []int32{buffer.DopplerCorrectionScale}
	sensor_data.VehicleDepth = []float32{float32(buffer.VehicleDepth) / float32(1000)}
	sensor_data.TransmitSectors = []int16{buffer.TransmitSectors}

	// populate sector info
	sensor_data.TiltAngle = [][]float32{sector_buffer.TiltAngle}
	sensor_data.FocusRange = [][]float32{sector_buffer.FocusRange}
	sensor_data.SignalLength = [][]float64{sector_buffer.SignalLength}
	sensor_data.TransmitDelay = [][]float64{sector_buffer.TransmitDelay}
	sensor_data.CenterFrequency = [][]float64{sector_buffer.CenterFrequency}
	sensor_data.MeanAbsorption = [][]float32{sector_buffer.MeanAbsorption}
	sensor_data.WaveformId = [][]uint8{sector_buffer.WaveformId}
	sensor_data.SectorNumber = [][]uint8{sector_buffer.SectorNumber}
	sensor_data.SignalBandwith = [][]float64{sector_buffer.SignalBandwith}

	// populate runtime info
	sensor_data.RunTimeModelNumber = []int16{runtime_buffer.RunTimeModelNumber}
	sensor_data.RunTimeDatagramTime = []time.Time{time.Unix(
		int64(runtime_buffer.RunTimeDatagramTime_sec),
		int64(runtime_buffer.RunTimeDatagramTime_nsec),
	)}
	sensor_data.RunTimePingCounter = []int16{runtime_buffer.RunTimePingCounter}
	sensor_data.RunTimeSerialNumber = []int16{runtime_buffer.RunTimeSerialNumber}
	sensor_data.RunTimeOperatorStationStatus = []uint8{runtime_buffer.RunTimeOperatorStationStatus}
	sensor_data.RunTimeProcessingUnitStatus = []uint8{runtime_buffer.RunTimeProcessingUnitStatus}
	sensor_data.RunTimeBspStatus = []uint8{runtime_buffer.RunTimeBspStatus}
	sensor_data.RunTimeHeadTransceiverStatus = []uint8{runtime_buffer.RunTimeHeadTransceiverStatus}
	sensor_data.RunTimeMode = []uint8{runtime_buffer.RunTimeMode}
	sensor_data.RunTimeFilterId = []uint8{runtime_buffer.RunTimeFilterId}
	sensor_data.RunTimeMinDepth = []float32{float32(runtime_buffer.RunTimeMinDepth)}
	sensor_data.RunTimeMaxDepth = []float32{float32(runtime_buffer.RunTimeMaxDepth)}
	sensor_data.RunTimeAbsorption = []float32{float32(runtime_buffer.RunTimeAbsorption) / float32(100)}
	sensor_data.RunTimeTransmitPulseLength = []float32{float32(runtime_buffer.RunTimeTransmitPulseLength)}
	sensor_data.RunTimeTransmitBeamWidth = []float32{float32(runtime_buffer.RunTimeTransmitBeamWidth) / float32(10)}
	sensor_data.RunTimeTransmitPowerReduction = []uint8{runtime_buffer.RunTimeTransmitPowerReduction}
	sensor_data.RunTimeReceiveBeamWidth = []float32{float32(runtime_buffer.RunTimeReceiveBeamWidth) / float32(10)}
	sensor_data.RunTimeReceiveBandwidth = []float32{float32(runtime_buffer.RunTimeReceiveBandwidth) * float32(50)}
	sensor_data.RunTimeReceiveFixedGain = []uint8{runtime_buffer.RunTimeReceiveFixedGain}
	sensor_data.RunTimeTvgCrossOverAngle = []uint8{runtime_buffer.RunTimeTvgCrossOverAngle}
	sensor_data.RunTimeSsvSource = []uint8{runtime_buffer.RunTimeSsvSource}
	sensor_data.RunTimeMaxPortSwathWidth = []int16{runtime_buffer.RunTimeMaxPortSwathWidth}
	sensor_data.RunTimeBeamSpacing = []uint8{runtime_buffer.RunTimeBeamSpacing}
	sensor_data.RunTimeMaxPortCoverage = []uint8{runtime_buffer.RunTimeMaxPortCoverage}
	sensor_data.RunTimeStabilization = []uint8{runtime_buffer.RunTimeStabilization}
	sensor_data.RunTimeMaxStbdCoverage = []uint8{runtime_buffer.RunTimeMaxStbdCoverage}
	sensor_data.RunTimeMaxStdbSwathWidth = []int16{runtime_buffer.RunTimeMaxStdbSwathWidth}
	sensor_data.RunTimeTransmitAlongTilt = []float32{float32(runtime_buffer.RunTimeTransmitAlongTilt) / float32(100)}
	sensor_data.RunTimeFilterId2 = []uint8{runtime_buffer.RunTimeFilterId2}

	// populate processor unit info
	sensor_data.ProcessorUnitCpuLoad = []uint8{proc_buffer.ProcessorUnitCpuLoad}
	sensor_data.ProcessorUnitSensorStatus = []uint16{proc_buffer.ProcessorUnitSensorStatus}
	sensor_data.ProcessorUnitAchievedPortCoverage = []uint8{proc_buffer.ProcessorUnitAchievedPortCoverage}
	sensor_data.ProcessorUnitAchievedStbdCoverage = []uint8{proc_buffer.ProcessorUnitAchievedStbdCoverage}
	sensor_data.ProcessorUnitYawStabilization = []float32{float32(proc_buffer.ProcessorUnitYawStabilization) / float32(100)}

	return sensor_data
}
