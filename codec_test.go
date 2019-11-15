package xmlrpc

import "testing"

//TestDecodeResponse ...
func TestDecodeResponse(t *testing.T) {
	theXML := `<methodResponse>
    <params>
        <param>
            <value>
                <struct>
                    <member>
                        <name>code</name>
                        <value>
                            <string>[code]</string>
                        </value>
                    </member>
                    <member>
                        <name>message</name>
                        <value>
                            <string>[msg]</string>
                        </value>
                    </member>
                    <member>
                        <name>refid</name>
                        <value>
                            <string>[refid]</string>
                        </value>
                    </member>
                    <member>
                        <name>trxid</name>
                        <value>
                            <string>[trxid]</string>
                        </value>
                    </member>
                </struct>
            </value>
        </param>
    </params>
</methodResponse>`
	var result Response
	Decode(theXML, &result)
	if &result == nil {
		t.Fail()
	}
	if result.Fault != nil {
		t.Fail()
	}
}

func TestDecodeResponseWithFault(t *testing.T) {
	theXML := `<methodResponse>
    <fault>
        <value>
            <struct>
                <member>
                    <name>faultCode</name>
                    <value>
                        <int>4</int>
                    </value>
                </member>
                <member>
                    <name>faultString</name>
                    <value>
                        <string>Too many parameters.</string>
                    </value>
                </member>
            </struct>
        </value>
    </fault>
</methodResponse>`
	var result Response
	Decode(theXML, &result)

	if result.Fault == nil {
		t.Fail()
	}
}
