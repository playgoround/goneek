package parser

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/playgoround/goneek/core/logger"
	"github.com/playgoround/goneek/core/types"
	log "github.com/sirupsen/logrus"
	"strings"
)

type nnBLockType string

const (
	nnBlockNothing = nnBLockType(iota)
	nnBlockTitle
	nnBlockContents
	nnBlockComment
	nnBlock10Min
	nnBlockFollowUp
)

var (
	nnStyleTitle           = []string{"font-size: 26px;"}
	nnStyleSubTitle        = []string{"font-size: 18px;"}
	nnStyleCommentBlock    = []string{"background: #e6e6e6;"}
	nnStyleCommentSubTitle = []string{"font-weight: bold;"}
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

func getBlockType(selection *goquery.Selection) (nnBLockType, string) {
	blockStyle := selection.AttrOr("style", "")

	// check if its comment block
	if checkStyle(blockStyle, nnStyleCommentBlock) {
		return nnBlockComment, ""
	}

	// check if its title/contents block
	var title, subTitle string
	selection.Find(".stb-text-box span").Each(
		func(i int, selection *goquery.Selection) {
			spanStyle := selection.AttrOr("style", "")
			if checkStyle(spanStyle, nnStyleTitle) {
				title += selection.Text()
			}
			if checkStyle(spanStyle, nnStyleSubTitle) {
				subTitle += selection.Text()
			}
		})

	if title != "" {
		if strings.Contains(title, "10분 더 있다면 읽어 볼 거리") {
			return nnBlock10Min, title
		}
		return nnBlockTitle, title
	}

	if subTitle != "" {
		return nnBlockContents, ""
	}
	return nnBlockNothing, ""
}

type blockInfo struct {
	blockType nnBLockType
	blockData string
}

func (nn *newneek) Parse() ([]byte, error) {
	var blockIndex []blockInfo
	blocks := nn.docs[len(nn.docs)-1].Find(".stb-block")
	blocks.Each(
		func(i int, block *goquery.Selection) {
			blockType, blockData := getBlockType(block)
			blockIndex = append(blockIndex, blockInfo{blockType, blockData})

			switch blockType {
			case nnBlockTitle:
				nn.parseBlockTitle(block, i)
			case nnBlockContents:
				// previous block type
				prevBlock := blockIndex[i-1]
				if prevBlock.blockType == nnBlock10Min {
					return
				}
				if prevBlock.blockType != nnBlockTitle {
					nn.log.Error(prevBlock.blockType)
					return
				}
				nn.parseBlockContents(prevBlock.blockData, block, i)
			case nnBlockComment:
				if i == len(blocks.Nodes)-2 {
					return
				}
				nn.log.WithFields(log.Fields{
					"block": i,
					"type":  "comment",
				}).Debug("block comment")
				nn.parseBlockComment(block, i)
			case nnBlock10Min:
			default:
				nn.log.WithFields(log.Fields{
					"block": i,
					"type":  "nothing",
				}).Debug("block nothing")
			}
		})

	return nil, nil
}

func (nn *newneek) parseBlockTitle(
	block *goquery.Selection,
	blockIndex int,
) (title string) {
	block.Find(".stb-text-box span").Each(
		func(i int, span *goquery.Selection) {
			spanStyle := span.AttrOr("style", "")
			if checkStyle(spanStyle, nnStyleTitle) {
				title += span.Text()
			}
		})
	return
}

type contentsInfo struct {
	title    string
	contents string
}

func (nn *newneek) parseBlockContents(
	title string,
	block *goquery.Selection,
	blockIndex int,
) (subtitle string, contents []contentsInfo) {
	var contentsBuilder strings.Builder

	subTitle := &sepIdxArr{index: 0}

	box := block.Find(".stb-text-box").Children()
	textBoxDiv := block.Find(".stb-text-box > div")
	if len(textBoxDiv.Nodes) != 0 {
		box = textBoxDiv.Children()
	}

	box.Each(
		func(boxIndex int, boxContents *goquery.Selection) {
			// get full contents
			boxContentsTagName := goquery.NodeName(boxContents)
			switch boxContentsTagName {
			case "span":
				fallthrough
			case "div":
				fallthrough
			case "ul":
				contentsBuilder.WriteString(boxContents.Text())
			}
			boxContents.Children().Each(
				nn.parseBlockContentsElement(boxContents, subTitle))
		})

	contentsText := contentsBuilder.String()
	contents = append(contents, contentsInfo{
		title:    title,
		contents: strings.Split(contentsText, subTitle.array[0])[0],
	})

	for index, sub := range subTitle.array {
		info := contentsInfo{title: sub}
		if index == len(subTitle.array)-1 {
			info.contents = strings.Split(contentsText, sub)[1]
		} else {
			info.contents = strings.Split(
				strings.Split(contentsText, sub)[1],
				subTitle.array[index+1],
			)[0]
		}
		contents = append(contents, info)
	}

	for _, c := range contents {
		nn.log.WithFields(log.Fields{
			"type": "title",
		}).Info(c.title)
		nn.log.WithFields(log.Fields{
			"type": "contents",
		}).Debug(c.contents)
	}
	return
}

func (nn *newneek) parseBlockContentsElement(
	boxContents *goquery.Selection,
	subTitle *sepIdxArr,
) func(int, *goquery.Selection) {
	return func(idx int, elem *goquery.Selection) {
		if goquery.NodeName(elem) == "b" {
			elem = elem.Find("span")
		}

		if goquery.NodeName(elem) == "span" {
			depth := 0
		CheckStyle:
			spanStyle := elem.AttrOr("style", "")
			if !checkStyle(spanStyle, nnStyleSubTitle) {
				if depth < 1 {
					elem = elem.Find("span")
					depth++
					goto CheckStyle
				}
				return
			}

			if idx != 0 {
				prevNode := boxContents.Children().Get(idx - 1)
				prevNodeStyle := getAttrFromNode(prevNode, "style")
				if prevNode.Data == "span" {
					if checkStyle(prevNodeStyle, nnStyleSubTitle) {
						subTitle.array[subTitle.index-1] += elem.Text()
						return
					}
				}
			}

			subTitle.array = append(subTitle.array, elem.Text())
			subTitle.index += 1
		}
	}
}

func (nn *newneek) parseBlockComment(
	block *goquery.Selection,
	blockIndex int,
) (subtitle string, contents map[string]interface{}) {
	block.Find(".stb-text-box").Each(
		func(i int, textBox *goquery.Selection) {
			textBox.Children().Each(
				func(i int, selection *goquery.Selection) {
					nn.log.Debug(goquery.NodeName(selection))
				})
		})
	return
}
