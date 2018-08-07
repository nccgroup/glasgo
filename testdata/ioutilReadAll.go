package main

import(
	"io/ioutil"
	"strings"
)

func testReadAll() string {
	r := strings.NewReader("this is a test for use of ioutil.ReadAll");

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return ""
	}
	return string(b)
}
