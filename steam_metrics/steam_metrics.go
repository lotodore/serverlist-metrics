// Example lambda function returning prometheus metrics for steam game servers.
package main

import (
	"errors"
	"net/http"
	"os"
	"steam_metrics/webrequest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalNumPlayers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "steam",
			Name:      "total_num_players",
			Help:      "Total number of players on all servers",
		})
	totalNumServers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "steam",
			Name:      "total_num_servers",
			Help:      "Total number of servers",
		})
	numNonEmptyServers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "steam",
			Name:      "num_non_empty_servers",
			Help:      "Number of non empty servers",
		})
)

// Implementation of http.ResponseWriter to map Lambda to Prometheus http
type LambdaResponse struct {
	Body       string
	StatusCode int
}

// The incoming http Header is not used by Prometheus.
func (r LambdaResponse) Header() http.Header {
	return http.Header{}
}

// Append Prometheus data to Lambda response variable.
func (r *LambdaResponse) Write(data []byte) (int, error) {
	r.Body += string(data)
	return len(data), nil
}

// Set lambda status code from Prometheus.
func (r *LambdaResponse) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
}

// Implementation of Lambda function handler. Use Prometheus http handler by mapping callbacks to lambda.
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	steamAppId := os.Getenv("STEAM_APP_ID")
	if steamAppId == "" {
		return events.APIGatewayProxyResponse{}, errors.New("steam app id not set")
	}

	// Example: Filter only by app id.
	servers, err := webrequest.RequestGameServers(webrequest.GetFilterForAppId(steamAppId))

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Do not use default prometheus metrics, they are pretty useless within a lambda enviroment.
	// Create a new prometheus registry instead.
	registry := prometheus.NewRegistry()
	registry.MustRegister(totalNumPlayers)
	registry.MustRegister(totalNumServers)
	registry.MustRegister(numNonEmptyServers)

	// Example metrics for game servers: Total servers, non-empty servers and total players.
	tmpTotalPlayers := 0
	tmpTotalServers := 0
	tmpNumNonEmptyServers := 0
	for _, s := range servers {
		tmpTotalPlayers += s.Players
		if s.Players > 0 {
			tmpNumNonEmptyServers++
		}
		tmpTotalServers++
	}
	totalNumPlayers.Set(float64(tmpTotalPlayers))
	totalNumServers.Set(float64(tmpTotalServers))
	numNonEmptyServers.Set(float64(tmpNumNonEmptyServers))

	// Map prometheus http handler to Lambda.
	var httpRequest http.Request
	httpRequest.Header = request.MultiValueHeaders
	var lambdaResponse LambdaResponse

	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(&lambdaResponse, &httpRequest)

	// Return Prometheus handler data as Lambda result.
	return events.APIGatewayProxyResponse{
		Body:       lambdaResponse.Body,
		StatusCode: lambdaResponse.StatusCode,
	}, nil
}

func main() {
	lambda.Start(handler)
}
