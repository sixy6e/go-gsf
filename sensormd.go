package gsf

import (
	"bytes"
	"encoding/binary"
	"time"
)

// It's problematic to determine the correct read type in cases
// where the spec differs to the c-code.
// In many instances the spec says int16, but the code reads it as a uint16
// and then converts to an int32.
// If the spec says int16, then we should be safe in reading as int16.
// However, I've come across many instances where the code does something
// very different to the spec.
// Best attempts will be made to infer the "more correct" type if something
// doesn't look right:
// (differs wildly in the spec vs code, as well as what the data represents).
// PingNumber, especially when stating "Sequential ping counter, 0 through 65535"
// should be unsigned anyway (why the spec is signed, no idea). But weird when
// some cases use uint16 to decode, then promote to int32. For those that do
// we'll keep as uint16 and not promote to int32.
// TODO; look at the int8 assignments and confirm if *p is unsigned

type Seabeam struct {
	EclipseTime []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
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
	PingNumber    []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Resolution    []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PingQuality   []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Mode          []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm12Specific(reader *bytes.Reader) (sensor_data Em12) {
	var buffer struct {
		PingNumber    uint16
		Resolution    uint8
		PingQuality   uint8
		SoundVelocity uint16
		Mode          uint8
		Spare         [4]int32
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.Resolution = []uint8{buffer.Resolution}
	sensor_data.PingQuality = []uint8{buffer.PingQuality}
	sensor_data.SoundVelocity = []float32{float32(buffer.SoundVelocity) / 10.0}
	sensor_data.Mode = []uint8{buffer.Mode}

	return sensor_data
}

type Em100 struct {
	ShipPitch       []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransducerPitch []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Mode            []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Power           []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Attenuation     []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Tvg             []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PulseLength     []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Counter         []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm100Specific(reader *bytes.Reader) (sensor_data Em100) {
	var buffer struct {
		ShipPitch       int16
		TransducerPitch int16
		Mode            uint8
		Power           uint8
		Attenuation     uint8
		Tvg             uint8
		PulseLength     uint8
		Counter         uint16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.ShipPitch = []float32{float32(buffer.ShipPitch) / SCALE2}
	sensor_data.TransducerPitch = []float32{float32(buffer.TransducerPitch) / SCALE2}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.Power = []uint8{buffer.Power}
	sensor_data.Attenuation = []uint8{buffer.Attenuation}
	sensor_data.Tvg = []uint8{buffer.Tvg}
	sensor_data.PulseLength = []uint8{buffer.PulseLength}
	sensor_data.Counter = []uint16{buffer.Counter}

	return sensor_data
}

type Em950 struct {
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Quality              []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ShipPitch            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransducerPitch      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
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
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ValidBeams           []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PulseLength          []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	BeamWidth            []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TransmitPower        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TransmitStatus       []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ReceiveStatus        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm121ASpecific(reader *bytes.Reader) (sensor_data Em121A) {
	var buffer struct {
		PingNumber           uint16
		Mode                 uint8
		ValidBeams           uint8
		PulseLength          uint8
		BeamWidth            uint8
		TransmitPower        uint8
		TransmitStatus       uint8
		ReceiveStatus        uint8
		SurfaceSoundVelocity uint16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.ValidBeams = []uint8{buffer.ValidBeams}
	sensor_data.PulseLength = []uint8{buffer.PulseLength}
	sensor_data.BeamWidth = []uint8{buffer.BeamWidth}
	sensor_data.TransmitPower = []uint8{buffer.TransmitPower}
	sensor_data.TransmitStatus = []uint8{buffer.TransmitStatus}
	sensor_data.ReceiveStatus = []uint8{buffer.ReceiveStatus}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.ReceiveStatus) / 10.0}

	return sensor_data
}

type Em121 struct {
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ValidBeams           []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PulseLength          []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	BeamWidth            []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TransmitPower        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TransmitStatus       []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ReceiveStatus        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm121Specific(reader *bytes.Reader) (sensor_data Em121) {
	var buffer struct {
		PingNumber           uint16
		Mode                 uint8
		ValidBeams           uint8
		PulseLength          uint8
		BeamWidth            uint8
		TransmitPower        uint8
		TransmitStatus       uint8
		ReceiveStatus        uint8
		SurfaceSoundVelocity uint16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.ValidBeams = []uint8{buffer.ValidBeams}
	sensor_data.PulseLength = []uint8{buffer.PulseLength}
	sensor_data.BeamWidth = []uint8{buffer.BeamWidth}
	sensor_data.TransmitPower = []uint8{buffer.TransmitPower}
	sensor_data.TransmitStatus = []uint8{buffer.TransmitStatus}
	sensor_data.ReceiveStatus = []uint8{buffer.ReceiveStatus}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.ReceiveStatus) / 10.0}

	return sensor_data
}

type Sass struct {
	LeftMostBeam       []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RightMostBeam      []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TotalNumverOfBeams []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NavigationMode     []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingNumber         []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	MissionNumber      []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSassSpecfic(reader *bytes.Reader) (sensor_data Sass) {
	var buffer struct {
		LeftMostBeam       uint16
		RightMostBeam      uint16
		TotalNumverOfBeams uint16
		NavigationMode     uint16
		PingNumber         uint16
		MissionNumber      uint16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.LeftMostBeam = []uint16{buffer.LeftMostBeam}
	sensor_data.RightMostBeam = []uint16{buffer.RightMostBeam}
	sensor_data.TotalNumverOfBeams = []uint16{buffer.TotalNumverOfBeams}
	sensor_data.NavigationMode = []uint16{buffer.NavigationMode}
	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.MissionNumber = []uint16{buffer.MissionNumber}

	return sensor_data
}

type Seamap struct {
	PortTransmit1        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	PortTransmit2        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	StarboardTransmit1   []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	StarboardTransmit2   []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	PortGain             []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	StarboardGain        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	PortPulseLength      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	StarboardPulseLength []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	PressureDepth        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"` // only present in GSF >= 2.08
	Altitude             []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Temperature          []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSeamapSpecific(reader *bytes.Reader, gsfd GsfDetails) (sensor_data Seamap) {
	var (
		buffer1 struct {
			PortTransmit1        uint16
			PortTransmit2        uint16
			StarboardTransmit1   uint16
			StarboardTransmit2   uint16
			PortGain             uint16
			StarboardGain        uint16
			PortPulseLength      uint16
			StarboardPulseLength uint16
		}
		pressure_depth uint16 // only present in GSF >= 2.08
		buffer2        struct {
			Altitude    uint16
			Temperature uint16
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
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Range                []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TransmitPower        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ReceiveGain          []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSeabatSpecific(reader *bytes.Reader, gsfd GsfDetails) (sensor_data Seabat) {
	var buffer struct {
		PingNumber           uint16
		SurfaceSoundVelocity uint16
		Mode                 uint8
		Range                uint8
		TransmitPower        uint8
		ReceiveGain          uint8
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / 10.0}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.Range = []uint8{buffer.Range}
	sensor_data.TransmitPower = []uint8{buffer.TransmitPower}
	sensor_data.ReceiveGain = []uint8{buffer.ReceiveGain}

	return sensor_data
}

type Em1000 struct {
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Quality              []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ShipPitch            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransducerPitch      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm1000Specific(reader *bytes.Reader) (sensor_data Em1000) {
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

type TypeIIISeabeam struct {
	LeftMostBeam       []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RightMostBeam      []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TotalNumverOfBeams []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NavigationMode     []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingNumber         []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	MissionNumber      []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeTypeIIISeabeamSpecific(reader *bytes.Reader) (sensor_data TypeIIISeabeam) {
	var buffer struct {
		LeftMostBeam       uint16
		RightMostBeam      uint16
		TotalNumverOfBeams uint16
		NavigationMode     uint16
		PingNumber         uint16
		MissionNumber      uint16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.LeftMostBeam = []uint16{buffer.LeftMostBeam}
	sensor_data.RightMostBeam = []uint16{buffer.RightMostBeam}
	sensor_data.TotalNumverOfBeams = []uint16{buffer.TotalNumverOfBeams}
	sensor_data.NavigationMode = []uint16{buffer.NavigationMode}
	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.MissionNumber = []uint16{buffer.MissionNumber}

	return sensor_data
}

type SbAmp struct {
	Hour         []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Minute       []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Second       []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Hundredths   []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	BlockNumber  []uint32 `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	AvgGateDepth []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSbAmpSeabeamSpecific(reader *bytes.Reader) (sensor_data SbAmp) {
	var buffer struct {
		Hour         uint8
		Minute       uint8
		Second       uint8
		Hundredths   uint8
		BlockNumber  uint32
		AvgGateDepth uint16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.Hour = []uint8{buffer.Hour}
	sensor_data.Minute = []uint8{buffer.Minute}
	sensor_data.Second = []uint8{buffer.Second}
	sensor_data.Hundredths = []uint8{buffer.Hundredths}
	sensor_data.BlockNumber = []uint32{buffer.BlockNumber}
	sensor_data.AvgGateDepth = []uint16{buffer.AvgGateDepth}

	return sensor_data
}

type SeabatII struct {
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SonarRange           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TransmitPower        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ReceiveGain          []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ForeAftBandwidth     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	AthwartBandwidth     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSeabatIISpecific(reader *bytes.Reader) (sensor_data SeabatII) {
	var buffer struct {
		PingNumber           uint16
		SurfaceSoundVelocity uint16
		Mode                 uint16
		SonarRange           uint16
		TransmitPower        uint16
		ReceiveGain          uint16
		ForeAftBandwidth     uint8
		AthwartBandwidth     uint8
		Spare                int32
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / 10.0}
	sensor_data.Mode = []uint16{buffer.Mode}
	sensor_data.SonarRange = []uint16{buffer.SonarRange}
	sensor_data.TransmitPower = []uint16{buffer.TransmitPower}
	sensor_data.ReceiveGain = []uint16{buffer.ReceiveGain}
	sensor_data.ForeAftBandwidth = []float32{float32(buffer.ForeAftBandwidth) / 10.0}
	sensor_data.AthwartBandwidth = []float32{float32(buffer.AthwartBandwidth) / 10.0}

	return sensor_data
}

type Seabat8101 struct {
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Range                []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TransmitPower        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RecieveGain          []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PulseWidth           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TvgSpreading         []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TvgAbsorption        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ForeAftBandwidth     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	AthwartBandwidth     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RangeFilterMin       []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RangeFilterMax       []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	DepthFilterMin       []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	DepthFilterMax       []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorType        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSeabat8101Specific(reader *bytes.Reader) (sensor_data Seabat8101) {
	var buffer struct {
		PingNumber           uint16
		SurfaceSoundVelocity uint16
		Mode                 uint16
		Range                uint16
		TransmitPower        uint16
		RecieveGain          uint16
		PulseWidth           uint16
		TvgSpreading         uint8
		TvgAbsorption        uint8
		ForeAftBandwidth     uint8
		AthwartBandwidth     uint8
		RangeFilterMin       uint16
		RangeFilterMax       uint16
		DepthFilterMin       uint16
		DepthFilterMax       uint16
		ProjectorType        uint8
		Spare                int32
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / 10.0}
	sensor_data.Mode = []uint16{buffer.Mode}
	sensor_data.Range = []uint16{buffer.Range}
	sensor_data.TransmitPower = []uint16{buffer.TransmitPower}
	sensor_data.RecieveGain = []uint16{buffer.RecieveGain}
	sensor_data.PulseWidth = []uint16{buffer.PulseWidth}
	sensor_data.TvgSpreading = []uint8{buffer.TvgSpreading}
	sensor_data.TvgAbsorption = []uint8{buffer.TvgAbsorption}
	sensor_data.ForeAftBandwidth = []float32{float32(buffer.ForeAftBandwidth) / 10.0}
	sensor_data.AthwartBandwidth = []float32{float32(buffer.AthwartBandwidth) / 10.0}
	sensor_data.RangeFilterMin = []float32{float32(buffer.RangeFilterMin)}
	sensor_data.RangeFilterMax = []float32{float32(buffer.RangeFilterMax)}
	sensor_data.DepthFilterMin = []float32{float32(buffer.DepthFilterMin)}
	sensor_data.DepthFilterMax = []float32{float32(buffer.DepthFilterMax)}
	sensor_data.ProjectorType = []uint8{buffer.ProjectorType}

	return sensor_data
}

type Seabeam2112 struct {
	Mode                   []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity   []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SsvSource              []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PingGain               []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PulseWidth             []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TransmitterAttenuation []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	NumberAlgorithms       []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	AlgorithmOrder         []string  `tiledb:"dtype=string,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSeabeam2112Specific(reader *bytes.Reader) (sensor_data Seabeam2112) {
	var buffer struct {
		Mode                   uint8
		SurfaceSoundVelocity   uint16
		SsvSource              uint8
		PingGain               uint8
		PulseWidth             uint8
		TransmitterAttenuation uint8
		NumberAlgorithms       uint8
		AlgorithmOrder         [5]byte
		Spare                  int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.SurfaceSoundVelocity = []float32{(float32(buffer.SurfaceSoundVelocity) + 130000.0) / 100.0}
	sensor_data.SsvSource = []uint8{buffer.SsvSource}
	sensor_data.PingGain = []uint8{buffer.PingGain}
	sensor_data.PulseWidth = []uint8{buffer.PulseWidth}
	sensor_data.TransmitterAttenuation = []uint8{buffer.TransmitterAttenuation}
	sensor_data.NumberAlgorithms = []uint8{buffer.NumberAlgorithms}
	sensor_data.AlgorithmOrder = []string{string(buffer.AlgorithmOrder[:])}

	return sensor_data
}

type ElacMkII struct {
	Mode                  []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PingNumber            []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity  []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PulseLength           []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ReceiverGainStarboard []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ReceiverGainPort      []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeElacMkIISpecific(reader *bytes.Reader) (sensor_data ElacMkII) {
	var buffer struct {
		Mode                  uint8
		PingNumber            uint16
		SurfaceSoundVelocity  uint16
		PulseLength           uint16
		ReceiverGainStarboard uint8
		ReceiverGainPort      uint8
		Spare                 int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []uint16{buffer.SurfaceSoundVelocity}
	sensor_data.PulseLength = []uint16{buffer.PulseLength}
	sensor_data.ReceiverGainStarboard = []uint8{buffer.ReceiverGainStarboard}
	sensor_data.ReceiverGainPort = []uint8{buffer.ReceiverGainPort}

	return sensor_data
}

type CmpSass struct {
	Lfreq  []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Lntens []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeCmpSass(reader *bytes.Reader) (sensor_data CmpSass) {
	var buffer struct {
		Lfreq  uint16
		Lntens uint16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.Lfreq = []float32{float32(buffer.Lfreq) / 10.0}
	sensor_data.Lntens = []float32{float32(buffer.Lntens) / 10.0}

	return sensor_data
}

type Reson8100 struct {
	Latency              []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SonarID              []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SonarModel           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Frequency            []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SampleRate           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingRate             []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Range                []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TransmitPower        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ReceiveGain          []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PulseWidth           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TvgSpreading         []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TvgAbsorption        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ForeAftBandwidth     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	AthwartBandwidth     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorType        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ProjectorAngle       []int16   `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	RangeFilterMin       []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RangeFilterMax       []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	DepthFilterMin       []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	DepthFilterMax       []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	FiltersActive        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Temperature          []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	BeamSpacing          []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeReson8100(reader *bytes.Reader) (sensor_data Reson8100) {
	var buffer struct {
		Latency              uint16
		PingNumber           uint16
		SonarID              uint16
		SonarModel           uint16
		Frequency            uint16
		SurfaceSoundVelocity uint16
		SampleRate           uint16
		PingRate             uint16
		Mode                 uint16
		Range                uint16
		TransmitPower        uint16
		ReceiveGain          uint16
		PulseWidth           uint16
		TvgSpreading         uint8
		TvgAbsorption        uint8
		ForeAftBandwidth     uint8
		AthwartBandwidth     uint8
		ProjectorType        uint8
		ProjectorAngle       int16
		RangeFilterMin       uint16
		RangeFilterMax       uint16
		DepthFilterMin       uint16
		DepthFilterMax       uint16
		FiltersActive        uint8
		Temperature          uint16
		BeamSpacing          uint16
		Spare                int16
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.Latency = []uint16{buffer.Latency}
	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SonarID = []uint16{buffer.SonarID}
	sensor_data.SonarModel = []uint16{buffer.SonarModel}
	sensor_data.Frequency = []uint16{buffer.Frequency}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / 10.0}
	sensor_data.SampleRate = []uint16{buffer.SampleRate}
	sensor_data.PingRate = []uint16{buffer.PingRate}
	sensor_data.Mode = []uint16{buffer.Mode}
	sensor_data.Range = []uint16{buffer.Range}
	sensor_data.TransmitPower = []uint16{buffer.TransmitPower}
	sensor_data.ReceiveGain = []uint16{buffer.ReceiveGain}
	sensor_data.PulseWidth = []uint16{buffer.PulseWidth}
	sensor_data.TvgSpreading = []uint8{buffer.TvgSpreading}
	sensor_data.TvgAbsorption = []uint8{buffer.TvgAbsorption}
	sensor_data.ForeAftBandwidth = []float32{float32(buffer.ForeAftBandwidth) / 10.0}
	sensor_data.AthwartBandwidth = []float32{float32(buffer.AthwartBandwidth) / 10.0}
	sensor_data.ProjectorType = []uint8{buffer.ProjectorType}
	sensor_data.ProjectorAngle = []int16{buffer.ProjectorAngle}
	sensor_data.RangeFilterMin = []float32{float32(buffer.RangeFilterMin)}
	sensor_data.RangeFilterMax = []float32{float32(buffer.RangeFilterMax)}
	sensor_data.DepthFilterMin = []float32{float32(buffer.DepthFilterMin)}
	sensor_data.DepthFilterMax = []float32{float32(buffer.DepthFilterMax)}
	sensor_data.FiltersActive = []uint8{buffer.FiltersActive}
	sensor_data.Temperature = []uint16{buffer.Temperature}
	sensor_data.BeamSpacing = []float32{float32(buffer.BeamSpacing) / 10000.0}

	return sensor_data
}

type Em3 struct {
	ModelNumber          []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SerialNumber         []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransducerDepth      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ValidBeams           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SampleRate           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	DepthDifference      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	OffsetMultiplier     []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	// RunTimeID                      []uint32 // not stored
	RunTimeModelNumber             [][]uint16    `tiledb:"dtype=uint16,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeDgTime                  [][]time.Time `tiledb:"dtype=time,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimePingNumber              [][]uint16    `tiledb:"dtype=uint16,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeSerialNumber            [][]uint16    `tiledb:"dtype=uint16,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeSystemStatus            [][]uint32    `tiledb:"dtype=uint32,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeMode                    [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeFilterID                [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeMinDepth                [][]float32   `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeMaxDepth                [][]float32   `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeAbsorption              [][]float32   `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeTransmitPulseLength     [][]float32   `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeTransmitBeamWidth       [][]float32   `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimePowerReduction          [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeReceiveBeamWidth        [][]float32   `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeReceiveBandwidth        [][]int16     `tiledb:"dtype=int16,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeReceiveGain             [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeCrossOverAngle          [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeSsvSource               [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimePortSwathWidth          [][]uint16    `tiledb:"dtype=uint16,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeBeamSpacing             [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimePortCoverageSector      [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeStabilization           [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeStarboardCoverageSector [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeStarboardSwathWidth     [][]uint16    `tiledb:"dtype=uint16,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeHiloFreqAbsorpRatio     [][]uint8     `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeSwathWidth              [][]uint16    `tiledb:"dtype=uint16,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeCoverageSector          [][]uint16    `tiledb:"dtype=uint16,ftype=attr,var" filters:"zstd(level=16)"`
}

func DecodeEm3(reader *bytes.Reader) (sensor_data Em3) {
	var (
		buffer struct {
			ModelNumber          uint16
			PingNumber           uint16
			SerialNumber         uint16
			SurfaceSoundVelocity uint16
			TransducerDepth      uint16
			ValidBeams           uint16
			SampleRate           uint16
			DepthDifference      int16
			OffsetMultiplier     uint8
			RunTimeID            uint32
		}
		rt1 struct {
			ModelNumber             uint16
			TvSec                   uint32
			TvNSec                  uint32
			PingNumber              uint16
			SerialNumber            uint16
			SystemStatus            uint32
			Mode                    uint8
			FilterID                uint8
			MinDepth                uint16
			MaxDepth                uint16
			Absorption              uint16
			TransmitPulseLength     uint16
			TransmitBeamWidth       uint16
			PowerReduction          uint8
			ReceiveBeamWidth        uint8
			ReceiveBandwidth        uint8
			ReceiveGain             uint8
			CrossOverAnlge          uint8
			SsvSource               uint8
			PortSwathWidth          uint16
			BeamSpacing             uint8
			PortCoverageSector      uint8
			Stabilization           uint8
			StarboardCoverageSector uint8
			StarboardSwathWidth     uint16
			HiloFreqAbsorpRatio     uint8
			Spare                   int32
		}
		rt2 struct {
			ModelNumber             uint16
			TvSec                   uint32
			TvNSec                  uint32
			PingNumber              uint16
			SerialNumber            uint16
			SystemStatus            uint32
			Mode                    uint8
			FilterID                uint8
			MinDepth                uint16
			MaxDepth                uint16
			Absorption              uint16
			TransmitPulseLength     uint16
			TransmitBeamWidth       uint16
			PowerReduction          uint8
			ReceiveBeamWidth        uint8
			ReceiveBandwidth        uint8
			ReceiveGain             uint8
			CrossOverAnlge          uint8
			SsvSource               uint8
			PortSwathWidth          uint16
			BeamSpacing             uint8
			PortCoverageSector      uint8
			Stabilization           uint8
			StarboardCoverageSector uint8
			StarboardSwathWidth     uint16
			HiloFreqAbsorpRatio     uint8
			Spare                   int32
		}
	)
	model_number := make([]uint16, 0, 2)
	dg_time := make([]time.Time, 0, 2)
	ping_number := make([]uint16, 0, 2)
	serial_number := make([]uint16, 0, 2)
	system_status := make([]uint32, 0, 2)
	mode := make([]uint8, 0, 2)
	filter_id := make([]uint8, 0, 2)
	min_depth := make([]float32, 0, 2)
	max_depth := make([]float32, 0, 2)
	absorption := make([]float32, 0, 2)
	transmit_pulse_length := make([]float32, 0, 2)
	transmit_beam_width := make([]float32, 0, 2)
	power_reduction := make([]uint8, 0, 2)
	receive_beamwidth := make([]float32, 0, 2)
	receive_bandwidth := make([]int16, 0, 2)
	receive_gain := make([]uint8, 0, 2)
	cross_over_angle := make([]uint8, 0, 2)
	ssv_source := make([]uint8, 0, 2)
	port_swath_width := make([]uint16, 0, 2)
	beam_spacing := make([]uint8, 0, 2)
	port_coverage_sector := make([]uint8, 0, 2)
	stabilization := make([]uint8, 0, 2)
	starboard_coverage_sector := make([]uint8, 0, 2)
	starboard_swath_width := make([]uint16, 0, 2)
	hilo_freq_absorp_ratio := make([]uint8, 0, 2)
	swath_width := make([]uint16, 0, 2)
	coverage_sector := make([]uint16, 0, 2)

	_ = binary.Read(reader, binary.BigEndian, &buffer)

	// base set of values
	sensor_data.ModelNumber = []uint16{buffer.ModelNumber}
	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SerialNumber = []uint16{buffer.SerialNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SerialNumber) / 10.0}
	sensor_data.TransducerDepth = []float32{float32(buffer.TransducerDepth) / SCALE2}
	sensor_data.ValidBeams = []uint16{buffer.ValidBeams}
	sensor_data.SampleRate = []uint16{buffer.SampleRate}
	sensor_data.DepthDifference = []float32{float32(buffer.DepthDifference) / SCALE2}
	sensor_data.OffsetMultiplier = []uint8{buffer.OffsetMultiplier}

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
		absorption = append(absorption, float32(rt2.Absorption)/SCALE2)
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
			coverage_sector = append(coverage_sector, uint16(rt1.PortCoverageSector)+uint16(rt1.StarboardCoverageSector))
			port_coverage_sector = append(port_coverage_sector, rt1.PortCoverageSector)
			starboard_coverage_sector = append(starboard_coverage_sector, rt1.StarboardCoverageSector)
		} else {
			coverage_sector = append(coverage_sector, uint16(rt1.PortCoverageSector))
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
			absorption = append(absorption, float32(rt2.Absorption)/SCALE2)
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
			coverage_sector = append(coverage_sector, uint16(rt2.PortCoverageSector)+uint16(rt2.StarboardCoverageSector))
			port_coverage_sector = append(port_coverage_sector, rt2.PortCoverageSector)
			starboard_coverage_sector = append(starboard_coverage_sector, rt2.StarboardCoverageSector)
		} else {
			coverage_sector = append(coverage_sector, uint16(rt2.PortCoverageSector))
			port_coverage_sector = append(port_coverage_sector, rt2.PortCoverageSector/2)
			starboard_coverage_sector = append(starboard_coverage_sector, rt2.PortCoverageSector/2)
		}
	}

	// insert runtime data
	sensor_data.RunTimeModelNumber = [][]uint16{model_number}
	sensor_data.RunTimeDgTime = [][]time.Time{dg_time}
	sensor_data.RunTimePingNumber = [][]uint16{ping_number}
	sensor_data.RunTimeSerialNumber = [][]uint16{serial_number}
	sensor_data.RunTimeMode = [][]uint8{mode}
	sensor_data.RunTimeFilterID = [][]uint8{filter_id}
	sensor_data.RunTimeMinDepth = [][]float32{min_depth}
	sensor_data.RunTimeMaxDepth = [][]float32{max_depth}
	sensor_data.RunTimeAbsorption = [][]float32{absorption}
	sensor_data.RunTimeTransmitPulseLength = [][]float32{transmit_pulse_length}
	sensor_data.RunTimeTransmitBeamWidth = [][]float32{transmit_beam_width}
	sensor_data.RunTimePowerReduction = [][]uint8{power_reduction}
	sensor_data.RunTimeReceiveBeamWidth = [][]float32{receive_beamwidth}
	sensor_data.RunTimeReceiveBandwidth = [][]int16{receive_bandwidth}
	sensor_data.RunTimeReceiveGain = [][]uint8{receive_gain}
	sensor_data.RunTimeCrossOverAngle = [][]uint8{cross_over_angle}
	sensor_data.RunTimeSsvSource = [][]uint8{ssv_source}
	sensor_data.RunTimePortSwathWidth = [][]uint16{port_swath_width}
	sensor_data.RunTimeBeamSpacing = [][]uint8{beam_spacing}
	sensor_data.RunTimePortCoverageSector = [][]uint8{port_coverage_sector}
	sensor_data.RunTimeStabilization = [][]uint8{stabilization}
	sensor_data.RunTimeStarboardCoverageSector = [][]uint8{starboard_coverage_sector}
	sensor_data.RunTimeStarboardSwathWidth = [][]uint16{starboard_swath_width}
	sensor_data.RunTimeHiloFreqAbsorpRatio = [][]uint8{hilo_freq_absorp_ratio}
	sensor_data.RunTimeSwathWidth = [][]uint16{swath_width}
	sensor_data.RunTimeCoverageSector = [][]uint16{coverage_sector}

	return sensor_data
}

// float64 may be overkill
// where scale factors are applied, float32 is used
// where it's confident float32 is enough to represent the value
// TODO; align into 64bit chunks
// the spec says binary integers are stored as either 1-byte unsigned, 2-byte signed or unsigned, or 4-byte signed
type Em4 struct {
	ModelNumber                       []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingCounter                       []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SerialNumber                      []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceVelocity                   []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransducerDepth                   []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ValidDetections                   []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SamplingFrequency                 []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	DopplerCorrectionScale            []uint32    `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	VehicleDepth                      []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransmitSectors                   []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TiltAngle                         [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	FocusRange                        [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	SignalLength                      [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	TransmitDelay                     [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	CenterFrequency                   [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	MeanAbsorption                    [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	WaveformId                        [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	SectorNumber                      [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	SignalBandwith                    [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeModelNumber                []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeDatagramTime               []time.Time `tiledb:"dtype=datetime_ns,ftype=attr" filters:"zstd(level=16)"`
	RunTimePingCounter                []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeSerialNumber               []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeOperatorStationStatus      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeProcessingUnitStatus       []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeBspStatus                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeHeadTransceiverStatus      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMode                       []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeFilterID                   []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
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
	RunTimeMaxStdbSwathWidth          []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTransmitAlongTilt          []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeFilterID2                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ProcessorUnitCpuLoad              []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ProcessorUnitSensorStatus         []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ProcessorUnitAchievedPortCoverage []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ProcessorUnitAchievedStbdCoverage []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ProcessorUnitYawStabilization     []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm4Specific(reader *bytes.Reader) (sensor_data Em4) {

	var (
		buffer struct {
			ModelNumber            uint16
			PingCounter            uint16
			SerialNumber           uint16
			SurfaceVelocity        uint16
			TransducerDepth        int32
			ValidDetections        uint16
			SamplingFrequency1     uint32
			SamplingFrequency2     uint32
			DopplerCorrectionScale uint32
			VehicleDepth           int32
			Spare                  [4]int32
			TransmitSectors        uint16
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
			FocusRange      uint16
			SignalLength    uint32
			TransmitDelay   uint32
			CenterFrequency uint32
			MeanAbsorption  uint16
			WaveformId      uint8
			SectorNumber    uint8
			SignalBandwith  uint32
			Spare           [4]int32
		} // 40 bytes
		spare_buffer struct {
			Spare [4]int32
		} // 16 bytes
		runtime_buffer struct {
			RunTimeModelNumber            uint16
			RunTimeDatagramTime_sec       uint32
			RunTimeDatagramTime_nsec      uint32
			RunTimePingCounter            uint16
			RunTimeSerialNumber           uint16
			RunTimeOperatorStationStatus  uint8
			RunTimeProcessingUnitStatus   uint8
			RunTimeBspStatus              uint8
			RunTimeHeadTransceiverStatus  uint8
			RunTimeMode                   uint8
			RunTimeFilterID               uint8
			RunTimeMinDepth               uint16
			RunTimeMaxDepth               uint16
			RunTimeAbsorption             uint16
			RunTimeTransmitPulseLength    uint16
			RunTimeTransmitBeamWidth      uint16
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
			RunTimeMaxStdbSwathWidth      uint16
			RunTimeTransmitAlongTilt      int16
			RunTimeFilterID2              uint8
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
	for i := uint16(0); i < buffer.TransmitSectors; i++ {
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
	sensor_data.ModelNumber = []uint16{buffer.ModelNumber}
	sensor_data.PingCounter = []uint16{buffer.PingCounter}
	sensor_data.SerialNumber = []uint16{buffer.SerialNumber}
	sensor_data.SurfaceVelocity = []float32{float32(buffer.SurfaceVelocity) / float32(10)}
	sensor_data.TransducerDepth = []float64{float64(buffer.TransducerDepth) / float64(20000)}
	sensor_data.ValidDetections = []uint16{buffer.ValidDetections}
	sensor_data.SamplingFrequency = []float64{float64(buffer.SamplingFrequency1) + float64(buffer.SamplingFrequency2)/float64(4_000_000_000)}
	sensor_data.DopplerCorrectionScale = []uint32{buffer.DopplerCorrectionScale}
	sensor_data.VehicleDepth = []float32{float32(buffer.VehicleDepth) / float32(1000)}
	sensor_data.TransmitSectors = []uint16{buffer.TransmitSectors}

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
	sensor_data.RunTimeModelNumber = []uint16{runtime_buffer.RunTimeModelNumber}
	sensor_data.RunTimeDatagramTime = []time.Time{time.Unix(
		int64(runtime_buffer.RunTimeDatagramTime_sec),
		int64(runtime_buffer.RunTimeDatagramTime_nsec),
	)}
	sensor_data.RunTimePingCounter = []uint16{runtime_buffer.RunTimePingCounter}
	sensor_data.RunTimeSerialNumber = []uint16{runtime_buffer.RunTimeSerialNumber}
	sensor_data.RunTimeOperatorStationStatus = []uint8{runtime_buffer.RunTimeOperatorStationStatus}
	sensor_data.RunTimeProcessingUnitStatus = []uint8{runtime_buffer.RunTimeProcessingUnitStatus}
	sensor_data.RunTimeBspStatus = []uint8{runtime_buffer.RunTimeBspStatus}
	sensor_data.RunTimeHeadTransceiverStatus = []uint8{runtime_buffer.RunTimeHeadTransceiverStatus}
	sensor_data.RunTimeMode = []uint8{runtime_buffer.RunTimeMode}
	sensor_data.RunTimeFilterID = []uint8{runtime_buffer.RunTimeFilterID}
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
	sensor_data.RunTimeMaxStdbSwathWidth = []uint16{runtime_buffer.RunTimeMaxStdbSwathWidth}
	sensor_data.RunTimeTransmitAlongTilt = []float32{float32(runtime_buffer.RunTimeTransmitAlongTilt) / float32(100)}
	sensor_data.RunTimeFilterID2 = []uint8{runtime_buffer.RunTimeFilterID2}

	// populate processor unit info
	sensor_data.ProcessorUnitCpuLoad = []uint8{proc_buffer.ProcessorUnitCpuLoad}
	sensor_data.ProcessorUnitSensorStatus = []uint16{proc_buffer.ProcessorUnitSensorStatus}
	sensor_data.ProcessorUnitAchievedPortCoverage = []uint8{proc_buffer.ProcessorUnitAchievedPortCoverage}
	sensor_data.ProcessorUnitAchievedStbdCoverage = []uint8{proc_buffer.ProcessorUnitAchievedStbdCoverage}
	sensor_data.ProcessorUnitYawStabilization = []float32{float32(proc_buffer.ProcessorUnitYawStabilization) / float32(100)}

	return sensor_data
}

// GeoSwathPlus TODO; change DataSource and Side types from uint16 to uint8.
// Seems a waste to store 2 bytes for data that is only a 0 or 1.
type GeoSwathPlus struct {
	// (0 = CBF, 1 = RDF) why 2bytes? why not uint8? could convert to string ...
	DataSource []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	// again why 2bytes for (0 port, 1 = stbd)
	Side                  []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ModelNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Frequency             []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	EchosounderType       []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingNumber            []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	NumNavSamples         []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumAttitudeSamples    []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumHeadingSamples     []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumMiniSvsSamples     []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumEchosounderSamples []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumRaaSamples         []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	MeanSv                []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SurfaceVelocity       []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ValidBeams            []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SampleRate            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	PulseLength           []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	PingLength            []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TransmitPower         []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SidescanGainChannel   []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Stabilization         []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	GpsQuality            []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RangeUncertainty      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	AngleUncertainty      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeGeoSwathPlusSpecific(reader *bytes.Reader) (sensor_data GeoSwathPlus) {
	var buffer struct {
		DataSource            uint16
		Side                  uint16
		ModelNumber           uint16
		Frequency             uint16
		EchosounderType       uint16
		PingNumber            uint32
		NumNavSamples         uint16
		NumAttitudeSamples    uint16
		NumHeadingSamples     uint16
		NumMiniSvsSamples     uint16
		NumEchosounderSamples uint16
		NumRaaSamples         uint16
		MeanSv                uint16
		SurfaceVelocity       uint16
		ValidBeams            uint16
		SampleRate            float32
		PulseLength           float32
		PingLength            uint16
		TransmitPower         uint16
		SidescanGainChannel   uint16
		Stabilization         uint16
		GpsQuality            uint16
		RangeUncertainty      float32
		AngleUncertainty      float32
		Spare                 [4]int32
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.DataSource = []uint16{buffer.DataSource}
	sensor_data.Side = []uint16{buffer.Side}
	sensor_data.ModelNumber = []uint16{buffer.ModelNumber}
	sensor_data.Frequency = []float32{float32(buffer.Frequency) / 10.0}
	sensor_data.EchosounderType = []uint16{buffer.EchosounderType}
	sensor_data.PingNumber = []uint32{buffer.PingNumber}
	sensor_data.NumNavSamples = []uint16{buffer.NumNavSamples}
	sensor_data.NumAttitudeSamples = []uint16{buffer.NumAttitudeSamples}
	sensor_data.NumHeadingSamples = []uint16{buffer.NumHeadingSamples}
	sensor_data.NumMiniSvsSamples = []uint16{buffer.NumMiniSvsSamples}
	sensor_data.NumEchosounderSamples = []uint16{buffer.NumEchosounderSamples}
	sensor_data.NumRaaSamples = []uint16{buffer.NumRaaSamples}
	sensor_data.MeanSv = []float32{float32(buffer.MeanSv) / 20.0}
	sensor_data.SurfaceVelocity = []float32{float32(buffer.SurfaceVelocity) / 20.0}
	sensor_data.ValidBeams = []uint16{buffer.ValidBeams}
	sensor_data.SampleRate = []float32{float32(buffer.SampleRate) / 10.0}
	sensor_data.PulseLength = []float32{float32(buffer.PulseLength)}
	sensor_data.PingLength = []uint16{buffer.PingLength}
	sensor_data.TransmitPower = []uint16{buffer.TransmitPower}
	sensor_data.SidescanGainChannel = []uint16{buffer.SidescanGainChannel}
	sensor_data.Stabilization = []uint16{buffer.Stabilization}
	sensor_data.GpsQuality = []uint16{buffer.GpsQuality}
	sensor_data.RangeUncertainty = []float32{float32(buffer.RangeUncertainty) / SCALE3}
	sensor_data.AngleUncertainty = []float32{float32(buffer.AngleUncertainty) / SCALE2}

	return sensor_data
}

// DecodeKlein5410Bss TODO; change DataSource and Side types from uint16 to uint8.
// Seems a waste to store 2 bytes for data that is only a 0 or 1.
type Klein5410Bss struct {
	DataSource        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Side              []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ModelNumber       []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	AcousticFrequency []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SamplingFrequency []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	PingNumber        []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	NumSamples        []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	NumRaaSamples     []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ErrorFlags        []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	Range             []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	FishDepth         []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	FishAltitude      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SoundSpeed        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransmitWaveform  []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Altimeter         []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RawDataConfig     []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeKlein5410BssSpecific(reader *bytes.Reader) (sensor_data Klein5410Bss) {
	var buffer struct {
		DataSource        uint16
		Side              uint16
		ModelNumber       uint16
		AcousticFrequency uint32
		SamplingFrequency uint32
		PingNumber        uint32
		NumSamples        uint32
		NumRaaSamples     uint32
		ErrorFlags        uint32
		Range             uint32
		FishDepth         uint32
		FishAltitude      uint32
		SoundSpeed        uint32
		TransmitWaveform  uint16
		Altimeter         uint16
		RawDataConfig     uint32
		Spare             [4]int32
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.DataSource = []uint16{buffer.DataSource}
	sensor_data.Side = []uint16{buffer.Side}
	sensor_data.ModelNumber = []uint16{buffer.ModelNumber}
	sensor_data.AcousticFrequency = []float32{float32(buffer.AcousticFrequency) / SCALE3}
	sensor_data.SamplingFrequency = []float32{float32(buffer.SamplingFrequency) / SCALE3}
	sensor_data.PingNumber = []uint32{buffer.PingNumber}
	sensor_data.NumSamples = []uint32{buffer.NumSamples}
	sensor_data.NumRaaSamples = []uint32{buffer.NumRaaSamples}
	sensor_data.ErrorFlags = []uint32{buffer.ErrorFlags}
	sensor_data.Range = []uint32{buffer.Range}
	sensor_data.FishDepth = []float32{float32(buffer.FishDepth) / SCALE3}
	sensor_data.FishAltitude = []float32{float32(buffer.FishAltitude) / SCALE3}
	sensor_data.SoundSpeed = []float32{float32(buffer.SoundSpeed) / SCALE3}
	sensor_data.TransmitWaveform = []uint16{buffer.TransmitWaveform}
	sensor_data.Altimeter = []uint16{buffer.Altimeter}
	sensor_data.RawDataConfig = []uint32{buffer.RawDataConfig}

	return sensor_data
}

type Reson7100 struct {
	ProtocolVersion                   []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	DeviceID                          []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	MajorSerialNumber                 []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	MinorSerialNumber                 []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	PingNumber                        []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	MultiPingSequence                 []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Frequency                         []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SampleRate                        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ReceiverBandwidth                 []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TxPulseWidth                      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TxPulseTypeID                     []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	TxPulseEnvlpID                    []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	TxPulseEnvlpParam                 []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	MaxPingRate                       []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	PingPeriod                        []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Range                             []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Power                             []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Gain                              []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ControlFlags                      []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorID                       []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorSteerAnglVert            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorSteerAnglHorz            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamWidthVert            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamWidthHorz            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamFocalPt              []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamWeightingWindowType  []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamWeightingWindowParam []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	TransmitFlags                     []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	HydrophoneID                      []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ReceivingBeamWeightingWindowType  []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ReceivingBeamWeightingWindowParam []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ReceiveFlags                      []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ReceiveBeamWidth                  []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RangeFiltMin                      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RangeFiltMax                      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	DepthFiltMin                      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	DepthFiltMax                      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Absorption                        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SoundVelocity                     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Spreading                         []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RawDataFrom7027                   []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SvSource                          []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	LayerCompFlag                     []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	// TxPulseReserved                   []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeReson7100Specific(reader *bytes.Reader) (sensor_data Reson7100) {
	var buffer struct {
		ProtocolVersion                   uint16
		DeviceID                          uint32
		Reserved1                         [16]byte
		MajorSerialNumber                 uint32
		MinorSerialNumber                 uint32
		PingNumber                        uint32
		MultiPingSequence                 uint16
		Frequency                         uint32
		SampleRate                        uint32
		ReceiverBandwidth                 uint32
		TxPulseWidth                      uint32
		TxPulseTypeID                     uint32
		TxPulseEnvlpID                    uint32
		TxPulseEnvlpParam                 uint32
		TxPulseReserved                   uint32
		MaxPingRate                       uint32
		PingPeriod                        uint32
		Range                             uint32
		Power                             uint32
		Gain                              int32
		ControlFlags                      uint32
		ProjectorID                       uint32
		ProjectorSteerAnglVert            int32
		ProjectorSteerAnglHorz            int32
		ProjectorBeamWidthVert            uint16
		ProjectorBeamWidthHorz            uint16
		ProjectorBeamFocalPt              uint32
		ProjectorBeamWeightingWindowType  uint32
		ProjectorBeamWeightingWindowParam uint32
		TransmitFlags                     uint32
		HydrophoneID                      uint32
		ReceivingBeamWeightingWindowType  uint32
		ReceivingBeamWeightingWindowParam uint32
		ReceiveFlags                      uint32
		ReceiveBeamWidth                  uint16
		RangeFiltMin                      uint16
		RangeFiltMax                      uint16
		DepthFiltMin                      uint16
		DepthFiltMax                      uint16
		Absorption                        uint32
		SoundVelocity                     uint16
		Spreading                         uint32
		RawDataFrom7027                   uint8
		Reserved2                         [15]byte
		SvSource                          uint8
		LayerCompFlag                     uint8
		Reserved3                         [8]byte
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.ProtocolVersion = []uint16{buffer.ProtocolVersion}
	sensor_data.DeviceID = []uint32{buffer.DeviceID}
	sensor_data.MajorSerialNumber = []uint32{buffer.MajorSerialNumber}
	sensor_data.MinorSerialNumber = []uint32{buffer.MinorSerialNumber}
	sensor_data.PingNumber = []uint32{buffer.PingNumber}
	sensor_data.MultiPingSequence = []uint16{buffer.MultiPingSequence}
	sensor_data.Frequency = []float32{float32(buffer.Frequency) / SCALE3}
	sensor_data.SampleRate = []float32{float32(buffer.SampleRate) / 10_000.0}
	sensor_data.ReceiverBandwidth = []float32{float32(buffer.ReceiverBandwidth) / 10_000.0}
	sensor_data.TxPulseWidth = []float32{float32(buffer.TxPulseWidth) / SCALE1}
	sensor_data.TxPulseTypeID = []uint32{buffer.TxPulseTypeID}
	sensor_data.TxPulseEnvlpID = []uint32{buffer.TxPulseEnvlpID}
	sensor_data.TxPulseEnvlpParam = []float32{float32(buffer.TxPulseEnvlpParam) / SCALE2}
	sensor_data.MaxPingRate = []float64{float64(buffer.MaxPingRate) / 1_000_000.0}
	sensor_data.PingPeriod = []float64{float64(buffer.PingPeriod) / 1_000_000.0}
	sensor_data.Range = []float32{float32(buffer.Range) / SCALE2}
	sensor_data.Power = []float32{float32(buffer.Power) / SCALE2}
	sensor_data.Gain = []float32{float32(buffer.Gain) / SCALE2}
	sensor_data.ControlFlags = []uint32{buffer.ControlFlags}
	sensor_data.ProjectorID = []uint32{buffer.ProjectorID}
	sensor_data.ProjectorSteerAnglVert = []float32{float32(buffer.ProjectorSteerAnglVert) / SCALE3}
	sensor_data.ProjectorSteerAnglHorz = []float32{float32(buffer.ProjectorSteerAnglHorz) / SCALE3}
	sensor_data.ProjectorBeamWidthVert = []float32{float32(buffer.ProjectorBeamWidthVert) / SCALE2}
	sensor_data.ProjectorBeamWidthHorz = []float32{float32(buffer.ProjectorBeamWidthHorz) / SCALE2}
	sensor_data.ProjectorBeamFocalPt = []float32{float32(buffer.ProjectorBeamFocalPt) / SCALE2}
	sensor_data.ProjectorBeamWeightingWindowType = []uint32{buffer.ProjectorBeamWeightingWindowType}
	sensor_data.ProjectorBeamWeightingWindowParam = []uint32{buffer.ProjectorBeamWeightingWindowParam}
	sensor_data.TransmitFlags = []uint32{buffer.TransmitFlags}
	sensor_data.HydrophoneID = []uint32{buffer.HydrophoneID}
	sensor_data.ReceivingBeamWeightingWindowType = []uint32{buffer.ReceivingBeamWeightingWindowType}
	sensor_data.ReceivingBeamWeightingWindowParam = []uint32{buffer.ReceivingBeamWeightingWindowParam}
	sensor_data.ReceiveFlags = []uint32{buffer.ReceiveFlags}
	sensor_data.ReceiveBeamWidth = []float32{float32(buffer.ReceiveBeamWidth) / SCALE2}
	sensor_data.RangeFiltMin = []float32{float32(buffer.RangeFiltMin) / 10.0}
	sensor_data.RangeFiltMax = []float32{float32(buffer.RangeFiltMax) / 10.0}
	sensor_data.DepthFiltMin = []float32{float32(buffer.DepthFiltMin) / 10.0}
	sensor_data.DepthFiltMax = []float32{float32(buffer.DepthFiltMax) / 10.0}
	sensor_data.Absorption = []float32{float32(buffer.Absorption) / SCALE3}
	sensor_data.SoundVelocity = []float32{float32(buffer.SoundVelocity) / 10.0}
	sensor_data.Spreading = []float32{float32(buffer.Spreading) / SCALE3}
	sensor_data.RawDataFrom7027 = []uint8{buffer.RawDataFrom7027}
	sensor_data.SvSource = []uint8{buffer.SvSource}
	sensor_data.LayerCompFlag = []uint8{buffer.LayerCompFlag}

	return sensor_data
}

type Em3Raw struct {
	ModelNumber                       []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingCounter                       []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SerialNumber                      []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceVelocity                   []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransducerDepth                   []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ValidDetections                   []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SamplingFrequency                 []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	VehicleDepth                      []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	DepthDifference                   []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	OffsetMultiplier                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TransmitSectors                   []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TiltAngle                         [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	FocusRange                        [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	SignalLength                      [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	TransmitDelay                     [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	CenterFrequency                   [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	WaveformID                        [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	SectorNumber                      [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	SignalBandwidth                   [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeModelNumber                []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeDgTime                     []time.Time `tiledb:"dtype=time,ftype=attr" filters:"zstd(level=16)"`
	RunTimePingCounter                []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeSerialNumber               []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeOperatorStationStatus      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeProcessingUnitStatus       []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeBspStatus                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeHeadTransceiverStatus      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMode                       []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeFilterID                   []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMinDepth                   []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMaxDepth                   []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeAbsorption                 []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTxPulseLength              []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTxBeamWidth                []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTxPowerReMax               []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeRxBeamWidth                []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeRxBandwidth                []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeRxFixedGain                []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTvgCrossOverAngle          []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeSsvSource                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMaxPortSwathWidth          []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeBeamSpacing                []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMaxPortCoverage            []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeStabilization              []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMaxStarboardCoverage       []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeMaxStarboardSwathWidth     []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeDurotongSpeed              []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeTxAlongTilt                []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RunTimeHiLoAbsorptionRatio        []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PuStatusPuCpuLoad                 []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PuStatusSensorStatus              []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PuStatusAchievedPortCoverage      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PuStatusAchievedStarboardCoverage []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PuStatusYawStabilization          []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm3RawSpecific(reader *bytes.Reader) (sensor_data Em3Raw) {
	var (
		buffer struct {
			ModelNumber        uint16
			PingCounter        uint16
			SerialNumber       uint16
			SurfaceVelocity    uint16
			TransducerDepth    int32
			ValidDetections    uint16
			SamplingFrequency1 uint32
			SamplingFrequency2 uint32
			VehicleDepth       int32
			DepthDifference    int16
			OffsetMultiplier   uint8
			Spare              [16]byte
			TransmitSectors    uint16
		}
		var_buff struct {
			TiltAngle       int16
			FocusRange      uint16
			SignalLength    uint32
			TransmitDelay   uint32
			CenterFrequency uint32
			WaveformID      uint8
			SectorNumber    uint8
			SignalBandwidth uint32
			Spare           [16]byte
		}
		rt_buff struct {
			Spare                         [16]byte
			RunTimeModelNumber            uint16
			RunTimeDgTimeSec              uint32
			RunTimeDgTimeNSec             uint32
			RunTimePingCounter            uint16
			RunTimeSerialNumber           uint16
			RunTimeOperatorStationStatus  uint8
			RunTimeProcessingUnitStatus   uint8
			RunTimeBspStatus              uint8
			RunTimeHeadTransceiverStatus  uint8
			RunTimeMode                   uint8
			RunTimeFilterID               uint8
			RunTimeMinDepth               uint16
			RunTimeMaxDepth               uint16
			RunTimeAbsorption             uint16
			RunTimeTxPulseLength          uint16
			RunTimeTxBeamWidth            uint16
			RunTimeTxPowerReMax           uint8
			RunTimeRxBeamWidth            uint8
			RunTimeRxBandwidth            uint8
			RunTimeRxFixedGain            uint8
			RunTimeTvgCrossOverAngle      uint8
			RunTimeSsvSource              uint8
			RunTimeMaxPortSwathWidth      int16
			RunTimeBeamSpacing            uint8
			RunTimeMaxPortCoverage        uint8
			RunTimeStabilization          uint8
			RunTimeMaxStarboardCoverage   uint8
			RunTimeMaxStarboardSwathWidth uint16
		}
		RunTimeDurotongSpeed       uint16
		RunTimeTxAlongTilt         int16
		Spare                      [2]byte
		RunTimeHiLoAbsorptionRatio uint8
		pu_buff                    struct {
			Spare1                            [16]byte
			PuStatusPuCpuLoad                 uint8
			PuStatusSensorStatus              uint16
			PuStatusAchievedPortCoverage      uint8
			PuStatusAchievedStarboardCoverage uint8
			PuStatusYawStabilization          int16
			Spare2                            [16]byte
		}
	)

	// first block
	_ = binary.Read(reader, binary.BigEndian, &buffer)
	sensor_data.ModelNumber = []uint16{buffer.ModelNumber}
	sensor_data.PingCounter = []uint16{buffer.PingCounter}
	sensor_data.SerialNumber = []uint16{buffer.SerialNumber}
	sensor_data.SurfaceVelocity = []float32{float32(buffer.SurfaceVelocity) / 10.0}
	sensor_data.TransducerDepth = []float32{float32(buffer.TransducerDepth) / 20_000.0}
	sensor_data.ValidDetections = []uint16{buffer.ValidDetections}
	sensor_data.SamplingFrequency = []float64{float64(buffer.SamplingFrequency1) + float64(buffer.SamplingFrequency2)/4_000_000_000.0}
	sensor_data.VehicleDepth = []float32{float32(buffer.VehicleDepth) / SCALE3}
	sensor_data.DepthDifference = []float32{float32(buffer.DepthDifference) / SCALE2}
	sensor_data.OffsetMultiplier = []uint8{buffer.OffsetMultiplier}
	sensor_data.TransmitSectors = []uint16{buffer.TransmitSectors}

	// second block (variable length arrays)
	nsectors := int(buffer.TransmitSectors)
	tilt_angle := make([]float32, 0, nsectors)
	focus_range := make([]float32, 0, nsectors)
	signal_length := make([]float64, 0, nsectors)
	transmit_delay := make([]float64, 0, nsectors)
	centre_frequency := make([]float32, 0, nsectors)
	waveformID := make([]uint8, 0, nsectors)
	sector_number := make([]uint8, 0, nsectors)
	signal_bandwidth := make([]float32, 0, nsectors)

	for i := 0; i < nsectors; i++ {
		_ = binary.Read(reader, binary.BigEndian, &var_buff)
		tilt_angle = append(tilt_angle, float32(var_buff.TiltAngle)/SCALE2)
		focus_range = append(focus_range, float32(var_buff.FocusRange)/10.0)
		signal_length = append(signal_length, float64(var_buff.SignalLength)/1_000_000.0)
		transmit_delay = append(transmit_delay, float64(var_buff.TransmitDelay)/1_000_000.0)
		centre_frequency = append(centre_frequency, float32(var_buff.CenterFrequency)/SCALE3)
		waveformID = append(waveformID, var_buff.WaveformID)
		sector_number = append(sector_number, var_buff.SectorNumber)
		signal_bandwidth = append(signal_bandwidth, float32(var_buff.SignalBandwidth)/SCALE3)
	}

	sensor_data.TiltAngle = [][]float32{tilt_angle}
	sensor_data.FocusRange = [][]float32{focus_range}
	sensor_data.SignalLength = [][]float64{signal_length}
	sensor_data.TransmitDelay = [][]float64{transmit_delay}
	sensor_data.CenterFrequency = [][]float32{centre_frequency}
	sensor_data.WaveformID = [][]uint8{waveformID}
	sensor_data.SectorNumber = [][]uint8{sector_number}
	sensor_data.SignalBandwidth = [][]float32{signal_bandwidth}

	// third block (runtime)
	_ = binary.Read(reader, binary.BigEndian, &rt_buff)
	sensor_data.RunTimeModelNumber = []uint16{rt_buff.RunTimeModelNumber}
	sensor_data.RunTimeDgTime = []time.Time{time.Unix(int64(rt_buff.RunTimeDgTimeSec), int64(rt_buff.RunTimeDgTimeNSec)).UTC()}
	sensor_data.RunTimePingCounter = []uint16{rt_buff.RunTimePingCounter}
	sensor_data.RunTimeSerialNumber = []uint16{rt_buff.RunTimeSerialNumber}
	sensor_data.RunTimeOperatorStationStatus = []uint8{rt_buff.RunTimeOperatorStationStatus}
	sensor_data.RunTimeProcessingUnitStatus = []uint8{rt_buff.RunTimeProcessingUnitStatus}
	sensor_data.RunTimeBspStatus = []uint8{rt_buff.RunTimeBspStatus}
	sensor_data.RunTimeHeadTransceiverStatus = []uint8{rt_buff.RunTimeHeadTransceiverStatus}
	sensor_data.RunTimeMode = []uint8{rt_buff.RunTimeMode}
	sensor_data.RunTimeFilterID = []uint8{rt_buff.RunTimeFilterID}
	sensor_data.RunTimeMinDepth = []uint16{rt_buff.RunTimeMinDepth}
	sensor_data.RunTimeMaxDepth = []uint16{rt_buff.RunTimeMaxDepth}
	sensor_data.RunTimeAbsorption = []float32{float32(rt_buff.RunTimeAbsorption) / SCALE2}
	sensor_data.RunTimeTxPulseLength = []uint16{rt_buff.RunTimeTxPulseLength}
	sensor_data.RunTimeTxBeamWidth = []float32{float32(rt_buff.RunTimeTxBeamWidth) / 10.0}
	sensor_data.RunTimeTxPowerReMax = []uint8{rt_buff.RunTimeTxPowerReMax}
	sensor_data.RunTimeRxBeamWidth = []float32{float32(rt_buff.RunTimeRxBeamWidth) / 10.0}
	sensor_data.RunTimeRxBandwidth = []float32{float32(rt_buff.RunTimeRxBandwidth) * 50.0}
	sensor_data.RunTimeRxFixedGain = []uint8{rt_buff.RunTimeRxFixedGain}
	sensor_data.RunTimeTvgCrossOverAngle = []uint8{rt_buff.RunTimeTvgCrossOverAngle}
	sensor_data.RunTimeSsvSource = []uint8{rt_buff.RunTimeSsvSource}
	sensor_data.RunTimeMaxPortSwathWidth = []int16{rt_buff.RunTimeMaxPortSwathWidth}
	sensor_data.RunTimeBeamSpacing = []uint8{rt_buff.RunTimeBeamSpacing}
	sensor_data.RunTimeMaxPortCoverage = []uint8{rt_buff.RunTimeMaxPortCoverage}
	sensor_data.RunTimeStabilization = []uint8{rt_buff.RunTimeStabilization}
	sensor_data.RunTimeMaxStarboardCoverage = []uint8{rt_buff.RunTimeMaxStarboardCoverage}
	sensor_data.RunTimeMaxStarboardSwathWidth = []uint16{rt_buff.RunTimeMaxStarboardSwathWidth}

	switch rt_buff.RunTimeModelNumber {
	case 1002:
		_ = binary.Read(reader, binary.BigEndian, &RunTimeDurotongSpeed)
		sensor_data.RunTimeDurotongSpeed = []float32{float32(RunTimeDurotongSpeed) / 10.0}
		sensor_data.RunTimeTxAlongTilt = []float32{NULL_FLOAT32_ZERO}
	case 300:
		sensor_data.RunTimeDurotongSpeed = []float32{NULL_FLOAT32_ZERO}
		sensor_data.RunTimeTxAlongTilt = []float32{NULL_FLOAT32_ZERO}
	case 120:
		sensor_data.RunTimeDurotongSpeed = []float32{NULL_FLOAT32_ZERO}
		sensor_data.RunTimeTxAlongTilt = []float32{NULL_FLOAT32_ZERO}
	case 3020:
		sensor_data.RunTimeDurotongSpeed = []float32{NULL_FLOAT32_ZERO}
		_ = binary.Read(reader, binary.BigEndian, &RunTimeTxAlongTilt)
		sensor_data.RunTimeTxAlongTilt = []float32{float32(RunTimeTxAlongTilt) / SCALE2}
	default:
		_ = binary.Read(reader, binary.BigEndian, &Spare)
	}

	// appears that this piece is incomplete in the C-code and awaiting info from KM
	// regarding final datagram documentation.
	// This was captured back in 2009, and it is now 2024 with no updates
	// So merely replicating what they've constructed
	switch rt_buff.RunTimeModelNumber {
	default:
		_ = binary.Read(reader, binary.BigEndian, &RunTimeHiLoAbsorptionRatio)
		sensor_data.RunTimeHiLoAbsorptionRatio = []uint8{RunTimeHiLoAbsorptionRatio}
	}

	// fourth block (process unit)
	_ = binary.Read(reader, binary.BigEndian, &pu_buff)
	sensor_data.PuStatusPuCpuLoad = []uint8{pu_buff.PuStatusPuCpuLoad}
	sensor_data.PuStatusSensorStatus = []uint16{pu_buff.PuStatusSensorStatus}
	sensor_data.PuStatusAchievedPortCoverage = []uint8{pu_buff.PuStatusAchievedPortCoverage}
	sensor_data.PuStatusAchievedStarboardCoverage = []uint8{pu_buff.PuStatusAchievedStarboardCoverage}
	sensor_data.PuStatusYawStabilization = []float32{float32(pu_buff.PuStatusYawStabilization) / SCALE2}

	return sensor_data
}

type DeltaT struct {
	FileExtension        []string
	Version              []uint8
	PingByteSize         []uint16
	InterrogationTime    []time.Time
	SamplesPerBeam       []uint16
	SectorSize           []uint16
	StartAngle           []float32
	AngleIncrement       []float32
	AcousticRange        []uint16
	AcousticFrequency    []uint16
	SoundVelocity        []float32
	RangeResolution      []uint16
	ProfileTiltAngle     []float32
	RepetitionRate       []uint16
	PingNumber           []uint32
	IntensityFlag        []uint8
	PingLatency          []float32
	DataLatency          []float32
	SampleRateFlag       []uint8
	OptionsFlag          []uint8
	NumberPingsAveraged  []uint8
	CenterPingTimeOffset []float32
	UserDefinedByte      []uint8
	Altitude             []float32
	ExternalSensorFlags  []uint8
	PulseLength          []float64
	ForeAftBeamwidth     []float32
	AthwartBeamwidth     []float32
}

func DecodeDeltaTSpecific(reader *bytes.Reader) (sensor_data DeltaT) {
	var buffer struct {
		FileExtension        [4]byte
		Version              uint8
		PingByteSize         uint16
		TvSec                uint32
		TvNsec               uint32
		SamplesPerBeam       uint16
		SectorSize           uint16
		StartAngle           uint16
		AngleIncrement       uint16
		AcousticRange        uint16
		AcousticFrequency    uint16
		SoundVelocity        uint16
		RangeResolution      uint16
		ProfileTiltAngle     uint16
		RepetitionRate       uint16
		PingNumber           uint32
		IntensityFlag        uint8
		PingLatency          uint16
		DataLatency          uint16
		SampleRateFlag       uint8
		OptionsFlag          uint8
		NumberPingsAveraged  uint8
		CenterPingTimeOffset uint16
		UserDefinedByte      uint8
		Altitude             uint32
		ExternalSensorFlags  uint8
		PulseLength          uint32
		ForeAftBeamwidth     uint8
		AthwartBeamwidth     uint8
		Spare                [32]byte
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.FileExtension = []string{string(buffer.FileExtension[:])}
	sensor_data.Version = []uint8{buffer.Version}
	sensor_data.PingByteSize = []uint16{buffer.PingByteSize}
	sensor_data.InterrogationTime = []time.Time{time.Unix(int64(buffer.TvSec), int64(buffer.TvNsec)).UTC()}
	sensor_data.SamplesPerBeam = []uint16{buffer.SamplesPerBeam}
	sensor_data.SectorSize = []uint16{buffer.SectorSize}
	sensor_data.StartAngle = []float32{(float32(buffer.StartAngle) / SCALE2) - 180.0}
	sensor_data.AngleIncrement = []float32{float32(buffer.AngleIncrement) / SCALE2}
	sensor_data.AcousticRange = []uint16{buffer.AcousticRange}
	sensor_data.AcousticFrequency = []uint16{buffer.AcousticFrequency}
	sensor_data.SoundVelocity = []float32{float32(buffer.SoundVelocity) / 10.0}
	sensor_data.RangeResolution = []uint16{buffer.RangeResolution}
	sensor_data.ProfileTiltAngle = []float32{float32(buffer.ProfileTiltAngle) - 180.0}
	sensor_data.RepetitionRate = []uint16{buffer.RepetitionRate}
	sensor_data.PingNumber = []uint32{buffer.PingNumber}
	sensor_data.IntensityFlag = []uint8{buffer.IntensityFlag}
	sensor_data.PingLatency = []float32{float32(buffer.PingLatency) / 10_000.0}
	sensor_data.UserDefinedByte = []uint8{buffer.UserDefinedByte}
	sensor_data.Altitude = []float32{float32(buffer.Altitude) / SCALE2}
	sensor_data.ExternalSensorFlags = []uint8{buffer.ExternalSensorFlags}
	sensor_data.PulseLength = []float64{float64(buffer.PulseLength) / 1_000_000.0}
	sensor_data.ForeAftBeamwidth = []float32{float32(buffer.ForeAftBeamwidth) / 10.0}
	sensor_data.AthwartBeamwidth = []float32{float32(buffer.AthwartBeamwidth) / 10.0}

	return sensor_data
}

type R2Sonic struct {
	ModelNumber      []string
	SerialNumber     []string
	DgTime           []time.Time
	PingNumber       []uint32
	PingPeriod       []float64
	SoundSpeed       []float32
	Frequency        []float32
	TxPower          []float32
	TxPulseWidth     []float64
	TxBeamWidthVert  []float64
	TxBeamWidthHoriz []float64
	TxSteeringVert   []float64
	TxSteeringHoriz  []float64
	TxMiscInfo       []uint32
	RxBandwidth      []float32
	RxSampleRate     []float32
	RxRange          []float64
	RxGain           []float32
	RxSpreading      []float32
	RxAbsorption     []float32
	RxMountTilt      []float64
	RxMiscInfo       []uint32
	NumberBeams      []uint16
	A0MoreInfo       [][]float64
	A2MoreInfo       [][]float64
	G0DepthGateMin   []float64
	G0DepthGateMax   []float64
	G0DepthGateSlope []float64
}

func DecodeR2SonicSpecific(reader *bytes.Reader) (sensor_data R2Sonic) {
	var (
		buffer1 struct {
			ModelNumber      [12]byte
			SerialNumber     [12]byte
			TvSec            uint32
			TvNsec           uint32
			PingNumber       uint32
			PingPeriod       uint32
			SoundSpeed       uint32
			Frequency        uint32
			TxPower          uint32
			TxPulseWidth     uint32
			TxBeamWidthVert  uint32
			TxBeamWidthHoriz uint32
			TxSteeringVert   uint32
			TxSteeringHoriz  uint32
			TxMiscInfo       uint32
			RxBandwidth      uint32
			RxSampleRate     uint32
			RxRange          uint32
			RxGain           uint32
			RxSpreading      uint32
			RxAbsorption     uint32
			RxMountTilt      uint32
			RxMiscInfo       uint32
			Reserved         uint16
			NumberBeams      uint16
		}
		var_buf struct {
			A0MoreInfo [6]int32
			A2MoreInfo [6]int32
		}
		buffer2 struct {
			G0DepthGateMin   uint32
			G0DepthGateMax   uint32
			G0DepthGateSlope int32
			Spare            [32]byte
		}
	)

	// block one
	_ = binary.Read(reader, binary.BigEndian, &buffer1)
	sensor_data.ModelNumber = []string{string(buffer1.ModelNumber[:])}
	sensor_data.SerialNumber = []string{string(buffer1.SerialNumber[:])}
	sensor_data.DgTime = []time.Time{time.Unix(int64(buffer1.TvSec), int64(buffer1.TvNsec)).UTC()}
	sensor_data.PingNumber = []uint32{buffer1.PingNumber}
	sensor_data.PingPeriod = []float64{float64(buffer1.PingPeriod) / 1_000_000.0}
	sensor_data.SoundSpeed = []float32{float32(buffer1.SoundSpeed) / SCALE2}
	sensor_data.Frequency = []float32{float32(buffer1.Frequency) / SCALE3}
	sensor_data.TxPower = []float32{float32(buffer1.TxPower) / SCALE2}
	sensor_data.TxPulseWidth = []float64{float64(buffer1.TxPulseWidth) / 10_000_000.0}
	sensor_data.TxBeamWidthVert = []float64{float64(buffer1.TxBeamWidthVert) / 1_000_000.0}
	sensor_data.TxBeamWidthHoriz = []float64{float64(buffer1.TxBeamWidthHoriz) / 1_000_000.0}
	sensor_data.TxSteeringVert = []float64{float64(buffer1.TxSteeringVert) / 1_000_000.0}
	sensor_data.TxSteeringHoriz = []float64{float64(buffer1.TxSteeringHoriz) / 1_000_000.0}
	sensor_data.TxMiscInfo = []uint32{buffer1.TxMiscInfo}
	sensor_data.RxBandwidth = []float32{float32(buffer1.RxBandwidth) / 10_000.0}
	sensor_data.RxSampleRate = []float32{float32(buffer1.RxSampleRate) / SCALE3}
	sensor_data.RxRange = []float64{float64(buffer1.RxRange) / 100_000.0}
	sensor_data.RxGain = []float32{float32(buffer1.RxGain) / SCALE2}
	sensor_data.RxSpreading = []float32{float32(buffer1.RxSpreading) / SCALE3}
	sensor_data.RxAbsorption = []float32{float32(buffer1.RxAbsorption) / SCALE3}
	sensor_data.RxMountTilt = []float64{float64(buffer1.RxMountTilt) / 1_000_000.0}
	sensor_data.RxMiscInfo = []uint32{buffer1.RxMiscInfo}
	sensor_data.NumberBeams = []uint16{buffer1.NumberBeams}

	// block two (var length arrays)
	_ = binary.Read(reader, binary.BigEndian, &var_buf)
	A0MoreInfo := make([]float64, 0, 6)
	A2MoreInfo := make([]float64, 0, 6)

	for i := 0; i < 6; i++ {
		A0MoreInfo = append(A0MoreInfo, float64(var_buf.A0MoreInfo[i])/1_000_000.0)
		A2MoreInfo = append(A2MoreInfo, float64(var_buf.A2MoreInfo[i])/1_000_000.0)
	}

	sensor_data.A0MoreInfo = [][]float64{A0MoreInfo}
	sensor_data.A2MoreInfo = [][]float64{A2MoreInfo}

	// block three
	_ = binary.Read(reader, binary.BigEndian, &buffer2)
	sensor_data.G0DepthGateMin = []float64{float64(buffer2.G0DepthGateMin) / 1_000_000.0}
	sensor_data.G0DepthGateMax = []float64{float64(buffer2.G0DepthGateMax) / 1_000_000.0}
	sensor_data.G0DepthGateSlope = []float64{float64(buffer2.G0DepthGateSlope) / 1_000_000.0}

	return sensor_data
}

type ResonTSeries struct {
	ProtocolVersion                   []uint16
	DeviceID                          []uint32
	NumberDevices                     []uint32
	SystemEnumerator                  []uint16
	MajorSerialNumber                 []uint32
	MinorSerialNumber                 []uint32
	PingNumber                        []uint32
	MultiPingSequence                 []uint16
	Frequency                         []float32
	SampleRate                        []float32
	ReceiverBandwidth                 []float32
	TxPulseWidth                      []float64
	TxPulseTypeID                     []uint32
	TxPulseEnvlpID                    []uint32
	TxPulseEnvlpParam                 []float32
	TxPulseMode                       []uint16
	MaxPingRate                       []float64
	PingPeriod                        []float64
	Range                             []float32
	Power                             []float32
	Gain                              []float32
	ControlFlags                      []uint32
	ProjectorID                       []uint32
	ProjectorSteerAnglVert            []float32
	ProjectorSteerAnglHorz            []float32
	ProjectorBeamWidthVert            []float32
	ProjectorBeamWidthHorz            []float32
	ProjectorBeamFocalPt              []float32
	ProjectorBeamWeightingWindowType  []uint32
	ProjectorBeamWeightingWindowParam []uint32
	TransmitFlags                     []uint32
	HydrophoneID                      []uint32
	ReceivingBeamWeightingWindowType  []uint32
	ReceivingBeamWeightingWindowParam []uint32
	ReceiveFlags                      []uint32
	ReceiveBeamWidth                  []float32
	RangeFiltMin                      []float32
	RangeFiltMax                      []float32
	DepthFiltMin                      []float32
	DepthFiltMax                      []float32
	Absorption                        []float32
	SoundVelocity                     []float64
	SvSource                          []uint8
	Spreading                         []float32
	BeamSpacingMode                   []uint16
	SonarSourceMode                   []uint16
	CoverageMode                      []uint8
	CoverageAngle                     []float32
	HorizontalReceiverSteeringAngle   []float32
	UncertaintyType                   []uint32
	TransmitterSteeringAngle          []float32
	AppliedRoll                       []float32
	DetectionAlgorithm                []uint16
	DetectionFlags                    []uint32
	DeviceDescription                 []string
}

func DecodeResonTSeriesSonicSpecific(reader *bytes.Reader) (sensor_data ResonTSeries) {
	var buffer struct {
		ProtocolVersion                   uint16
		DeviceID                          uint32
		NumberDevices                     uint32
		SystemEnumerator                  uint16
		Reserved1                         [10]byte
		MajorSerialNumber                 uint32
		MinorSerialNumber                 uint32
		PingNumber                        uint32
		MultiPingSequence                 uint16
		Frequency                         uint32
		SampleRate                        uint32
		ReceiverBandwidth                 uint32
		TxPulseWidth                      uint32
		TxPulseTypeID                     uint32
		TxPulseEnvlpID                    uint32
		TxPulseEnvlpParam                 uint32
		TxPulseMode                       uint16
		TxPulseReserved                   uint16
		MaxPingRate                       uint32
		PingPeriod                        uint32
		Range                             uint32
		Power                             uint32
		Gain                              int32
		ControlFlags                      uint32
		ProjectorID                       uint32
		ProjectorSteerAnglVert            int32
		ProjectorSteerAnglHorz            int32
		ProjectorBeamWidthVert            uint16
		ProjectorBeamWidthHorz            uint16
		ProjectorBeamFocalPt              uint32
		ProjectorBeamWeightingWindowType  uint32
		ProjectorBeamWeightingWindowParam uint32
		TransmitFlags                     uint32
		HydrophoneID                      uint32
		ReceivingBeamWeightingWindowType  uint32
		ReceivingBeamWeightingWindowParam uint32
		ReceiveFlags                      uint32
		ReceiveBeamWidth                  uint16
		RangeFiltMin                      int32
		RangeFiltMax                      int32
		DepthFiltMin                      int32
		DepthFiltMax                      int32
		Absorption                        uint32
		SoundVelocity                     uint16
		SvSource                          uint8
		Spreading                         uint32
		BeamSpacingMode                   uint16
		SonarSourceMode                   uint16
		CoverageMode                      uint8
		CoverageAngle                     uint32
		HorizontalReceiverSteeringAngle   int32
		Reserved2                         [3]byte
		UncertaintyType                   uint32
		TransmitterSteeringAngle          int32
		AppliedRoll                       int32
		DetectionAlgorithm                uint16
		DetectionFlags                    uint32
		DeviceDescription                 [60]byte
		SoundVelocity2                    uint32
		Reserved7027                      [416]byte
		Reserved3                         [32]byte
	}
	_ = binary.Read(reader, binary.BigEndian, &buffer)

	sensor_data.ProtocolVersion = []uint16{buffer.ProtocolVersion}
	sensor_data.DeviceID = []uint32{buffer.DeviceID}
	sensor_data.NumberDevices = []uint32{buffer.NumberDevices}
	sensor_data.SystemEnumerator = []uint16{buffer.SystemEnumerator}
	sensor_data.MajorSerialNumber = []uint32{buffer.MajorSerialNumber}
	sensor_data.MinorSerialNumber = []uint32{buffer.MinorSerialNumber}
	sensor_data.PingNumber = []uint32{buffer.PingNumber}
	sensor_data.MultiPingSequence = []uint16{buffer.MultiPingSequence}
	sensor_data.Frequency = []float32{float32(buffer.Frequency) / SCALE3}
	sensor_data.SampleRate = []float32{float32(buffer.SampleRate) / 10_000.0}
	sensor_data.ReceiverBandwidth = []float32{float32(buffer.ReceiverBandwidth) / 10_000.0}
	sensor_data.TxPulseWidth = []float64{float64(buffer.TxPulseWidth) / 10_000_000.0}
	sensor_data.TxPulseTypeID = []uint32{buffer.TxPulseTypeID}
	sensor_data.TxPulseEnvlpID = []uint32{buffer.TxPulseEnvlpID}
	sensor_data.TxPulseEnvlpParam = []float32{float32(buffer.TxPulseEnvlpParam) / SCALE2}
	sensor_data.TxPulseMode = []uint16{buffer.TxPulseMode}
	sensor_data.MaxPingRate = []float64{float64(buffer.MaxPingRate) / 1_000_000.0}
	sensor_data.PingPeriod = []float64{float64(buffer.PingPeriod) / 1_000_000.0}
	sensor_data.Range = []float32{float32(buffer.Range) / SCALE2}
	sensor_data.Power = []float32{float32(buffer.Power) / SCALE2}
	sensor_data.Gain = []float32{float32(buffer.Gain) / SCALE2}
	sensor_data.ControlFlags = []uint32{buffer.ControlFlags}
	sensor_data.ProjectorID = []uint32{buffer.ProjectorID}
	sensor_data.ProjectorSteerAnglVert = []float32{float32(buffer.ProjectorSteerAnglVert) / SCALE3}
	sensor_data.ProjectorSteerAnglHorz = []float32{float32(buffer.ProjectorSteerAnglHorz) / SCALE3}
	sensor_data.ProjectorBeamWidthVert = []float32{float32(buffer.ProjectorBeamWidthVert) / SCALE2}
	sensor_data.ProjectorBeamWidthHorz = []float32{float32(buffer.ProjectorBeamWidthHorz) / SCALE2}
	sensor_data.ProjectorBeamFocalPt = []float32{float32(buffer.ProjectorBeamFocalPt) / SCALE2}
	sensor_data.ProjectorBeamWeightingWindowType = []uint32{buffer.ProjectorBeamWeightingWindowType}
	sensor_data.ProjectorBeamWeightingWindowParam = []uint32{buffer.ProjectorBeamWeightingWindowParam}
	sensor_data.TransmitFlags = []uint32{buffer.TransmitFlags}
	sensor_data.HydrophoneID = []uint32{buffer.HydrophoneID}
	sensor_data.ReceivingBeamWeightingWindowType = []uint32{buffer.ReceivingBeamWeightingWindowType}
	sensor_data.ReceivingBeamWeightingWindowParam = []uint32{buffer.ReceivingBeamWeightingWindowParam}
	sensor_data.ReceiveFlags = []uint32{buffer.ReceiveFlags}
	sensor_data.ReceiveBeamWidth = []float32{float32(buffer.ReceiveBeamWidth) / SCALE2}
	sensor_data.RangeFiltMin = []float32{float32(buffer.RangeFiltMin) / 10.0}
	sensor_data.RangeFiltMax = []float32{float32(buffer.RangeFiltMax) / 10.0}
	sensor_data.DepthFiltMin = []float32{float32(buffer.DepthFiltMin) / 10.0}
	sensor_data.DepthFiltMax = []float32{float32(buffer.DepthFiltMax) / 10.0}
	sensor_data.Absorption = []float32{float32(buffer.Absorption) / SCALE3}
	sensor_data.SoundVelocity = []float64{float64(buffer.SoundVelocity) / 10.0}
	sensor_data.SvSource = []uint8{buffer.SvSource}
	sensor_data.Spreading = []float32{float32(buffer.Spreading) / SCALE3}
	sensor_data.BeamSpacingMode = []uint16{buffer.BeamSpacingMode}
	sensor_data.SonarSourceMode = []uint16{buffer.SonarSourceMode}
	sensor_data.CoverageMode = []uint8{buffer.CoverageMode}
	sensor_data.CoverageAngle = []float32{float32(buffer.CoverageAngle) / SCALE2}
	sensor_data.HorizontalReceiverSteeringAngle = []float32{float32(buffer.HorizontalReceiverSteeringAngle) / SCALE2}
	sensor_data.UncertaintyType = []uint32{buffer.UncertaintyType}
	sensor_data.TransmitterSteeringAngle = []float32{float32(buffer.TransmitterSteeringAngle) / 100_000.0}
	sensor_data.AppliedRoll = []float32{float32(buffer.AppliedRoll) / 100_000.0}
	sensor_data.DetectionAlgorithm = []uint16{buffer.DetectionAlgorithm}
	sensor_data.DetectionFlags = []uint32{buffer.DetectionFlags}
	sensor_data.DeviceDescription = []string{string(buffer.DeviceDescription[:])}

	// higher precision sound velocity
	if buffer.SoundVelocity2 > 0 {
		sensor_data.SoundVelocity = []float64{float64(buffer.SoundVelocity2) / 1_000_000.0}
	}

	return sensor_data
}

type Kmall struct {
	KmallVersion                       []uint8
	DgmType                            []uint8
	DgmVersion                         []uint8
	SystemID                           []uint8
	EchoSounderID                      []uint16
	NumBytesCmnPart                    []uint16
	PingCounter                        []uint16
	RxFansPerRing                      []uint8
	RxFansIndex                        []uint8
	SwathsPerRing                      []uint8
	SwathAlongPosition                 []uint8
	TxTransducerIndex                  []uint8
	RxTransducerIndex                  []uint8
	NumRxTransducers                   []uint8
	AlgorithmType                      []uint8
	NumBytesInfoData                   []uint16
	PingRateHz                         []float64
	BeamSpacing                        []uint8
	DepthMode                          []uint8
	SubDepthMode                       []uint8
	DistanceBetweenSwath               []uint8
	DetectionMode                      []uint8
	PulseForm                          []uint8
	FrequencyModeHz                    []int32
	FrequencyRangeLowLimHz             []float64
	FrequencyRangeHighLimHz            []float64
	MaxTotalTxPulseLengthSector        []float64
	MaxEffectiveTxPulseLengthSector    []float64
	MaxEffectiveTxBandWidthHz          []float64
	AbsCoeffDbPerKm                    []float64
	PortSectorEdgeDeg                  []float32
	StarboardSectorEdgeDeg             []float32
	PortMeanCoverageDeg                []float32
	StarboardMeanCoverageDeg           []float32
	PortMeanCoverageMetres             []int16
	StarboardMeanCoverageMetres        []int16
	ModeAndStabilisation               []uint8
	RunTimeFilter1                     []uint8
	RunTimeFilter2                     []uint8
	PipeTrackingStatus                 []uint32
	TransmitArraySizeUsedDeg           []float32
	ReceiveArraySizeUsedDeg            []float32
	TransmitPowerDb                    []float32
	SlRampUpTimeRemaining              []uint16
	YawAngleDeg                        []float64
	NumTxSectors                       []uint16
	NumBytesPerTxSector                []uint16
	HeadingVesselDeg                   []float64
	SoundSpeedAtTxDepthMetresPerSecond []float64
	TxTransducerDepthMetres            []float64
	ZwaterLevelReRefPointMetres        []float64
	XKmallToAllMetres                  []float64
	YKmallToAllMetres                  []float64
	LatLonInfo                         []uint8
	PositionSensorStatus               []uint8
	AttitudeSensorStatus               []uint8
	LatitudeDeg                        []float64
	LongitudeDeg                       []float64
	EllipsoidHeightReRefPointMetres    []float64
	TxSectorNumber                     [][]uint8
	TxArrayNumber                      [][]uint8
	TxSubArray                         [][]uint8
	SectorTransmitDelaySec             [][]float64
	TiltAngleReTxDeg                   [][]float64
	TxNominalSourceLevelDb             [][]float64
	TxFocusRangeMetres                 [][]float64
	CentreFrequencyHz                  [][]float64
	SignalBandWidthHz                  [][]float64
	TotalSignalLengthSec               [][]float64
	PulseShading                       [][]uint8
	SignalWaveForm                     [][]uint8
	NumBytesRxInfo                     []uint16
	NumSoundingsMaxMain                []uint16
	NumSoundingsValidMain              []uint16
	NumBytesPerSounding                []uint16
	WcSampleRate                       []float64
	SeabedImageSampleRate              []float64
	BackscatterNormalDb                []float64
	BackscatterObliqueDb               []float64
	ExtraDetectionAlarmFlag            []uint16
	NumExtraDetections                 []uint16
	NumExtraDetectionClasses           []uint16
	NumBytesPerClass                   []uint16
	NumExtraDetectionInClass           [][]uint16
	AlarmFlag                          [][]uint8
}

func DecodeKmallSpecific(reader *bytes.Reader) (sensor_data Kmall) {
	var (
		buffer struct {
			KmallVersion  uint8
			DgmType       uint8
			DgmVersion    uint8
			SystemID      uint8
			EchoSounderID uint16
			Spare         [8]byte
		}
		cmn_buf struct {
			NumBytesCmnPart    uint16
			PingCounter        uint16
			RxFansPerRing      uint8
			RxFansIndex        uint8
			SwathsPerRing      uint8
			SwathAlongPosition uint8
			TxTransducerIndex  uint8
			RxTransducerIndex  uint8
			NumRxTransducers   uint8
			AlgorithmType      uint8
			Spare              [16]byte
		}
		ping_buf struct {
			NumBytesInfoData                uint16
			PingRateHz                      uint32
			BeamSpacing                     uint8
			DepthMode                       uint8
			SubDepthMode                    uint8
			DistanceBetweenSwath            uint8
			DetectionMode                   uint8
			PulseForm                       uint8
			FrequencyModeHz                 int32
			FrequencyRangeLowLimHz          int32
			FrequencyRangeHighLimHz         int32
			MaxTotalTxPulseLengthSector     int32
			MaxEffectiveTxPulseLengthSector int32
			MaxEffectiveTxBandWidthHz       int32
			AbsCoeffDbPerKm                 int32
			PortSectorEdgeDeg               int16
			StarboardSectorEdgeDeg          int16
			PortMeanCoverageDeg             int16
			StarboardMeanCoverageDeg        int16
			// PortMeanCoverageDeg2 int16 // the C-code reads PortMeanCoverageDeg twice, spec doesn't indicate this. potential error????
			// StarboardMeanCoverageDeg2 int16 // the C-code reads StarboardMeanCoverageDeg twice, spec doesn't indicate this. potential error????
			PortMeanCoverageMetres             int16
			StarboardMeanCoverageMetres        int16
			ModeAndStabilisation               uint8
			RunTimeFilter1                     uint8
			RunTimeFilter2                     uint8
			PipeTrackingStatus                 uint32
			TransmitArraySizeUsedDeg           uint16
			ReceiveArraySizeUsedDeg            uint16
			TransmitPowerDb                    int16
			SlRampUpTimeRemaining              uint16
			YawAngleDeg                        uint32
			NumTxSectors                       uint16
			NumBytesPerTxSector                uint16
			HeadingVesselDeg                   int32
			SoundSpeedAtTxDepthMetresPerSecond int32
			TxTransducerDepthMetres            int32
			ZwaterLevelReRefPointMetres        int32
			XKmallToAllMetres                  int32
			YKmallToAllMetres                  int32
			LatLonInfo                         uint8
			PositionSensorStatus               uint8
			AttitudeSensorStatus               uint8
			LatitudeDeg                        int32
			LongitudeDeg                       int32
			EllipsoidHeightReRefPointMetres    int32
			Spare                              [32]byte
		}
		sec_buf struct {
			TxSectorNumber         uint8
			TxArrayNumber          uint8
			TxSubArray             uint8
			SectorTransmitDelaySec int32
			TiltAngleReTxDeg       int32
			TxNominalSourceLevelDb int32
			TxFocusRangeMetres     int32
			CentreFrequencyHz      int32
			SignalBandWidthHz      int32
			TotalSignalLengthSec   int32
			PulseShading           uint8
			SignalWaveForm         uint8
			Spare                  [20]byte
		}
		rx_buf struct {
			NumBytesRxInfo           uint16
			NumSoundingsMaxMain      uint16
			NumSoundingsValidMain    uint16
			NumBytesPerSounding      uint16
			WcSampleRate1            int32
			WcSampleRate2            uint32
			SeabedImageSampleRate1   int32
			SeabedImageSampleRate2   uint32
			BackscatterNormalDb      int32
			BackscatterObliqueDb     int32
			ExtraDetectionAlarmFlag  uint16
			NumExtraDetections       uint16
			NumExtraDetectionClasses uint16
			NumBytesPerClass         uint16
			Spare                    [32]byte
		}
		cls_buf struct {
			NumExtraDetectionInClass uint16
			AlarmFlag                uint8
			Spare                    [32]byte
		}
		final_spare [32]byte
	)

	// block one
	_ = binary.Read(reader, binary.BigEndian, &buffer)
	sensor_data.KmallVersion = []uint8{buffer.KmallVersion}
	sensor_data.DgmType = []uint8{buffer.DgmType}
	sensor_data.DgmVersion = []uint8{buffer.DgmVersion}
	sensor_data.EchoSounderID = []uint16{buffer.EchoSounderID}

	// block two (Cmn part)
	_ = binary.Read(reader, binary.BigEndian, &cmn_buf)
	sensor_data.NumBytesCmnPart = []uint16{cmn_buf.NumBytesCmnPart}
	sensor_data.PingCounter = []uint16{cmn_buf.PingCounter}
	sensor_data.RxFansPerRing = []uint8{cmn_buf.RxFansPerRing}
	sensor_data.RxFansIndex = []uint8{cmn_buf.RxFansIndex}
	sensor_data.SwathsPerRing = []uint8{cmn_buf.SwathsPerRing}
	sensor_data.SwathAlongPosition = []uint8{cmn_buf.SwathAlongPosition}
	sensor_data.TxTransducerIndex = []uint8{cmn_buf.TxTransducerIndex}
	sensor_data.RxTransducerIndex = []uint8{cmn_buf.RxTransducerIndex}
	sensor_data.NumRxTransducers = []uint8{cmn_buf.NumRxTransducers}
	sensor_data.AlgorithmType = []uint8{cmn_buf.AlgorithmType}

	// block three (ping data)
	_ = binary.Read(reader, binary.BigEndian, &cmn_buf)
	sensor_data.NumBytesInfoData = []uint16{ping_buf.NumBytesInfoData}
	sensor_data.PingRateHz = []float64{float64(ping_buf.PingRateHz) / 100_000.0}
	sensor_data.BeamSpacing = []uint8{ping_buf.BeamSpacing}
	sensor_data.DepthMode = []uint8{ping_buf.DepthMode}
	sensor_data.SubDepthMode = []uint8{ping_buf.SubDepthMode}
	sensor_data.DistanceBetweenSwath = []uint8{ping_buf.DistanceBetweenSwath}
	sensor_data.DetectionMode = []uint8{ping_buf.DetectionMode}
	sensor_data.PulseForm = []uint8{ping_buf.PulseForm}
	sensor_data.FrequencyModeHz = []int32{ping_buf.FrequencyModeHz}
	sensor_data.FrequencyRangeLowLimHz = []float64{float64(ping_buf.FrequencyRangeLowLimHz) / 1_000.0}
	sensor_data.FrequencyRangeHighLimHz = []float64{float64(ping_buf.FrequencyRangeHighLimHz) / 1_000.0}
	sensor_data.MaxTotalTxPulseLengthSector = []float64{float64(ping_buf.MaxTotalTxPulseLengthSector) / 1_000_000.0}
	sensor_data.MaxEffectiveTxPulseLengthSector = []float64{float64(ping_buf.MaxEffectiveTxPulseLengthSector) / 1_000_0.0}
	sensor_data.MaxEffectiveTxBandWidthHz = []float64{float64(ping_buf.MaxEffectiveTxBandWidthHz) / 1_000.0}
	sensor_data.AbsCoeffDbPerKm = []float64{float64(ping_buf.AbsCoeffDbPerKm) / 1_000.0}
	sensor_data.PortSectorEdgeDeg = []float32{float32(ping_buf.PortSectorEdgeDeg) / SCALE2}
	sensor_data.StarboardSectorEdgeDeg = []float32{float32(ping_buf.StarboardSectorEdgeDeg) / SCALE2}
	sensor_data.PortMeanCoverageDeg = []float32{float32(ping_buf.PortMeanCoverageDeg) / SCALE2}
	sensor_data.StarboardMeanCoverageDeg = []float32{float32(ping_buf.StarboardMeanCoverageDeg) / SCALE2}
	sensor_data.PortMeanCoverageMetres = []int16{ping_buf.PortMeanCoverageMetres}
	sensor_data.StarboardMeanCoverageMetres = []int16{ping_buf.StarboardMeanCoverageMetres}
	sensor_data.ModeAndStabilisation = []uint8{ping_buf.ModeAndStabilisation}
	sensor_data.RunTimeFilter1 = []uint8{ping_buf.RunTimeFilter1}
	sensor_data.RunTimeFilter2 = []uint8{ping_buf.RunTimeFilter2}
	sensor_data.PipeTrackingStatus = []uint32{ping_buf.PipeTrackingStatus}
	sensor_data.TransmitArraySizeUsedDeg = []float32{float32(ping_buf.TransmitArraySizeUsedDeg) / SCALE3}
	sensor_data.ReceiveArraySizeUsedDeg = []float32{float32(ping_buf.ReceiveArraySizeUsedDeg) / SCALE3}
	sensor_data.TransmitPowerDb = []float32{float32(ping_buf.TransmitPowerDb) / SCALE2}
	sensor_data.SlRampUpTimeRemaining = []uint16{ping_buf.SlRampUpTimeRemaining}
	sensor_data.YawAngleDeg = []float64{float64(ping_buf.YawAngleDeg) / 1_000_000.0}
	sensor_data.NumTxSectors = []uint16{ping_buf.NumTxSectors}
	sensor_data.NumBytesPerTxSector = []uint16{ping_buf.NumBytesPerTxSector}
	sensor_data.HeadingVesselDeg = []float64{float64(ping_buf.HeadingVesselDeg) / 1_000_000.0}
	sensor_data.SoundSpeedAtTxDepthMetresPerSecond = []float64{float64(ping_buf.SoundSpeedAtTxDepthMetresPerSecond) / 1_000_000.0}
	sensor_data.TxTransducerDepthMetres = []float64{float64(ping_buf.TxTransducerDepthMetres) / 1_000_000.0}
	sensor_data.ZwaterLevelReRefPointMetres = []float64{float64(ping_buf.ZwaterLevelReRefPointMetres) / 1_000_000.0}
	sensor_data.XKmallToAllMetres = []float64{float64(ping_buf.XKmallToAllMetres) / 1_000_000.0}
	sensor_data.YKmallToAllMetres = []float64{float64(ping_buf.YKmallToAllMetres) / 1_000_000.0}
	sensor_data.LatLonInfo = []uint8{ping_buf.LatLonInfo}
	sensor_data.PositionSensorStatus = []uint8{ping_buf.PositionSensorStatus}
	sensor_data.AttitudeSensorStatus = []uint8{ping_buf.AttitudeSensorStatus}
	sensor_data.LatitudeDeg = []float64{float64(ping_buf.LatitudeDeg) / 10_000_000.0}
	sensor_data.LongitudeDeg = []float64{float64(ping_buf.LongitudeDeg) / 10_000_000.0}
	sensor_data.EllipsoidHeightReRefPointMetres = []float64{float64(ping_buf.EllipsoidHeightReRefPointMetres) / 1_000.0}

	// block four (sector specific info)
	nsectors := int(ping_buf.NumTxSectors)
	TxSectorNumber := make([]uint8, 0, nsectors)
	TxArrayNumber := make([]uint8, 0, nsectors)
	TxSubArray := make([]uint8, 0, nsectors)
	SectorTransmitDelaySec := make([]float64, 0, nsectors)
	TiltAngleReTxDeg := make([]float64, 0, nsectors)
	TxNominalSourceLevelDb := make([]float64, 0, nsectors)
	TxFocusRangeMetres := make([]float64, 0, nsectors)
	CentreFrequencyHz := make([]float64, 0, nsectors)
	SignalBandWidthHz := make([]float64, 0, nsectors)
	TotalSignalLengthSec := make([]float64, 0, nsectors)
	PulseShading := make([]uint8, 0, nsectors)
	SignalWaveForm := make([]uint8, 0, nsectors)

	for i := uint16(0); i < ping_buf.NumTxSectors; i++ {
		_ = binary.Read(reader, binary.BigEndian, &sec_buf)
		TxSectorNumber = append(TxSectorNumber, sec_buf.TxSectorNumber)
		TxArrayNumber = append(TxArrayNumber, sec_buf.TxArrayNumber)
		TxSubArray = append(TxSubArray, sec_buf.TxSubArray)
		SectorTransmitDelaySec = append(SectorTransmitDelaySec, float64(sec_buf.SectorTransmitDelaySec)/1_000_000.0)
		TiltAngleReTxDeg = append(TiltAngleReTxDeg, float64(sec_buf.TiltAngleReTxDeg)/1_000_000.0)
		TxNominalSourceLevelDb = append(TxNominalSourceLevelDb, float64(sec_buf.TxNominalSourceLevelDb)/1_000_000.0)
		TxFocusRangeMetres = append(TxFocusRangeMetres, float64(sec_buf.TxFocusRangeMetres)/1_000.0)
		CentreFrequencyHz = append(CentreFrequencyHz, float64(sec_buf.CentreFrequencyHz)/1_000.0)
		SignalBandWidthHz = append(SignalBandWidthHz, float64(sec_buf.SignalBandWidthHz)/1_000.0)
		TotalSignalLengthSec = append(TotalSignalLengthSec, float64(sec_buf.TotalSignalLengthSec)/1_000_000.0)
		PulseShading = append(PulseShading, sec_buf.PulseShading)
		SignalWaveForm = append(SignalWaveForm, sec_buf.SignalWaveForm)
	}

	sensor_data.TxSectorNumber = [][]uint8{TxSectorNumber}
	sensor_data.TxArrayNumber = [][]uint8{TxArrayNumber}
	sensor_data.TxSubArray = [][]uint8{TxSubArray}
	sensor_data.SectorTransmitDelaySec = [][]float64{SectorTransmitDelaySec}
	sensor_data.TiltAngleReTxDeg = [][]float64{TiltAngleReTxDeg}
	sensor_data.TxNominalSourceLevelDb = [][]float64{TxNominalSourceLevelDb}
	sensor_data.TxFocusRangeMetres = [][]float64{TxFocusRangeMetres}
	sensor_data.CentreFrequencyHz = [][]float64{CentreFrequencyHz}
	sensor_data.SignalBandWidthHz = [][]float64{SignalBandWidthHz}
	sensor_data.TotalSignalLengthSec = [][]float64{TotalSignalLengthSec}
	sensor_data.PulseShading = [][]uint8{PulseShading}
	sensor_data.SignalWaveForm = [][]uint8{SignalWaveForm}

	// block five (rx info)
	_ = binary.Read(reader, binary.BigEndian, &rx_buf)
	sensor_data.NumBytesRxInfo = []uint16{rx_buf.NumBytesRxInfo}
	sensor_data.NumSoundingsMaxMain = []uint16{rx_buf.NumSoundingsMaxMain}
	sensor_data.NumSoundingsValidMain = []uint16{rx_buf.NumSoundingsValidMain}
	sensor_data.NumBytesPerSounding = []uint16{rx_buf.NumBytesPerSounding}
	sensor_data.WcSampleRate = []float64{float64(rx_buf.WcSampleRate1) + (float64(rx_buf.WcSampleRate2) / 1_000_000_000.0)}
	sensor_data.SeabedImageSampleRate = []float64{float64(rx_buf.SeabedImageSampleRate1) + (float64(rx_buf.SeabedImageSampleRate2) / 1_000_000_000.0)}
	sensor_data.BackscatterNormalDb = []float64{float64(rx_buf.BackscatterNormalDb) / 1_000_000.0}
	sensor_data.BackscatterObliqueDb = []float64{float64(rx_buf.BackscatterObliqueDb) / 1_000_000.0}
	sensor_data.ExtraDetectionAlarmFlag = []uint16{rx_buf.ExtraDetectionAlarmFlag}
	sensor_data.NumExtraDetections = []uint16{rx_buf.NumExtraDetections}
	sensor_data.NumExtraDetectionClasses = []uint16{rx_buf.NumExtraDetectionClasses}
	sensor_data.NumBytesPerClass = []uint16{rx_buf.NumBytesPerClass}

	// block six (extra detection classes
	nclasses := int(rx_buf.NumExtraDetectionClasses)
	NumExtraDetectionInClass := make([]uint16, 0, nclasses)
	AlarmFlag := make([]uint8, 0, nclasses)

	for i := 0; i < nclasses; i++ {
		_ = binary.Read(reader, binary.BigEndian, &cls_buf)
		NumExtraDetectionInClass = append(NumExtraDetectionInClass, cls_buf.NumExtraDetectionInClass)
		AlarmFlag = append(AlarmFlag, cls_buf.AlarmFlag)
	}

	sensor_data.NumExtraDetectionInClass = [][]uint16{NumExtraDetectionInClass}
	sensor_data.AlarmFlag = [][]uint8{AlarmFlag}

	_ = binary.Read(reader, binary.BigEndian, &final_spare)

	return sensor_data
}
