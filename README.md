# ocmcp

## Setup

### Run local server (HTTP)

```bash
bazel run //:ocmcp -- -http_address=":8080"
```

### Use with Gemini CLI with Bazel

Update Gemini CLI config in ~/.gemini/settings.json.  Replace the value of cwd
with the path to the ocmcp repository.
```bash
$ cat ~/.gemini/settings.json 
{
  "mcpServers": {
    "ocmcp": {
      "command": "bazel",
      "args": ["run", "//:ocmcp"],
      "cwd": "/home/user/ocmcp",
      "trust": true
    }
  }
}
```
