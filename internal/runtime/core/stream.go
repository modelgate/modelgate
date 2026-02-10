package core

// Stream 流式处理
type Stream interface {
	Recv() (*StreamChunk, error)
	Close() error
}

// StreamChunk 流式处理chunk
type StreamChunk struct {
	Data   string
	Finish bool
}

// StreamWriter 流式写入器
type StreamWriter interface {
	Open() error
	Write(chunk *StreamChunk) error
	Close() error
}
