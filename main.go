package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/hashicorp/logutils"
	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/command"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/config/module"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	log.Printf("[INFO] Go runtime version: %s", runtime.Version())
	log.Printf("[INFO] CLI args: %#v", os.Args)

	configPath := "/Users/sebastian/gocode/src/poc"
	configFile := configPath + "/main.tf"

	var pluginOverrides command.PluginOverrides
	meta := command.Meta{
		Color:               true,
		PluginOverrides:     &pluginOverrides,
		RunningInAutomation: true}

	p := command.PlanCommand{
		Meta: meta}

	var plan *terraform.Plan

	// plan, err := p.Plan("configPath")
	// if err != nil {
	// 	panic(err)
	// }

	conf, err := config.LoadFile(configFile)
	if err != nil {
		panic(err)
	}

	bgArgs := command.BackendOpts{
		Config:     conf,
		Plan:       plan,
		ForceLocal: false}

	// Load the backend
	b, err := p.Backend(&bgArgs)
	if err != nil {
		panic(err)
	}

	var (
		destroy bool
		refresh bool
		outPath string
	)

	var mod *module.Tree
	if plan == nil {
		mod, err = p.Module(configPath)
		if err != nil {
			panic(err)
		}
	}

	// Build the operation
	opReq := p.Operation()
	opReq.Destroy = destroy
	opReq.Module = mod
	opReq.Plan = plan
	opReq.PlanRefresh = refresh
	opReq.PlanOutPath = outPath
	opReq.Type = backend.OperationTypePlan

	// Perform the operation
	op, err := b.Operation(context.Background(), opReq)
	if err != nil {
		panic(err)
	}

	// Wait for the operation to complete
	<-op.Done()
	if err := op.Err; err != nil {
		panic(err)
	}
}
