package main

import (
	"context"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/go-redis/redis/v9"
	"os"
	"strings"
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

func prettyprintSlots(conf *redis.ClusterOptions, nodes []string) string {
	var res strings.Builder
	client := redis.NewClient(&redis.Options{Addr: nodes[0]})
	slots, err := client.ClusterSlots(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	for _, slot := range slots {
		res.WriteString(
			fmt.Sprintf(
				"Slot: %5d - %5d Primary: %+v Secondaries: ",
				slot.Start, slot.End, slot.Nodes[0].Addr))
		for _, n := range slot.Nodes[1:] {
			res.WriteString(fmt.Sprintf("%+v ", n.Addr))
		}
		res.WriteString("\n")
	}

	return res.String()
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

	cmd := make([]interface{}, len(args.Command))
	for i, v := range args.Command {
		cmd[i] = v
	}

	// if the command is cluster slots intercept it and return an prettier output
	if strings.ToLower(cmd[0].(string)) == "cluster" && strings.ToLower(cmd[1].(string)) == "slots" {
		fmt.Println(prettyprintSlots(conf, nodes))
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
