package article

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.x2ox.com/sorbifolia/httputils"
	"mime/multipart"
	"path"
	"strings"
	"time"
)

func (art *Article) replaceImage() {
	for i, v := range art.Images {
		newURL, err := downloadAndUpload(v)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if i == 0 {
			art.Avatar = newURL
		}
		art.Content = strings.ReplaceAll(art.Content, v, newURL)
	}
}

func downloadAndUpload(s string) (string, error) {
	var (
		body  []byte
		resp  string
		retry = 5
	)

RetryDownload:
	if err := httputils.Get(s).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36").
		Request(3, nil, 30*time.Second).
		ParserData(func(resp *fasthttp.Response) error {
			data := resp.Body()
			if len(data) == 0 {
				return errors.New("null body")
			}

			body = make([]byte, len(data))
			copy(body, data)
			return nil
		}).DoRelease(); err != nil {
		if retry > 0 {
			retry--
			goto RetryDownload
		}

		return "", err
	}

	time.Now().Unix()

RetryUpload:
	if err := httputils.Post("https://s3.hera.show/api/v1/xos/blog/upload").
		SetContentType("multipart/form-data; charset=utf-8; boundary=__X_PAW_BOUNDARY__").
		SetBodyWithEncoder(httputils.FormData(func(w *multipart.Writer) error {
			if err := w.SetBoundary("__X_PAW_BOUNDARY__"); err != nil {
				return err
			}

			if err := w.WriteField("auth", "em7"); err != nil {
				return err
			}

			fw, err := w.CreateFormFile("file", path.Base(s))
			if err != nil {
				return err
			}
			_, err = fw.Write(body)

			return err
		}), nil).
		Request(5, nil, 30*time.Second).
		ParserData(func(res *fasthttp.Response) error {
			if res.StatusCode() != 200 {
				return errors.New("status code not 200")
			}
			if len(res.Body()) == 0 {
				return errors.New("null body")
			}
			resp = string(res.Body())
			return nil
		}).
		DoRelease(); err != nil {
		if retry > 0 {
			retry--
			goto RetryUpload
		}
		return "", fmt.Errorf("upload file err %s", err)
	}

	return resp, nil
}
