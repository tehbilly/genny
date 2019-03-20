package bugreports

import "github.com/mauricelam/genny/generic"

// TA will be replaced in tests
type TA generic.Type

// Receiver is the struct used for tests.
type Receiver struct{}

//genny:start

// TB will be replaced in tests
type TB generic.Type

// TAsToTBs converts a []TA to a []TB
func (Receiver) TAsToTBs(ta []TA) []TB {
    // returning an empty TB slice is sufficient for this test.
    return []TB{}
}
