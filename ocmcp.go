// ocmcp is a Model Context Protocol server to allow AI models to resolve
// authoritative details on OpenConfig related components.
package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/bstoll/ocmcp/internal/yang"
	log "github.com/golang/glog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	httpAddr = flag.String("http_address", "", "If set, listen on this address for HTTP requests instead of stdin/stdout.")
	yangDirs = flag.String("yang_dirs", "public/release/models,public/third_party/ietf", "Comma separated list of YANG directories to recursively load.")
)

func main() {
	flag.Parse()

	server := mcp.NewServer(&mcp.Implementation{Name: "ocmcp", Version: "v0.0.0"}, nil)
	if (*yangDirs) == "" {
		log.Errorf("No YANG dirs provided.")
		os.Exit(1)
	}
	yangIndex := yang.YangModelIndex{
		SrcDirs: strings.Split(*yangDirs, ","),
	}
	err := yangIndex.RegisterTools(server)
	if err != nil {
		log.Errorf("Error registering yang index: %v", err)
		os.Exit(1)
	}

	if *httpAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Infof("Listening on address %q", *httpAddr)
		if err := http.ListenAndServe(*httpAddr, handler); err != nil {
			log.Errorf("Server failed: %v", err)
			os.Exit(1)
		}
	} else {
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Errorf("Server failed: %v", err)
			os.Exit(1)
		}
	}
}
