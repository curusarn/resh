package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/binary"
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
	"strconv"
	"strings"
	"unicode"

	"github.com/curusarn/resh/common"
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
	sanitizer := sanitizer{hashLength: 4}
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
		err = sanitizer.sanitizeRecord(&record)
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
	hashLength int
	whitelist  map[string]bool
}

func (s *sanitizer) init(dataPath string) error {
	globalData := path.Join(dataPath, "whitelist.txt")
	s.whitelist = loadData(globalData)
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

func (s *sanitizer) sanitizeRecord(record *common.Record) error {
	record.Pwd = s.sanitizePath(record.Pwd)
	record.RealPwd = s.sanitizePath(record.RealPwd)
	record.PwdAfter = s.sanitizePath(record.PwdAfter)
	record.RealPwdAfter = s.sanitizePath(record.RealPwdAfter)
	record.GitDir = s.sanitizePath(record.GitDir)
	record.GitRealDir = s.sanitizePath(record.GitRealDir)
	record.Home = s.sanitizePath(record.Home)
	record.ShellEnv = s.sanitizePath(record.ShellEnv)

	record.Host = s.hashToken(record.Host)
	record.Login = s.hashToken(record.Login)
	record.MachineId = s.hashToken(record.MachineId)

	var err error
	// this changes git url a bit but I'm still happy with the result
	// e.g. "git@github.com:curusarn/resh" becomes "ssh://git@github.com/3385162f14d7/5a7b2909005c"
	// 		notice the "ssh://" prefix
	record.GitOriginRemote, err = s.sanitizeGitURL(record.GitOriginRemote)
	if err != nil {
		log.Println("Error while snitizing GitOriginRemote url", record.GitOriginRemote, ":", err)
		return err
	}

	// sanitization destroys original CmdLine length -> save it
	record.CmdLength = len(record.CmdLine)

	record.CmdLine, err = s.sanitizeCmdLine(record.CmdLine)
	if err != nil {
		log.Fatal("Cmd:", record.CmdLine, "; sanitization error:", err)
	}
	return nil
}

func (s *sanitizer) sanitizeCmdLine(cmdLine string) (string, error) {
	sanCmdLine := ""
	buff := ""

	// simple options shouldn't be sanitized
	// 1) whitespace 2) "-" or "--" 3) letters, digits, "-", "_" 4) ending whitespace or "="
	var optionDetected bool

	prevR3 := ' '
	prevR2 := ' '
	prevR := ' '
	for _, r := range cmdLine {
		switch optionDetected {
		case true:
			if unicode.IsSpace(r) || r == '=' || r == ';' {
				// whitespace, "=" or ";" ends the option
				// => add option unsanitized
				optionDetected = false
				if len(buff) > 0 {
					sanCmdLine += buff
					buff = ""
				}
				sanCmdLine += string(r)
			} else if unicode.IsLetter(r) == false && unicode.IsDigit(r) == false && r != '-' && r != '_' {
				// r is not any of allowed chars for an option: letter, digit, "-" or "_"
				// => sanitize
				if len(buff) > 0 {
					sanToken, err := s.sanitizeCmdToken(buff)
					if err != nil {
						return cmdLine, err
					}
					sanCmdLine += sanToken
					buff = ""
				}
				sanCmdLine += string(r)
			} else {
				buff += string(r)
			}
		case false:
			// split command on all non-letter and non-digit characters
			if unicode.IsLetter(r) == false && unicode.IsDigit(r) == false {
				// split token
				if len(buff) > 0 {
					sanToken, err := s.sanitizeCmdToken(buff)
					if err != nil {
						return cmdLine, err
					}
					sanCmdLine += sanToken
					buff = ""
				}
				sanCmdLine += string(r)
			} else {
				if (unicode.IsSpace(prevR2) && prevR == '-') ||
					(unicode.IsSpace(prevR3) && prevR2 == '-' && prevR == '-') {
					optionDetected = true
				}
				buff += string(r)
			}
		}
		prevR3 = prevR2
		prevR2 = prevR
		prevR = r
	}
	if len(buff) <= 0 {
		// nothing in the buffer => work is done
		return sanCmdLine, nil
	}
	if optionDetected {
		// option detected => dont sanitize
		sanCmdLine += buff
		return sanCmdLine, nil
	}
	// sanitize
	sanToken, err := s.sanitizeCmdToken(buff)
	if err != nil {
		return cmdLine, err
	}
	sanCmdLine += sanToken
	return sanCmdLine, nil
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
		if s.whitelist[token] != true {
			token = s.hashToken(token)
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

func (s *sanitizer) sanitizeCmdToken(token string) (string, error) {
	// there shouldn't be tokens with letters or digits mixed together with symbols
	if len(token) <= 0 {
		return token, nil
	}
	if s.whitelist[token] == true {
		return token, nil
	}

	isLettersOrDigits := true
	isDigits := true
	isOtherCharacters := true
	for _, r := range token {
		if unicode.IsDigit(r) == false && unicode.IsLetter(r) == false {
			isLettersOrDigits = false
			isDigits = false
		}
		if unicode.IsDigit(r) == false {
			isDigits = false
		}
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			isOtherCharacters = false
		}
	}
	if isDigits {
		return s.hashNumericToken(token), nil
	}
	if isLettersOrDigits {
		return s.hashToken(token), nil
	}
	if isOtherCharacters {
		return token, nil
	}
	log.Println("token:", token)
	return token, errors.New("cmd token is made of mix of letters or digits and other characters")
}

func (s *sanitizer) sanitizeToken(token string) string {
	if len(token) <= 0 {
		return token
	}
	if s.whitelist[token] {
		return token
	}
	return s.hashToken(token)
}

func (s *sanitizer) hashToken(token string) string {
	if len(token) <= 0 {
		return token
	}
	// hash with sha1
	// trim to 12 characters
	h := sha1.New()
	h.Write([]byte(token))
	sum := h.Sum(nil)
	// TODO: extend hashes to 12
	return s.trimHash(hex.EncodeToString(sum))
}

func (s *sanitizer) hashNumericToken(token string) string {
	if len(token) <= 0 {
		return token
	}
	// hash with fnv
	// trim to 12 characters
	h := sha1.New()
	h.Write([]byte(token))
	sum := h.Sum(nil)
	sumInt := int(binary.LittleEndian.Uint64(sum))
	if sumInt < 0 {
		return strconv.Itoa(sumInt * -1)
	}
	return s.trimHash(strconv.Itoa(sumInt))
}

func (s *sanitizer) trimHash(hash string) string {
	length := s.hashLength
	if length <= 0 || len(hash) < length {
		length = len(hash)
	}
	return hash[:length]
}
