package site

import (
	"bytes"
	"errors"
	"io"
	"time"

	"github.com/valyala/fasthttp"
	"go.x2ox.com/sorbifolia/httputils"
)

func ioParser() (io.Reader, httputils.Parser) {
	var buf bytes.Buffer
	return &buf, func(resp *fasthttp.Response) error {
		arr := resp.Body()
		_, err := buf.Write(arr)
		return err
	}
}

func getSite(url string) (io.Reader, error) {
	var (
		r, p   = ioParser()
		status int
		err    = httputils.Get(url).
			SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36").
			Request(3, func(err error, response *fasthttp.Response) bool {
				status = response.StatusCode()
				return err != nil || response.StatusCode() != 200
			}, 10*time.Second).
			ParserData(p).DoRelease()
	)

	if err != nil {
		return nil, err
	}
	if status != 200 || r.(*bytes.Buffer).Len() == 0 {
		return nil, errors.New("page is null")
	}

	return r, nil
}
