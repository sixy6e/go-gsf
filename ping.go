package gsf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	tiledb "github.com/TileDB-Inc/TileDB-Go"
	"github.com/samber/lo"
)

// PingHeader contains the base information recorded for every SWATH_BATHYMETRY_PING
// record.
type PingHeader struct {
	Timestamp          time.Time
	Longitude          float64
	Latitude           float64
	Number_beams       uint16
	Centre_beam        uint16
	Tide_corrector     float32
	Depth_corrector    float64
	Heading            float32
	Pitch              float32
	Roll               float32
	Heave              float32
	Course             float32
	Speed              float32
	Height             float64
	Separation         float64
	GPS_tide_corrector float64
	Ping_flags         uint16
}

// PingHeaders extends the PingHeader type by converting every attribute into
// a slice. Mostly used as a convienience for directly serialising the data to
// a TileDB array, rather than building the slice upon write from a []PingHeader.
type PingHeaders struct {
	Timestamp          []time.Time `tiledb:"dtype=datetime_ns,ftype=attr" filters:"zstd(level=16)"`
	Longitude          []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Latitude           []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Number_beams       []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Centre_beam        []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	Tide_corrector     []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Depth_corrector    []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Heading            []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Pitch              []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Roll               []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Heave              []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Course             []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Speed              []float32   `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	Height             []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Separation         []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	GPS_tide_corrector []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Ping_flags         []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
}

// newPingHeaders is a helper func for initialising PingHeaders where
// the it will contain slices initialised to the number of pings required.
// This func is only utilised when processing groups of pings to form a single
// cohesive block of data.
func newPingHeaders(number_pings int) (ping_headers PingHeaders) {
	ping_headers = PingHeaders{}
	chunkedStructSlices(&ping_headers, number_pings)

	return ping_headers
}

// ScaleOffset defines the scale and offset used in the decoding process
// that most of the BeamArray SubRecord types are stored as within the GSF file.
type ScaleOffset struct {
	Scale  float64
	Offset float64
}

// ScaleFactor contains the relevant information for a BeamArray SubRecord such
// as the SubRecordID and the ScaleOffset and whether or not it was compressed.
type ScaleFactor struct {
	Id SubRecordID
	ScaleOffset
	Compression_flag uint32
	Compressed       bool // if true, then the associated array is compressed
	Field_size       uint32
}

// BeamArray is the main type for holding all the BeamArray data contained within the SubRecords
// of the SWATH_BATHYMETRY_PING record.
// The data stored in this type have already been converted to their normal units, i.e.
// applied any relevant scale and offsets in the decoding process.
type BeamArray struct {
	Z                    []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	NominalDepth         []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	AcrossTrack          []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	AlongTrack           []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TravelTime           []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	BeamAngle            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	MeanCalAmplitude     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	MeanRelAmplitude     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	EchoWidth            []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	QualityFactor        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RecieveHeave         []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	DepthError           []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"` // obsolete
	AcrossTrackError     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"` // obsolete
	AlongTrackError      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"` // obsolete
	BeamFlags            []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	QualityFlags         []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	SignalToNoise        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	BeamAngleForward     []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	VerticalError        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	HorizontalError      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	IntensitySeries      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"` // TODO; check that this field can be removed
	SectorNumber         []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	DetectionInfo        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	IncidentBeamAdj      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SystemCleaning       []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	DopplerCorrection    []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SonarVertUncertainty []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	SonarHorzUncertainty []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	DetectionWindow      []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	MeanAbsCoef          []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TvgDb                []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
}

// newBeamArray is a helper func for initialising BeamArray where
// the specific sensor will contain slices initialised to the number of pings
// required.
// This func is only utilised when processing groups of pings to form a single
// cohesive block of data.
func newBeamArray(number_beams int, beam_names []string) (beam_array BeamArray) {
	beam_array = BeamArray{}
	chunkedBeamArray(&beam_array, number_beams, beam_names)

	return beam_array
}

// PingBeamNumbers contains the ping and beam number ids for a group of Pings decoded
// from the SWATH_BATHYMETRY_PING record.
// Essentially, slices of n length, where the ping number will be repeated n times
// for as many beams are in a given ping.
// Eg PingBeamNumbers{[0, 0, 0, 1, 1, 1, 2, 2, 2], [0, 1, 2, 0, 1, 2, 0, 1, 2]}
type PingBeamNumbers struct {
	PingNumber []uint64 `tiledb:"dtype=uint64,ftype=attr" filters:"zstd(level=16)"`
	BeamNumber []uint64 `tiledb:"dtype=uint64,ftype=attr" filters:"zstd(level=16)"`
}

// appendPingBeam appends the ping and beam numbers to for a given ping
// onto the slice containing a chunk of pings.
// The ping id will be the same for all beams.
// The GSF spec only indicates that beams ids start at the first array element
// and end at the last element in the array.
// The issue that arises is when the number of beams differ between pings,
// meaning that beam IDs are not likely to be consistent in terms of ordering
// across the pings over the whole file.
// We can only guarantee ping ids within a single file, assigning Ping-0 for
// the first ping record in the file, and Ping-Nminus1 for the last ping record
// in the file. Some GSF in the sensor metadata populated the global ping id
// (global in the sense for the entire survey), and some did not.
// Another annoying inconsistency, meaning we can't rely on global ordering of
// ping id's across all GSF files within a given survey.
func (pb *PingBeamNumbers) appendPingBeam(ping_id uint64, number_beams uint16) error {
	for i := uint16(0); i < number_beams; i++ {
		pb.PingNumber = append(pb.PingNumber, ping_id)
		pb.BeamNumber = append(pb.BeamNumber, uint64(i))
	}

	return nil
}

// padPingBeam pads the ping and beam number slices by a given size.
// This padding is only applied when the beam data array is to be output as
// a dense TileDB array.
// Why is padding required? Because the number of beams for each ping
// can differ from one another. Padding makes the BeamNumber dimensional axis
// consistent across all pings.
func (pb *PingBeamNumbers) padPingBeam(ping_id uint64, number_beams uint16, size uint16) error {
	n := number_beams + size
	for i := number_beams; i < n; i++ {
		pb.PingNumber = append(pb.PingNumber, ping_id)
		pb.BeamNumber = append(pb.BeamNumber, uint64(i))
	}

	return nil
}

// newPingBeamNumbers initialises the PingBeamNumbers type.
func newPingBeamNumbers(number_beams int) (ping_beam_ids PingBeamNumbers) {
	ping_id := make([]uint64, 0, number_beams)
	beam_id := make([]uint64, 0, number_beams)
	ping_beam_ids = PingBeamNumbers{ping_id, beam_id}

	return ping_beam_ids
}

// PingGroup defines a group of pings defined by Start and Stop ping IDs.
// It contains the total count of the number of beams across all pings which is
// used to setup the total length of the BeamArray slices. It also contains the
// the scale factors for the SubRecords associated with this group of Pings.
// Ping groups are defined by their commonality of shared scale factors.
// So across the entire GSF, the length of each Ping groups varies depending
// on which pings use shared scale factors.
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

