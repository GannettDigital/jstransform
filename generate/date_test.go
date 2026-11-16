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

func TestDate_Parse(t *testing.T) {
	tests := []struct {
		description string
		input       string
		wantYear    int
		wantMonth   time.Month
		wantDay     int
		wantErr     bool
	}{
		{
			description: "valid date",
			input:       "2026-11-13",
			wantYear:    2026,
			wantMonth:   time.November,
			wantDay:     13,
		},
		{
			description: "leap year date",
			input:       "2024-02-29",
			wantYear:    2024,
			wantMonth:   time.February,
			wantDay:     29,
		},
		{
			description: "empty string returns zero date",
			input:       "",
		},
		{
			description: "invalid month",
			input:       "2026-13-01",
			wantErr:     true,
		},
		{
			description: "invalid day",
			input:       "2026-02-30",
			wantErr:     true,
		},
		{
			description: "date with time",
			input:       "2026-11-13T12:00:00Z",
			wantErr:     true,
		},
		{
			description: "malformed text",
			input:       "not-a-date",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.input == "" {
				assert.True(t, got.IsZero())
				return
			}
			assert.Equal(t, tt.wantYear, got.Year())
			assert.Equal(t, tt.wantMonth, got.Month())
			assert.Equal(t, tt.wantDay, got.Day())
			assert.Equal(t, tt.input, got.String())
		})
	}
}

func TestDate_JSON(t *testing.T) {
	type wrapper struct {
		EventDate Date  `json:"eventDate"`
		OptDate   *Date `json:"optDate,omitempty"`
	}

	date := NewDate(2026, time.November, 13)
	w := wrapper{
		EventDate: date,
		OptDate:   &date,
	}

	data, err := json.Marshal(w)
	require.NoError(t, err)
	assert.JSONEq(t, `{"eventDate":"2026-11-13","optDate":"2026-11-13"}`, string(data))

	var unmarshaled wrapper
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, date.Year(), unmarshaled.EventDate.Year())
	assert.Equal(t, date.Month(), unmarshaled.EventDate.Month())
	assert.Equal(t, date.Day(), unmarshaled.EventDate.Day())
	require.NotNil(t, unmarshaled.OptDate)
	assert.Equal(t, date.Day(), unmarshaled.OptDate.Day())

	// Nullable and empty handling
	var emptyWrapper wrapper
	err = json.Unmarshal([]byte(`{"eventDate":null}`), &emptyWrapper)
	require.NoError(t, err)
	assert.True(t, emptyWrapper.EventDate.IsZero())

	err = json.Unmarshal([]byte(`{"eventDate":"invalid"}`), &emptyWrapper)
	assert.Error(t, err)
}

func TestDate_Text(t *testing.T) {
	date := NewDate(2026, time.November, 13)

	text, err := date.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "2026-11-13", string(text))

	var parsed Date
	err = parsed.UnmarshalText(text)
	require.NoError(t, err)
	assert.Equal(t, date.String(), parsed.String())

	var zero Date
	zeroText, err := zero.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "0001-01-01", string(zeroText))

	err = parsed.UnmarshalText([]byte(""))
	require.NoError(t, err)
	assert.True(t, parsed.IsZero())
}

func TestDate_IsZero(t *testing.T) {
	var zero Date
	assert.True(t, zero.IsZero())

	nonZero := NewDate(2026, time.November, 13)
	assert.False(t, nonZero.IsZero())
}

func TestDate_Msgp(t *testing.T) {
	date := NewDate(2026, time.November, 13)

	// Test MarshalMsg / UnmarshalMsg
	bytesOut, err := date.MarshalMsg(nil)
	require.NoError(t, err)

	var readDate Date
	rem, err := readDate.UnmarshalMsg(bytesOut)
	require.NoError(t, err)
	assert.Empty(t, rem)
	assert.Equal(t, date.String(), readDate.String())

	// Test EncodeMsg / DecodeMsg
	var buf bytes.Buffer
	writer := msgp.NewWriter(&buf)
	err = date.EncodeMsg(writer)
	require.NoError(t, err)
	err = writer.Flush()
	require.NoError(t, err)

	var decodedDate Date
	reader := msgp.NewReader(&buf)
	err = decodedDate.DecodeMsg(reader)
	require.NoError(t, err)
	assert.Equal(t, date.String(), decodedDate.String())

	assert.Positive(t, date.Msgsize())
}

func TestDateFromTime(t *testing.T) {
	tm := time.Date(2026, time.November, 13, 15, 30, 45, 123456, time.FixedZone("EST", -5*3600))
	date := DateFromTime(tm)

	assert.Equal(t, 2026, date.Year())
	assert.Equal(t, time.November, date.Month())
	assert.Equal(t, 13, date.Day())
	assert.Equal(t, 0, date.Hour())
	assert.Equal(t, 0, date.Minute())
	assert.Equal(t, 0, date.Second())

	zeroDate := DateFromTime(time.Time{})
	assert.True(t, zeroDate.IsZero())
}
