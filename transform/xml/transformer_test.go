package xmlTransform

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/GannettDigital/jstransform/jsonschema"
)

func TestNewTransformer(t *testing.T) {
	tests := []struct {
		description    string
		schemaFilePath string
		xmlFilePath    string
	}{
		//{
		//	description:    "teams NBA",
		//	schemaFilePath: "./test_data/sports/teams/teams.json",
		//	xmlFilePath:    "./test_data/sports/teams/teams_NBA.xml",
		//},
		//{
		//	description:    "teams MLB",
		//	schemaFilePath: "./test_data/sports/teams/teams.json",
		//	xmlFilePath:    "./test_data/sports/teams/teams_MLB.xml",
		//},
		{
			description:    "boxscores NBA",
			schemaFilePath: "./test_data/sports/boxscores/boxscoresBasketball.json",
			xmlFilePath:    "./test_data/sports/boxscores/boxscores_NBA.xml",
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

		b, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			t.Fatalf(err.Error())
		}
		os.Stdout.Write(b)
	}

}
