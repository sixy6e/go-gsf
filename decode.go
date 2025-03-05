package gsf

import (
	"github.com/samber/lo"
)

type RecordID uint32
type SubRecordID int32
type SensorID uint16

const (
	NEXT_RECORD        RecordID = 0
	BEAM_WIDTH_UNKNOWN float32  = -1.0
)

// float32 and float64 scale factors
// in general I think they're a bit more readable than 1.0e7
// they're easy to change if there is some desire to in future
const (
	SCALE_1_F32  float32 = 10.0
	SCALE_1_F64  float64 = 10.0
	SCALE_2_F32  float32 = 100.0
	SCALE_2_F64  float64 = 100.0
	SCALE_3_F32  float32 = 1_000.0
	SCALE_3_F64  float64 = 1_000.0
	SCALE_4_F32  float32 = 10_000.0
	SCALE_4_F64  float64 = 10_000.0
	SCALE_5_F32  float32 = 100_000.0
	SCALE_5_F64  float64 = 100_000.0
	SCALE_6_F32  float32 = 1_000_000.0
	SCALE_6_F64  float64 = 1_000_000.0
	SCALE_7_F32  float32 = 10_000_000.0
	SCALE_7_F64  float64 = 10_000_000.0
	SCALE_8_F32  float32 = 100_000_000.0
	SCALE_8_F64  float64 = 100_000_000.0
	SCALE_9_F32  float32 = 1_000_000_000.0
	SCALE_9_F64  float64 = 1_000_000_000.0
	SCALE_10_F32 float32 = 10_000_000_000.0
	SCALE_10_F64 float64 = 10_000_000_000.0
)

const (
	MAX_BEAM_ARRAY_SUBRECORD_ID SubRecordID = 31
)

// Base record IDs.
const (
	HEADER RecordID = 1 + iota
	SWATH_BATHYMETRY_PING
	SOUND_VELOCITY_PROFILE
	PROCESSING_PARAMETERS
	SENSOR_PARAMETERS
	COMMENT
	HISTORY
	NAVIGATION_ERROR // obsolete
	SWATH_BATHY_SUMMARY
	SINGLE_BEAM_PING    // use discouraged
	HV_NAVIGATION_ERROR // replaces navigation error
	ATTITUDE            // 12
)

// Swath bathy ping subrecord IDs.
const (
	DEPTH SubRecordID = 1 + iota
	ACROSS_TRACK
	ALONG_TRACK
	TRAVEL_TIME
	BEAM_ANGLE
	MEAN_CAL_AMPLITUDE
	MEAN_REL_AMPLITUDE
	ECHO_WIDTH
	QUALITY_FACTOR // replaces quality flags
	RECEIVE_HEAVE
	DEPTH_ERROR        // obsolete
	ACROSS_TRACK_ERROR // obsolete
	ALONG_TRACK_ERROR  // obsolete
	NOMINAL_DEPTH
	QUALITY_FLAGS // considered obsolete
	BEAM_FLAGS
	SIGNAL_TO_NOISE
	BEAM_ANGLE_FORWARD
	VERTICAL_ERROR   // replaces depth error
	HORIZONTAL_ERROR // replaces across track error
	INTENSITY_SERIES
	SECTOR_NUMBER
	DETECTION_INFO
	INCIDENT_BEAM_ADJ
	SYSTEM_CLEANING
	DOPPLER_CORRECTION
	SONAR_VERT_UNCERTAINTY
	SONAR_HORZ_UNCERTAINTY
	DETECTION_WINDOW
	MEAN_ABS_COEF // 30
	TVG_DB        // 31 (added in GSF v3.10)
)

// General subrecord IDs.
const (
	UNKNOWN       SubRecordID = 0
	SCALE_FACTORS SubRecordID = 100
	SB_UNKNOWN    SubRecordID = 0 // single_beam
)

