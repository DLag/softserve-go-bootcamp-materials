package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type idGenServerAppMysqlTest struct {
	cfg                 cfgModel
	idgen               idGenModel
	heartbeat, generate http.Handler
}

func (app *idGenServerAppMysqlTest) Run() error {
	app.idgen.Alive()
	return nil
}

func newIdGenServerAppMysqlTest() *idGenServerAppMysqlTest {
	app := new(idGenServerAppMysqlTest)
	app.cfg = &cfgModelMock{data: map[string]string{
		"DSN": "Predefined",
	}}
	app.idgen = NewIdGenModelAtomicTest()
	app.heartbeat = newHeartbeatHandler(app.idgen)
	app.generate = newIdGenHandler(app.idgen)
	return app
}

func (app *idGenServerAppMysqlTest) TestWebServer(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(newIdGenHandler(app.idgen))

	defer server.Close()

	for i := 1; i <= 10; i++ {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
		}
		expected := fmt.Sprint(i)
		actual, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		if expected != string(actual) {
			t.Errorf("Expected the message %q, recieved %q\n", expected, string(actual))
		}
	}
	gen, cur, alive := app.idgen.(modelStub).GetCounters()
	assert.Equal([]uint32{gen, cur, alive}, []uint32{10, 0, 0}, "Wrong counters")

}

func TestWebServer(t *testing.T) {
	app := newIdGenServerAppMysqlTest()
	app.TestWebServer(t)
}
