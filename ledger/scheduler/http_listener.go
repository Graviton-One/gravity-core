package scheduler

import (
	"net/http"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func startHttpListener() {
	e := echo.New()
	e.POST("/webhooks", HttpUpdateOraclesHandler)
	e.Logger.Debug(e.Start("127.0.0.1:3501"))
}

func HttpUpdateOraclesHandler(c echo.Context) error {
	payload := struct {
		NebulaKey string            `mapstructure:"nebula_key"`
		ChainType account.ChainType `mapstructure:"chain_type"`
		RoundId   int64             `mapstructure:"round_id"`
	}{}
	if err := c.Bind(&payload); err != nil {
		zap.L().Error(err.Error())
		return err
	}

	adaptor, ok := GlobalScheduler.Adaptors[payload.ChainType]
	if !ok {
		zap.L().Debug("adaptor not exists")
	}

	eventPayload := map[string]interface{}{
		"nebula_key": payload.NebulaKey,
		"round_id":   payload.RoundId,
		"sender":     adaptor.PubKey(),
		"is_sender":  true,
		"chain_type": payload.ChainType,
	}
	PublishMessage("ledger.events", SchedulerEvent{
		Name:   "update_oracles",
		Params: eventPayload,
	})
	return c.String(http.StatusOK, "OK")
}
