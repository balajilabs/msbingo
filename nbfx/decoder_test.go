package nbfx

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
)

//https://golang.org/pkg/testing/

func TestDecodeExampleEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x01},
		"<doc></doc>")
}

func TestDecodeExampleComment(t *testing.T) {
	testDecode(t,
		[]byte{0x02, 0x07, 0x63, 0x6F, 0x6D, 0x6D, 0x65, 0x6E, 0x74},
		"<!--comment-->")
}

func TestDecodeExampleArray(t *testing.T) {
	testDecode(t,
		[]byte{0x03, 0x40, 0x03, 0x61, 0x72, 0x72, 0x01, 0x8B, 0x03, 0x33, 0x33, 0x88, 0x88, 0xDD, 0xDD},
		"<arr>13107</arr><arr>-30584</arr><arr>-8739</arr>")
}

func TestDecodeShortAttribute(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x04, 0x61, 0x74, 0x74, 0x72, 0x84, 0x01},
		"<doc attr=\"false\"></doc>")
}

func TestDecodeExampleAttribute(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x03, 0x70, 0x72, 0x65, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x05, 0x03, 0x70, 0x72, 0x65, 0x04, 0x61, 0x74, 0x74, 0x72, 0x84, 0x01},
		"<doc xmlns:pre=\"http://abc\" pre:attr=\"false\"></doc>")
}

func TestDecodeExampleShortDictionaryAttribute(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x08, 0x86, 0x01},
		"<doc str8=\"true\"></doc>")
}

func TestDecodeExampleDictionaryAttribute(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x03, 0x70, 0x72, 0x65, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x07, 0x03, 0x70, 0x72, 0x65, 0x00, 0x86, 0x01},
		"<doc xmlns:pre=\"http://abc\" pre:str0=\"true\"></doc>")
}

func TestDecodeExampleShortXmlnsAttribute(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x08, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<doc xmlns=\"http://abc\"></doc>")
}

func TestDecodeExampleXmlnsAttribute(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x01, 0x70, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<doc xmlns:p=\"http://abc\"></doc>")
}

func TestDecodeExampleShortDictionaryXmlnsAttribute(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x0A, 0x04, 0x01},
		"<doc xmlns=\"str4\"></doc>")
}

func TestDecodeExampleDictionaryXmlnsAttribute(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x0B, 0x01, 0x70, 0x04, 0x01},
		"<doc xmlns:p=\"str4\"></doc>")
}

func TestDecodeExamplePrefixDictionaryAttributeF(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x01, 0x66, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x11, 0x0B, 0x98, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x01},
		"<doc xmlns:f=\"http://abc\" f:str11=\"hello\"></doc>")
}

func TestDecodeExamplePrefixDictionaryAttributeX(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x01, 0x78, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x23, 0x15, 0x98, 0x05, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x01},
		"<doc xmlns:x=\"http://abc\" x:str21=\"world\"></doc>")
}

func TestDecodeExamplePrefixAttributeK(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x01, 0x6B, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x30, 0x04, 0x61, 0x74, 0x74, 0x72, 0x86, 0x01},
		"<doc xmlns:k=\"http://abc\" k:attr=\"true\"></doc>")
}

func TestDecodeExamplePrefixAttributeZ(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x01, 0x7A, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x3F, 0x03, 0x61, 0x62, 0x63, 0x98, 0x03, 0x78, 0x79, 0x7A, 0x01},
		"<doc xmlns:z=\"http://abc\" z:abc=\"xyz\"></doc>")
}

func TestDecodeExampleShortElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x01},
		"<doc></doc>")
}

