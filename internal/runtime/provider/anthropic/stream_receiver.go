package anthropic

import (
	"bufio"
	"io"
	"strings"

	"github.com/modelgate/modelgate/internal/runtime/core"
)

// StreamReceiver 流式接收器
type StreamReceiver struct {
	reader *bufio.Reader
	body   io.ReadCloser
}

func newStreamReceiver(body io.ReadCloser) core.Stream {
	return &StreamReceiver{
		reader: bufio.NewReader(body),
		body:   body,
	}
}

// Recv 接收
func (s *StreamReceiver) Recv() (*core.StreamChunk, error) {
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return &core.StreamChunk{Finish: true}, err
			}
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}
		if line == "data: [DONE]" {
			return &core.StreamChunk{Finish: true}, io.EOF
		}
		payload := strings.TrimPrefix(line, "data: ")
		if payload == "" {
			continue
		}
		return &core.StreamChunk{
			Data: payload,
		}, nil
	}
}

// Close 关闭
func (s *StreamReceiver) Close() error {
	return s.body.Close()
}