// PingData is the main type for holding all information relevant to n pings worth
// of SWATH_BATHYMETRY_PING records.
type PingData struct {
	Ping_headers            PingHeaders
	Beam_array              BeamArray
	Brb_intensity           BrbIntensity
	Sensor_metadata         SensorMetadata
	Sensor_imagery_metadata SensorImageryMetadata
	Lon_lat                 LonLat
	n_pings                 uint64
	ba_subrecords           []string
}

// appendPingData is used when combining chunks of pings together into
// a single cohesive data block ready for writing to TileDB.
// As the schema can be inconsistent between pings, the global schema
// (defined by the whole GSF file), only the beam array records that
// have been read for the SWATH_BATHYMETRY_PING record will be appended.
// A separate method will need to be used to append null data for the ping
// missing required beam array records defined as a Set of all Sub_Records
// from all SWATH_BATHYMETRY_PING records.
func (pd *PingData) appendPingData(singlePing *PingData, contains_intensity bool, sensor_id SubRecordID, beam_names []string) error {
	// TODO; look into functionalising the appending mechanism where we use reflect
	// Ping_headers
	rf_pd := reflect.ValueOf(&pd.Ping_headers).Elem()
	rf_sp := reflect.ValueOf(&singlePing.Ping_headers).Elem()
	types := rf_pd.Type()

	for i := 0; i < rf_pd.NumField(); i++ {
		name := types.Field(i).Name
		field_pd := rf_pd.FieldByName(name)
		field_sp := rf_sp.FieldByName(name)
		field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
	}

	// Beam_array
	rf_pd = reflect.ValueOf(&pd.Beam_array).Elem()
	rf_sp = reflect.ValueOf(&singlePing.Beam_array).Elem()

	for _, name := range beam_names {
		if name == "IntensitySeries" {
			continue
		}
		field_pd := rf_pd.FieldByName(name)
		field_sp := rf_sp.FieldByName(name)
		field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
	}

	// Lon_lat
	pd.Lon_lat.Longitude = append(pd.Lon_lat.Longitude, singlePing.Lon_lat.Longitude...)
	pd.Lon_lat.Latitude = append(pd.Lon_lat.Latitude, singlePing.Lon_lat.Latitude...)

	// Sensor_metadata
	err := pd.Sensor_metadata.appendSensorMetadata(&singlePing.Sensor_metadata, sensor_id)
	if err != nil {
		return errors.Join(err, errors.New("Error appending SensorMetadata"))
	}

	if contains_intensity {
		// Brb_intensity
		pd.Brb_intensity.TimeSeries = append(pd.Brb_intensity.TimeSeries, singlePing.Brb_intensity.TimeSeries...)
		// pd.Brb_intensity.BottomDetect = append(singlePing.Brb_intensity.BottomDetect, singlePing.Brb_intensity.BottomDetect...)
		pd.Brb_intensity.BottomDetectIndex = append(pd.Brb_intensity.BottomDetectIndex, singlePing.Brb_intensity.BottomDetectIndex...)
		pd.Brb_intensity.StartRange = append(pd.Brb_intensity.StartRange, singlePing.Brb_intensity.StartRange...)
		pd.Brb_intensity.TsMean = append(pd.Brb_intensity.TsMean, singlePing.Brb_intensity.TsMean...)
		pd.Brb_intensity.sample_count = append(pd.Brb_intensity.sample_count, singlePing.Brb_intensity.sample_count...)

		// Sensor_imagery_metadata
		err := pd.Sensor_imagery_metadata.appendSensorImageryMetadata(&singlePing.Sensor_imagery_metadata, sensor_id)
		if err != nil {
			return errors.Join(err, errors.New("Error appending SensorImageryMetadata"))
		}
	}

	return nil
}

