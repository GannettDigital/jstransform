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

var dateTimeDurationRegex = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)W)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?)?$`)

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
// Valid examples: "P1Y2M3D", "P1W", "P1W1D", "P1D", "PT1H", "PT15M", "PT30S", "PT0S", "P1Y2M3DT4H5M6S", "-P1D".
func ParseDuration(s string) (Duration, error) {
	if s == "" {
		return Duration{}, nil
	}

	original := s
	var negate bool
	if strings.HasPrefix(s, "-") {
		negate = true
		s = s[1:]
	}

	matches := dateTimeDurationRegex.FindStringSubmatch(s)
	if len(matches) != 8 {
		return Duration{}, fmt.Errorf("invalid ISO 8601 duration %q", original)
	}

	yearsStr, monthsStr, weeksStr, daysStr := matches[1], matches[2], matches[3], matches[4]
	hoursStr, minutesStr, secondsStr := matches[5], matches[6], matches[7]

	hasDate := yearsStr != "" || monthsStr != "" || weeksStr != "" || daysStr != ""
	hasTime := hoursStr != "" || minutesStr != "" || secondsStr != ""

	if strings.Contains(s, "T") && !hasTime {
		return Duration{}, fmt.Errorf("invalid ISO 8601 duration %q: missing time elements after T", original)
	}
	if !hasDate && !hasTime {
		return Duration{}, fmt.Errorf("invalid ISO 8601 duration %q: empty duration", original)
	}

	var total time.Duration

	if yearsStr != "" {
		v, err := strconv.ParseInt(yearsStr, 10, 64)
		if err != nil {
			return Duration{}, fmt.Errorf("invalid duration years %q: %w", original, err)
		}
		total += time.Duration(v) * 365 * 24 * time.Hour
	}
	if monthsStr != "" {
		v, err := strconv.ParseInt(monthsStr, 10, 64)
		if err != nil {
			return Duration{}, fmt.Errorf("invalid duration months %q: %w", original, err)
		}
		total += time.Duration(v) * 30 * 24 * time.Hour
	}
	if weeksStr != "" {
		v, err := strconv.ParseInt(weeksStr, 10, 64)
		if err != nil {
			return Duration{}, fmt.Errorf("invalid duration weeks %q: %w", original, err)
		}
		total += time.Duration(v) * 7 * 24 * time.Hour
	}
	if daysStr != "" {
		v, err := strconv.ParseInt(daysStr, 10, 64)
		if err != nil {
			return Duration{}, fmt.Errorf("invalid duration days %q: %w", original, err)
		}
		total += time.Duration(v) * 24 * time.Hour
	}
	if hoursStr != "" {
		v, err := strconv.ParseInt(hoursStr, 10, 64)
		if err != nil {
			return Duration{}, fmt.Errorf("invalid duration hours %q: %w", original, err)
		}
		total += time.Duration(v) * time.Hour
	}
	if minutesStr != "" {
		v, err := strconv.ParseInt(minutesStr, 10, 64)
		if err != nil {
			return Duration{}, fmt.Errorf("invalid duration minutes %q: %w", original, err)
		}
		total += time.Duration(v) * time.Minute
	}
	if secondsStr != "" {
		sec, err := strconv.ParseFloat(secondsStr, 64)
		if err != nil {
			return Duration{}, fmt.Errorf("invalid duration seconds %q: %w", original, err)
		}
		total += time.Duration(sec * float64(time.Second))
	}

	if negate {
		total = -total
	}

	return Duration{
		Duration: total,
		raw:      original,
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
	return formatDuration(d.Duration)
}

// MarshalJSON marshals the duration as an ISO 8601 JSON string.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON unmarshals an ISO 8601 JSON string into a Duration.
func (d *Duration) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = Duration{}
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	var err error
	*d, err = ParseDuration(s)
	return err
}

// MarshalText implements encoding.TextMarshaler.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	*d, err = ParseDuration(string(text))
	return err
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
	var negative string
	if d < 0 {
		negative = "-"
		d = -d
	}
	days := d / (24 * time.Hour)
	d %= 24 * time.Hour
	hours := d / time.Hour
	d %= time.Hour
	minutes := d / time.Minute
	d %= time.Minute
	seconds := float64(d) / float64(time.Second)

	var datePart string
	if days > 0 {
		datePart = strconv.FormatInt(int64(days), 10) + "D"
	}

	var timePart string
	if hours > 0 {
		timePart += strconv.FormatInt(int64(hours), 10) + "H"
	}
	if minutes > 0 {
		timePart += strconv.FormatInt(int64(minutes), 10) + "M"
	}
	if seconds > 0 || (hours == 0 && minutes == 0 && days == 0) {
		timePart += strconv.FormatFloat(seconds, 'f', -1, 64) + "S"
	}
	if timePart != "" {
		timePart = "T" + timePart
	}

	return negative + "P" + datePart + timePart
}
