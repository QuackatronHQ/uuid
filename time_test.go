package uuid

import (
	"testing"
	"time"
)

func TestGetTime(t *testing.T) {
	now := time.Now()
	tt := map[string]struct {
		input        func() *time.Time
		expectedTime int64
	}{
		"it should return the current time": {
			input: func() *time.Time {
				return nil
			},
			expectedTime: now.Unix(),
		},
		"it should return the provided time": {
			input: func() *time.Time {
				parsed, err := time.Parse(time.RFC3339, "2024-10-15T09:32:23Z")
				if err != nil {
					t.Errorf("timeParse unexpected error: %v", err)
				}
				return &parsed
			},
			expectedTime: 1728984743,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result, _, err := getTime(tc.input())
			if err != nil {
				t.Errorf("getTime unexpected error: %v", err)
			}
			sec, _ := result.UnixTime()
			if sec != tc.expectedTime {
				t.Errorf("expected %v, got %v", tc.expectedTime, result)
			}
		})
	}
}

func TestTimestamp(t *testing.T) {
	when := time.Date(2024, 10, 15, 9, 32, 23, 0, time.UTC)
	v1, err := NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	v6, err := NewV6WithTime(&when)
	if err != nil {
		t.Fatal(err)
	}
	v7, err := NewV7WithTime(when)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name   string
		uuid   UUID
		want   time.Time
		wantOK bool
	}{
		{"v1 is close to now", v1, time.Now(), true},
		{"v6 keeps custom time", v6, when, true},
		{"v7 keeps custom time", v7, when, true},
		{"v4 has no time", Must(NewRandom()), time.Time{}, false},
		{"v5 has no time", NewSHA1(NameSpaceDNS, []byte("example.com")), time.Time{}, false},
		{"nil has no time", Nil, time.Time{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.uuid.Timestamp()
			if ok != tc.wantOK {
				t.Fatalf("Timestamp() ok = %v, want %v", ok, tc.wantOK)
			}
			if !ok {
				return
			}
			if d := got.Sub(tc.want); d < -time.Second || d > time.Second {
				t.Errorf("Timestamp() = %v, want within 1s of %v", got.UTC(), tc.want.UTC())
			}
		})
	}
}