// newPingData initialises the PingData type with a set number of beams (total number of beams across n pings).
func newPingData(npings int, number_beams uint64, sensor_id SubRecordID, beam_names []string, contains_intensity bool) (pdata PingData) {
	var (
		brb        BrbIntensity
		sen_img_md SensorImageryMetadata
	)
	inumber_beams := int(number_beams)
	ping_headers := newPingHeaders(npings)
	beam_array := newBeamArray(inumber_beams, beam_names)
	sen_md := newSensorMetadata(npings, sensor_id)
	lonlat := LonLat{make([]float64, 0, number_beams), make([]float64, 0, number_beams)}

	// only allocate intensity and img_metadata slices if the GSF contains intensity
	if contains_intensity {
		brb = newBrbIntensity(int(number_beams))
		sen_img_md = newSensorImageryMetadata(npings, sensor_id)
	} else {
		brb = BrbIntensity{}
		sen_img_md = SensorImageryMetadata{}
	}

	pdata.Ping_headers = ping_headers
	pdata.Beam_array = beam_array
	pdata.Brb_intensity = brb
	pdata.Sensor_metadata = sen_md
	pdata.Sensor_imagery_metadata = sen_img_md
	pdata.Lon_lat = lonlat
	pdata.n_pings = uint64(npings)
	pdata.ba_subrecords = beam_names

	return pdata
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

// decode_ping_hdr reads and unscales the ping header information.
// Most fields are unscaled into float64 and then converted to float32.
// Fields that were encoded as 4bytes (int32, uint32) are returned as float64.
func decode_ping_hdr(reader *bytes.Reader, gsfd GsfDetails) PingHeader {
	var (
		hdr_base struct {
			Seconds         uint32
			Nano_seconds    uint32
			Longitude       uint32
			Latitude        uint32
			Number_beams    uint16
			Centre_beam     uint16
			Ping_flags      uint16
			Reserved        uint16
			Tide_corrector  uint16
			Depth_corrector uint32
			Heading         uint16
			Pitch           uint16
			Roll            uint16
			Heave           uint16
			Course          uint16
			Speed           uint16
		}
		hdr_xtra struct {
			Height             uint32
			Separation         uint32
			GPS_tide_corrector uint32
			Spare              [2]byte
		}
		hdr    PingHeader
		height float64
		sep    float64
		gps_tc float64
	)

	_ = binary.Read(reader, binary.BigEndian, &hdr_base)

	major, _ := gsfd.MajorMinor()
	if major > 2 {
		_ = binary.Read(reader, binary.BigEndian, &hdr_xtra)
		height = float64(int32(hdr_xtra.Height)) / SCALE_3_F64
		sep = float64(int32(hdr_xtra.Separation)) / SCALE_3_F64
		gps_tc = float64(int32(hdr_xtra.GPS_tide_corrector)) / SCALE_3_F64
	} else {
		height = NULL_HEIGHT
		sep = NULL_SEP
		gps_tc = NULL_GPS_TIDE_CORRECTOR
	}

	hdr.Timestamp = time.Unix(int64(hdr_base.Seconds), int64(hdr_base.Nano_seconds)).UTC()
	hdr.Longitude = float64(int32(hdr_base.Longitude)) / SCALE_7_F64
	hdr.Latitude = float64(int32(hdr_base.Latitude)) / SCALE_7_F64
	hdr.Number_beams = hdr_base.Number_beams
	hdr.Centre_beam = hdr_base.Centre_beam
	hdr.Ping_flags = hdr_base.Ping_flags
	hdr.Tide_corrector = float32(float64(int16(hdr_base.Tide_corrector)) / SCALE_2_F64)
	hdr.Depth_corrector = float64(int32(hdr_base.Depth_corrector)) / SCALE_2_F64
	hdr.Heading = float32(float64(hdr_base.Heading) / SCALE_2_F64)
	hdr.Pitch = float32(float64(int16(hdr_base.Pitch)) / SCALE_2_F64)
	hdr.Roll = float32(float64(int16(hdr_base.Roll)) / SCALE_2_F64)
	hdr.Heave = float32(float64(int16(hdr_base.Heave)) / SCALE_2_F64)
	hdr.Course = float32(float64(hdr_base.Course) / SCALE_2_F64)
	hdr.Speed = float32(float64(hdr_base.Speed) / SCALE_2_F64)
	hdr.Height = height
	hdr.Separation = sep
	hdr.GPS_tide_corrector = gps_tc

	return hdr
}

// SubRecHdr decodes the header for a SubRecord and constructs a SubRecord.
func SubRecHdr(reader *bytes.Reader, offset int64) SubRecord {
	var subrecord_hdr uint32

	_ = binary.Read(reader, binary.BigEndian, &subrecord_hdr)

	subrecord_id := (subrecord_hdr & 0xFF000000) >> 24
	subrecord_size := subrecord_hdr & 0x00FFFFFF

	byte_index := offset + 4

	subhdr := SubRecord{SubRecordID(subrecord_id), subrecord_size, byte_index} // include a byte_index??

	return subhdr
}

// scale_factors_rec decodes the scale factors subrecord defined by SCALE_FACTORS.
func scale_factors_rec(reader *bytes.Reader) (scale_factors map[SubRecordID]ScaleFactor, nbytes int64) {
	var (
		i            uint32
		num_factors  uint32
		scale_factor ScaleFactor
	)
	data := make([]uint32, 3) // id, scale, offset
	scale_factors = map[SubRecordID]ScaleFactor{}

	_ = binary.Read(reader, binary.BigEndian, &num_factors)
	nbytes = 4

	for i = 0; i < num_factors; i++ {
		_ = binary.Read(reader, binary.BigEndian, &data)

		subid := (data[0] & 0xFF000000) >> 24
		comp_flag := (data[0] & 0x00FF0000) >> 16
		comp := (comp_flag & 0x0F) == 1
		cnvrt_subid := SubRecordID(subid)
		field_size := comp_flag & 0xF0

		scale_factor = ScaleFactor{
			Id:               cnvrt_subid,
			ScaleOffset:      ScaleOffset{float64(data[1]), float64(int32(data[2]))},
			Compression_flag: comp_flag, // TODO; implement compression decoder
			Compressed:       comp,
			Field_size:       field_size, // this field doesn't appear to be used in the C code ???
		}

		nbytes += 12

		scale_factors[cnvrt_subid] = scale_factor
	}

	return scale_factors, nbytes
}

// ping_info decodes the SWATH_BATHYMETRY_PING record such as the header,
// SubRecord's and constructs the PingInfo type.
func ping_info(reader *bytes.Reader, rec RecordHdr, gsfd GsfDetails) PingInfo {
	var (
		idx     int64 = 0
		pinfo   PingInfo
		records      = make([]SubRecordID, 0, 32)
		sf      bool = false
		scl_fac map[SubRecordID]ScaleFactor
		nbytes  int64
	)

	datasize := int64(rec.Datasize)

	hdr := decode_ping_hdr(reader, gsfd)
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
func SwathBathymetryPingRec(buffer []byte, rec RecordHdr, pinfo PingInfo, sensor_id SubRecordID, gsfd GsfDetails) (PingData, error) {
	var (
		idx        int64 = 0
		beam_data  []float64
		ping_data  PingData
		beam_array BeamArray
		img_md     SensorImageryMetadata
		intensity  BrbIntensity
		sen_md     SensorMetadata
		ba_read    []string // keep track of which beam array records have been read
		err        error
	)
	ba_read = make([]string, 0, 30)

	reader := bytes.NewReader(buffer)

	hdr := decode_ping_hdr(reader, gsfd)
	idx += 56 // 56 bytes read for ping header

	for reader.Len() > 4 {

		// subrecord header
		// offset is used to track the start of the subrecord from the start
		// of the file as given by Record.Byte_index
		// Incase we wish to serialise the subrecord info along with the record info
		sub_rec := SubRecHdr(reader, rec.Byte_index+idx)
		idx += 4

		// the next two blocks of comments are kept for historical reference
		// TLDR; the INTENSITY_SERIES subrecords, based on testing (admittedly)
		// a half dozen or so files, recorded an incorrect datasize value.
		// As such, relying on reading the exact size into a memory buffer is futile
		// as in the examples it required reading past the memory buffer containing
		// the full record.

		// read the whole subrecord and form a new reader
		// i think this is easier than passing around how many
		// bytes are read from each func associated with decoding a subrecord
		// sr_buff = make([]byte, sub_rec.Datasize)
		// err := binary.Read(reader, binary.BigEndian, &sr_buff)
		// if err != nil {
		// 	// i've come across a subrecord, specifically INTENSITY_SERIES,
		// 	// where the subrecord would read read past the full record size
		// 	// by one byte. i.e. subrecord size is 10471, record size is 18332
		// 	// and current position is 7862.
		// 	// 7862 + 10471 = 18333
		// 	errn := errors.Join(
		// 		errors.New("Binary Read Failed"),
		// 		errors.New("Attempting to read SubRecord: "+SubRecordNames[sub_rec.Id]),
		// 		errors.New("SubRecord Datasize: "+strconv.Itoa(int(sub_rec.Datasize))),
		// 		errors.New("Record Datasize: "+strconv.Itoa(int(rec.Datasize))),
		// 		errors.New("Current byte location: "+strconv.Itoa(int(idx))),
		// 		errors.New("Record byte location: "+strconv.Itoa(int(rec.Byte_index))),
		// 		errors.New("SubRecord byte location: "+strconv.Itoa(int(sub_rec.Byte_index))),
		// 	)
		// 	return ping_data, errors.Join(err, errn)
		// }
		// sr_reader = bytes.NewReader(sr_buff)

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
			_, _ = scale_factors_rec(reader)

		// beam array subrecords
		case DEPTH:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)

			// converting to Z-axis domain (integrate with elevation)
			// TODO; loop over length, as range may copy the array
			for k, v := range beam_data {
				beam_data[k] = v * float64(-1.0)
			}
			beam_array.Z = beam_data
			ba_read = append(ba_read, pascalCase(SubRecordNames[DEPTH]))
		case ACROSS_TRACK:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				true,
			)
			beam_array.AcrossTrack = beam_data
			ba_read = append(ba_read, pascalCase(SubRecordNames[ACROSS_TRACK]))
		case ALONG_TRACK:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				true,
			)
			beam_array.AlongTrack = beam_data
			ba_read = append(ba_read, pascalCase(SubRecordNames[ALONG_TRACK]))
		case TRAVEL_TIME:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.TravelTime = beam_data
			ba_read = append(ba_read, pascalCase(SubRecordNames[TRAVEL_TIME]))
		case BEAM_ANGLE:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_TWO,
				true,
			)
			beam_array.BeamAngle = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[BEAM_ANGLE]))
		case MEAN_CAL_AMPLITUDE:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				true,
			)
			beam_array.MeanCalAmplitude = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[MEAN_CAL_AMPLITUDE]))
		case MEAN_REL_AMPLITUDE:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.MeanRelAmplitude = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[MEAN_REL_AMPLITUDE]))
		case ECHO_WIDTH:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.EchoWidth = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[ECHO_WIDTH]))
		case QUALITY_FACTOR:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_ONE,
				false,
			)
			beam_array.QualityFactor = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[QUALITY_FACTOR]))
		case RECEIVE_HEAVE:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_ONE,
				true,
			)
			beam_array.RecieveHeave = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[RECEIVE_HEAVE]))
		case DEPTH_ERROR:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_TWO,
				false,
			)
			beam_array.DepthError = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[DEPTH_ERROR]))
		case ACROSS_TRACK_ERROR:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_TWO,
				false,
			)
			beam_array.AcrossTrackError = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[ACROSS_TRACK_ERROR]))
		case ALONG_TRACK_ERROR:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_TWO,
				false,
			)
			beam_array.AlongTrackError = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[ALONG_TRACK_ERROR]))
		case NOMINAL_DEPTH:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.NominalDepth = beam_data
			ba_read = append(ba_read, pascalCase(SubRecordNames[NOMINAL_DEPTH]))
		case BEAM_FLAGS:
			beam_array.BeamFlags = DecodeBeamFlagsArray(
				reader,
				pinfo.Number_Beams,
			)
			ba_read = append(ba_read, pascalCase(SubRecordNames[BEAM_FLAGS]))
		case QUALITY_FLAGS:
			// obselete
			// TODO; has specific decoder
			panic("QUALITY_FLAGS subrecord has been superceded")
		case SIGNAL_TO_NOISE:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_ONE,
				true,
			)
			beam_array.SignalToNoise = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[SIGNAL_TO_NOISE]))
		case BEAM_ANGLE_FORWARD:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_TWO,
				false,
			)
			beam_array.BeamAngleForward = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[BEAM_ANGLE_FORWARD]))
		case VERTICAL_ERROR:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_TWO,
				false,
			)
			beam_array.VerticalError = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[VERTICAL_ERROR]))
		case HORIZONTAL_ERROR:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_TWO,
				false,
			)
			beam_array.HorizontalError = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[HORIZONTAL_ERROR]))
		case INTENSITY_SERIES:
			intensity, img_md, err = DecodeBrbIntensity(reader, pinfo.Number_Beams, sensor_id, pinfo.scale_factors[sub_rec.Id])
			if err != nil {
				return ping_data, err
			}
			ba_read = append(ba_read, pascalCase(SubRecordNames[INTENSITY_SERIES]))
		case SECTOR_NUMBER:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_ONE,
				false,
			)
			beam_array.SectorNumber = f64_to_u16(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[SECTOR_NUMBER]))
		case DETECTION_INFO:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_ONE,
				false,
			)
			beam_array.DetectionInfo = f64_to_u16(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[DETECTION_INFO]))
		case INCIDENT_BEAM_ADJ:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_ONE,
				true,
			)
			beam_array.IncidentBeamAdj = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[INCIDENT_BEAM_ADJ]))
		case SYSTEM_CLEANING:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_ONE,
				false,
			)
			beam_array.SystemCleaning = f64_to_u16(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[SYSTEM_CLEANING]))
		case DOPPLER_CORRECTION:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_ONE,
				true,
			)
			beam_array.DopplerCorrection = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[DOPPLER_CORRECTION]))
		case SONAR_VERT_UNCERTAINTY:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_TWO,
				false,
			)
			beam_array.SonarVertUncertainty = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[SONAR_VERT_UNCERTAINTY]))
		case SONAR_HORZ_UNCERTAINTY:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				BYTES_PER_BEAM_TWO,
				false,
			)
			beam_array.SonarHorzUncertainty = f64_to_f32(beam_data)
			ba_read = append(ba_read, pascalCase(SubRecordNames[SONAR_HORZ_UNCERTAINTY]))
		case DETECTION_WINDOW:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.DetectionWindow = beam_data
			ba_read = append(ba_read, pascalCase(SubRecordNames[DETECTION_WINDOW]))
		case MEAN_ABS_COEF:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.MeanAbsCoef = beam_data
			ba_read = append(ba_read, pascalCase(SubRecordNames[MEAN_ABS_COEF]))
		case TVG_DB:
			beam_data = sub_rec.DecodeSubRecArray(
				reader,
				pinfo.Number_Beams,
				pinfo.scale_factors[sub_rec.Id],
				bytes_per_beam,
				false,
			)
			beam_array.TvgDb = beam_data
			ba_read = append(ba_read, pascalCase(SubRecordNames[TVG_DB]))

		// sensor specific subrecords
		case SEABEAM:
			// DecodeSeabeam
			sen_md.Seabeam, err = DecodeSeabeamSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case EM12:
			// DecodeEM12
			sen_md.Em12, err = DecodeEm12Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case EM100:
			// DecodeEM100
			sen_md.Em100, err = DecodeEm100Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case EM950:
			// DecodeEM950
			sen_md.Em950, err = DecodeEm950Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case EM121A:
			// DecodeEM121A
			sen_md.Em121A, err = DecodeEm121ASpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case EM121:
			// DecodeEM121
			sen_md.Em121, err = DecodeEm121Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case SASS: // obsolete
			// DecodeSASS
			sen_md.Sass, err = DecodeSassSpecfic(reader)
			if err != nil {
				return ping_data, err
			}
		case SEAMAP:
			// DecodeSeaMap
			sen_md.SeaMap, err = DecodeSeaMapSpecific(reader, gsfd)
			if err != nil {
				return ping_data, err
			}
		case SEABAT:
			// DecodeSeaBat
			sen_md.SeaBat, err = DecodeSeaBatSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case EM1000:
			// DecodeEM1000
			sen_md.Em1000, err = DecodeEm1000Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case TYPEIII_SEABEAM: // obsolete
			// DecodeTypeIII
			sen_md.TypeIIISeabeam, err = DecodeTypeIIISeabeamSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case SB_AMP:
			// DecodeSBAmp
			sen_md.SbAmp, err = DecodeSbAmpSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case SEABAT_II:
			// DecodeSeaBatII
			sen_md.SeaBatII, err = DecodeSeaBatIISpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case SEABAT_8101:
			// DecodeSeaBat8101
			sen_md.SeaBat8101, err = DecodeSeaBat8101Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case SEABEAM_2112:
			// DecodeSeaBeam2112
			sen_md.Seabeam2112, err = DecodeSeabeam2112Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case ELAC_MKII:
			// DecodeElacMkII
			sen_md.ElacMkII, err = DecodeElacMkIISpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case CMP_SAAS: // CMP (compressed), should be used in place of SASS
			// DecodeCmpSass
			sen_md.CmpSass, err = DecodeCmpSass(reader)
			if err != nil {
				return ping_data, err
			}
		case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
			// DecodeReson8100
			sen_md.Reson8100, err = DecodeReson8100(reader)
			if err != nil {
				return ping_data, err
			}
		case EM120, EM300, EM1002, EM2000, EM3000, EM3002, EM3000D, EM3002D, EM121A_SIS:
			// DecodeEM3
			sen_md.Em3, err = DecodeEm3Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case EM710, EM302, EM122, EM2040, ME70BO:
			// DecodeEM4
			sen_md.Em4, err = DecodeEm4Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case GEOSWATH_PLUS:
			// DecodeGeoSwathPlus
			sen_md.GeoSwathPlus, err = DecodeGeoSwathPlusSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case KLEIN_5410_BSS:
			// DecodeKlein5410Bss
			sen_md.Klein5410Bss, err = DecodeKlein5410BssSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case RESON_7125:
			// DecodeReson7100
			sen_md.Reson7100, err = DecodeReson7100Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case EM300_RAW, EM1002_RAW, EM2000_RAW, EM3000_RAW, EM120_RAW, EM3002_RAW, EM3000D_RAW, EM3002D_RAW, EM121A_SIS_RAW:
			// DecodeEM3Raw
			sen_md.Em3Raw, err = DecodeEm3RawSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case DELTA_T:
			// DecodeDeltaT
			sen_md.DeltaT, err = DecodeDeltaTSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case R2SONIC_2022, R2SONIC_2024, R2SONIC_2020:
			// DecodeR2Sonic
			sen_md.R2Sonic, err = DecodeR2SonicSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case SR_NOT_DEFINED: // the spec makes no mention of ID 154
			errn := errors.Join(
				errors.New("Error, Subrecord ID 154 is not defined."),
				errors.New("Record index: "+strconv.Itoa(int(rec.Byte_index))),
				errors.New("Record datasize: "+strconv.Itoa(int(rec.Datasize))),
				errors.New("SubRecord index: "+strconv.Itoa(int(sub_rec.Byte_index))),
				errors.New("SubRecord datasize: "+strconv.Itoa(int(sub_rec.Datasize))),
				errors.New("Current byte location (relative to current record): "+strconv.Itoa(int(idx))),
			)
			panic(errn)
		case RESON_TSERIES:
			// DecodeResonTSeries
			sen_md.ResonTSeries, err = DecodeResonTSeriesSonicSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case KMALL:
			// DecodeKMALL
			sen_md.Kmall, err = DecodeKmallSpecific(reader)
			if err != nil {
				return ping_data, err
			}

			// single beam swath sensor specific subrecords
		case SWATH_SB_ECHOTRAC, SWATH_SB_BATHY2000, SWATH_SB_PDD:
			// DecodeSBEchotrac
			sen_md.SwathSbEchotrac, err = DecodeSwathSbEchotracSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case SWATH_SB_MGD77:
			// DecodeSBMGD77
			sen_md.SwathSbMgd77, err = DecodeSwathSbMGD77Specific(reader)
			if err != nil {
				return ping_data, err
			}
		case SWATH_SB_BDB:
			// DecodeSBBDB
			sen_md.SwathSbBdb, err = DecodeSwathSbBdbSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case SWATH_SB_NOSHDB:
			// DecodeSBNOSHDB
			sen_md.SwathSbNoShDb, err = DecodeSwathSbNoShDbSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		case SWATH_SB_NAVISOUND:
			// DecodeSBNavisound
			sen_md.SwathSbNavisound, err = DecodeSwathSbNavisoundSpecific(reader)
			if err != nil {
				return ping_data, err
			}
		}
	}

	geocoef := NewCoefWgs84()
	lonlat := beam_array.BeamsLonLat(hdr.Longitude, hdr.Latitude, hdr.Heading, geocoef)

	ping_headers := PingHeaders{
		[]time.Time{hdr.Timestamp},
		[]float64{hdr.Longitude},
		[]float64{hdr.Latitude},
		[]uint16{hdr.Number_beams},
		[]uint16{hdr.Centre_beam},
		[]float32{hdr.Tide_corrector},
		[]float64{hdr.Depth_corrector},
		[]float32{hdr.Heading},
		[]float32{hdr.Pitch},
		[]float32{hdr.Roll},
		[]float32{hdr.Heave},
		[]float32{hdr.Course},
		[]float32{hdr.Speed},
		[]float64{hdr.Height},
		[]float64{hdr.Separation},
		[]float64{hdr.GPS_tide_corrector},
		[]uint16{hdr.Ping_flags},
	}

	ping_data.Ping_headers = ping_headers
	ping_data.Beam_array = beam_array
	ping_data.Brb_intensity = intensity
	ping_data.Sensor_imagery_metadata = img_md
	ping_data.Sensor_metadata = sen_md
	ping_data.Lon_lat = lonlat
	ping_data.n_pings = uint64(1)
	ping_data.ba_subrecords = ba_read

	return ping_data, err
}

