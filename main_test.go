package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_convertToTimeLog(t *testing.T) {
	tests := []struct {
		name      string
		inputTime string
		want      time.Duration
		wantErr   bool
	}{
		{name: "1h -> 1 hour", inputTime: "1h", want: time.Hour},
		{name: "1 -> 1 hour", inputTime: "1h", want: time.Hour},
		{name: "60m -> 60 minutes", inputTime: "60m", want: time.Hour},
		{name: "30m -> 30 minutes", inputTime: "30m", want: 30 * time.Minute},
		{name: "1d -> error", inputTime: "1d", wantErr: true}, // this is not supported yet
		{name: "ahaha -> error", inputTime: "ahaha", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToTimeLog(tt.inputTime)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_convertToDay(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{name: "empty", input: "", want: time.Now()},
		{name: "today", input: "today", want: time.Now()},
		{name: "yesterday", input: "yesterday", want: time.Now().Add(-24 * time.Hour)},
		{name: "monday", input: "monday", want: time.Now().Add(daysDiff(time.Monday, time.Now().Weekday()))},
		{name: "tuesday", input: "tuesday", want: time.Now().Add(daysDiff(time.Tuesday, time.Now().Weekday()))},
		{name: "wednesday", input: "wednesday", want: time.Now().Add(daysDiff(time.Wednesday, time.Now().Weekday()))},
		{name: "thursday", input: "thursday", want: time.Now().Add(daysDiff(time.Thursday, time.Now().Weekday()))},
		{name: "friday", input: "friday", want: time.Now().Add(daysDiff(time.Friday, time.Now().Weekday()))},
		{name: "saturday", input: "saturday", want: time.Now().Add(daysDiff(time.Saturday, time.Now().Weekday()))},
		{name: "sunday", input: "sunday", want: time.Now().Add(daysDiff(time.Weekday(7), time.Now().Weekday()))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToDay(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.Truncate(24*time.Hour).UTC(), got)
		})
	}
}

func daysDiff(a, b time.Weekday) time.Duration {
	diff := a - b
	return time.Duration(diff*24) * time.Hour
}
