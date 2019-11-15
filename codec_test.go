package xmlrpc

import "testing"

//TestDecodeResponse ...
func TestDecodeResponse(t *testing.T) {
	theXML := `<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
<params>
<param>
<value>
<struct>
<member>
<name>code</name>
<value>
<string>00</string>
</value>
</member>
<member>
<name>message</name>
<value>
<string>idpel:02122992515|nmpel:Aldi|periode:NOP2019|tagihan:1056699|billqty:1|admbank:2500|charge:1054799|saldo:0</string>
</value>
</member>
<member>
<name>refid</name>
<value>
<string>498081</string>
</value>
</member>
<member>
<name>trxid</name>
<value>
<string>19111503546</string>
</value>
</member>
</struct></value>
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

	var inquiryResponse InquiryResponse
	CopyAllParam(result.Params, &inquiryResponse)
	if inquiryResponse.Code != "00" {
		t.Fail()
	}

	if inquiryResponse.Message != "idpel:02122992515|nmpel:Aldi|periode:NOP2019|tagihan:1056699|billqty:1|admbank:2500|charge:1054799|saldo:0" {
		t.Fail()
	}

	if inquiryResponse.RefID != "498081" {
		t.Fail()
	}

	if inquiryResponse.TrxID != "19111503546" {
		t.Fail()
	}
}

//InquiryResponse ...
type InquiryResponse struct {
	Code    string `rpc:"code"`
	Message string `rpc:"message"`
	RefID   string `rpc:"refid"`
	TrxID   string `rpc:"trxid"`
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