// writeBeamData serialises the beam data to a sparse TileDB array
// using longitude and latitude as the dimensional axes.
func (pd *PingData) writeBeamData(ctx *tiledb.Context, array *tiledb.Array, ping_beam_ids *PingBeamNumbers) error {
	schema, err := array.Schema()
	if err != nil {
		errn := errors.New("Error retrieving array schema")
		return errors.Join(err, errn)
	}
	defer schema.Free()

	arr_type, err := schema.Type()
	if err != nil {
		return err
	}

	// query construction
	query, err := tiledb.NewQuery(ctx, array)
	if err != nil {
		errn := errors.New("Error creating TileDB query")
		return errors.Join(err, errn)
	}
	defer query.Free()

	if arr_type == tiledb.TILEDB_DENSE {
		err = query.SetLayout(tiledb.TILEDB_ROW_MAJOR)
		if err != nil {
			errn := errors.New("Error setting TileDB layout")
			return errors.Join(err, errn)
		}

		ping_start := ping_beam_ids.PingNumber[0]
		end_idx := len(ping_beam_ids.PingNumber) - 1
		ping_end := ping_beam_ids.PingNumber[end_idx]

		// this may require rethinking so that we ensure to get max beams as end beam
		// want to ensure that max beams is written, not the number of beams for given ping
		beam_end_idx := len(ping_beam_ids.BeamNumber) - 1
		beam_end := ping_beam_ids.BeamNumber[beam_end_idx]

		rng_ping := tiledb.MakeRange(ping_start, ping_end)
		rng_beam := tiledb.MakeRange(uint64(0), beam_end)

		subarr, err := array.NewSubarray()
		if err != nil {
			errn := errors.New("Error creating TileDB NewSubarray")
			return errors.Join(err, errn)
		}
		defer subarr.Free()

		subarr.AddRangeByName("PingNumber", rng_ping)
		subarr.AddRangeByName("BeamNumber", rng_beam)

		err = query.SetSubarray(subarr)
		if err != nil {
			errn := errors.New("Error setting TileDB Subarray")
			return errors.Join(err, errn)
		}
	} else {
		err = query.SetLayout(tiledb.TILEDB_UNORDERED)
		if err != nil {
			errn := errors.New("Error setting TileDB layout")
			return errors.Join(err, errn)
		}
	}

	// should make for simpler code, if reflect is used to get the type's
	// names and values (slice)
	// For the time being, using a case switch and being explicit works just
	// fine, albeit more code
	// TODO; look at replacing most of the following with reflect

	// dimensional axes buffers (or attributes depending on sparse/dense array)
	// X & Y for sparse array
	_, err = query.SetDataBuffer("X", pd.Lon_lat.Longitude)
	if err != nil {
		errn := errors.New("Error setting TileDB data buffer for dimension/attribute: X")
		return errors.Join(err, errn)
	}

	_, err = query.SetDataBuffer("Y", pd.Lon_lat.Latitude)
	if err != nil {
		errn := errors.New("Error setting TileDB data buffer for dimension/attribute: Y")
		return errors.Join(err, errn)
	}

	if arr_type != tiledb.TILEDB_DENSE {
		// ping and beam ids buffers are only set for sparse arrays
		_, err = query.SetDataBuffer("PingNumber", ping_beam_ids.PingNumber)
		if err != nil {
			errn := errors.New("Error setting TileDB data buffer for dimension/attribute: PingNumber")
			return errors.Join(err, errn)
		}

		_, err = query.SetDataBuffer("BeamNumber", ping_beam_ids.BeamNumber)
		if err != nil {
			errn := errors.New("Error setting TileDB data buffer for dimension/attribute: BeamNumber")
			return errors.Join(err, errn)
		}
	}

	// beam array buffers
	for _, name := range pd.ba_subrecords {
		subr_id := BeamDataName2SubRecordID[name]

		switch subr_id {
		case DEPTH:
			_, err = query.SetDataBuffer(name, pd.Beam_array.Z)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: Z")
				return errors.Join(err, errn)
			}
		case ACROSS_TRACK:
			_, err = query.SetDataBuffer(name, pd.Beam_array.AcrossTrack)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: AcrossTrack")
				return errors.Join(err, errn)
			}
		case ALONG_TRACK:
			_, err = query.SetDataBuffer(name, pd.Beam_array.AlongTrack)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: AlongTrack")
				return errors.Join(err, errn)
			}
		case TRAVEL_TIME:
			_, err = query.SetDataBuffer(name, pd.Beam_array.TravelTime)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: TravelTime")
				return errors.Join(err, errn)
			}
		case BEAM_ANGLE:
			_, err = query.SetDataBuffer(name, pd.Beam_array.BeamAngle)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: BeamAngle")
				return errors.Join(err, errn)
			}
		case MEAN_CAL_AMPLITUDE:
			_, err = query.SetDataBuffer(name, pd.Beam_array.MeanCalAmplitude)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: MeanCalAmplitude")
				return errors.Join(err, errn)
			}
		case MEAN_REL_AMPLITUDE:
			_, err = query.SetDataBuffer(name, pd.Beam_array.MeanRelAmplitude)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: MeanRelAmplitude")
				return errors.Join(err, errn)
			}
		case ECHO_WIDTH:
			_, err = query.SetDataBuffer(name, pd.Beam_array.EchoWidth)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: EchoWidth")
				return errors.Join(err, errn)
			}
		case QUALITY_FACTOR:
			_, err = query.SetDataBuffer(name, pd.Beam_array.QualityFactor)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: QualityFactor")
				return errors.Join(err, errn)
			}
		case RECEIVE_HEAVE:
			_, err = query.SetDataBuffer(name, pd.Beam_array.RecieveHeave)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: RecieveHeave")
				return errors.Join(err, errn)
			}
		case DEPTH_ERROR:
			_, err = query.SetDataBuffer(name, pd.Beam_array.DepthError)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: DepthError")
				return errors.Join(err, errn)
			}
		case ACROSS_TRACK_ERROR:
			_, err = query.SetDataBuffer(name, pd.Beam_array.AcrossTrackError)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: AcrossTrackError")
				return errors.Join(err, errn)
			}
		case ALONG_TRACK_ERROR:
			_, err = query.SetDataBuffer(name, pd.Beam_array.AlongTrackError)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: AlongTrackError")
				return errors.Join(err, errn)
			}
		case NOMINAL_DEPTH:
			_, err = query.SetDataBuffer(name, pd.Beam_array.NominalDepth)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: NominalDepth")
				return errors.Join(err, errn)
			}
		case QUALITY_FLAGS:
			_, err = query.SetDataBuffer(name, pd.Beam_array.QualityFlags)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: QualityFlags")
				return errors.Join(err, errn)
			}
		case BEAM_FLAGS:
			_, err = query.SetDataBuffer(name, pd.Beam_array.BeamFlags)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: BeamFlags")
				return errors.Join(err, errn)
			}
		case SIGNAL_TO_NOISE:
			_, err = query.SetDataBuffer(name, pd.Beam_array.SignalToNoise)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: SignalToNoise")
				return errors.Join(err, errn)
			}
		case BEAM_ANGLE_FORWARD:
			_, err = query.SetDataBuffer(name, pd.Beam_array.BeamAngleForward)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: BeamAngleForward")
				return errors.Join(err, errn)
			}
		case VERTICAL_ERROR:
			_, err = query.SetDataBuffer(name, pd.Beam_array.VerticalError)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: VerticalError")
				return errors.Join(err, errn)
			}
		case HORIZONTAL_ERROR:
			_, err = query.SetDataBuffer(name, pd.Beam_array.HorizontalError)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: HorizontalError")
				return errors.Join(err, errn)
			}
		case INTENSITY_SERIES:
			// offset buffer (intensity timeseries is variable length)
			n_obs := uint64(len(pd.Lon_lat.Longitude))
			arr_offset := make([]uint64, n_obs)
			offset := uint64(0)
			bytes_val := uint64(8) // may look confusing with uint64, so 8*bytes for float64

			for i := uint64(0); i < n_obs; i++ {
				arr_offset[i] = offset
				sample := uint64(pd.Brb_intensity.sample_count[i])

				// handle case with no sample counts, as we've inserted a NaN
				if sample == uint64(0) {
					offset += uint64(1) * bytes_val
				} else {
					offset += sample * bytes_val
				}
			}

			_, err = query.SetOffsetsBuffer("TimeSeries", arr_offset)
			if err != nil {
				errn := errors.New("Error setting TileDB offsets data buffer for attribute: TimeSeries")
				return errors.Join(err, errn)
			}

			_, err = query.SetDataBuffer("TimeSeries", pd.Brb_intensity.TimeSeries)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: TimeSeries")
				return errors.Join(err, errn)
			}

			// other non-var-length fields
			// _, err = query.SetDataBuffer("BottomDetect", pd.Brb_intensity.BottomDetect)
			// if err != nil {
			// 	return errors.Join(ErrWriteBdTdb, err)
			// }

			_, err = query.SetDataBuffer("BottomDetectIndex", pd.Brb_intensity.BottomDetectIndex)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: BottomDetectIndex")
				return errors.Join(err, errn)
			}

			_, err = query.SetDataBuffer("StartRange", pd.Brb_intensity.StartRange)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: StartRange")
				return errors.Join(err, errn)
			}

			_, err = query.SetDataBuffer("TsMean", pd.Brb_intensity.TsMean)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: TsMean")
				return errors.Join(err, errn)
			}
		case SECTOR_NUMBER:
			_, err = query.SetDataBuffer(name, pd.Beam_array.SectorNumber)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: SectorNumber")
				return errors.Join(err, errn)
			}
		case DETECTION_INFO:
			_, err = query.SetDataBuffer(name, pd.Beam_array.DetectionInfo)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: DetectionInfo")
				return errors.Join(err, errn)
			}
		case INCIDENT_BEAM_ADJ:
			_, err = query.SetDataBuffer(name, pd.Beam_array.IncidentBeamAdj)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: IncidentBeamAdj")
				return errors.Join(err, errn)
			}
		case SYSTEM_CLEANING:
			_, err = query.SetDataBuffer(name, pd.Beam_array.SystemCleaning)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: SystemCleaning")
				return errors.Join(err, errn)
			}
		case DOPPLER_CORRECTION:
			_, err = query.SetDataBuffer(name, pd.Beam_array.DopplerCorrection)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: DopplerCorrection")
				return errors.Join(err, errn)
			}
		case SONAR_VERT_UNCERTAINTY:
			_, err = query.SetDataBuffer(name, pd.Beam_array.SonarVertUncertainty)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: SonarVertUncertainty")
				return errors.Join(err, errn)
			}
		case SONAR_HORZ_UNCERTAINTY:
			_, err = query.SetDataBuffer(name, pd.Beam_array.SonarHorzUncertainty)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: SonarHorzUncertainty")
				return errors.Join(err, errn)
			}
		case DETECTION_WINDOW:
			_, err = query.SetDataBuffer(name, pd.Beam_array.DetectionWindow)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: DetectionWindow")
				return errors.Join(err, errn)
			}
		case MEAN_ABS_COEF:
			_, err = query.SetDataBuffer(name, pd.Beam_array.MeanAbsCoef)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: MeanAbsCoef")
				return errors.Join(err, errn)
			}
		case TVG_DB:
			_, err = query.SetDataBuffer(name, pd.Beam_array.TvgDb)
			if err != nil {
				errn := errors.New("Error setting TileDB data buffer for attribute: TvgDb")
				return errors.Join(err, errn)
			}
		}
	}

	// write the data and flush
	err = query.Submit()
	if err != nil {
		errn := errors.New("Error submitting TileDB query")
		return errors.Join(err, errn)
	}

	// not applicable, as layout is tiledb.TILEDB_UNORDERED
	// (tiledb lib will reorder it)
	err = query.Finalize()
	if err != nil {
		errn := errors.New("Error finalising TileDB query")
		return errors.Join(err, errn)
	}

	return nil
}

