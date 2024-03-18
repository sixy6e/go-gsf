package gsf

import (
	"bytes"
	"encoding/binary"
	"time"
)

// float64 may be overkill
// where scale factors are applied, float32 is used
// where it's confident float32 is enough to represent the value
// TODO; look into defining tiledb compression struct tags
// TODO; align into 64bit chunks
// the spec says binary integers are stored as either 1-byte unsigned, 2-byte signed or unsigned, or 4-byte signed
type EM4 struct {
	ModelNumber                       []int16
	PingCounter                       []int16
	SerialNumber                      []int16
	SurfaceVelocity                   []float32
	TransducerDepth                   []float64
	ValidDetections                   []int16
	SamplingFrequency                 []float64
	DopplerCorrectionScale            []int32
	VehicleDepth                      []float32
	TransmitSectors                   []int16
	TiltAngle                         [][]float32
	FocusRange                        [][]float32
	SignalLength                      [][]float64
	TransmitDelay                     [][]float64
	CenterFrequency                   [][]float64
	MeanAbsorption                    [][]float32
	WaveformId                        [][]uint8
	SectorNumber                      [][]uint8
	SignalBandwith                    [][]float64
	RunTimeModelNumber                []int16
	RunTimeDatagramTime               []time.Time
	RunTimePingCounter                []int16
	RunTimeSerialNumber               []int16
	RunTimeOperatorStationStatus      []uint8
	RunTimeProcessingUnitStatus       []uint8
	RunTimeBspStatus                  []uint8
	RunTimeHeadTransceiverStatus      []uint8
	RunTimeMode                       []uint8
	RunTimeFilterId                   []uint8
	RunTimeMinDepth                   []float32
	RunTimeMaxDepth                   []float32
	RunTimeAbsorption                 []float32
	RunTimeTransmitPulseLength        []float32
	RunTimeTransmitBeamWidth          []float32
	RunTimeTransmitPowerReduction     []uint8
	RunTimeReceiveBeamWidth           []float32
	RunTimeReceiveBandwidth           []float32
	RunTimeReceiveFixedGain           []uint8
	RunTimeTvgCrossOverAngle          []uint8
	RunTimeSsvSource                  []uint8
	RunTimeMaxPortSwathWidth          []int16
	RunTimeBeamSpacing                []uint8
	RunTimeMaxPortCoverage            []uint8
	RunTimeStabilization              []uint8
	RunTimeMaxStbdCoverage            []uint8
	RunTimeMaxStdbSwathWidth          []int16
	RunTimeTransmitAlongTilt          []float32
	RunTimeFilterId2                  []uint8
	ProcessorUnitCpuLoad              []uint8
	ProcessorUnitSensorStatus         []uint16
	ProcessorUnitAchievedPortCoverage []uint8
	ProcessorUnitAchievedStbdCoverage []uint8
	ProcessorUnitYawStabilization     []float32
}

type EM4Imagery struct {
	SamplingFrequency   []float64
	MeanAbsorption      []float32
	TransmitPulseLength []float32
	RangeNorm           []uint16
	StartTvgRamp        []uint16
	StopTvgRamp         []uint16
	BackscatterN        []float32
	BackscatterO        []float32
	TransmitBeamWidth   []float32
	TvgCrossOver        []float32
}

// func (e EM4Imagery) Serialise() bool {
// 	return true
// }

func DecodeEM4Specific(reader *bytes.Reader) (sensor_data EM4) {

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

	// n_bytes = 0

	// first 46 bytes
	_ = binary.Read(reader, binary.BigEndian, &buffer)
	// n_bytes += 46

	// sector arrays
	// what if TransmitSectors == 0???
	for i := int16(0); i < buffer.TransmitSectors; i++ {
		_ = binary.Read(reader, binary.BigEndian, &sector_buffer_base)
		// n_bytes += 40
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
	// n_bytes += 16

	// next 63 bytes for the RunTime info
	_ = binary.Read(reader, binary.BigEndian, &runtime_buffer)
	// n_bytes += 63

	// next 23 bytes for the processing unit info
	_ = binary.Read(reader, binary.BigEndian, &proc_buffer)
	// n_bytes += 23

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

	return sensor_data // , n_bytes
}

func (g *GsfFile) EM4SpecificRecords(fi *FileInfo, start uint64, stop uint64) (sensor_data EM4) {

	// var (
	//     // nrecs uint64
	//     rec RecordHdr
	//     recs []RecordHdr
	// )

	// nrecs = fi.Metadata.Record_Counts[RecordNames[SWATH_BATHYMETRY_PING]]
	// recs = fi.Record_Index[RecordNames[SWATH_BATHYMETRY_PING]][start:stop]

	// retrieve and process each ping
	// for i := uint64(0); i < nrecs; i++ {
	//
	// }

	// for idx, rec := range recs {
	//
	// }

	return sensor_data
}

func DecodeEM4Imagery(reader *bytes.Reader) (em4_md EM4Imagery, scl_off ScaleOffset) {

	var (
		base struct {
			SamplingFrequency1  int32
			SamplingFrequency2  int32
			MeanAbsorption      uint16
			TransmitPulseLength uint16
			RangeNorm           uint16
			StartTvgRamp        uint16
			StopTvgRamp         uint16
			BackscatterN        int16
			BackscatterO        int16
			TransmitBeamWidth   uint16
			TvgCrossOver        uint16
			Offset              int16
			Scale               int16
			Spare               [5]uint32 // 20 bytes spare
		} // 50 bytes
	)
	// n_bytes = 0

	_ = binary.Read(reader, binary.BigEndian, &base)
	// n_bytes += 50

	em4_md.SamplingFrequency = []float64{
		float64(base.SamplingFrequency1) +
			float64(base.SamplingFrequency2)/
				float64(4_000_000_000)}
	em4_md.MeanAbsorption = []float32{float32(base.MeanAbsorption)}
	em4_md.TransmitPulseLength = []float32{float32(base.TransmitPulseLength)}
	em4_md.RangeNorm = []uint16{base.RangeNorm}
	em4_md.StartTvgRamp = []uint16{base.StartTvgRamp}
	em4_md.StopTvgRamp = []uint16{base.StopTvgRamp}
	em4_md.BackscatterN = []float32{float32(base.BackscatterN) / float32(10.0)}
	em4_md.BackscatterO = []float32{float32(base.BackscatterO) / float32(10.0)}
	em4_md.TransmitBeamWidth = []float32{float32(base.TransmitBeamWidth) / float32(10.0)}
	em4_md.TvgCrossOver = []float32{float32(base.TvgCrossOver) / float32(10.0)}

	scl_off = ScaleOffset{float32(base.Scale), float32(base.Offset)}

	return em4_md, scl_off // , n_bytes
}
