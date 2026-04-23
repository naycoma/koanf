package fetch

import (
	"errors"
	"io"
	"net/http"

	"github.com/knadh/koanf/v2"
)

// 指定されたURLをFetchする Provider
func Provider(url string) *Fetch {
	return &Fetch{Url: url}
}

var _ koanf.Provider = (*Fetch)(nil)

type Fetch struct {
	Url string
}

func (p *Fetch) ReadBytes() ([]byte, error) {
	resp, err := http.Get(p.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (p *Fetch) Read() (map[string]any, error) {
	return nil, errors.New("Fetch does not support Read()")
}
