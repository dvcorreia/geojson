package geojson

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/tidwall/geojson/geometry"
)

func TestMultiLineString(t *testing.T) {
	expectJSON(t, `{"type":"MultiLineString","coordinates":[[[1,2,3]]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiLineString","coordinates":[[[1,2]]],"bbox":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiLineString"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiLineString","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiLineString","coordinates":[1,null]}`, errCoordinatesInvalid)
}

func TestMultiLineStringUnmarshal(t *testing.T) {
	expectUnmarshalMultiLineString(t, `{"type":"MultiLineString","coordinates":[[[1,2,3]]]}`, nil, errCoordinatesInvalid)
	expectUnmarshalMultiLineString(t, `{"type":"MultiLineString","coordinates":[[[1,2]]],"bbox":null}`, nil, errCoordinatesInvalid)
	expectUnmarshalMultiLineString(t, `{"type":"MultiLineString"}`, nil, errCoordinatesMissing)
	expectUnmarshalMultiLineString(t, `{"type":"MultiLineString","coordinates":null}`, nil, errCoordinatesInvalid)
	expectUnmarshalMultiLineString(t, `{"type":"MultiLineString","coordinates":[1,null]}`, nil, errCoordinatesInvalid)

	json := `{"type":"MultiLineString","coordinates":[
		[[10,10],[120,190]],
		[[50,50],[100,100]]
	]}`
	tc := &MultiLineString{collection: collection{
		children: []Object{
			NewLineString(geometry.NewLine([]geometry.Point{
				{X: 10, Y: 10},
				{X: 120, Y: 190},
			}, nil)),
			NewLineString(geometry.NewLine([]geometry.Point{
				{X: 50, Y: 50},
				{X: 100, Y: 100},
			}, nil)),
		},
	}}
	tc.parseInitRectIndex(DefaultParseOptions)
	expectUnmarshalMultiLineString(t, json, tc, nil)
}

func TestMultiLineStringValid(t *testing.T) {
	json := `{"type":"MultiLineString","coordinates":[
		[[10,10],[120,190]],
		[[50,50],[100,100]]
	]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true})
}

func TestMultiLineStringPoly(t *testing.T) {
	p := expectJSON(t, `{"type":"MultiLineString","coordinates":[
		[[10,10],[20,20]],
		[[50,50],[100,100]]
	]}`, nil)
	expect(t, p.Intersects(PO(15, 15)))
	expect(t, p.Contains(PO(15, 15)))
	expect(t, p.Contains(PO(70, 70)))
	expect(t, !p.Contains(PO(40, 40)))
}

func expectUnmarshalMultiLineString(t *testing.T, input string, expected *MultiLineString, exerr error) {
	var mls MultiLineString
	err := json.Unmarshal([]byte(input), &mls)
	if err != exerr {
		t.Fatalf("expected error '%v', got '%v'", exerr, err)
	}

	if exerr != nil {
		return
	}

	if !reflect.DeepEqual(expected, &mls) {
		t.Fatalf("expected '%#v', got '%#v'", expected, &mls)
	}
}
