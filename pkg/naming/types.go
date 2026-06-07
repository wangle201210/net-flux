package naming

type DiscoSetting struct {
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	Namespace    string `mapstructure:"namespace" json:"namespace" yaml:"namespace"`
	LogDir       string `mapstructure:"log_dir" json:"log_dir" yaml:"log_dir"`
	CacheDir     string `mapstructure:"cache_dir" json:"cache_dir" yaml:"cache_dir"`
	PreloadCache bool   `mapstructure:"preload_cache" json:"preload_cache" yaml:"preload_cache"`
	Timeout      int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	GroupName    string `mapstructure:"group" json:"group" yaml:"group"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Node         int    `mapstructure:"node" json:"node" yaml:"node"`
}

type ConfigItem struct {
	Id      string `param:"id"`
	DataId  string `param:"dataId"`
	Group   string `param:"group"`
	Content string `param:"content"`
	Md5     string `param:"md5"`
	Tenant  string `param:"tenant"`
	Appname string `param:"appname"`
}

type ConfigPage struct {
	TotalCount     int          `param:"totalCount"`
	PageNumber     int          `param:"pageNumber"`
	PagesAvailable int          `param:"pagesAvailable"`
	PageItems      []ConfigItem `param:"pageItems"`
}