// writePingHeaders is a helper to serialise the PingHeaders
// to the respective TileDB array.
func (ph *PingHeaders) writePingHeaders(ctx *tiledb.Context, array *tiledb.Array, ping_start, ping_end uint64) error {
	// query construction
	query, err := tiledb.NewQuery(ctx, array)
	if err != nil {
		return errors.Join(ErrWriteMdTdb, err)
	}
	defer query.Free()

	err = query.SetLayout(tiledb.TILEDB_ROW_MAJOR)
	if err != nil {
		return errors.Join(ErrWriteMdTdb, err)
	}

	// define the subarray (dim coordinates that we'll write into)
	subarr, err := array.NewSubarray()
	if err != nil {
		errn := errors.New("Error defining subarray for writing PingHeaders")
		return errors.Join(err, errn)
	}
	defer subarr.Free()

	rng := tiledb.MakeRange(ping_start, ping_end)
	subarr.AddRangeByName("PING_ID", rng)
	err = query.SetSubarray(subarr)
	if err != nil {
		return errors.Join(ErrWriteMdTdb, err)
	}

	err = setStructFieldBuffers(query, ph)
	if err != nil {
		return errors.Join(err, errors.New("Error writing PingHeaders"))
	}

	// write the data flush
	err = query.Submit()
	if err != nil {
		errn := errors.New("Error submitting TileDB query")
		return errors.Join(err, errn)
	}

	err = query.Finalize()
	if err != nil {
		errn := errors.New("Error finalising TileDB query")
		return errors.Join(err, errn)
	}

	return nil
}

