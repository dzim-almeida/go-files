package upload

const (
	MaxFileSize = 10 * 1024 * 1024
)

var AllowedFileTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",

	".pdf": "application/pdf",

	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",

	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",

	".txt": "text/plain",

	".csv": "text/csv",

	".xml": "application/xml",

	".json": "application/json",
}