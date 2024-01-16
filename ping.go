package gsf

import (
	// "os"
	"bytes"
	"encoding/binary"
	"time"
)

type PingHeader struct {
	Timestamp          time.Time
	Longitude          float64
	Latitude           float64
	Number_beams       uint16
	Centre_beam        uint16
	Tide_corrector     float32
	Depth_corrector    float32
	Heading            float32
	Pitch              float32
	Roll               float32
	Heave              float32
	Course             float32
	Speed              float32
	Height             float32
	Separation         float32
	GPS_tide_corrector float32
	Ping_flags         int16
}

type PingHeaders struct {
	Timestamp          []time.Time
	Longitude          []float64
	Latitude           []float64
	Number_beams       []uint16
	Centre_beam        []uint16
	Tide_corrector     []float32
	Depth_corrector    []float32
	Heading            []float32
	Pitch              []float32
	Roll               []float32
	Heave              []float32
	Course             []float32
	Speed              []float32
	Height             []float32
	Separation         []float32
	GPS_tide_corrector []float32
	Ping_flags         []int16
}

type ScaleFactor struct {
	Id               SubRecordID
	Scale            float32 // TODO float32?
	Offset           float32
	Compression_flag int
	Compressed       bool // if true, then the associated array is compressed
	Field_size       int
}

type BeamArray struct {
	Depth                []float32
	AcrossTrack          []float32
	AlongTrack           []float32
	TravelTime           []float32
	BeamAngle            []float32
	MeanCalAmplitude     []float32
	MeanRelAmplitude     []float32
	EchoWidth            []float32
	QualityFactor        []float32
	RecieveHeave         []float32
	DepthError           []float32 // obsolete
	AcrossTrackError     []float32 // obsolete
	AlongTrackError      []float32 // obsolete
	NominalDepth         []float32
	QualityFlags         []float32
	BeamFlags            []uint8
	SignalToNoise        []float32
	BeamAngleForward     []float32
	VerticalError        []float32
	HorizontalError      []float32
	IntensitySeries      []float32
	SectorNumber         []float32
	DetectionInfo        []float32
	IncidentBeamAdj      []float32
	SystemCleaning       []float32
	DopplerCorrection    []float32
	SonarVertUncertainty []float32
	SonarHorzUncertainty []float32
	DetectionWindow      []float32
	MeanAbsCoef          []float32
}

type PingGroup struct {
	Start         uint64
	Stop          uint64
	Number_Beams  uint64
	Scale_Factors map[SubRecordID]ScaleFactor
}

// PingInfo contains some basic information regarding the ping such as
// the number of beams, what sub-records are populated.
// The initial reasoning behind why, is to provide a basic descriptor
// to inform a global schema across all pings, and derive max(n_beams) to
// inform a global [ping, beam] dimensional array structure.
type PingInfo struct {
	Timestamp     time.Time
	Number_Beams  uint16
	Sub_Records   []SubRecordID
	Scale_Factors bool
	scale_factors map[SubRecordID]ScaleFactor
}

type PingData struct {
	Ping_headers             []PingHeader
	Beam_array               BeamArray
	Brb_intensity            BrbIntensity
	Sensor_metadata          SensorMetadata
	Sensory_imagery_metadata SensorImageryMetadata
	Lon_lat                  LonLat
	n_pings                  uint64
}

// PingGroups combines pings together based on their presence or absence of
// scale factors. It is a forward linear search, and if a given ping is missing
// scale factors, then it is included as part of the ping group where the previous
// set of scale factors were found.
// For example; [0, 10] indicates that the ping group contains pings 0 up to and
// including ping 9. It is a [start, stop) index based on the linear ordering
// of pings found in the GSF file.
func (fi *FileInfo) PGroups() {
	var (
		start      int
		ping_group PingGroup
		groups     []PingGroup
		sf         map[SubRecordID]ScaleFactor
		beam_count uint64
	)

	groups = make([]PingGroup, 0)
	beam_count = uint64(0)

	for i, ping := range fi.Ping_Info {
		if ping.Scale_Factors {
			if i > 0 {
				// new group
				ping_group = PingGroup{
					uint64(start), uint64(i), beam_count, sf,
				}
				groups = append(groups, ping_group)
			}
			// update with latest sf dependency and reset counters
			start = i
			beam_count = uint64(0)
			sf = fi.Ping_Info[start].scale_factors
		} else {
			// set scale factors based on the last read scale factors
			fi.Ping_Info[i].scale_factors = sf
			beam_count += uint64(ping_group.Number_Beams)
		}
	}

	fi.Index.Ping_Groups = groups
}