// Additional swath bathymetry subrecord IDs; Sensor specific records.
// The scale factors contained within the SCALE_FACTORS subrecord do not apply here.
// (Possibly in relation to the intensity data).
const (
	SEABEAM SubRecordID = 102 + iota
	EM12
	EM100
	EM950
	EM121A
	EM121
	SASS // obsolete
	SEAMAP
	SEABAT
	EM1000
	TYPEIII_SEABEAM // obsolete
	SB_AMP
	SEABAT_II
	SEABAT_8101
	SEABEAM_2112
	ELAC_MKII
	EM3000
	EM1002
	EM300
	CMP_SAAS // CMP (compressed), should be used in place of SASS
	RESON_8101
	RESON_8111
	RESON_8124
	RESON_8125
	RESON_8150
	RESON_8160
	EM120
	EM3002
	EM3000D
	EM3002D
	EM121A_SIS
	EM710
	EM302
	EM122
	GEOSWATH_PLUS
	KLEIN_5410_BSS
	RESON_7125
	EM2000
	EM300_RAW
	EM1002_RAW
	EM2000_RAW
	EM3000_RAW
	EM120_RAW
	EM3002_RAW
	EM3000D_RAW
	EM3002D_RAW
	EM121A_SIS_RAW
	EM2040
	DELTA_T
	R2SONIC_2022
	R2SONIC_2024
	R2SONIC_2020
	SR_NOT_DEFINED // the spec makes no mention of ID 154
	RESON_TSERIES
	KMALL // 156
)

// Single beam sensor subrecord IDs.
const (
	SB_ECHOTRAC SubRecordID = 201 + iota
	SB_BATHY2000
	SB_MGD77
	SB_BDB
	SB_NOSHDB
	SWATH_SB_ECHOTRAC
	SWATH_SB_BATHY2000
	SWATH_SB_MGD77
	SWATH_SB_BDB
	SWATH_SB_NOSHDB
	SWATH_SB_PDD
	SWATH_SB_NAVISOUND // 212
)

// Null values for missing data.
// Defining both float32 and float64 types for now, whereas
// in future only one specific type may be needed.
const (
	NULL_LATITUDE_F32               float32 = 91.0
	NULL_LONGITUDE_F32              float32 = 181.0
	NULL_HEADING_F32                float32 = 361.0
	NULL_COURSE_F32                 float32 = 361.0
	NULL_SPEED_F32                  float32 = 99.0
	NULL_PITCH_F32                  float32 = 99.0
	NULL_ROLL_F32                   float32 = 99.0
	NULL_HEAVE_F32                  float32 = 99.0
	NULL_DRAFT_F32                  float32 = 0.0
	NULL_DEPTH_CORRECTOR_F32        float32 = 99.99
	NULL_TIDE_CORRECTOR_F32         float32 = 99.99
	NULL_SOUND_SPEED_CORRECTION_F32 float32 = 99.99
	NULL_HORIZONTAL_ERROR_F32       float32 = -1.00
	NULL_VERTICAL_ERROR_F32         float32 = -1.00
	NULL_HEIGHT                     float64 = 9999.99
	NULL_SEP                        float64 = 9999.99
	NULL_SEP_UNCERTAINTY            float32 = 0.0
	NULL_GPS_TIDE_CORRECTOR         float64 = 99.99
	NULL_LATITUDE_F64               float64 = 91.0
	NULL_LONGITUDE_F64              float64 = 181.0
	NULL_HEADING_F64                float64 = 361.0
	NULL_COURSE_F64                 float64 = 361.0
	NULL_SPEED_F64                  float64 = 99.0
	NULL_PITCH_F64                  float64 = 99.0
	NULL_ROLL_F64                   float64 = 99.0
	NULL_HEAVE_F64                  float64 = 99.0
	NULL_DRAFT_F64                  float64 = 0.0
	NULL_DEPTH_CORRECTOR_F64        float64 = 99.99
	NULL_TIDE_CORRECTOR_F64         float64 = 99.99
	NULL_SOUND_SPEED_CORRECTION_F64 float64 = 99.99
	NULL_HORIZONTAL_ERROR_F64       float64 = -1.00
	NULL_VERTICAL_ERROR_F64         float64 = -1.00
	NULL_SEP_UNCERTAINTY_F64        float64 = 0.0
)

