package moviedownloader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {
	configPath = "./testdata/cfg.yaml"
	config, err := ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "testDownloadingFolder", config.DownloadFolder)
	assert.Equal(t, 2, len(config.Items))
	assert.Equal(t, "name1", config.Items[0].Name)
	assert.Equal(t, "type1", config.Items[0].Type)
	assert.Equal(t, "info1", config.Items[0].Info)
	assert.Equal(t, 1, config.Items[0].Cond["start-session"])
	assert.Equal(t, 1, config.Items[0].Cond["start-episode"])
	assert.Equal(t, "name2", config.Items[1].Name)
	assert.Equal(t, "type2", config.Items[1].Type)
	assert.Equal(t, "info2", config.Items[1].Info)
	assert.Equal(t, 2, config.Items[1].Cond["start-session"])
	assert.Equal(t, 2, config.Items[1].Cond["start-episode"])
}