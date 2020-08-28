package helpers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wavesplatform/gowaves/pkg/crypto"

	"github.com/wavesplatform/gowaves/pkg/client"
)

const (
	GetStateByAddressPath = "addresses/data"

	TxWaitCount    = 10
	BlockWaitCount = 30
)

type ClientHelper struct {
	client *client.Client
}

func NewClientHelper(client *client.Client) ClientHelper {
	return ClientHelper{client: client}
}

func (helper *ClientHelper) GetStateByAddressAndKey(address string, key string, ctx context.Context) (*State, *client.Response, error) {
	url := fmt.Sprintf("%s/%s/%s?key=%s", helper.client.GetOptions().BaseUrl, GetStateByAddressPath, address, key)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []State
	response, err := helper.client.Do(ctx, req, &out)
	if err != nil {
		return nil, response, err
	}

	return &out[0], response, nil
}

func (helper *ClientHelper) WaitTx(id string, ctx context.Context) <-chan error {
	out := make(chan error)
	idDig := crypto.MustDigestFromBase58(id)
	go func() {
		defer close(out)
		for i := 0; i <= TxWaitCount; i++ {
			_, res, err := helper.client.Transactions.UnconfirmedInfo(ctx, idDig)
			if err != nil && res == nil {
				out <- err
				break
			}

			if res.StatusCode != http.StatusOK {
				_, res, err := helper.client.Transactions.Info(ctx, idDig)
				if err != nil && res == nil {
					out <- err
					break
				}

				if res.StatusCode != http.StatusOK {
					out <- errors.New("tx not found")
					break
				} else {
					break
				}
			}

			if TxWaitCount == i {
				out <- errors.New("tx not found")
				break
			}
			time.Sleep(time.Second)
		}
	}()
	return out
}
func (helper *ClientHelper) WaitByHeight(height uint64, ctx context.Context) <-chan error {
	out := make(chan error)
	go func() {
		defer close(out)
		for i := 0; i <= BlockWaitCount; i++ {
			currentHeight, _, err := helper.client.Blocks.Height(ctx)
			if err != nil {
				out <- err
				break
			}

			if currentHeight.Height >= height {
				break
			}

			if BlockWaitCount == i {
				out <- errors.New("block not found")
				break
			}

			time.Sleep(10 * time.Second)
		}
	}()
	return out
}
