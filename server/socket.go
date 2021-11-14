// nolint:errcheck
package server

import (
	"context"
	"net/http"

	"github.com/blewater/zh/log"
	"github.com/blewater/zh/types"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	SubReqMsgType         = "subscribe"
	SubAckMsgType         = "subscriptions"
	MatchesChannelMsgType = "matches"
	MatchMsgType          = "match"
	MatchLastMsgType      = "last_match"
	ErrorMsgType          = "error"
)

func Connect(ctx context.Context, socketAddr string) (*websocket.Conn, error) {
	logger := log.FromContext(ctx)

	logger.Debug("connecting", zap.String("host", socketAddr))

	conn, resp, err := websocket.DefaultDialer.Dial(socketAddr, nil)
	if err != nil {
		logger.Error(
			"attempting to connect erred:",
			zap.String("url", socketAddr),
			zap.Error(err),
		)
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		logger.Error("new connection failed to switch:", zap.String("resp", resp.Status))
		return nil, err
	}

	return conn, err
}

func Subscribe(ctx context.Context, conn *websocket.Conn, productIDs []string) error {
	logger := log.FromContext(ctx)

	err := conn.WriteJSON(
		&types.SubReq{
			Type:       SubReqMsgType,
			ProductIds: productIDs,
			Channels:   []string{MatchesChannelMsgType},
		},
	)
	logger.Error("Sending a subscribe msg erred", zap.Error(err))

	return err
}
