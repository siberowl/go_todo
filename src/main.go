package main

import "fmt"
import "flag"

func main() {

  addPtr := flag.String("add", "", "Todo to add")
  removePtr := flag.Int("remove", 0, "Todo ID to remove")

  flag.Parse()

  if flag.NFlag() == 0 {
    flag.PrintDefaults()
  }

  fmt.Println("add:", *addPtr)
  fmt.Println("remove:", *removePtr)

  if *addPtr != "" {

  }
}
