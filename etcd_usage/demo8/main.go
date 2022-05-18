package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		cli    *clientv3.Client
		err    error
		kv     clientv3.KV
		putOp  clientv3.Op
		opResp clientv3.OpResponse
		getOp  clientv3.Op
	)
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		return
	}

	defer cli.Close()

	kv = clientv3.NewKV(cli)

	putOp = clientv3.OpPut("/cron/jobs/job3", "")
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}

	getOp = clientv3.OpGet("/cron/jobs/job3")
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Revision:", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("value:", string(opResp.Get().Kvs[0].Value))

}
