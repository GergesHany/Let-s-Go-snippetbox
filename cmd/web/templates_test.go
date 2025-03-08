package main

import (
	"snippetbox.alexedwards.net/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2021, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2021 at 10:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2021, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 60*60)),
			want: "17 Dec 2021 at 10:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}

}
