package arkruntime

import (
	"context"
	"net/http"

	"github.com/byteplus-sdk/byteplus-go-sdk-v2/service/arkruntime/model"
)

const batchChatCompletionsSuffix = "/batch/chat/completions"

func newBatchHTTPClient(maxParallel int) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: maxParallel,
		},
	}
}

// CreateBatchChatCompletion — API call to Create a batch completion for the chat message.
func (c *Client) CreateBatchChatCompletion(
	ctx context.Context,
	request model.ChatRequest,
	setters ...requestOption,
) (response model.ChatCompletionResponse, err error) {
	if request.IsStream() {
		err = model.ErrChatCompletionStreamNotSupported
		return
	}
	requestOptions := append(setters, withBody(request))
	err = c.DoBatch(ctx, http.MethodPost, c.fullURL(batchChatCompletionsSuffix), resourceTypeEndpoint, request.GetModel(), &response, requestOptions...)
	if err != nil {
		return
	}
	return
}
