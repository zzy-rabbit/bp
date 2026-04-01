package internal

import (
	"context"
	"github.com/zzy-rabbit/xtools/xerror"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

func (s *service) startExpireMonitor(ctx context.Context) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				s.ILogger.Error(ctx, "tusd handler panic %v %s", err, debug.Stack())
			}
		}()
		for {
			select {
			case <-ctx.Done():
				s.ILogger.Info(ctx, "expire monitor exist")
				return
			case <-time.After(time.Second * time.Duration(s.config.Interval)):
				err := filepath.Walk(s.config.RootPath, func(path string, info fs.FileInfo, err error) error {
					if xerror.Error(err) {
						s.ILogger.Error(ctx, "expire monitor walk path %s fail %v", path, err)
						return err
					}
					if info.IsDir() {
						return nil
					}
					if int(time.Now().Unix()-info.ModTime().Unix()) > s.config.Expire {
						s.ILogger.Info(ctx, "expire monitor delete path %s", path)
						err = os.RemoveAll(path)
						if xerror.Error(err) {
							s.ILogger.Error(ctx, "expire monitor delete path %s fail %v", path, err)
						}
					}
					return nil
				})
				if xerror.Error(err) {
					s.ILogger.Error(ctx, "expire monitor walk path %s fail %v", s.config.RootPath, err)
				}
			}
		}
	}()
}