// Null values for swath bathymetry ping subrecords.
// Defining both float32 and float64 types for now, whereas
// in future only one specific type may be needed.
const (
	NULL_DEPTH_F32              float32 = 0.0
	NULL_ACROSS_TRACK_F32       float32 = 0.0
	NULL_ALONG_TRACK_F32        float32 = 0.0
	NULL_TRAVEL_TIME_F32        float32 = 0.0
	NULL_BEAM_ANGLE_F32         float32 = 0.0
	NULL_MC_AMPLITUDE_F32       float32 = 0.0
	NULL_MR_AMPLITUDE_F32       float32 = 0.0
	NULL_ECHO_WIDTH_F32         float32 = 0.0
	NULL_QUALITY_FACTOR_F32     float32 = 0.0
	NULL_RECEIVE_HEAVE_F32      float32 = 0.0
	NULL_DEPTH_ERROR_F32        float32 = 0.0
	NULL_ACROSS_TRACK_ERROR_F32 float32 = 0.0
	NULL_ALONG_TRACK_ERROR_F32  float32 = 0.0
	NULL_NAP_POS_ERROR_F32      float32 = 0.0
	NULL_FLOAT32_ZERO           float32 = 0.0 // Would NaN be better?
	NULL_DEPTH_F64              float64 = 0.0
	NULL_ACROSS_TRACK_F64       float64 = 0.0
	NULL_ALONG_TRACK_F64        float64 = 0.0
	NULL_TRAVEL_TIME_F64        float64 = 0.0
	NULL_BEAM_ANGLE_F64         float64 = 0.0
	NULL_MC_AMPLITUDE_F64       float64 = 0.0
	NULL_MR_AMPLITUDE_F64       float64 = 0.0
	NULL_ECHO_WIDTH_F64         float64 = 0.0
	NULL_QUALITY_FACTOR_F64     float64 = 0.0
	NULL_RECEIVE_HEAVE_F64      float64 = 0.0
	NULL_DEPTH_ERROR_F64        float64 = 0.0
	NULL_ACROSS_TRACK_ERROR_F64 float64 = 0.0
	NULL_ALONG_TRACK_ERROR_F64  float64 = 0.0
	NULL_NAP_POS_ERROR_F64      float64 = 0.0
	NULL_FLOAT64_ZERO           float64 = 0.0 // Would NaN be better?
	NULL_UINT8_ZERO             uint8   = 0
	NULL_UINT16_ZERO            uint16  = 0
	NULL_UINT32_ZERO            uint32  = 0
	NULL_UINT64_ZERO            uint64  = 0
)

// Field sizes for ping subarrays
const (
	FIELD_SIZE_DEFAULT     uint32 = 0x00
	FIELD_SIZE_ONE         uint32 = 0x10
	FIELD_SIZE_TWO         uint32 = 0x20
	FIELD_SIZE_FOUR        uint32 = 0x40
	BYTES_PER_BEAM_DEFAULT uint32 = 1
	BYTES_PER_BEAM_ONE     uint32 = 1
	BYTES_PER_BEAM_TWO     uint32 = 2
	BYTES_PER_BEAM_FOUR    uint32 = 4
)

