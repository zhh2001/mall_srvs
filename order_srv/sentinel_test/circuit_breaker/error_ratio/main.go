package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
)

type stateChangeTestListener struct{}

func (listener *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf(
		"rule.steategy: %+v, From %s to Closed, time: %d\n",
		rule.Strategy, prev.String(), util.CurrentTimeMillis(),
	)
}

func (listener *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf(
		"rule.steategy: %+v, From %s to Open, snapshot: %.2f, time: %d\n",
		rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis(),
	)
}

func (listener *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf(
		"rule.steategy: %+v, From %s to Half-Open, time: %d\n",
		rule.Strategy, prev.String(), util.CurrentTimeMillis(),
	)
}

func main() {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:                     "abc",
			Strategy:                     circuitbreaker.ErrorRatio,
			RetryTimeoutMs:               3000,
			MinRequestAmount:             10,
			StatIntervalMs:               5000,
			StatSlidingWindowBucketCount: 10,
			Threshold:                    0.4,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	logging.Info("[CircuitBreaker ErrorRatio] Sentinel Go circuit breaking demo is running. You may see the pass/block metric in the metric log.")
	go func() {
		for {
			e, b := sentinel.Entry("abc")
			if b != nil {
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				if rand.Uint64()%20 > 6 {
					sentinel.TraceError(e, errors.New("biz error"))
				}
				time.Sleep(time.Duration(rand.Uint64()%80+20) * time.Millisecond)
				e.Exit()
			}
		}
	}()
	go func() {
		for {
			e, b := sentinel.Entry("abc")
			if b != nil {
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				time.Sleep(time.Duration(rand.Uint64()%80+40) * time.Millisecond)
				e.Exit()
			}
		}
	}()
	<-ch
}
