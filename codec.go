package xmlrpc

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"reflect"
	"strconv"
)

type (
	//MethodCall ServerRequest ...
	MethodCall struct {
		XMLName    xml.Name `xml:"methodCall"`
		MethodName string   `xml:"methodName"`
		Params     []Param  `xml:"params>param"`
	}
	//Response ...
	Response struct {
		XMLName xml.Name `xml:"methodResponse"`
		Params  []Param  `xml:"params>param"`
		Fault   *Fault   `xml:"fault,omitempty"`
	}
	//Param ...
	Param struct {
		Value Value `xml:"value"`
	}

	//Fault ...
	Fault struct {
		Value Value `xml:"value"`
	}
	//Value ...
	Value struct {
		I4       string  `xml:"i4,omitempty"`
		Int      string  `xml:"int,omitempty"`
		Double   string  `xml:"double,omitempty"`
		Boolean  string  `xml:"boolean,omitempty"`
		String   string  `xml:"string,omitempty"`
		DateTime string  `xml:"dateTime.iso8601,omitempty"`
		Base64   string  `xml:"base64,omitempty"`
		Struct   *Struct `xml:"struct,omitempty"`
	}
	//Struct ...
	Struct struct {
		//XMLName xml.Name `xml:"struct"`
		Member []Member `xml:"member,omitempty"`
	}

	//Member ...
	Member struct {
		Name  string `xml:"name,omitempty"`
		Value Value  `xml:"value,omitempty"`
	}
)

//EncodeResponse ...
func EncodeResponse(args ...interface{}) ([]byte, error) {
	response := Response{}
	params, _ := Encode(args)
	response.Params = params
	return xml.Marshal(response)
}

//EncodeFault ...
func EncodeFault(val Value) ([]byte, error) {
	response := Response{}
	f := &Fault{Value: val}
	response.Fault = f
	return xml.Marshal(response)
}

//Encode ...
func Encode(args ...interface{}) ([]Param, error) {
	var params []Param
	//Do other things
	for _, arg := range args {
		argVal := reflect.ValueOf(arg)
		argVal = getPtrValue(argVal)

		if argVal.Kind() == reflect.String {
			param := Param{
				Value: Value{
					String: argVal.String(),
					Struct: nil,
				},
			}
			params = append(params, param)
		} else if argVal.Kind() == reflect.Bool {
			boolStr := strconv.FormatBool(argVal.Bool())
			param := Param{
				Value: Value{
					Boolean: boolStr,
					Struct:  nil,
				},
			}
			params = append(params, param)
		} else if argVal.Kind() == reflect.Struct {
			param := Param{
				Value: EncodeStruct(argVal),
			}
			params = append(params, param)
		}
	}
	return params, nil
}

//EncodeRequest ...
func EncodeRequest(method string, args ...interface{}) ([]byte, error) {
	m := MethodCall{}
	if method == "" {
		return nil, errors.New("Method name required")
	}
	m.MethodName = method
	params, _ := Encode(args)
	m.Params = params
	return xml.Marshal(m)
}

//EncodeStruct ...
func EncodeStruct(arg reflect.Value) Value {
	val := Value{}
	argVal := getPtrValue(arg)
	argType := argVal.Type()
	var members []Member
	for i := 0; i < argType.NumField(); i++ {
		field := argVal.Field(i)
		ft := argType.Field(i)
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if isZero(field) {
			continue
		}

		tag, ok := ft.Tag.Lookup("rpc")
		if ok {
			members = append(members, Member{Name: tag, Value: getValue(field)})
		} else {
			members = append(members, Member{Name: ft.Name, Value: getValue(field)})
		}
	}
	val.Struct = &Struct{Member: members}
	return val
}

/*
https://stackoverflow.com/questions/23555241/how-to-get-zero-value-of-a-field-type
*/
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

func getValue(val reflect.Value) Value {
	if val.Kind() == reflect.Bool {
		return Value{
			Boolean: strconv.FormatBool(val.Bool()),
		}
	} else if val.Kind() == reflect.Float32 || val.Kind() == reflect.Float64 {
		return Value{
			Double: strconv.FormatFloat(val.Float(), 'E', -1, 64),
		}
	} else if val.Kind() == reflect.Int || val.Kind() == reflect.Int16 || val.Kind() == reflect.Int32 || val.Kind() == reflect.Int8 {
		return Value{
			Int: strconv.FormatInt(val.Int(), 32),
		}
	} else if val.Kind() == reflect.String {
		return Value{
			String: val.String(),
		}
	} else if val.Kind() == reflect.Int64 {
		return Value{
			I4: strconv.FormatInt(val.Int(), 32),
		}
	}
	return Value{}
}

func getPtrValue(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Ptr {
		return val.Elem()
	}
	return val
}

//Decode ...
func Decode(rawXML string, result interface{}) {
	xml.Unmarshal([]byte(rawXML), result)
}

//CopyAllParam ...
func CopyAllParam(m MethodCall, args ...interface{}) {
	/*if len(m.Params) > len(args) {
		return
	}*/
	for i, v := range m.Params {
		if i >= len(args) {
			break
		}
		val := reflect.ValueOf(args[i])
		CopyValue(v.Value, val)
	}
}

//CopyValue ...
func CopyValue(value Value, field reflect.Value) {
	var el reflect.Value
	if field.Kind() == reflect.Ptr {
		el = field.Elem()
	} else {
		el = field
	}
	if !el.CanSet() {
		return
	}
	ft := el.Type()
	switch {
	case value.Boolean != "":
		var boolValue bool
		boolValue, err := strconv.ParseBool(value.Boolean)
		if err != nil {
			return
		}
		if el.Kind() == reflect.Bool {
			el.SetBool(boolValue)
		}
	case value.String != "":
		if el.Kind() == reflect.String {
			el.SetString(value.String)
		}
	case value.Int != "":
		if el.Kind() == reflect.Int {
			intValue, err := strconv.ParseInt(value.Int, 10, 32)
			if err != nil {
				return
			}
			el.SetInt(intValue)
		}
	case value.Base64 != "":
		bytes, err := base64.StdEncoding.DecodeString(value.Base64)
		if err != nil {
			return
		}
		el.SetBytes(bytes)
	case len(value.Struct.Member) > 0:
		if ft.Kind() != reflect.Struct {
			return
		}
		for _, member := range value.Struct.Member {
			for i := 0; i < ft.NumField(); i++ {
				ftf := ft.Field(i)
				ftv := el.Field(i)
				if ftf.Tag.Get("rpc") == member.Name {
					CopyValue(member.Value, ftv)
				} else if ftf.Name == member.Name {
					CopyValue(member.Value, ftv)
				}
			}
		}
	}
}
