package gitop

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing/format/gitattributes"
)

type Matcher struct {
	attrFilePath string
	mas          []gitattributes.MatchAttribute
}

func NewMatcher(attrFile string) (*Matcher, error) {
	f, err := os.Open(attrFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	mas, err := gitattributes.ReadAttributes(f, nil, true)
	if err != nil {
		return nil, err
	}
	return &Matcher{
		attrFilePath: attrFile,
		mas:          mas,
	}, nil
}

func (m *Matcher) MatchLfs(checkPath string) (bool, error) {
	relativePath, err := filepath.Rel(filepath.Dir(m.attrFilePath), checkPath)
	if err != nil {
		return false, err
	}

	for _, ma := range m.mas {
		if ma.Pattern.Match([]string{relativePath}) {
			for _, attr := range ma.Attributes {
				if attr.Name() == "filter" && attr.Value() == "lfs" {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