// Subrecord labels. Used for defining the output schema
var SubRecordNames = map[SubRecordID]string{
	DEPTH:                  "Z", // 1
	ACROSS_TRACK:           "ACROSS_TRACK",
	ALONG_TRACK:            "ALONG_TRACK",
	TRAVEL_TIME:            "TRAVEL_TIME",
	BEAM_ANGLE:             "BEAM_ANGLE",
	MEAN_CAL_AMPLITUDE:     "MEAN_CAL_AMPLITUDE",
	MEAN_REL_AMPLITUDE:     "MEAN_REL_AMPLITUDE",
	ECHO_WIDTH:             "ECHO_WIDTH",
	QUALITY_FACTOR:         "QUALITY_FACTOR",
	RECEIVE_HEAVE:          "RECEIVE_HEAVE",
	DEPTH_ERROR:            "DEPTH_ERROR",        // obsolete
	ACROSS_TRACK_ERROR:     "ACROSS_TRACK_ERROR", // obsolete
	ALONG_TRACK_ERROR:      "ALONG_TRACK_ERROR",  // obsolete
	NOMINAL_DEPTH:          "NOMINAL_DEPTH",
	QUALITY_FLAGS:          "QUALITY_FLAGS",
	BEAM_FLAGS:             "BEAM_FLAGS",
	SIGNAL_TO_NOISE:        "SIGNAL_TO_NOISE",
	BEAM_ANGLE_FORWARD:     "BEAM_ANGLE_FORWARD",
	VERTICAL_ERROR:         "VERTICAL_ERROR",   // replaces depth error
	HORIZONTAL_ERROR:       "HORIZONTAL_ERROR", // replaces across track error
	INTENSITY_SERIES:       "INTENSITY_SERIES",
	SECTOR_NUMBER:          "SECTOR_NUMBER",
	DETECTION_INFO:         "DETECTION_INFO",
	INCIDENT_BEAM_ADJ:      "INCIDENT_BEAM_ADJ",
	SYSTEM_CLEANING:        "SYSTEM_CLEANING",
	DOPPLER_CORRECTION:     "DOPPLER_CORRECTION",
	SONAR_VERT_UNCERTAINTY: "SONAR_VERT_UNCERTAINTY",
	SONAR_HORZ_UNCERTAINTY: "SONAR_HORZ_UNCERTAINTY",
	DETECTION_WINDOW:       "DETECTION_WINDOW",
	MEAN_ABS_COEF:          "MEAN_ABS_COEF", // 30, general subrecords
	UNKNOWN:                "UNKNOWN",       // 0
	SCALE_FACTORS:          "SCALE_FACTORS", // 100
	// SB_UNKNOWN: "SB_UNKNOWN",  // 0, single_beam
	SEABEAM:            "SEABEAM", // 102, multi beam sensor specific
	EM12:               "EM12",
	EM100:              "EM100",
	EM950:              "EM950",
	EM121A:             "EM121A",
	EM121:              "EM121",
	SASS:               "SASS", // obsolete
	SEAMAP:             "SEAMAP",
	SEABAT:             "SEABAT",
	EM1000:             "EM1000",
	TYPEIII_SEABEAM:    "TYPEIII_SEABEAM", // obsolete
	SB_AMP:             "SB_AMP",
	SEABAT_II:          "SEABAT_II",
	SEABAT_8101:        "SEABAT_8101",
	SEABEAM_2112:       "SEABEAM_2112",
	ELAC_MKII:          "ELAC_MKII",
	EM3000:             "EM3000",
	EM1002:             "EM1002",
	EM300:              "EM300",
	CMP_SAAS:           "CMP_SAAS", // CMP (compressed), should be used in place of SASS
	RESON_8101:         "RESON_8101",
	RESON_8111:         "RESON_8111",
	RESON_8124:         "RESON_8124",
	RESON_8125:         "RESON_8125",
	RESON_8150:         "RESON_8150",
	RESON_8160:         "RESON_8160",
	EM120:              "EM120",
	EM3002:             "EM3002",
	EM3000D:            "EM3000D",
	EM3002D:            "EM3002D",
	EM121A_SIS:         "EM121A_SIS",
	EM710:              "EM710",
	EM302:              "EM302",
	EM122:              "EM122",
	GEOSWATH_PLUS:      "GEOSWATH_PLUS",
	KLEIN_5410_BSS:     "KLEIN_5410_BSS",
	RESON_7125:         "RESON_7125",
	EM2000:             "EM2000",
	EM300_RAW:          "EM300_RAW",
	EM1002_RAW:         "EM1002_RAW",
	EM2000_RAW:         "EM2000_RAW",
	EM3000_RAW:         "EM3000_RAW",
	EM120_RAW:          "EM120_RAW",
	EM3002_RAW:         "EM3002_RAW",
	EM3000D_RAW:        "EM3000D_RAW",
	EM3002D_RAW:        "EM3002D_RAW",
	EM121A_SIS_RAW:     "EM121A_SIS_RAW",
	EM2040:             "EM2040",
	DELTA_T:            "DELTA_T",
	R2SONIC_2022:       "R2SONIC_2022",
	R2SONIC_2024:       "R2SONIC_2024",
	R2SONIC_2020:       "R2SONIC_2020",
	SR_NOT_DEFINED:     "SR_NOT_DEFINED", // the spec makes no mention of ID 154
	RESON_TSERIES:      "RESON_TSERIES",
	KMALL:              "KMALL",       // 156
	SB_ECHOTRAC:        "SB_ECHOTRAC", // 201, single beam sensor specific
	SB_BATHY2000:       "SB_BATHY2000",
	SB_MGD77:           "SB_MGD77",
	SB_BDB:             "SB_BDB",
	SB_NOSHDB:          "SB_NOSHDB",
	SWATH_SB_ECHOTRAC:  "SWATH_SB_ECHOTRAC",
	SWATH_SB_BATHY2000: "SWATH_SB_BATHY2000",
	SWATH_SB_MGD77:     "SWATH_SB_MGD77",
	SWATH_SB_BDB:       "SWATH_SB_BDB",
	SWATH_SB_NOSHDB:    "SWATH_SB_NOSHDB",
	SWATH_SB_PDD:       "SWATH_SB_PDD",
	SWATH_SB_NAVISOUND: "SWATH_SB_NAVISOUND", // 212
}

