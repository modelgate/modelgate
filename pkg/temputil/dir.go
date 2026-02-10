package temputil

import (
	"os"

	"github.com/modelgate/modelgate/internal/config"
	"github.com/yearnfar/gokit/fsutil"
)

func Mkdir(prefix string) (string, error) {
	cfg := config.GetConfig()
	if !fsutil.IsDir(cfg.Temp.Dir) {
		if err := os.MkdirAll(cfg.Temp.Dir, 0755); err != nil {
			return "", err
		}
	}
	return os.MkdirTemp(cfg.Temp.Dir, prefix)
}
