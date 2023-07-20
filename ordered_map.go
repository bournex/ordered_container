package orderedmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

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

	kToken, err := d.Token()
	if err == io.EOF || (kToken != json.Delim('}') && kToken != json.Delim(']')) {
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

	if t != json.Delim('[') {
		// log.Print("unexpected start of array")
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
			// 立即值
			val = token
		}

		if err != nil {
			return err
		}

		*om = append(*om, val)
	}

	kToken, err := d.Token()
	if err == io.EOF || kToken != json.Delim(']') {
		// log.Print("unexpected EOF")
		return errors.New("unexpected EOF")
	}

	if *om == nil {
		*om = OrderedArray(make([]interface{}, 0))
	}

	return nil
}
