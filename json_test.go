package parcon_test

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Fuzz_json(f *testing.F) {
	var tests = []string{
		`0`,
		`1`,
		`123`,
		`123.456`,
		`1e-10`,
		`5e+3`,
		`9E123`,
		`-1`,
		`-123`,
		`-123.456`,
		`""`,
		`"hello world"`,
		`"hello\"\t\r\nworld"`,
		`"\u3042\u4e9C"`,
		`"\uD800\udc00"`,
		`"\uD801\ud801"`,
		`"\udc02\uDC02"`,
		` true `,
		`false`,
		`null`,
		`[]`,
		`["hello"]`,
		`["hello","world"]`,
		`[ "hello" , "world" ]`,
		`{}`,
		`{"hello":"world"}`,
		`{"hello":"world","foo":"bar"}`,
		`{ "hello" : "world", "foo" : "bar" }`,
		`{"hello": 123, "foo": 456}`,
		`{"hello": [123, true], "foo": ["bar", false]}`,
		`[null, {"hello": "world"}]`,
		`{
			"Image": {
				"Width":  800,
				"Height": 600,
				"Title":  "View from 15th Floor",
				"Thumbnail": {
					"Url":    "http://www.example.com/image/481989943",
					"Height": 125,
					"Width":  100
				},
				"Animated" : false,
				"IDs": [116, 943, 234, 38793]
			}
		}`,
		`[
			{
				"precision": "zip",
				"Latitude":  37.7668,
				"Longitude": -122.3959,
				"Address":   "",
				"City":      "SAN FRANCISCO",
				"State":     "CA",
				"Zip":       "94107",
				"Country":   "US"
			},
			{
				"precision": "zip",
				"Latitude":  37.371991,
				"Longitude": -122.026020,
				"Address":   "",
				"City":      "SUNNYVALE",
				"State":     "CA",
				"Zip":       "94085",
				"Country":   "US"
			}
		]`,
	}

	for _, tt := range tests {
		f.Add(string(tt))
	}

	f.Fuzz(func(t *testing.T, input string) {
		var want interface{}
		shouldBeError := false
		if json.Unmarshal([]byte(input), &want) != nil {
			shouldBeError = true
		}

		got, err := ParseJson(input)

		if shouldBeError {
			if err == nil {
				t.Fatalf("should be error but succeed to parse")
			}
		} else {
			if err != nil {
				t.Fatalf("failed to parse: %s", err)
			}

			if !reflect.DeepEqual(want, got) {
				t.Errorf("unexpected output\nwant: %#v\n got: %#v", want, got)
			}
		}
	})
}