func decode_ping_hdr(reader *bytes.Reader) PingHeader {
	var (
		hdr_base struct {
			Seconds            int32
			Nano_seconds       int32
			Longitude          int32
			Latitude           int32
			Number_beams       uint16
			Centre_beam        uint16
			Ping_flags         int16
			Reserved           int16
			Tide_corrector     int16
			Depth_corrector    int32
			Heading            uint16
			Pitch              int16
			Roll               int16
			Heave              int16
			Course             uint16
			Speed              uint16
			Height             int32
			Separation         int32
			GPS_tide_corrector int32
			Spare              int16
		}
		hdr PingHeader
	)

	_ = binary.Read(reader, binary.BigEndian, &hdr_base)

	hdr.Timestamp = time.Unix(int64(hdr_base.Seconds), int64(hdr_base.Nano_seconds)).UTC()
	hdr.Longitude = float64(float32(hdr_base.Longitude) / SCALE1)
	hdr.Latitude = float64(float32(hdr_base.Latitude) / SCALE1)
	hdr.Number_beams = hdr_base.Number_beams
	hdr.Centre_beam = hdr_base.Centre_beam
	hdr.Ping_flags = hdr_base.Ping_flags
	hdr.Tide_corrector = float32(hdr_base.Tide_corrector) / SCALE2
	hdr.Depth_corrector = float32(hdr_base.Depth_corrector) / SCALE2
	hdr.Heading = float32(hdr_base.Heading) / SCALE2
	hdr.Pitch = float32(hdr_base.Pitch) / SCALE2
	hdr.Roll = float32(hdr_base.Roll) / SCALE2
	hdr.Heave = float32(hdr_base.Heave) / SCALE2
	hdr.Course = float32(hdr_base.Course) / SCALE2
	hdr.Speed = float32(hdr_base.Speed) / SCALE2
	hdr.Height = float32(hdr_base.Height) / SCALE3
	hdr.Separation = float32(hdr_base.Separation) / SCALE3
	hdr.GPS_tide_corrector = float32(hdr_base.GPS_tide_corrector) / SCALE3

	return hdr
}

func SubRecHdr(reader *bytes.Reader, offset int64) SubRecord {
	var subrecord_hdr int32

	_ = binary.Read(reader, binary.BigEndian, &subrecord_hdr)

	subrecord_id := (int(subrecord_hdr) & 0xFF000000) >> 24 // TODO; define a const as int64
	subrecord_size := int(subrecord_hdr) & 0x00FFFFFF       // TODO; define a const as int64

	byte_index := offset + 4

	subhdr := SubRecord{SubRecordID(subrecord_id), uint32(subrecord_size), byte_index} // include a byte_index??

	return subhdr
}

func scale_factors_rec(reader *bytes.Reader) (scale_factors map[SubRecordID]ScaleFactor, nbytes int64) {
	var (
		i            int32
		num_factors  int32
		scale_factor ScaleFactor
		// scale_factors map[int32]scale_factor
	)
	data := make([]int32, 3) // id, scale, offset
	scale_factors = map[SubRecordID]ScaleFactor{}
	// scale_factors := make(map[SubRecordID]ScaleFactor)

	_ = binary.Read(reader, binary.BigEndian, &num_factors)
	nbytes = 4

	for i = 0; i < num_factors; i++ {
		_ = binary.Read(reader, binary.BigEndian, &data)

		subid := (int64(data[0]) & 0xFF000000) >> 24   // TODO; define const for 0xFF000000
		comp_flag := (int(data[0]) & 0x00FF0000) >> 16 // TODO; define const for 0x00FF0000
		comp := (comp_flag & 0x0F) == 1                // TODO; define const for 0x00FF0000
		cnvrt_subid := SubRecordID(subid)
		field_size := comp_flag & 0xF0

		scale_factor = ScaleFactor{
			Id:               cnvrt_subid,
			Scale:            float32(data[1]),
			Offset:           float32(data[2]),
			Compression_flag: comp_flag, // TODO; implement compression decoder
			Compressed:       comp,
			Field_size:       field_size, // this field doesn't appear to be used in the C code ???
		}

		nbytes += 12

		// scale_factors[SubRecordID(subid)] = scale_factor
		scale_factors[cnvrt_subid] = scale_factor
	}

	return scale_factors, nbytes
}

