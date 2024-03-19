package gsf

import (
	"github.com/samber/lo"
)

// left, right := lo.Difference([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 6})
// []int{1, 3, 4, 5}, []int{6}

func (pd *PingData) fillNulls(singlePing *PingData) error {
	left, _ := lo.Difference(pd.ba_subrecords, singlePing.ba_subrecords)
	for _, name := range left {
		subr_id := InvSubRecordNames[name]

		switch subr_id {
		case DEPTH:

		}
	}
	return nil
}
