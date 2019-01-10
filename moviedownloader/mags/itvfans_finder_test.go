package mags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestITVFansFinder_ListAvailableMagLinks(t *testing.T) {
	finder := ITVFansFinder{Page: "http://www.itvfans.com/ziyuan/1076002.html"}
	magLinks, err := finder.ListAvailableItems()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, len(magLinks) > 0)
}

func TestItvfans_fill_size_info(t *testing.T) {
	targetItem := &TargetItem{}
	err := itvfans_fill_size_info(targetItem,"1003.74MB")
	assert.Nil(t, err)
	assert.Equal(t, 1003, targetItem.Size)

	err = itvfans_fill_size_info(targetItem,"2.61GB")
	assert.Nil(t, err)
	assert.Equal(t, 2672, targetItem.Size)
}

func TestItvfans_extract_magLink(t *testing.T) {
	targetItem := &TargetItem{}
	err := itvfans_fill_episode_resolution_info(targetItem,"NCIS.Los.Angeles.S10E01.1080p.AMZN.WEB-DL.DDP5.1.H.264-ViSUM.mkv")
	assert.Nil(t, err)
	assert.Equal(t, "1080p", targetItem.Resolution)
	assert.Equal(t, 1, targetItem.Episode)
	assert.Equal(t, 10, targetItem.Session)
}