func ping_info(reader *bytes.Reader, rec RecordHdr) PingInfo {
	var (
		idx     int64 = 0
		pinfo   PingInfo
		records      = make([]SubRecordID, 0, 32)
		sf      bool = false
		scl_fac map[SubRecordID]ScaleFactor
		nbytes  int64
	)

	datasize := int64(rec.Datasize)

	hdr := decode_ping_hdr(reader)
	idx += 56 // 56 bytes read for ping header
	offset := rec.Byte_index + idx

	// read through each subrecord
	for (datasize - idx) > 4 {
		sub_rec := SubRecHdr(reader, offset)
		srec_dsize := int64(sub_rec.Datasize)
		idx += 4 // bytes read from header

		if sub_rec.Id == SCALE_FACTORS {
			sf = true
			scl_fac, nbytes = scale_factors_rec(reader)
			idx += nbytes
		} else {
			// prep for the next record
			_, _ = reader.Seek(srec_dsize, 1)
			idx += srec_dsize
		}

		records = append(records, sub_rec.Id)
	}

	pinfo.Timestamp = hdr.Timestamp
	pinfo.Number_Beams = hdr.Number_beams
	pinfo.Sub_Records = records[:]
	pinfo.Scale_Factors = sf

	if sf {
		pinfo.scale_factors = scl_fac
	}

	return pinfo
}

