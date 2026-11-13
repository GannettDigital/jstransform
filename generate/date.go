package generate

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/GannettDigital/msgp/msgp"
)

// DateFormat is the RFC 3339 Section 5.6 full-date format (YYYY-MM-DD).
const DateFormat = "2006-01-02"

// Date represents an RFC 3339 Section 5.6 full-date (YYYY-MM-DD).
// It wraps time.Time with UTC date normalization and JSON/text/msgp serialization.
type Date struct {
	time.Time
}

// NewDate returns a Date with the given year, month, and day in UTC.
func NewDate(year int, month time.Month, day int) Date {
	return Date{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// DateFromTime converts a time.Time into a Date truncated to midnight UTC.
func DateFromTime(t time.Time) Date {
	if t.IsZero() {
		return Date{}
	}
	utc := t.UTC()
	return NewDate(utc.Year(), utc.Month(), utc.Day())
}

// ParseDate parses an RFC 3339 full-date string (YYYY-MM-DD).
func ParseDate(s string) (Date, error) {
	if s == "" {
		return Date{}, nil
	}
	t, err := time.Parse(DateFormat, s)
	if err != nil {
		return Date{}, fmt.Errorf("invalid date format %q: %w", s, err)
	}
	return Date{Time: t}, nil
}

// String returns the date formatted as YYYY-MM-DD, or an empty string if zero.
func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.Format(DateFormat)
}

// MarshalJSON marshals the date as a JSON string formatted as "YYYY-MM-DD", or null if zero.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Format(DateFormat))
}

// UnmarshalJSON unmarshals a JSON string into a Date.
func (d *Date) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil || *s == "" {
		d.Time = time.Time{}
		return nil
	}
	parsed, err := ParseDate(*s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (d Date) MarshalText() ([]byte, error) {
	if d.IsZero() {
		return nil, nil
	}
	return []byte(d.Format(DateFormat)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Date) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" {
		d.Time = time.Time{}
		return nil
	}
	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// MarshalMsg implements msgp.Marshaler.
func (d Date) MarshalMsg(b []byte) ([]byte, error) {
	return msgp.AppendString(b, d.String()), nil
}

// UnmarshalMsg implements msgp.Unmarshaler.
func (d *Date) UnmarshalMsg(bts []byte) ([]byte, error) {
	s, bts, err := msgp.ReadStringBytes(bts)
	if err != nil {
		return bts, err
	}
	parsed, err := ParseDate(s)
	if err != nil {
		return bts, err
	}
	*d = parsed
	return bts, nil
}

// EncodeMsg implements msgp.Encodable.
func (d Date) EncodeMsg(en *msgp.Writer) error {
	return en.WriteString(d.String())
}

// DecodeMsg implements msgp.Decodable.
func (d *Date) DecodeMsg(dc *msgp.Reader) error {
	s, err := dc.ReadString()
	if err != nil {
		return err
	}
	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// Msgsize implements msgp.Sizer.
func (d Date) Msgsize() int {
	return msgp.StringPrefixSize + len(DateFormat)
}
