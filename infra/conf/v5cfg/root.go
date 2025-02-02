package v5cfg

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	core "github.com/v2fly/v2ray-core/v4"
	"github.com/v2fly/v2ray-core/v4/app/dispatcher"
	"github.com/v2fly/v2ray-core/v4/app/proxyman"
	"github.com/v2fly/v2ray-core/v4/common/platform"
	"github.com/v2fly/v2ray-core/v4/common/serial"
	"github.com/v2fly/v2ray-core/v4/infra/conf/cfgcommon"
	"github.com/v2fly/v2ray-core/v4/infra/conf/geodata"
	"github.com/v2fly/v2ray-core/v4/infra/conf/synthetic/log"
	"google.golang.org/protobuf/types/known/anypb"
)

func (c RootConfig) BuildV5(ctx context.Context) (proto.Message, error) {
	config := &core.Config{
		App: []*anypb.Any{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	}

	var logConfMsg *anypb.Any
	if c.LogConfig != nil {
		logConfMsgUnpacked, err := loadHeterogeneousConfigFromRawJson("service", "log", c.LogConfig)
		if err != nil {
			return nil, err
		}
		logConfMsg = serial.ToTypedMessage(logConfMsgUnpacked)
	} else {
		logConfMsg = serial.ToTypedMessage(log.DefaultLogConfig())
	}
	// let logger module be the first App to start,
	// so that other modules could print log during initiating
	config.App = append([]*anypb.Any{logConfMsg}, config.App...)

	if c.RouterConfig != nil {
		routerConfig, err := loadHeterogeneousConfigFromRawJson("service", "router", c.RouterConfig)
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(routerConfig))
	}

	if c.DNSConfig != nil {
		dnsApp, err := loadHeterogeneousConfigFromRawJson("service", "dns", c.DNSConfig)
		if err != nil {
			return nil, newError("failed to parse DNS config").Base(err)
		}
		config.App = append(config.App, serial.ToTypedMessage(dnsApp))
	}

	for _, rawInboundConfig := range c.Inbounds {
		ic, err := rawInboundConfig.BuildV5(ctx)
		if err != nil {
			return nil, err
		}
		config.Inbound = append(config.Inbound, ic.(*core.InboundHandlerConfig))
	}

	for _, rawOutboundConfig := range c.Outbounds {
		ic, err := rawOutboundConfig.BuildV5(ctx)
		if err != nil {
			return nil, err
		}
		config.Outbound = append(config.Outbound, ic.(*core.OutboundHandlerConfig))
	}

	for serviceName, service := range c.Services {
		servicePackedConfig, err := loadHeterogeneousConfigFromRawJson("service", serviceName, service)
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(servicePackedConfig))
	}
	return config, nil
}

func loadJsonConfig(data []byte) (*core.Config, error) {
	rootConfig := &RootConfig{}

	err := json.Unmarshal(data, rootConfig)
	if err != nil {
		return nil, newError("unable to load json").Base(err)
	}

	buildctx := cfgcommon.NewConfigureLoadingContext(context.Background())

	geoloadername := platform.NewEnvFlag("v2ray.conf.geoloader").GetValue(func() string {
		return "standard"
	})

	if loader, err := geodata.GetGeoDataLoader(geoloadername); err == nil {
		cfgcommon.SetGeoDataLoader(buildctx, loader)
	} else {
		return nil, newError("unable to create geo data loader ").Base(err)
	}

	message, err := rootConfig.BuildV5(buildctx)
	if err != nil {
		return nil, newError("unable to build config").Base(err)
	}
	return message.(*core.Config), nil
}
