package main

import (
	"fmt"

	"github.com/fajrirahmat/xmlrpc-codec/codec"
)

type tag struct {
	Name  string
	Value string
}

type example struct {
	Who string `rpc:"who"`
}

type object struct {
	Bool   bool
	String string `rpc:"user"`
	Reffid string
	Int    int
	Int64  int64
	Double float64
	Base64 []byte
}

var userXML = `<param>
			<value>
				<i4>193949</i4>
			</value>
		</param>
		<param>
			<value>
				<int>10</int>
			</value>
		</param>
		<param>
			<value>
				<double>10.4</double>
			</value>
		</param>
		<param>
			<value>
				<boolean>true</boolean>
			</value>
		</param>
		<param>
			<value>
				<string>makan apa?</string>
			</value>
		</param>
		<param>
			<value>
				<dateTime.iso8601>19980717T14:08:55</dateTime.iso8601>
			</value>
		</param>`

var theString = `<methodCall>
	<methodName>Fajri.RPC</methodName>
	<params>
		<param>
			<value>
				<string>makan apa?</string>
			</value>
		</param>
	</params>
</methodCall>
`

var theBytes = `<methodCall>
	<methodName>Fajri.RPC</methodName>
	<params>
		<param>
			<value>
				<base64>bWFrYW4gYXBh</base64>
			</value>
		</param>
	</params>
</methodCall>
`

var theXML = `<methodCall>
    <methodName>Fajri.RPC</methodName>
	<params>
		<param>
			<value>
				<string>makan apa?</string>
			</value>
		</param>
        <param>
            <value>
                <struct>
                    <member>
                        <name>user</name>
                        <value>
                            <string>[user]</string>
                        </value>
                    </member>
                    <member>
                        <name>waktu</name>
                        <value>
                            <string>[time]</string>
                        </value>
                    </member>
                    <member>
                        <name>produk</name>
                        <value>
                            <string>[produk]</string>
                        </value>
                    </member>
                    <member>
                        <name>idpel</name>
                        <value>
                            <string>[idpel]</string>
                        </value>
                    </member>
                    <member>
                        <name>hppel</name>
                        <value>
                            <string>[hppel]</string>
                        </value>
                    </member>
                    <member>
                        <name>reffid</name>
                        <value>
                            <string>[trxid]</string>
                        </value>
                    </member>
                    <member>
                        <name>signature</name>
                        <value>
                            <string>[sign]</string>
                        </value>
					</member>
					<member>
                        <name>Reffid</name>
                        <value>
                            <string>[trxid]</string>
                        </value>
                    </member>
                </struct>
            </value>
		</param>
		<param>
			<value>
				<base64>bWFrYW4gYXBh</base64>
			</value>
		</param>
    </params>
</methodCall>`

func main() {
	theXMLObj := codec.DecodeRequest(theXML)
	obj := object{}
	var a string
	codec.CopyAllParam(*theXMLObj, &a, &obj)
	fmt.Println(obj)
	fmt.Println(a)

	buf, err := codec.EncodeRequest("Testing.hello", "testing", true, object{Bool: true, String: "testing"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(buf))

}

/*
func valueOpo(val reflect.Value, obj *object) {
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr && !val.IsNil() {
		if val.Elem().IsValid() && val.Elem().Kind() == reflect.Int64 {
			codec.CopyValue(val.Elem(), &obj.Int64)
		}
		if val.Elem().IsValid() && val.Elem().Kind() == reflect.Int {
			codec.CopyValue(val.Elem(), &obj.Int)
		}
		if val.Elem().IsValid() && val.Elem().Kind() == reflect.Float64 {
			//fmt.Println(val.Elem().Float())
			codec.CopyValue(val.Elem(), &obj.Double)
		}
		if val.Elem().IsValid() && val.Elem().Kind() == reflect.Bool {
			codec.CopyValue(val.Elem(), &obj.Bool)
		}
		if val.Elem().IsValid() && val.Elem().Kind() == reflect.String {
			codec.CopyValue(val.Elem(), &obj.String)
		}
	} else if val.Kind() == reflect.Slice && val.Len() > 0 {
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			if v.Kind() == reflect.Struct {
				//assume always member
				for j := 0; j < v.NumField(); j++ {
					vfield := v.Field(j)
					if vfield.Kind() == reflect.String {
						//	fmt.Println(vfield.String())
					} else if vfield.Kind() == reflect.Struct {
						for k := 0; k < vfield.NumField(); k++ {
							val2 := vfield.Field(k)
							valueOpo(val2, obj)
						}
					}
				}
			}
		}
	}
}*/
