package parser

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/playgoround/goneek/core/logger"
	"github.com/playgoround/goneek/core/types"
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	newneekStyleContentTitle = []string{
		"font-size: 26px;",
	}
	newneekStyleContentSubTitle = []string{
		"font-size: 18px;",
	}
	newneekStyleCommentBlock = []string{
		"background: #e6e6e6;",
	}
)

type newneek struct {
	docs []*goquery.Document
	log  *log.Entry
}

func NewNeekParser(articles []types.ArticleInfo) (Parser, error) {
	newneekLogger := logger.New().WithField("parser", "newneek")
	newneekLogger.Debug("Initializing newneek parser")

	docs := make([]*goquery.Document, len(articles))
	for index, article := range articles {
		doc, err := goquery.NewDocumentFromReader(
			bytes.NewReader(article.GetHtml()),
		)
		if err != nil {
			return nil, err
		}
		docs[index] = doc
	}

	parser := &newneek{
		docs: docs,
		log:  newneekLogger,
	}
	parser.log.Info("Initialize Newneek parser")
	return parser, nil
}

func checkStyle(style string, compare []string) bool {
	ok := true
	for _, s := range compare {
		if !strings.Contains(style, s) {
			ok = false
		}
	}
	return ok
}

func (nn newneek) parseCommentBlock(
	blockIndex int,
	selection *goquery.Selection,
) {
	nn.log.WithFields(log.Fields{
		"block": blockIndex,
		"type":  "comment block",
	}).Debug("TODO : parse comment block")
}

func (nn newneek) parsePlainBlock(
	blockIndex int,
	selection *goquery.Selection,
) {
	var contentTitle string
	var contentSubTitle []string
	selection.Find(".stb-text-box span").Each(func(i int, selection *goquery.Selection) {
		spanStyle := selection.AttrOr("style", "")
		isContentTitle := checkStyle(spanStyle, newneekStyleContentTitle)
		isContentSubTitle := checkStyle(spanStyle, newneekStyleContentSubTitle)

		switch {
		case isContentTitle:
			contentTitle += selection.Text()
		case isContentSubTitle:
			contentSubTitle = append(contentSubTitle, selection.Text())
		}
	})

	if contentTitle != "" {
		nn.log.WithFields(log.Fields{
			"block": blockIndex,
			"type":  "content title",
		}).Debug(contentTitle)
	}
	if contentSubTitle != nil {
		for index, subTitle := range contentSubTitle {
			nn.log.WithFields(log.Fields{
				"block": blockIndex,
				"index": index,
				"type":  "content subtitle",
			}).Debug(subTitle)
		}
	}
}

func (nn *newneek) Parse() ([]byte, error) {
	nn.docs[len(nn.docs)-1].Find(".stb-block").Each(func(i int, selection *goquery.Selection) {
		blockStyle := selection.AttrOr("style", "")
		isCommentBlock := checkStyle(blockStyle, newneekStyleCommentBlock)

		if isCommentBlock {
			nn.parseCommentBlock(i, selection)
		} else {
			nn.parsePlainBlock(i, selection)
		}

	})
	return nil, nil
}
