package meta

import (
	mysqldb "filestore-server/db"
	"sort"
)

type FileMetadata struct {
	FileSha1 string
	FileName string
	FileSize int64
	FileAddr string
	UploadAt string
}

var fileMetadataStore map[string]FileMetadata

func init() {
	fileMetadataStore = make(map[string]FileMetadata)
}

func UpdateFileMetadata(fmeta FileMetadata) {
	fileMetadataStore[fmeta.FileSha1] = fmeta
}

func UpdateFileMetadataDB(fmeta FileMetadata) bool {
	return mysqldb.OnFileUploadFinishedToDB(fmeta.FileSha1, fmeta.FileName,
		fmeta.FileSize, fmeta.FileAddr)
}

func GetFileMeta(filesha1 string) FileMetadata {
	return fileMetadataStore[filesha1]
}

func GetFileMetaFromDB(filesha1 string) (FileMetadata, error) {
	tableFile, err := mysqldb.GetFileMetadataFromDB(filesha1)
	if err != nil {
		return FileMetadata{}, err
	}
	filemeta := FileMetadata{
		FileSha1: tableFile.FileHash,
		FileName: tableFile.FileName.String,
		FileSize: tableFile.FileSize.Int64,
		FileAddr: tableFile.FileAddr.String,
	}
	return filemeta, nil
}

func GetLastFileMetadata(count int) []FileMetadata {
	fMetaArray := make([]FileMetadata, len(fileMetadataStore))
	for _, v := range fileMetadataStore {
		fMetaArray = append(fMetaArray, v)
	}
	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}
