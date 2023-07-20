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

// func Test_Marshal(t *testing.T) {
// 	testCases := []struct {
// 		name   string
// 		input  OrderedMap
// 		expect string
// 	}{
// 		{
// 			name: "基本正序场景",
// 			input: OrderedMap{
// 				Values: []OrderedValue{
// 					{
// 						Key:   "a",
// 						Value: 1,
// 					},
// 					{
// 						Key:   "b",
// 						Value: 2,
// 					},
// 					{
// 						Key:   "c",
// 						Value: 3,
// 					},
// 				},
// 			},
// 			expect: `{"a":1,"b":2,"c":3}`,
// 		},
// 		{
// 			name: "基本逆序场景",
// 			input: OrderedMap{
// 				Values: []OrderedValue{
// 					{
// 						Key:   "c",
// 						Value: 3,
// 					},
// 					{
// 						Key:   "b",
// 						Value: 2,
// 					},
// 					{
// 						Key:   "a",
// 						Value: 1,
// 					},
// 				},
// 			},
// 			expect: `{"c":3,"b":2,"a":1}`,
// 		},
// 		{
// 			name: "嵌套map逆序场景",
// 			input: OrderedMap{
// 				Values: []OrderedValue{
// 					{
// 						Key:   "c",
// 						Value: 3,
// 					},
// 					{
// 						Key: "b",
// 						Value: OrderedMap{
// 							Values: []OrderedValue{
// 								{
// 									Key:   "y",
// 									Value: "*",
// 								},
// 								{
// 									Key:   "z",
// 									Value: "#",
// 								},
// 								{
// 									Key:   "x",
// 									Value: "$",
// 								},
// 							},
// 						},
// 					},
// 					{
// 						Key:   "a",
// 						Value: 1,
// 					},
// 				},
// 			},
// 			expect: `{"c":3,"b":{"y":"*","z":"#","x":"$"},"a":1}`,
// 		},
// 	}
// 	for _, v := range testCases {
// 		t.Run(v.name, func(t *testing.T) {
// 			got, err := json.Marshal(v.input)
// 			if err != nil {
// 				t.Errorf("marshal failed: \n\tin: %+v\n\tex: %+v\n\terr: %+v", v.input, v.expect, err)
// 				return
// 			}

// 			if string(got) != v.expect {
// 				t.Errorf("miss match: \n\tin: %+v\n\tex: %+v\n\tgot: %+v", v.input, v.expect, string(got))
// 			}
// 		})
// 	}
// }

// func Test_Unmarshal_Marshal(t *testing.T) {
// 	testCases := []struct {
// 		name   string
// 		input  []byte
// 		expect OrderedMap
// 	}{
// 		{
// 			name:  "基本正序场景",
// 			input: []byte(`{"a":1,"b":2,"c":3}`),
// 			expect: OrderedMap{
// 				Values: []OrderedValue{
// 					{
// 						Key:   "a",
// 						Value: json.Number("1"),
// 					},
// 					{
// 						Key:   "b",
// 						Value: json.Number("2"),
// 					},
// 					{
// 						Key:   "c",
// 						Value: json.Number("3"),
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:  "基本逆序场景",
// 			input: []byte(`{"c":3,"b":2,"a":1}`),
// 			expect: OrderedMap{
// 				Values: []OrderedValue{
// 					{
// 						Key:   "c",
// 						Value: json.Number("3"),
// 					},
// 					{
// 						Key:   "b",
// 						Value: json.Number("2"),
// 					},
// 					{
// 						Key:   "a",
// 						Value: json.Number("1"),
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:  "嵌套object逆序场景",
// 			input: []byte(`{"c":{"y":"*","z":"#","x":"$"},"b":2,"a":1}`),
// 			expect: OrderedMap{
// 				Values: []OrderedValue{
// 					{
// 						Key: "c",
// 						Value: OrderedMap{
// 							Values: []OrderedValue{
// 								{
// 									Key:   "y",
// 									Value: "*",
// 								},
// 								{
// 									Key:   "z",
// 									Value: "#",
// 								},
// 								{
// 									Key:   "x",
// 									Value: "$",
// 								},
// 							},
// 						},
// 					},
// 					{
// 						Key:   "b",
// 						Value: json.Number("2"),
// 					},
// 					{
// 						Key:   "a",
// 						Value: json.Number("1"),
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:  "嵌套array场景",
// 			input: []byte(`{"c":["hello","world","hello","golang"],"b":2,"a":1}`),
// 			expect: OrderedMap{
// 				Values: []OrderedValue{
// 					{
// 						Key: "c",
// 						Value: OrderedArray{
// 							"hello",
// 							"world",
// 							"hello",
// 							"golang",
// 						},
// 					},
// 					{
// 						Key:   "b",
// 						Value: json.Number("2"),
// 					},
// 					{
// 						Key:   "a",
// 						Value: json.Number("1"),
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:  "嵌套array/object场景",
// 			input: []byte(`{"c":[{"hello":"world"},{"hello":"golang"}],"b":2,"a":1}`),
// 			expect: OrderedMap{
// 				Values: []OrderedValue{
// 					{
// 						Key: "c",
// 						Value: OrderedArray{
// 							OrderedMap{
// 								Values: []OrderedValue{
// 									{
// 										Key:   "hello",
// 										Value: "world",
// 									},
// 								},
// 							},
// 							OrderedMap{
// 								Values: []OrderedValue{
// 									{
// 										Key:   "hello",
// 										Value: "golang",
// 									},
// 								},
// 							},
// 						},
// 					},
// 					{
// 						Key:   "b",
// 						Value: json.Number("2"),
// 					},
// 					{
// 						Key:   "a",
// 						Value: json.Number("1"),
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for _, v := range testCases {
// 		t.Run(v.name, func(t *testing.T) {
// 			var got OrderedMap
// 			err := json.Unmarshal([]byte(v.input), &got)
// 			if err != nil {
// 				t.Errorf("unmarshal failed: \n\tin: %+v\n\tex: %+v\n\terr: %+v", v.input, v.expect, err)
// 				return
// 			}

// 			if !reflect.DeepEqual(v.expect, got) {
// 				t.Errorf("miss match: \n\tin: %+v\n\tex: %+v\n\tgot: %+v", v.input, v.expect, got)
// 				return
// 			}

// 			b, _ := got.MarshalJSON()
// 			if !bytes.Equal(b, v.input) {
// 				t.Errorf("miss match: \n\tin: %+v\n\tgot: %+v", string(v.input), string(b))
// 				return
// 			}
// 		})
// 	}
// }
