package main

import (
	"context"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/go-redis/redis/v9"
	"github.com/maguec/RedisClusterAdmin/commands"
	"os"
	"strings"
)

var args struct {
	ClusterServer string   `help:"Cluster Server Host" default:"localhost" arg:"--server, -s, env:CLUSTER_SERVER"`
	ClusterPort   string   `help:"Cluster Server Port" default:"6379" arg:"--port, -p, env:CLUSTER_PORT"`
	Verbose       bool     `help:"Verbose" arg:"--verbose, -v"`
	Keyspace      bool     `help:"Get the cluster Keyspace stats" arg:"--keyspace, -k"`
	Summit        bool     `help:"Sum the stat returned" arg:"--sum, -m"`
	Command       []string `help:"Command" arg:"positional" required:"true"`
}

func main() {
	arg.MustParse(&args)
	conf := &redis.ClusterOptions{
		Addrs: []string{fmt.Sprintf("%s:%s", args.ClusterServer, args.ClusterPort)},
	}
	nodes, err := commands.GetMasterNodes(conf)
	if err != nil {
		panic(err)
	}

	cmd := make([]interface{}, len(args.Command))
	for i, v := range args.Command {
		cmd[i] = v
	}

	// if the command is cluster slots intercept it and return an prettier output
	if strings.ToLower(cmd[0].(string)) == "cluster" && strings.ToLower(cmd[1].(string)) == "slots" {
		fmt.Println(commands.PrettyPrintSlots(conf, nodes))
		return
	}

	if strings.ToLower(cmd[0].(string)) == "info" && args.Keyspace {
		err := commands.PrintKeyspace(conf, nodes, args.Summit)
		if err != nil {
			panic(err)
		}
		return
	}

	for _, n := range nodes {
		if args.Verbose {
			fmt.Fprintf(os.Stderr, "# %s\n", n)
		}
		rdb := redis.NewClient(&redis.Options{Addr: n})
		res, err := rdb.Do(context.Background(), cmd...).Result()
		if err != nil {
			panic(err)
		}
		if _, ok := res.(string); ok {
			fmt.Fprintf(os.Stdout, "%s\n", res)
		} else if _, ok := res.([]interface{}); ok {
			for _, v := range res.([]interface{}) {
				fmt.Fprintf(os.Stdout, "%+v\n", v)
			}
		} else {
			fmt.Fprintf(os.Stdout, "%+v\n", res)
		}
	}

}
