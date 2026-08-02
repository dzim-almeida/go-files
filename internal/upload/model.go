package upload

type UploadModel struct {
	OrinialFileName string
	StoredFileName  string
	StoragePath     string
	ContentType     string
	Size            int64
}
