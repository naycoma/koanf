package expand

import (
	"os"

	"github.com/knadh/koanf/v2"
)

type EnvExpander struct {
	koanf.Parser
}

// koanf.ParserをラップしてUnmarshal前に環境変数を展開する汎用パーサー
func Parser(p koanf.Parser) koanf.Parser {
	return &EnvExpander{Parser: p}
}

func (p *EnvExpander) Unmarshal(b []byte) (map[string]any, error) {
	// `${VAR}` や `$VAR` 形式の文字列を環境変数の値で置換
	expanded := os.ExpandEnv(string(b))

	return p.Parser.Unmarshal([]byte(expanded))
}

func (p *EnvExpander) Marshal(m map[string]any) ([]byte, error) {
	return p.Parser.Marshal(m)
}