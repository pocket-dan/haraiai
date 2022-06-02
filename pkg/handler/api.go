package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/line"
	"github.com/raahii/haraiai/pkg/config"
	"github.com/raahii/haraiai/pkg/log"
)

const (
	MESSAGE_SUBJECT = "[問い合わせがありました]"
)

type ApiHandler interface {
	NotifyInquiry(http.ResponseWriter, *http.Request)
}

type ApiHandlerImpl struct {
	config   config.ApiConfig
	notifier notify.Notifier
}

func NewApiHandler() (*ApiHandlerImpl, error) {
	c, err := config.NewApiConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize api config: %w", err)
	}

	receiverToken := os.Getenv("LINE_NOTIFY_TOKEN")
	if receiverToken == "" {
		return nil, errors.New("$LINE_NOTIFY_TOKEN required.")
	}

	n := line.NewNotify()
	n.AddReceivers(receiverToken)

	return &ApiHandlerImpl{
		config:   c,
		notifier: n,
	}, nil
}

type NotifyInquiryRequest struct {
	Text string
}

func (ah *ApiHandlerImpl) NotifyInquiry(w http.ResponseWriter, req *http.Request) {
	logger := log.NewLogger(req)

	// Allow OPTIONS(preflight) request and POST request only.
	if !(req.Method == http.MethodPost || req.Method == http.MethodOptions) {
		logger.Warnf("unsuppoed request method: %s", req.Method)
		http.Error(w, "Method Not Allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Set CORS headers for the preflight request.
	// HACK: CORS code should be common.
	if req.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", ah.config.GetFrontOrigin())
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", ah.config.GetFrontOrigin())

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("failed to read request body: %s", err)
		http.Error(w, "Bad Request.", http.StatusBadRequest)
		return
	}

	notifyInquiryRequest := new(NotifyInquiryRequest)
	err = json.Unmarshal(bodyBytes, notifyInquiryRequest)
	if err != nil {
		logger.Errorf("failed to unmarshal request json: %s", err)
		http.Error(w, "Bad Request.", http.StatusBadRequest)
		return
	}

	err = ah.notifier.Send(context.Background(), MESSAGE_SUBJECT, notifyInquiryRequest.Text)
	if err != nil {
		logger.Errorf("failed to send line notify message: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
