package method

import (
	"testing"
)

func TestFetchLanguage(t *testing.T) {

	_, _, err := FetchLanguage("rkxsrahul", "A")
	if err != nil {
		t.Error("test case fail")
	}

}
