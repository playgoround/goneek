package crawler

import (
	"encoding/binary"
	"github.com/playgoround/goneek/core/types"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playgoround/goneek/core/logger"
	log "github.com/sirupsen/logrus"
)

const (
	// hard-coded url
	NewNeekLibrary = "https://newneek.co/library"
)

type newneek struct {
	url    string
	client *http.Client
	log    *log.Entry
}

// NewNeekCrawler generates articles crawler for newneek
func NewNeekCrawler(
	url string,
	client *http.Client,
) (Crawler, error) {
	newneekLogger := logger.New().WithField("crawler", "newneek")
	newneekLogger.Debug("Initializing newneek crawler")

	// if client is nil, it sets default http client
	if client == nil {
		client = http.DefaultClient
	}

	// initialize newneek struct
	crawler := &newneek{
		url:    url,
		client: client,
		log:    newneekLogger,
	}
	crawler.log.Info("Initialized Newneek crawler")
	return crawler, nil
}

// Get crawls articles in given newneek library page link.
func (nn *newneek) Get() ([]types.ArticleInfo, error) {
	nn.log.WithField("url", nn.url).Debug("Fetching document from given url")
	resp, err := nn.client.Get(nn.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	nn.log.Info("Successfully fetched document")

	nn.log.Debug("Decoding document from response body")
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	nn.log.Debug("Crawling article infos from decoded document")
	var crawlInfos []struct {
		title string
		url   string
	}
	doc.Find(".text-table").Each(func(i int, selection *goquery.Selection) {
		var title, url string
		selection.Find("h5 span").Each(
			func(i int, selection *goquery.Selection) {
				style := selection.AttrOr("style", "")
				if strings.Contains(style, "font-size: 26px") {
					title = selection.Text()
					nn.log.WithFields(log.Fields{
						"index": i,
						"title": title,
					}).Debug("Got article title")
				}
			})

		selection.Find("p .btn").Each(
			func(i int, selection *goquery.Selection) {
				href := selection.AttrOr("href", "")
				if strings.HasPrefix(href, "https://stib.ee/") {
					url = href
					nn.log.WithFields(log.Fields{
						"index": i,
						"url":   url,
					}).Debug("Got article url")
				}
			})

		if title != "" && url != "" {
			crawlInfos = append(crawlInfos, struct {
				title string
				url   string
			}{title: title, url: url})
		}
	})
	nn.log.Info("Successfully got article informations")

	nn.log.WithField("count", len(crawlInfos)).Debug("Fetching articles from given urls")
	infos := make([]types.ArticleInfo, len(crawlInfos))
	for index, info := range crawlInfos {
		resp, err := nn.client.Get(info.url)
		if err != nil {
			return nil, err
		}

		article, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}
		nn.log.WithFields(log.Fields{
			"size": binary.Size(article),
			"url":  info.url,
		}).Debug("Fetched article")

		infos[index] = types.NewArticleInfo(info.title, info.url, article)

		if err := resp.Body.Close(); err != nil {
			return nil, err
		}
	}
	nn.log.WithField("count", len(infos)).Info("Successfully fetched articles")
	return infos, nil
}
