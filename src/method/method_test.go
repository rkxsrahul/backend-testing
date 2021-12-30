package method

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func TestHeaderFileSlug(t *testing.T) {
	if headerFileSlug("Hello") == "" {
		t.Error("test case fail")
	}

}

func TestProjectNamebyEmail(t *testing.T) {
	if ProjectNamebyEmail("") == "" {
		t.Error("test case fail")
	}
	if ProjectNamebyEmail("test@xenonstack.com") == "" {
		t.Error("test case fail")
	}
}

func TestCheckURL(t *testing.T) {
	_, err := CheckURL("www.xenonstack.com")
	if err != nil {
		t.Error("test case fail")
	}
	_, err = CheckURL("http://www.xenonstack.com")
	if err != nil {
		t.Error("test case fail")
	}
	_, err = CheckURL("www.eirbgfibeifbifbi")
	if err == nil {
		t.Error("test case fail")
	}
	_, err = CheckURL("eirbgfibeifbifbi")
	if err == nil {
		t.Error("test case fail")
	}
	_, err = CheckURL("https://eirbgfibeifbifbi")
	if err == nil {
		t.Error("test case fail")
	}
}
