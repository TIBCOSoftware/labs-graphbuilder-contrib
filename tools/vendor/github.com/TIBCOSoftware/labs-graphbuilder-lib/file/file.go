/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package file

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//-====================-//
//    File Watcher
//-====================-//

func NewFileWatcher(handlerId int, filename string, emitPerLine bool, maxNumberOfLine int) (FileWatcher, error) {
	var fileWatcher FileWatcher
	if stat, err := os.Stat(filename); err == nil {
		switch mode := stat.Mode(); {
		case mode.IsDir():
			fmt.Println("directory")
			fileWatcher = FolderReader{
				filename:        filename,
				maxNumberOfLine: maxNumberOfLine,
				handlerId:       handlerId,
			}
		case mode.IsRegular():
			fmt.Println("file")
			if strings.HasSuffix(strings.ToLower(filename), ".zip") {
				fileWatcher = ZipFileReader{
					filename:        filename,
					emitPerLine:     emitPerLine,
					maxNumberOfLine: maxNumberOfLine,
					handlerId:       handlerId,
				}
			} else {
				fileWatcher = FileReader{
					filename:        filename,
					emitPerLine:     emitPerLine,
					maxNumberOfLine: maxNumberOfLine,
					handlerId:       handlerId,
				}
			}
		}
	} else if os.IsNotExist(err) {
		return nil, err

	} else {
		return nil, err
	}
	return fileWatcher, nil
}

type FileContentHandler interface {
	HandleContent(handlerId int, id string, content string, time int64, lineNumber int, endOfFile bool)
}

type FileWatcher interface {
	Start(handler FileContentHandler) error
}

type FileReader struct {
	handlerId       int
	filename        string
	emitPerLine     bool
	maxNumberOfLine int
}

func (this FileReader) Start(handler FileContentHandler) error {

	file, err := os.Open(this.filename)
	if nil != err {
		return err
	}

	modified, err := file.Stat()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	counter := 1
	var buffer bytes.Buffer
	for scanner.Scan() {
		if this.emitPerLine {
			handler.HandleContent(this.handlerId, "", scanner.Text(), modified.ModTime().Unix(), counter, false)
			counter += 1
		} else {
			buffer.WriteString(scanner.Text())
			buffer.WriteString("\r\n")
		}
	}

	if !this.emitPerLine {
		handler.HandleContent(this.handlerId, "", buffer.String(), modified.ModTime().Unix(), counter, true)
	} else {
		handler.HandleContent(this.handlerId, "", "", -1, -1, true)
	}

	file.Close()
	return nil
}

type FolderReader struct {
	handlerId       int
	filename        string
	maxNumberOfLine int
}

func (this FolderReader) Start(handler FileContentHandler) error {
	files, _ := ioutil.ReadDir(this.filename)
	counter := 1
	for index := range files {
		content, err := readFile(files[index].Name())
		if nil != err {
			continue
		}
		handler.HandleContent(this.handlerId, files[index].Name(), content, files[index].ModTime().Unix(), counter, false)
		counter += 1
	}
	handler.HandleContent(this.handlerId, "", "", -1, -1, true)

	return nil
}

func readFile(filename string) (string, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}
	//fmt.Println("Contents of file:", string(fileContent))
	return string(fileContent), nil
}
