package main

import (
    "log"
    "os"
    "github.com/rivo/tview"
    "github.com/gdamore/tcell"
)

func addListFolder() []string {
    files, err := os.ReadDir(".")
    if err != nil {
        log.Fatal(err)
    }

    element := []string{}

    for _, file := range files {
        if file.Type().IsDir() {
            join := " "+file.Name()+"/"
            element = append(element, join)
        }else{
            element = append(element, file.Name())
        }
    }

    return element
}


func main() {
    app := tview.NewApplication()

    list := tview.NewList().ShowSecondaryText(false)
    list.Clear()
    for _, listValue := range addListFolder() {
        list.AddItem(listValue, "", 0, nil)
    }

    flex := tview.NewFlex().
    AddItem(list, 0, 1, false).
    AddItem(tview.NewBox().SetBorder(true), 0, 1, false)

    if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(flex).Run(); err != nil {
    	panic(err)
    }
}