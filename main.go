package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	host = "http://srv.msk01.gigacorp.local/_stats"
)

var errorCount = 0

func main() {
	for {
		responseInSplit := requestToServer(host)

		proccessingStatsData(parseAndGetArrayFromResponse(responseInSplit))
		time.Sleep(1 * time.Second)
	}
}

func requestToServer(host string) []string {
	response, error := http.Get(host)

	if error != nil {
		fmt.Println(error)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	if response.StatusCode != http.StatusOK {
		errorCount++
	}
	response.Body.Close()

	return strings.Split(strings.Trim(string(body), "  \n"), ",")
}

func parseAndGetArrayFromResponse(responseInSplit []string) [7]int64 {
	loadAverage, err := strconv.ParseInt(responseInSplit[0], 10, 64)
	if err != nil {
		errorCount++
	}

	memoryValue, err := strconv.ParseInt(responseInSplit[1], 10, 64)
	if err != nil {
		errorCount++
	}

	memoryUsage, err := strconv.ParseInt(responseInSplit[2], 10, 64)
	if err != nil {
		errorCount++
	}
	diskValue, err := strconv.ParseInt(responseInSplit[3], 10, 64)
	if err != nil {
		errorCount++
	}

	diskUsage, err := strconv.ParseInt(responseInSplit[4], 10, 64)
	if err != nil {
		errorCount++
	}

	networkValue, err := strconv.ParseInt(responseInSplit[5], 10, 64)
	if err != nil {
		errorCount++
	}

	networkUsage, err := strconv.ParseInt(responseInSplit[6], 10, 64)
	if err != nil {
		errorCount++
	}

	return [7]int64{loadAverage, memoryValue, memoryUsage, diskValue, diskUsage, networkValue, networkUsage}
}

func proccessingStatsData(responseInInt64 [7]int64) {
	if errorCount >= 3 {
		fmt.Println("Unable to fetch server statistic")

		return
	}

	if responseInInt64[0] > 30 {
		fmt.Println("Load Average is too high:", responseInInt64[0])
	}

	capacityUsageMemory := float64(responseInInt64[2]) / float64(responseInInt64[1])

	if capacityUsageMemory > 0.8 {
		fmt.Printf("Memory usage too high: %d%\n", (int)(capacityUsageMemory*100))
	}
	capacityUsageDisk := float64(responseInInt64[4]) / float64(responseInInt64[3])
	freeDiskSpace := responseInInt64[3] - responseInInt64[4]
	if 0 > freeDiskSpace {
		freeDiskSpace = 0
	}

	if capacityUsageDisk > 0.9 {
		fmt.Printf("Free disk space is too low %d Mb left\n", (int)(freeDiskSpace/1024/1024))
	}

	networkCapacity := float64(responseInInt64[6]) / float64(responseInInt64[5])

	freeNetworkValue := responseInInt64[5] - responseInInt64[6]

	if 0 > freeNetworkValue {
		freeNetworkValue = 0
	}

	if networkCapacity > 0.9 {
		fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", (int)(freeNetworkValue/1000/1000))
	}
	errorCount = 0
}
