package gsf

import (
	"math"

	"github.com/samber/lo"
)

// left, right := lo.Difference([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 6})
// []int{1, 3, 4, 5}, []int{6}

// fillNulls is specifically when we're building a chunk of pings together as
// a single cohesive unit. As the schema can be (annoyingly) different between
// SWATH_BATHYMETRY_PING records, we need to fill in the blanks for the attributes
// missing from the newly read ping.
// Assuming that the Lon_lat field won't have nulls, as they're computed upon reading the
// ping. Sensor_metadata and Sensor_imagery_metadata fields are also assumed to be
// always populated with content. Haven't investigated whether the schema of each can
// differ between pings though.
// If issues arise later, or a case of inconsistencies in schemas occur for these
// other fields is found, then something can be done.
func (pd *PingData) fillNulls(singlePing *PingData) error {
	nbeams := singlePing.Ping_headers.Number_beams[0]
	left, _ := lo.Difference(pd.ba_subrecords, singlePing.ba_subrecords)
	_ = pd.beamArrayNulls(uint64(nbeams), left)
	return nil
}

// padDense pads the BeamArray, LatLon and BrbIntensity structures and their
// numerous slices by a given size.
// This padding is only applied when the beam data array is to be output as
// a dense TileDB array.
// Appropriate null values are inserted for each attribute.
// Why is padding required? Because the number of beams for each ping
// can differ from one another. Padding makes the BeamNumber dimensional axis
// consistent across all pings.
func (pd *PingData) padDense(size uint16) error {
	_ = pd.beamArrayNulls(uint64(size), pd.ba_subrecords)

	// longitude and latitude
	for i := uint16(0); i < size; i++ {
		pd.Lon_lat.Longitude = append(pd.Lon_lat.Longitude, NULL_LONGITUDE_F64)
		pd.Lon_lat.Latitude = append(pd.Lon_lat.Latitude, NULL_LATITUDE_F64)
	}
	return nil
}

