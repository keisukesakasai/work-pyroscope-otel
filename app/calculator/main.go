package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"

	_ "time/tzdata"

	"github.com/gin-gonic/gin"
	"github.com/grafana/pyroscope-go"
)

func main() {
	// 環境変数
	appVersion := os.Getenv("APP_VERSION")

	// profiling 設定
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
	pyroscope.Start(pyroscope.Config{
		ApplicationName: "calculator",
		ServerAddress:   "http://pyroscope.pyroscope.svc.cluster.local:4040",
		Logger:          pyroscope.StandardLogger,

		// you can provide static tags via a map:
		Tags: map[string]string{
			"hostname": "calculator",
			"version":  appVersion,
		},

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
	})

	// Http server
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		// リクエストのボディを取得します
		_, err := io.ReadAll(c.Request.Body)
		if err != nil {
			http.Error(c.Writer, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// ロジック
		count := calcTargetLogic(appVersion)
		fmt.Println("count: ", count)
	})

	r.Run(":8080")
}

func calcTargetLogic(appVersion string) (total int) {
	dummyData, err := read("./data/input.txt")
	if err != nil {
		fmt.Println(err.Error())
	}

	total = count(dummyData, appVersion)

	return total
}

func read(filename string) ([]int, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var dummyData []int
	for _, char := range string(content) {
		value, err := strconv.Atoi(string(char))
		if err != nil {
			return nil, fmt.Errorf("error converting character to int: %v", err)
		}
		if value != 0 && value != 1 {
			return nil, fmt.Errorf("invalid value in file: %d", value)
		}
		dummyData = append(dummyData, value)
	}

	return dummyData, nil
}

func count(dummyData []int, appVersion string) (total int) {
	switch appVersion {
	case "v1.0.0":
		n := len(dummyData)
		for i := 0; i < n; i++ {
			for j := 0; j < n-i-1; j++ {
				if dummyData[j] > dummyData[j+1] {
					dummyData[j], dummyData[j+1] = dummyData[j+1], dummyData[j]
				}
			}
		}
		index := sort.SearchInts(dummyData, 1)
		fmt.Println("index: ", index)

		return len(dummyData) - index

	case "v2.0.0":
		sort.Ints(dummyData)
		index := sort.SearchInts(dummyData, 1)
		fmt.Println("index: ", index)

		return len(dummyData) - index

	default:
		return 0
	}
}
