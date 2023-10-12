package gsf

import (
    // "os"
    "bytes"
    "encoding/binary"
    "time"
)

type PingHeader struct {
    Timestamp time.Time
    Longitude float64
    Latitude float64
    Number_beams uint16
    Centre_beam uint16
    Tide_corrector float32
    Depth_corrector float32
    Heading float32
    Pitch float32
    Roll float32
    Heave float32
    Course float32
    Speed float32
    Height float32
    Separation float32
    GPS_tide_corrector float32
    Ping_flags int16
}

type SubRecord struct {
    Id SubRecordID
    Datasize uint32
    Byte_index int64
}

type ScaleFactor struct {
    Id SubRecordID
    Scale float32  // TODO float32?
    Offset float32
    Compression_flag bool  // if true, then the associated array is compressed
}

// PingInfo contains some basic information regarding the ping such as
// the number of beams, what sub-records are populated.
// The initial reasoning behind why, is to provide a basic descriptor
// to inform a global schema across all pings, and derive max(n_beams) to
// inform a global [ping, beam] dimensional array structure.
type PingInfo struct {
    Timestamp time.Time
    Number_Beams uint16
    Sub_Records []SubRecordID
    Scale_Factors bool
    scale_factors map[SubRecordID]ScaleFactor
}

func decode_ping_hdr(reader *bytes.Reader) PingHeader {
    var (
        hdr_base struct {
            Seconds int32
            Nano_seconds int32
            Longitude int32
            Latitude int32
            Number_beams uint16
            Centre_beam uint16
            Ping_flags int16
            Reserved int16
            Tide_corrector int16
            Depth_corrector int32
            Heading uint16
            Pitch int16
            Roll int16
            Heave int16
            Course uint16
            Speed uint16
            Height int32
            Separation int32
            GPS_tide_corrector int32
            Spare int16
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

    subrecord_id := (int(subrecord_hdr) & 0xFF000000) >> 24  // TODO; define a const as int64
    subrecord_size := int(subrecord_hdr) & 0x00FFFFFF  // TODO; define a const as int64

    byte_index := offset + 4

    subhdr := SubRecord{SubRecordID(subrecord_id), uint32(subrecord_size), byte_index} // include a byte_index??

    return subhdr
}

func scale_factors_rec(reader *bytes.Reader) (scale_factors map[SubRecordID]ScaleFactor, nbytes int64) {
    var (
        i int32
        num_factors int32
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

        subid := (int64(data[0]) & 0xFF000000) >> 24 // TODO; define const for 0xFF000000
        comp_flag := (data[0] & 0x00FF0000) >> 16 == 1 // TODO; define const for 0x00FF0000
        cnvrt_subid := SubRecordID(subid)

        scale_factor = ScaleFactor{
            Id: cnvrt_subid,
            Scale: float32(data[1]),
            Offset: float32(data[2]),
            Compression_flag: comp_flag,  // TODO; implement compression decoder
        }

        nbytes += 12

        // scale_factors[SubRecordID(subid)] = scale_factor
        scale_factors[cnvrt_subid] = scale_factor
    }

    return scale_factors, nbytes
}

func ping_info(reader *bytes.Reader, rec RecordHdr) PingInfo {
    var (
        idx int64 = 0
        pinfo PingInfo
        records = make([]SubRecordID, 0, 32)
        sf bool = false
        scl_fac map[SubRecordID]ScaleFactor
        nbytes int64
    )

    datasize := int64(rec.Datasize)

    hdr := decode_ping_hdr(reader)
    idx += 56 // 56 bytes read for ping header
    offset := rec.Byte_index + idx

    // read through each subrecord
    for (datasize - idx) > 4 {
        sub_rec := SubRecHdr(reader, offset)
        srec_dsize := int64(sub_rec.Datasize)
        idx += 4  // bytes read from header

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
func SwathBathymetryPingRec(buffer []byte, rec RecordHdr) PingHeader {
    var (
        idx int64 = 0
        // sf map[SubRecordID]ScaleFactor
        nbytes int64
        // subrecord_hdr int32
    )

    // buffer := make([]byte, rec.Datasize)
    // _ = binary.Read(stream, binary.BigEndian, &buffer)
    reader := bytes.NewReader(buffer)

    hdr := decode_ping_hdr(reader)
    idx += 56 // 56 bytes read for ping header
    offset := rec.Byte_index + idx

    // first subrecord
    //reader := bytes.NewReader(buffer[56:])
    _, _ = reader.Seek(idx, 0)
    // _ = binary.Read(reader, binary.BigEndian, &subrecord_hdr)
    // subrecord_id := (int(subrecord_hdr) & 0xFF000000) >> 24
    // subrecord_size := int(subrecord_hdr) & 0x00FFFFFF
    sub_rec := SubRecHdr(reader, offset)
    idx += 4

    // case switching; SCALE_FACTORS == 100
    // if scale factor else get scale factor
    if sub_rec.Id == SCALE_FACTORS {
        // read and structure the scale factors
        // however, we'll rely on the scale factors from PingInfo.scale_factors
        _, nbytes = scale_factors_rec(reader)
        idx += nbytes
    }

    return hdr
}
