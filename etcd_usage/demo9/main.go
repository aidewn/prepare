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
		leaseId      clientv3.LeaseID
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		keepResp     *clientv3.LeaseKeepAliveResponse
		cancelFunc   context.CancelFunc
		txn          clientv3.Txn
		txnResp      *clientv3.TxnResponse
		ctx          context.Context
	)
	// 客户端配置
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		return
	}

	defer cli.Close()

	// 1, 上锁 (创建租约, 自动续租, 拿着租约去抢占一个key)
	lease = clientv3.NewLease(cli)

	// 申请一个5秒的租约
	leaseGrantResp, _ = lease.Grant(context.TODO(), 5)

	// 拿到租约的ID
	leaseId = leaseGrantResp.ID

	// 准备一个用于取消自动续租的context
	ctx, cancelFunc = context.WithCancel(context.TODO())

	// 确保函数退出后, 自动续租会停止
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	// 5秒后会取消自动续租
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
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

	txn = kv.Txn(context.TODO())

	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job3"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job3", "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job3"))

	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	if !txnResp.Succeeded {
		fmt.Println("锁被占用：", txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		return
	}

	fmt.Println("处理事务")
	time.Sleep(time.Second * 5)

}
