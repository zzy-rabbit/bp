package internal

import (
	"context"
	"encoding/json"
	"github.com/zzy-rabbit/bp/protocol/upload/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xfile"
	"github.com/zzy-rabbit/xtools/xtrace"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

type tusFileInfo struct {
	ID             string `json:"ID"`
	Size           int    `json:"Size"`
	SizeIsDeferred bool   `json:"SizeIsDeferred"`
	Offset         int    `json:"Offset"`
	MetaData       struct {
		Filename string `json:"filename"`
		Filetype string `json:"filetype"`
	} `json:"MetaData"`
	IsPartial      bool        `json:"IsPartial"`
	IsFinal        bool        `json:"IsFinal"`
	PartialUploads interface{} `json:"PartialUploads"`
	Storage        struct {
		Path string `json:"Path"`
		Type string `json:"Type"`
	} `json:"Storage"`
}

func (s *service) MoveFile(ctx context.Context, id string, path string) xerror.IError {
	defer xtrace.Trace(ctx)(id, path)

	if xfile.IsExist(ctx, path) {
		s.ILogger.Error(ctx, "file %s already exist", path)
		return xerror.Extend(xerror.ErrAlreadyExists, path)
	}

	s.FileLock(ctx, id)
	defer s.FileUnlock(ctx, id)

	srcPath := filepath.Join(s.config.RootPath, id)
	if !xfile.IsExist(ctx, srcPath) {
		s.ILogger.Error(ctx, "file %s not exist", srcPath)
		return xerror.Extend(xerror.ErrNotFound, srcPath)
	}

	err := os.Rename(srcPath, path)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "move file %s to %s fail %v", srcPath, path, err)
		return xerror.Extend(xerror.ErrInternalError, "move file %s to %s", srcPath, path)
	}
	err = os.RemoveAll(srcPath + ".info")
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "move file %s to %s success, delete %s fail %v", id, path, srcPath+".info", err)
	}
	return nil
}

func (s *service) CopyFile(ctx context.Context, id string, w io.Writer) (api.FileInfo, xerror.IError) {
	defer xtrace.Trace(ctx)(id)

	s.FileRLock(ctx, id)
	defer s.FileRUnlock(ctx, id)

	fileInfo, xerr := s.GetFileInfo(ctx, id)
	if xerror.Error(xerr) {
		s.ILogger.Error(ctx, "get file info %s fail %v", id, xerror.Error(xerr))
		return api.FileInfo{}, xerr
	}

	filePath := filepath.Join(s.config.RootPath, id)
	file, err := os.Open(filePath)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "open file %s fail %v", filePath, err)
		return api.FileInfo{}, xerror.Extend(xerror.ErrInternalError, "open file "+filePath)
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "copy file %s fail %v", filePath, err)
		return api.FileInfo{}, xerror.Extend(xerror.ErrInternalError, "copy file "+filePath)
	}
	return fileInfo, nil
}

func (s *service) GetFileInfo(ctx context.Context, id string) (api.FileInfo, xerror.IError) {
	defer xtrace.Trace(ctx)(id)

	s.FileRLock(ctx, id)
	defer s.FileRUnlock(ctx, id)

	filePath := filepath.Join(s.config.RootPath, id)
	infoPath := filepath.Join(s.config.RootPath, id+".info")
	if !xfile.IsExist(ctx, filePath) {
		s.ILogger.Error(ctx, "file %s not exist", filePath)
		return api.FileInfo{}, xerror.Extend(xerror.ErrNotFound, filePath)
	}
	if !xfile.IsExist(ctx, infoPath) {
		s.ILogger.Error(ctx, "file info %s not exist", infoPath)
		return api.FileInfo{}, xerror.Extend(xerror.ErrNotFound, infoPath)
	}

	content, err := os.ReadFile(infoPath)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "read file %s fail %v", infoPath, err)
		return api.FileInfo{}, xerror.Extend(xerror.ErrInternalError, "read file "+infoPath)
	}
	var info tusFileInfo
	err = json.Unmarshal(content, &info)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "json unmarshal fail %v", err)
		return api.FileInfo{}, xerror.Extend(xerror.ErrInternalError, "json unmarshal file info fail")
	}

	return api.FileInfo{
		ID:   id,
		Size: info.Size,
		Name: info.MetaData.Filename,
		Type: info.MetaData.Filetype,
		Path: info.Storage.Path,
	}, nil
}

func (s *service) DeleteFile(ctx context.Context, id string) xerror.IError {
	defer xtrace.Trace(ctx)(id)

	defer s.deleteFileSync(ctx, id)

	s.FileLock(ctx, id)
	defer s.FileUnlock(ctx, id)

	filePath := filepath.Join(s.config.RootPath, id)
	if xfile.IsExist(ctx, filePath) {
		s.ILogger.Info(ctx, "delete file %s", filePath)
		err := os.RemoveAll(filePath)
		if xerror.Error(err) {
			s.ILogger.Error(ctx, "delete file %s fail %v", filePath, err)
		}
	}
	infoPath := filepath.Join(s.config.RootPath, id+".info")
	if xfile.IsExist(ctx, infoPath) {
		s.ILogger.Info(ctx, "delete file %s", infoPath)
		err := os.RemoveAll(infoPath)
		if xerror.Error(err) {
			s.ILogger.Error(ctx, "delete file %s fail %v", infoPath, err)
		}
	}
	return nil
}

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
					if filepath.Ext(path) != ".info" {
						return nil
					}
					if int(time.Now().Unix()-info.ModTime().Unix()) <= s.config.Expire {
						return nil
					}

					content, err := os.ReadFile(path)
					if xerror.Error(err) {
						s.ILogger.Error(ctx, "read file %s fail %v", path, err)
						return xerror.Extend(xerror.ErrInternalError, "read file "+path)
					}
					var fileInfo tusFileInfo
					err = json.Unmarshal(content, &fileInfo)
					if xerror.Error(err) {
						s.ILogger.Error(ctx, "json unmarshal fail %v", err)
						return xerror.Extend(xerror.ErrInternalError, "json unmarshal file info fail")
					}

					// 文件正在被使用
					if s.IsFileLocked(ctx, fileInfo.ID) {
						return nil
					}

					defer s.deleteFileSync(ctx, fileInfo.ID)

					s.FileLock(ctx, fileInfo.ID)
					defer s.FileUnlock(ctx, fileInfo.ID)

					s.ILogger.Info(ctx, "expire monitor delete path %s", path)
					err = os.RemoveAll(path)
					if xerror.Error(err) {
						s.ILogger.Error(ctx, "expire monitor delete path %s fail %v", path, err)
					}
					err = os.RemoveAll(path + ".info")
					if xerror.Error(err) {
						s.ILogger.Error(ctx, "expire monitor delete path %s fail %v", path+".info", err)
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
