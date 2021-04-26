package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	node1, address := startNode1()
	node2 := startNode2(address)

	<-ch
	//停止节点
	fmt.Println("shut the node down")
	if err := node1.Close(); err != nil {
		panic(err)
	}
	if err := node2.Close(); err != nil {
		panic(err)
	}

}

func startNode1() (host.Host, string) {
	ctx := context.Background()

	//创建一个节点
	node, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/2000"),
	)
	if err != nil {
		panic(err)
	}

	//配置ping协议消息的处理器
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	//打印节点监听的地址
	fmt.Println("Listen addresses:", node.Addrs())

	//打印当前节点地址
	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrList, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	fmt.Println("libp2p node address:", addrList)
	return node, addrList[0].String()
}

func startNode2(address string) host.Host {
	ctx := context.Background()

	//创建一个节点
	node, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/3000"),
	)
	if err != nil {
		panic(err)
	}

	//配置ping协议消息的处理器
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	//打印节点监听的地址
	fmt.Println("Listen addresses:", node.Addrs())

	//打印当前节点地址
	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrList, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	fmt.Println("libp2p node address:", addrList)

	addr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		panic(err)
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	if err := node.Connect(ctx, *peer); err != nil {
		panic(err)
	}
	fmt.Println("sending ping messages to", addr)
	ch := pingService.Ping(ctx, peer.ID)
	for {
		res := <-ch
		fmt.Println("got ping response!", "RTT:", res.RTT)
	}
	return node
}
