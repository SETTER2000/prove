package entity

type File struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Bytes []byte `json:"file"`
	Size  int64  `json:"size"`
}
