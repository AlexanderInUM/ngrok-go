package libngrok

import (
	"net"

	"github.com/ngrok/libngrok-go/internal/pb_agent"
)

type CommonConfig struct {
	CIDRRestrictions *CIDRRestriction
	ProxyProto       ProxyProtoVersion
	Metadata         string
	ForwardsTo       string
}

type ProxyProtoVersion int32

const (
	ProxyProtoV1 = ProxyProtoVersion(1)
	ProxyProtoV2 = ProxyProtoVersion(2)
)

func (cfg *CommonConfig) WithProxyProto(version ProxyProtoVersion) *CommonConfig {
	cfg.ProxyProto = version
	return cfg
}

func (cfg *CommonConfig) WithMetadata(meta string) *CommonConfig {
	cfg.Metadata = meta
	return cfg
}

func (cfg *CommonConfig) WithForwardsTo(address string) *CommonConfig {
	cfg.ForwardsTo = address
	return cfg
}

type CIDRRestriction struct {
	Allowed []string
	Denied  []string
}

func CIDRSet() *CIDRRestriction {
	return &CIDRRestriction{}
}

func (cr *CIDRRestriction) AllowString(cidr ...string) *CIDRRestriction {
	cr.Allowed = append(cr.Allowed, cidr...)
	return cr
}

func (cr *CIDRRestriction) Allow(net ...*net.IPNet) *CIDRRestriction {
	for _, n := range net {
		cr.AllowString(n.String())
	}
	return cr
}

func (cr *CIDRRestriction) DenyString(cidr ...string) *CIDRRestriction {
	cr.Denied = append(cr.Denied, cidr...)
	return cr
}

func (cr *CIDRRestriction) Deny(net ...*net.IPNet) *CIDRRestriction {
	for _, n := range net {
		cr.DenyString(n.String())
	}
	return cr
}

func (ir *CIDRRestriction) toProtoConfig() *pb_agent.MiddlewareConfiguration_IPRestriction {
	if ir == nil {
		return nil
	}

	return &pb_agent.MiddlewareConfiguration_IPRestriction{
		AllowCIDRs: ir.Allowed,
		DenyCIDRs:  ir.Denied,
	}
}

func (cfg *CommonConfig) WithCIDRRestriction(set ...*CIDRRestriction) *CommonConfig {
	if cfg.CIDRRestrictions == nil {
		cfg.CIDRRestrictions = CIDRSet()
	}

	for _, s := range set {
		if s != nil {
			cfg.CIDRRestrictions.AllowString(s.Allowed...)
			cfg.CIDRRestrictions.DenyString(s.Denied...)
		}
	}
	return cfg
}