// Contains the main data of the acquisition such as depth, across track, along track.
// The header contains the time, position, attitude, heading, course, speed and the number
// of beams. The position in lon/lat for every beam needs to be calculated.
// This record also contains sub-records, such as scale factors, sensor specifics, and the
// beam data such as depth.
// In the sample data provided, there has been occurrences of inconsistencies between pings,
// for example sub-records containing MEAN_CAL_AMPLITUDE information in one ping but not
// another. Cases like that and bringing all pings into a single data structure requires
// missing data be filled with nulls, or drop fields/sub-records that aren't in every ping.
// In one case, there was an instance of inconsistency in the number of beams across pings.
// The case that occurred was something like ~90000 pings had 400 beams, and 1 ping had 399.
// Data providers had no idea how, but possibly a beam was removed manually from the file.
// Another instance was a duplicate ping. Same timestamp, location, depth, but zero values
// for supporting attributes/sub-records/fields (heading, course, +others). Again, this
// appeared to have never been encountered before (or never looked).
func SwathBathymetryPingRec(buffer []byte, rec RecordHdr, pinfo PingInfo, sensor_id SubRecordID) PingHeader {
	var (
		idx        int64 = 0
		beam_data  []float32
		ping_data  PingData
		beam_array BeamArray
		img_md     SensorImageryMetadata
		intensity  BrbIntensity
		sr_buff    []byte
		sr_reader  *bytes.Reader
		sen_md     SensorMetadata
		// nbytes    int64
		// sf map[SubRecordID]ScaleFactor
		// beams     BeamArray
		// subrecord_hdr int32
	)

	reader := bytes.NewReader(buffer)

	hdr := decode_ping_hdr(reader)
	idx += 56 // 56 bytes read for ping header
	offset := rec.Byte_index + idx

	// TODO; move loop to here
	for (int64(rec.Datasize) - idx) > 4 {

		// subrecord header
		// _, _ = reader.Seek(idx, 0) // shouldn't be needed
		sub_rec := SubRecHdr(reader, offset)
		idx += 4

		// read the whole subrecord and form a new reader
		// i think this is easier than passing around how many
		// bytes are read from each func associated with decoding a subrecord
		sr_buff = make([]byte, sub_rec.Datasize)
		_ = binary.Read(reader, binary.BigEndian, &sr_buff)
		sr_reader = bytes.NewReader(sr_buff)

		// offset is used to track the start of the subrecord from the start
		// of the file as given by Record.Byte_index
		// Incase we wish to serialise the subrecord info along with the record info
		offset += int64(sub_rec.Datasize)
		idx += int64(sub_rec.Datasize)

		// only relevant for the ping sub-record arrays and not the sensor specific
		bytes_per_beam := sub_rec.Datasize / uint32(pinfo.Number_Beams)

		// decode each sub-record beam array
		// also need to handle sensor specific records as well
		switch sub_rec.Id {

		// scale factors
		case SCALE_FACTORS:
			// read and structure the scale factors
			// however, we'll rely on the scale factors from PingInfo.scale_factors
			_, _ = scale_factors_rec(sr_reader)
			// idx += nbytes

		// beam array subrecords
		case DEPTH:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.Depth = beam_data
			// idx += nbytes
		case ACROSS_TRACK:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				true,
			)
			beam_array.AcrossTrack = beam_data
			// idx += nbytes
		case ALONG_TRACK:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				true,
			)
			beam_array.AlongTrack = beam_data
			// idx += nbytes
		case TRAVEL_TIME:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.TravelTime = beam_data
			// idx += nbytes
		case BEAM_ANGLE:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_TWO,
				true,
			)
			beam_array.BeamAngle = beam_data
			// idx += nbytes
		case MEAN_CAL_AMPLITUDE:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				true,
			)
			beam_array.MeanCalAmplitude = beam_data
			// idx += nbytes
		case MEAN_REL_AMPLITUDE:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.MeanRelAmplitude = beam_data
			// idx += nbytes
		case ECHO_WIDTH:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.EchoWidth = beam_data
			// idx += nbytes
		case QUALITY_FACTOR:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_ONE,
				false,
			)
			beam_array.QualityFactor = beam_data
			// idx += nbytes
		case RECEIVE_HEAVE:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_ONE,
				true,
			)
			beam_array.RecieveHeave = beam_data
			// idx += nbytes
		case DEPTH_ERROR:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_TWO,
				false,
			)
			beam_array.DepthError = beam_data
			// idx += nbytes
		case ACROSS_TRACK_ERROR:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_TWO,
				false,
			)
			beam_array.AcrossTrackError = beam_data
			// idx += nbytes
		case ALONG_TRACK_ERROR:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_TWO,
				false,
			)
			beam_array.AlongTrackError = beam_data
			// idx += nbytes
		case NOMINAL_DEPTH:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.NominalDepth = beam_data
			// idx += nbytes
		case QUALITY_FLAGS:
			// obselete
			// TODO; has specific decoder
			panic("QUALITY_FLAGS subrecord has been superceded")
		case BEAM_FLAGS:
			beam_array.BeamFlags = DecodeBeamFlagsArray(
				sr_reader,
				pinfo.Number_Beams,
			)
			// idx += nbytes
		case SIGNAL_TO_NOISE:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_ONE,
				true,
			)
			beam_array.SignalToNoise = beam_data
			// idx += nbytes
		case BEAM_ANGLE_FORWARD:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_TWO,
				false,
			)
			beam_array.BeamAngleForward = beam_data
			// idx += nbytes
		case VERTICAL_ERROR:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_TWO,
				false,
			)
			beam_array.VerticalError = beam_data
			// idx += nbytes
		case HORIZONTAL_ERROR:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_TWO,
				false,
			)
			beam_array.HorizontalError = beam_data
			// idx += nbytes
		case INTENSITY_SERIES:
			intensity, img_md = DecocdeBrbIntensity(sr_reader, pinfo.Number_Beams, sensor_id)
			// idx += nbytes
		case SECTOR_NUMBER:
			// should be fine to just use DecodeSubRecArray and specify
			// 1-byte per beam
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_ONE,
				false,
			)
			beam_array.SectorNumber = beam_data
			// idx += nbytes
		case DETECTION_INFO:
			// should be fine to just use DecodeSubRecArray and specify
			// 1-byte per beam
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_ONE,
				false,
			)
			beam_array.DetectionInfo = beam_data
			// idx += nbytes
		case INCIDENT_BEAM_ADJ:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_ONE,
				true,
			)
			beam_array.IncidentBeamAdj = beam_data
			// idx += nbytes
		case SYSTEM_CLEANING:
			// should be fine to just use DecodeSubRecArray and specify
			// 1-byte per beam
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_ONE,
				false,
			)
			beam_array.SystemCleaning = beam_data
			// idx += nbytes
		case DOPPLER_CORRECTION:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_ONE,
				true,
			)
			beam_array.DopplerCorrection = beam_data
			// idx += nbytes
		case SONAR_VERT_UNCERTAINTY:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_TWO,
				false,
			)
			beam_array.SonarVertUncertainty = beam_data
			// idx += nbytes
		case SONAR_HORZ_UNCERTAINTY:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				FIELD_SIZE_TWO,
				false,
			)
			beam_array.SonarHorzUncertainty = beam_data
			// idx += nbytes
		case DETECTION_WINDOW:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.DetectionWindow = beam_data
			// idx += nbytes
		case MEAN_ABS_COEF:
			beam_data = sub_rec.DecodeSubRecArray(
				sr_reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.MeanAbsCoef = beam_data
			// idx += nbytes

		// sensor specific subrecords
		case SEABEAM:
			// DecodeSeabeam
		case EM12:
			// DecodeEM12
		case EM100:
			// DecodeEM100
		case EM950:
			// DecodeEM950
		case EM121A:
			// DecodeEM121A
		case EM121:
			// DecodeEM121
		case SASS: // obsolete
			// DecodeSASS
		case SEAMAP:
			// DecodeSeaMap
		case SEABAT:
			// DecodeSeaBat
		case EM1000:
			// DecodeEM1000
		case TYPEIII_SEABEAM: // obsolete
			// DecodeTypeIII
		case SB_AMP:
			// DecodeSBAmp
		case SEABAT_II:
			// DecodeSeaBatII
		case SEABAT_8101:
			// DecodeSeaBat8101
		case SEABEAM_2112:
			// DecodeSeaBeam2112
		case ELAC_MKII:
			// DecodeElacMkII
		case CMP_SAAS: // CMP (compressed), should be used in place of SASS
			// DecodeCmpSass
		case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
			// DecodeReson8100
		case EM120, EM300, EM1002, EM2000, EM3000, EM3002, EM3000D, EM3002D, EM121A_SIS:
			// DecodeEM3
		case EM710, EM302, EM122, EM2040:
			// DecodeEM4
			sen_md.EM_4 = DecodeEM4Specific(sr_reader)
			ping_data.Sensor_metadata.EM_4 = DecodeEM4Specific(sr_reader)
			// idx += nbytes
		case GEOSWATH_PLUS:
			// DecodeGeoSwathPlus
		case KLEIN_5410_BSS:
			// DecodeKlein5410Bss
		case RESON_7125:
			// DecodeReson7100
		case EM300_RAW, EM1002_RAW, EM2000_RAW, EM3000_RAW, EM120_RAW, EM3002_RAW, EM3000D_RAW, EM3002D_RAW, EM121A_SIS_RAW:
			// DecodeEM3Raw
		case DELTA_T:
			// DecodeDeltaT
		case R2SONIC_2022, R2SONIC_2024, R2SONIC_2020:
			// DecodeR2Sonic
		case SR_NOT_DEFINED: // the spec makes no mention of ID 154
			panic("Subrecord ID 154 is not defined.")
		case RESON_TSERIES:
			// DecodeResonTSeries
		case KMALL:
			// DecodeKMALL

			// single beam swath sensor specific subrecords
		case SWATH_ECHOTRAC, SWATH_BATHY2000, SWATH_PDD:
			// DecodeSBEchotrac
		case SWATH_MGD77:
			// DecodeSBMGD77
		case SWATH_BDB:
			// DecodeSBBDB
		case SWATH_NOSHDB:
			// DecodeSBNOSHDB
		case SWATH_NAVISOUND:
			// DecodeSBNavisound
		}
	}

	geocoef := NewCoefWgs84()
	lonlat := beam_array.BeamsLonLat(hdr.Longitude, hdr.Latitude, hdr.Heading, geocoef)

	ping_data.Ping_headers = []PingHeader{hdr}
	ping_data.Beam_array = beam_array
	ping_data.Brb_intensity = intensity
	ping_data.Sensory_imagery_metadata = img_md
	ping_data.Sensor_metadata = sen_md
	ping_data.Lon_lat = lonlat
	ping_data.n_pings = uint64(1)

	return hdr
}

func (g *GsfFile) SwathBathymetryPingRecords() (ping_data PingData) {
	return ping_data
}
