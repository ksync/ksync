package main

import (


  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"

  "github.com/vapor-ware/ksync/pkg/cli"
  "github.com/vapor-ware/ksync/pkg/input"
  "github.com/vapor-ware/ksync/pkg/ksync"
)

type execCmd struct {
  cli.FinderCmd
}

func (e *execCmd) new() *cobra.Command {
  long := `Create a spec then execute a command within it. This combines two steps:
      1. Create a spec with given parameters (similar to ksync create)
      2. Connect to cluster within the spec context and execute a command, optionally opening a pseudo tty and attaching.`
  example := ``

  e.Init("ksync", &cobra.Command{
    Use: "",
    Short: "Execute a command within a spec",
    Long: long,
    Example: example,
    Aliases: []string{"e"},
    Args: cobra.MinimumNArgs(2),
    Run e.run,
  })

  if err := e.DefaultFlags(); err != nil {
    log.Fatal(err)
  }

  flags := e.Cmd.Flags()

  flags.BoolP(
    "use-existing",
    "e",
    false,
    "Use an existing spec, rather than creating a new one (default: false).")
  if err := e.BindFlag("use-existing"); err != nil {
    log.Fatal(err)
  }

  flags.Bool(
    "copy-after",
    false,
    "Execute command before copying files to the spec (default: false).")
  if err := e.BindFlag("copy-after"); err != nil {
    log.Fatal(err)
  }

  flags.Bool(
    "force",
    false,
    "Force creation, ignoring similarity.")
  if err := e.BindFlag("force"); err != nil {
    log.Fatal(err)
  }

  flags.Bool(
    "reload",
    true,
    "Reload the remote container on file update.")
  if err := e.BindFlag("reload"); err != nil {
    log.Fatal(err)
  }

  flags.Bool(
    "local-read-only",
    false,
    "Set the local folder to read-only.")
  if err := e.BindFlag("local-read-only"); err != nil {
    log.Fatal(err)
  }

  flags.Bool(
    "remote-read-only",
    false,
    "Set the remote folder to read-only.")
  if err := e.BindFlag("remote-read-only"); err != nil {
    log.Fatal(err)
  }
}
