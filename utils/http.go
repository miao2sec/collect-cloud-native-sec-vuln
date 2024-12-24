package utils

import (
	"crypto/rand"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
	"math"
	"math/big"
	"time"
)

const (
	Get  = "GET"
	Post = "POST"
)

type HTTP struct {
	URL    string
	Method string
	Retry  int
	Body   interface{}
	Params map[string]string
}

func (h *HTTP) Fetch(res any) (err error) {
	for i := 0; i <= h.Retry; i++ {
		// 若首次访问失败，则计算出一个随机数作为间隔，然后再次访问
		if i > 0 {
			//wait = i^2+[0-9中的一个随机数]+[0,MaxInt64)的整数
			wait := math.Pow(float64(i), 2) + float64(RandInt()%10)
			log.Debug().Msgf("retry after %f seconds", wait)
			time.Sleep(time.Duration(wait) * time.Second)
		}

		r, err := h.fetch()
		if err == nil {
			if err = json.Unmarshal(r, res); err != nil {
				return err
			}
			return nil
		}

	}
	return xerrors.Errorf("failed to fetch %s: %w", h.URL, err)
}

func (h *HTTP) fetch() ([]byte, error) {
	var (
		client = resty.New()
		req    = client.R()
	)

	req.Method = h.Method
	req.URL = h.URL
	if h.Method == Post {
		req.SetBody(h.Body)
	}
	if h.Method == Get {
		req.SetQueryParams(h.Params)
	}

	resp, err := req.Send()
	if err != nil {
		return nil, xerrors.Errorf("failed to send %s:%w", req.URL, err)
	}
	if resp.StatusCode() != 200 {
		return nil, xerrors.Errorf("[%v]request %s:%v", resp.StatusCode(), req.URL, string(resp.Body()))
	}

	return resp.Body(), nil
}

// RandInt 返回一个[0,MaxInt64)的一个随机整数
func RandInt() int {
	seed, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	return int(seed.Int64())
}
