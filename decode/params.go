package decode

import (
    "os"
    "bytes"
    "encoding/binary"
    "time"
    "strings"
    "strconv"

    "github.com/soniakeys/meeus/v3/julian"
)

type params_base struct {
    Seconds int32
    Nano_seconds int32
    N_params int16
}

func parse_reftime(date_str string) time.Time {
    // format is (according to spec) yyyy/ddd hh:mm:ss (eg 1970/001 00:00:00)
    split := strings.Split(date_str, " ")
    split2 := strings.Split(split[0], "/")

    year, _ := strconv.Atoi(split2[0])
    doy, _ := strconv.Atoi(split2[1])
    month, day := julian.DayOfYearToCalendar(doy, julian.LeapYearGregorian(year))

    // hour, min, sec
    split3 := strings.Split(split[1], ":")
    hms := make([]int, len(split3))

    for i, val := range split3 {
        hms[i], _ = strconv.Atoi(val)
    }

    date := time.Date(year, time.Month(month), day, hms[0], hms[1], hms[2], 0, time.UTC)

    return date
}

// ProcessingParametersRec decodes the PROCESSING_PARAMETERS record.
// It contains important scalar or vector values that describe the overall survey
// conditions or operational values.
// Typical parameters include items uch as the navigation sensor's antenna location or the
// reference ellipsoid for the geographic position.
// This record could contain pretty much anything, of any type. We'll try to detect
// as many types as possible and convert them from strings.
func ProcessingParametersRec(stream *os.File, rec Record) map[string]interface{} {
    var (
        param_size int16
        param string
        split []string
        key string
        val string
        svals []string
        base params_base
        i int16
        j int64
    )

    // some fields contain a mix of strings that imply a boolean condition
    bools := map[string]bool{
        "yes": true,
        "no": false,
        "true": true,
        "false": false,
    }

    // standardise the spelling for consistency (potentially change the type to nil)
    unkn := map[string]string{
        "unknwn": "unknown",
        "unknown": "unknown",
    }

    params := make(map[string]interface{})

    buffer := make([]byte, rec.Datasize)

    _ = binary.Read(stream, binary.BigEndian, &buffer)
    reader := bytes.NewReader(buffer)
    _ = binary.Read(reader, binary.BigEndian, &base)

    start_idx, end_idx := 10, 12 // the var `base` contains the first 10 bytes read

    // params are deciphered by an int16 indicating the length of the string param value
    // and the param value containing "=" eg "22APPLIED_ROLL_BIAS=0.03" where 22 is string length
    // rather than retaing the raw string, parse all values to proper types
    // with the intent on outputing the data as a json doc
    for i = 0; i < base.N_params; i++ {

        // size of param (length of string)
        param_size = int16(binary.BigEndian.Uint16(buffer[start_idx:end_idx]))
        start_idx += 2
        end_idx += int(param_size)

        // the param string ("key=value")
        param = string(buffer[start_idx:end_idx])
        start_idx += int(param_size)
        end_idx += 2

        // establish key and value; standardise keys (remove spaces, lowercase); strip chars
        split = strings.Split(strings.TrimSpace(param), "=")
        key = strings.ReplaceAll(strings.ToLower(split[0]), " ", "_")
        val = strings.Trim(strings.ToLower(split[1]), "\x00")

        // this whole next section is a slightly complicated mess of unwrapping ...
        // TODO; define a cleaner & intelligent method than this brute force approach
        if strings.Contains(val, ",") == true {  // ',' implies an array of data
            svals = strings.Split(val, ",")
            length := int64(len(svals))

            if strings.Contains(val, ".") == true {  // assumption on period being a decimal point
                fvals := make([]float32, length)
                for j = 0; j < length; j++ {
                    fval, err := strconv.ParseFloat(svals[j], 32)
                    if err != nil {
                        panic(err)  // something bad, and we want to know why
                    } else {
                        fvals[j] = float32(fval)
                    }
                }
                params[key] = fvals
            } else {  // could be dealing with an array of unknwn or unknown
                for j = 0; j < length; j++ {
                    svals[j] = "unknown"
                }
            }
        } else if strings.Contains(val, ".") == true {  // again, assume float
            fval, err := strconv.ParseFloat(val, 32)
            if err != nil {
                panic(err)
            } else {
                params[key] = float32(fval)
            }
        } else if _, exists := bools[val]; exists {  // convert to bool
            params[key] = bools[val]
        } else if _, exists := unkn[val]; exists {  // unknwn to unknown
            params[key] = unkn[val]
        } else if key == "reference_time" {
            params[key] = parse_reftime(val)
        } else {  // most likely an integer or generic string
            ival, err := strconv.Atoi(val)
            if err != nil {
                params[key] = val  // string
            } else {
                params[key] = ival
            }
        }
    }

    // add the processed time (additional field not defined in the GSF spec)
    params["processed_time"] = time.Unix(int64(base.Seconds), int64(base.Nano_seconds)).UTC()

    return params
}