func TestDecodeExampleElement(t *testing.T) {
	testDecode(t,
		[]byte{0x41, 0x03, 0x70, 0x72, 0x65, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x03, 0x70, 0x72, 0x65, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<pre:doc xmlns:pre=\"http://abc\"></pre:doc>")
}

func TestDecodeExampleShortDictionaryElement(t *testing.T) {
	testDecode(t,
		[]byte{0x42, 0x0E, 0x01},
		"<str14></str14>")
}

func TestDecodeExampleDictionaryElement(t *testing.T) {
	testDecode(t,
		[]byte{0x43, 0x03, 0x70, 0x72, 0x65, 0x0E, 0x09, 0x03, 0x70, 0x72, 0x65, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<pre:str14 xmlns:pre=\"http://abc\"></pre:str14>")
}

func TestDecodeExamplePrefixDictionaryElementA(t *testing.T) {
	testDecode(t,
		[]byte{0x44, 0x0A, 0x09, 0x01, 0x61, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<a:str10 xmlns:a=\"http://abc\"></a:str10>")
}

func TestDecodeExamplePrefixDictionaryElementS(t *testing.T) {
	testDecode(t,
		[]byte{0x56, 0x26, 0x09, 0x01, 0x73, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<s:str38 xmlns:s=\"http://abc\"></s:str38>")
}

func TestDecodeExamplePrefixElementA(t *testing.T) {
	testDecode(t,
		[]byte{0x5E, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x09, 0x01, 0x61, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<a:hello xmlns:a=\"http://abc\"></a:hello>")
}

func TestDecodeExamplePrefixElementS(t *testing.T) {
	testDecode(t,
		[]byte{0x70, 0x09, 0x4D, 0x79, 0x4D, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x09, 0x01, 0x73, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<s:MyMessage xmlns:s=\"http://abc\"></s:MyMessage>")
}

func TestDecodeExampleZeroText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xA0, 0x03, 0x80, 0x01},
		"<doc str416=\"0\"></doc>")
}

func TestDecodeExampleZeroTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x61, 0x62, 0x63, 0x81},
		"<abc>0</abc>")
}

func TestDecodeExampleOneText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0x82, 0x01},
		"<doc str0=\"1\"></doc>")
}

func TestDecodeExampleOneTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x61, 0x62, 0x63, 0x83},
		"<abc>1</abc>")
}

func TestDecodeExampleFalseText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0x84, 0x01},
		"<doc str0=\"false\"></doc>")
}

func TestDecodeExampleFalseTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x61, 0x62, 0x63, 0x85},
		"<abc>false</abc>")
}

func TestDecodeExampleTrueText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0x86, 0x01},
		"<doc str0=\"true\"></doc>")
}

func TestDecodeExampleTrueTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x61, 0x62, 0x63, 0x87},
		"<abc>true</abc>")
}

func TestDecodeExampleInt8Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xEC, 0x01, 0x88, 0xDE, 0x01},
		"<doc str236=\"-34\"></doc>")
}

func TestDecodeExampleInt8TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x42, 0x9A, 0x01, 0x89, 0x7F},
		"<str154>127</str154>")
}

func TestDecodeExampleInt16Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xEC, 0x01, 0x8A, 0x00, 0x80, 0x01},
		"<doc str236=\"-32768\"></doc>")
}

func TestDecodeExampleInt16TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x42, 0x9A, 0x01, 0x8B, 0xFF, 0x7F},
		"<str154>32767</str154>")
}

func TestDecodeExampleInt32Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xEC, 0x01, 0x8C, 0x15, 0xCD, 0x5B, 0x07, 0x01},
		"<doc str236=\"123456789\"></doc>")
}

func TestDecodeExampleInt32TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x42, 0x9A, 0x01, 0x8D, 0xFF, 0xFF, 0xFF, 0x7F},
		"<str154>2147483647</str154>")
}

func TestDecodeExampleInt64Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xEC, 0x01, 0x8E, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x01},
		"<doc str236=\"2147483648\"></doc>")
}

func TestDecodeExampleInt64TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x42, 0x9A, 0x01, 0x8F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00},
		"<str154>1099511627776</str154>")
}

func TestDecodeExampleFloatText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x61, 0x90, 0xCD, 0xCC, 0x8C, 0x3F, 0x01},
		"<doc a=\"1.1\"></doc>")
}

func TestDecodeExampleFloatTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x91, 0xCD, 0xCC, 0x01, 0x42},
		"<Price>32.45</Price>")
}

func TestDecodeExampleDoubleText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x61, 0x92, 0x74, 0x57, 0x14, 0x8B, 0x0A, 0xBF, 0x05, 0x40, 0x01},
		"<doc a=\"2.71828182845905\"></doc>")
}

func TestDecodeExampleDoubleTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x02, 0x50, 0x49, 0x93, 0x11, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0x40},
		"<PI>3.14159265358979</PI>")
}

func TestDecodeExampleDecimalText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x03, 0x69, 0x6E, 0x74, 0x94, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x2D, 0x4E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		"<doc int=\"5.123456\"></doc>")
}

func TestDecodeExampleDecimalTextNegative(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x03, 0x69, 0x6E, 0x74, 0x94, 0x00, 0x00, 0x06, 0x80, 0x00, 0x00, 0x00, 0x00, 0x80, 0x2D, 0x4E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		"<doc int=\"-5.123456\"></doc>")
}

func TestDecodeExampleDecimalTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x08, 0x4D, 0x61, 0x78, 0x56, 0x61, 0x6C, 0x75, 0x65, 0x95, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		"<MaxValue>79228162514264337593543950335</MaxValue>")
}

func TestDecodeExampleDateTimeText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x6E, 0x96, 0xFF, 0x3F, 0x37, 0xF4, 0x75, 0x28, 0xCA, 0x2B, 0x01},
		"<doc str110=\"9999-12-31T23:59:59.9999999\"></doc>")
}

func TestDecodeExampleDateTimeTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x42, 0x6C, 0x97, 0x00, 0x40, 0x8E, 0xF9, 0x5B, 0x47, 0xC8, 0x08},
		"<str108>2006-05-17T00:00:00</str108>")
}

func TestDecodeExampleChars8Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x98, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x01},
		"<doc>hello</doc>")
}

func TestDecodeExampleChars8TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x01, 0x61, 0x99, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F},
		"<a>hello</a>")
}

func TestDecodeExampleChars16Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x9A, 0x05, 0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x01},
		"<doc>hello</doc>")
}

func TestDecodeExampleChars16TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x01, 0x61, 0x9B, 0x05, 0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F},
		"<a>hello</a>")
}

func TestDecodeExampleChars32Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x9C, 0x05, 0x00, 0x00, 0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x01},
		"<doc>hello</doc>")
}

func TestDecodeExampleChars32TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x01, 0x61, 0x9D, 0x05, 0x00, 0x00, 0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F},
		"<a>hello</a>")
}

func TestDecodeExampleBytes8Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x9E, 0x08, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x01},
		"<doc>AAECAwQFBgc=</doc>")
}

func TestDecodeExampleBytes8TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x06, 0x42, 0x61, 0x73, 0x65, 0x36, 0x34, 0x9F, 0x08, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
		"<Base64>AAECAwQFBgc=</Base64>")
}

func TestDecodeExampleBytes16Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xA0, 0x08, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x01},
		"<doc>AAECAwQFBgc=</doc>")
}

func TestDecodeExampleBytes16TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x06, 0x42, 0x61, 0x73, 0x65, 0x36, 0x34, 0xA1, 0x08, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
		"<Base64>AAECAwQFBgc=</Base64>")
}

func TestDecodeExampleBytes32Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xA2, 0x08, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x01},
		"<doc>AAECAwQFBgc=</doc>")
}

func TestDecodeExampleBytes32TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x06, 0x42, 0x61, 0x73, 0x65, 0x36, 0x34, 0xA3, 0x08, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
		"<Base64>AAECAwQFBgc=</Base64>")
}

func TestDecodeExampleStartListText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x61, 0xA4, 0x88, 0x7B, 0x98, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x86, 0xA6, 0x01},
		"<doc a=\"123 hello true\"></doc>")
}

func TestDecodeExampleEndListText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x61, 0xA4, 0x88, 0x7B, 0x98, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x86, 0xA6, 0x01},
		"<doc a=\"123 hello true\"></doc>")
}

