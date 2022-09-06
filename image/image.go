package image

const Prefix = "Dockerfile."

type Image struct {
	Name     string
	Tag      string
	Hash     string
	Registry string
}
