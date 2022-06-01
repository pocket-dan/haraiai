package api

import (
	"net/http"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/raahii/haraiai/pkg/handler"
)

var apiHandler handler.ApiHandler

func init() {
	var err error

	apiHandler, err = handler.NewApiHandler()
	if err != nil {
		panic(err)
	}
}

func NotifyInquiry(w http.ResponseWriter, req *http.Request) {
	apiHandler.NotifyInquiry(w, req)
}
