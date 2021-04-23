package log

import (
	"runtime"
	"testing"
	"time"

	"github.com/francoispqt/onelog"
	"github.com/stretchr/testify/assert"
)

var log *Log

type TestWriter struct {
	b      []byte
	called bool
}

func (t *TestWriter) Write(b []byte) (int, error) {
	t.called = true
	if len(t.b) < len(b) {
		t.b = make([]byte, len(b))
	}
	copy(t.b, b)
	return len(t.b), nil
}

func newWriter() *TestWriter {
	return &TestWriter{make([]byte, 0, 512), false}
}

func Test_XXXMsg(t *testing.T) {
	writer := newWriter()
	l := GetLogConfig()
	l.CallersLevels = 7
	l.Level = "debug"
	l.Writer = writer
	log = NewLog(l)
	_, currentFilePath, _, _ := runtime.Caller(0)

	log.ErrorMsg("Hello").String("Name", "John").WithCaller().Submit()
	json := `{"level":"error","message":"Hello","time":"` + time.Now().Format(time.RFC3339) + `","Name":"John","CALLER":"` + currentFilePath + `:41 /usr/local/go/src/testing/testing.go:777 /usr/local/go/src/runtime/asm_amd64.s:2361 "}` + "\n"

	assert.Equal(t, json, string(writer.b))

	log.DebugMsg("Hello").String("Name", "John").WithCaller().Submit()
	json = `{"level":"debug","message":"Hello","time":"` + time.Now().Format(time.RFC3339) + `","Name":"John","CALLER":"` + currentFilePath + `:46 /usr/local/go/src/testing/testing.go:777 /usr/local/go/src/runtime/asm_amd64.s:2361 "}` + "\n"

	assert.Equal(t, json, string(writer.b))

	log.InfoMsg("Hello").String("Name", "John").WithCaller().Submit()
	json = `{"level":"info","message":"Hello","time":"` + time.Now().Format(time.RFC3339) + `","Name":"John","CALLER":"` + currentFilePath + `:51 /usr/local/go/src/testing/testing.go:777 /usr/local/go/src/runtime/asm_amd64.s:2361 "}` + "\n\n"

	assert.Equal(t, json, string(writer.b))
}

func Test_Mock(t *testing.T) {
	writer := newWriter()
	l := GetLogConfig()
	l.CallersLevels = 7
	l.Level = "debug"
	l.Writer = writer
	log = NewLog(l)
	log.SwitchToMockMode()

	log.Error("Print Error")
	json := ""
	assert.Equal(t, json, string(writer.b))

	log.Info("Print Error")
	assert.Equal(t, json, string(writer.b))

	log.Debug("Print Error")
	assert.Equal(t, json, string(writer.b))
}

func Benchmark_ErrorMsg(b *testing.B) {
	writer := newWriter()
	l := GetLogConfig()
	l.CallersLevels = 7
	l.Level = "debug"
	l.Writer = writer
	log = NewLog(l)
	for i := 0; i < b.N; i++ {
		log.ErrorMsg("Hello").String("Name", "John").WithCaller().Submit()
	}
}

func Benchmark_ErrorWithFields(b *testing.B) {
	writer := newWriter()
	l := GetLogConfig()
	l.CallersLevels = 7
	l.Level = "debug"
	l.Writer = writer
	log = NewLog(l)
	for i := 0; i < b.N; i++ {
		log.ErrorWithFields("Hello", func(e onelog.Entry) {
			e.String("Name", "John")
			e.String("CALLER", generateCallerList(7))
		})
	}
}
