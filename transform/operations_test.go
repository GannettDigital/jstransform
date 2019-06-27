package transform

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type testOp struct {
	args   map[string]string
	called interface{}
	fail   bool
}

func (op *testOp) init(args map[string]string) error {
	fail, err := strconv.ParseBool(args["fail"])
	if err != nil {
		return err
	}
	op.fail = fail
	op.args = args
	return nil
}

func (op *testOp) transform(in interface{}) (interface{}, error) {
	if op.fail {
		return nil, errors.New("fail")
	}
	op.called = in
	return op.args["out"], nil
}

type opTests struct {
	description string
	args        map[string]string
	in          interface{}
	want        interface{}
	wantErr     bool
	wantInitErr bool
}

// A common test runner for all the operations tests
func runOpTests(t *testing.T, opType func() transformOperation, tests []opTests) {

	for _, test := range tests {
		runOpTest(t, opType, test)
	}
}

func runOpTest(t *testing.T, opType func() transformOperation, test opTests) {
	op, err := runOpTestInit(opType, test)
	if test.wantInitErr && err != nil {
		return
	} else if err != nil {
		t.Error(err)
	}

	_, err = runOpTestTransform(op, test)
	if test.wantErr && err != nil {
		return
	} else if err != nil {
		t.Error(err)
	}
}

func runOpTestInit(opType func() transformOperation, test opTests) (transformOperation, error) {
	op := opType()
	err := op.init(test.args)
	switch {
	case test.wantInitErr && err != nil:
		return op, err
	case test.wantInitErr && err == nil:
		return op, fmt.Errorf("Test %q - got init error nil, want error", test.description)
	case !test.wantErr && err != nil:
		return op, fmt.Errorf("Test %q - got init error, want nil: %v", test.description, err)
	}
	return op, nil
}

func runOpTestTransform(op transformOperation, test opTests) (interface{}, error) {
	got, err := op.transform(test.in)

	switch {
	case test.wantErr && err != nil:
		return got, err
	case test.wantErr && err == nil:
		return got, fmt.Errorf("Test %q - got nil, want error", test.description)
	case !test.wantErr && err != nil:
		return got, fmt.Errorf("Test %q - got error, want nil: %v", test.description, err)
	case !reflect.DeepEqual(got, test.want):
		return got, fmt.Errorf("Test %q - got\n%s\nwant\n%s", test.description, got, test.want)
	}

	return got, nil
}

func TestDuration(t *testing.T) {
	tests := []opTests{
		{
			description: "working MM:SS",
			in:          "01:13",
			want:        73,
		},
		{
			description: "working HH:MM:SS",
			in:          "01:01:13",
			want:        3673,
		},
		{
			description: "zero minutes",
			in:          "01:00:13",
			want:        3613,
		},
		{
			description: "zero seconds",
			in:          "01:00",
			want:        60,
		},
		{
			description: "seconds only",
			in:          ":13",
			want:        13,
		},
		{
			description: "not a string",
			in:          13,
			wantErr:     true,
		},
		{
			description: "Invalid string",
			in:          "2hours",
			wantErr:     true,
		},
		{
			description: "Invalid hour",
			in:          "blah:01:13",
			wantErr:     true,
		},
		{
			description: "Invalid minute",
			in:          "1:!:13",
			wantErr:     true,
		},
		{
			description: "Invalid second",
			in:          "1:00:ab",
			wantErr:     true,
		},
	}
	runOpTests(t, func() transformOperation { return &duration{} }, tests)
}

func TestChangeCase(t *testing.T) {
	tests := []opTests{
		{
			description: "Simple working case, lower",
			args:        map[string]string{"to": "lower"},
			in:          "MixedCase",
			want:        "mixedcase",
		},
		{
			description: "Simple working case, upper",
			args:        map[string]string{"to": "upper"},
			in:          "MixedCase",
			want:        "MIXEDCASE",
		},
		{
			description: "Too many args",
			args:        map[string]string{"to": "lower", "from": "?"},
			in:          "MixedCase",
			wantInitErr: true,
		},
		{
			description: "Missing to arg",
			args:        map[string]string{"from": "?"},
			in:          "MixedCase",
			wantInitErr: true,
		},
		{
			description: "Missing all arg",
			args:        map[string]string{},
			in:          "MixedCase",
			wantInitErr: true,
		},
		{
			description: "Invalid to",
			args:        map[string]string{"to": "?"},
			in:          "MixedCase",
			wantErr:     true,
		},
		{
			description: "Non-string input",
			args:        map[string]string{"to": "lower"},
			in:          5,
			wantErr:     true,
		},
	}

	runOpTests(t, func() transformOperation { return &changeCase{} }, tests)
}

