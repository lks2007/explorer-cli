package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

func addListFolder(dir string) [][]string {
    files, err := os.ReadDir(dir)
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

func initialize() {
    if _, err := os.Stat("initialize.txt"); errors.Is(err, os.ErrNotExist) {
        f, _ := os.Create("initialize.txt")
        f.WriteString("0")
        defer f.Close()
    }

    fileSecond, _ := os.Open("initialize.txt")

    buf := make([]byte, 1024)
    for {
        n, _ := fileSecond.Read(buf)

        if string(buf[:n]) == "0"{
            defer fileSecond.Close()

            mode := int(0777)
            os.Mkdir("explorer", os.FileMode(mode))
            os.Chdir("explorer")

            baseHttp := "wget"
            argHttp0 := "https://github.com/lks2007/icons-in-terminal/archive/refs/heads/master.zip"
            cmd := exec.Command(baseHttp, argHttp0)
            cmd.Output()

            baseUnZip := "unzip"
            argUnZip0 := "-q"
            argUnZip1 := "master.zip"
            cmdZip := exec.Command(baseUnZip, argUnZip0, argUnZip1)
            cmdZip.Output()
        
            os.Chdir("icons-in-terminal-master")
            arg := "./install-autodetect.sh"
            cmdIcon := exec.Command(arg)
            cmdIcon.Output()

            os.Chdir("../")

            baseMv := "mv"
            argMv0 := "../v1"
            argMv1 := "../initialize.txt"
            argMv2 := "."
            cmdMv := exec.Command(baseMv, argMv0, argMv1, argMv2)
            cmdMv.Output()

            os.Remove("initialize.txt")
            f, _ := os.Create("initialize.txt")
            f.WriteString("1")
            defer f.Close()

            break
        }
        break
    }
}

func main() {
    initialize()


    app := tview.NewApplication()

    box := tview.NewTable().
    SetBorders(true)
    box.SetBorderPadding(1,0,2,1)
    box.Clear()
    box.InsertColumn(1)

    table := tview.NewTable().SetBorders(false)
    table.SetBorderPadding(1,0,2,1)
    table.Clear()
    table.InsertColumn(1)

    for index, listValue := range addListFolder(".") {
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

    table.SetSelectable(true, false)

    table.Select(0, 0).SetFixed(1, 1).
    SetSelectedFunc(func(row int, column int) {
    	box.Clear()
        // table.GetCell(row, 1).Text
    	table.SetSelectable(true, false)
    }).
    SetDoneFunc(func(key tcell.Key) {
    	if key == tcell.KeyEscape {
    		app.Stop()
    	}
    	if key == tcell.KeyEnter {
    		table.SetSelectable(true, false)
    	}
    })

    flex := tview.NewFlex().
    AddItem(table, 0, 1, true).
    AddItem(box, 0, 1, false)

    if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(flex).Run(); err != nil {
    	panic(err)
    }
}