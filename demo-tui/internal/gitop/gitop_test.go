package gitop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	matcher, err := NewMatcher("/home/lv/lvsoso/tools/demo-tui/test/demo-repo/.gitattributes")
	assert.Nil(t, err)
	r, err := matcher.MatchLfs("/home/lv/lvsoso/tools/demo-tui/test/demo-repo/1.bin")
	assert.Nil(t, err)
	assert.True(t, r)
}
