package test

import (
	"testing"

	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
	"github.com/degatedev/degate-sdk-golang/log"
)

func TestNewListenerKey(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.NewListenKey()
	if err != nil {
		panic(err)
	}
	if r.Success() {
		log.Print("NewListenKey success %v", lib.String(r.Data))
	} else {
		log.Print("NewListenKey fail %v", lib.String(r))
	}
}

func TestDelayListenerKey(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.ReNewListenKey(&model.ListenKeyParam{
		ListenKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoyLCJleHAiOjEwODYyNTIyOTk3LCJvcmlnX2F0IjoxNjM5MTUwOTYwLCJzb3VyY2VfaXAiOiIxMDEuODcuNjQuMjQ2In0.8hShzearBLi5qLQEjxCdxtFWGSQSUd3DaU77rs_HCe4",
	})
	if err != nil {
		panic(err)
	}
	if r.Success() {
		log.Print("DelayListenerKey success")
	} else {
		log.Print("DelayListenerKey fail %v", lib.String(r))
	}
}

func TestDeleteListenerKey(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.DeleteListenKey(&model.ListenKeyParam{
		ListenKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoyLCJleHAiOjEwODYyNTIyOTk3LCJvcmlnX2F0IjoxNjM5MTUwOTYwLCJzb3VyY2VfaXAiOiIxMDEuODcuNjQuMjQ2In0.8hShzearBLi5qLQEjxCdxtFWGSQSUd3DaU77rs_HCe4",
	})
	if err != nil {
		panic(err)
	}
	if r.Success() {
		log.Print("DeleteListenKey success")
	} else {
		log.Print("DeleteListenKey fail %v", lib.String(r))
	}
}
