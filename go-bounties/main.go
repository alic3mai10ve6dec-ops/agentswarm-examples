package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Bounty struct {
	ID     string                 `json:"id"`
	Title  string                 `json:"title"`
	Status string                 `json:"status"`
	Bounty map[string]interface{} `json:"bounty"`
}

type Response struct {
	Bounties []Bounty `json:"bounties"`
}

func main() {
	base := os.Getenv("AGENTSWARM_BASE_URL")
	if base == "" {
		base = "https://agentswarm.space"
	}
	base = strings.TrimRight(base, "/")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(base + "/api/bounties")
	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "HTTP error: status %d %s\n", resp.StatusCode, resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP error: reading body: %v\n", err)
		os.Exit(1)
	}

	var payload Response
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Fprintf(os.Stderr, "JSON error: %v\n", err)
		os.Exit(1)
	}

	var out strings.Builder
	for _, b := range payload.Bounties {
		if b.Status != "open" {
			continue
		}
		reward, ok := b.Bounty["reward"].(string)
		if !ok {
			fmt.Fprintln(os.Stderr, "JSON error: expected string reward field")
			os.Exit(1)
		}
		out.WriteString(b.ID + "\t" + b.Title + "\t" + reward + "\t" + base + "/p/" + b.ID + "\n")
	}

	io.WriteString(os.Stdout, out.String())
}
