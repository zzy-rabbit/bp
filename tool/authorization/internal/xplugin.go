package internal

import (
	"context"
	"encoding/json"
	"github.com/zzy-rabbit/bp/protocol/authorization/api"
	encryptApi "github.com/zzy-rabbit/bp/tool/encrypt/api"
	logApi "github.com/zzy-rabbit/bp/tool/log/api"
	"github.com/zzy-rabbit/xtools/xerror"
)

type service struct {
	ILogger   logApi.IPlugin     `xplugin:"bp.tool.log"`
	IEncrypt  encryptApi.IPlugin `xplugin:"bp.tool.encrypt"`
	encoder   encryptApi.IProcessor
	decoder   encryptApi.IProcessor
	config    api.Config
	secretKey []byte
}

func New(ctx context.Context) api.IPlugin {
	return &service{}
}

func (s *service) GetName(ctx context.Context) string {
	return api.PluginName
}

func (s *service) Init(ctx context.Context, initParam string) error {
	err := json.Unmarshal([]byte(initParam), &s.config)
	if xerror.Error(err) {
		s.ILogger.Error(ctx, "plugin %s init fail %v", s.GetName(ctx), err)
		return err
	}
	s.secretKey = generateRandomKey(ctx, s.config.SecretKeyLength)

	s.encoder = s.IEncrypt.Workflow(ctx,
		s.IEncrypt.NewAESEncoder(ctx, []byte(s.secretKey)),
		s.IEncrypt.NewBase64Encoder(ctx),
	)
	s.decoder = s.IEncrypt.Workflow(ctx,
		s.IEncrypt.NewBase64Decoder(ctx),
		s.IEncrypt.NewAESDecoder(ctx, []byte(s.secretKey)),
	)

	s.ILogger.Info(ctx, "plugin %s init success", s.GetName(ctx))
	return nil
}

func (s *service) Run(ctx context.Context, runParam string) error {
	s.ILogger.Info(ctx, "plugin %s run success", s.GetName(ctx))
	return nil
}

func (s *service) Stop(ctx context.Context, stopParam string) error {
	s.ILogger.Info(ctx, "plugin %s stop success", s.GetName(ctx))
	return nil
}
