package meta

import (
	"sort"
	mysqldb "filestore-server/db"
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
	return mysqldb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName,
		fmeta.FileSize, fmeta.FileAddr)
}

func GetFileMeta(filesha1 string) FileMetadata {
	return fileMetadataStore[filesha1]
}

func GetLastFileMetadata(count int) []FileMetadata {
	fMetaArray := make([]FileMetadata, len(fileMetadataStore))
	for _, v := range fileMetadataStore {
		fMetaArray = append(fMetaArray, v)
	}
	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}
