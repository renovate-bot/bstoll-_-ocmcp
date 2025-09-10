// Package yang provides YANG model MCP server tools.
package yang

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// YangModelIndex provides YANG model MCP server tools.
type YangModelIndex struct {
	SrcDirs []string

	db *sql.DB
}

type searchRequest struct {
	Query string `json:"query" jsonschema:"the search query"`
}

type searchResponse struct {
	Results []searchResult `json:"results" jsonschema:"results of search query"`
}

type searchResult struct {
	Path        string `json:"path" jsonschema:"YANG path"`
	Description string `json:"description" jsonschema:"description"`
	Type        string `json:"type" jsonschema:"data type"`
}

// loadYangModels processes YANG models from disk and loads them into a sqlite
// database.
func (y *YangModelIndex) loadYangModels() error {
	srcFiles, err := filePaths(y.SrcDirs)
	if err != nil {
		return err
	}

	entries, err := processModules(srcFiles)
	if err != nil {
		return err
	}
	records, err := processEntries(entries)
	if err != nil {
		return err
	}

	db, err := initDB(records)
	if err != nil {
		return err
	}
	y.db = db

	return nil
}

// RegisterTools adds the YangModelIndex tools to the provided MCP server.
func (y *YangModelIndex) RegisterTools(server *mcp.Server) error {
	if err := y.loadYangModels(); err != nil {
		return err
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "openconfig_yang_model_search",
		Description: "Searches OpenConfig YANG Models for a provided query string",
	}, y.searchModel)

	return nil
}

func (y *YangModelIndex) searchModel(ctx context.Context, req *mcp.CallToolRequest, input searchRequest) (*mcp.CallToolResult, searchResponse, error) {
	stmt, err := y.db.Prepare("SELECT path, description, type FROM path_index WHERE path_index MATCH ? ORDER BY rank LIMIT 100")
	if err != nil {
		return nil, searchResponse{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(sanitizeQuery(input.Query))
	if err != nil {
		return nil, searchResponse{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	results := []searchResult{}
	for rows.Next() {
		var path, description, typ string
		if err := rows.Scan(&path, &description, &typ); err != nil {
			return nil, searchResponse{}, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, searchResult{Path: path, Description: description, Type: typ})
	}
	return nil, searchResponse{
		Results: results,
	}, nil
}

func sanitizeQuery(query string) string {
	words := strings.Fields(query)
	sanitizedWords := make([]string, len(words))
	for i, word := range words {
		if strings.Contains(word, "-") {
			sanitizedWords[i] = fmt.Sprintf(`"%s"`, word)
		} else {
			sanitizedWords[i] = word
		}
	}
	return strings.Join(sanitizedWords, " ")
}
