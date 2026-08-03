package upload

type UploadModel struct {
	ID              int
	OriginalFileName string
	StoredFileName  string
	StoragePath     string
	ContentType     string
	Size            int64
}
