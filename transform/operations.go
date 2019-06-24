package transform

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
)

var durationRe = regexp.MustCompile(`^([\d]*?):?([\d]*):([\d]*)$`)

// duration is a transformOperation which changes from a string duration like "MM:SS" to a number of seconds as
// an integer.
type duration struct {
	re *regexp.Regexp
}

func (c *duration) init(args map[string]string) error {
	c.re = durationRe
	return nil
}

func (c *duration) transform(raw interface{}) (interface{}, error) {
	if array, ok := raw.([]interface{}); ok && len(array) == 1 {
		raw = array[0]
	}

	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("duration only supports strings")
	}

	matches := c.re.FindStringSubmatch(in)
	if matches == nil || len(matches) != 4 {
		return nil, errors.New("duration transform input did not match 'MM:SS' or 'HH:MM:SS'")
	}
	hours, err := strconv.Atoi(matches[1])
	if err != nil && matches[1] != "" {
		return nil, err
	}
	minutes, err := strconv.Atoi(matches[2])
	if err != nil && matches[2] != "" {
		return nil, err
	}
	seconds, err := strconv.Atoi(matches[3])
	if err != nil && matches[3] != "" {
		return nil, err
	}

	minutes += 60 * hours
	seconds += 60 * minutes

	return seconds, nil
}

// changeCase is a transformOperation which changes the case of strings.
type changeCase struct {
	Args map[string]string
}

func (c *changeCase) init(args map[string]string) error {
	if err := requiredArgs([]string{"to"}, args); err != nil {
		return err
	}
	value := args["to"]
	if value != "lower" && value != "upper" {
		return errors.New("the argument 'to' is required and must be either 'lower' or 'upper'")
	}

	c.Args = args
	return nil
}

func (c *changeCase) transform(raw interface{}) (interface{}, error) {
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("changeCase only supports strings")
	}

	switch c.Args["to"] {
	case "lower":
		return strings.ToLower(in), nil
	case "upper":
		return strings.ToUpper(in), nil
	}
	return nil, errors.New("unknown error in changeCase")
}

// inverse is a transformOperation which flips the value of a boolean.
type inverse struct {
	args map[string]string
}

func (i *inverse) init(args map[string]string) error {
	return nil
}

func (i *inverse) transform(raw interface{}) (interface{}, error) {
	in, ok := raw.(bool)
	if !ok {
		return nil, errors.New("inverse only supports booleans")
	}
	return !in, nil
}

// max is a transformOperation which retrieves a field from the maximum item in an array.
// The maxiumum item is determined by comparing values in a defined number field on the array items.
type max struct {
	Args map[string]string
}

func (m *max) init(args map[string]string) error {
	if err := requiredArgs([]string{"by", "return"}, args); err != nil {
		return err
	}
	m.Args = args
	return nil
}

func (m *max) transform(in interface{}) (interface{}, error) {
	inArray, ok := in.([]interface{})
	if !ok {
		return nil, errors.New("input must be an array")
	}
	byArg := strings.Replace(m.Args["by"], "@", "$", 1)
	returnArg := strings.Replace(m.Args["return"], "@", "$", 1)

	var largest float64
	var largestIndex int
	for i, item := range inArray {
		byRaw, err := jsonpath.Get(byArg, item)
		if err != nil {
			return nil, fmt.Errorf("failed extracting 'by' field: %v", err)
		}
		by, ok := byRaw.(float64)
		if !ok {
			byInt, ok := byRaw.(int)
			if !ok {
				return nil, errors.New("by field is not a number")
			}
			by = float64(byInt)
		}
		if by > largest {
			largest = by
			largestIndex = i
		}
	}

	rawReturn, err := jsonpath.Get(returnArg, inArray[largestIndex])
	if err != nil {
		return nil, fmt.Errorf("failed extracting 'return' field: %v", err)
	}

	return rawReturn, nil
}

// replace is a transformOperation which performs a regex based find/replace on a string value.
type replace struct {
	Args  map[string]string
	regex *regexp.Regexp
}

func (r *replace) init(args map[string]string) error {
	if err := requiredArgs([]string{"regex", "new"}, args); err != nil {
		return err
	}
	re, err := regexp.Compile(args["regex"])
	if err != nil {
		return fmt.Errorf("failed to parse regex %q: %v", args["regex"], err)
	}

	r.regex = re
	r.Args = args
	return nil
}

func (r *replace) transform(raw interface{}) (interface{}, error) {
	if r.regex == nil {
		return nil, errors.New("init was not run")
	}
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("replace only supports strings")
	}

	return r.regex.ReplaceAllString(in, r.Args["new"]), nil
}

// split is a transformOperation which splits a string based on a given split string.
type split struct {
	Args map[string]string
}

func (s *split) init(args map[string]string) error {
	if err := requiredArgs([]string{"on"}, args); err != nil {
		return err
	}

	s.Args = args
	return nil
}

func (s *split) transform(raw interface{}) (interface{}, error) {
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("split only supports strings")
	}

	splits := strings.Split(in, s.Args["on"])

	// Return []interface{} to avoid messing up type casts later in the process
	interfaceSplits := []interface{}{}
	for _, s := range splits {
		interfaceSplits = append(interfaceSplits, s)
	}
	return interfaceSplits, nil
}

// timeParse is a transformOperation which formats a date string into the layout
type timeParse struct {
	Args map[string]string
}

func (t *timeParse) init(args map[string]string) error {
	if err := requiredArgs([]string{"format", "layout"}, args); err != nil {
		return err
	}

	t.Args = args
	return nil
}

func (t *timeParse) transform(raw interface{}) (interface{}, error) {
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("timeParse only supports strings")
	}
	parsedTime, err := time.Parse(t.Args["format"], in)
	if err != nil {
		return nil, fmt.Errorf("time could not be parsed using supplied format")
	}
	return parsedTime.Format(t.Args["layout"]), nil
}

// stringToInteger is a transformOperation which takes a string and converts it into an Int
type stringToInteger struct {
}

func (s *stringToInteger) init(args map[string]string) error {
	return nil
}

func (s *stringToInteger) transform(raw interface{}) (interface{}, error) {
	str, ok := raw.(string)
	if !ok {
		return nil, errors.New("stringToInteger only supports strings")
	}

	result, err := strconv.Atoi(str)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error converting to an Integer from a string: %s", err))
	} 

	return result, nil
}

// requiredArgs checks the given args map to make sure it contains the required args and only the required args.
func requiredArgs(required []string, args map[string]string) error {
	if len(args) != len(required) {
		return fmt.Errorf("expected args %v but got %d args", required, len(args))
	}
	for _, arg := range required {
		if _, ok := args[arg]; !ok {
			return fmt.Errorf("argument %q is required", arg)
		}
	}
	return nil
}