func TestDecodeExampleEmptyText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x61, 0xA8, 0x01},
		"<doc a=\"\"></doc>")
}

func TestDecodeExampleEmptyTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xA9},
		"<doc></doc>")
}

func TestDecodeExampleDictionaryText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x02, 0x6E, 0x73, 0xAA, 0x38, 0x01},
		"<doc ns=\"str56\"></doc>")
}

func TestDecodeExampleDictionaryTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x04, 0x54, 0x79, 0x70, 0x65, 0xAB, 0xC4, 0x01},
		"<Type>str196</Type>")
}

func TestDecodeExampleUniqueIdText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xAC, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x01},
		"<doc>urn:uuid:33221100-5544-7766-8899-aabbccddeeff</doc>")
}

func TestDecodeExampleUniqueIdTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x42, 0x1A, 0xAD, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
		"<str26>urn:uuid:33221100-5544-7766-8899-aabbccddeeff</str26>")
}

func TestDecodeExampleTimeSpanText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xAE, 0x00, 0xC4, 0xF5, 0x32, 0xFF, 0xFF, 0xFF, 0xFF, 0x01},
		"<doc>-PT5M44S</doc>")
}

func TestDecodeExampleTimeSpanTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x42, 0x94, 0x07, 0xAF, 0x00, 0xB0, 0x8E, 0xF0, 0x1B, 0x00, 0x00, 0x00},
		"<str916>PT3H20M</str916>")
}

func TestDecodeExampleUuidText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xB0, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x01},
		"<doc>03020100-0504-0706-0809-0a0b0c0d0e0f</doc>")
}

func TestDecodeExampleUuidTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x02, 0x49, 0x44, 0xB1, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
		"<ID>03020100-0504-0706-0809-0a0b0c0d0e0f</ID>")
}

func TestDecodeExampleUInt64Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xB2, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01},
		"<doc>18446744073709551615</doc>")
}

func TestDecodeExampleUInt64TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x42, 0x9A, 0x01, 0xB3, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		"<str154>18446744073709551614</str154>")
}

func TestDecodeExampleBoolText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xB4, 0x01, 0x01},
		"<doc>true</doc>")
}

func TestDecodeExampleBoolTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x03, 0x40, 0x03, 0x61, 0x72, 0x72, 0x01, 0xB5, 0x05, 0x01, 0x00, 0x01, 0x00, 0x01},
		"<arr>true</arr><arr>false</arr><arr>true</arr><arr>false</arr><arr>true</arr>")
}

func TestDecodeExampleUnicodeChars8Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x75, 0xB6, 0x06, 0x75, 0x00, 0x6E, 0x00, 0x69, 0x00, 0x01},
		"<doc u=\"uni\"></doc>")
}

func TestDecodeExampleUnicodeChars8TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x01, 0x55, 0xB7, 0x06, 0x75, 0x00, 0x6E, 0x00, 0x69, 0x00},
		"<U>uni</U>")
}

func TestDecodeExampleUnicodeChars16Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x03, 0x75, 0x31, 0x36, 0xB8, 0x08, 0x00, 0x75, 0x00, 0x6E, 0x00, 0x69, 0x00, 0x32, 0x00, 0x01},
		"<doc u16=\"uni2\"></doc>")
}

func TestDecodeExampleUnicodeChars16TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x55, 0x31, 0x36, 0xB9, 0x08, 0x00, 0x75, 0x00, 0x6E, 0x00, 0x69, 0x00, 0x32, 0x00},
		"<U16>uni2</U16>")
}

func TestDecodeExampleUnicodeChars32Text(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x03, 0x75, 0x33, 0x32, 0xBA, 0x04, 0x00, 0x00, 0x00, 0x33, 0x00, 0x32, 0x00, 0x01},
		"<doc u32=\"32\"></doc>")
}

