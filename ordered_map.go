package orderedmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

// json token类型
//	Delim, for the four JSON delimiters [ ] { }
//	bool, for JSON booleans
//	float64, for JSON numbers
//	Number, for JSON numbers
//	string, for JSON string literals
//	nil, for JSON null
// OrderedContainer目标是保持顺序，因此只需要对[]、{}中的内容按照顺序序列化、反序列化，其他的字面量直接原样序列化即可

var (
	// change to false, if you don't need json.Number to represent Numbers
	UseNumber bool = true
)

type OrderedValue struct {
	Key   string
	Value interface{}
}

type OrderedMap struct {
	Values []OrderedValue
}

type OrderedArray []interface{}

// implementation of interface Marshaler
func (om OrderedMap) MarshalJSON() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.WriteRune('{')

	for i, val := range om.Values {
		// write key
		b, e := json.Marshal(val.Key)
		if e != nil {
			return nil, e
		}
		buf.Write(b)

		// write delimiter
		buf.WriteRune(':')

		// write value
		b, e = json.Marshal(val.Value)
		if e != nil {
			return nil, e
		}
		buf.Write(b)

		// write delimiter
		if i+1 < len(om.Values) {
			buf.WriteRune(',')
		}
	}

	buf.WriteRune('}')

	return buf.Bytes(), nil
}

// implementation of interface Unmarshaler
func (om *OrderedMap) UnmarshalJSON(b []byte) error {
	d := json.NewDecoder(bytes.NewReader(b))
	if UseNumber {
		d.UseNumber()
	}

	t, err := d.Token()
	if err == io.EOF {
		return nil
	}

	// 识别是否为对象
	if t != json.Delim('{') {
		// log.Print("unexpected start of object")
		return errors.New("unexpected start of object")
	}
	return om.unmarshalEmbededObject(d)
}

func (om *OrderedMap) unmarshalEmbededObject(d *json.Decoder) error {
	for d.More() {
		kToken, err := d.Token()
		if err == io.EOF || (err == nil && kToken == json.Delim('}')) {
			// log.Print("unexpected EOF")
			return errors.New("unexpected EOF")
		}

		vToken, err := d.Token()
		if err == io.EOF {
			// log.Print("unexpected EOF")
			return errors.New("unexpected EOF")
		}

		var val interface{}
		switch vToken {
		case json.Delim('{'):
			var obj OrderedMap
			if err = obj.unmarshalEmbededObject(d); err != nil {
				return err
			}
			val = obj
		case json.Delim('['):
			var arr OrderedArray
			err = arr.unmarshalEmbededArray(d)
			val = arr
		default:
			val = vToken
		}

		if err != nil {
			return err
		}

		om.Values = append(om.Values, OrderedValue{kToken.(string), val})
	}

	// 读取对象结束token '}'
	kToken, err := d.Token()
	if err == io.EOF || kToken != json.Delim('}') {
		// log.Print("unexpected EOF")
		return errors.New("unexpected EOF")
	}

	return err
}

// implementation of interface Marshaler
func (arr OrderedArray) MarshalJSON() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.WriteRune('[')

	for i, val := range arr {
		// write key
		b, e := json.Marshal(val)
		if e != nil {
			return nil, e
		}
		buf.Write(b)

		// write delimiter
		if i+1 < len(arr) {
			buf.WriteRune(',')
		}
	}

	buf.WriteRune(']')

	return buf.Bytes(), nil
}

func (arr *OrderedArray) UnmarshalJSON(b []byte) error {
	d := json.NewDecoder(bytes.NewReader(b))
	if UseNumber {
		d.UseNumber()
	}

	t, err := d.Token()
	if err == io.EOF {
		return nil
	}

	// 识别是否为数组
	if t != json.Delim('[') {
		return errors.New("unexpected start of array")
	}
	return arr.unmarshalEmbededArray(d)
}

func (om *OrderedArray) unmarshalEmbededArray(d *json.Decoder) error {
	for d.More() {
		token, err := d.Token()
		if err == io.EOF || (err == nil && token == json.Delim(']')) {
			// log.Print("unexpected EOF")
			return errors.New("unexpected EOF")
		}

		var val interface{}
		switch token {
		case json.Delim('{'):
			var obj OrderedMap
			if err = obj.unmarshalEmbededObject(d); err != nil {
				return err
			}
			val = obj
		case json.Delim('['):
			var arr OrderedArray
			err = arr.unmarshalEmbededArray(d)
			val = arr
		default:
			// 字面量 literial
			val = token
		}

		if err != nil {
			return err
		}

		*om = append(*om, val)
	}

	// 读取数组结束token ']'
	kToken, err := d.Token()
	if err == io.EOF || kToken != json.Delim(']') {
		return errors.New("unexpected EOF")
	}

	if *om == nil {
		*om = OrderedArray(make([]interface{}, 0))
	}

	return nil
}
