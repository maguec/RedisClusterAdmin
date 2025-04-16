package commands

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
  "github.com/mpvl/unique"
	"strings"
)

func GetMasterNodes(conf *redis.ClusterOptions) ([]string, error) {
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
  unique.Strings(&nodes)
	return nodes, nil
}

func PrettyPrintSlots(conf *redis.ClusterOptions, nodes []string) string {
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
