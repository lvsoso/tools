package lfscache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLfsCache(t *testing.T) {
	src := "/home/lv/lvsoso/tools/demo-tui/test/demo-src"
	demo := "/home/lv/lvsoso/tools/demo-tui/test/demo-repo"
	err := LfsCache(src, demo, 1)
	assert.Nil(t, err)
}