func TestInverse(t *testing.T) {
	tests := []opTests{
		{
			description: "Simple working case, true",
			in:          true,
			want:        false,
		},
		{
			description: "Simple working case, false",
			in:          false,
			want:        true,
		},
		{
			description: "string input",
			in:          "a",
			wantErr:     true,
		},
		{
			description: "number input",
			in:          1,
			wantErr:     true,
		},
	}

	runOpTests(t, func() transformOperation { return &inverse{} }, tests)
}

func TestMax(t *testing.T) {
	tests := []opTests{
		{
			description: "Simple working case",
			args:        map[string]string{"by": "@.encodingRate", "return": "@.url"},
			in: []interface{}{
				map[string]interface{}{"url": "max", "encodingRate": 10},
				map[string]interface{}{"url": "min", "encodingRate": 2},
			},
			want: "max",
		},
		{
			description: "Extra args",
			args:        map[string]string{"by": "@.encodingRate", "return": "@.url", "bye": "bye"},
			in: []interface{}{
				map[string]interface{}{"url": "max", "encodingRate": 10},
				map[string]interface{}{"url": "min", "encodingRate": 2},
			},
			wantInitErr: true,
		},
		{
			description: "Missing by arg",
			args:        map[string]string{"return": "@.url"},
			in: []interface{}{
				map[string]interface{}{"url": "max", "encodingRate": 10},
				map[string]interface{}{"url": "min", "encodingRate": 2},
			},
			wantInitErr: true,
		},
		{
			description: "Missing return arg",
			args:        map[string]string{"by": "@.encodingRate"},
			in: []interface{}{
				map[string]interface{}{"url": "max", "encodingRate": 10},
				map[string]interface{}{"url": "min", "encodingRate": 2},
			},
			wantInitErr: true,
		},
		{
			description: "by field is not a number",
			args:        map[string]string{"by": "@.encodingRate", "return": "@.url"},
			in: []interface{}{
				map[string]interface{}{"url": "max", "encodingRate": "10"},
				map[string]interface{}{"url": "min", "encodingRate": "2"},
			},
			wantErr: true,
		},
		{
			description: "return field does not exist",
			args:        map[string]string{"by": "@.encodingRate", "return": "@.url"},
			in: []interface{}{
				map[string]interface{}{"encodingRate": 10},
				map[string]interface{}{"encodingRate": 2},
			},
			wantErr: true,
		},
		{
			description: "in is not an array",
			args:        map[string]string{"by": "@.encodingRate", "return": "@.url"},
			in:          map[string]interface{}{"url": "max", "encodingRate": 10},
			wantErr:     true,
		},
	}

	runOpTests(t, func() transformOperation { return &max{} }, tests)
}

func TestReplace(t *testing.T) {
	tests := []opTests{
		{
			description: "Simple working case",
			args:        map[string]string{"regex": `foo\.com`, "new": `media.gannett-cdn.com`},
			in:          "http://foo.com",
			want:        "http://media.gannett-cdn.com",
		},
		{
			description: "Missing regex arg",
			args:        map[string]string{"new": `media.gannett-cdn.com`},
			in:          "http://foo.com",
			wantInitErr: true,
		},
		{
			description: "Missing new arg",
			args:        map[string]string{"regex": `foo\.com`},
			in:          "http://foo.com",
			wantInitErr: true,
		},
		{
			description: "Extra args",
			args:        map[string]string{"regex": `foo\.com`, "new": `media.gannett-cdn.com`, "alt": "a"},
			in:          "http://foo.com",
			wantInitErr: true,
		},
		{
			description: "Non string input",
			args:        map[string]string{"regex": `foo\.com`, "new": `media.gannett-cdn.com`},
			in:          5,
			wantErr:     true,
		},
	}

	runOpTests(t, func() transformOperation { return &replace{} }, tests)
}

