package base

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/spf13/viper"
)

const (
	LinkModeTUN = "tun"
	LinkModeTAP = "tap"
)

var (
	Cfg = &ServerConfig{}
)

// # ReKey time (in seconds)
// rekey-time = 172800
// # ReKey method
// # Valid options: ssl, new-tunnel
// #  ssl: Will perform an efficient rehandshake on the channel allowing
// #       a seamless connection during rekey.
// #  new-tunnel: Will instruct the client to discard and re-establish the channel.
// #       Use this option only if the connecting clients have issues with the ssl
// #       option.
// rekey-method = ssl

type ServerConfig struct {
	// LinkAddr      string `json:"link_addr"`
	ServerAddr    string `json:"server_addr"`
	AdminAddr     string `json:"admin_addr"`
	ProxyProtocol bool   `json:"proxy_protocol"`
	DbFile        string `json:"db_file"`
	CertFile      string `json:"cert_file"`
	CertKey       string `json:"cert_key"`
	UiPath        string `json:"ui_path"`
	FilesPath     string `json:"files_path"`
	LogPath       string `json:"log_path"`
	LogLevel      string `json:"log_level"`
	Issuer        string `json:"issuer"`
	AdminUser     string `json:"admin_user"`
	AdminPass     string `json:"admin_pass"`
	JwtSecret     string `json:"jwt_secret"`

	LinkMode    string `json:"link_mode"` // tun tap
	Ipv4CIDR    string `json:"ipv4_cidr"` // 192.168.1.0/24
	Ipv4Gateway string `json:"ipv4_gateway"`
	Ipv4Start   string `json:"ipv4_start"` // 192.168.1.100
	Ipv4End     string `json:"ipv4_end"`   // 192.168.1.200
	IpLease     int    `json:"ip_lease"`

	MaxClient       int    `json:"max_client"`
	MaxUserClient   int    `json:"max_user_client"`
	DefaultGroup    string `json:"default_group"`
	CstpKeepalive   int    `json:"cstp_keepalive"` // in seconds
	CstpDpd         int    `json:"cstp_dpd"`       // Dead peer detection in seconds
	MobileKeepalive int    `json:"mobile_keepalive"`
	MobileDpd       int    `json:"mobile_dpd"`

	SessionTimeout int `json:"session_timeout"` // in seconds
	AuthTimeout    int `json:"auth_timeout"`    // in seconds
}

func initServerCfg() {

	sf, _ := filepath.Abs(cfgFile)
	base := filepath.Dir(sf)

	// 转换成绝对路径
	Cfg.DbFile = getAbsPath(base, Cfg.DbFile)
	Cfg.CertFile = getAbsPath(base, Cfg.CertFile)
	Cfg.CertKey = getAbsPath(base, Cfg.CertKey)
	Cfg.UiPath = getAbsPath(base, Cfg.UiPath)
	Cfg.FilesPath = getAbsPath(base, Cfg.FilesPath)
	Cfg.LogPath = getAbsPath(base, Cfg.LogPath)

	if len(Cfg.JwtSecret) < 20 {
		fmt.Println("请设置 jwt_secret 长度20位以上")
		os.Exit(0)
	}

	fmt.Printf("ServerCfg: %+v \n", Cfg)
}

func getAbsPath(base, cfile string) string {
	if cfile == "" {
		return ""
	}

	abs := filepath.IsAbs(cfile)
	if abs {
		return cfile
	}
	return filepath.Join(base, cfile)
}

func initCfg() {
	ref := reflect.ValueOf(Cfg)
	s := ref.Elem()

	typ := s.Type()
	numFields := s.NumField()
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		value := s.Field(i)
		tag := field.Tag.Get("json")

		for _, v := range configs {
			if v.Name == tag {
				if v.Typ == cfgStr {
					value.SetString(viper.GetString(v.Name))
				}
				if v.Typ == cfgInt {
					value.SetInt(int64(viper.GetInt(v.Name)))
				}
				if v.Typ == cfgBool {
					value.SetBool(viper.GetBool(v.Name))
				}
			}
		}
	}

	initServerCfg()
}

type SCfg struct {
	Name string      `json:"name"`
	Env  string      `json:"env"`
	Info string      `json:"info"`
	Data interface{} `json:"data"`
}

func ServerCfg2Slice() []SCfg {
	ref := reflect.ValueOf(Cfg)
	s := ref.Elem()

	var datas []SCfg

	typ := s.Type()
	numFields := s.NumField()
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		value := s.Field(i)
		tag := field.Tag.Get("json")
		usage, env := getUsageEnv(tag)
		if usage == "" {
			continue
		}

		datas = append(datas, SCfg{Name: tag, Env: env, Info: usage, Data: value.Interface()})
	}

	return datas
}

func getUsageEnv(name string) (usage, env string) {
	for _, v := range configs {
		if v.Name == name {
			usage = v.Usage
		}
	}

	if e, ok := envs[name]; ok {
		env = e
	}

	return
}
