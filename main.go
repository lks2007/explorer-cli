package main

import (
	"fmt"
    "log"
    "os"
    "github.com/rivo/tview"
)

func main() {
    box := tview.NewBox()
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}

    button := tview.NewButton("OK").SetSelectedFunc(func() {
        println("ok")
    })

	files, err := os.ReadDir(".")
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        if file.Type().IsDir() {
            fmt.Println("\033[34m "+file.Name()+"/\033[37m")
        }else{
            fmt.Println(file.Name())
        }
    }
}