// toTileDB is a helper routine to serialise the beam data and the ping metadata to
// TileDB arrays.
// The beam data consists of the BeamArray, BrbIntensity (if intensity exists),
// and PingBeamNumbers.
// The ping metadata consists of the PingHeaders, sensor metadata, and
// sensor imagery (if intensity exists)
func (pd *PingData) toTileDB(ph_array, s_md_array, si_md_array, bd_array *tiledb.Array, ctx *tiledb.Context, ping_beam_ids *PingBeamNumbers, sensor_id SubRecordID, contains_intensity bool) error {
	ping_start := ping_beam_ids.PingNumber[0]
	end_idx := len(ping_beam_ids.PingNumber) - 1
	ping_end := ping_beam_ids.PingNumber[end_idx]

	// PingHeaders
	err := pd.Ping_headers.writePingHeaders(ctx, ph_array, ping_start, ping_end)
	if err != nil {
		errn := errors.New("Error writing PingHeaders")
		return errors.Join(err, errn)
	}

	// SensorMetadata
	err = pd.Sensor_metadata.writeSensorMetadata(ctx, s_md_array, sensor_id, ping_start, ping_end)
	if err != nil {
		errn := errors.New("Error writing SensorMetadata")
		return errors.Join(err, errn)
	}

	// SensorImageryMetadata
	if contains_intensity {
		err = pd.Sensor_imagery_metadata.writeSensorImageryMetadata(ctx, si_md_array, sensor_id, ping_start, ping_end)
		if err != nil {
			errn := errors.New("Error writing SensorImageryMetadata")
			return errors.Join(err, errn)
		}
	}

	// beam array data; BeamArray, PingBeamNumbers, LonLat, BrbIntensity
	err = pd.writeBeamData(ctx, bd_array, ping_beam_ids)
	if err != nil {
		errn := errors.New("Error writing beam data")
		return errors.Join(err, errn)
	}

	return nil
}

