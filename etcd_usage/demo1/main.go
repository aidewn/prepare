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

		putResp *clientv3.PutResponse
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
	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job", "321", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
		}
	}
}
