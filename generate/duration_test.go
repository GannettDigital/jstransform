package generate

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/GannettDigital/msgp/msgp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration_Parse(t *testing.T) {
	tests := []struct {
		description  string
		input        string
		wantDuration time.Duration
		wantErr      bool
	}{
		{
			description:  "week duration",
			input:        "P1W",
			wantDuration: 7 * 24 * time.Hour,
		},
		{
			description:  "multiple weeks",
			input:        "P2W",
			wantDuration: 14 * 24 * time.Hour,
		},
		{
			description:  "single day",
			input:        "P1D",
			wantDuration: 24 * time.Hour,
		},
		{
			description:  "hours only",
			input:        "PT1H",
			wantDuration: time.Hour,
		},
		{
			description:  "minutes only",
			input:        "PT15M",
			wantDuration: 15 * time.Minute,
		},
		{
			description:  "seconds only",
			input:        "PT30S",
			wantDuration: 30 * time.Second,
		},
		{
			description:  "zero duration",
			input:        "PT0S",
			wantDuration: 0,
		},
		{
			description:  "day and time",
			input:        "P1DT12H",
			wantDuration: 36 * time.Hour,
		},
		{
			description:  "combined date and time",
			input:        "P1Y2M3DT4H5M6S",
			wantDuration: (365+60+3)*24*time.Hour + 4*time.Hour + 5*time.Minute + 6*time.Second,
		},
		{
			description:  "negative day duration",
			input:        "-P1D",
			wantDuration: -24 * time.Hour,
		},
		{
			description:  "negative time duration",
			input:        "-PT15M",
			wantDuration: -15 * time.Minute,
		},
		{
			description:  "negative week duration",
			input:        "-P1W",
			wantDuration: -7 * 24 * time.Hour,
		},
		{
			description:  "decimal seconds duration",
			input:        "PT1.5S",
			wantDuration: 1500 * time.Millisecond,
		},
		{
			description:  "sub-second duration",
			input:        "PT0.5S",
			wantDuration: 500 * time.Millisecond,
		},
		{
			description: "empty string returns zero duration",
			input:       "",
		},
		{
			description: "overflow integer component",
			input:       "P999999999999999999999999999999Y",
			wantErr:     true,
		},
		{
			description: "bare P without elements",
			input:       "P",
			wantErr:     true,
		},
		{
			description: "bare PT without time elements",
			input:       "PT",
			wantErr:     true,
		},
		{
			description: "time unit without T",
			input:       "P1H",
			wantErr:     true,
		},
		{
			description: "seconds without T",
			input:       "P30S",
			wantErr:     true,
		},
		{
			description: "wrong element order",
			input:       "P1M2Y",
			wantErr:     true,
		},
		{
			description:  "week combined with days",
			input:        "P1W1D",
			wantDuration: 8 * 24 * time.Hour,
		},
		{
			description: "standard Go duration rejected",
			input:       "1h30m",
			wantErr:     true,
		},
		{
			description: "malformed text",
			input:       "invalid-duration",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.input == "" {
				assert.True(t, got.IsZero())
				return
			}
			assert.Equal(t, tt.wantDuration, got.TimeDuration())
			assert.Equal(t, tt.input, got.String())
		})
	}
}

func TestDuration_JSON(t *testing.T) {
	type wrapper struct {
		Duration    Duration  `json:"duration"`
		OptDuration *Duration `json:"optDuration,omitempty"`
	}

	dur, err := ParseDuration("PT15M")
	require.NoError(t, err)
	w := wrapper{
		Duration:    dur,
		OptDuration: &dur,
	}

	data, err := json.Marshal(w)
	require.NoError(t, err)
	assert.JSONEq(t, `{"duration":"PT15M","optDuration":"PT15M"}`, string(data))

	var unmarshaled wrapper
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, dur.TimeDuration(), unmarshaled.Duration.TimeDuration())
	require.NotNil(t, unmarshaled.OptDuration)
	assert.Equal(t, dur.TimeDuration(), unmarshaled.OptDuration.TimeDuration())

	// Nullable and empty handling
	var emptyWrapper wrapper
	err = json.Unmarshal([]byte(`{"duration":null}`), &emptyWrapper)
	require.NoError(t, err)
	assert.True(t, emptyWrapper.Duration.IsZero())

	err = json.Unmarshal([]byte(`{"duration":"invalid"}`), &emptyWrapper)
	assert.Error(t, err)
}

func TestDuration_Text(t *testing.T) {
	dur, err := ParseDuration("PT1H30M")
	require.NoError(t, err)

	text, err := dur.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "PT1H30M", string(text))

	var parsed Duration
	err = parsed.UnmarshalText(text)
	require.NoError(t, err)
	assert.Equal(t, dur.String(), parsed.String())

	var zero Duration
	zeroText, err := zero.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "PT0S", string(zeroText))

	err = parsed.UnmarshalText([]byte(""))
	require.NoError(t, err)
	assert.True(t, parsed.IsZero())
}

func TestDuration_Msgp(t *testing.T) {
	dur, err := ParseDuration("PT45S")
	require.NoError(t, err)

	// Test MarshalMsg / UnmarshalMsg
	bytesOut, err := dur.MarshalMsg(nil)
	require.NoError(t, err)

	var readDur Duration
	rem, err := readDur.UnmarshalMsg(bytesOut)
	require.NoError(t, err)
	assert.Empty(t, rem)
	assert.Equal(t, dur.String(), readDur.String())

	// Test EncodeMsg / DecodeMsg
	var buf bytes.Buffer
	writer := msgp.NewWriter(&buf)
	err = dur.EncodeMsg(writer)
	require.NoError(t, err)
	err = writer.Flush()
	require.NoError(t, err)

	var decodedDur Duration
	reader := msgp.NewReader(&buf)
	err = decodedDur.DecodeMsg(reader)
	require.NoError(t, err)
	assert.Equal(t, dur.String(), decodedDur.String())

	assert.Positive(t, dur.Msgsize())
}

func TestDurationFromTime(t *testing.T) {
	tests := []struct {
		description string
		input       time.Duration
		wantString  string
	}{
		{
			description: "zero duration",
			input:       0,
			wantString:  "PT0S",
		},
		{
			description: "30 seconds",
			input:       30 * time.Second,
			wantString:  "PT30S",
		},
		{
			description: "15 minutes",
			input:       15 * time.Minute,
			wantString:  "PT15M",
		},
		{
			description: "1 hour 30 minutes",
			input:       1*time.Hour + 30*time.Minute,
			wantString:  "PT1H30M",
		},
		{
			description: "1 day",
			input:       24 * time.Hour,
			wantString:  "P1D",
		},
		{
			description: "1 day 2 hours 3 minutes 4 seconds",
			input:       24*time.Hour + 2*time.Hour + 3*time.Minute + 4*time.Second,
			wantString:  "P1DT2H3M4S",
		},
		{
			description: "1.5 seconds (sub-second precision)",
			input:       1500 * time.Millisecond,
			wantString:  "PT1.5S",
		},
		{
			description: "0.5 seconds (sub-second precision)",
			input:       500 * time.Millisecond,
			wantString:  "PT0.5S",
		},
		{
			description: "negative duration",
			input:       -24 * time.Hour,
			wantString:  "-P1D",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			d := DurationFromTime(tt.input)
			assert.Equal(t, tt.input, d.TimeDuration())
			assert.Equal(t, tt.wantString, d.String())
		})
	}
}
