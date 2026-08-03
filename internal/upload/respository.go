package upload

import (
	"database/sql"
)

type UploadRepository interface {
	CreateUpload(upload UploadModel) (int, error)
	GetUploads() ([]UploadModel, error)
	FindUploadById(id int) (UploadModel, error)
	FindUploadByOriginalFilename(originalFilename string) ([]UploadModel, error)
	DeleteUploadById(id int) error
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
		return -1, NewSaveError()
	}

	defer q.Close()

	err = q.QueryRow(
		upload.OriginalFileName,
		upload.StoredFileName,
		upload.StoragePath,
		upload.ContentType,
		upload.Size,
	).Scan(&id)

	if err != nil {
		return -1, NewSaveError()
	}
	
	return id, nil
}

func (r *repository) GetUploads() ([]UploadModel, error) {
	var uploads []UploadModel

	rows, err := r.conn.Query(
		"SELECT * FROM files",
	)

	if err != nil {
		return []UploadModel{}, NewInternalServerError()
	}

	defer rows.Close()

	for rows.Next() {
		var upload UploadModel

		err = rows.Scan(
			&upload.ID,
			&upload.OriginalFileName,
			&upload.StoredFileName,
			&upload.StoragePath,
			&upload.ContentType,
			&upload.Size,
		)

		if err != nil {
			return []UploadModel{}, NewInternalServerError()
		}

		uploads = append(uploads, upload)
	}

	return uploads, nil
}

func (r *repository) FindUploadById(upload_id int) (UploadModel, error) {
	var upload UploadModel

	q, err := r.conn.Prepare(
		"SELECT * FROM files WHERE id = $1",
	)

	if err != nil {
		return UploadModel{}, NewInternalServerError()
	}

	defer q.Close()

	err = q.QueryRow(upload_id).Scan(
		&upload.ID,
		&upload.OriginalFileName,
		&upload.StoredFileName,
		&upload.StoragePath,
		&upload.ContentType,
		&upload.Size,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return UploadModel{}, NewNotFoundError()
		}

		return UploadModel{}, NewInternalServerError()
	}

	return upload, nil
}

func (r *repository) FindUploadByOriginalFilename(originalFilename string) ([]UploadModel, error) {
	var uploads []UploadModel

	q, err := r.conn.Prepare(
		"SELECT * FROM files WHERE original_filename LIKE $1",
	)

	if err != nil {
		return []UploadModel{}, NewInternalServerError()
	}

	defer q.Close()

	rows, err := q.Query(originalFilename + "%")

	if err != nil {
		return []UploadModel{}, NewInternalServerError()
	}

	defer rows.Close()

	for rows.Next() {
		var upload UploadModel

		err = rows.Scan(
			&upload.ID,
			&upload.OriginalFileName,
			&upload.StoredFileName,
			&upload.StoragePath,
			&upload.ContentType,
			&upload.Size,
		)

		if err != nil {
			return []UploadModel{}, NewInternalServerError()
		}

		uploads = append(uploads, upload)
	}
	
	return uploads, nil
}

func (r *repository) DeleteUploadById(upload_id int) error {
	_, err := r.conn.Exec("DELETE FROM files WHERE id = $1", upload_id)

	if err != nil {
		return NewDeleteError()
	}

	return nil
}
