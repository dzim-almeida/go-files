package upload

import (
	"database/sql"
	"fmt"
)

type UploadRepository interface {
	CreateUpload(upload UploadModel) (int, error)
}

type repository struct {
	conn *sql.DB
}

func NewRepository(conn *sql.DB) UploadRepository {
	return &repository{
		conn: conn,
	}
}

func (r *repository) CreateUpload(upload UploadModel) (int, error) {
	var id int

	q, err := r.conn.Prepare(
		"INSERT INTO files (original_filename, stored_filename, storage_path, content_type, size_bytes) " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING id",
	)

	if err != nil {
		fmt.Print(err)
		return -1, err
	}

	err = q.QueryRow(
		upload.OrinialFileName,
		upload.StoredFileName,
		upload.StoragePath,
		upload.ContentType,
		upload.Size,
	).Scan(&id)

	if err != nil {
		fmt.Print(err)
		return -1, err
	}

	q.Close()
	
	return id, nil
}
