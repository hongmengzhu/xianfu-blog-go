package modAttachment

type PutFileDto struct {
	Name string `label:"原始文件名称"`
	Size int64  `label:"文件大小"`
}

func PutFile(Name string, Size int64) PutFileDto {
	return PutFileDto{Name: Name, Size: Size}
}
