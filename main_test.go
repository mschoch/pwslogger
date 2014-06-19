package main

import (
	"testing"
)

func TestDewpoint(t *testing.T) {
	tests := []struct {
		temp     float64
		humidity float64
		dewpoint float64
	}{
		{
			temp:     23,
			humidity: 79,
			dewpoint: 19,
		},
		{
			temp:     30,
			humidity: 75,
			dewpoint: 25,
		},
		{
			temp:     10,
			humidity: 50,
			dewpoint: 0,
		},
		{
			temp:     -5,
			humidity: 25,
			dewpoint: -22,
		},
	}

	for _, test := range tests {
		approximateDewpoint := ApproximateDewpoint(test.temp, test.humidity)
		if approximateDewpoint < (test.dewpoint-0.5) || approximateDewpoint > (test.dewpoint+0.5) {
			t.Errorf("expected %f, got %f for %.1f C and %.1f%% RH", test.dewpoint, approximateDewpoint, test.temp, test.humidity)
		}
	}
}
