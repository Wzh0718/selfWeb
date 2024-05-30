package structs

type VersionConfiguration struct {
	Version    string `json:"version"`    // 主版本号
	UpdateDate string `json:"updateDate"` // 更新时间
	Id         int    `json:"id"`         // 副版本号
}
