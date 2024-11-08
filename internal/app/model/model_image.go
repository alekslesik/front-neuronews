package model

type Image struct {
	ImageID   int64
	ImagePath string
	Size      int64
	Name      string
	Alt       string
}

type ImageModel interface {
	InsertImage(*Image) error
	SelectGalleryImage(limit int) ([]Image, error)
}
