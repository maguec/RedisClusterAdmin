package commands

// This is heavily ripped off from
// https://github.com/oliver006/redis_exporter/blob/d9ccc72f2061857f41585768d2318ae5b765cfdf/exporter/info.go

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"os"
)

func PrintKeyspace(conf *redis.ClusterOptions, nodes []string, summit bool) error {
  var keys []KeyspaceInfo
	r, err := getInfo(nodes, false, "keyspace", summit)
	if err != nil {
		return err
	}
	for _, v := range r {
    k, err := extractInfoMetrics(v, "keyspace")
    if err != nil {
      panic(err)
    }

    keys = append(keys, k.([]KeyspaceInfo)...)
  }
  if summit {
    var t int64
    for _, k := range keys {
      t += k.Keys
    }
    fmt.Printf("%d\n", t)
  } else {
    for _, k := range keys {
      fmt.Printf("%d\n", k.Keys)
    }
  }
	return nil
}

func getInfo(nodes []string, verbose bool, filter string, summit bool) ([]string, error) {
	var res []string
	for _, n := range nodes {
		if verbose {
			fmt.Fprintf(os.Stderr, "# %s\n", n)
		}
		rdb := redis.NewClient(&redis.Options{Addr: n})
		r, err := rdb.Info(context.Background(), filter).Result()
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func extractInfoMetrics(info, filter string) (interface{}, error) {
	var data interface{}
  var err error
	switch filter {
	case "keyspace":
		data, err = parseKeyspaceDBLines(info)
		if err != nil {
			return nil, err
		}

	default:
	}
	return data, nil
}