var InvSubRecordNames = lo.Invert(SubRecordNames)

// Record labels.
var RecordNames = map[RecordID]string{
	HEADER:                 "HEADER", // 1
	SWATH_BATHYMETRY_PING:  "SWATH_BATHYMETRY_PING",
	SOUND_VELOCITY_PROFILE: "SOUND_VELOCITY_PROFILE",
	PROCESSING_PARAMETERS:  "PROCESSING_PARAMETERS",
	SENSOR_PARAMETERS:      "SENSOR_PARAMETERS",
	COMMENT:                "COMMENT",
	HISTORY:                "HISTORY",
	NAVIGATION_ERROR:       "NAVIGATION_ERROR", // obsolete
	SWATH_BATHY_SUMMARY:    "SWATH_BATHY_SUMMARY",
	SINGLE_BEAM_PING:       "SINGLE_BEAM_PING",
	HV_NAVIGATION_ERROR:    "HV_NAVIGATION_ERROR", // replaces navigation error
	ATTITUDE:               "ATTITUDE",            // 12
}

// Schema mapping for beam data; TODO; rework this back and forth mapping
var BeamDataName2SubRecordID = map[string]SubRecordID{
	"Z":                    DEPTH,
	"AcrossTrack":          ACROSS_TRACK,
	"AlongTrack":           ALONG_TRACK,
	"TravelTime":           TRAVEL_TIME,
	"BeamAngle":            BEAM_ANGLE,
	"MeanCalAmplitude":     MEAN_CAL_AMPLITUDE,
	"MeanRelAmplitude":     MEAN_REL_AMPLITUDE,
	"EchoWidth":            ECHO_WIDTH,
	"QualityFactor":        QUALITY_FACTOR,
	"ReceiveHeave":         RECEIVE_HEAVE,
	"DepthError":           DEPTH_ERROR,        // obsolete
	"AcrossTrackError":     ACROSS_TRACK_ERROR, // obsolete
	"AlongTrackError":      ALONG_TRACK_ERROR,  // obsolete
	"NominalDepth":         NOMINAL_DEPTH,
	"QualityFlags":         QUALITY_FLAGS,
	"BeamFlags":            BEAM_FLAGS,
	"SignalToNoise":        SIGNAL_TO_NOISE,
	"BeamAngleForward":     BEAM_ANGLE_FORWARD,
	"VerticalError":        VERTICAL_ERROR,   // replaces depth error
	"HorizontalError":      HORIZONTAL_ERROR, // replaces across track error
	"IntensitySeries":      INTENSITY_SERIES,
	"SectorNumber":         SECTOR_NUMBER,
	"DetectionInfo":        DETECTION_INFO,
	"IncidentBeamAdj":      INCIDENT_BEAM_ADJ,
	"SystemCleaning":       SYSTEM_CLEANING,
	"DopplerCorrection":    DOPPLER_CORRECTION,
	"SonarVertUncertainty": SONAR_VERT_UNCERTAINTY,
	"SonarHorzUncertainty": SONAR_HORZ_UNCERTAINTY,
	"DetectionWindow":      DETECTION_WINDOW,
	"MeanAbsCoef":          MEAN_ABS_COEF, // 30, general subrecords
}
