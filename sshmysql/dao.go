package sshmysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"time"
)

var sshcon *ssh.Client
var sshConfig *ssh.ClientConfig
var dbSourceStr string
var sshSourceStr string
var Db *sql.DB
var err error

func Sshinit() {
	sshHost := viper.Get("ssh_host").(string)               // SSH Server Hostname/IP
	sshPort := viper.Get("ssh_port").(string)               // SSH Port
	sshUser := viper.Get("ssh_user").(string)               // SSH Username
	sshPubKeyPath := viper.Get("ssh_pub_key_path").(string) // SSH publickey path
	dbUser := viper.Get("db_user").(string)                 // DB username
	dbPass := viper.Get("db_pass").(string)                 // DB Password
	dbHost := viper.Get("db_host").(string)                 // DB Hostname/IP
	dbName := viper.Get("db_name").(string)                 // DB database name
	Db, err = GetDbHandler(sshHost, sshPort, sshUser, sshPubKeyPath, dbUser, dbPass, dbHost, dbName)
	if err != nil {
		log.Println("获取数据库失败")
		panic(err)
	}
}

func GetDbHandler(sshHost, sshPort, sshUser, sshPubKeyPath, dbUser, dbPass, dbHost, dbName string) (*sql.DB, error) {
	if err != nil {
		log.Println("获取ssh mysql配置文件失败")
		panic(err)
	}
	sshAuthMethod := PublicKeyFile(sshPubKeyPath)
	// The client configuration with configuration option to use the ssh-agent
	sshConfig = &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			sshAuthMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         24 * 3600 * time.Second,
	}
	dbSourceStr = fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s", dbUser, dbPass, dbHost, dbName)
	sshSourceStr = fmt.Sprintf("%s:%s", sshHost, sshPort)
	// Connect to the SSH Server
	sshcon, err = ssh.Dial("tcp", sshSourceStr, sshConfig)
	if err != nil {
		log.Println("ssh服务连接失败")
		panic(err)
	}
	return connect()
}

type sSHDialer struct {
	client *ssh.Client
}

func (sd *sSHDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return sd.client.Dial("tcp", addr)
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func connect() (db *sql.DB, err error) {
	// Now we register the sSHDialer with the ssh connection as a parameter
	mysql.RegisterDialContext("mysql+tcp", (&sSHDialer{sshcon}).Dial)
	return sql.Open("mysql", dbSourceStr)
}
