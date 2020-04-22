package youtube

import (
	"fmt"
	"github.com/lithdew/bytesutil"
	"github.com/lithdew/youtube/sig"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

type Cipher struct {
	URL             string `json:"url"`
	Signature       string `json:"s"`
	SignaturePolicy string `json:"sp"`
}

func ParseCipherJSON(v *fastjson.Value) Cipher {
	a := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(a)

	a.ParseBytes(v.GetStringBytes())

	return Cipher{
		URL:             string(a.Peek("url")),
		Signature:       string(a.Peek("s")),
		SignaturePolicy: string(a.Peek("sp")),
	}
}

func (c Cipher) DecodeURL(script string) (string, error) {
	factory, err := sig.LookupCipherFactory(script)
	if err != nil {
		return "", fmt.Errorf("failed to lookup cipher factory in script: %w", err)
	}

	cipher, err := sig.LookupCipher(factory, script)
	if err != nil {
		return "", fmt.Errorf("failed to lookup cipher steps in script: %w", err)
	}

	decoded := cipher.Decode(c.Signature)

	uri := fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(uri)

	uri.Parse(nil, bytesutil.Slice(c.URL))

	if c.SignaturePolicy == "" {
		uri.QueryArgs().Add("signature", decoded)
	} else {
		uri.QueryArgs().Add(c.SignaturePolicy, decoded)
	}

	return uri.String(), nil
}
