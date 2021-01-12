/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package file

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("trigger-filereader")

type ZipFileReader struct {
	handlerId       int
	filename        string
	emitPerLine     bool
	maxNumberOfLine int
}

func (this ZipFileReader) Start(handler FileContentHandler) error {

	zipFile, err := zip.OpenReader(this.filename)
	if nil != err {
		return err
	}
	defer zipFile.Close()

	counter := 0
	for _, file := range zipFile.File {
		if 0 != this.maxNumberOfLine && this.maxNumberOfLine <= counter {
			break
		}

		if file.FileInfo().IsDir() || strings.Contains(file.FileInfo().Name(), "DS_Store") {
			continue
		}

		content, err := readContent(file)
		if nil != err {
			continue
		}

		fmt.Println("\n\n\n")
		log.Info("(Zip entry)Process entry name begin -> ", file.FileInfo().Name())
		handler.HandleContent(this.handlerId, file.FileInfo().Name(), string(content), file.FileInfo().ModTime().Unix(), 1, false)
		counter++
		log.Info("(Zip entry)Process entry name end -> ", file.FileInfo().Name(), "\n\n\n")
	}
	handler.HandleContent(this.handlerId, "", "", -1, -1, true)

	return nil
}

func readContent(file *zip.File) ([]byte, error) {
	zipEntry, err := file.Open()
	if nil != err {
		return nil, err
	}

	defer zipEntry.Close()

	content, err := ioutil.ReadAll(zipEntry)
	if nil != err {
		return nil, err
	}

	return content, nil
}
