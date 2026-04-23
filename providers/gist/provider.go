package gist

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/knadh/koanf/v2"
)

// GitHub Gistのファイルを取得する Provider
func Provider(user string, id string, file string) *Gist {
	return &Gist{
		User: user,
		Id:   id,
		File: file,
	}
}

var GIST_CONTENT_URL *url.URL = func() *url.URL {
	u, err := url.Parse("https://gist.githubusercontent.com/")
	if err != nil {
		panic(err)
	}
	return u
}()

var _ koanf.Provider = (*Gist)(nil)

type Gist struct {
	User string
	Id   string
	File string
}

func (p *Gist) Url() *url.URL {
	return GIST_CONTENT_URL.JoinPath(p.User, p.Id, "raw", p.File)
}

func (p *Gist) ReadBytes() ([]byte, error) {
	resp, err := http.Get(p.Url().String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (p *Gist) Read() (map[string]any, error) {
	return nil, errors.New("Gist does not support Read()")
}
