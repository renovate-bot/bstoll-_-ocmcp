package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	httpAddr = flag.String("http_address", "http://localhost:8080", "address of MCP server")
	query    = flag.String("query", "", "search query to send to server")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)
	cs, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: *httpAddr}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cs.Close()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name:      "openconfig_yang_model_search",
		Arguments: map[string]any{"query": *query},
	})
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RES: ", string(data))
}
