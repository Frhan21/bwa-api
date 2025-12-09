package entity

type FileUploadRequest struct {
	Name string
	Path string
}

type FileUploadResponse struct {
	Url      string
	PublicId string
}
