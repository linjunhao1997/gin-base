package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
)

type Attachment struct {
	Filename    string
	ReadSetting []gomail.FileSetting
}

func FilenameSetting(filename string) gomail.FileSetting {
	return gomail.SetHeader(map[string][]string{
		"Content-Disposition": {
			fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", filename)),
		}})
}

func CopyFileSetting(bytes []byte) gomail.FileSetting {
	return gomail.SetCopyFunc(func(w io.Writer) error {
		if _, err := w.Write(bytes); err != nil {
			return err
		}
		return nil
	})
}

func CopyFileByHttp(url string) gomail.FileSetting {
	return gomail.SetCopyFunc(func(w io.Writer) error {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if _, err := w.Write(body); err != nil {
			return err
		}

		if err := resp.Body.Close(); err != nil {
			return err
		}

		return nil
	})
}
