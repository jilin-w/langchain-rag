package tika_service

import (
	"context"
	"io"

	"github.com/google/go-tika/tika"
)

func ReadFile(reader io.ReadCloser) (res io.ReadCloser, err error) {
	client := tika.NewClient(nil, "http://127.0.0.1:9998")
	res, err = client.ParseReader(context.Background(), reader)
	return res, err
}
