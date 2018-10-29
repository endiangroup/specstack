package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Prefix_AppendsPartsOntoPrefix(t *testing.T) {
	testCases := []struct {
		parts    []string
		expected string
	}{
		{parts: []string{"testPostfix"}, expected: "testPrefix.testPostfix"},
		{parts: []string{"testPostfix1", "testPostfix2"}, expected: "testPrefix.testPostfix1.testPostfix2"},
	}
	pre := prefix("testPrefix")

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, pre.Append(testCase.parts...))
	}
}

func Test_FetchPrefx_FetchesAPrefixFromAString(t *testing.T) {
	for _, key := range []string{"testPrefix", "testPrefix.testPostfix", "testPrefix.testPostfix1.testPostfix2"} {
		assert.Equal(t, prefix("testPrefix"), fetchPrefix(key))
	}
}

func Test_FetchPostfix_FetchesAPostfixFromAString(t *testing.T) {
	testCases := []struct {
		key      string
		expected string
	}{
		{key: "testPostfix", expected: "testPostfix"},
		{key: "testPrefix.testPostfix", expected: "testPostfix"},
		{key: "testPrefix.testPostfix1.testPostfix2", expected: "testPostfix1.testPostfix2"},
	}
	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, fetchPostfix(testCase.key))
	}
}
