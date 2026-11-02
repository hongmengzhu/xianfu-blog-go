package entityRam

type RamAccountExtraData struct {
	Tag []string `gorm:"column:tag;comment:标签" json:"tag"`
}
