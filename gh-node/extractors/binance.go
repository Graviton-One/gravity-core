package extractors

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type BinanceExtractor struct{}

func (p *BinanceExtractor) GetData() ([]byte, error) {
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

	price := (priceWavesUsdt + (priceWavesBtc * priceBtcUsdt)) / 2

	var b []byte
	binary.BigEndian.PutUint64(b, uint64(price))

	return b, nil
}

func (p *BinanceExtractor) Aggregate(values [][]byte) ([]byte, error) {
	var intValues []uint64
	for _, b := range values {
		intValues = append(intValues, binary.BigEndian.Uint64(b))
	}

	var result uint64
	for _, v := range intValues {
		result += v
	}
	result = result / uint64(len(intValues))

	var b []byte
	binary.BigEndian.PutUint64(b, uint64(price))

	return b, nil //TODO invalid convert to byte (contracts)
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
