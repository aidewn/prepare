package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"time"
)

func main() {
	var (
		cli *clientv3.Client
		err error
		kv  clientv3.KV

		getResp            *clientv3.GetResponse
		watchStartRevision int64
		watcher            clientv3.Watcher
		wacherRespChan     <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
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

	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/job3", "i am job3")
			kv.Delete(context.TODO(), "/cron/jobs/job3")
			time.Sleep(time.Second * 1)
		}
	}()

	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job3"); err != nil {
		fmt.Println(err)
		return
	}

	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值", string(getResp.Kvs[0].Value))
	}

	watchStartRevision = getResp.Header.Revision + 1

	watcher = clientv3.NewWatcher(cli)

	fmt.Println("从该版本向后监听： ", watchStartRevision)

	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(time.Second*5, func() {
		cancelFunc()
	})
	wacherRespChan = watcher.Watch(ctx, "/cron/jobs/job3", clientv3.WithRev(watchStartRevision))

	for watchResp = range wacherRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}

}
