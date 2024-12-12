package utils

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// # awk -F: '$2 != "!" && $2 != "*" { print $1 $2 }' /etc/shadow
// # cryptpw -S MbwWn6tXkWdSkNtA test123

const (
	userFile     string = "/etc/passwd"
	groupFile    string = "/etc/group"
	passwordFile string = "/etc/shadow"
)

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`
	Group     string `json:group`
	Shell     string `json:shell`
}

// Read json file and return slice of byte.
func ReadUsers(f string) []byte {

	jsonFile, err := os.Open(f)

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	data, _ := io.ReadAll(jsonFile)
	return data
}

// Read file /etc/passwd and return slice of users
func ReadEtcPasswd(f string) (list []string) {

	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := bufio.NewScanner(file)

	for r.Scan() {
		lines := r.Text()
		parts := strings.Split(lines, ":")
		list = append(list, parts[0])
	}
	return list
}

// Check if user on the host
func check(s []string, u string) bool {
	for _, w := range s {
		if u == w {
			return true
		}
	}
	return false
}

// Return securely generated random bytes

func CreateRandom(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	return string(b)
}

// User is created by executing shell command useradd
func AddNewUser(u *User) (bool, string) {

	encrypt := base64.StdEncoding.EncodeToString([]byte(CreateRandom(9)))

	argUser := []string{"-m", "-d", u.Directory, "-G", u.Group, "-s", u.Shell, u.Name}
	argPass := []string{"-c", fmt.Sprintf("echo %s:%s | chpasswd", u.Name, encrypt)}

	userCmd := exec.Command("useradd", argUser...)
	passCmd := exec.Command("/bin/sh", argPass...)

	if out, err := userCmd.Output(); err != nil {
		fmt.Println(err, "There was an error by adding user", u.Name)
		return false, ""
	} else {

		fmt.Printf("Output: %s\n", out)

		if _, err := passCmd.Output(); err != nil {
			fmt.Println(err)
			return false, ""
		}
		return true, encrypt
	}
}

func UserExists(name string) bool {
	return check(ReadEtcPasswd(userFile), name)
}
