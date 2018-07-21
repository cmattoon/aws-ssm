package secret

import (
	"reflect"
	"testing"
)


func TestParseStringList(t *testing.T) {
	s := &Secret{
		Name: "test_secret",
		Namespace: "test",
		ParamName: "FOO_PARAM",
		ParamType: "StringList",
		ParamKey: "foo-param",
		ParamValue:"key1=val1,key2=val2,key3=val3,key4=val4=true",
		Data: map[string]string{},
	}
	
	expected := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4=true",
	}
	
	data, err := s.ParseStringList()
	if err != nil {
		t.Fail()
	}
	
	eq := reflect.DeepEqual(data, expected)
	if ! eq {
		t.Fail()
	}
}

// Should set the key/value pair
func TestSet(t *testing.T) {
	s := &Secret{
		Name: "test_secret",
		Namespace: "test",
		ParamName: "FOO_PARAM",
		ParamType: "StringList",
		ParamKey: "foo-param",
		ParamValue:"key1=val1,key2=val2,key3=val3,key4=val4=true",
		Data: map[string]string{},
	}
	s.Set("foo", "bar")
	if s.Secret.StringData["foo"] != "bar" {
		t.Fail()
	}
}

