package generate

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/GannettDigital/msgp/msgp"
)

var (
	weekDurationRegex     = regexp.MustCompile(`^P(\d+)W$`)
	dateTimeDurationRegex = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?$`)
)

// Duration represents an ISO 8601 duration according to RFC 3339 Appendix A.
type Duration struct {
	time.Duration

	raw string
}

// DurationFromTime creates a Duration from a standard time.Duration.
func DurationFromTime(d time.Duration) Duration {
	return Duration{Duration: d}
}

// ParseDuration parses an ISO 8601 duration string per RFC 3339 Appendix A.
// Valid examples: "P1Y2M3D", "P1W", "P1D", "PT1H", "PT15M", "PT30S", "PT0S", "P1Y2M3DT4H5M6S".
func ParseDuration(s string) (Duration, error) {
	if s == "" {
		return Duration{}, nil
	}

	if matches := weekDurationRegex.FindStringSubmatch(s); len(matches) == 2 {
		weeks, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return Duration{}, fmt.Errorf("invalid duration weeks %q: %w", s, err)
		}
		return Duration{
			Duration: time.Duration(weeks) * 7 * 24 * time.Hour,
			raw:      s,
		}, nil
	}

	matches := dateTimeDurationRegex.FindStringSubmatch(s)
	if len(matches) != 7 {
		return Duration{}, fmt.Errorf("invalid ISO 8601 duration %q", s)
	}

	yearsStr, monthsStr, daysStr := matches[1], matches[2], matches[3]
	hoursStr, minutesStr, secondsStr := matches[4], matches[5], matches[6]

	hasDate := yearsStr != "" || monthsStr != "" || daysStr != ""
	hasTime := hoursStr != "" || minutesStr != "" || secondsStr != ""

	if strings.Contains(s, "T") && !hasTime {
		return Duration{}, fmt.Errorf("invalid ISO 8601 duration %q: missing time elements after T", s)
	}
	if !hasDate && !hasTime {
		return Duration{}, fmt.Errorf("invalid ISO 8601 duration %q: empty duration", s)
	}

	var total time.Duration

	if yearsStr != "" {
		v, _ := strconv.ParseInt(yearsStr, 10, 64)
		total += time.Duration(v) * 365 * 24 * time.Hour
	}
	if monthsStr != "" {
		v, _ := strconv.ParseInt(monthsStr, 10, 64)
		total += time.Duration(v) * 30 * 24 * time.Hour
	}
	if daysStr != "" {
		v, _ := strconv.ParseInt(daysStr, 10, 64)
		total += time.Duration(v) * 24 * time.Hour
	}
	if hoursStr != "" {
		v, _ := strconv.ParseInt(hoursStr, 10, 64)
		total += time.Duration(v) * time.Hour
	}
	if minutesStr != "" {
		v, _ := strconv.ParseInt(minutesStr, 10, 64)
		total += time.Duration(v) * time.Minute
	}
	if secondsStr != "" {
		v, _ := strconv.ParseInt(secondsStr, 10, 64)
		total += time.Duration(v) * time.Second
	}

	return Duration{
		Duration: total,
		raw:      s,
	}, nil
}

// TimeDuration returns the underlying standard time.Duration.
func (d Duration) TimeDuration() time.Duration {
	return d.Duration
}

// IsZero returns true if the duration is unset or zero value.
func (d Duration) IsZero() bool {
	return d.Duration == 0 && d.raw == ""
}

// String returns the duration formatted as an ISO 8601 duration string.
func (d Duration) String() string {
	if d.raw != "" {
		return d.raw
	}
	if d.Duration == 0 {
		return ""
	}
	return formatDuration(d.Duration)
}

// MarshalJSON marshals the duration as an ISO 8601 JSON string, or null if zero.
func (d Duration) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.String())
}

// UnmarshalJSON unmarshals an ISO 8601 JSON string into a Duration.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil || *s == "" {
		*d = Duration{}
		return nil
	}
	parsed, err := ParseDuration(*s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (d Duration) MarshalText() ([]byte, error) {
	if d.IsZero() {
		return nil, nil
	}
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Duration) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" {
		*d = Duration{}
		return nil
	}
	parsed, err := ParseDuration(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// MarshalMsg implements msgp.Marshaler.
func (d Duration) MarshalMsg(b []byte) ([]byte, error) {
	return msgp.AppendString(b, d.String()), nil
}

// UnmarshalMsg implements msgp.Unmarshaler.
func (d *Duration) UnmarshalMsg(bts []byte) ([]byte, error) {
	s, bts, err := msgp.ReadStringBytes(bts)
	if err != nil {
		return bts, err
	}
	parsed, err := ParseDuration(s)
	if err != nil {
		return bts, err
	}
	*d = parsed
	return bts, nil
}

// EncodeMsg implements msgp.Encodable.
func (d Duration) EncodeMsg(en *msgp.Writer) error {
	return en.WriteString(d.String())
}

// DecodeMsg implements msgp.Decodable.
func (d *Duration) DecodeMsg(dc *msgp.Reader) error {
	s, err := dc.ReadString()
	if err != nil {
		return err
	}
	parsed, err := ParseDuration(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// Msgsize implements msgp.Sizer.
func (d Duration) Msgsize() int {
	return msgp.StringPrefixSize + len(d.String())
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "PT0S"
	}
	var negative bool
	if d < 0 {
		negative = true
		d = -d
	}
	days := d / (24 * time.Hour)
	d %= 24 * time.Hour
	hours := d / time.Hour
	d %= time.Hour
	minutes := d / time.Minute
	d %= time.Minute
	seconds := d / time.Second

	var b strings.Builder
	if negative {
		b.WriteString("-")
	}
	b.WriteString("P")
	if days > 0 {
		b.WriteString(strconv.FormatInt(int64(days), 10))
		b.WriteString("D")
	}
	if hours > 0 || minutes > 0 || seconds > 0 || days == 0 {
		b.WriteString("T")
		if hours > 0 {
			b.WriteString(strconv.FormatInt(int64(hours), 10))
			b.WriteString("H")
		}
		if minutes > 0 {
			b.WriteString(strconv.FormatInt(int64(minutes), 10))
			b.WriteString("M")
		}
		if seconds > 0 || (hours == 0 && minutes == 0 && days == 0) {
			b.WriteString(strconv.FormatInt(int64(seconds), 10))
			b.WriteString("S")
		}
	}
	return b.String()
}
