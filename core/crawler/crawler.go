package crawler

import "github.com/playgoround/goneek/core/types"

// Crawler is interface of all type of crawlers (newneek etc..)
type Crawler interface {
	Get() ([]types.ArticleInfo, error)
}
