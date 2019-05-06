package json

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
)

var (
	testRawJSON = json.RawMessage(`
{
	"group1": {
		"item1": {
			"itemA": "A"
		}
	},
	"group2": {
		"item1": 0,
		"item2": "two"
	},
	"group3": [
		"item1",
		"item2"
	]
}`)
	testRaw = interface{}(nil)
)

func TestMain(m *testing.M) {
	if err := json.Unmarshal(testRawJSON, &testRaw); err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestTransformInstruction(t *testing.T) {
	tests := []struct {
		description string
		ti          transformInstruction
		in          interface{}
		want        interface{}
		wantErr     bool
	}{
		{
			description: "Simple Instruction",
			ti: transformInstruction{
				jsonPath:   "$.group1.item1.itemA",
				Operations: []transformOperation{&testOp{args: map[string]string{"out": "out"}}},
			},
			in:   testRaw,
			want: "out",
		},
		{
			description: "Chained operations",
			ti: transformInstruction{
				jsonPath: "$.group1.item1.itemA",
				Operations: []transformOperation{
					&testOp{args: map[string]string{"out": "out"}},
					&testOp{args: map[string]string{"out": "out2"}},
				},
			},
			in:   testRaw,
			want: "out2",
		},
		{
			description: "Value not found",
			in:          testRaw,
			ti: transformInstruction{
				jsonPath:   "$.group1.item10.itemA",
				Operations: []transformOperation{&testOp{args: map[string]string{"out": "out"}}},
			},
			want: nil,
		},
		{
			description: "failed operation",
			in:          testRaw,
			ti: transformInstruction{
				jsonPath:   "$.group1.item1.itemA",
				Operations: []transformOperation{&testOp{fail: true}},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		got, err := test.ti.transform(test.in, "string", nil)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil, want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error, want nil: %v", test.description, err)
		case !reflect.DeepEqual(got, test.want):
			t.Errorf("Test %q - got %v, want %v", test.description, got, test.want)
		}
	}
}

func TestTransformInstructions(t *testing.T) {

	tests := []struct {
		description string
		tis         transformInstructions
		in          interface{}
		want        interface{}
		wantErr     bool
	}{
		{
			description: "Basic - method first",
			tis: transformInstructions{
				From: []*transformInstruction{
					{
						jsonPath:   "$.group1.item1.itemA",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out"}}},
					},
				},
				Method: first,
			},
			in:   testRaw,
			want: "out",
		},
		{
			description: "multiple instructions - method first",
			tis: transformInstructions{
				From: []*transformInstruction{
					{
						jsonPath:   "$.group1.item1.itemA",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out"}}},
					},
					{
						jsonPath:   "$.group3[1]",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out2"}}},
					},
				},
				Method: first,
			},
			in:   testRaw,
			want: "out",
		},
		{
			description: "multiple instructions - method last",
			tis: transformInstructions{
				From: []*transformInstruction{
					{
						jsonPath:   "$.group1.item1.itemA",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out"}}},
					},
					{
						jsonPath:   "$.group3[1]",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out2"}}},
					},
				},
				Method: last,
			},
			in:   testRaw,
			want: "out2",
		},
		{
			description: "multiple instructions - method Concat",
			tis: transformInstructions{
				From: []*transformInstruction{
					{
						jsonPath:   "$.group1.item1.itemA",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out"}}},
					},
					{
						jsonPath:   "$.group3[1]",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out2"}}},
					},
				},
				Method: concatenate,
			},
			in:   testRaw,
			want: "outout2",
		},
		{
			description: "multiple instructions - method Concat with delimiter",
			tis: transformInstructions{
				From: []*transformInstruction{
					{
						jsonPath:   "$.group1.item1.itemA",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out"}}},
					},
					{
						jsonPath:   "$.group3[1]",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out2"}}},
					},
				},
				Method: concatenate,
				MethodOptions: methodOptions{
					ConcatenateDelimiter: "/",
				},
			},
			in:   testRaw,
			want: "out/out2",
		},
		{
			description: "multiple instructions - method Concat with delimiter, one path missing",
			tis: transformInstructions{
				From: []*transformInstruction{
					{
						jsonPath:   "$.group1.item1.itemA",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out"}}},
					},
					{
						jsonPath: "$.group3[5]",
					},
					{
						jsonPath:   "$.group3[1]",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out2"}}},
					},
				},
				Method: concatenate,
				MethodOptions: methodOptions{
					ConcatenateDelimiter: "/",
				},
			},
			in:   testRaw,
			want: "out/out2",
		},
		{
			description: "multiple instructions - method Concat with delimiter, all paths missing",
			tis: transformInstructions{
				From: []*transformInstruction{
					{
						jsonPath: "$.group1.item1.itemF",
					},
					{
						jsonPath: "$.group3[5]",
					},
					{
						jsonPath: "$.group3[10]",
					},
				},
				Method: concatenate,
				MethodOptions: methodOptions{
					ConcatenateDelimiter: "/",
				},
			},
			in:   testRaw,
			want: nil,
		},
		{
			description: "all paths are missing",
			tis: transformInstructions{
				From: []*transformInstruction{
					{
						jsonPath:   "$.group10.item1.itemA",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out"}}},
					},
					{
						jsonPath:   "$.group30[1]",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out2"}}},
					},
				},
				Method: first,
			},
			in:   testRaw,
			want: nil,
		},
		{
			description: "one paths is missing",
			tis: transformInstructions{
				From: []*transformInstruction{
					{
						jsonPath:   "$.group10.item1.itemA",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out"}}},
					},
					{
						jsonPath:   "$.group3[1]",
						Operations: []transformOperation{&testOp{args: map[string]string{"out": "out2"}}},
					},
				},
				Method: first,
			},
			in:   testRaw,
			want: "out2",
		},
		{
			description: "failed operation",
			tis: transformInstructions{
				From: []*transformInstruction{
					{
						jsonPath:   "$.group1.item1.itemA",
						Operations: []transformOperation{&testOp{fail: true}},
					},
				},
				Method: first,
			},
			in:      testRaw,
			wantErr: true,
		},
	}

	for _, test := range tests {
		got, err := test.tis.transform(test.in, "string", nil)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil, want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error, want nil: %v", test.description, err)
		case !reflect.DeepEqual(got, test.want):
			t.Errorf("Test %q - got %v, want %v", test.description, got, test.want)
		}
	}
}

func TestTransformUnmarshal(t *testing.T) {
	tests := []struct {
		description string
		value       json.RawMessage
		want        transform
		wantErr     bool
	}{
		{
			description: "Missing transform",
			value:       []byte(`{}`),
			want:        transform{},
		},
		{
			description: "Basic transform",
			value: []byte(`
{
	"cumulo": {
		"from": [
			{
				"jsonPath": "$.data.type"
			}
		]
	}
}`,
			),
			want: transform{"cumulo": transformInstructions{
				From: []*transformInstruction{
					{jsonPath: "$.data.type", Operations: []transformOperation{}},
				},
				Method: first,
			},
			},
		},
		{
			description: "Basic transform, last method",
			value: []byte(`
{
	"cumulo": {
		"from": [
			{
				"jsonPath": "$.data.type"
			}
		],
		"method": "last"
	}
}`,
			),
			want: transform{"cumulo": transformInstructions{
				From: []*transformInstruction{
					{jsonPath: "$.data.type", Operations: []transformOperation{}},
				},
				Method: last,
			},
			},
		},
		{
			description: "Basic transform, concatenate method",
			value: []byte(`
{
	"cumulo": {
		"from": [
			{
				"jsonPath": "$.data.type"
			}
		],
		"method": "concatenate"
	}
}`,
			),
			want: transform{"cumulo": transformInstructions{
				From: []*transformInstruction{
					{jsonPath: "$.data.type", Operations: []transformOperation{}},
				},
				Method: concatenate,
			},
			},
		},
		{
			description: "Basic transform, concatenate method with delimiter",
			value: []byte(`
{
	"cumulo": {
		"from": [
			{
				"jsonPath": "$.data.type"
			}
		],
		"method": "concatenate",
		"methodOptions": {
			"concatenateDelimiter": "/"
		}
	}
}`,
			),
			want: transform{"cumulo": transformInstructions{
				From: []*transformInstruction{
					{jsonPath: "$.data.type", Operations: []transformOperation{}},
				},
				Method: concatenate,
				MethodOptions: methodOptions{
					ConcatenateDelimiter: "/",
				},
			},
			},
		},
		{
			description: "Two options",
			value: []byte(`
{
	"cumulo": {
		"from": [
			{
				"jsonPath": "$.data.mobileBody[*]"
			}
		]
	},
	"presentationv4": {
		"from": [
			{
				"jsonPath": "$.mobileBody[*]"
			}
		]
	}
}`,
			),
			want: transform{
				"cumulo": transformInstructions{
					From: []*transformInstruction{
						{jsonPath: "$.data.mobileBody[*]", Operations: []transformOperation{}},
					},
					Method: first,
				},
				"presentationv4": transformInstructions{
					From: []*transformInstruction{
						{jsonPath: "$.mobileBody[*]", Operations: []transformOperation{}},
					},
					Method: first,
				},
			},
		},
		{
			description: "Many from",
			value: []byte(`
{
	"presentationv4": {
		"from": [
			{
				"jsonPath": "$.associatedAssetId"
			},
			{
				"jsonPath": "$._attributes.AssociatedAssetId"
			},
			{
				"jsonPath": "$._attributes.associatedassetid"
			}
		]
	}
}`,
			),
			want: transform{
				"presentationv4": transformInstructions{
					From: []*transformInstruction{
						{jsonPath: "$.associatedAssetId", Operations: []transformOperation{}},
						{jsonPath: "$._attributes.AssociatedAssetId", Operations: []transformOperation{}},
						{jsonPath: "$._attributes.associatedassetid", Operations: []transformOperation{}},
					},
					Method: first,
				},
			},
		},
		{
			description: "Advanced test, with operations",
			value: []byte(`
{
	"cumulo": {
		"from": [
			{
				"jsonPath": "$.data.renditions[*]",
			  	"operations": [
					{
				  		"type": "max",
				  		"args" : {
							"by": "@.encodingRate",
							"return": "@.url"
				  		}
					},
					{
						"type": "replace",
					  	"args": {
							"regex": "(http:\/\/.*net)\/",
							"new": "https://media.gannett-cdn.com"
					  	}
					}
			  	]
			}
		]
	},
	"presentationv4": {
	  	"from": [
		  	{
				"jsonPath": "$.renditions[*]",
				"operations": [
				  	{
						"type": "changeCase",
						"args" : {
					  		"to": "lower"
						}
				  	},
				  	{
						"type": "inverse"
				  	},
				  	{
						"type": "split",
						"args": {
					  		"on": "|"
						}
				  	}
				]
		  	}
	  	]
	}
}`,
			),
			want: transform{
				"cumulo": transformInstructions{
					From: []*transformInstruction{
						{jsonPath: "$.data.renditions[*]", Operations: []transformOperation{
							&max{Args: map[string]string{"by": "@.encodingRate", "return": "@.url"}},
							&replace{Args: map[string]string{"regex": `(http://.*net)/`, "new": "https://media.gannett-cdn.com"}},
						}},
					},
					Method: first,
				},
				"presentationv4": transformInstructions{
					From: []*transformInstruction{
						{jsonPath: "$.renditions[*]", Operations: []transformOperation{
							&changeCase{Args: map[string]string{"to": "lower"}},
							&inverse{},
							&split{Args: map[string]string{"on": "|"}},
						}},
					},
					Method: first,
				},
			},
		},
	}

	for _, test := range tests {
		var tr transform

		err := json.Unmarshal(test.value, &tr)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil error want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error: %v", test.description, err)
		case !reflect.DeepEqual(tr, test.want):
			if got, want := len(tr), len(test.want); got != want {
				t.Errorf("Test %q - got transform length %d, want %d", test.description, got, want)
			}
			for key, value := range tr {
				if got, want := value, test.want[key]; !reflect.DeepEqual(got, want) {
					gotJSON, err := json.Marshal(got)
					if err != nil {
						t.Errorf("Test %q - failed to marshal response for key %q: %v", test.description, key, err)
					}
					wantJSON, err := json.Marshal(want)
					if err != nil {
						t.Errorf("Test %q - failed to marshal want for key %q: %v", test.description, key, err)
					}
					if string(gotJSON) != string(wantJSON) {
						t.Errorf("Test %q - transform key %q, got\n%s\nwant\n%s", test.description, key, gotJSON, wantJSON)
					}
				}
			}
		}
	}
}
