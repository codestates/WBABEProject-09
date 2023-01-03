package main

import (
	conf "WBABEProject-09/config"
	ctl "WBABEProject-09/controller"
	log "WBABEProject-09/logger"
	md "WBABEProject-09/model"
	rt "WBABEProject-09/router"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	var configFlag = flag.String("config", "./config/config.toml", "toml file to use for configuration")
	flag.Parse()

	// controller와 model이 한개의 파일로 구성됨, order, menu, 주문자, 피주문자 등으로 나눠서 관리가 필요함 - TODO -
	//model 모듈 선언
	if cf, err := conf.NewConfig(*configFlag); err != nil { // config 모듈 설정
		fmt.Printf("init config failed, err:%v\n", err)
		return
	} else if err := log.InitLogger(cf); err != nil { // logger 모듈 설정
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	} else if mod, err := md.NewModel(cf); err != nil { // model 모듈 설정
		fmt.Printf("NewModel Error: %v\n", err)
	} else if controller, err := ctl.NewCTL(mod); err != nil { //controller 모듈 설정
		fmt.Printf("NewCTL Error: %v\n", err)
	} else if rt, err := rt.NewRouter(controller); err != nil { //router 모듈 설정
		fmt.Printf("NewRouter Error: %v\n", err)
	} else {
		mapi := &http.Server{
			Addr:           cf.Server.Port,
			Handler:        rt.Idx(),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		g.Go(func() error {
			return mapi.ListenAndServe()
		})

		stopSig := make(chan os.Signal)
		signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)
		<-stopSig

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := mapi.Shutdown(ctx); err != nil {
			fmt.Println("Server Shutdown Error:", err)
		}

		select {
		case <-ctx.Done():
			fmt.Println("context done.")
		}
		fmt.Println("Server stop")
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
