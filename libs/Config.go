package libs

type ConnConf struct {
	ConnectTimeout int
	ReadTimeout    int
	WriteTimeout   int
	RetryTimes     int
}

type RedisConf struct {
	ConnConf
	Host         []string
	InitConnSize int
	MaxConnSize  int
	MaxIdleSec   int
	ExpireSec    int
	Password     string
}

type DdmqConf struct {
	Env               string
	PoolSize          int
	ProxyTimeoutInMs  int64
	ClientTimeoutInMs int64
	RetryTime         int
	LogPath           string
	LogSizeInGb       int
	Topics            []string
}

type PassportConfig struct {
	Caller   string
	RetryCnt int
	Location string
	Address  string
}

type LogConfig struct {
	Path  string
	Level uint32
}
