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
	visitors      map[string]EndpointThrottlingLimiter
	visitorsMutex *sync.Mutex
}

type EndpointThrottlingConfig struct {
	endpoint      string
	method        string
	maxRequests   int
	timePeriodSec int
}

type EndpointThrottlingLimiter struct {
	limiter           *rate.Limiter
	config            EndpointThrottlingConfig
	endingTimeLimiter time.Time
}

func NewThrottlingEngine() ThrottlingEngine {
	return ThrottlingEngine{
		visitors:      map[string]EndpointThrottlingLimiter{},
		visitorsMutex: &sync.Mutex{},
	}
}

func NewEndpointThrottlingLimiter(req http.Request) EndpointThrottlingLimiter {
	var etc EndpointThrottlingConfig

	switch req.URL.Path {
	case platformsRoute:
		etc = createPlatformsThrottlingConfig()
	case leaksRoute:
	default:
		etc = createLeaksEndpointThrottlingConfig()
	}

	return EndpointThrottlingLimiter{
		limiter: rate.NewLimiter(),
		config:  etc,
	}
}

func (te *ThrottlingEngine) CanAllowRequest(req http.Request) bool {
	ipextr := echo.ExtractIPDirect()
	ip := ipextr(&req)
	id := fmt.Sprintf("%s%s", ip, req.Method)

	te.visitorsMutex.Lock()
	defer te.visitorsMutex.Unlock()

	etl, ok := te.visitors[id]

	if !ok {
		etl = NewEndpointThrottlingLimiter(req)
		te.visitors[id] = etl
	}

	return etl.limiter.Allow()
}

func createLeaksEndpointThrottlingConfig() EndpointThrottlingConfig {
	return EndpointThrottlingConfig{
		endpoint:      leaksRoute,
		method:        http.MethodGet,
		maxRequests:   leaksEndpointMaxRequestsPerTimePeriod,
		timePeriodSec: leaksEndpointMaxRequestsTimePeriod,
	}
}

func createPlatformsThrottlingConfig() EndpointThrottlingConfig {
	return EndpointThrottlingConfig{
		endpoint:      platformsRoute,
		method:        http.MethodGet,
		maxRequests:   platformsEndpointMaxRequestsPerTimePeriod,
		timePeriodSec: platformsEndpointMaxRequestsTimePeriod,
	}
}
