package http_proxy

import (
	"os"
	"testing"
)

func TestCreateCA(t *testing.T) {
	t.Logf("Begin")
	os.Chdir("../../../")
	dir, err := os.Getwd()
	if err != nil {
		t.Logf("Get current directory error %v", err)
		t.Fail()
	}
	t.Logf("Current directory %v", dir)
	err = createCA("./cert/cert.pem", "./cert/key.pem")
	if err != nil {
		t.Logf("Generate error %v", err)
		t.Fail()
	}
}
