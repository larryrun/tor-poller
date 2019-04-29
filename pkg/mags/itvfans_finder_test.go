package mags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestITVFansFinder_ListAvailableMagLinks(t *testing.T) {
	extra := make(map[string]interface{})
	extra["start-season"] = "9"
	extra["start-episode"] = "4"
	finder := ITVFansFinder{Page: "http://www.itvfans.com/ziyuan/1076002.html", Extra: extra}
	items, err := finder.ListAvailableItems()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, len(items) > 0)
}
