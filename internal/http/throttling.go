package http

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/palavrapasse/santos/internal/logging"
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

const throttlingEngineVisitorsCleanUpMinutes = 1

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
	Limiter  *rate.Limiter
	Config   EndpointThrottlingConfig
	LastSeen time.Time
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
		Limiter:  rate.NewLimiter(r, b),
		Config:   etc,
		LastSeen: time.Now(),
	}
}

func StartThrottlingEngineCleanUp(te *ThrottlingEngine) {
	go te.cleanupThrottlingEngine()
}

func (te *ThrottlingEngine) CanAllowRequest(req http.Request) bool {
	id := requestId(req)

	etl, ok := te.Visitors[id]

	if !ok {
		te.VisitorsMutex.Lock()
		defer te.VisitorsMutex.Unlock()

		etl = NewEndpointThrottlingLimiter(req)
	}

	etl.LastSeen = time.Now()

	te.Visitors[id] = etl

	return etl.Limiter.Allow()
}

func (te *ThrottlingEngine) cleanupThrottlingEngine() {
	for {
		time.Sleep(time.Minute)

		te.VisitorsMutex.Lock()

		for ip, v := range te.Visitors {
			if time.Since(v.LastSeen) > throttlingEngineVisitorsCleanUpMinutes*time.Minute {
				delete(te.Visitors, ip)
			}
		}

		te.VisitorsMutex.Unlock()
	}
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
