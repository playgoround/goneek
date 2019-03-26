package crawler

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/airbloc/logger"
	"net/http"
)

const (
	NewNeekLibrary = "https://newneek.co/library"
)

type NewNeek struct {
	url    string
	client *http.Client
	log    *log.Logger
}

func NewNeekCrawler() (Crawler, error) {
	crawler := &NewNeek{
		url:    NewNeekLibrary,
		client: http.DefaultClient,
		log:    log.New("newneek crawler"),
	}

	crawler.log.Info("Initialize Newneek crawler")
	return crawler, nil
}

func (nn *NewNeek) Get() ([][]byte, error) {
	nn.log.Debug("fetching document")
	resp, err := nn.client.Get(nn.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	nn.log.Debug("decoding document from reader")
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("btn btn-sm").Each(func(i int, selection *goquery.Selection) {
		selection.
		val, exists := selection.Attr("href")
		if !exists {

		}
	})

	return nil, nil
}
