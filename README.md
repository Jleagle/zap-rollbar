# zap-rollbar

## Installation

```bash
go get github.com/Jleagle/zap-rollbar
```

## Usage

```go
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	client := rollbar.NewAsync("YOUR_TOKEN", "production", "", "", "")
	defer client.Close()

	core := zaprollbar.NewCore(client, zaprollbar.WithMinLevel(zapcore.WarnLevel))

	logger := zap.New(zapcore.NewTee(
		core,
		// Add other cores like console or file here
	))
	defer logger.Sync()

	logger.Error("something went wrong", zap.String("foo", "bar"))
}

```
