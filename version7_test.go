// Copyright 2024 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"bytes"
	"errors"
	"sync"
	"testing"
	"time"
)

// uuid v7 time is only unix milliseconds, so
// uuid1.Time() == uuid2.Time() is right, but uuid1 must != uuid2
func TestVersion7(t *testing.T) {
	SetRand(nil)
	m := make(map[string]bool)
	for x := 1; x < 128; x++ {
		uuid, err := NewV7()
		if err != nil {
			t.Fatalf("could not create UUID: %v", err)
		}
		s := uuid.String()
		if m[s] {
			t.Errorf("NewV7 returned duplicated UUID %s", s)
		}
		m[s] = true
		if v := uuid.Version(); v != 7 {
			t.Errorf("UUID of version %s", v)
		}
		if uuid.Variant() != RFC4122 {
			t.Errorf("UUID is variant %d", uuid.Variant())
		}
	}
}

// uuid v7 time is only unix milliseconds, so
// uuid1.Time() == uuid2.Time() is right, but uuid1 must != uuid2
func TestVersion7_pooled(t *testing.T) {
	SetRand(nil)
	EnableRandPool()
	defer DisableRandPool()

	m := make(map[string]bool)
	for x := 1; x < 128; x++ {
		uuid, err := NewV7()
		if err != nil {
			t.Fatalf("could not create UUID: %v", err)
		}
		s := uuid.String()
		if m[s] {
			t.Errorf("NewV7 returned duplicated UUID %s", s)
		}
		m[s] = true
		if v := uuid.Version(); v != 7 {
			t.Errorf("UUID of version %s", v)
		}
		if uuid.Variant() != RFC4122 {
			t.Errorf("UUID is variant %d", uuid.Variant())
		}
	}
}

func TestVersion7FromReader(t *testing.T) {
	myString := "8059ddhdle77cb52"
	r := bytes.NewReader([]byte(myString))
	_, err := NewV7FromReader(r)
	if err != nil {
		t.Errorf("failed generating UUID from a reader")
	}
	_, err = NewV7FromReader(r)
	if err == nil {
		t.Errorf("expecting an error as reader has no more bytes. Got uuid. NewV7FromReader may not be using the provided reader")
	}
}

func TestVersion7Monotonicity(t *testing.T) {
	length := 10000
	u1 := Must(NewV7()).String()
	for i := 0; i < length; i++ {
		u2 := Must(NewV7()).String()
		if u2 <= u1 {
			t.Errorf("monotonicity failed at #%d: %s(next) < %s(before)", i, u2, u1)
			break
		}
		u1 = u2
	}
}

type fakeRand struct{}

func (g fakeRand) Read(bs []byte) (int, error) {
	for i, _ := range bs {
		bs[i] = 0x88
	}
	return len(bs), nil
}

func TestVersion7MonotonicityStrict(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2008, 8, 8, 8, 8, 8, 8, time.UTC)
	}
	defer func() {
		timeNow = time.Now
	}()

	SetRand(fakeRand{})
	defer SetRand(nil)

	length := 100000 // > 3906
	u1 := Must(NewV7())
	for i := 0; i < length; i++ {
		u2 := Must(NewV7())
		if Compare(u1, u2) >= 0 {
			t.Errorf("monotonicity failed at #%d: %s(next) < %s(before)", i, u2, u1)
			break
		}
		u1 = u2
	}
}

func TestNewV7WithTime(t *testing.T) {
	when := time.Date(2024, 10, 15, 9, 32, 23, 123456789, time.UTC)
	u, err := NewV7WithTime(when)
	if err != nil {
		t.Fatal(err)
	}
	if v := u.Version(); v != 7 {
		t.Errorf("version = %s, want 7", v)
	}
	if u.Variant() != RFC4122 {
		t.Errorf("variant = %s, want RFC4122", u.Variant())
	}
	got, ok := u.Timestamp()
	if !ok {
		t.Fatal("Timestamp() returned ok = false for a v7 UUID")
	}
	if want := when.Truncate(time.Millisecond); !got.Equal(want) {
		t.Errorf("Timestamp() = %v, want %v", got.UTC(), want)
	}
}

func TestNewV7WithTimeSortsByTime(t *testing.T) {
	base := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	prev := Must(NewV7WithTime(base))
	for i := 1; i < 1000; i++ {
		next := Must(NewV7WithTime(base.Add(time.Duration(i) * time.Millisecond)))
		if Compare(prev, next) >= 0 {
			t.Fatalf("#%d: %s does not sort after %s", i, next, prev)
		}
		prev = next
	}
}

func TestNewV7WithTimeBeforeEpoch(t *testing.T) {
	_, err := NewV7WithTime(time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC))
	if !errors.Is(err, ErrInvalidV7Time) {
		t.Errorf("error = %v, want ErrInvalidV7Time", err)
	}
}

func TestNewV7WithTimeDoesNotMoveSequence(t *testing.T) {
	before := Must(NewV7())
	// A time far in the future must not push NewV7 forward.
	if _, err := NewV7WithTime(time.Now().Add(24 * time.Hour)); err != nil {
		t.Fatal(err)
	}
	after := Must(NewV7())
	ts, _ := after.Timestamp()
	if time.Until(ts) > time.Minute {
		t.Errorf("NewV7 timestamp jumped to %v after NewV7WithTime", ts)
	}
	if Compare(before, after) >= 0 {
		t.Errorf("monotonicity failed: %s(next) <= %s(before)", after, before)
	}
}

func TestVersion7Concurrent(t *testing.T) {
	const workers, perWorker = 8, 2000
	var mu sync.Mutex
	seen := make(map[UUID]bool, workers*perWorker)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			local := make([]UUID, 0, perWorker)
			for i := 0; i < perWorker; i++ {
				local = append(local, Must(NewV7()))
			}
			mu.Lock()
			defer mu.Unlock()
			for _, u := range local {
				if seen[u] {
					t.Errorf("duplicate UUID %s", u)
				}
				seen[u] = true
			}
		}()
	}
	wg.Wait()
}

func BenchmarkNewV7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Must(NewV7())
	}
}

func BenchmarkNewV7WithTime(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		Must(NewV7WithTime(now))
	}
}
