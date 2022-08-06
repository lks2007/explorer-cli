package main

import (
	"bufio"
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

func findIcon(name string, dir string) (string){
    fileContent, _ := os.Open(dir+"/"+"langage.json")

    defer fileContent.Close()

    byteResult, _ := ioutil.ReadAll(fileContent)

    var res map[string]interface{}
    json.Unmarshal([]byte(byteResult), &res)
    
    contentType := filepath.Ext(name)
    
    result := fmt.Sprintf("%v", res[contentType])

    if result == "<nil>"{
        result = fmt.Sprintf("%v", res[name])
    }

    return result
}

func addListFolder(dir string, subdir string) [][]string {
    files, err := os.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }

    element := [][]string{}
    element = append(element, []string{"../", "136, 175, 255", ""})

    for _, file := range files {
        if file.Type().IsDir() {
            join := file.Name()+"/"
            icon := ""
            element = append(element, []string{join, "136, 175, 255", icon})
        }else{
            if findIcon(file.Name(), subdir) != "<nil>"{
                data := strings.Split(findIcon(file.Name(), subdir), "|")

                element = append(element, []string{file.Name(), string([]rune(data[1])[:len(data[1])-1]), data[0][1:]})
            }else {
                text := file.Name()
                icon := ""
                secondColor := "214, 214, 214"
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

            baseHttpJson := "wget"
            argHttpJson0 := "https://raw.githubusercontent.com/lks2007/explorer-go/main/langage.json"
            cmdJson := exec.Command(baseHttpJson, argHttpJson0)
            cmdJson.Output()

            baseUnZip := "unzip"
            argUnZip0 := "-q"
            argUnZip1 := "master.zip"
            cmdZip := exec.Command(baseUnZip, argUnZip0, argUnZip1)
            cmdZip.Output()
        
            os.Chdir("icons-in-terminal-master")
            arg0 := "sh"
            arg := "install-autodetect.sh"
            cmdIcon := exec.Command(arg0, arg)
            cmdIcon.Output()

            os.Chdir("..")

            baseMv := "mv"
            argMv0 := "../explorer_go-1.0.5"
            argMv2 := "."
            cmdMv := exec.Command(baseMv, argMv0, argMv2)
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

func showList(table *tview.Table, path string, dir string)  {
    for index, listValue := range addListFolder(path, dir) {
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
}

func addListCode(file string, dir string) map[int][]string{
    f, _ := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
    defer f.Close()

    sc := bufio.NewScanner(f)
    i := 0
    lines := make(map[int][]string) 
    for sc.Scan() {
        i+=1
        lines[i] = append(lines[i], sc.Text(), fmt.Sprint(i))
    }
    return lines
}

func showCode(table *tview.Table, path string, dir string)  {
    for index, value := range addListCode(path, dir){


        table.SetCell(index, 0, tview.NewTableCell(value[1]).
        SetAlign(tview.AlignLeft).
        SetTextColor(tcell.Color102).SetSelectable(false))
                
        table.SetCell(index, 1,
            tview.NewTableCell(value[0]).
            SetTextColor(tcell.ColorWhite).
            SetAlign(tview.AlignLeft))
    }
}

func main() {
    if os.Getenv("EXPLORE_ENV") != "development"{
        initialize()
    }

    projectDir, _ := os.Getwd()

    app := tview.NewApplication()

    box := tview.NewTable()
    box.SetBorder(true)
    box.SetBorderPadding(0,0,2,0)
    box.Clear()
    box.InsertColumn(1)
    style := tcell.Style{}
    style = style.Background(tcell.NewRGBColor(29, 31, 33))
    style = style.Foreground(tcell.ColorWhite)
    box.SetSelectedStyle(style)

    box.SetDoneFunc(func(key tcell.Key) {
        if key == tcell.KeyEscape {
    		app.Stop()
    	}
    })

    table := tview.NewTable()
    table.SetBorderPadding(1,0,2,0)
    table.Clear()
    table.InsertColumn(1)

    showList(table, ".", projectDir)

    table.SetSelectable(true, false)

    table.Select(0, 0).SetFixed(1, 1).
    SetSelectedFunc(func(row int, column int) {
    	box.Clear()
        path := table.GetCell(row, 1).Text
        fileInfo, _ := os.Stat(path)


        box.SetSelectable(true, true)

        if fileInfo.IsDir() {
            table.Clear()
            showList(table, path, projectDir)
            os.Chdir(path[:len(path)-1])
        } else {
            showCode(box, path, projectDir)
            box.Select(0, 1)
        }

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
    table.SetBackgroundColor(tcell.NewRGBColor(29, 31, 33))
    box.SetBackgroundColor(tcell.NewRGBColor(29, 31, 33))

    if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(flex).Run(); err != nil {
    	panic(err)
    }
}