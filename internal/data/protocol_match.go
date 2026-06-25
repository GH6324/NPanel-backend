package data

import (
	"strings"

	servermodel "github.com/npanel-dev/NPanel-backend/internal/model/server"
)

func matchNodeProtocolConfig(protocols []*servermodel.Protocol, nodeProtocol string, nodePort uint16) (matched, firstEnabled, firstAvailable *servermodel.Protocol) {
	targetType := strings.ToLower(strings.TrimSpace(nodeProtocol))
	var typeFallback *servermodel.Protocol
	var loosePortMatch *servermodel.Protocol

	for _, protocol := range protocols {
		if protocol == nil {
			continue
		}
		if firstAvailable == nil {
			firstAvailable = protocol
		}
		if protocol.Enable && firstEnabled == nil {
			firstEnabled = protocol
		}
		if strings.ToLower(strings.TrimSpace(protocol.Type)) != targetType {
			continue
		}
		if typeFallback == nil {
			typeFallback = protocol
		}
		if nodePort > 0 && protocol.Port > 0 && uint16(protocol.Port) == nodePort {
			return protocol, firstEnabled, firstAvailable
		}
		if loosePortMatch == nil && (nodePort == 0 || protocol.Port == 0) {
			loosePortMatch = protocol
		}
	}

	if loosePortMatch != nil {
		return loosePortMatch, firstEnabled, firstAvailable
	}
	if nodePort > 0 {
		return nil, firstEnabled, firstAvailable
	}
	return typeFallback, firstEnabled, firstAvailable
}
