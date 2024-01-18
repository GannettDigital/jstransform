package transform

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/antchfx/xmlquery"
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

// A common test runner for all the operations tests.
func runOpTests(t *testing.T, opType func() transformOperation, tests []opTests) {
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			runOpTest(t, opType, test)
		})
	}
}

func runOpTest(t *testing.T, opType func() transformOperation, test opTests) {
	op := opType()
	err := op.init(test.args)

	if err := compareWantErrs(err, test.wantInitErr); err != nil {
		t.Fatal(err)
	}
	if test.wantInitErr {
		return
	}
	got, err := op.transform(test.in)

	if err := compareWantErrs(err, test.wantErr); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, test.want) {
		t.Fatalf("got: %s, want: %s", got, test.want)
	}
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
		{
			description: "Null value in json",
			want:        0,
			wantErr:     false,
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
			wantInitErr: true,
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

func TestValueExists(t *testing.T) {
	tests := []opTests{
		{
			description: "Simple working string case, true",
			in:          "someValue",
			want:        true,
		},
		{
			description: "Simple working string case, false",
			in:          "",
			want:        false,
		},
		{
			description: "Simple working XML node array case, true",
			in:          []*xmlquery.Node{{}},
			want:        true,
		},
		{
			description: "Simple working XML node array case, false",
			in:          []*xmlquery.Node{},
			want:        false,
		},
	}

	runOpTests(t, func() transformOperation { return &valueExists{} }, tests)
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
		{
			description: "Empty input",
			args:        map[string]string{"on": "|"},
			in:          "",
			want:        []interface{}{},
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

func TestToCamelCase(t *testing.T) {
	tests := []opTests{
		{
			description: "Simple working case",
			args:        map[string]string{"delimiter": "-"},
			in:          "extra-base-hit",
			want:        "extraBaseHit",
		},
		{
			description: "Missing an arg",
			args:        map[string]string{},
			in:          "extra-base-hit",
			wantInitErr: true,
		},
		{
			description: "Non-string input",
			args:        map[string]string{"delimiter": "-"},
			in:          1234,
			wantErr:     true,
		},
		{
			description: "Too many args",
			args:        map[string]string{"delimiter": "-", "otherDelimiter": ","},
			in:          "extra-base-hit",
			wantInitErr: true,
		},
	}

	runOpTests(t, func() transformOperation { return &toCamelCase{} }, tests)
}

func TestCurrentTime(t *testing.T) {
	tests := []opTests{
		{
			description: "Simple working case",
			args:        map[string]string{"format": time.RFC3339},
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
			args:        map[string]string{"format": "Mon Jan 2 15:04:05 MST 2006"},
			want:        time.Now().Format("Mon Jan 2 15:04:05 MST 2006"),
		},
		{
			description: "Passing only time constant as a string",
			args:        map[string]string{"format": "RFC3339"},
			want:        time.Now().Format(time.RFC3339),
		},
	}

	runOpTests(t, func() transformOperation { return &currentTime{} }, tests)
}

func TestRemoveHTML(t *testing.T) {
	tests := []opTests{
		{
			description: "Working case one",
			in:          "Japanese man calls for a trial by combat when talking to his wife's <br>mother in law<br> &",
			want:        "Japanese man calls for a trial by combat when talking to his wife's mother in law &",
		},
		{
			description: "Working case two",
			in:          "<div><br>++Japanese <i>man</i> Collin D'Oliver Murasame <strong>calls for a trial by combat when talking</strong> to his wife<custom>",
			want:        "++Japanese man Collin D'Oliver Murasame calls for a trial by combat when talking to his wife",
		},
		{
			description: "Working case three",
			in:          "<p dir=\"ltr\">STILLWATER — Oklahoma State’s men’s new-look basketball team has versatility and a lot of hype.</p><p dir=\"ltr\">The Cowboys finalized their top-10 recruiting class on Wednesday, signing a talented trio from across the continent to join early signees Cade Cunningham, Montreal Pena and Rondel Walker.&nbsp;</p><p dir=\"ltr\">OSU added Canada’s top player, Matthew-Alexander Moncrieffe along with dynamic scorer Donovan Williams and sharp-shooting graduate transfer Ferron Flavors Jr.</p><p dir=\"ltr\">The class is ranked No. 4 by Rivals.com and No. 9 by 247Sports.com.</p><p dir=\"ltr\">It’s also the biggest step in coach Mike Boynton’s plan to rebuild the program.",
			want:        "STILLWATER — Oklahoma State’s men’s new-look basketball team has versatility and a lot of hype. The Cowboys finalized their top-10 recruiting class on Wednesday, signing a talented trio from across the continent to join early signees Cade Cunningham, Montreal Pena and Rondel Walker.  OSU added Canada’s top player, Matthew-Alexander Moncrieffe along with dynamic scorer Donovan Williams and sharp-shooting graduate transfer Ferron Flavors Jr. The class is ranked No. 4 by Rivals.com and No. 9 by 247Sports.com. It’s also the biggest step in coach Mike Boynton’s plan to rebuild the program.",
		},
	}

	runOpTests(t, func() transformOperation { return &removeHTML{} }, tests)
}

func TestConvertToFloat64(t *testing.T) {
	tests := []opTests{
		{
			description: "Working case one",
			in:          "0.96102143",
			want:        float64(0.96102143),
		},
		{
			description: "Working case two",
			in:          "0",
			want:        float64(0),
		},
		{
			description: "Working case three",
			in:          "124",
			want:        float64(124),
		},
		{
			description: "Handle an integer",
			in:          0,
			want:        float64(0),
		},
		{
			description: "Handle a float64",
			in:          float64(123),
			want:        float64(123),
			wantErr:     false,
		},
		{
			description: "Can't convert from string to float64",
			in:          "Hello",
			want:        float64(0),
			wantErr:     true,
		},
	}

	runOpTests(t, func() transformOperation { return &convertToFloat64{} }, tests)
}

func TestConvertToInt64(t *testing.T) {
	tests := []opTests{
		{
			description: "Working case one",
			in:          "1",
			want:        int64(1),
		},
		{
			description: "Working case two",
			in:          "0",
			want:        int64(0),
		},
		{
			description: "Handle an int",
			in:          int(123),
			want:        int64(123),
			wantErr:     false,
		},
		{
			description: "Handle an int64",
			in:          int64(123),
			want:        int64(123),
			wantErr:     false,
		},
		{
			description: "Handle a float64",
			in:          float64(3.14),
			want:        int64(3),
			wantErr:     false,
		},
		{
			description: "Handle a float32",
			in:          float32(3.14),
			want:        int64(3),
			wantErr:     false,
		},
		{
			description: "Can't convert from string to int64",
			in:          "Hello",
			want:        int64(0),
			wantErr:     true,
		},
	}

	runOpTests(t, func() transformOperation { return &convertToInt64{} }, tests)
}

func TestConvertToBool(t *testing.T) {
	tests := []opTests{
		{
			description: "Working string case",
			in:          "1",
			want:        true,
		},
		{
			description: "Failing string case",
			in:          "",
			want:        false,
		},
		{
			description: "Working int case",
			in:          int(123),
			want:        true,
			wantErr:     false,
		},
		{
			description: "Failing int case",
			in:          int(0),
			want:        false,
			wantErr:     false,
		},
		{
			description: "Working float32 case",
			in:          float32(3.14),
			want:        true,
			wantErr:     false,
		},
		{
			description: "Failing float32 case",
			in:          float32(0),
			want:        false,
			wantErr:     false,
		},
		{
			description: "Working float64 case",
			in:          float64(3.14),
			want:        true,
			wantErr:     false,
		},
		{
			description: "Failing float64 case",
			in:          float64(0),
			want:        false,
			wantErr:     false,
		},
		{
			description: "Can't convert from arbitrary object to bool",
			in: struct {
				Field string
			}{
				Field: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}

	runOpTests(t, func() transformOperation { return &convertToBool{} }, tests)
}

func compareWantErrs(gotErr error, wantErr bool) error {
	switch {
	case wantErr && gotErr == nil:
		return errors.New("expected error and didn't get one")
	case wantErr && gotErr != nil:
		return nil
	case !wantErr && gotErr == nil:
		return nil
	case !wantErr && gotErr != nil:
		return fmt.Errorf("got error unexpected error: %v", gotErr)
	}
	return nil
}
