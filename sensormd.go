package gsf

import (
	"bytes"
	"encoding/binary"
	"time"
)

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
