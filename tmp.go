package tmp

import (
	"sync/atomic"
	"time"

	"github.com/golang/rate"
)

const (
	// DEFAULTQPS _
	DEFAULTQPS = 1
	// DEFAULTCON 可以理解为并发，10000基本表示不限制
	DEFAULTCON = 10000
)

var (
	glimiter     atomic.Value
	updateTicker *time.Ticker
)

// GatewayLimiter _
type GatewayLimiter struct {
	Limiter   *rate.Limiter
	LimitTime *LimitTime
}

// LimitTime _
type LimitTime struct {
	StartTime int64
	EndTime   int64
}

// RateLimiterMgr _
type RateLimiterMgr struct {
	actionLimiters map[string]*GatewayLimiter
}

// Allow _
func (r *RateLimiterMgr) Allow(action string) bool {

	if len(action) <= 0 {
		return true
	}

	actionLimiter, ok := r.actionLimiters[action]
	if !ok {
		return true
	}

	if actionLimiter.Limiter == nil {
		return true
	}
	nowUnix := time.Now().Unix()
	// 设置了限制时间，但是不在限制时间内
	if actionLimiter.LimitTime != nil &&
		(actionLimiter.LimitTime.StartTime > nowUnix || actionLimiter.LimitTime.EndTime < nowUnix) {
		return true
	}

	return actionLimiter.Limiter.Allow()

}

// InitRate _
func InitRate(initTicker bool) {
	if initTicker {
		updateTicker = time.NewTicker(1 * time.Minute)
		go crondReload()
	}

	rateCfg, err := GlobalRateConfig()
	if err != nil {
		logs.Errorf("InitRate err %+v", err)
		return
	}

	limiter := &RateLimiterMgr{
		actionLimiters: map[string]*GatewayLimiter{},
	}

	for action, cfg := range rateCfg {
		if cfg.Enabled <= 0 {
			continue
		}

		if cfg.Strategy == nil {
			continue
		}

		aqps := cfg.Strategy.QPS
		if aqps <= 0 {
			aqps = DEFAULTQPS
		}

		acon := cfg.Strategy.Con
		if acon <= 0 {
			acon = DEFAULTCON
		}

		tmpGateLimiter := &GatewayLimiter{
			Limiter: rate.NewLimiter(rate.Limit(aqps), int(acon)),
		}

		if len(cfg.EnableTime) > 0 {
			start, end := GetListTimeUnix(cfg.EnableTime)
			if start == 0 || end == 0 {
				// 如果配置了时间，但是时间配错了，流控不生效
				logs.Errorf("ignore because of err enable time action %s", action)
				continue
			}
			tmpGateLimiter.LimitTime = &LimitTime{
				StartTime: start,
				EndTime:   end,
			}
		}

		logs.Debugf("init %s %+v %+v", action, tmpGateLimiter.LimitTime, tmpGateLimiter.Limiter)

		limiter.actionLimiters[action] = tmpGateLimiter
	}
	glimiter.Store(limiter)

	return
}

func crondReload() {

	defer func() {
		if r := recover(); r != nil {
			logs.Errorf("[Update rate err] recovered in crondReload, %v", r)
		}
	}()

	for {
		now := time.Now()
		InitRate(false)
		cost := time.Now().Sub(now)
		logs.Infof("gateway rate crond Update cost %v", cost)
		<-updateTicker.C
	}
}

// GetGlobalLimiter _
func GetGlobalLimiter() atomic.Value {
	return glimiter
}
