# ocmcp

[![License: Apache](https://img.shields.io/badge/license-Apache%202-blue)](https://opensource.org/licenses/Apache-2.0)
[![Lint Code Base](https://github.com/bstoll/ocmcp/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/bstoll/ocmcp/actions/workflows/lint.yml)
[![bazel build](https://github.com/bstoll/ocmcp/actions/workflows/bazel.yml/badge.svg?branch=main)](https://github.com/bstoll/ocmcp/actions/workflows/bazel.yml)

## Quick Start

You will need [Docker](https://docs.docker.com/get-docker/) installed to run the pre-built image.

### Gemini CLI

Update your Gemini CLI config in `~/.gemini/settings.json` to use the Docker container.

```json
{
  "mcpServers": {
    "ocmcp": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "ghcr.io/bstoll/ocmcp:latest"],
      "trust": true
    }
  }
}
```

### HTTP Server

To run the MCP server in HTTP mode:

```bash
docker run --rm -p 8080:8080 ghcr.io/bstoll/ocmcp:latest -http_address=":8080"
```

## Development

### Prerequisites

ocmcp uses Bazel for builds. We recommend using [Bazelisk](https://github.com/bazelbuild/bazelisk) to manage your Bazel installation.

#### Install Bazelisk

**macOS** (via [Homebrew](https://brew.sh/))

```bash
brew install bazelisk
```

**Linux**

```bash
sudo curl -L https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-amd64 -o /usr/local/bin/bazel
sudo chmod +x /usr/local/bin/bazel
```

**NPM (Any OS)**

```bash
npm install -g @bazel/bazelisk
```

### Building and Running

#### Run local server (HTTP)

```bash
bazel run //:ocmcp -- -http_address=":8080"
```

#### Use with Gemini CLI

If you are developing `ocmcp` and want to test your local changes with Gemini CLI, update your config in `~/.gemini/settings.json`:

```json
{
  "mcpServers": {
    "ocmcp": {
      "command": "bazel",
      "args": ["run", "//:ocmcp"],
      // IMPORTANT: Replace with the absolute path to this repository
      "cwd": "/path/to/your/ocmcp",
      "trust": true
    }
  }
}
```

#### Docker

##### Build and Load Local Image

To build the OCI image and load it into your local Docker daemon:

```bash
bazel run //:image_tarball
```

You can then run the image locally:

```bash
docker run --rm -it ocmcp:latest
```

##### Publish Image

To publish the image to GitHub Container Registry (GHCR):

1. Login to GHCR:
   ```bash
   docker login ghcr.io
   ```
2. Push the image:
   ```bash
   bazel run //:push
   ```
