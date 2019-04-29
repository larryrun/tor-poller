package mags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractResolutionInfo(t *testing.T) {
	title := "NCIS.Los.Angeles.S10E01.1080p.AMZN.WEB-DL.DDP5.1.H.264-ViSUM.mkv"
	resolution, err := ExtractResolutionInfo(title)
	assert.Nil(t, err)
	assert.Equal(t, "1080p", resolution)

}

func TestExtractEpisodeInfo(t *testing.T) {
	title := "NCIS.Los.Angeles.S10E01.1080p.AMZN.WEB-DL.DDP5.1.H.264-ViSUM.mkv"
	episode, err := ExtractEpisodeInfo(title)
	assert.Nil(t, err)
	assert.Equal(t, 10, episode.Season)
	assert.Equal(t, 1, episode.Episode)

	title = "NCIS.Los.Angeles.S10E1001.1080p.AMZN.WEB-DL.DDP5.1.H.264-ViSUM.mkv"
	episode, err = ExtractEpisodeInfo(title)
	assert.NotNil(t, err)
}

func TestExtractItemSize(t *testing.T) {
	size, err := ExtractItemSize("1003.74MB")
	assert.Nil(t, err)
	assert.Equal(t, 1003, size)

	size, err = ExtractItemSize("2.61GB")
	assert.Nil(t, err)
	assert.Equal(t, 2672, size)
}

func TestGetLatestEpisodeInTitles(t *testing.T) {
	titles := make([]string, 5)
	titles[0] = "NCIS.Los.Angeles.S1E01.1080p.AMZN.WEB-DL.DDP5.1.H.264-ViSUM.mkv"
	titles[1] = "NCIS.Los.Angeles.S7E03.1080p.AMZN.WEB-DL.DDP5.1.H.264-ViSUM.mkv"
	titles[2] = "NCIS.Los.Angeles.S10E04.1080p.AMZN.WEB-DL.DDP5.1.H.264-ViSUM.mkv"
	titles[3] = "NCIS.Los.Angeles.S10E05.1080p.AMZN.WEB-DL.DDP5.1.H.264-ViSUM.mkv"
	titles[4] = "NCIS.Los.Angeles.S10E08.1080p.AMZN.WEB-DL.DDP5.1.H.264-ViSUM.mkv"
	latestEpisode := GetLatestEpisodeInTitles(titles)
	assert.Equal(t, latestEpisode.Season, 10)
	assert.Equal(t, latestEpisode.Episode, 8)
}