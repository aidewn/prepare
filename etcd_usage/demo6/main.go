package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"

	"time"
)

func main() {
	var (
		cli            *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse

		kv           clientv3.KV
		putResp      *clientv3.PutResponse
		leaseId      clientv3.LeaseID
		getResp      *clientv3.GetResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		keepResp     *clientv3.LeaseKeepAliveResponse
	)
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		return
	}

	defer cli.Close()

	lease = clientv3.NewLease(cli)

	leaseGrantResp, _ = lease.Grant(context.TODO(), 10)

	leaseId = leaseGrantResp.ID

	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约失效了")
					goto END
				} else {
					fmt.Println("收到自动续租请求")
				}
			}
		}
	END:
	}()

	kv = clientv3.NewKV(cli)

	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(putResp.Header.Revision)
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}

		if len(getResp.Kvs) == 0 {
			fmt.Println("过期")
			break
		}

		fmt.Println("还没过期")
		time.Sleep(time.Second * 2)

	}

}
