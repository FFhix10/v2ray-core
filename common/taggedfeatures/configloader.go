package taggedfeatures

import (
	"context"
	"encoding/json"
	"github.com/v2fly/v2ray-core/v4/common/serial"
	"github.com/v2fly/v2ray-core/v4/infra/conf/v5cfg"
	"google.golang.org/protobuf/types/known/anypb"
)

func LoadJsonConfig(ctx context.Context, interfaceType, defaultImpl string, message json.RawMessage) (*Config, error) {
	type ItemStub struct {
		MemberType string          `json:"type"`
		Tag        string          `json:"tag"`
		Value      json.RawMessage `json:"settings"`
	}
	type namedStub []ItemStub
	var stub namedStub
	err := json.Unmarshal(message, &stub)
	if err != nil {
		return nil, err
	}
	config := &Config{Features: map[string]*anypb.Any{}}
	for _, v := range stub {
		if v.MemberType == "" {
			v.MemberType = defaultImpl
		}
		pack, err := v5cfg.LoadHeterogeneousConfigFromRawJson(ctx, interfaceType, v.MemberType, v.Value)
		if err != nil {
			return nil, err
		}
		config.Features[v.Tag] = serial.ToTypedMessage(pack)
	}
	return config, nil
}
