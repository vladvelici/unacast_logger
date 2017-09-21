package logger

import (
	"bytes"
	"testing"

	"encoding/json"

	"io/ioutil"
	"strings"

	"github.com/getsentry/raven-go"
	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"cloud.google.com/go/errorreporting"
)

type mockTransport struct {
	Packet chan raven.Packet
}

func (t *mockTransport) Send(url, authHeader string, packet *raven.Packet) error {
	t.Packet <- *packet
	close(t.Packet)
	return nil
}

func TestErrorLogging(t *testing.T) {
	var buf bytes.Buffer
	l := log.NewLogger3(&buf, "test", log.NewJSONFormatter("test"))

	err := errors.New("This is an error")
	sl := unaLogger{
		Logger: l,
	}
	msg := "Something is wrong"
	sl.Error(
		msg,
		err,
		"one", "1", "two", 2,
	)

	var obj map[string]interface{}
	jsonErr := json.Unmarshal(buf.Bytes(), &obj)
	if jsonErr != nil {
		t.Fatalf("Hmm, couldn't unmarshal the log buffer %v. %v", buf.String(), jsonErr)
	}
	if obj["_m"] != msg {
		t.Errorf("%v didn't match %v\n", obj["_m"], msg)
	}
	if obj["error"] != err.Error() {
		t.Errorf("%v didn't match %v\n", buf.String(), err.Error())
	}
	if obj["one"] != "1" {
		t.Errorf("%v didn't match %v\n", obj["one"], "1")
	}
	if obj["two"] != 2.0 {
		t.Errorf("%v didn't match %v\n", obj["two"], 2.0)
	}
}

func TestErrorLoggingToFile(t *testing.T) {

	err := errors.New("This is an error")
	testLogFile := "/tmp/test.log"
	sl := NewLogger(Config{
		Name:     "test",
		FileName: testLogFile,
	})
	msg := "Something is wrong"
	sl.Error(
		msg,
		err,
		"one", "1", "two", 2,
	)

	fileContent, err := ioutil.ReadFile(testLogFile)
	if err != nil {
		t.Fatal("Didn't expect an error, got: ", err)
	}

	if !strings.Contains(string(fileContent), msg) {
		t.Errorf("Expected %v to contain %v, did not", testLogFile, msg)
	}

}