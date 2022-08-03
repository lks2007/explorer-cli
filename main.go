package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func findIcon(name string) string{
    fileContent, _ := os.Open("langage.json")

    defer fileContent.Close()

    byteResult, _ := ioutil.ReadAll(fileContent)

    var res map[string]interface{}
    json.Unmarshal([]byte(byteResult), &res)
    
    contentType := filepath.Ext(name)
    
    result := fmt.Sprintf("%v", res[contentType])
    return result
}

func addListFolder() [][]string {
    files, err := os.ReadDir(".")
    if err != nil {
        log.Fatal(err)
    }

    element := [][]string{}

    for _, file := range files {
        if file.Type().IsDir() {
            join := " "+file.Name()+"/"
            element = append(element, []string{join, "120, 20, 200"})
        }else{
            if findIcon(file.Name()) != "<nil>"{
                element = append(element, []string{findIcon(file.Name())+" "+file.Name(), "190, 220, 10"})
            }else {
                text := " "+file.Name()
                secondColor := "50, 200, 255"
                element = append(element, []string{text, secondColor})
            }
        }
    }

    return element
}


func main() {
    app := tview.NewApplication()

    defaultColor := tcell.ColorBlue
    
    list := tview.NewList().ShowSecondaryText(false).SetMainTextColor(defaultColor)
    list.Clear()
    for _, listValue := range addListFolder() {
        // rgb := strings.Split(listValue[1], ", ")
        // fmt.Print(rgb)
        
        // rgbOneString, _ := strconv.ParseInt(rgb[0], 10, 32)
        // rgbOne := int32(rgbOneString)
        // rgbSecString, _ := strconv.ParseInt(rgb[1], 10, 32)
        // rgbSec := int32(rgbSecString)
        // rgbThiString, _ := strconv.ParseInt(rgb[2], 10, 32)
        // rgbThi := int32(rgbThiString)

        list.AddItem(listValue[0], "", 0, nil)
    }

    flex := tview.NewFlex().
    AddItem(list, 0, 1, false).
    AddItem(tview.NewBox().SetBorder(true), 0, 1, false)

    if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(flex).Run(); err != nil {
    	panic(err)
    }
}