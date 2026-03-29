package main

import "fmt"
import "os"
import "minibox/container"


func main() {
  if len(os.Args) < 2 {
	  fmt.Println("Usage: minibox run <command>")
	  os.Exit(1)
  }

  switch os.Args[1] {
  case "run":
	  container.Run(os.Args[2:])
  case "child":
	  container.Child(os.Args[2:])
  default:
	  fmt.Println("Unknown command:", os.Args[1])
  }
}
