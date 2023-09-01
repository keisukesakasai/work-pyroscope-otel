package main

import (
	"io"
	"math/rand"
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
	numRandomDataStr := os.Getenv("NUM_RUNDOM_DATA")
	numRandomData, _ := strconv.Atoi(numRandomDataStr)
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

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
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
		calcTargetLogic(numRandomData, appVersion)
	})

	r.Run(":8080")
}

func calcTargetLogic(numRandomData int, appVersion string) (total int) {
	dummyData := create(numRandomData)
	total = count(dummyData, appVersion)

	return total
}

func create(num int) []int {
	slice := make([]int, num)
	for i := range slice {
		slice[i] = rand.Intn(2)
	}
	return slice
}

func count(dummyData []int, appVersion string) (total int) {
	switch appVersion {
	case "v1.0.0":
		sort.Ints(dummyData)
		index := sort.SearchInts(dummyData, 1)
		return len(dummyData) - index

	case "v2.0.0":
		total := 0
		for _, value := range dummyData {
			if value == 1 {
				total++
			}
		}
		return total

	default:
		return 0
	}
}
