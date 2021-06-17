package ocrfreeapi

type Parse interface {
	ParseFromBase64(base64String string) (string, error)
	ParseFromUrl(imageUrl string) (string, error)
	ParseFromLocal(imageFilePath string) (string, error)
}