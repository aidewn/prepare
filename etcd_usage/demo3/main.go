package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"

	"time"
)

func main() {
	var (
		cli *clientv3.Client
		err error
		kv  clientv3.KV

		getResp *clientv3.GetResponse
	)
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.0.16:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		return
	}

	defer cli.Close()

	kv = clientv3.NewKV(cli)

	kv.Put(context.TODO(), "/cron/jobs/job2", "123")
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Kvs)
	}
}
