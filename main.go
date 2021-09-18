package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	datadog "github.com/DataDog/datadog-api-client-go/api/v2/datadog"
)

func main() {
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)

	var (
		next *string
		i    int
	)
	for {
		fmt.Printf("starting page %d\n", i)
		resp, r, err := apiClient.LogsApi.ListLogsGet(ctx, getParams(next))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `LogsApi.ListLogsGet`: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
			return
		}

		f := getExportFile(i)
		err = json.NewEncoder(f).Encode(resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing to file`: %v\n", err)
			return
		}
		err = f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error closing file`: %v\n", err)
			return
		}
		if resp.Meta == nil || resp.Meta.Page == nil || resp.Meta.Page.After == nil {
			break
		}
		next = resp.Meta.Page.After
		fmt.Printf("finshed page %d\n", i)
		i++
	}
}

func getParams(next *string) datadog.ListLogsGetOptionalParameters {
	filterQuery := `service:aaas-ingestion @msg:"error publishing protobuf payload; dropping"` // string | Search query following logs syntax. (optional)
	filterFrom := time.Date(2021, time.September, 14, 0, 0, 0, 0, time.UTC)                    // time.Time | Minimum timestamp for requested logs. (optional)
	filterTo := time.Date(2021, time.September, 17, 0, 0, 0, 0, time.UTC)                      // time.Time | Maximum timestamp for requested logs. (optional)
	pageLimit := int32(1000)                                                                   // int32 | Maximum number of logs in the response. (optional) (default to 10)
	optionalParams := datadog.ListLogsGetOptionalParameters{
		FilterQuery: &filterQuery,
		FilterFrom:  &filterFrom,
		FilterTo:    &filterTo,
		PageCursor:  next,
		PageLimit:   &pageLimit,
	}

	return optionalParams
}

func getExportFile(i int) io.WriteCloser {
	fileName := fmt.Sprintf("export-%d.json", i)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("failed to open export file", err)
	}
	return file
}
