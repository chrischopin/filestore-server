package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internal Server error, cannot load index.html file")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get file stream data, err: %s", err.Error())
			return
		}
		defer file.Close()

		fileMetadata := meta.FileMetadata{
			FileName: header.Filename,
			FileAddr: "/tmp/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMetadata.FileAddr)
		if err != nil {
			fmt.Printf("Failed to created new file at local disk by os, err: %s", err.Error())
			return
		}
		defer newFile.Close()

		fileMetadata.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to copy file to new file handle, err: %s", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMetadata.FileSha1 = util.FileSha1(newFile)
		// meta.UpdateFileMetadata(fileMetadata)
		_ = meta.UpdateFileMetadataDB(fileMetadata)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

func GetFileMetadataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	// fileMetadata := meta.GetFileMeta(filehash)
	fileMetadata, err := meta.GetFileMetaFromDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fileMetadata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	fileMetas := meta.GetLastFileMetadata(limitCnt)

	data, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fileMetadata := meta.GetFileMeta(fsha1)

	file, err := os.Open(fileMetadata.FileAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	//TODO: large file need more work
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fileMetadata.FileName+"\"")

	w.Write(data)
}
