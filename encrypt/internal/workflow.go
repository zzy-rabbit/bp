package internal

import (
	"context"
	"github.com/zzy-rabbit/bp/encrypt/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

type workflow struct {
	processors []api.IProcessor
}

func (s *service) Workflow(ctx context.Context, processors ...api.IProcessor) api.IProcessor {
	return &workflow{processors: processors}
}

func (w *workflow) Process(ctx context.Context, content []byte) ([]byte, xerror.IError) {
	var err xerror.IError
	for _, processor := range w.processors {
		content, err = processor.Process(ctx, content)
		if err != nil {
			return nil, err
		}
	}
	return content, nil
}
