package main

import (
	"fmt"
	"net/http"

	"github.com/NeoZis/link-preview/handlers"
)

func main() {
	data, err := Preview("https://youtrack.jcdev.net/issue/RL-99/Backoffice-Company-members-page-doesnt-reload-after-deleting-a-company-member-with-2FA", nil)

	if nil != err {
		panic(err)
	}

	fmt.Println(data.ImageURL)
}

func Preview(link string, extraContent *http.Request) (*handlers.LinkPreviewContext, error) {
	return PreviewLink(link, extraContent)
}

func PreviewLink(link string, extraClient *http.Request) (*handlers.LinkPreviewContext, error) {
	cxt := &handlers.LinkPreviewContext{
		Link:       link,
		TargetType: handlers.StandardMetaTags,
	}

	if nil != extraClient {
		cxt.Client = extraClient
	}

	handler, handlerErr := handlers.GetPreviewHandler(cxt)
	if nil != handlerErr {
		return nil, handlerErr
	}

	cxt, previewErr := handler.Preview()
	if nil != previewErr {
		return nil, previewErr
	}

	return cxt, nil
}
