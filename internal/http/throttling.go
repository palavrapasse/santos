package http

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

const (
	leaksEndpointMaxRequestsPerTimePeriod     = 3
	platformsEndpointMaxRequestsPerTimePeriod = 3
)

const (
	leaksEndpointMaxRequestsTimePeriod     = 5
	platformsEndpointMaxRequestsTimePeriod = 5
)

type ThrottlingEngine struct {
	Visitors      map[string]EndpointThrottlingLimiter
	VisitorsMutex *sync.Mutex
}

type EndpointThrottlingConfig struct {
	Endpoint      string
	Method        string
	MaxRequests   int
	TimePeriodSec int
}

type EndpointThrottlingLimiter struct {
	Limiter *rate.Limiter
	Config  EndpointThrottlingConfig
}

func NewThrottlingEngine() ThrottlingEngine {
	return ThrottlingEngine{
		Visitors:      map[string]EndpointThrottlingLimiter{},
		VisitorsMutex: &sync.Mutex{},
	}
}

func NewEndpointThrottlingLimiter(req http.Request) EndpointThrottlingLimiter {
	etc := requestThrottlingConfig(req)

	r := rate.Every(time.Second * time.Duration(etc.TimePeriodSec))
	b := etc.MaxRequests

	return EndpointThrottlingLimiter{
		Limiter: rate.NewLimiter(r, b),
		Config:  etc,
	}
}

func (te *ThrottlingEngine) CanAllowRequest(req http.Request) bool {
	id := requestId(req)

	etl, ok := te.Visitors[id]

	if !ok {
		te.VisitorsMutex.Lock()
		defer te.VisitorsMutex.Unlock()

		etl = NewEndpointThrottlingLimiter(req)
		te.Visitors[id] = etl
	}

	return etl.Limiter.Allow()
}

func requestId(req http.Request) string {
	ipextr := echo.ExtractIPDirect()
	ip := ipextr(&req)

	id := fmt.Sprintf("%s%s", ip, req.Method)

	return id
}

func requestThrottlingConfig(req http.Request) EndpointThrottlingConfig {
	var etc EndpointThrottlingConfig

	switch req.URL.Path {
	case platformsRoute:
		etc = createPlatformsThrottlingConfig()
	case leaksRoute:
		etc = createLeaksEndpointThrottlingConfig()
	}

	return etc
}

func createLeaksEndpointThrottlingConfig() EndpointThrottlingConfig {
	return EndpointThrottlingConfig{
		Endpoint:      leaksRoute,
		Method:        http.MethodGet,
		MaxRequests:   leaksEndpointMaxRequestsPerTimePeriod,
		TimePeriodSec: leaksEndpointMaxRequestsTimePeriod,
	}
}

func createPlatformsThrottlingConfig() EndpointThrottlingConfig {
	return EndpointThrottlingConfig{
		Endpoint:      platformsRoute,
		Method:        http.MethodGet,
		MaxRequests:   platformsEndpointMaxRequestsPerTimePeriod,
		TimePeriodSec: platformsEndpointMaxRequestsTimePeriod,
	}
}
