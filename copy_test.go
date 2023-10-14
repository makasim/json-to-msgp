package jsontomsgp_test

import (
	"bytes"
	"testing"

	"github.com/makasim/jsontomsgp"
	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"
)

func TestCopyBytes(t *testing.T) {
	// JSON data as bytes
	src := []byte(`{"level": "debug", "event_timestamp": "2022-09-15T10:20:30Z", "msg": "the message", "caller": "the caller", "custom": "customVal", "array": ["foo", "bar"], "true": true, "false": false, "null": null}`)
	buf := bytes.NewBuffer(make([]byte, 0, len(src)))
	w := msgp.NewWriter(buf)

	err := jsontomsgp.CopyBytes(src, w)
	require.NoError(t, err)

	msgpB := buf.Bytes()
	jsonB := bytes.NewBuffer(make([]byte, 0, len(src)))

	_, err = msgp.UnmarshalAsJSON(jsonB, msgpB)
	require.NoError(t, err)

	require.JSONEq(t, string(src), jsonB.String())
}

func BenchmarkCopyJSONToMsgp(b *testing.B) {
	src := []byte(`{"level": "debug", "event_timestamp": "2022-09-15T10:20:30Z", "msg": "the message", "caller": "the caller", "custom": "customVal", "array": ["foo", "bar"], "true": true, "false": false, "null": null}`)
	dst := make([]byte, 0, len(src))

	var err error

	buf := bytes.NewBuffer(dst)
	w := msgp.NewWriter(buf)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err = jsontomsgp.CopyBytes(src, w); err != nil {
			b.Fatal(err)
		}
		buf.Reset()
		w.Reset(buf)
	}
}