// beamArrayNulls appends null data, mostly zero defined at varied datatypes
// float32, float64, signed/unsigned integers.
// NaN is used for Brb_intensity.TimeSeries
func (pd *PingData) beamArrayNulls(size uint64, fields []string) error {
	for _, name := range fields {
		subr_id := BeamDataName2SubRecordID[name]

		switch subr_id {
		case DEPTH:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.Z = append(pd.Beam_array.Z, NULL_DEPTH_F64)
			}
		case ACROSS_TRACK:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.AcrossTrack = append(pd.Beam_array.AcrossTrack, NULL_ACROSS_TRACK_F64)
			}
		case ALONG_TRACK:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.AlongTrack = append(pd.Beam_array.AlongTrack, NULL_ALONG_TRACK_F64)
			}
		case TRAVEL_TIME:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.TravelTime = append(pd.Beam_array.TravelTime, NULL_TRAVEL_TIME_F64)
			}
		case BEAM_ANGLE:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.BeamAngle = append(pd.Beam_array.BeamAngle, NULL_BEAM_ANGLE_F32)
			}
		case MEAN_CAL_AMPLITUDE:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.MeanCalAmplitude = append(pd.Beam_array.MeanCalAmplitude, NULL_MC_AMPLITUDE_F32)
			}
		case MEAN_REL_AMPLITUDE:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.MeanRelAmplitude = append(pd.Beam_array.MeanRelAmplitude, NULL_MR_AMPLITUDE_F32)
			}
		case ECHO_WIDTH:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.EchoWidth = append(pd.Beam_array.EchoWidth, NULL_ECHO_WIDTH_F32)
			}
		case QUALITY_FACTOR:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.QualityFactor = append(pd.Beam_array.QualityFactor, NULL_QUALITY_FACTOR_F32)
			}
		case RECEIVE_HEAVE:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.RecieveHeave = append(pd.Beam_array.RecieveHeave, NULL_RECEIVE_HEAVE_F32)
			}
		case DEPTH_ERROR:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.DepthError = append(pd.Beam_array.DepthError, NULL_DEPTH_ERROR_F32)
			}
		case ACROSS_TRACK_ERROR:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.AcrossTrackError = append(pd.Beam_array.AcrossTrackError, NULL_ACROSS_TRACK_ERROR_F32)
			}
		case ALONG_TRACK_ERROR:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.AlongTrackError = append(pd.Beam_array.AlongTrackError, NULL_ALONG_TRACK_ERROR_F32)
			}
		case NOMINAL_DEPTH:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.NominalDepth = append(pd.Beam_array.NominalDepth, NULL_FLOAT64_ZERO)
			}
		case QUALITY_FLAGS:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.QualityFlags = append(pd.Beam_array.QualityFlags, NULL_UINT8_ZERO)
			}
		case BEAM_FLAGS:
			// TODO; look at what an ideal null would be for the beam bitflag
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.BeamFlags = append(pd.Beam_array.BeamFlags, NULL_UINT8_ZERO)
			}
		case SIGNAL_TO_NOISE:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.SignalToNoise = append(pd.Beam_array.SignalToNoise, NULL_FLOAT32_ZERO)
			}
		case BEAM_ANGLE_FORWARD:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.BeamAngleForward = append(pd.Beam_array.BeamAngleForward, NULL_FLOAT32_ZERO)
			}
		case VERTICAL_ERROR:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.VerticalError = append(pd.Beam_array.VerticalError, NULL_FLOAT32_ZERO)
			}
		case HORIZONTAL_ERROR:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.HorizontalError = append(pd.Beam_array.HorizontalError, NULL_FLOAT32_ZERO)
			}
		case INTENSITY_SERIES:
			// each beam has a series of values, so not sure what the best approach is
			// maybe a single value of 0?
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Brb_intensity.TimeSeries = append(pd.Brb_intensity.TimeSeries, math.NaN())
				// pd.Brb_intensity.BottomDetect = append(pd.Brb_intensity.BottomDetect, NULL_FLOAT32_ZERO)
				pd.Brb_intensity.BottomDetectIndex = append(pd.Brb_intensity.BottomDetectIndex, NULL_UINT16_ZERO)
				pd.Brb_intensity.StartRange = append(pd.Brb_intensity.StartRange, NULL_UINT16_ZERO)
				pd.Brb_intensity.sample_count = append(pd.Brb_intensity.sample_count, uint16(0))
			}
		case SECTOR_NUMBER:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.SectorNumber = append(pd.Beam_array.SectorNumber, NULL_UINT16_ZERO)
			}
		case DETECTION_INFO:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.DetectionInfo = append(pd.Beam_array.DetectionInfo, NULL_UINT16_ZERO)
			}
		case INCIDENT_BEAM_ADJ:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.IncidentBeamAdj = append(pd.Beam_array.IncidentBeamAdj, NULL_FLOAT32_ZERO)
			}
		case SYSTEM_CLEANING:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.SystemCleaning = append(pd.Beam_array.SystemCleaning, NULL_UINT16_ZERO)
			}
		case DOPPLER_CORRECTION:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.DopplerCorrection = append(pd.Beam_array.DopplerCorrection, NULL_FLOAT32_ZERO)
			}
		case SONAR_VERT_UNCERTAINTY:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.SonarVertUncertainty = append(pd.Beam_array.SonarVertUncertainty, NULL_FLOAT32_ZERO)
			}
		case SONAR_HORZ_UNCERTAINTY:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.SonarHorzUncertainty = append(pd.Beam_array.SonarHorzUncertainty, NULL_FLOAT32_ZERO)
			}
		case DETECTION_WINDOW:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.DetectionWindow = append(pd.Beam_array.DetectionWindow, NULL_FLOAT64_ZERO)
			}
		case MEAN_ABS_COEF:
			for i := NULL_UINT64_ZERO; i < size; i++ {
				pd.Beam_array.MeanAbsCoef = append(pd.Beam_array.MeanAbsCoef, NULL_FLOAT64_ZERO)
			}

		}
	}
	return nil
}
