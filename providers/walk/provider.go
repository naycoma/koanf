package walk

import (
	"errors"

	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// 指定されたディレクトリとファイル名を用いてwalkする Provider
//  1. name1.yml
//  2. child/name1.yml
//  3. name1@[suf].yml
//  4. child/name1@[suf].yml
//  5. name2.yml
func Provider(parser koanf.Parser, root string, fileNames ...string) *Walk {
	return &Walk{
		Parser: parser,
		Root:   root,
		Names:  fileNames,
	}
}

var _ koanf.Provider = (*Walk)(nil)

type Walk struct {
	Parser koanf.Parser
	Root   string
	Names  []string
}

func (p *Walk) ReadBytes() ([]byte, error) {
	return nil, errors.New("Walk does not support ReadBytes()")
}

func (p *Walk) Read() (map[string]any, error) {
	sortedPriorityWalks, err := sortedPriorityWalk(p.Root, p.Names...)
	if err != nil {
		return nil, err
	}
	if len(sortedPriorityWalks) == 0 {
		return nil, errors.New("no matching files found")
	}
	k := koanf.New(".")
	for _, item := range sortedPriorityWalks {
		if err := k.Load(file.Provider(item.Path), p.Parser); err != nil {
			return nil, err
		}
	}
	return k.Raw(), nil
}
