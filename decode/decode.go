//
package decode

import (
    "time"
)

type RecordID uint16
type SubRecordID uint16
type SensorID uint16

const (
    NEXT_RECORD RecordID = 0
    BEAM_WIDTH_UNKNOWN float32 = -1.0
    SCALE1 float64 = 10_000_000
    SCALE2 float32 = 100
    SCALE3 float32 = 1000
    SCALE4 int = 1_000_000
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
    NAVIGATION_ERROR  // obselete
    SWATH_BATHY_SUMMARY
    SINGLE_BEAM_PING
    HV_NAVIGATION_ERROR  // replaces navigation error
    ATTITUDE  // 12
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
    QUALITY_FACTOR
    RECEIVE_HEAVE
    DEPTH_ERROR  // obselete
    ACROSS_TRACK_ERROR // obselete
    ALONG_TRACK_ERROR // obselete
    NOMINAL_DEPTH
    QUALITY_FLAGS
    BEAM_FLAGS
    SIGNAL_TO_NOISE
    BEAM_ANGLE_FORWARD
    VERTICAL_ERROR  // replaces depth error
    HORIZONTAL_ERROR  // replaces across track error
    INTENSITY_SERIES
    SECTOR_NUMBER
    DETECTION_INFO
    INCIDENT_BEAM_ADJ
    SYSTEM_CLEANING
    DOPPLER_CORRECTION
    SONAR_VERT_UNCERTAINTY
    SONAR_HORZ_UNCERTAINTY
    DETECTION_WINDOW
    MEAN_ABS_COEF  // 30
)

// General subrecord IDs.
const (
    UNKNOWN SubRecordID = 0
    SCALE_FACTORS SubRecordID = 100
    SB_UNKNOWN SubRecordID = 0  // single_beam
)

// Additional swath bathymetry subrecord IDs; Sensor specific records.
// The scale factors contained within the SCALE_FACTORS subrecord do not apply here.
// (Possibly in relation to the intensity data).
const (
    SEABEAM SensorID = 102 + iota
    EM12
    EM100
    EM950
    EM121A
    EM121
    SASS  // obselete
    SEAMAP
    SEABAT
    EM1000
    TYPEIII_SEABEAM  // obselete
    SB_AMP
    SEABAT_II
    SEABAT_8101
    SEABEAM_2112
    ELAC_MKII
    EM3000
    EM1002
    EM300
    CMP_SAAS  // CMP (compressed), should be used in place of SASS
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
    EM121a_SIS
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
    RESON_TSERIES
    KMALL  // 156
)

// Single beam sensor subrecord IDs.
const (
    ECHOTRAC SensorID = 201 + iota
    BATHY2000
    MGD77
    BDB
    NOSHDB
    SWATH_ECHOTRAC
    SWATH_BATHY2000
    SWATH_MGD77
    SWATH_BDB
    SWATH_NOSHDB
    SWATH_PDD
    SWATH_NAVISOUND  // 212
)

// Null values for missing data.
const (
    NULL_LATITUDE float32 = 91.0
    NULL_LONGITUDE float32 = 181.0
    NULL_HEADING float32 = 361.0
    NULL_COURSE float32 = 361.0
    NULL_SPEED float32 = 99.0
    NULL_PITCH float32 = 99.0
    NULL_ROLL float32 = 99.0
    NULL_HEAVE float32 = 99.0
    NULL_DRAFT float32 = 0.0
    NULL_DEPTH_CORRECTOR float32 = 99.99
    NULL_TIDE_CORRECTOR float32 = 99.99
    NULL_SOUND_SPEED_CORRECTION float32 = 99.99
    NULL_HORIZONTAL_ERROR float32 = -1.00
    NULL_VERTICAL_ERROR float32 = -1.00
    NULL_HEIGHT float32 = 9999.99
    NULL_SEP float32 = 9999.99
    NULL_SEP_UNCERTAINTY float32 = 0.0
)

// Null values for swath bathymetry ping subrecords.
const (
    NULL_DEPTH float32 = 0.0
    NULL_ACROSS_TRACK float32 = 0.0
    NULL_ALONG_TRACK float32 = 0.0
    NULL_TRAVEL_TIME float32 = 0.0
    NULL_BEAM_ANGLE float32 = 0.0
    NULL_MC_AMPLITUDE float32 = 0.0
    NULL_MR_AMPLITUDE float32 = 0.0
    NULL_ECHO_WIDTH float32 = 0.0
    NULL_QUALITY_FACTOR float32 = 0.0
    NULL_RECEIVE_HEAVE float32 = 0.0
    NULL_DEPTH_ERROR float32 = 0.0
    NULL_ACROSS_TRACK_ERROR float32 = 0.0
    NULL_ALONG_TRACK_ERROR float32 = 0.0
    NULL_NAP_POS_ERROR float32 = 0.0
)

// Subrecord labels. Used for defining the output schema
var SubRecordNames = map[SubRecordID]string{
    DEPTH: "DEPTH",  // 1
    ACROSS_TRACK: "ACROSS_TRACK",
    ALONG_TRACK: "ALONG_TRACK",
    TRAVEL_TIME: "TRAVEL_TIME",
    BEAM_ANGLE: "BEAM_ANGLE",
    MEAN_CAL_AMPLITUDE: "MEAN_CAL_AMPLITUDE",
    MEAN_REL_AMPLITUDE: "MEAN_REL_AMPLITUDE",
    ECHO_WIDTH: "ECHO_WIDTH",
    QUALITY_FACTOR: "QUALITY_FACTOR",
    RECEIVE_HEAVE: "RECEIVE_HEAVE",
    DEPTH_ERROR: "DEPTH_ERROR",  // obselete
    ACROSS_TRACK_ERROR: "ACROSS_TRACK_ERROR", // obselete
    ALONG_TRACK_ERROR: "ALONG_TRACK_ERROR", // obselete
    NOMINAL_DEPTH: "NOMINAL_DEPTH",
    QUALITY_FLAGS: "QUALITY_FLAGS",
    BEAM_FLAGS: "BEAM_FLAGS",
    SIGNAL_TO_NOISE: "SIGNAL_TO_NOISE",
    BEAM_ANGLE_FORWARD: "BEAM_ANGLE_FORWARD",
    VERTICAL_ERROR: "VERTICAL_ERROR",  // replaces depth error
    HORIZONTAL_ERROR: "HORIZONTAL_ERROR", // replaces across track error
    INTENSITY_SERIES: "INTENSITY_SERIES",
    SECTOR_NUMBER: "SECTOR_NUMBER",
    DETECTION_INFO: "DETECTION_INFO",
    INCIDENT_BEAM_ADJ: "INCIDENT_BEAM_ADJ",
    SYSTEM_CLEANING: "SYSTEM_CLEANING",
    DOPPLER_CORRECTION: "DOPPLER_CORRECTION",
    SONAR_VERT_UNCERTAINTY: "SONAR_VERT_UNCERTAINTY",
    SONAR_HORZ_UNCERTAINTY: "SONAR_HORZ_UNCERTAINTY",
    DETECTION_WINDOW: "DETECTION_WINDOW",
    MEAN_ABS_COEF: "MEAN_ABS_COEF", // 30
}
