package banner

import (
	"bytes"
	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
)

func start() {
	isEnabled := true
	isColorEnabled := true
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString("My Custom Banner"))
}
