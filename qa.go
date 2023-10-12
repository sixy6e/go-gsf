package gsf

import (
    "time"
    "github.com/samber/lo"
)

type QualityInfo struct {
    Min_Max_Beams []uint16
    Consistent_Beams bool
    Coincident_Pings bool
    Duplicate_Pings bool
    Duplicates []time.Time
    Consistent_Schema bool
}

func (fi *FileInfo) QInfo() {
    var (
        nbeams []uint16
        timestamps []time.Time
        qa QualityInfo
        // sub_rec_counts_str map[string]uint64
    )

    coincident_pings := false
    dup_pings := false
    npings := len(fi.Ping_Info)
    // sub_rec_counts_str = make(map[string]uint64)

    // there have been instances where the number of beams was inconsistent between pings
    // the general idea is to know whether we're dealing with a consistent number of beams
    nbeams = make([]uint16, npings)

    // duplicate pings. one of the samples we were given had duplicate timestamps
    // UPDATE; the sensor configuration could be dual-head or dual-swath
    // dual-swath: two sensors slightly offset that ping at the time can produce
    // something that looks like duplicate pings with the same timestamp.
    // One of the samples had both records with the same ping counter
    // eg ping 1168, ping 1168, then ping 1170, ping 1170.
    // MBSystem still treated them as ping 1, 2, 3, 4.
    timestamps = make([]time.Time, npings)

    for i, ping := range(fi.Ping_Info) {
        nbeams[i] = ping.Number_Beams
        timestamps[i] = ping.Timestamp
    }

    // domain for number of beams
    max := lo.Max(nbeams)
    min := lo.Min(nbeams)
    min_max_beams := []uint16{min, max}
    consistent_beams := min == max

    // potential duplicates
    // to cater for dual swath; the general logic is (n_pings / 2) != len(duplicates)
    // however, this will still fail to ignore dual swath pings and there is
    // an actual duplicate, meaning dual swaths pings will be reported if the
    // dual acquisition isn't consistent
    // or fail to identify if the dual swath recording of a "ping" is in two
    // different files
    // eg ping 1a in file1, ping1b in file2, which potentially means
    // we won't really know if it is dual swath
    // either way, this results in reporting all duplicate timestamps meaning it is
    // still hard to separate valid from invalid duplicate timestamps
    duplicates := lo.FindDuplicates(timestamps)
    if len(duplicates) > 0 {
        dup_pings = (float32(npings) / float32(2)) != float32(len(duplicates))
    }

    // consistent schema; we've had cases where the schema is inconsistent between pings
    vals := make([]uint64, 0)
    for key, val := range fi.SubRecord_Counts {
        // sub_rec_counts_str[SubRecordNames[key]] = val
        // scale factors are not required to be stored in every ping :(
        if key != "SCALE_FACTORS" {
            vals = append(vals, val)
        }
    }
    set := lo.Union(vals)

    qa.Min_Max_Beams = min_max_beams
    qa.Consistent_Beams = consistent_beams
    qa.Duplicate_Pings = dup_pings
    // qa.Duplicates = duplicates
    qa.Consistent_Schema = len(set) == 1

    if dup_pings {
        qa.Duplicates = duplicates
    } else {
        qa.Duplicates = make([]time.Time, 0)

        // we may have a dual sensor configuration (dual swath, or dual head)
        if len(duplicates) > 0 {
            // qa.Coincident_Pings = true
            coincident_pings = true
        }
    }

    // we may have a dual sensor configuration (dual swath, or dual head)
    // if len(duplicates) > 0 && dup_pings == false {
    //     qa.Coincident_Pings = true
    // }
    qa.Coincident_Pings = coincident_pings

    fi.Metadata.Quality_Info = qa
}
