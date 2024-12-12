package server

import (
	"easy-share/utils"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

var DataDir string
var ConfigDir string
var VarDir string
var ServerPidFile string
var Testing bool

var (
	SHUTDOWN = syscall.SIGTERM
)

type SambaGlobalConfig struct {
	ServerString  string   `json:"serverString"`
	Workgroup     string   `json:"workgroup"`
	VfsObjects    []string `json:"vfsObjects"`
	AllowedHosts  []string `json:"allowedHosts"`
	GuestUser     string   `json:"guestUser"`
	EnableRecycle bool     `json:"enableRecycle"`
	EnableNmbd    bool     `json:"enableNmbd"`
}

type SambaShareConfig struct {
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Browsable     bool     `json:"browsable"`
	Writable      bool     `json:"writable"`
	GuestOk       bool     `json:"guest"`
	Users         []string `json:"users"`
	Admins        []string `json:"admins"`
	Writelist     []string `json:"writelist"`
	Comment       string   `json:"comment"`
	Veto          bool     `json:"veto"`
	CatiaMappings string   `json:"catiaMappings"`
	CreateMask    string   `json:"createMask"`
	CustomOptions string   `json:"customOptions"`
}

type NfsShareConfig struct {
	Path          string `json:"path"`
	Host          string `json:"host"`
	Secure        bool   `json:"secure"`
	Writable      bool   `json:"writable"`
	Sync          bool   `json:"sync"`
	Mapping       string `json:"mapping"`
	Anonuid       int    `json:"anonuid"`
	Anongid       int    `json:"anongid"`
	CustomOptions string `json:"customOptions"`
}

type UserConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	Users             map[string]UserConfig       `json:"users"`
	SambaGlobalConfig SambaGlobalConfig           `json:"sambaGlobalConfig"`
	SambaShares       map[string]SambaShareConfig `json:"sambaShares"`
	NfsShares         map[string][]NfsShareConfig `json:"nfsShares"`
}

var config Config

func LoadConfig() (Config, error) {
	var config Config
	err := utils.ReadJson(filepath.Join(ConfigDir, "config.json"), &config)
	return config, err
}

func SaveConfig(config Config) error {
	return utils.WriteJson(filepath.Join(ConfigDir, "config.json"), config)
}

func GetConfig() Config {
	return config
}

func Init(dataDir string) error {
	utils.InitSignals([]os.Signal{SHUTDOWN})

	DataDir = dataDir
	ConfigDir = filepath.Join(dataDir, "config")
	VarDir = filepath.Join(dataDir, "var")
	ServerPidFile = filepath.Join(VarDir, "easy-share.pid")

	err := os.MkdirAll(ConfigDir, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(VarDir, 0755)
	if err != nil {
		return err
	}

	// Delete all log/pid files in var dir
	for _, pattern := range []string{"*.log*", "*.pid"} {
		files, _ := filepath.Glob(VarDir + "/" + pattern)
		for _, file := range files {
			os.Remove(file)
		}
	}

	utils.InitLog(filepath.Join(VarDir, "easy-share.log"))

	err = os.WriteFile(ServerPidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if err != nil {
		return err
	}

	config, err = LoadConfig()
	return err
}

func RefreshShares() {
	UpdateSambaConfig(config.SambaGlobalConfig, config.SambaShares)

	UpdateNfsConfig(config.NfsShares)

	// restart smbd

	// refresh exportfs
}

func IsNmbdRunning() bool {
	return utils.SignalRunning(filepath.Join(VarDir, "nmbd.pid"), syscall.SIGTERM)
}

func IsSmbdRunning() bool {
	return utils.SignalRunning(filepath.Join(VarDir, "smbd.pid"), syscall.SIGTERM)
}

func EnableNmbd() error {
	_, err := utils.RunCommand(utils.UseSudo, "nmbd", "-D", "-s", VarDir)
	return err
}

func DisableNmbd() error {
	utils.SignalRunning(filepath.Join(VarDir, "nmbd.pid"), syscall.SIGTERM)
	return nil
}

func AddUser(config UserConfig) error {
	return nil
}

func DelUser(config UserConfig) error {
	return nil
}

func AddSambaShare(config SambaShareConfig) error {
	return nil
}

func DelSambaShare(config SambaShareConfig) error {
	return nil
}

func AddNfsShare(config NfsShareConfig) error {
	return nil
}

func DelNfsShare(config NfsShareConfig) error {
	return nil
}
