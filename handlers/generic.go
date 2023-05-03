package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	StandardMetaTags = iota
	WeChatMP
)

type PreviewHandler interface {
	PreviewContext() *LinkPreviewContext
	Preview() (*LinkPreviewContext, error)
}

func GetPreviewHandler(c *LinkPreviewContext) (PreviewHandler, error) {
	if nil == c {
		return nil, errors.New("bad link preview cxt, nil given")
	}

	if nil == c.Client {
		c.initClient()
	}

	var handler PreviewHandler

	switch c.TargetType {
	case StandardMetaTags:
		handler = &StandardLinkPreview{
			c,
		}
	default:
		return nil, errors.New("unknown target type")
	}

	return handler, nil
}

type HTMLMetaAttr struct {
	Key   string
	Value string
}

type LinkPreviewContext struct {
	TargetType  int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image"`
	Link        string `json:"website"`

	Client *http.Request     `json:"-"`
	Parsed *goquery.Document `json:"-"`
}

func (p *LinkPreviewContext) PreviewContext() *LinkPreviewContext {
	return p
}

func (p *LinkPreviewContext) initClient() {
	client, _ := http.NewRequest("GET", p.Link, nil)
	p.Client = client
}

func (p *LinkPreviewContext) checkAccessToLink(link string) bool {
	client, _ := http.NewRequest("GET", link, nil)
	res, err := http.DefaultClient.Do(client)
	if nil != err || res.StatusCode != 200 {
		return false
	}
	defer res.Body.Close()

	return true
}

func (p *LinkPreviewContext) request() error {
	res, err := http.DefaultClient.Do(p.Client)
	if nil != err {
		return err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if nil != err {
		return err
	}

	p.Parsed = doc
	return nil
}

func (p *LinkPreviewContext) parseFavicon(node *html.Node) bool {
	var link string

	for _, attr := range node.Attr {
		switch strings.ToLower(attr.Key) {
		case "href":
			link = attr.Val
			break
		default:
			continue
		}
	}

	if "" == link {
		return false
	}

	if "" == p.ImageURL {
		if preparedLink := p.prepareLink(link); "" != preparedLink && p.checkAccessToLink(preparedLink) {
			p.ImageURL = preparedLink

			return true
		}
	}

	return false
}

func (p *LinkPreviewContext) prepareLink(link string) string {
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return link
	}

	parsedURL, _ := url.Parse(p.Link)
	joinedURL := url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
		Path:   link,
	}

	link = joinedURL.String()

	return link
}
