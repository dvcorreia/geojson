package geojson

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/tidwall/geojson/geometry"
)

func TestMultiPoint(t *testing.T) {
	p := expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2,3]]}`, nil)
	expect(t, p.Center() == P(1, 2))
	expectJSON(t, `{"type":"MultiPoint","coordinates":[1,null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2]],"bbox":null}`, nil)
	expectJSON(t, `{"type":"MultiPoint"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiPoint","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2,3],[4,5,6]],"bbox":[1,2,3,4]}`, nil)
}

// func TestMultiPointPoly(t *testing.T) {
// 	p := expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2],[2,2]]}`, nil)
// 	expect(t, p.Intersects(PO(1, 2)))
// 	expect(t, p.Contains(PO(1, 2)))
// 	expect(t, p.Contains(PO(2, 2)))
// 	expect(t, !p.Contains(PO(3, 2)))
// }

func TestMultiPointUnmarshal(t *testing.T) {
	tc := &MultiPoint{collection: collection{
		children: []Object{
			&Point{base: geometry.Point{X: 1, Y: 2}, extra: &extra{dims: 1, values: []float64{3}}},
		},
	}}
	tc.parseInitRectIndex(DefaultParseOptions)
	expectUnmarshalMultiPoint(t, `{"type":"MultiPoint","coordinates":[[1,2,3]]}`, tc, nil)
	expectUnmarshalMultiPoint(t, `{"type":"MultiPoint","coordinates":[1,null]}`, nil, errCoordinatesInvalid)

	tc = &MultiPoint{collection: collection{
		children: []Object{
			&Point{base: geometry.Point{X: 1, Y: 2}},
		},
		extra: &extra{members: `{"bbox":null}`},
	}}
	tc.parseInitRectIndex(DefaultParseOptions)
	expectUnmarshalMultiPoint(t, `{"type":"MultiPoint","coordinates":[[1,2]],"bbox":null}`, tc, nil)
	expectUnmarshalMultiPoint(t, `{"type":"MultiPoint"}`, nil, errCoordinatesMissing)
	expectUnmarshalMultiPoint(t, `{"type":"MultiPoint","coordinates":null}`, nil, errCoordinatesInvalid)

	tc = &MultiPoint{collection: collection{
		children: []Object{
			&Point{base: geometry.Point{X: 1, Y: 2}, extra: &extra{dims: 1, values: []float64{3}}},
			&Point{base: geometry.Point{X: 4, Y: 5}, extra: &extra{dims: 1, values: []float64{6}}},
		},
		extra: &extra{members: `{"bbox":[1,2,3,4]}`},
	}}
	tc.parseInitRectIndex(DefaultParseOptions)
	expectUnmarshalMultiPoint(t, `{"type":"MultiPoint","coordinates":[[1,2,3],[4,5,6]],"bbox":[1,2,3,4]}`, tc, nil)

}

func expectUnmarshalMultiPoint(t *testing.T, input string, expected *MultiPoint, exerr error) {
	var mp MultiPoint
	err := json.Unmarshal([]byte(input), &mp)
	if err != exerr {
		t.Fatalf("expected error '%v', got '%v'", exerr, err)
	}

	if exerr != nil {
		return
	}

	if !reflect.DeepEqual(expected, &mp) {
		t.Fatalf("expected '%#v', got '%#v'", expected, &mp)
	}
}