// SbpToTileDB converts SwathBathymetryPing Records to TileDB arrays.
// Beam array data will be converted to a sparse point cloud using
// longitude and latitude (named as X and Y) dimensional axes.
// SensorMetadata and SensorImageryMetadata subrecords will be written to a dense
// array along with the PingHeader data. This dense array will be single axis, using
// pings [0, n] as the axis units, akin to a table of data with n-rows where n is
// the number of pings.
// As there potentially are a lot of ping records, this process will be chunked
// into roughly chunks of 1000 pings in size given by:
// github.com/samber/lo.Chunk([]ping_records, 1000).
// In time (and interest), this chunk size can be made configurable.
// There is potential to create the beam arrays as a 2D dense array using [ping, beam]
// as the dimensional axes. The rationale is for input into algorithms that require
// input based on the sensor configuration; such as a beam adjacency filter that
// operates on a ping by ping basis.
func (g *GsfFile) SbpToTileDB(fi *FileInfo, ctx *tiledb.Context, grp *tiledb.Group, outdir_uri string, dense_bd bool) error {
	var (
		ping_data       PingData
		ping_data_chunk PingData
		ping_beam_ids   PingBeamNumbers
		err             error
		number_beams    uint64

		// declaring these so they can be passed through to various
		// funcs, even if no intensity data is present
		si_md_array *tiledb.Array
	)
	number_beams = 0

	rec_name := RecordNames[SWATH_BATHYMETRY_PING]
	total_pings := fi.Record_Counts[rec_name]
	ping_records := fi.Index.Record_Index[rec_name]

	// schema and cleanup subrecord names to match the BeamArray fields names
	sr_schema := make([]string, 0, len(fi.SubRecord_Schema))
	for _, v := range fi.SubRecord_Schema {
		sr_schema = append(sr_schema, v)
	}
	sr_schema_c := make([]string, len(sr_schema))
	for k, v := range sr_schema {
		sr_schema_c[k] = pascalCase(v)
	}

	contains_intensity := lo.Contains(sr_schema, SubRecordNames[INTENSITY_SERIES])
	sensor_id := SubRecordID(fi.Metadata.Sensor_Info.Sensor_ID)

	// output locations
	ph_name := "PingHeader.tiledb"
	ph_aname := "PingHeader"
	s_md_name := "SensorMetadata.tiledb"
	s_md_aname := "SensorMetadata"
	si_md_name := "SensorImageryMetadata.tiledb"
	si_md_aname := "SensorImageryMetadata"
	bd_name := "BeamData.tiledb"
	bd_aname := "BeamData"
	ph_uri := filepath.Join(outdir_uri, ph_name)
	s_md_uri := filepath.Join(outdir_uri, s_md_name)
	si_md_uri := filepath.Join(outdir_uri, si_md_name)
	bd_uri := filepath.Join(outdir_uri, bd_name)

	err = fi.pingTdbArrays(ctx, ph_uri, s_md_uri, si_md_uri, bd_uri, dense_bd)
	if err != nil {
		return errors.Join(err, errors.New("Error creating PingData TileDB arrays"))
	}

	// add arrays to tiledb group
	err = grp.AddMember(ph_name, ph_aname, true)
	if err != nil {
		return errors.Join(err, errors.New("Error adding ping headers to group"))
	}
	err = grp.AddMember(s_md_name, s_md_aname, true)
	if err != nil {
		return errors.Join(err, errors.New("Error adding sensor metadata to group"))
	}
	err = grp.AddMember(bd_name, bd_aname, true)
	if err != nil {
		return errors.Join(err, errors.New("Error adding beam data to group"))
	}

	// open the arrays for writing

	// PingHeaders
	ph_array, err := ArrayOpenWrite(ctx, ph_uri)
	if err != nil {
		return errors.Join(err, ErrWriteBdTdb, errors.New("Error opening (w) PingHeaders TileDB array"))
	}
	defer ph_array.Free()
	defer ph_array.Close()

	// SensorMetadata
	s_md_array, err := ArrayOpenWrite(ctx, s_md_uri)
	if err != nil {
		return errors.Join(err, ErrWriteBdTdb, errors.New("Error opening (w) SensorMetadata TileDB array"))
	}
	defer s_md_array.Free()
	defer s_md_array.Close()

	// SensorImageryMetadata (only exists if intensity exists)
	if contains_intensity {
		err = grp.AddMember(si_md_name, si_md_aname, true)
		if err != nil {
			return errors.Join(err, errors.New("Error adding sensor imagery metadata to group"))
		}

		si_md_array, err = ArrayOpenWrite(ctx, si_md_uri)
		if err != nil {
			return errors.Join(err, ErrWriteBdTdb, errors.New("Error opening (w) SensorImageryMetadata TileDB array"))
		}
		defer si_md_array.Free()
		defer si_md_array.Close()
	}

	// beam data; BeamArray, LonLat, PingBeamNumbers, BrbIntensity
	bd_array, err := ArrayOpenWrite(ctx, bd_uri)
	if err != nil {
		return errors.Join(err, ErrWriteBdTdb, errors.New("Error opening (w) TileDB beam array"))
	}
	defer bd_array.Free()
	defer bd_array.Close()

	// setup the chunks to process
	idxs := make([]uint64, total_pings)
	for i := uint64(0); i < total_pings; i++ {
		idxs[i] = i
	}
	chunks := lo.Chunk(idxs, int(1000))

	// need some info to initialise arrays that will get written into
	// also need to cater for intensity, which at the moment are stored
	// as 1-D, with count offsets (to define var length)
	for _, chunk := range chunks {

		n_pings := len(chunk)

		// initialise beam arrays, backscatter, lonlat
		// arrays for ping and beam numbers
		if dense_bd {
			number_beams = uint64(n_pings) * uint64(fi.Metadata.Quality_Info.Min_Max_Beams[1])
			ping_data_chunk = newPingData(n_pings, number_beams, sensor_id, sr_schema_c, contains_intensity)
			ping_beam_ids = newPingBeamNumbers(int(number_beams))
		} else {
			number_beams = 0
			for _, idx := range chunk {
				number_beams += uint64(fi.Ping_Info[idx].Number_Beams)
			}
			ping_data_chunk = newPingData(n_pings, number_beams, sensor_id, sr_schema_c, contains_intensity)
			ping_beam_ids = newPingBeamNumbers(int(number_beams))
		}

		// for dense_ba, need to account for failed ping read and fill with nulls
		// also need to account for adding null data for additional beams if ping.nbeams < max_beams

		// loop over each ping for this chunk of pings
		for _, idx := range chunk {
			rec := ping_records[idx]
			pinfo := fi.Ping_Info[idx]

			// seek to record
			_, _ = g.Stream.Seek(rec.Byte_index, 0)

			buffer := make([]byte, rec.Datasize)
			_ = binary.Read(g.Stream, binary.BigEndian, &buffer)
			ping_data, err = SwathBathymetryPingRec(buffer, rec, pinfo, sensor_id, fi.Metadata.GSF_Details)
			if err != nil {
				// for the time being, rather than stop and return,
				// log an issue, and keep processing
				errn := errors.New("Error reading ping: " + strconv.Itoa(int(idx)))
				// return errors.Join(err, errn)
				log.Println(errors.Join(err, errn))
				log.Println("Skipping PingID: ", idx)
				continue
			}

			// appending and null filling
			_ = ping_beam_ids.appendPingBeam(idx, pinfo.Number_Beams)
			_ = ping_data_chunk.appendPingData(&ping_data, contains_intensity, sensor_id, sr_schema_c)
			_ = ping_data_chunk.fillNulls(&ping_data, sensor_id)

			if dense_bd {
				pad_size := fi.Metadata.Quality_Info.Min_Max_Beams[1] - pinfo.Number_Beams
				if pad_size > uint16(0) {
					_ = ping_data_chunk.padDense(pad_size)
					_ = ping_beam_ids.padPingBeam(idx, pinfo.Number_Beams, pad_size)
				}
			}
		}

		// serialise chunk to the TileDB array
		err = ping_data_chunk.toTileDB(
			ph_array,
			s_md_array,
			si_md_array,
			bd_array,
			ctx,
			&ping_beam_ids,
			sensor_id,
			contains_intensity,
		)
		if err != nil {
			return errors.Join(err, errors.New("Error writing PingData chunk"))
		}
	}

	return nil
}
