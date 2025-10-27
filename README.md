# ocmcp

## Setup

ocmcp uses Bazel for compiling. [Bazelisk](https://github.com/bazelbuild/bazelisk)
is the easiest way to get started using Bazel.

#### Linux (x86)
```bash
sudo curl -L https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-amd64 -o /usr/local/bin/bazel
sudo chmod +x /usr/local/bin/bazel
```

### Run local server (HTTP)

```bash
bazel run //:ocmcp -- -http_address=":8080"
```

### Use with Gemini CLI

Update Gemini CLI config in `~/.gemini/settings.json`.  Replace the value of cwd
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
