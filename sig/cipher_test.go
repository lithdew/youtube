package sig

import (
	"github.com/lithdew/bytesutil"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestCipher(t *testing.T) {
	buf, err := ioutil.ReadFile("testdata/sec.js")
	require.NoError(t, err)

	script := bytesutil.String(buf)

	factory, err := LookupCipherFactory(script)
	require.NoError(t, err)

	steps, err := LookupCipher(factory, script)
	require.NoError(t, err)

	require.Len(t, steps, 3)
}
