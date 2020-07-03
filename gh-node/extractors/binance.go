package extractors

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type BinanceExtractor struct{}

func (p *BinanceExtractor) PriceNow() (float64, error) {
	priceWavesUsdt, err := p.PriceNowByPair("WAVESUSDT")
	if err != nil {
		return 0, err
	}

	priceWavesBtc, err := p.PriceNowByPair("WAVESBTC")
	if err != nil {
		return 0, err
	}

	priceBtcUsdt, err := p.PriceNowByPair("BTCUSDT")
	if err != nil {
		return 0, err
	}

	price := (priceWavesUsdt + (priceWavesBtc * priceBtcUsdt)) / 2

	return price, nil
}

func (p *BinanceExtractor) PriceNowByPair(pair string) (float64, error) {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=" + pair)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	jsonResponse := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(jsonResponse["price"].(string), 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}
