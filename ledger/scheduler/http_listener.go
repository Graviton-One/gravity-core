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
		NebulaKey string            `json:"nebula_key"`
		ChainType account.ChainType `json:"chain_type"`
		//RoundId   int64             `json:"round_id"`
	}{}
	if err := c.Bind(&payload); err != nil {
		zap.L().Error(err.Error())
		return err
	}
	zap.L().Sugar().Debug("payload ", payload)

	ManualUpdate.Active = true
	ManualUpdate.UpdateQueue = append(ManualUpdate.UpdateQueue, NebulaToUpdate{
		Id:        payload.NebulaKey,
		ChainType: payload.ChainType,
	})
	zap.L().Sugar().Debugf("Added manual update nebula: ", ManualUpdate)
	return c.String(http.StatusOK, "OK")
}
