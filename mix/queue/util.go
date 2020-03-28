package queue

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/hkdf"
)

const maxIntUint64 = uint64(^uint(0) >> 1)

func uint64ToInt(l uint64) int {
	// Overflow mitigation.
	if l > maxIntUint64 {
		return int(maxIntUint64)
	}
	return int(l)
}

var hkdfInfo = []byte("HKDF Random Source Info")

var unixTimeSource = unixTimeSourceF
var globalRandomSource = rand.Reader

func unixTimeSourceF() uint64 { return uint64(time.Now().Unix()) }

type hkdfReseedSource struct {
	parentSource    io.Reader
	expansionSource io.Reader
}

func (source *hkdfReseedSource) init() (*hkdfReseedSource, error) {
	pseudorandomKey := make([]byte, 64)
	if n, err := io.ReadFull(source.parentSource, pseudorandomKey); err != nil {
		return source, err
	} else if n != 64 {
		return source, io.ErrShortWrite
	}
	source.expansionSource = hkdf.Expand(sha512.New, pseudorandomKey, hkdfInfo)
	return source, nil
}

func (source *hkdfReseedSource) Read(p []byte) (n int, err error) {
	if n, err = source.expansionSource.Read(p); err != nil {
		source.init()
		n, err = source.expansionSource.Read(p)
	}
	return
}

// randomSourceHKDF returns a random source based on a HKDF-SHA512 with initial seeding from system randomsource if randomSource is nil.
func randomSourceHKDF(randomSource io.Reader) (io.Reader, error) {
	if randomSource == nil {
		randomSource = globalRandomSource
	}
	return (&hkdfReseedSource{
		parentSource: randomSource,
	}).init()
}

// randomUint64 returns a random uint64 read from randomSource.
func randomUint64(randomSource io.Reader) uint64 {
	q := make([]byte, 8)
	if _, err := io.ReadFull(randomSource, q); err != nil {
		panic(fmt.Sprintf("mix/queue/randomUint64: %s", err))
	}
	return binary.BigEndian.Uint64(q)
}

// randomUint64Max returns 0 <= random < max. Calling with max==0, 0 may be returned.
func randomUint64Max(randomSource io.Reader, max uint64) uint64 {
	if max == 0 {
		return 0
	}
	return (randomUint64(randomSource) % max)
}

// randomPermutation returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n].
func randomPermutation(randomSource io.Reader, n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := randomUint64Max(randomSource, uint64(i+1))
		m[i] = m[j]
		m[j] = i
	}
	return m
}
