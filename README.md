# zap-rollbar

## Installation

```bash
go get github.com/Jleagle/zap-rollbar
```

## Usage

```go
package main

import (
	"github.com/Jleagle/zap-rollbar"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	core := zaprollbar.NewCore("YOUR_TOKEN", zaprollbar.WithEnvironment("production"))

	logger := zap.New(zapcore.NewTee(
		core,
		// Add other cores like console or file here
	))
	defer logger.Sync()

	logger.Error("something went wrong", zap.String("foo", "bar"))
}
```
