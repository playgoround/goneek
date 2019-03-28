package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"testing"

	"github.com/playgoround/goneek/core/types"
	"github.com/stretchr/testify/require"
)

const (
	testNewneekDataPath = "testdata/newneek"
)

func TestNewNeek_Parse(t *testing.T) {
	fileInfos, err := ioutil.ReadDir(testNewneekDataPath)
	require.NoError(t, err)

	articleInfos := make([]types.ArticleInfo, len(fileInfos))
	for index, fileInfo := range fileInfos {
		html, err := ioutil.ReadFile(path.Join(testNewneekDataPath, fileInfo.Name()))
		require.NoError(t, err)

		articleInfos[index] = types.NewArticleInfo(
			fmt.Sprintf("newneek-%d", index),
			"", html,
		)
	}

	parser, err := NewNeekParser(articleInfos)
	require.NoError(t, err)

	articles, err := parser.Parse()
	require.NoError(t, err)
	log.Println(articles)
}
