package server

import (
	"github.com/Createitv/agc-cli/internal/managementtest"
	"os"
	"testing"
)

func TestMain(m *testing.M) { managementtest.ExitIfHDCStub(); os.Exit(m.Run()) }
