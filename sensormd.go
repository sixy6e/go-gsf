package gsf

import (
	"bytes"
	"encoding/binary"
	"errors"
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

var ErrSensorMetadata = errors.New("Error reading Sensor Metadata")

type Seabeam struct {
	EclipseTime []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSeabeamSpecific(reader *bytes.Reader) (sensor_data Seabeam, err error) {
	var buffer struct {
		EclipseTime uint16
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Seabeam sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.EclipseTime = []uint16{buffer.EclipseTime}

	return sensor_data, err
}

type Em12 struct {
	PingNumber    []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Resolution    []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PingQuality   []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Mode          []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm12Specific(reader *bytes.Reader) (sensor_data Em12, err error) {
	var buffer struct {
		PingNumber    uint16
		Resolution    uint8
		PingQuality   uint8
		SoundVelocity uint16
		Mode          uint8
		Spare         [4]int32
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("EM12 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.Resolution = []uint8{buffer.Resolution}
	sensor_data.PingQuality = []uint8{buffer.PingQuality}
	sensor_data.SoundVelocity = []float32{float32(buffer.SoundVelocity) / 10.0}
	sensor_data.Mode = []uint8{buffer.Mode}

	return sensor_data, err
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

func DecodeEm100Specific(reader *bytes.Reader) (sensor_data Em100, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("EM100 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.ShipPitch = []float32{float32(buffer.ShipPitch) / SCALE_2_F32}
	sensor_data.TransducerPitch = []float32{float32(buffer.TransducerPitch) / SCALE_2_F32}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.Power = []uint8{buffer.Power}
	sensor_data.Attenuation = []uint8{buffer.Attenuation}
	sensor_data.Tvg = []uint8{buffer.Tvg}
	sensor_data.PulseLength = []uint8{buffer.PulseLength}
	sensor_data.Counter = []uint16{buffer.Counter}

	return sensor_data, err
}

type Em950 struct {
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Quality              []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ShipPitch            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransducerPitch      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm950Specific(reader *bytes.Reader) (sensor_data Em950, err error) {
	var buffer struct {
		PingNumber           uint16
		Mode                 uint8
		Quality              uint8
		ShipPitch            int16
		TransducerPitch      int16
		SurfaceSoundVelocity uint16
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("EM950 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.Quality = []uint8{buffer.Quality}
	sensor_data.ShipPitch = []float32{float32(buffer.ShipPitch) / SCALE_2_F32}
	sensor_data.TransducerPitch = []float32{float32(buffer.TransducerPitch) / SCALE_2_F32}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / SCALE_1_F32}

	return sensor_data, err
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

func DecodeEm121ASpecific(reader *bytes.Reader) (sensor_data Em121A, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("EM121A sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.ValidBeams = []uint8{buffer.ValidBeams}
	sensor_data.PulseLength = []uint8{buffer.PulseLength}
	sensor_data.BeamWidth = []uint8{buffer.BeamWidth}
	sensor_data.TransmitPower = []uint8{buffer.TransmitPower}
	sensor_data.TransmitStatus = []uint8{buffer.TransmitStatus}
	sensor_data.ReceiveStatus = []uint8{buffer.ReceiveStatus}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.ReceiveStatus) / SCALE_1_F32}

	return sensor_data, err
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

func DecodeEm121Specific(reader *bytes.Reader) (sensor_data Em121, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("EM121 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.ValidBeams = []uint8{buffer.ValidBeams}
	sensor_data.PulseLength = []uint8{buffer.PulseLength}
	sensor_data.BeamWidth = []uint8{buffer.BeamWidth}
	sensor_data.TransmitPower = []uint8{buffer.TransmitPower}
	sensor_data.TransmitStatus = []uint8{buffer.TransmitStatus}
	sensor_data.ReceiveStatus = []uint8{buffer.ReceiveStatus}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.ReceiveStatus) / SCALE_1_F32}

	return sensor_data, err
}

type Sass struct {
	LeftMostBeam       []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RightMostBeam      []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TotalNumverOfBeams []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NavigationMode     []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingNumber         []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	MissionNumber      []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSassSpecfic(reader *bytes.Reader) (sensor_data Sass, err error) {
	var buffer struct {
		LeftMostBeam       uint16
		RightMostBeam      uint16
		TotalNumverOfBeams uint16
		NavigationMode     uint16
		PingNumber         uint16
		MissionNumber      uint16
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("SASS sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.LeftMostBeam = []uint16{buffer.LeftMostBeam}
	sensor_data.RightMostBeam = []uint16{buffer.RightMostBeam}
	sensor_data.TotalNumverOfBeams = []uint16{buffer.TotalNumverOfBeams}
	sensor_data.NavigationMode = []uint16{buffer.NavigationMode}
	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.MissionNumber = []uint16{buffer.MissionNumber}

	return sensor_data, err
}

type SeaMap struct {
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

func DecodeSeaMapSpecific(reader *bytes.Reader, gsfd GsfDetails) (sensor_data SeaMap, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer1)
	if err != nil {
		errn := errors.New("SeaMap sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	major, minor := gsfd.MajorMinor()

	if major > 2 || (major == 2 && minor > 7) {
		err = binary.Read(reader, binary.BigEndian, &pressure_depth)
		if err != nil {
			errn := errors.New("SeaMap sensor")
			err = errors.Join(err, ErrSensorMetadata, errn)
			return sensor_data, err
		}
	} else {
		pressure_depth = 0 // treating as null
	}

	err = binary.Read(reader, binary.BigEndian, &buffer2)
	if err != nil {
		errn := errors.New("SeaMap sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.PortTransmit1 = []float32{float32(buffer1.PortTransmit1) / SCALE_1_F32}
	sensor_data.PortTransmit2 = []float32{float32(buffer1.PortTransmit2) / SCALE_1_F32}
	sensor_data.StarboardTransmit1 = []float32{float32(buffer1.StarboardTransmit1) / SCALE_1_F32}
	sensor_data.StarboardTransmit2 = []float32{float32(buffer1.StarboardTransmit2) / SCALE_1_F32}
	sensor_data.PortGain = []float32{float32(buffer1.PortGain) / SCALE_1_F32}
	sensor_data.StarboardGain = []float32{float32(buffer1.StarboardGain) / SCALE_1_F32}
	sensor_data.PortPulseLength = []float32{float32(buffer1.PortPulseLength) / SCALE_1_F32}
	sensor_data.StarboardPulseLength = []float32{float32(buffer1.StarboardPulseLength) / SCALE_1_F32}
	sensor_data.PressureDepth = []float32{float32(pressure_depth) / SCALE_1_F32}
	sensor_data.Altitude = []float32{float32(buffer2.Altitude) / SCALE_1_F32}
	sensor_data.Temperature = []float32{float32(buffer2.Temperature) / SCALE_1_F32}

	return sensor_data, err
}

type SeaBat struct {
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Range                []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TransmitPower        []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ReceiveGain          []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSeaBatSpecific(reader *bytes.Reader) (sensor_data SeaBat, err error) {
	var buffer struct {
		PingNumber           uint16
		SurfaceSoundVelocity uint16
		Mode                 uint8
		Range                uint8
		TransmitPower        uint8
		ReceiveGain          uint8
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("SeaBat sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / SCALE_1_F32}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.Range = []uint8{buffer.Range}
	sensor_data.TransmitPower = []uint8{buffer.TransmitPower}
	sensor_data.ReceiveGain = []uint8{buffer.ReceiveGain}

	return sensor_data, err
}

type Em1000 struct {
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Quality              []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ShipPitch            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransducerPitch      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeEm1000Specific(reader *bytes.Reader) (sensor_data Em1000, err error) {
	var buffer struct {
		PingNumber           uint16
		Mode                 uint8
		Quality              uint8
		ShipPitch            int16
		TransducerPitch      int16
		SurfaceSoundVelocity uint16
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("EM1000 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.Quality = []uint8{buffer.Quality}
	sensor_data.ShipPitch = []float32{float32(buffer.ShipPitch) / SCALE_2_F32}
	sensor_data.TransducerPitch = []float32{float32(buffer.TransducerPitch) / SCALE_2_F32}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / SCALE_1_F32}

	return sensor_data, err
}

type TypeIIISeabeam struct {
	LeftMostBeam       []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RightMostBeam      []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TotalNumverOfBeams []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NavigationMode     []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingNumber         []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	MissionNumber      []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeTypeIIISeabeamSpecific(reader *bytes.Reader) (sensor_data TypeIIISeabeam, err error) {
	var buffer struct {
		LeftMostBeam       uint16
		RightMostBeam      uint16
		TotalNumverOfBeams uint16
		NavigationMode     uint16
		PingNumber         uint16
		MissionNumber      uint16
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("TypeIIISeaBeam sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.LeftMostBeam = []uint16{buffer.LeftMostBeam}
	sensor_data.RightMostBeam = []uint16{buffer.RightMostBeam}
	sensor_data.TotalNumverOfBeams = []uint16{buffer.TotalNumverOfBeams}
	sensor_data.NavigationMode = []uint16{buffer.NavigationMode}
	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.MissionNumber = []uint16{buffer.MissionNumber}

	return sensor_data, err
}

type SbAmp struct {
	Hour         []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Minute       []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Second       []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Hundredths   []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	BlockNumber  []uint32 `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	AvgGateDepth []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSbAmpSpecific(reader *bytes.Reader) (sensor_data SbAmp, err error) {
	var buffer struct {
		Hour         uint8
		Minute       uint8
		Second       uint8
		Hundredths   uint8
		BlockNumber  uint32
		AvgGateDepth uint16
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("SBAmp sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.Hour = []uint8{buffer.Hour}
	sensor_data.Minute = []uint8{buffer.Minute}
	sensor_data.Second = []uint8{buffer.Second}
	sensor_data.Hundredths = []uint8{buffer.Hundredths}
	sensor_data.BlockNumber = []uint32{buffer.BlockNumber}
	sensor_data.AvgGateDepth = []uint16{buffer.AvgGateDepth}

	return sensor_data, err
}

type SeaBatII struct {
	PingNumber           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Mode                 []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SonarRange           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TransmitPower        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ReceiveGain          []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ForeAftBandwidth     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	AthwartBandwidth     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSeaBatIISpecific(reader *bytes.Reader) (sensor_data SeaBatII, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("SeaBatII sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / SCALE_1_F32}
	sensor_data.Mode = []uint16{buffer.Mode}
	sensor_data.SonarRange = []uint16{buffer.SonarRange}
	sensor_data.TransmitPower = []uint16{buffer.TransmitPower}
	sensor_data.ReceiveGain = []uint16{buffer.ReceiveGain}
	sensor_data.ForeAftBandwidth = []float32{float32(buffer.ForeAftBandwidth) / SCALE_1_F32}
	sensor_data.AthwartBandwidth = []float32{float32(buffer.AthwartBandwidth) / SCALE_1_F32}

	return sensor_data, err
}

type SeaBat8101 struct {
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

func DecodeSeaBat8101Specific(reader *bytes.Reader) (sensor_data SeaBat8101, err error) {
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
		Spare                [4]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("SeaBat8101 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / SCALE_1_F32}
	sensor_data.Mode = []uint16{buffer.Mode}
	sensor_data.Range = []uint16{buffer.Range}
	sensor_data.TransmitPower = []uint16{buffer.TransmitPower}
	sensor_data.RecieveGain = []uint16{buffer.RecieveGain}
	sensor_data.PulseWidth = []uint16{buffer.PulseWidth}
	sensor_data.TvgSpreading = []uint8{buffer.TvgSpreading}
	sensor_data.TvgAbsorption = []uint8{buffer.TvgAbsorption}
	sensor_data.ForeAftBandwidth = []float32{float32(buffer.ForeAftBandwidth) / SCALE_1_F32}
	sensor_data.AthwartBandwidth = []float32{float32(buffer.AthwartBandwidth) / SCALE_1_F32}
	sensor_data.RangeFilterMin = []float32{float32(buffer.RangeFilterMin)}
	sensor_data.RangeFilterMax = []float32{float32(buffer.RangeFilterMax)}
	sensor_data.DepthFilterMin = []float32{float32(buffer.DepthFilterMin)}
	sensor_data.DepthFilterMax = []float32{float32(buffer.DepthFilterMax)}
	sensor_data.ProjectorType = []uint8{buffer.ProjectorType}

	return sensor_data, err
}

type Seabeam2112 struct {
	Mode                   []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity   []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SsvSource              []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PingGain               []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PulseWidth             []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TransmitterAttenuation []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	NumberAlgorithms       []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	AlgorithmOrder         []string  `tiledb:"dtype=string,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSeabeam2112Specific(reader *bytes.Reader) (sensor_data Seabeam2112, err error) {
	var buffer struct {
		Mode                   uint8
		SurfaceSoundVelocity   uint16
		SsvSource              uint8
		PingGain               uint8
		PulseWidth             uint8
		TransmitterAttenuation uint8
		NumberAlgorithms       uint8
		AlgorithmOrder         [5]byte
		Spare                  [2]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Seabeam2112 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.SurfaceSoundVelocity = []float64{(float64(buffer.SurfaceSoundVelocity) + 130000.0) / SCALE_2_F64}
	sensor_data.SsvSource = []uint8{buffer.SsvSource}
	sensor_data.PingGain = []uint8{buffer.PingGain}
	sensor_data.PulseWidth = []uint8{buffer.PulseWidth}
	sensor_data.TransmitterAttenuation = []uint8{buffer.TransmitterAttenuation}
	sensor_data.NumberAlgorithms = []uint8{buffer.NumberAlgorithms}
	sensor_data.AlgorithmOrder = []string{string(buffer.AlgorithmOrder[:])}

	return sensor_data, err
}

type ElacMkII struct {
	Mode                  []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PingNumber            []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceSoundVelocity  []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PulseLength           []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ReceiverGainStarboard []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ReceiverGainPort      []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeElacMkIISpecific(reader *bytes.Reader) (sensor_data ElacMkII, err error) {
	var buffer struct {
		Mode                  uint8
		PingNumber            uint16
		SurfaceSoundVelocity  uint16
		PulseLength           uint16
		ReceiverGainStarboard uint8
		ReceiverGainPort      uint8
		Spare                 int16
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("ElacMkII sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.Mode = []uint8{buffer.Mode}
	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SurfaceSoundVelocity = []uint16{buffer.SurfaceSoundVelocity}
	sensor_data.PulseLength = []uint16{buffer.PulseLength}
	sensor_data.ReceiverGainStarboard = []uint8{buffer.ReceiverGainStarboard}
	sensor_data.ReceiverGainPort = []uint8{buffer.ReceiverGainPort}

	return sensor_data, err
}

type CmpSass struct {
	Lfreq  []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Lntens []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeCmpSass(reader *bytes.Reader) (sensor_data CmpSass, err error) {
	var buffer struct {
		Lfreq  uint16
		Lntens uint16
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("CmpSASS sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.Lfreq = []float32{float32(buffer.Lfreq) / SCALE_1_F32}
	sensor_data.Lntens = []float32{float32(buffer.Lntens) / SCALE_1_F32}

	return sensor_data, err
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

func DecodeReson8100(reader *bytes.Reader) (sensor_data Reson8100, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Reson8100 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.Latency = []uint16{buffer.Latency}
	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SonarID = []uint16{buffer.SonarID}
	sensor_data.SonarModel = []uint16{buffer.SonarModel}
	sensor_data.Frequency = []uint16{buffer.Frequency}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SurfaceSoundVelocity) / SCALE_1_F32}
	sensor_data.SampleRate = []uint16{buffer.SampleRate}
	sensor_data.PingRate = []uint16{buffer.PingRate}
	sensor_data.Mode = []uint16{buffer.Mode}
	sensor_data.Range = []uint16{buffer.Range}
	sensor_data.TransmitPower = []uint16{buffer.TransmitPower}
	sensor_data.ReceiveGain = []uint16{buffer.ReceiveGain}
	sensor_data.PulseWidth = []uint16{buffer.PulseWidth}
	sensor_data.TvgSpreading = []uint8{buffer.TvgSpreading}
	sensor_data.TvgAbsorption = []uint8{buffer.TvgAbsorption}
	sensor_data.ForeAftBandwidth = []float32{float32(buffer.ForeAftBandwidth) / SCALE_1_F32}
	sensor_data.AthwartBandwidth = []float32{float32(buffer.AthwartBandwidth) / SCALE_1_F32}
	sensor_data.ProjectorType = []uint8{buffer.ProjectorType}
	sensor_data.ProjectorAngle = []int16{buffer.ProjectorAngle}
	sensor_data.RangeFilterMin = []float32{float32(buffer.RangeFilterMin)}
	sensor_data.RangeFilterMax = []float32{float32(buffer.RangeFilterMax)}
	sensor_data.DepthFilterMin = []float32{float32(buffer.DepthFilterMin)}
	sensor_data.DepthFilterMax = []float32{float32(buffer.DepthFilterMax)}
	sensor_data.FiltersActive = []uint8{buffer.FiltersActive}
	sensor_data.Temperature = []uint16{buffer.Temperature}
	sensor_data.BeamSpacing = []float32{float32(buffer.BeamSpacing) / SCALE_4_F32}

	return sensor_data, err
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
	RunTimeDgTime                  [][]time.Time `tiledb:"dtype=datetime_ns,ftype=attr,var" filters:"zstd(level=16)"`
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
	RunTimeReceiveBandwidth        [][]uint16    `tiledb:"dtype=uint16,ftype=attr,var" filters:"zstd(level=16)"`
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

func DecodeEm3Specific(reader *bytes.Reader) (sensor_data Em3, err error) {
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
	receive_bandwidth := make([]uint16, 0, 2)
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

	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("EM3 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	// base set of values
	sensor_data.ModelNumber = []uint16{buffer.ModelNumber}
	sensor_data.PingNumber = []uint16{buffer.PingNumber}
	sensor_data.SerialNumber = []uint16{buffer.SerialNumber}
	sensor_data.SurfaceSoundVelocity = []float32{float32(buffer.SerialNumber) / SCALE_1_F32}
	sensor_data.TransducerDepth = []float32{float32(buffer.TransducerDepth) / SCALE_2_F32}
	sensor_data.ValidBeams = []uint16{buffer.ValidBeams}
	sensor_data.SampleRate = []uint16{buffer.SampleRate}
	sensor_data.DepthDifference = []float32{float32(buffer.DepthDifference) / SCALE_2_F32}
	sensor_data.OffsetMultiplier = []uint8{buffer.OffsetMultiplier}

	// runtime values
	if (buffer.RunTimeID & 0x00000001) != 0 {
		err = binary.Read(reader, binary.BigEndian, &rt1)
		if err != nil {
			errn := errors.New("EM3 sensor")
			err = errors.Join(err, ErrSensorMetadata, errn)
			return sensor_data, err
		}
		model_number = append(model_number, rt1.ModelNumber)
		dg_time = append(dg_time, time.Unix(int64(rt1.TvSec), int64(rt1.TvNSec)).UTC())
		ping_number = append(ping_number, rt1.PingNumber)
		serial_number = append(serial_number, rt1.SerialNumber)
		system_status = append(system_status, rt1.SystemStatus)
		mode = append(mode, rt1.Mode)
		filter_id = append(filter_id, rt1.FilterID)
		min_depth = append(min_depth, float32(rt1.MinDepth))
		max_depth = append(max_depth, float32(rt1.MaxDepth))
		absorption = append(absorption, float32(rt2.Absorption)/SCALE_2_F32)
		transmit_pulse_length = append(transmit_pulse_length, float32(rt1.TransmitPulseLength))
		transmit_beam_width = append(transmit_beam_width, float32(rt1.TransmitBeamWidth)/SCALE_1_F32)
		power_reduction = append(power_reduction, rt1.PowerReduction)
		receive_beamwidth = append(receive_beamwidth, float32(rt1.ReceiveBeamWidth)/SCALE_1_F32)
		receive_bandwidth = append(receive_bandwidth, uint16(rt1.ReceiveBandwidth)*50)
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
			err = binary.Read(reader, binary.BigEndian, &rt2)
			if err != nil {
				errn := errors.New("EM3 sensor")
				err = errors.Join(err, ErrSensorMetadata, errn)
				return sensor_data, err
			}
			model_number = append(model_number, rt2.ModelNumber)
			dg_time = append(dg_time, time.Unix(int64(rt2.TvSec), int64(rt2.TvNSec)).UTC())
			ping_number = append(ping_number, rt2.PingNumber)
			serial_number = append(serial_number, rt2.SerialNumber)
			system_status = append(system_status, rt2.SystemStatus)
			mode = append(mode, rt2.Mode)
			filter_id = append(filter_id, rt2.FilterID)
			min_depth = append(min_depth, float32(rt2.MinDepth))
			max_depth = append(max_depth, float32(rt2.MaxDepth))
			absorption = append(absorption, float32(rt2.Absorption)/SCALE_2_F32)
			transmit_pulse_length = append(transmit_pulse_length, float32(rt2.TransmitPulseLength))
			transmit_beam_width = append(transmit_beam_width, float32(rt2.TransmitBeamWidth)/SCALE_1_F32)
			power_reduction = append(power_reduction, rt2.PowerReduction)
			receive_beamwidth = append(receive_beamwidth, float32(rt2.ReceiveBeamWidth)/SCALE_1_F32)
			receive_bandwidth = append(receive_bandwidth, uint16(rt2.ReceiveBandwidth)*50)
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
	sensor_data.RunTimeSystemStatus = [][]uint32{system_status}
	sensor_data.RunTimeMode = [][]uint8{mode}
	sensor_data.RunTimeFilterID = [][]uint8{filter_id}
	sensor_data.RunTimeMinDepth = [][]float32{min_depth}
	sensor_data.RunTimeMaxDepth = [][]float32{max_depth}
	sensor_data.RunTimeAbsorption = [][]float32{absorption}
	sensor_data.RunTimeTransmitPulseLength = [][]float32{transmit_pulse_length}
	sensor_data.RunTimeTransmitBeamWidth = [][]float32{transmit_beam_width}
	sensor_data.RunTimePowerReduction = [][]uint8{power_reduction}
	sensor_data.RunTimeReceiveBeamWidth = [][]float32{receive_beamwidth}
	sensor_data.RunTimeReceiveBandwidth = [][]uint16{receive_bandwidth}
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

	return sensor_data, err
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
	VehicleDepth                      []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
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

func DecodeEm4Specific(reader *bytes.Reader) (sensor_data Em4, err error) {

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
			Spare                  [16]byte
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
			Spare           []byte
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
			Spare           [16]byte
		} // 40 bytes
		spare_buffer struct {
			Spare [16]byte
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
			Spare                         [16]byte
		} // 63 bytes
		proc_buffer struct {
			ProcessorUnitCpuLoad              uint8
			ProcessorUnitSensorStatus         uint16
			ProcessorUnitAchievedPortCoverage uint8
			ProcessorUnitAchievedStbdCoverage uint8
			ProcessorUnitYawStabilization     int16
			Spare                             [16]byte
		} // 23 bytes
	)

	// first 46 bytes
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("EM4 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
	// n_bytes += 46

	// sector arrays
	// what if TransmitSectors == 0???
	for i := uint16(0); i < buffer.TransmitSectors; i++ {
		err = binary.Read(reader, binary.BigEndian, &sector_buffer_base)
		if err != nil {
			errn := errors.New("EM4 sensor")
			err = errors.Join(err, ErrSensorMetadata, errn)
			return sensor_data, err
		}
		sector_buffer.TiltAngle = append(
			sector_buffer.TiltAngle,
			float32(sector_buffer_base.TiltAngle)/SCALE_2_F32,
		)
		sector_buffer.FocusRange = append(
			sector_buffer.FocusRange,
			float32(sector_buffer_base.FocusRange)/SCALE_1_F32,
		)
		sector_buffer.SignalLength = append(
			sector_buffer.SignalLength,
			float64(sector_buffer_base.SignalLength)/SCALE_6_F64,
		)
		sector_buffer.TransmitDelay = append(
			sector_buffer.TransmitDelay,
			float64(sector_buffer_base.TransmitDelay)/SCALE_6_F64,
		)
		sector_buffer.CenterFrequency = append(
			sector_buffer.CenterFrequency,
			float64(sector_buffer_base.CenterFrequency)/SCALE_3_F64,
		)
		sector_buffer.MeanAbsorption = append(
			sector_buffer.MeanAbsorption,
			float32(sector_buffer_base.MeanAbsorption)/SCALE_2_F32,
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
			float64(sector_buffer_base.SignalBandwith)/SCALE_3_F64,
		)
	}

	// spare 16 bytes
	err = binary.Read(reader, binary.BigEndian, &spare_buffer)
	if err != nil {
		errn := errors.New("EM4 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	// next 63 bytes for the RunTime info
	err = binary.Read(reader, binary.BigEndian, &runtime_buffer)
	if err != nil {
		errn := errors.New("EM4 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	// next 23 bytes for the processing unit info
	err = binary.Read(reader, binary.BigEndian, &proc_buffer)
	if err != nil {
		errn := errors.New("EM4 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	// populate generic
	sensor_data.ModelNumber = []uint16{buffer.ModelNumber}
	sensor_data.PingCounter = []uint16{buffer.PingCounter}
	sensor_data.SerialNumber = []uint16{buffer.SerialNumber}
	sensor_data.SurfaceVelocity = []float32{float32(buffer.SurfaceVelocity) / SCALE_1_F32}
	sensor_data.TransducerDepth = []float64{float64(buffer.TransducerDepth) / float64(20000)}
	sensor_data.ValidDetections = []uint16{buffer.ValidDetections}
	sensor_data.SamplingFrequency = []float64{float64(buffer.SamplingFrequency1) + float64(buffer.SamplingFrequency2)/float64(4_000_000_000)}
	sensor_data.DopplerCorrectionScale = []uint32{buffer.DopplerCorrectionScale}
	sensor_data.VehicleDepth = []float64{float64(buffer.VehicleDepth) / SCALE_3_F64}
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
	sensor_data.RunTimeAbsorption = []float32{float32(runtime_buffer.RunTimeAbsorption) / SCALE_2_F32}
	sensor_data.RunTimeTransmitPulseLength = []float32{float32(runtime_buffer.RunTimeTransmitPulseLength)}
	sensor_data.RunTimeTransmitBeamWidth = []float32{float32(runtime_buffer.RunTimeTransmitBeamWidth) / SCALE_1_F32}
	sensor_data.RunTimeTransmitPowerReduction = []uint8{runtime_buffer.RunTimeTransmitPowerReduction}
	sensor_data.RunTimeReceiveBeamWidth = []float32{float32(runtime_buffer.RunTimeReceiveBeamWidth) / SCALE_1_F32}
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
	sensor_data.RunTimeTransmitAlongTilt = []float32{float32(runtime_buffer.RunTimeTransmitAlongTilt) / SCALE_2_F32}
	sensor_data.RunTimeFilterID2 = []uint8{runtime_buffer.RunTimeFilterID2}

	// populate processor unit info
	sensor_data.ProcessorUnitCpuLoad = []uint8{proc_buffer.ProcessorUnitCpuLoad}
	sensor_data.ProcessorUnitSensorStatus = []uint16{proc_buffer.ProcessorUnitSensorStatus}
	sensor_data.ProcessorUnitAchievedPortCoverage = []uint8{proc_buffer.ProcessorUnitAchievedPortCoverage}
	sensor_data.ProcessorUnitAchievedStbdCoverage = []uint8{proc_buffer.ProcessorUnitAchievedStbdCoverage}
	sensor_data.ProcessorUnitYawStabilization = []float32{float32(proc_buffer.ProcessorUnitYawStabilization) / SCALE_2_F32}

	return sensor_data, err
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

func DecodeGeoSwathPlusSpecific(reader *bytes.Reader) (sensor_data GeoSwathPlus, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("GeoSwathPlus sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.DataSource = []uint16{buffer.DataSource}
	sensor_data.Side = []uint16{buffer.Side}
	sensor_data.ModelNumber = []uint16{buffer.ModelNumber}
	sensor_data.Frequency = []float32{float32(buffer.Frequency) * SCALE_1_F32}
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
	sensor_data.SampleRate = []float32{float32(buffer.SampleRate) * SCALE_1_F32}
	sensor_data.PulseLength = []float32{float32(buffer.PulseLength)}
	sensor_data.PingLength = []uint16{buffer.PingLength}
	sensor_data.TransmitPower = []uint16{buffer.TransmitPower}
	sensor_data.SidescanGainChannel = []uint16{buffer.SidescanGainChannel}
	sensor_data.Stabilization = []uint16{buffer.Stabilization}
	sensor_data.GpsQuality = []uint16{buffer.GpsQuality}
	sensor_data.RangeUncertainty = []float32{float32(buffer.RangeUncertainty) / SCALE_3_F32}
	sensor_data.AngleUncertainty = []float32{float32(buffer.AngleUncertainty) / SCALE_2_F32}

	return sensor_data, err
}

// DecodeKlein5410Bss TODO; change DataSource and Side types from uint16 to uint8.
// Seems a waste to store 2 bytes for data that is only a 0 or 1.
type Klein5410Bss struct {
	DataSource        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Side              []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ModelNumber       []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	AcousticFrequency []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SamplingFrequency []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	PingNumber        []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	NumSamples        []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	NumRaaSamples     []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ErrorFlags        []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	Range             []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	FishDepth         []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	FishAltitude      []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SoundSpeed        []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TransmitWaveform  []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Altimeter         []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RawDataConfig     []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeKlein5410BssSpecific(reader *bytes.Reader) (sensor_data Klein5410Bss, err error) {
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
		Spare             [32]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Klein5410BSS sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.DataSource = []uint16{buffer.DataSource}
	sensor_data.Side = []uint16{buffer.Side}
	sensor_data.ModelNumber = []uint16{buffer.ModelNumber}
	sensor_data.AcousticFrequency = []float64{float64(buffer.AcousticFrequency) / SCALE_3_F64}
	sensor_data.SamplingFrequency = []float64{float64(buffer.SamplingFrequency) / SCALE_3_F64}
	sensor_data.PingNumber = []uint32{buffer.PingNumber}
	sensor_data.NumSamples = []uint32{buffer.NumSamples}
	sensor_data.NumRaaSamples = []uint32{buffer.NumRaaSamples}
	sensor_data.ErrorFlags = []uint32{buffer.ErrorFlags}
	sensor_data.Range = []uint32{buffer.Range}
	sensor_data.FishDepth = []float64{float64(buffer.FishDepth) / SCALE_3_F64}
	sensor_data.FishAltitude = []float64{float64(buffer.FishAltitude) / SCALE_3_F64}
	sensor_data.SoundSpeed = []float64{float64(buffer.SoundSpeed) / SCALE_3_F64}
	sensor_data.TransmitWaveform = []uint16{buffer.TransmitWaveform}
	sensor_data.Altimeter = []uint16{buffer.Altimeter}
	sensor_data.RawDataConfig = []uint32{buffer.RawDataConfig}

	return sensor_data, err
}

type Reson7100 struct {
	ProtocolVersion                   []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	DeviceID                          []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	MajorSerialNumber                 []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	MinorSerialNumber                 []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	PingNumber                        []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	MultiPingSequence                 []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Frequency                         []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SampleRate                        []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ReceiverBandwidth                 []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxPulseWidth                      []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxPulseTypeID                     []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	TxPulseEnvlpID                    []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	TxPulseEnvlpParam                 []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	MaxPingRate                       []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	PingPeriod                        []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Range                             []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Power                             []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Gain                              []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ControlFlags                      []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorID                       []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorSteerAnglVert            []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ProjectorSteerAnglHorz            []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamWidthVert            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamWidthHorz            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamFocalPt              []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
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
	Absorption                        []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SoundVelocity                     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Spreading                         []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RawDataFrom7027                   []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SvSource                          []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	LayerCompFlag                     []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	// TxPulseReserved                   []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeReson7100Specific(reader *bytes.Reader) (sensor_data Reson7100, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Reson7100 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.ProtocolVersion = []uint16{buffer.ProtocolVersion}
	sensor_data.DeviceID = []uint32{buffer.DeviceID}
	sensor_data.MajorSerialNumber = []uint32{buffer.MajorSerialNumber}
	sensor_data.MinorSerialNumber = []uint32{buffer.MinorSerialNumber}
	sensor_data.PingNumber = []uint32{buffer.PingNumber}
	sensor_data.MultiPingSequence = []uint16{buffer.MultiPingSequence}
	sensor_data.Frequency = []float64{float64(buffer.Frequency) / SCALE_3_F64}
	sensor_data.SampleRate = []float64{float64(buffer.SampleRate) / SCALE_4_F64}
	sensor_data.ReceiverBandwidth = []float64{float64(buffer.ReceiverBandwidth) / SCALE_4_F64}
	sensor_data.TxPulseWidth = []float64{float64(buffer.TxPulseWidth) / SCALE_7_F64}
	sensor_data.TxPulseTypeID = []uint32{buffer.TxPulseTypeID}
	sensor_data.TxPulseEnvlpID = []uint32{buffer.TxPulseEnvlpID}
	sensor_data.TxPulseEnvlpParam = []float64{float64(buffer.TxPulseEnvlpParam) / SCALE_2_F64}
	sensor_data.MaxPingRate = []float64{float64(buffer.MaxPingRate) / SCALE_6_F64}
	sensor_data.PingPeriod = []float64{float64(buffer.PingPeriod) / SCALE_6_F64}
	sensor_data.Range = []float64{float64(buffer.Range) / SCALE_2_F64}
	sensor_data.Power = []float64{float64(buffer.Power) / SCALE_2_F64}
	sensor_data.Gain = []float64{float64(buffer.Gain) / SCALE_2_F64}
	sensor_data.ControlFlags = []uint32{buffer.ControlFlags}
	sensor_data.ProjectorID = []uint32{buffer.ProjectorID}
	sensor_data.ProjectorSteerAnglVert = []float64{float64(buffer.ProjectorSteerAnglVert) / SCALE_3_F64}
	sensor_data.ProjectorSteerAnglHorz = []float64{float64(buffer.ProjectorSteerAnglHorz) / SCALE_3_F64}
	sensor_data.ProjectorBeamWidthVert = []float32{float32(buffer.ProjectorBeamWidthVert) / SCALE_2_F32}
	sensor_data.ProjectorBeamWidthHorz = []float32{float32(buffer.ProjectorBeamWidthHorz) / SCALE_2_F32}
	sensor_data.ProjectorBeamFocalPt = []float64{float64(buffer.ProjectorBeamFocalPt) / SCALE_2_F64}
	sensor_data.ProjectorBeamWeightingWindowType = []uint32{buffer.ProjectorBeamWeightingWindowType}
	sensor_data.ProjectorBeamWeightingWindowParam = []uint32{buffer.ProjectorBeamWeightingWindowParam}
	sensor_data.TransmitFlags = []uint32{buffer.TransmitFlags}
	sensor_data.HydrophoneID = []uint32{buffer.HydrophoneID}
	sensor_data.ReceivingBeamWeightingWindowType = []uint32{buffer.ReceivingBeamWeightingWindowType}
	sensor_data.ReceivingBeamWeightingWindowParam = []uint32{buffer.ReceivingBeamWeightingWindowParam}
	sensor_data.ReceiveFlags = []uint32{buffer.ReceiveFlags}
	sensor_data.ReceiveBeamWidth = []float32{float32(buffer.ReceiveBeamWidth) / SCALE_2_F32}
	sensor_data.RangeFiltMin = []float32{float32(buffer.RangeFiltMin) / SCALE_1_F32}
	sensor_data.RangeFiltMax = []float32{float32(buffer.RangeFiltMax) / SCALE_1_F32}
	sensor_data.DepthFiltMin = []float32{float32(buffer.DepthFiltMin) / SCALE_1_F32}
	sensor_data.DepthFiltMax = []float32{float32(buffer.DepthFiltMax) / SCALE_1_F32}
	sensor_data.Absorption = []float64{float64(buffer.Absorption) / SCALE_3_F64}
	sensor_data.SoundVelocity = []float32{float32(buffer.SoundVelocity) / SCALE_1_F32}
	sensor_data.Spreading = []float64{float64(buffer.Spreading) / SCALE_3_F64}
	sensor_data.RawDataFrom7027 = []uint8{buffer.RawDataFrom7027}
	sensor_data.SvSource = []uint8{buffer.SvSource}
	sensor_data.LayerCompFlag = []uint8{buffer.LayerCompFlag}

	return sensor_data, err
}

type Em3Raw struct {
	ModelNumber                       []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingCounter                       []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SerialNumber                      []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SurfaceVelocity                   []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransducerDepth                   []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ValidDetections                   []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SamplingFrequency                 []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	VehicleDepth                      []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	DepthDifference                   []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	OffsetMultiplier                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TransmitSectors                   []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TiltAngle                         [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	FocusRange                        [][]float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	SignalLength                      [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	TransmitDelay                     [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	CenterFrequency                   [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	WaveformID                        [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	SectorNumber                      [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	SignalBandwidth                   [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	RunTimeModelNumber                []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RunTimeDgTime                     []time.Time `tiledb:"dtype=datetime_ns,ftype=attr" filters:"zstd(level=16)"`
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

func DecodeEm3RawSpecific(reader *bytes.Reader) (sensor_data Em3Raw, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("EM3Raw sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
	sensor_data.ModelNumber = []uint16{buffer.ModelNumber}
	sensor_data.PingCounter = []uint16{buffer.PingCounter}
	sensor_data.SerialNumber = []uint16{buffer.SerialNumber}
	sensor_data.SurfaceVelocity = []float32{float32(buffer.SurfaceVelocity) / SCALE_1_F32}
	sensor_data.TransducerDepth = []float64{float64(buffer.TransducerDepth) / 20_000.0}
	sensor_data.ValidDetections = []uint16{buffer.ValidDetections}
	sensor_data.SamplingFrequency = []float64{float64(buffer.SamplingFrequency1) + float64(buffer.SamplingFrequency2)/4_000_000_000.0}
	sensor_data.VehicleDepth = []float64{float64(buffer.VehicleDepth) / SCALE_3_F64}
	sensor_data.DepthDifference = []float32{float32(buffer.DepthDifference) / SCALE_2_F32}
	sensor_data.OffsetMultiplier = []uint8{buffer.OffsetMultiplier}
	sensor_data.TransmitSectors = []uint16{buffer.TransmitSectors}

	// second block (variable length arrays)
	nsectors := int(buffer.TransmitSectors)
	tilt_angle := make([]float32, 0, nsectors)
	focus_range := make([]float32, 0, nsectors)
	signal_length := make([]float64, 0, nsectors)
	transmit_delay := make([]float64, 0, nsectors)
	centre_frequency := make([]float64, 0, nsectors)
	waveformID := make([]uint8, 0, nsectors)
	sector_number := make([]uint8, 0, nsectors)
	signal_bandwidth := make([]float64, 0, nsectors)

	for i := 0; i < nsectors; i++ {
		err = binary.Read(reader, binary.BigEndian, &var_buff)
		if err != nil {
			errn := errors.New("EM3Raw sensor")
			err = errors.Join(err, ErrSensorMetadata, errn)
			return sensor_data, err
		}
		tilt_angle = append(tilt_angle, float32(var_buff.TiltAngle)/SCALE_2_F32)
		focus_range = append(focus_range, float32(var_buff.FocusRange)/SCALE_1_F32)
		signal_length = append(signal_length, float64(var_buff.SignalLength)/SCALE_6_F64)
		transmit_delay = append(transmit_delay, float64(var_buff.TransmitDelay)/SCALE_6_F64)
		centre_frequency = append(centre_frequency, float64(var_buff.CenterFrequency)/SCALE_3_F64)
		waveformID = append(waveformID, var_buff.WaveformID)
		sector_number = append(sector_number, var_buff.SectorNumber)
		signal_bandwidth = append(signal_bandwidth, float64(var_buff.SignalBandwidth)/SCALE_3_F64)
	}

	sensor_data.TiltAngle = [][]float32{tilt_angle}
	sensor_data.FocusRange = [][]float32{focus_range}
	sensor_data.SignalLength = [][]float64{signal_length}
	sensor_data.TransmitDelay = [][]float64{transmit_delay}
	sensor_data.CenterFrequency = [][]float64{centre_frequency}
	sensor_data.WaveformID = [][]uint8{waveformID}
	sensor_data.SectorNumber = [][]uint8{sector_number}
	sensor_data.SignalBandwidth = [][]float64{signal_bandwidth}

	// third block (runtime)
	err = binary.Read(reader, binary.BigEndian, &rt_buff)
	if err != nil {
		errn := errors.New("EM3Raw sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
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
	sensor_data.RunTimeAbsorption = []float32{float32(rt_buff.RunTimeAbsorption) / SCALE_2_F32}
	sensor_data.RunTimeTxPulseLength = []uint16{rt_buff.RunTimeTxPulseLength}
	sensor_data.RunTimeTxBeamWidth = []float32{float32(rt_buff.RunTimeTxBeamWidth) / SCALE_1_F32}
	sensor_data.RunTimeTxPowerReMax = []uint8{rt_buff.RunTimeTxPowerReMax}
	sensor_data.RunTimeRxBeamWidth = []float32{float32(rt_buff.RunTimeRxBeamWidth) / SCALE_1_F32}
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
		err = binary.Read(reader, binary.BigEndian, &RunTimeDurotongSpeed)
		if err != nil {
			errn := errors.New("EM3Raw sensor")
			err = errors.Join(err, ErrSensorMetadata, errn)
			return sensor_data, err
		}
		sensor_data.RunTimeDurotongSpeed = []float32{float32(RunTimeDurotongSpeed) / SCALE_1_F32}
		sensor_data.RunTimeTxAlongTilt = []float32{NULL_FLOAT32_ZERO}
	case 300:
		sensor_data.RunTimeDurotongSpeed = []float32{NULL_FLOAT32_ZERO}
		sensor_data.RunTimeTxAlongTilt = []float32{NULL_FLOAT32_ZERO}
	case 120:
		sensor_data.RunTimeDurotongSpeed = []float32{NULL_FLOAT32_ZERO}
		sensor_data.RunTimeTxAlongTilt = []float32{NULL_FLOAT32_ZERO}
	case 3020:
		sensor_data.RunTimeDurotongSpeed = []float32{NULL_FLOAT32_ZERO}
		err = binary.Read(reader, binary.BigEndian, &RunTimeTxAlongTilt)
		if err != nil {
			errn := errors.New("EM3Raw sensor")
			err = errors.Join(err, ErrSensorMetadata, errn)
			return sensor_data, err
		}
		sensor_data.RunTimeTxAlongTilt = []float32{float32(RunTimeTxAlongTilt) / SCALE_2_F32}
	default:
		err = binary.Read(reader, binary.BigEndian, &Spare)
		if err != nil {
			errn := errors.New("EM3Raw sensor")
			err = errors.Join(err, ErrSensorMetadata, errn)
			return sensor_data, err
		}
	}

	// appears that this piece is incomplete in the C-code and awaiting info from KM
	// regarding final datagram documentation.
	// This was captured back in 2009, and it is now 2024 with no updates
	// So merely replicating what they've constructed
	switch rt_buff.RunTimeModelNumber {
	default:
		err = binary.Read(reader, binary.BigEndian, &RunTimeHiLoAbsorptionRatio)
		if err != nil {
			errn := errors.New("EM3Raw sensor")
			err = errors.Join(err, ErrSensorMetadata, errn)
			return sensor_data, err
		}
		sensor_data.RunTimeHiLoAbsorptionRatio = []uint8{RunTimeHiLoAbsorptionRatio}
	}

	// fourth block (process unit)
	err = binary.Read(reader, binary.BigEndian, &pu_buff)
	if err != nil {
		errn := errors.New("EM3Raw sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
	sensor_data.PuStatusPuCpuLoad = []uint8{pu_buff.PuStatusPuCpuLoad}
	sensor_data.PuStatusSensorStatus = []uint16{pu_buff.PuStatusSensorStatus}
	sensor_data.PuStatusAchievedPortCoverage = []uint8{pu_buff.PuStatusAchievedPortCoverage}
	sensor_data.PuStatusAchievedStarboardCoverage = []uint8{pu_buff.PuStatusAchievedStarboardCoverage}
	sensor_data.PuStatusYawStabilization = []float32{float32(pu_buff.PuStatusYawStabilization) / SCALE_2_F32}

	return sensor_data, err
}

type DeltaT struct {
	FileExtension        []string    `tiledb:"dtype=string,ftype=attr" filters:"zstd(level=16)"`
	Version              []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PingByteSize         []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	InterrogationTime    []time.Time `tiledb:"dtype=datetime_ns,ftype=attr" filters:"zstd(level=16)"`
	SamplesPerBeam       []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SectorSize           []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	StartAngle           []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	AngleIncrement       []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	AcousticRange        []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	AcousticFrequency    []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SoundVelocity        []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RangeResolution      []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	ProfileTiltAngle     []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RepetitionRate       []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingNumber           []uint32    `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	IntensityFlag        []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PingLatency          []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	DataLatency          []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SampleRateFlag       []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	OptionsFlag          []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	NumberPingsAveraged  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	CenterPingTimeOffset []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	UserDefinedByte      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Altitude             []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ExternalSensorFlags  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PulseLength          []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ForeAftBeamwidth     []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	AthwartBeamwidth     []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeDeltaTSpecific(reader *bytes.Reader) (sensor_data DeltaT, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("DeltaT sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.FileExtension = []string{string(buffer.FileExtension[:])}
	sensor_data.Version = []uint8{buffer.Version}
	sensor_data.PingByteSize = []uint16{buffer.PingByteSize}
	sensor_data.InterrogationTime = []time.Time{time.Unix(int64(buffer.TvSec), int64(buffer.TvNsec)).UTC()}
	sensor_data.SamplesPerBeam = []uint16{buffer.SamplesPerBeam}
	sensor_data.SectorSize = []uint16{buffer.SectorSize}
	sensor_data.StartAngle = []float32{(float32(buffer.StartAngle) / SCALE_2_F32) - 180.0}
	sensor_data.AngleIncrement = []float32{float32(buffer.AngleIncrement) / SCALE_2_F32}
	sensor_data.AcousticRange = []uint16{buffer.AcousticRange}
	sensor_data.AcousticFrequency = []uint16{buffer.AcousticFrequency}
	sensor_data.SoundVelocity = []float32{float32(buffer.SoundVelocity) / SCALE_1_F32}
	sensor_data.RangeResolution = []uint16{buffer.RangeResolution}
	sensor_data.ProfileTiltAngle = []float32{float32(buffer.ProfileTiltAngle) - 180.0}
	sensor_data.RepetitionRate = []uint16{buffer.RepetitionRate}
	sensor_data.PingNumber = []uint32{buffer.PingNumber}
	sensor_data.IntensityFlag = []uint8{buffer.IntensityFlag}
	sensor_data.PingLatency = []float32{float32(buffer.PingLatency) / SCALE_4_F32}
	sensor_data.DataLatency = []float32{float32(buffer.DataLatency) / SCALE_4_F32}
	sensor_data.SampleRateFlag = []uint8{buffer.SampleRateFlag}
	sensor_data.OptionsFlag = []uint8{buffer.OptionsFlag}
	sensor_data.NumberPingsAveraged = []uint8{buffer.NumberPingsAveraged}
	sensor_data.CenterPingTimeOffset = []float32{float32(buffer.CenterPingTimeOffset) / SCALE_4_F32}
	sensor_data.UserDefinedByte = []uint8{buffer.UserDefinedByte}
	sensor_data.Altitude = []float64{float64(buffer.Altitude) / SCALE_2_F64}
	sensor_data.ExternalSensorFlags = []uint8{buffer.ExternalSensorFlags}
	sensor_data.PulseLength = []float64{float64(buffer.PulseLength) / SCALE_6_F64}
	sensor_data.ForeAftBeamwidth = []float32{float32(buffer.ForeAftBeamwidth) / SCALE_1_F32}
	sensor_data.AthwartBeamwidth = []float32{float32(buffer.AthwartBeamwidth) / SCALE_1_F32}

	return sensor_data, err
}

type R2Sonic struct {
	ModelNumber      []string    `tiledb:"dtype=string,ftype=attr" filters:"zstd(level=16)"`
	SerialNumber     []string    `tiledb:"dtype=string,ftype=attr" filters:"zstd(level=16)"`
	DgTime           []time.Time `tiledb:"dtype=datetime_ns,ftype=attr" filters:"zstd(level=16)"`
	PingNumber       []uint32    `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	PingPeriod       []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SoundSpeed       []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Frequency        []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxPower          []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxPulseWidth     []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxBeamWidthVert  []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxBeamWidthHoriz []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxSteeringVert   []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxSteeringHoriz  []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxMiscInfo       []uint32    `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	RxBandwidth      []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxSampleRate     []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxRange          []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxGain           []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxSpreading      []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxAbsorption     []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxMountTilt      []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxMiscInfo       []uint32    `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	NumberBeams      []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	A0MoreInfo       [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	A2MoreInfo       [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	G0DepthGateMin   []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	G0DepthGateMax   []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	G0DepthGateSlope []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeR2SonicSpecific(reader *bytes.Reader) (sensor_data R2Sonic, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer1)
	if err != nil {
		errn := errors.New("R2Sonic sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
	sensor_data.ModelNumber = []string{string(buffer1.ModelNumber[:])}
	sensor_data.SerialNumber = []string{string(buffer1.SerialNumber[:])}
	sensor_data.DgTime = []time.Time{time.Unix(int64(buffer1.TvSec), int64(buffer1.TvNsec)).UTC()}
	sensor_data.PingNumber = []uint32{buffer1.PingNumber}
	sensor_data.PingPeriod = []float64{float64(buffer1.PingPeriod) / SCALE_6_F64}
	sensor_data.SoundSpeed = []float64{float64(buffer1.SoundSpeed) / SCALE_2_F64}
	sensor_data.Frequency = []float64{float64(buffer1.Frequency) / SCALE_3_F64}
	sensor_data.TxPower = []float64{float64(buffer1.TxPower) / SCALE_2_F64}
	sensor_data.TxPulseWidth = []float64{float64(buffer1.TxPulseWidth) / SCALE_7_F64}
	sensor_data.TxBeamWidthVert = []float64{float64(buffer1.TxBeamWidthVert) / SCALE_6_F64}
	sensor_data.TxBeamWidthHoriz = []float64{float64(buffer1.TxBeamWidthHoriz) / SCALE_6_F64}
	sensor_data.TxSteeringVert = []float64{float64(buffer1.TxSteeringVert) / SCALE_6_F64}
	sensor_data.TxSteeringHoriz = []float64{float64(buffer1.TxSteeringHoriz) / SCALE_6_F64}
	sensor_data.TxMiscInfo = []uint32{buffer1.TxMiscInfo}
	sensor_data.RxBandwidth = []float64{float64(buffer1.RxBandwidth) / SCALE_4_F64}
	sensor_data.RxSampleRate = []float64{float64(buffer1.RxSampleRate) / SCALE_3_F64}
	sensor_data.RxRange = []float64{float64(buffer1.RxRange) / SCALE_5_F64}
	sensor_data.RxGain = []float64{float64(buffer1.RxGain) / SCALE_2_F64}
	sensor_data.RxSpreading = []float64{float64(buffer1.RxSpreading) / SCALE_3_F64}
	sensor_data.RxAbsorption = []float64{float64(buffer1.RxAbsorption) / SCALE_3_F64}
	sensor_data.RxMountTilt = []float64{float64(buffer1.RxMountTilt) / SCALE_6_F64}
	sensor_data.RxMiscInfo = []uint32{buffer1.RxMiscInfo}
	sensor_data.NumberBeams = []uint16{buffer1.NumberBeams}

	// block two (var length arrays)
	err = binary.Read(reader, binary.BigEndian, &var_buf)
	if err != nil {
		errn := errors.New("R2Sonic sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
	A0MoreInfo := make([]float64, 0, 6)
	A2MoreInfo := make([]float64, 0, 6)

	for i := 0; i < 6; i++ {
		A0MoreInfo = append(A0MoreInfo, float64(var_buf.A0MoreInfo[i])/SCALE_6_F64)
		A2MoreInfo = append(A2MoreInfo, float64(var_buf.A2MoreInfo[i])/SCALE_6_F64)
	}

	sensor_data.A0MoreInfo = [][]float64{A0MoreInfo}
	sensor_data.A2MoreInfo = [][]float64{A2MoreInfo}

	// block three
	err = binary.Read(reader, binary.BigEndian, &buffer2)
	if err != nil {
		errn := errors.New("R2Sonic sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
	sensor_data.G0DepthGateMin = []float64{float64(buffer2.G0DepthGateMin) / SCALE_6_F64}
	sensor_data.G0DepthGateMax = []float64{float64(buffer2.G0DepthGateMax) / SCALE_6_F64}
	sensor_data.G0DepthGateSlope = []float64{float64(buffer2.G0DepthGateSlope) / SCALE_6_F64}

	return sensor_data, err
}

type ResonTSeries struct {
	ProtocolVersion                   []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	DeviceID                          []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	NumberDevices                     []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	SystemEnumerator                  []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	MajorSerialNumber                 []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	MinorSerialNumber                 []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	PingNumber                        []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	MultiPingSequence                 []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Frequency                         []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SampleRate                        []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ReceiverBandwidth                 []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxPulseWidth                      []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxPulseTypeID                     []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	TxPulseEnvlpID                    []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	TxPulseEnvlpParam                 []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxPulseMode                       []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	MaxPingRate                       []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	PingPeriod                        []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Range                             []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Power                             []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Gain                              []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ControlFlags                      []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorID                       []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorSteerAnglVert            []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ProjectorSteerAnglHorz            []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamWidthVert            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamWidthHorz            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamFocalPt              []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamWeightingWindowType  []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ProjectorBeamWeightingWindowParam []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	TransmitFlags                     []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	HydrophoneID                      []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ReceivingBeamWeightingWindowType  []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ReceivingBeamWeightingWindowParam []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ReceiveFlags                      []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	ReceiveBeamWidth                  []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RangeFiltMin                      []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RangeFiltMax                      []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	DepthFiltMin                      []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	DepthFiltMax                      []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Absorption                        []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SoundVelocity                     []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SvSource                          []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	Spreading                         []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	BeamSpacingMode                   []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	SonarSourceMode                   []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	CoverageMode                      []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	CoverageAngle                     []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	HorizontalReceiverSteeringAngle   []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	UncertaintyType                   []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	TransmitterSteeringAngle          []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	AppliedRoll                       []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	DetectionAlgorithm                []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	DetectionFlags                    []uint32  `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	DeviceDescription                 []string  `tiledb:"dtype=string,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeResonTSeriesSonicSpecific(reader *bytes.Reader) (sensor_data ResonTSeries, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("ResonTSeries sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.ProtocolVersion = []uint16{buffer.ProtocolVersion}
	sensor_data.DeviceID = []uint32{buffer.DeviceID}
	sensor_data.NumberDevices = []uint32{buffer.NumberDevices}
	sensor_data.SystemEnumerator = []uint16{buffer.SystemEnumerator}
	sensor_data.MajorSerialNumber = []uint32{buffer.MajorSerialNumber}
	sensor_data.MinorSerialNumber = []uint32{buffer.MinorSerialNumber}
	sensor_data.PingNumber = []uint32{buffer.PingNumber}
	sensor_data.MultiPingSequence = []uint16{buffer.MultiPingSequence}
	sensor_data.Frequency = []float64{float64(buffer.Frequency) / SCALE_3_F64}
	sensor_data.SampleRate = []float64{float64(buffer.SampleRate) / SCALE_4_F64}
	sensor_data.ReceiverBandwidth = []float64{float64(buffer.ReceiverBandwidth) / SCALE_4_F64}
	sensor_data.TxPulseWidth = []float64{float64(buffer.TxPulseWidth) / SCALE_7_F64}
	sensor_data.TxPulseTypeID = []uint32{buffer.TxPulseTypeID}
	sensor_data.TxPulseEnvlpID = []uint32{buffer.TxPulseEnvlpID}
	sensor_data.TxPulseEnvlpParam = []float64{float64(buffer.TxPulseEnvlpParam) / SCALE_2_F64}
	sensor_data.TxPulseMode = []uint16{buffer.TxPulseMode}
	sensor_data.MaxPingRate = []float64{float64(buffer.MaxPingRate) / SCALE_6_F64}
	sensor_data.PingPeriod = []float64{float64(buffer.PingPeriod) / SCALE_6_F64}
	sensor_data.Range = []float64{float64(buffer.Range) / SCALE_2_F64}
	sensor_data.Power = []float64{float64(buffer.Power) / SCALE_2_F64}
	sensor_data.Gain = []float64{float64(buffer.Gain) / SCALE_2_F64}
	sensor_data.ControlFlags = []uint32{buffer.ControlFlags}
	sensor_data.ProjectorID = []uint32{buffer.ProjectorID}
	sensor_data.ProjectorSteerAnglVert = []float64{float64(buffer.ProjectorSteerAnglVert) / SCALE_3_F64}
	sensor_data.ProjectorSteerAnglHorz = []float64{float64(buffer.ProjectorSteerAnglHorz) / SCALE_3_F64}
	sensor_data.ProjectorBeamWidthVert = []float32{float32(buffer.ProjectorBeamWidthVert) / SCALE_2_F32}
	sensor_data.ProjectorBeamWidthHorz = []float32{float32(buffer.ProjectorBeamWidthHorz) / SCALE_2_F32}
	sensor_data.ProjectorBeamFocalPt = []float64{float64(buffer.ProjectorBeamFocalPt) / SCALE_2_F64}
	sensor_data.ProjectorBeamWeightingWindowType = []uint32{buffer.ProjectorBeamWeightingWindowType}
	sensor_data.ProjectorBeamWeightingWindowParam = []uint32{buffer.ProjectorBeamWeightingWindowParam}
	sensor_data.TransmitFlags = []uint32{buffer.TransmitFlags}
	sensor_data.HydrophoneID = []uint32{buffer.HydrophoneID}
	sensor_data.ReceivingBeamWeightingWindowType = []uint32{buffer.ReceivingBeamWeightingWindowType}
	sensor_data.ReceivingBeamWeightingWindowParam = []uint32{buffer.ReceivingBeamWeightingWindowParam}
	sensor_data.ReceiveFlags = []uint32{buffer.ReceiveFlags}
	sensor_data.ReceiveBeamWidth = []float64{float64(buffer.ReceiveBeamWidth) / SCALE_2_F64}
	sensor_data.RangeFiltMin = []float64{float64(buffer.RangeFiltMin) / SCALE_1_F64}
	sensor_data.RangeFiltMax = []float64{float64(buffer.RangeFiltMax) / SCALE_1_F64}
	sensor_data.DepthFiltMin = []float64{float64(buffer.DepthFiltMin) / SCALE_1_F64}
	sensor_data.DepthFiltMax = []float64{float64(buffer.DepthFiltMax) / SCALE_1_F64}
	sensor_data.Absorption = []float64{float64(buffer.Absorption) / SCALE_3_F64}
	sensor_data.SoundVelocity = []float64{float64(buffer.SoundVelocity) / SCALE_1_F64}
	sensor_data.SvSource = []uint8{buffer.SvSource}
	sensor_data.Spreading = []float64{float64(buffer.Spreading) / SCALE_3_F64}
	sensor_data.BeamSpacingMode = []uint16{buffer.BeamSpacingMode}
	sensor_data.SonarSourceMode = []uint16{buffer.SonarSourceMode}
	sensor_data.CoverageMode = []uint8{buffer.CoverageMode}
	sensor_data.CoverageAngle = []float64{float64(buffer.CoverageAngle) / SCALE_2_F64}
	sensor_data.HorizontalReceiverSteeringAngle = []float64{float64(buffer.HorizontalReceiverSteeringAngle) / SCALE_2_F64}
	sensor_data.UncertaintyType = []uint32{buffer.UncertaintyType}
	sensor_data.TransmitterSteeringAngle = []float64{float64(buffer.TransmitterSteeringAngle) / SCALE_5_F64}
	sensor_data.AppliedRoll = []float64{float64(buffer.AppliedRoll) / SCALE_5_F64}
	sensor_data.DetectionAlgorithm = []uint16{buffer.DetectionAlgorithm}
	sensor_data.DetectionFlags = []uint32{buffer.DetectionFlags}
	sensor_data.DeviceDescription = []string{string(buffer.DeviceDescription[:])}

	// higher precision sound velocity
	if buffer.SoundVelocity2 > 0 {
		sensor_data.SoundVelocity = []float64{float64(buffer.SoundVelocity2) / SCALE_6_F64}
	}

	return sensor_data, err
}

type Kmall struct {
	KmallVersion                       []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	DgmType                            []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	DgmVersion                         []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SystemID                           []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	EchoSounderID                      []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumBytesCmnPart                    []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingCounter                        []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	RxFansPerRing                      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RxFansIndex                        []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SwathsPerRing                      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SwathAlongPosition                 []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TxTransducerIndex                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RxTransducerIndex                  []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	NumRxTransducers                   []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	AlgorithmType                      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	NumBytesInfoData                   []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PingRateHz                         []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	BeamSpacing                        []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	DepthMode                          []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SubDepthMode                       []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	DistanceBetweenSwath               []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	DetectionMode                      []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PulseForm                          []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	FrequencyModeHz                    []int32     `tiledb:"dtype=int32,ftype=attr" filters:"zstd(level=16)"`
	FrequencyRangeLowLimHz             []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	FrequencyRangeHighLimHz            []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	MaxTotalTxPulseLengthSector        []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	MaxEffectiveTxPulseLengthSector    []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	MaxEffectiveTxBandWidthHz          []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	AbsCoeffDbPerKm                    []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	PortSectorEdgeDeg                  []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	StarboardSectorEdgeDeg             []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	PortMeanCoverageDeg                []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	StarboardMeanCoverageDeg           []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	PortMeanCoverageMetres             []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	StarboardMeanCoverageMetres        []int16     `tiledb:"dtype=int16,ftype=attr" filters:"zstd(level=16)"`
	ModeAndStabilisation               []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeFilter1                     []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	RunTimeFilter2                     []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PipeTrackingStatus                 []uint32    `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	TransmitArraySizeUsedDeg           []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	ReceiveArraySizeUsedDeg            []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransmitPowerDb                    []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SlRampUpTimeRemaining              []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	YawAngleDeg                        []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	NumTxSectors                       []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumBytesPerTxSector                []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	HeadingVesselDeg                   []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SoundSpeedAtTxDepthMetresPerSecond []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxTransducerDepthMetres            []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ZwaterLevelReRefPointMetres        []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	XKmallToAllMetres                  []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	YKmallToAllMetres                  []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	LatLonInfo                         []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PositionSensorStatus               []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	AttitudeSensorStatus               []uint8     `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	LatitudeDeg                        []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	LongitudeDeg                       []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	EllipsoidHeightReRefPointMetres    []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxSectorNumber                     [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	TxArrayNumber                      [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	TxSubArray                         [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	SectorTransmitDelaySec             [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	TiltAngleReTxDeg                   [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	TxNominalSourceLevelDb             [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	TxFocusRangeMetres                 [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	CentreFrequencyHz                  [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	SignalBandWidthHz                  [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	TotalSignalLengthSec               [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	PulseShading                       [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	SignalWaveForm                     [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
	NumBytesRxInfo                     []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumSoundingsMaxMain                []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumSoundingsValidMain              []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumBytesPerSounding                []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	WcSampleRate                       []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SeabedImageSampleRate              []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	BackscatterNormalDb                []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	BackscatterObliqueDb               []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	ExtraDetectionAlarmFlag            []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumExtraDetections                 []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumExtraDetectionClasses           []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumBytesPerClass                   []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	NumExtraDetectionInClass           [][]uint16  `tiledb:"dtype=uint16,ftype=attr,var" filters:"zstd(level=16)"`
	AlarmFlag                          [][]uint8   `tiledb:"dtype=uint8,ftype=attr,var" filters:"zstd(level=16)"`
}

func DecodeKmallSpecific(reader *bytes.Reader) (sensor_data Kmall, err error) {
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
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("KMALL sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
	sensor_data.KmallVersion = []uint8{buffer.KmallVersion}
	sensor_data.DgmType = []uint8{buffer.DgmType}
	sensor_data.DgmVersion = []uint8{buffer.DgmVersion}
	sensor_data.EchoSounderID = []uint16{buffer.EchoSounderID}

	// block two (Cmn part)
	err = binary.Read(reader, binary.BigEndian, &cmn_buf)
	if err != nil {
		errn := errors.New("KMALL sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
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
	err = binary.Read(reader, binary.BigEndian, &cmn_buf)
	if err != nil {
		errn := errors.New("KMALL sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
	sensor_data.NumBytesInfoData = []uint16{ping_buf.NumBytesInfoData}
	sensor_data.PingRateHz = []float64{float64(ping_buf.PingRateHz) / SCALE_5_F64}
	sensor_data.BeamSpacing = []uint8{ping_buf.BeamSpacing}
	sensor_data.DepthMode = []uint8{ping_buf.DepthMode}
	sensor_data.SubDepthMode = []uint8{ping_buf.SubDepthMode}
	sensor_data.DistanceBetweenSwath = []uint8{ping_buf.DistanceBetweenSwath}
	sensor_data.DetectionMode = []uint8{ping_buf.DetectionMode}
	sensor_data.PulseForm = []uint8{ping_buf.PulseForm}
	sensor_data.FrequencyModeHz = []int32{ping_buf.FrequencyModeHz}
	sensor_data.FrequencyRangeLowLimHz = []float64{float64(ping_buf.FrequencyRangeLowLimHz) / SCALE_3_F64}
	sensor_data.FrequencyRangeHighLimHz = []float64{float64(ping_buf.FrequencyRangeHighLimHz) / SCALE_3_F64}
	sensor_data.MaxTotalTxPulseLengthSector = []float64{float64(ping_buf.MaxTotalTxPulseLengthSector) / SCALE_6_F64}
	sensor_data.MaxEffectiveTxPulseLengthSector = []float64{float64(ping_buf.MaxEffectiveTxPulseLengthSector) / SCALE_6_F64}
	sensor_data.MaxEffectiveTxBandWidthHz = []float64{float64(ping_buf.MaxEffectiveTxBandWidthHz) / SCALE_3_F64}
	sensor_data.AbsCoeffDbPerKm = []float64{float64(ping_buf.AbsCoeffDbPerKm) / SCALE_3_F64}
	sensor_data.PortSectorEdgeDeg = []float32{float32(ping_buf.PortSectorEdgeDeg) / SCALE_2_F32}
	sensor_data.StarboardSectorEdgeDeg = []float32{float32(ping_buf.StarboardSectorEdgeDeg) / SCALE_2_F32}
	sensor_data.PortMeanCoverageDeg = []float32{float32(ping_buf.PortMeanCoverageDeg) / SCALE_2_F32}
	sensor_data.StarboardMeanCoverageDeg = []float32{float32(ping_buf.StarboardMeanCoverageDeg) / SCALE_2_F32}
	sensor_data.PortMeanCoverageMetres = []int16{ping_buf.PortMeanCoverageMetres}
	sensor_data.StarboardMeanCoverageMetres = []int16{ping_buf.StarboardMeanCoverageMetres}
	sensor_data.ModeAndStabilisation = []uint8{ping_buf.ModeAndStabilisation}
	sensor_data.RunTimeFilter1 = []uint8{ping_buf.RunTimeFilter1}
	sensor_data.RunTimeFilter2 = []uint8{ping_buf.RunTimeFilter2}
	sensor_data.PipeTrackingStatus = []uint32{ping_buf.PipeTrackingStatus}
	sensor_data.TransmitArraySizeUsedDeg = []float32{float32(ping_buf.TransmitArraySizeUsedDeg) / SCALE_3_F32}
	sensor_data.ReceiveArraySizeUsedDeg = []float32{float32(ping_buf.ReceiveArraySizeUsedDeg) / SCALE_3_F32}
	sensor_data.TransmitPowerDb = []float32{float32(ping_buf.TransmitPowerDb) / SCALE_2_F32}
	sensor_data.SlRampUpTimeRemaining = []uint16{ping_buf.SlRampUpTimeRemaining}
	sensor_data.YawAngleDeg = []float64{float64(ping_buf.YawAngleDeg) / SCALE_6_F64}
	sensor_data.NumTxSectors = []uint16{ping_buf.NumTxSectors}
	sensor_data.NumBytesPerTxSector = []uint16{ping_buf.NumBytesPerTxSector}
	sensor_data.HeadingVesselDeg = []float64{float64(ping_buf.HeadingVesselDeg) / SCALE_6_F64}
	sensor_data.SoundSpeedAtTxDepthMetresPerSecond = []float64{float64(ping_buf.SoundSpeedAtTxDepthMetresPerSecond) / SCALE_6_F64}
	sensor_data.TxTransducerDepthMetres = []float64{float64(ping_buf.TxTransducerDepthMetres) / SCALE_6_F64}
	sensor_data.ZwaterLevelReRefPointMetres = []float64{float64(ping_buf.ZwaterLevelReRefPointMetres) / SCALE_6_F64}
	sensor_data.XKmallToAllMetres = []float64{float64(ping_buf.XKmallToAllMetres) / SCALE_6_F64}
	sensor_data.YKmallToAllMetres = []float64{float64(ping_buf.YKmallToAllMetres) / SCALE_6_F64}
	sensor_data.LatLonInfo = []uint8{ping_buf.LatLonInfo}
	sensor_data.PositionSensorStatus = []uint8{ping_buf.PositionSensorStatus}
	sensor_data.AttitudeSensorStatus = []uint8{ping_buf.AttitudeSensorStatus}
	sensor_data.LatitudeDeg = []float64{float64(ping_buf.LatitudeDeg) / SCALE_7_F64}
	sensor_data.LongitudeDeg = []float64{float64(ping_buf.LongitudeDeg) / SCALE_7_F64}
	sensor_data.EllipsoidHeightReRefPointMetres = []float64{float64(ping_buf.EllipsoidHeightReRefPointMetres) / SCALE_3_F64}

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
		err = binary.Read(reader, binary.BigEndian, &sec_buf)
		if err != nil {
			errn := errors.New("KMALL sensor")
			err = errors.Join(err, ErrSensorMetadata, errn)
			return sensor_data, err
		}
		TxSectorNumber = append(TxSectorNumber, sec_buf.TxSectorNumber)
		TxArrayNumber = append(TxArrayNumber, sec_buf.TxArrayNumber)
		TxSubArray = append(TxSubArray, sec_buf.TxSubArray)
		SectorTransmitDelaySec = append(SectorTransmitDelaySec, float64(sec_buf.SectorTransmitDelaySec)/SCALE_6_F64)
		TiltAngleReTxDeg = append(TiltAngleReTxDeg, float64(sec_buf.TiltAngleReTxDeg)/SCALE_6_F64)
		TxNominalSourceLevelDb = append(TxNominalSourceLevelDb, float64(sec_buf.TxNominalSourceLevelDb)/SCALE_6_F64)
		TxFocusRangeMetres = append(TxFocusRangeMetres, float64(sec_buf.TxFocusRangeMetres)/SCALE_3_F64)
		CentreFrequencyHz = append(CentreFrequencyHz, float64(sec_buf.CentreFrequencyHz)/SCALE_3_F64)
		SignalBandWidthHz = append(SignalBandWidthHz, float64(sec_buf.SignalBandWidthHz)/SCALE_3_F64)
		TotalSignalLengthSec = append(TotalSignalLengthSec, float64(sec_buf.TotalSignalLengthSec)/SCALE_6_F64)
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
	err = binary.Read(reader, binary.BigEndian, &rx_buf)
	if err != nil {
		errn := errors.New("KMALL sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}
	sensor_data.NumBytesRxInfo = []uint16{rx_buf.NumBytesRxInfo}
	sensor_data.NumSoundingsMaxMain = []uint16{rx_buf.NumSoundingsMaxMain}
	sensor_data.NumSoundingsValidMain = []uint16{rx_buf.NumSoundingsValidMain}
	sensor_data.NumBytesPerSounding = []uint16{rx_buf.NumBytesPerSounding}
	sensor_data.WcSampleRate = []float64{float64(rx_buf.WcSampleRate1) + (float64(rx_buf.WcSampleRate2) / SCALE_9_F64)}
	sensor_data.SeabedImageSampleRate = []float64{float64(rx_buf.SeabedImageSampleRate1) + (float64(rx_buf.SeabedImageSampleRate2) / SCALE_9_F64)}
	sensor_data.BackscatterNormalDb = []float64{float64(rx_buf.BackscatterNormalDb) / SCALE_6_F64}
	sensor_data.BackscatterObliqueDb = []float64{float64(rx_buf.BackscatterObliqueDb) / SCALE_6_F64}
	sensor_data.ExtraDetectionAlarmFlag = []uint16{rx_buf.ExtraDetectionAlarmFlag}
	sensor_data.NumExtraDetections = []uint16{rx_buf.NumExtraDetections}
	sensor_data.NumExtraDetectionClasses = []uint16{rx_buf.NumExtraDetectionClasses}
	sensor_data.NumBytesPerClass = []uint16{rx_buf.NumBytesPerClass}

	// block six (extra detection classes
	nclasses := int(rx_buf.NumExtraDetectionClasses)
	NumExtraDetectionInClass := make([]uint16, 0, nclasses)
	AlarmFlag := make([]uint8, 0, nclasses)

	for i := 0; i < nclasses; i++ {
		err = binary.Read(reader, binary.BigEndian, &cls_buf)
		if err != nil {
			errn := errors.New("KMALL sensor")
			err = errors.Join(err, ErrSensorMetadata, errn)
			return sensor_data, err
		}
		NumExtraDetectionInClass = append(NumExtraDetectionInClass, cls_buf.NumExtraDetectionInClass)
		AlarmFlag = append(AlarmFlag, cls_buf.AlarmFlag)
	}

	sensor_data.NumExtraDetectionInClass = [][]uint16{NumExtraDetectionInClass}
	sensor_data.AlarmFlag = [][]uint8{AlarmFlag}

	err = binary.Read(reader, binary.BigEndian, &final_spare)
	if err != nil {
		errn := errors.New("KMALL sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	return sensor_data, err
}

// Single beam types (specified as swath as they're specific to reading from the SWATH_BATHYMETRY_PING record

type SwathSbEchotrac struct {
	NavigationError []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	MppSource       []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TideSource      []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	DynamicDraft    []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSwathSbEchotracSpecific(reader *bytes.Reader) (sensor_data SwathSbEchotrac, err error) {
	var buffer struct {
		NavigationError uint16
		MppSource       uint8
		TideSource      uint8
		DynamicDraft    int16
		Spare           [4]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Swath SbEchotrac or Swath Bathy2000 or Swath PDD sensors")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.NavigationError = []uint16{buffer.NavigationError}
	sensor_data.MppSource = []uint8{buffer.MppSource}
	sensor_data.TideSource = []uint8{buffer.TideSource}
	sensor_data.DynamicDraft = []float32{float32(buffer.DynamicDraft) / SCALE_2_F32}

	return sensor_data, err
}

type SwathSbMgd77 struct {
	TimeZoneCorrection []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PositionTypeCode   []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	CorrectionCode     []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	BathyTypeCode      []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	QualityCode        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TravelTime         []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSwathSbMGD77Specific(reader *bytes.Reader) (sensor_data SwathSbMgd77, err error) {
	var buffer struct {
		TimeZoneCorrection uint16
		PositionTypeCode   uint16
		CorrectionCode     uint16
		BathyTypeCode      uint16
		QualityCode        uint16
		TravelTime         uint32
		Spare              [4]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Swath SBMGD77 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.TimeZoneCorrection = []uint16{buffer.TimeZoneCorrection}
	sensor_data.PositionTypeCode = []uint16{buffer.PositionTypeCode}
	sensor_data.CorrectionCode = []uint16{buffer.CorrectionCode}
	sensor_data.BathyTypeCode = []uint16{buffer.BathyTypeCode}
	sensor_data.QualityCode = []uint16{buffer.QualityCode}
	sensor_data.TravelTime = []float64{float64(buffer.TravelTime) / SCALE_4_F64}

	return sensor_data, err
}

type SwathSbBdb struct {
	TravelTime           []uint32 `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	EvaluationFlag       []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ClassificationFlag   []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TrackAdjustmentFlag  []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SourceFlag           []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PointOrTrackLineFlag []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	DatumFlag            []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSwathSbBdbSpecific(reader *bytes.Reader) (sensor_data SwathSbBdb, err error) {
	var buffer struct {
		TravelTime           uint32
		EvaluationFlag       uint8
		ClassificationFlag   uint8
		TrackAdjustmentFlag  uint8
		SourceFlag           uint8
		PointOrTrackLineFlag uint8
		DatumFlag            uint8
		Spare                [4]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Swath SBBDB sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.TravelTime = []uint32{buffer.TravelTime}
	sensor_data.EvaluationFlag = []uint8{buffer.EvaluationFlag}
	sensor_data.ClassificationFlag = []uint8{buffer.ClassificationFlag}
	sensor_data.TrackAdjustmentFlag = []uint8{buffer.TrackAdjustmentFlag}
	sensor_data.SourceFlag = []uint8{buffer.SourceFlag}
	sensor_data.PointOrTrackLineFlag = []uint8{buffer.PointOrTrackLineFlag}
	sensor_data.DatumFlag = []uint8{buffer.DatumFlag}

	return sensor_data, err
}

type SwathSbNoShDb struct {
	TypeCode         []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	CartographicCode []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSwathSbNoShDbSpecific(reader *bytes.Reader) (sensor_data SwathSbNoShDb, err error) {
	var buffer struct {
		TypeCode         uint16
		CartographicCode uint16
		Spare            [4]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Swath SBNOSHDB sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.TypeCode = []uint16{buffer.TypeCode}
	sensor_data.CartographicCode = []uint16{buffer.CartographicCode}

	return sensor_data, err
}

type SwathSbNavisound struct {
	PulseLength []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSwathSbNavisoundSpecific(reader *bytes.Reader) (sensor_data SwathSbNavisound, err error) {
	var buffer struct {
		PulseLength uint16
		Spare       [8]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Swath SBNavisound sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.PulseLength = []float32{float32(buffer.PulseLength) / SCALE_2_F32}

	return sensor_data, err
}

// single beam types, that are intended for decoding via the SINGLE_BEAM_PING record
// Some of the types are, for some reason, not explicitly the same as the swath equivalent.
// Others are a duplicate. But for the time being, and consistency, separate types will be
// defined. At a future date, duplicates types and associated decoders may be considered for removal.

type SbEchotrac struct {
	NavigationError []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	MppSource       []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TideSource      []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSbEchotracSpecific(reader *bytes.Reader) (sensor_data SbEchotrac, err error) {
	var buffer struct {
		NavigationError uint16
		MppSource       uint8
		TideSource      uint8
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("SbEchotrac or Bathy2000 sensors")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.NavigationError = []uint16{buffer.NavigationError}
	sensor_data.MppSource = []uint8{buffer.MppSource}
	sensor_data.TideSource = []uint8{buffer.TideSource}

	return sensor_data, err
}

type SbMgd77 struct {
	TimeZoneCorrection []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	PositionTypeCode   []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	CorrectionCode     []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	BathyTypeCode      []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	QualityCode        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TravelTime         []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSbMGD77Specific(reader *bytes.Reader) (sensor_data SbMgd77, err error) {
	var buffer struct {
		TimeZoneCorrection uint16
		PositionTypeCode   uint16
		CorrectionCode     uint16
		BathyTypeCode      uint16
		QualityCode        uint16
		TravelTime         uint32
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("SBMGD77 sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.TimeZoneCorrection = []uint16{buffer.TimeZoneCorrection}
	sensor_data.PositionTypeCode = []uint16{buffer.PositionTypeCode}
	sensor_data.CorrectionCode = []uint16{buffer.CorrectionCode}
	sensor_data.BathyTypeCode = []uint16{buffer.BathyTypeCode}
	sensor_data.QualityCode = []uint16{buffer.QualityCode}
	sensor_data.TravelTime = []float64{float64(buffer.TravelTime) / SCALE_4_F64}

	return sensor_data, err
}

type SbBdb struct {
	TravelTime           []uint32 `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	EvaluationFlag       []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	ClassificationFlag   []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	TrackAdjustmentFlag  []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SourceFlag           []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	PointOrTrackLineFlag []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	DatumFlag            []uint8  `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSbBdbSpecific(reader *bytes.Reader) (sensor_data SbBdb, err error) {
	var buffer struct {
		TravelTime           uint32
		EvaluationFlag       uint8
		ClassificationFlag   uint8
		TrackAdjustmentFlag  uint8
		SourceFlag           uint8
		PointOrTrackLineFlag uint8
		DatumFlag            uint8
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("SBBDB sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.TravelTime = []uint32{buffer.TravelTime}
	sensor_data.EvaluationFlag = []uint8{buffer.EvaluationFlag}
	sensor_data.ClassificationFlag = []uint8{buffer.ClassificationFlag}
	sensor_data.TrackAdjustmentFlag = []uint8{buffer.TrackAdjustmentFlag}
	sensor_data.SourceFlag = []uint8{buffer.SourceFlag}
	sensor_data.PointOrTrackLineFlag = []uint8{buffer.PointOrTrackLineFlag}
	sensor_data.DatumFlag = []uint8{buffer.DatumFlag}

	return sensor_data, err
}

type SbNoShDb struct {
	TypeCode         []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	CartographicCode []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

func DecodeSbNoShDbSpecific(reader *bytes.Reader) (sensor_data SbNoShDb, err error) {
	var buffer struct {
		TypeCode         uint16
		CartographicCode uint16
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("SBNOSHDB sensor")
		err = errors.Join(err, ErrSensorMetadata, errn)
		return sensor_data, err
	}

	sensor_data.TypeCode = []uint16{buffer.TypeCode}
	sensor_data.CartographicCode = []uint16{buffer.CartographicCode}

	return sensor_data, err
}
