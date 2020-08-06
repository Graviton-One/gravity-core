package extractors

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type BinanceExtractor struct{}

func (p *BinanceExtractor) Extract() (interface{}, error) {
	priceWavesUsdt, err := p.priceNowByPair("WAVESUSDT")
	if err != nil {
		return nil, err
	}

	priceWavesBtc, err := p.priceNowByPair("WAVESBTC")
	if err != nil {
		return nil, err
	}

	priceBtcUsdt, err := p.priceNowByPair("BTCUSDT")
	if err != nil {
		return nil, err
	}

	price := int64(((priceWavesUsdt + (priceWavesBtc * priceBtcUsdt)) / 2) * 100)

	return price, nil
}

func (p *BinanceExtractor) Aggregate(values []interface{}) (interface{}, error) {
	var intValues []int64
	for _, b := range values {
		intValues = append(intValues, b.(int64))
	}

	var result int64
	for _, v := range intValues {
		result += v
	}
	result = result / int64(len(intValues))

	return result, nil
}

func (p *BinanceExtractor) priceNowByPair(pair string) (float64, error) {
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
