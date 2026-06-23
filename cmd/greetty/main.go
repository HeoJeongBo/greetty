// Command greetty is the entry point for the greetty CLI. The commands live in
// internal/cli; this file just wires up and runs them.
package main

import "github.com/HeoJeongBo/greetty/internal/cli"

func main() {
	cli.Execute()
}
