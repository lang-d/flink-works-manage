package main_test

import (
	"github.com/land-d/flink-works-amage/core"
	"testing"
)

func TestStart(t *testing.T) {
	core.NewApp().SetModel(core.DEV_MODEL).SetHost("0.0.0.0").SetPort(10130).Run()
}
