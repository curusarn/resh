package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/curusarn/resh/common"
	"github.com/mattn/go-shellwords"
	giturls "github.com/whilp/git-urls"
)

// Version from git set during build
var Version string

// Revision from git set during build
var Revision string

func main() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	historyPath := filepath.Join(dir, ".resh_history.json")
	// outputPath := filepath.Join(dir, "resh_history_sanitized.json")
	sanitizerDataPath := filepath.Join(dir, ".resh", "sanitizer_data")

	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")
	// outputToStdout := flag.Bool("stdout", false, "Print output to stdout instead of file")

	flag.Parse()

	if *showVersion == true {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *showRevision == true {
		fmt.Println(Revision)
		os.Exit(0)
	}
	sanitizer := sanitizer{}
	err := sanitizer.init(sanitizerDataPath)
	if err != nil {
		log.Fatal("Sanitizer init() error:", err)
	}

	file, err := os.Open(historyPath)
	if err != nil {
		log.Fatal("Open() resh history file error:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := common.Record{}
		line := scanner.Text()
		err = json.Unmarshal([]byte(line), &record)
		if err != nil {
			log.Println("Decoding error:", err)
			log.Println("Line:", line)
			return
		}
		err = sanitizer.sanitize(&record)
		if err != nil {
			log.Println("Sanitization error:", err)
			log.Println("Line:", line)
			return
		}
		outLine, err := json.Marshal(&record)
		if err != nil {
			log.Println("Encoding error:", err)
			log.Println("Line:", line)
			return
		}
		fmt.Println(string(outLine))
	}
}

type sanitizer struct {
	GlobalWhitelist map[string]bool
	PathWhitelist   map[string]bool
	// CmdWhitelist []string
}

func (s *sanitizer) init(dataPath string) error {
	globalData := path.Join(dataPath, "whitelist.txt")
	s.GlobalWhitelist = loadData(globalData)
	pathData := path.Join(dataPath, "path_whitelist.txt")
	s.PathWhitelist = loadData(pathData)
	return nil
}

func loadData(fname string) map[string]bool {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal("Open() file error:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	data := make(map[string]bool)
	for scanner.Scan() {
		line := scanner.Text()
		data[line] = true
	}
	return data
}

func (s *sanitizer) sanitize(record *common.Record) error {
	record.Pwd = s.sanitizePath(record.Pwd)
	record.RealPwd = s.sanitizePath(record.RealPwd)
	record.PwdAfter = s.sanitizePath(record.PwdAfter)
	record.RealPwdAfter = s.sanitizePath(record.RealPwdAfter)
	record.GitDir = s.sanitizePath(record.GitDir)
	record.GitRealDir = s.sanitizePath(record.GitRealDir)
	record.Home = s.sanitizePath(record.Home)
	record.ShellEnv = s.sanitizePath(record.ShellEnv)

	record.Host = s.sanitizeTokenDontUseWhitelist(record.Host)
	record.Uname = s.sanitizeTokenDontUseWhitelist(record.Uname)
	record.Login = s.sanitizeTokenDontUseWhitelist(record.Login)
	record.MachineId = s.sanitizeTokenDontUseWhitelist(record.MachineId)

	var err error
	record.GitOriginRemote, err = s.sanitizeGitURL(record.GitOriginRemote)
	if err != nil {
		log.Println("Error while snitizing GitOriginRemote url", record.GitOriginRemote, ":", err)
		return err
	}

	fmt.Println("....")
	parser := shellwords.NewParser()

	args, err := parser.Parse(record.CmdLine)
	if err != nil {
		log.Println("Parsing error @ position", parser.Position, ":", err)
		log.Println("CmdLine:", record.CmdLine)
		return err
	}
	fmt.Println(args)

	return nil

	//	var tokens []string
	//	word := ""
	//	for _, char := range strings.Split(, "") {
	//		if unicode.IsSpace([]rune(char)[0]) {
	//			if len(word) > 0 {
	//				tokens = append(tokens, word)
	//				word = ""
	//			}
	//			tokens = append(tokens, char)
	//		} else {
	//			word += char
	//		}
	//	}
	//	if len(word) > 0 {
	//		tokens = append(tokens, word)
	//	}
	//	for _, token := range tokens {
	//		fmt.Println(token)
	//	}
	//	return nil
}

func (s *sanitizer) sanitizeGitURL(rawURL string) (string, error) {
	parsedURL, err := giturls.Parse(rawURL)
	if err != nil {
		return rawURL, err
	}
	return s.sanitizeParsedURL(parsedURL)
}

func (s *sanitizer) sanitizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, err
	}
	return s.sanitizeParsedURL(parsedURL)
}

func (s *sanitizer) sanitizeParsedURL(parsedURL *url.URL) (string, error) {
	// Scheme string
	parsedURL.Opaque = s.sanitizeToken(parsedURL.Opaque)

	userinfo := parsedURL.User.Username() // only get username => password won't even make it to the sanitized data
	if len(userinfo) > 0 {
		parsedURL.User = url.User(s.sanitizeToken(userinfo))
	} else {
		// we need to do this because `gitUrls.Parse()` sets `User` to `url.User("")` instead of `nil`
		parsedURL.User = nil
	}
	var err error
	parsedURL.Host, err = s.sanitizeTwoPartToken(parsedURL.Host, ":")
	if err != nil {
		return parsedURL.String(), err
	}
	parsedURL.Path = s.sanitizePath(parsedURL.Path)
	// ForceQuery bool
	parsedURL.RawQuery = s.sanitizeToken(parsedURL.RawQuery)
	parsedURL.Fragment = s.sanitizeToken(parsedURL.Fragment)

	return parsedURL.String(), nil
}

func (s *sanitizer) sanitizePath(path string) string {
	var sanPath string
	for _, token := range strings.Split(path, "/") {
		if s.PathWhitelist[token] != true {
			token = s.sanitizeToken(token)
		}
		sanPath += token + "/"
	}
	if len(sanPath) > 0 {
		sanPath = sanPath[:len(sanPath)-1]
	}
	return sanPath
}

func (s *sanitizer) sanitizeTwoPartToken(token string, delimeter string) (string, error) {
	tokenParts := strings.Split(token, delimeter)
	if len(tokenParts) <= 1 {
		return s.sanitizeToken(token), nil
	}
	if len(tokenParts) == 2 {
		return s.sanitizeToken(tokenParts[0]) + delimeter + s.sanitizeToken(tokenParts[1]), nil
	}
	return token, errors.New("Token has more than two parts")
}

func (s *sanitizer) sanitizeToken(token string) string {
	return s._sanitizeToken(token, true)
}

func (s *sanitizer) sanitizeTokenDontUseWhitelist(token string) string {
	return s._sanitizeToken(token, false)
}

func (s *sanitizer) _sanitizeToken(token string, useWhitelist bool) string {
	if len(token) <= 0 {
		return token
	}
	if useWhitelist == true && s.GlobalWhitelist[token] == true {
		return token
	}
	// hash with sha1
	// trim to 12 characters
	h := sha1.New()
	h.Write([]byte(token))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)[:12]
}
