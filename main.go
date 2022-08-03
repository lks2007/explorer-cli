package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func findIcon(name string) (string){
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
            join := file.Name()+"/"
            icon := ""
            element = append(element, []string{join, "136, 175, 255", icon})
        }else{
            if findIcon(file.Name()) != "<nil>"{
                data := strings.Split(findIcon(file.Name()), "|")

                element = append(element, []string{file.Name(), string([]rune(data[1])[:len(data[1])-1]), data[0][1:]})
            }else {
                text := file.Name()
                icon := ""
                secondColor := "5, 191, 90"
                element = append(element, []string{text, secondColor, icon})
            }
        }
    }

    return element
}


func main() {
    app := tview.NewApplication()

    table := tview.NewTable().SetBorders(false)
    table.SetBorderPadding(1,0,2,1)
    table.Clear()
    table.InsertColumn(1)

    for index, listValue := range addListFolder() {
        rgb := strings.Split(listValue[1], ", ")
        
        rgbOneString, _ := strconv.ParseInt(rgb[0], 10, 32)
        rgbOne := int32(rgbOneString)
        rgbSecString, _ := strconv.ParseInt(rgb[1], 10, 32)
        rgbSec := int32(rgbSecString)
        rgbThiString, _ := strconv.ParseInt(rgb[2], 10, 32)
        rgbThi := int32(rgbThiString)

        table.SetCell(index, 0,
            tview.NewTableCell(listValue[2]).
                SetTextColor(tcell.NewRGBColor(rgbOne, rgbSec, rgbThi)).
                SetAlign(tview.AlignLeft))
                
                table.SetCell(index, 1,
                    tview.NewTableCell(listValue[0]).
                        SetTextColor(tcell.ColorWhite).
                        SetAlign(tview.AlignLeft))
    }

    flex := tview.NewFlex().
    AddItem(table, 0, 1, false).
    AddItem(tview.NewBox().SetBorder(true), 0, 1, false)

    if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(flex).Run(); err != nil {
    	panic(err)
    }
}