func TestSplit(t *testing.T) {
	tests := []opTests{
		{
			description: "Simple working case",
			args:        map[string]string{"on": "|"},
			in:          "a|b|c",
			want:        []interface{}{"a", "b", "c"},
		},
		{
			description: "Missing on arg",
			args:        map[string]string{},
			in:          "a|b|c",
			wantInitErr: true,
		},
		{
			description: "Too many args",
			args:        map[string]string{"on": "|", "off": "|"},
			in:          "a|b|c",
			wantInitErr: true,
		},
		{
			description: "Non string input",
			args:        map[string]string{"on": "|"},
			in:          8,
			wantErr:     true,
		},
	}

	runOpTests(t, func() transformOperation { return &split{} }, tests)
}

func TestTimeParse(t *testing.T) {
	tests := []opTests{
		{
			description: "Simple working case",
			args:        map[string]string{"format": time.RFC3339, "layout": "2006-01-02"},
			in:          "2019-05-16T21:00:00-04:00",
			want:        "2019-05-16",
		},
		{
			description: "Missing arg",
			args:        map[string]string{"format": time.RFC3339},
			in:          "2019-05-16T21:00:00-04:00",
			wantInitErr: true,
		},
		{
			description: "Too many args",
			args:        map[string]string{"format": time.RFC3339, "layout": "2006-01-02", "cookies": "failure"},
			in:          "2019-05-16T21:00:00-04:00",
			wantInitErr: true,
		},
		{
			description: "Non string input",
			args:        map[string]string{"format": time.RFC3339, "layout": "2006-01-02"},
			in:          -1,
			wantErr:     true,
		},
		{
			description: "Invalid format",
			args:        map[string]string{"format": "foobar", "layout": "2006-01-02"},
			in:          "2019-05-16T21:00:00-04:00",
			wantErr:     true,
		},
	}
	runOpTests(t, func() transformOperation { return &timeParse{} }, tests)
}

func TestCurrentTime(t *testing.T) {
	tests := []opTests{
		{
			description: "Simple working case",
			args:        map[string]string{"format": "RFC3339"},
			want:        time.Now().Format(time.RFC3339),
		},
		{
			description: "Missing arg",
			args:        map[string]string{},
			wantInitErr: true,
		},
		{
			description: "Too many args",
			args:        map[string]string{"format": "RFC3339", "cookies": "failure"},
			wantInitErr: true,
		},
		{
			description: "Format not predefined",
			args:		 map[string]string{"format": "Mon Jan 2 15:04:05 MST 2006"},
			want:		time.Now().Format("Mon Jan 2 15:04:05 MST 2006")	,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			op, err := runOpTestInit(func() transformOperation { return &currentTime{} }, test)
			if err != nil && test.wantInitErr {
				return
			}
			if err != nil && !test.wantInitErr {
				t.Fatal(err)
			}
			got, err := runOpTestTransform(op, test)

			if err != nil {
				t.Fatal(err)
				return
			}
			now := time.Now()
			result, ok := got.(string)
			if !ok {
				t.Fatal("function must return string")
			}

			var gotParse time.Time

			if test.args["format"] == "RFC3339" {
				gotParse, err = time.Parse(time.RFC3339, result)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				gotParse, err = time.Parse(test.args["format"], result)
				if err != nil {
					t.Fatal(err)
				}
			}
			if result := compareTimeStamps(now, gotParse); !result {
				t.Fatal("time returned not close enough to current time")
			}
		})

	}
}

func TestStringToInteger(t *testing.T) {
	tests := []opTests{

		{
			description: "Simple working case",
			in:          "237754",
			want:        237754,
		},

		{
			description: "Boolean input type",
			in:          true,
			wantErr:     true,
		},
	}
	runOpTests(t, func() transformOperation { return &stringToInteger{} }, tests)
}

func compareTimeStamps(time1 time.Time, time2 time.Time) bool {
	maxTimeDifference := time.Duration(300) * time.Second
	actualDiff := absValue(time1.Sub(time2))
	return actualDiff < maxTimeDifference
}

func absValue (x time.Duration) time.Duration {
	if x < 0 {
		return -x
	}
	return x
}

//left this here so that we can optimize the test strategy to include comparing the error to the want error
//func compareErrs (gotErr error, wantErr error) error {
//	if !reflect.DeepEqual(gotErr, wantErr) {
//		return errors.New(fmt.Sprintf("Errors did not match, got %s, want %s", gotErr, wantErr))
//	}
//	return nil
//}

