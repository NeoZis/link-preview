package handlers

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

type StandardLinkPreview struct {
	*LinkPreviewContext
}

func (p *StandardLinkPreview) Preview() (*LinkPreviewContext, error) {
	err := p.request()

	if nil != err {
		return nil, err
	}

	err = p.readTags()
	if nil != err {
		return nil, err
	}

	return p.PreviewContext(), nil
}

func (p *StandardLinkPreview) readTags() error {
	titleNode := p.Parsed.Find("html > head > title")
	if 0 == titleNode.Length() {
		return errors.New("title not found")
	}
	p.Title = titleNode.Text()
	// Parse <meta> tags.
	metaNodes := p.Parsed.Find("html > head > meta")
	for _, node := range metaNodes.Nodes {
		for _, attr := range node.Attr {
			switch strings.ToLower(attr.Key) {
			case "property":
				p.parseMetaProperties(attr.Val, node)
				break
			case "itemprop":
				if "image" == strings.ToLower(attr.Val) && "" == p.ImageURL {
					content := p.parseMetaContent(node)
					if preparedLink := p.prepareLink(content); "" != preparedLink && p.checkAccessToLink(preparedLink) {
						p.ImageURL = preparedLink
					}

					break
				}
			case "name":
				if "description" == strings.ToLower(attr.Val) && "" == p.Description {
					content := p.parseMetaContent(node)
					p.Description = content
					break
				}
			default:
				continue
			}
		}
	}

	// Find `favicon.ico`.
	linkNodes := p.Parsed.Find("html > head > link")
	for _, node := range linkNodes.Nodes {
		for _, attr := range node.Attr {
			switch strings.ToLower(attr.Key) {
			case "rel":
				if (!strings.Contains(attr.Val, "icon") || attr.Val == "mask-icon") && attr.Val != "apple-touch-icon-precomposed" {
					break
				}
				status := p.parseFavicon(node)
				if status {
					break
				}
				// need to break after success
			default:
				continue
			}
		}
	}

	return nil
}

func (p *StandardLinkPreview) parseMetaContent(node *html.Node) string {
	var content string
	for _, attr := range node.Attr {
		switch strings.ToLower(attr.Key) {
		case "content":
			content = attr.Val
			break
		default:
			continue
		}
	}

	return content
}

func (p *StandardLinkPreview) parseMetaProperties(nodeType string, node *html.Node) {
	nodeType = strings.ToLower(nodeType)

	if !strings.HasPrefix(nodeType, "og:") {
		return
	}

	slices := strings.Split(nodeType, ":")
	if 2 != len(slices) {
		return
	}

	nodeType = slices[1]
	content := p.parseMetaContent(node)

	switch nodeType {
	case "description":
		p.Description = content
	case "image":
		if preparedLink := p.prepareLink(content); "" != preparedLink && p.checkAccessToLink(preparedLink) {
			p.ImageURL = preparedLink
		}
	case "title":
		p.Title = content
	}
}
