package queue

import (
	"crypto/rand"
	"testing"
)

func TestRunable(t *testing.T) {
	if r, err := randomSourceHKDF(nil); err != nil {
		t.Errorf("randomSourceHKDF(nil): %s", err)
	} else if r == nil {
		t.Error("randomSourceHKDF(nil): returns nil")
	}
	if r, err := randomSourceHKDF(rand.Reader); err != nil {
		t.Errorf("randomSourceHKDF(rand.Reader): %s", err)
	} else if r == nil {
		t.Error("randomSourceHKDF(rand.Reader): returns nil")
	}
	// ----
	a, b, c := randomUint64(rand.Reader), randomUint64(rand.Reader), randomUint64(rand.Reader)
	if a == b && b == c && c == a { // 1:2^67
		t.Error("randomUint64 not random")
	}
	// ----
	a, b, c = randomUint64Max(rand.Reader, 10000), randomUint64Max(rand.Reader, 10000), randomUint64Max(rand.Reader, 10000)
	if a == b && b == c && c == a {
		t.Error("randomUint64Max not random")
	}
	for i := 0; i < 100; i++ {
		j := randomUint64Max(rand.Reader, uint64(i))
		if j >= uint64(i) && i != 0 {
			t.Errorf("randomUint64Max outside boundary: %d %d", i, j)
		}
	}
	// ----
	testNum := 1000
	if len(randomPermutation(rand.Reader, testNum)) != testNum {
		t.Error("randomPermutation wrong number of elements")
	}
	d, e := randomPermutation(rand.Reader, testNum), randomPermutation(rand.Reader, testNum)
	same := 0
	for i, l := range d {
		if e[i] == l {
			same++
		}
	}
	if same > 900 {
		t.Error("Low randomization")
	}
}
