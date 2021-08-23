package windscribe

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (w *Windscribe) GetConnection(selection configuration.ServerSelection) (
	connection models.Connection, err error) {
	port := getPort(selection)
	protocol := utils.GetProtocol(selection)

	servers, err := w.filterServers(selection)
	if err != nil {
		return connection, err
	}

	var connections []models.Connection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connection := models.Connection{
				Type:     selection.VPN,
				IP:       IP,
				Port:     port,
				Protocol: protocol,
				Hostname: server.OvpnX509,
				PubKey:   server.WgPubKey,
			}
			connections = append(connections, connection)
		}
	}

	if selection.TargetIP != nil {
		return utils.GetTargetIPConnection(connections, selection.TargetIP)
	}

	return utils.PickRandomConnection(connections, w.randSource), nil
}

func getPort(selection configuration.ServerSelection) (port uint16) {
	switch selection.VPN {
	case constants.Wireguard:
		customPort := selection.Wireguard.CustomPort
		if customPort > 0 {
			return customPort
		}
		const defaultPort = 1194
		return defaultPort
	default: // OpenVPN
		customPort := selection.OpenVPN.CustomPort
		if customPort > 0 {
			return customPort
		}
		port = 1194
		if selection.OpenVPN.TCP {
			port = 443
		}
		return port
	}
}
