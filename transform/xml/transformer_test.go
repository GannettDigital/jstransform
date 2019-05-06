package xmlTransform

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/GannettDigital/jstransform/jsonschema"
)

func TestNewTransformer(t *testing.T) {
	tests := []struct {
		description    string
		schemaFilePath string
		xmlFilePath    string
		wantFilePath   string
	}{
		{
			description:    "teams NBA",
			schemaFilePath: "./test_data/sports/teams/teams.json",
			xmlFilePath:    "./test_data/sports/teams/teams_NBA.xml",
			wantFilePath:   "./test_data/sports/teams/teamsNBA.out.json",
		},
		{
			description:    "teams MLB",
			schemaFilePath: "./test_data/sports/teams/teams.json",
			xmlFilePath:    "./test_data/sports/teams/teams_MLB.xml",
			wantFilePath:   "./test_data/sports/teams/teamsMLB.out.json",
		},
		{
			description:    "array-transforms",
			schemaFilePath: "./test_data/array-transforms.json",
			xmlFilePath:    "./test_data/array-transforms.xml",
			wantFilePath:   "./test_data/array-transforms.out.json",
		},
		{
			description:    "multiple-array-transforms",
			schemaFilePath: "./test_data/multiple-arrays.json",
			xmlFilePath:    "./test_data/multiple-arrays.xml",
			wantFilePath:   "./test_data/multiple-arrays.out.json",
		},
		{
			description:    "conversion-transforms",
			schemaFilePath: "./test_data/conversion-transforms.json",
			xmlFilePath:    "./test_data/conversion-transforms.xml",
			wantFilePath:   "./test_data/conversion-transforms.out.json",
		},
		{
			description:    "attribute-selection",
			schemaFilePath: "./test_data/attribute-selection.json",
			xmlFilePath:    "./test_data/attribute-selection.xml",
			wantFilePath:   "./test_data/attribute-selection.out.json",
		},
		{
			description:    "operations",
			schemaFilePath: "./test_data/operations.json",
			xmlFilePath:    "./test_data/operations.xml",
			wantFilePath:   "./test_data/operations.out.json",
		},
	}

	for _, test := range tests {
		schema, err := jsonschema.SchemaFromFile(test.schemaFilePath, "")
		if err != nil {
			t.Fatal(err)
		}

		tr, err := NewTransformer(schema, "sport")
		if err != nil {
			t.Fatal(err)
		}

		rawXMLBytes, err := ioutil.ReadFile(test.xmlFilePath)
		if err != nil {
			t.Fatal(err)
		}

		output, err := tr.Transform(rawXMLBytes)
		if err != nil {
			t.Fatal(err)
		}

		want, err := ioutil.ReadFile(test.wantFilePath)

		var (
			outputMap map[string]interface{}
			wantMap   map[string]interface{}
		)

		if err := json.Unmarshal(output, &outputMap); err != nil {
			t.Fatal(err)
		}

		if err := json.Unmarshal(want, &wantMap); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(outputMap, wantMap) {
			t.Fatalf("test %s failed \n want:\n %s \n got:\n %s", test.description, want, output)
		}
	}

}