func TestDecodeExampleUnicodeChars32TextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x55, 0x33, 0x32, 0xBB, 0x04, 0x00, 0x00, 0x00, 0x33, 0x00, 0x32, 0x00},
		"<U32>32</U32>")
}

func TestDecodeExampleQNameDictionaryText(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xF0, 0x06, 0xBC, 0x08, 0x8E, 0x07, 0x01},
		"<doc str880=\"i:str910\"></doc>")
}

func TestDecodeExampleQNameDictionaryTextWithEndElement(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x04, 0x54, 0x79, 0x70, 0x65, 0xBD, 0x12, 0x90, 0x07},
		"<Type>s:str912</Type>")
}

func TestDecodeExampleUnicodeChars16TextWithChinese(t *testing.T) {
	testDecode(t,
		[]byte{0x40, 0x0c, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0xb7, 0x08, 0x91, 0x4E, 0x62, 0x88, 0x2D, 0x4E, 0x66, 0x5B, 0x5F},
		"<PositionName>云衢中学</PositionName>")
}

//----------------------------------------------------

func TestReadMultiByteInt31_17(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x11}, 17)
}

func TestReadMultiByteInt31_145(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x91, 0x01}, 145)
}

func TestReadMultiByteInt31_5521(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x91, 0x2B}, 5521)
}

func TestReadMultiByteInt31_16384(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x80, 0x80, 0x01}, 16384)
}

func TestReadMultiByteInt31_2097152(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x80, 0x80, 0x80, 0x01}, 2097152)
}

func TestReadMultiByteInt31_268435456(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x80, 0x80, 0x80, 0x80, 0x01}, 268435456)
}

func TestCharsEscapingInAttributes(t *testing.T) {
	t.Skip("TODO")
	t.Error("Not Implemented: Chars Escaping in Attributes")
}

func TestCharsEscapingInElements(t *testing.T) {
	t.Skip("TODO")
	t.Error("Not Implemented: Chars Escaping in Elements")
}

func TestReadFloatTextSpecialValues(t *testing.T) {
	testReadFloatTextSpecialValue(t, float32(math.Inf(1)), "INF")
	testReadFloatTextSpecialValue(t, float32(math.Inf(-1)), "-INF")
	testReadFloatTextSpecialValue(t, float32(math.NaN()), "NaN")
	testReadFloatTextSpecialValue(t, float32(math.Copysign(0, -1)), "-0")
}

func TestReadDoubleTextSpecialValues(t *testing.T) {
	testReadDoubleTextSpecialValue(t, math.Inf(1), "INF")
	testReadDoubleTextSpecialValue(t, math.Inf(-1), "-INF")
	testReadDoubleTextSpecialValue(t, math.NaN(), "NaN")
	testReadDoubleTextSpecialValue(t, math.Copysign(0, -1), "-0")
}

func testReadFloatTextSpecialValue(t *testing.T, num float32, expected string) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, &num)
	r := bytes.NewReader(buf.Bytes())
	d := &decoder{bin: r}
	actual, err := readFloatText(d)
	if err != nil {
		t.Error(err.Error())
	}
	assertStringEqual(t, actual, expected)
}

func testReadDoubleTextSpecialValue(t *testing.T, num float64, expected string) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, &num)
	r := bytes.NewReader(buf.Bytes())
	d := &decoder{bin: r}
	actual, err := readDoubleText(d)
	if err != nil {
		t.Error(err.Error())
	}
	assertStringEqual(t, actual, expected)
}

func testDecode(t *testing.T, bin []byte, expected string) {
	decoder := NewDecoder()
	actual, err := decoder.Decode(bytes.NewReader(bin))
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + actual)
	}
	assertStringEqual(t, actual, expected)
}

func testReadMultiByteInt31(t *testing.T, bin []byte, expected uint32) {
	reader := bytes.NewReader(bin)
	actual, err := readMultiByteInt31(reader)
	if err != nil {
		t.Error("Error: " + err.Error())
		return
	}
	if actual != expected {
		t.Errorf("Expected %d but got %d", expected, actual)
		return
	}
}
