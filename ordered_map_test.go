package orderedmap

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_OrderedMap(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect OrderedMap
	}{
		{
			name:   "empty map",
			input:  `{}`,
			expect: OrderedMap{},
		},
		{
			name:  "simple map",
			input: `{"name":"alice","age":5}`,
			expect: OrderedMap{
				Values: []OrderedValue{
					{
						Key:   "name",
						Value: "alice",
					},
					{
						Key:   "age",
						Value: json.Number("5"),
					},
				},
			},
		},
		{
			name:  "nested map",
			input: `{"girls":{"name":"alice","age":5},"boys":{"name":"bob","age":3}}`,
			expect: OrderedMap{
				Values: []OrderedValue{
					{
						Key: "girls",
						Value: OrderedMap{
							Values: []OrderedValue{
								{
									Key:   "name",
									Value: "alice",
								},
								{
									Key:   "age",
									Value: json.Number("5"),
								},
							},
						},
					},
					{
						Key: "boys",
						Value: OrderedMap{
							Values: []OrderedValue{
								{
									Key:   "name",
									Value: "bob",
								},
								{
									Key:   "age",
									Value: json.Number("3"),
								},
							},
						},
					},
				},
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			var om OrderedMap
			err := json.Unmarshal([]byte(v.input), &om)
			if err != nil {
				t.Errorf("unmarshal failed\n\tinput %+v\n\texpect %+v\n\terror %+v", v.input, v.expect, err)
				return
			}

			if !reflect.DeepEqual(om, v.expect) {
				t.Errorf("miss match object\n\tinput %+v\n\texpect %+v\n\tgot %+v", v.input, v.expect, om)
				return
			}

			b, err := json.Marshal(om)
			if err != nil {
				t.Errorf("unmarshal failed\n\tinput %+v\n\texpect %+v\n\terror %+v", v.input, v.expect, err)
				return
			}

			if string(b) != v.input {
				t.Errorf("miss match stream\n\tinput %+v\n\texpect %+v\n\tgot %+v", v.input, v.expect, om)
				return
			}
		})
	}
}

func Test_OrderedArray(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect OrderedArray
	}{
		{
			name:   "empty array",
			input:  "[]",
			expect: OrderedArray{},
		},
		{
			name:  "simple array",
			input: `["hello","world","alice",5]`,
			expect: OrderedArray{
				"hello",
				"world",
				"alice",
				json.Number("5"),
			},
		},
		{
			name:  "nested array",
			input: `[[[["alice"],["bob"]]]]`,
			expect: OrderedArray{
				OrderedArray{
					OrderedArray{
						OrderedArray{
							"alice",
						},
						OrderedArray{
							"bob",
						},
					},
				},
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			var oa OrderedArray
			err := json.Unmarshal([]byte(v.input), &oa)
			if err != nil {
				t.Errorf("unmarshal failed\n\tinput %+v\n\texpect %+v\n\terror %+v", v.input, v.expect, err)
				return
			}

			if !reflect.DeepEqual(oa, v.expect) {
				t.Errorf("miss match object\n\tinput %+v\n\texpect %+v\n\tgot %+v", v.input, v.expect, oa)
				return
			}

			b, err := json.Marshal(oa)
			if err != nil {
				t.Errorf("unmarshal failed\n\tinput %+v\n\texpect %+v\n\terror %+v", v.input, v.expect, err)
				return
			}

			if string(b) != v.input {
				t.Errorf("miss match stream\n\tinput %+v\n\texpect %+v\n\tgot %+v", v.input, v.expect, oa)
				return
			}
		})
	}
}

func Test_Mix(t *testing.T) {
	const (
		RootIsObject = 1
		RootIsArray  = 2
	)

	testCases := []struct {
		name     string
		rootType int // RootIsObject/RootIsArray
		input    string
		expect   interface{}
	}{
		{
			rootType: RootIsObject,
			input:    `{"a":1,"b":[true,false],"c":3.14}`,
			expect: OrderedMap{
				Values: []OrderedValue{
					{
						Key:   "a",
						Value: json.Number("1"),
					},
					{
						Key: "b",
						Value: OrderedArray{
							true,
							false,
						},
					},
					{
						Key:   "c",
						Value: json.Number("3.14"),
					},
				},
			},
		},
		{
			rootType: RootIsArray,
			input:    `[{"name":"alice","age":5},{"name":"bob","age":3}]`,
			expect: OrderedArray{
				OrderedMap{
					Values: []OrderedValue{
						{
							Key:   "name",
							Value: "alice",
						},
						{
							Key:   "age",
							Value: json.Number("5"),
						},
					},
				},
				OrderedMap{
					Values: []OrderedValue{
						{
							Key:   "name",
							Value: "bob",
						},
						{
							Key:   "age",
							Value: json.Number("3"),
						},
					},
				},
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			if v.rootType == RootIsObject {
				testObjectRoot(t, v.input, v.expect)
			} else if v.rootType == RootIsArray {
				testArrayRoot(t, v.input, v.expect)
			}
		})
	}
}

func testObjectRoot(t *testing.T, input string, expect interface{}) {
	var oa OrderedMap
	err := json.Unmarshal([]byte(input), &oa)
	if err != nil {
		t.Errorf("unmarshal failed\n\tinput %+v\n\texpect %+v\n\terror %+v", input, expect, err)
		return
	}

	if !reflect.DeepEqual(oa, expect) {
		t.Errorf("miss match object\n\tinput %+v\n\texpect %+v\n\tgot %+v", input, expect, oa)
		return
	}

	b, err := json.Marshal(oa)
	if err != nil {
		t.Errorf("unmarshal failed\n\tinput %+v\n\texpect %+v\n\terror %+v", input, expect, err)
		return
	}

	if string(b) != input {
		t.Errorf("miss match stream\n\tinput %+v\n\texpect %+v\n\tgot %+v", input, expect, oa)
		return
	}
}

func testArrayRoot(t *testing.T, input string, expect interface{}) {
	var oa OrderedArray
	err := json.Unmarshal([]byte(input), &oa)
	if err != nil {
		t.Errorf("unmarshal failed\n\tinput %+v\n\texpect %+v\n\terror %+v", input, expect, err)
		return
	}

	if !reflect.DeepEqual(oa, expect) {
		t.Errorf("miss match object\n\tinput %+v\n\texpect %+v\n\tgot %+v", input, expect, oa)
		return
	}

	b, err := json.Marshal(oa)
	if err != nil {
		t.Errorf("unmarshal failed\n\tinput %+v\n\texpect %+v\n\terror %+v", input, expect, err)
		return
	}

	if string(b) != input {
		t.Errorf("miss match stream\n\tinput %+v\n\texpect %+v\n\tgot %+v", input, expect, oa)
		return
	}
}
