/**
 * Copyright 2018 Curtis Mattoon
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package secret

import (
	//"reflect"
	"testing"

	"github.com/cmattoon/aws-ssm/pkg/provider"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestParseStringList(t *testing.T) {
	for _, tc := range []struct {
		title    string
		pvalue   string
		expected map[string]string
	}{
		{
			title:    "parse empty value",
			pvalue:   "",
			expected: map[string]string{},
		},
		{
			title:    "parse whitespace-only value",
			pvalue:   "    ",
			expected: map[string]string{},
		},
		{
			title:  "parse some simple values",
			pvalue: "key1=val1,key2=val2,key3=val3,key4=val4",
			expected: map[string]string{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
				"key4": "val4",
			},
		},
		{
			title:  "parse equal sign in value",
			pvalue: "key1=val1,key2=val2,key3=val3,key4=val4=true",
			expected: map[string]string{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
				"key4": "val4=true",
			},
		},
		{
			title:  "parse non-list",
			pvalue: "key1",
			expected: map[string]string{
				"key1": "",
			},
		},
		{
			title:  "parse trailing comma",
			pvalue: "key1=value1,",
			expected: map[string]string{
				"key1": "value1",
			},
		},
		{
			title:  "parse leading/trailing comma",
			pvalue: ",key1=value1,key2=val2,key3=val3,",
			expected: map[string]string{
				"key1": "value1",
				"key2": "val2",
				"key3": "val3",
			},
		},
		{
			title:  "not a StringList",
			pvalue: "asdfjkl",
			expected: map[string]string{
				"asdfjkl": "",
			},
		},
		{
			title:  "fine i guess",
			pvalue: "ThIsMiGhTBeBaSe64==",
			expected: map[string]string{
				"ThIsMiGhTBeBaSe64": "=",
			},
		},
		{
			title:  "another StringList",
			pvalue: "123,456,789",
			expected: map[string]string{
				"123": "",
				"456": "",
				"789": "",
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			s := &Secret{
				Name:       "test_secret",
				Namespace:  "test",
				ParamName:  "FOO_PARAM",
				ParamType:  "StringList",
				ParamKey:   "some-string-list",
				ParamValue: tc.pvalue,
				Data:       map[string]string{},
			}
			data := s.ParseStringList()
			assert.Equal(t, tc.expected, data)
		})
	}

}

// Should set the key/value pair
func TestSet(t *testing.T) {
	s := &Secret{
		Name:       "test_secret",
		Namespace:  "test",
		ParamName:  "FOO_PARAM",
		ParamType:  "StringList",
		ParamKey:   "foo-param",
		ParamValue: "key1=val1,key2=val2,key3=val3,key4=val4=true",
		Data:       map[string]string{},
	}
	s.Set("foo", "bar")
	if s.Secret.StringData["foo"] != "bar" {
		t.Fail()
	}
}

func TestSetRefusesToOverwriteKey(t *testing.T) {
	s := &Secret{
		Name:       "test_secret",
		Namespace:  "test",
		ParamName:  "FOO_PARAM",
		ParamType:  "StringList",
		ParamKey:   "foo-param",
		ParamValue: "key1=val1,key2=val2,key3=val3,key4=val4=true",
		Data:       map[string]string{},
	}
	err := s.Set("foo", "bar")
	if err != nil {
		t.Fail()
	}
	err = s.Set("foo", "baz")
	if err != nil {
		t.Fail()
	}
	if s.Secret.StringData["foo"] != "baz" {
		t.Fail()
	}
}

func TestSetsValue(t *testing.T) {
	p := provider.MockProvider{"FooBar123", "PlaintextIsAnError", make(map[string]string)}
	s := v1.Secret{}
	testSecret, err := NewSecret(s, p, "foo-secret", "namespace", "foo-param", "String", "")
	if err != nil {
		t.Fail()
	}
	if testSecret.ParamValue != "FooBar123" {
		t.Fail()
	}
}

// When the encryption key is defined, the decrypted value should be returned
func TestNewSecretDecryptsIfKeyIsSet(t *testing.T) {
	p := provider.MockProvider{"$@#*$(@)*$", "FooBar123", make(map[string]string)}
	s := v1.Secret{}
	testSecret, err := NewSecret(s, p, "foo-secret", "namespace", "foo-param", "String", "my/test/key")
	if err != nil {
		t.Fail()
	}
	if testSecret.ParamValue != p.DecryptedValue {
		t.Fail()
	}
}

func TestFromKubernetesSecretReturnsErrorIfIrrelevant(t *testing.T) {
	p := provider.MockProvider{"$@#*$(@)*$", "FooBar123", make(map[string]string)}
	s := v1.Secret{} // No annotations, so no params

	_, err := FromKubernetesSecret(p, s)
	if err.Error() != "Irrelevant Secret" {
		t.Fail()
	}
}

// If the parameter is of Type=SecureString, and no key is supplied,
// attempt to use the default key.
func TestFromKubernetesSecretUsesDefaultEncryptionKey(t *testing.T) {
	p := provider.MockProvider{"$@#*$(@)*$", "FooBar123", make(map[string]string)}

	s := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"alpha.ssm.cmattoon.com/aws-param-name": "foo-param",
				"alpha.ssm.cmattoon.com/aws-param-type": "SecureString",
			},
		},
	}

	ks, err := FromKubernetesSecret(p, s)

	if err != nil || ks.ParamKey != "alias/aws/ssm" || ks.ParamValue != "FooBar123" {
		t.Fail()
	}
}

func TestFromKubernetesSecretUsesSpecifiedEncryptionKey(t *testing.T) {
	p := provider.MockProvider{"$@#*$(@)*$", "FooBar123", make(map[string]string)}

	s := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"alpha.ssm.cmattoon.com/aws-param-name": "foo-param",
				"alpha.ssm.cmattoon.com/aws-param-type": "SecureString",
				"alpha.ssm.cmattoon.com/aws-param-key":  "foo/bar/baz",
			},
		},
	}

	ks, err := FromKubernetesSecret(p, s)

	if err != nil || ks.ParamKey != "foo/bar/baz" || ks.ParamValue != "FooBar123" {
		t.Fail()
	}
}

func TestSafeKeyName(t *testing.T) {
	keys := map[string]string{
		"/foo/bar":     "foo_bar",
		"/foo/bar/":    "foo_bar",
		"//foo/bar":    "_foo_bar",
		"//foo/bar/":   "_foo_bar",
		"/foo/bar/baz": "foo_bar_baz",
	}
	for path, exp := range keys {
		assert.Equal(t, safeKeyName(path), exp)
	}
}
