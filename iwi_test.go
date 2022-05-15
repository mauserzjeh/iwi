package iwi

import (
	"bytes"
	"io/ioutil"
	"testing"
)

// assertEqual fails if the two values are not equal
func assertEqual[T comparable](t testing.TB, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got: %v != want: %v", got, want)
	}
}

func TestReadIWI(t *testing.T) {

	t.Run("CoD2", func(t *testing.T) {
		iwi, err := ReadIWI("testdata/cod2.iwi")
		assertEqual(t, err == nil, true)

		data, err := ioutil.ReadFile("testdata/cod2iwi.data")
		assertEqual(t, err == nil, true)

		assertEqual(t, bytes.Compare(iwi.Data, data), 0)
	})
}
