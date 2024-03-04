package main

import (
	"context"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/go-redis/redis/v9"
	"os"
)

var args struct {
	ClusterServer string   `help:"Cluster Server Host" default:"localhost" arg:"--server, -s, env:CLUSTER_SERVER"`
	ClusterPort   string   `help:"Cluster Server Port" default:"6379" arg:"--port, -p, env:CLUSTER_PORT"`
	Verbose       bool     `help:"Verbose" arg:"--verbose, -v"`
	Command       []string `help:"Command" arg:"positional" required:"true"`
}

func getMasterNodes(conf *redis.ClusterOptions) ([]string, error) {
	var ctx = context.Background()
	nodes := []string{}
	client := redis.NewClusterClient(conf)
	slots, err := client.ClusterSlots(ctx).Result()
	if err != nil {
		return nil, err
	}
	for _, s := range slots {
		nodes = append(nodes, s.Nodes[0].Addr)
	}
	return nodes, nil
}

func main() {
	arg.MustParse(&args)
	conf := &redis.ClusterOptions{
		Addrs: []string{fmt.Sprintf("%s:%s", args.ClusterServer, args.ClusterPort)},
	}
	nodes, err := getMasterNodes(conf)
	if err != nil {
		panic(err)
	}
	for _, n := range nodes {
		if args.Verbose {
			fmt.Fprintf(os.Stderr, "# %s\n", n)
		}
		rdb := redis.NewClient(&redis.Options{Addr: n})
		cmd := make([]interface{}, len(args.Command))
		for i, v := range args.Command {
			cmd[i] = v
		}
		res, err := rdb.Do(context.Background(), cmd...).Result()
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(os.Stdout, res)
	}

}
