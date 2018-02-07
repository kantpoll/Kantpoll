package main

import (
	"strings"
	"net/http"
	"os"
	"log"
	"fmt"
	"os/exec"
	"io"
	"io/ioutil"
	"os/user"
	"runtime"
	"strconv"
	"path/filepath"
	"bufio"
	"flag"
	"time"
	"math/rand"
)

//Constants
const THE_AND = "--and--"
const QUESTION_MARK = "--question--"
const HASHTAG = "--hashtag--"
const QUOTE = "--quote--"
const DOUBLEQUOTE = "--doublequote--"

//Flags
var verbose bool
var customProvider string

//Set the commands according to the os
var executables map [string] map [string] string = make(map [string] map [string] string)
func add(m map[string]map[string]string, os, cmd, value string) {
	mm, ok := m[os]
	if !ok {
		mm = make(map[string]string)
		m[os] = mm
	}
	mm[cmd] = value
}

//In order to know if an instance of geth is already running
var geth_cmd *exec.Cmd

//To find peers
var enodes = ""

//The last chainid, which is currently active
var last_id = 0

func main() {
	//Loading the flags
	customProviderFlag := flag.String("provider", "", "Define a custom campaigns provider")
	verboseFlag := flag.Bool("verbose", false, "Activate the verbose mode")
	flag.Parse()
	customProvider = *customProviderFlag
	verbose = *verboseFlag

	//The map with executable paths
	loadExecutables()
	//Adding Content Security Policy
	changeHeaderThenServe := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path[1:]
			if strings.Index(path, "query") == 0 {
				if strings.Index(path, "queryMyIP=") == 0{
					whatIsMyIP(w)
					return
				}
				if strings.Index(path, "queryEnode=") == 0{
					whatIsMyEnode(w)
					return
				}
				if strings.Index(path, "queryAddProfile=") == 0{
					addProfile(path, w)
					return
				}
				if strings.Index(path, "queryAddCampaignIntoProvider=") == 0{
					addCampaignIntoProvider(path, w)
					return
				}
				if strings.Index(path, "queryProfileExists=") == 0{
					profileExists(path, w)
					return
				}
				if strings.Index(path, "queryGetProfile=") == 0{
					getProfile(path, w)
					return
				}
				if strings.Index(path, "querySetBlockchain=") == 0{
					setBlockchain(path, w)
					return
				}
				if strings.Index(path, "queryGetCampaign=") == 0{
					getCampaign(path, w)
					return
				}
				if strings.Index(path, "queryGetIPNS=") == 0{
					getIPNS(path, w)
					return
				}
				if strings.Index(path, "queryAddIPNSKey=") == 0{
					addIPNS(path, w)
					return
				}
				if strings.Index(path, "queryVerifyBlockchain=") == 0{
					verifyBlockchain(path, w)
					return
				}
				if strings.Index(path, "queryRunBlockchain=") == 0{
					runBlockchain(path, w)
					return
				}
				if strings.Index(path, "queryInsertAccountIntoBlockchain=") == 0{
					insertAccountIntoBlockchain(path, w)
					return
				}
				if strings.Index(path, "queryGetCustomProvider=") == 0{
					fmt.Fprint(w, *customProviderFlag)
					return
				}
				if strings.Index(path, "queryCheckUser=") == 0{
					checkUser(path, w)
					return
				}
				if strings.Index(path, "queryCreatePwdFile=") == 0{
					createPwdFile(path, w)
					return
				}

			}
			// Set some header
			w.Header().Add("Content-Security-Policy", "script-src http://localhost:1985; style-src http://localhost:1985 'unsafe-inline'; child-src 'none'; object-src 'none'; form-action http://localhost:1985; connect-src http://localhost:1985 http://localhost:8080 http://localhost:8545 https://ipfs.io; worker-src 'none'")
			// Serve with the actual handler
			h.ServeHTTP(w, r)
		}
	}

	//Checking if the .kantcoin directory was created. If it was not created, then create it
	initDirectory()

	//Initializing IPFS
	if !isIPFSRunning() {
		initIPFS(0)
	}

	//Just display geth version in order to run geth once before being used
	displayGethVersion()

	http.ListenAndServe(":1985", changeHeaderThenServe(http.FileServer(http.Dir("./website"))))
}

//It is important to execute this command before the user call it to avoid unnecessary delays (Firewall) which might provoke errors
func displayGethVersion(){
	_, err := exec.LookPath(executables[runtime.GOOS]["geth"])

	if err == nil {
		out, err := exec.Command(executables[runtime.GOOS]["geth"], "version").Output()
		outstr := string(out)
		if err == nil {
			if (verbose){
				log.Println("----------------------------------Command response ----------------------------------")
				log.Println(outstr)
				log.Println("------------------------------End of command response--------------------------------")
			}
		} else {
			log.Println("Geth error")
		}
	}
}

//Writing a new genesis.json file and starting a new geth instance
func setBlockchain(path string, w http.ResponseWriter){
	query := path[len("querySetBlockchain="):]
	query = strings.Replace(query, QUESTION_MARK, "?", -1)
	query = strings.Replace(query, HASHTAG, "#", -1)
	query = strings.Replace(query, DOUBLEQUOTE, "\"", -1)
	query = strings.Replace(query, QUOTE, "'", -1)

	//Obtaining the chainid and the address to be used in a specific campaign
	parts := strings.Split(query,THE_AND)
	chainid := parts[0]
	address := parts[1]
	enodes = parts[2]
	directory := parts[3]
	nonce := parts[4]

	//Creating the directory where the data of this specific blockchain will be placed
	os.Mkdir(getHome() + "/.kantcoin/blockchains/" + directory,0700)

	//Composing the genesis.json file
	genesis := "{ \"config\": {" +
		"\"chainId\": " + chainid + "," +
		"\"homesteadBlock\": 0," +
		"\"eip155Block\": 0," +
		"\"eip158Block\": 0," +
		"\"byzantiumBlock\": 0" +
		"}," +
		"\"difficulty\": \"200\"," +
		"\"gasLimit\": \"3100000000\"," +
		"\"nonce\": \""+ nonce + "\"," +
		"\"alloc\": {" +
		"\"" + address + "\": { \"balance\": \"10000000000000000000\" }" +
		"}}"


	//Writing genesis string into the file
	data := []byte(genesis)
	err := ioutil.WriteFile(getHome() + "/.kantcoin/blockchains/" + directory + "/genesis.json", data, 0700)
	if err == nil {

	} else {
		fmt.Fprint(w, "error")
		log.Println("Error while creating the genesis.json file")
		return
	}

	//Initializing Geth
	if runtime.GOOS == "windows" {
		initGethWin(w, directory)
	} else if runtime.GOOS == "linux" {
		initGethLin(w, directory)
	}
}

//Initialize the IPFS service to provide the users' pages
func initIPFS(times int){
	if times == 2 {
		log.Println("Limit of tries reached")
		return
	}
	_, err := exec.LookPath(executables[runtime.GOOS]["ipfs"])
	existDir := verifyIpfsDir()

	if err == nil && existDir {
		//This command 'daemon' provides user access to IPFS sites
		err := exec.Command(executables[runtime.GOOS]["ipfs"], "daemon").Start()
		if err == nil {
			log.Println("IPFS daemon was started")
		} else {
			log.Println("IPFS error")
		}
	} else if err == nil {
		//First we need to initialize the IPFS
		out, err := exec.Command(executables[runtime.GOOS]["ipfs"], "init").Output()
		outstr := string(out)
		if (verbose){
			log.Println("----------------------------------Command response ----------------------------------")
			log.Println(outstr)
			log.Println("------------------------------End of command response--------------------------------")
		}

		if err == nil {
			log.Println("IPFS was initialized")
			initIPFS(times + 1)
		} else {
			log.Println("IPFS was not initialized")
		}

	} else {
		log.Println("IPFS is not installed")
	}
}

//Initialize Geth node on Linux
func initGethLin(w http.ResponseWriter, directory string){
	_, err := exec.LookPath(executables[runtime.GOOS]["geth"])
	if err == nil {
		//filepath,_ := os.Getwd()
		out, err := exec.Command(executables[runtime.GOOS]["geth"], "--datadir", getHome() + "/.kantcoin/blockchains/" + directory, "init", getHome() + "/.kantcoin/blockchains/"+ directory + "/genesis.json").Output()
		outstr := string(out)

		if (verbose) {
			log.Println("----------------------------------Command response ----------------------------------")
			log.Println(outstr)
			log.Println("------------------------------End of command response--------------------------------")
		}

		if err == nil {
			log.Println("Geth was started")
			fmt.Fprint(w, "complete")
			//geth_init_time = time.Now().Second()
		} else {
			log.Println("Geth was not started")
			fmt.Fprint(w, "error")
		}
	} else {
		log.Println("Geth is not installed")
		fmt.Fprint(w, "error")
	}
}

//Initialize Geth node on Windows
func initGethWin(w http.ResponseWriter, directory string){
	var cmd *exec.Cmd
	_, err := exec.LookPath(executables[runtime.GOOS]["geth"])
	if err == nil {
		cmd = exec.Command("cmd")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		filepath,_ := os.Getwd()
		args := "( " + filepath + executables[runtime.GOOS]["geth_backslash"] + " --datadir \"" + getHome() + "\\.kantcoin\\blockchains\\" + directory + "\" init \"" + getHome() + "\\.kantcoin\\blockchains\\"+ directory + "\\genesis.json\" )"

		go func() {
			defer stdin.Close()
			io.WriteString(stdin, args)
		}()

		out, err := cmd.CombinedOutput()
		outstr := string(out)
		if (verbose) {
			log.Println("----------------------------------Command response ----------------------------------")
			log.Println(outstr)
			log.Println("------------------------------End of command response--------------------------------")
		}

		if err == nil {
			log.Println("Geth was started")
			fmt.Fprint(w, "complete")
			//geth_init_time = time.Now().Second()
		} else {
			log.Println("Geth was not started")
			log.Println(err)
			fmt.Fprint(w, "error")
		}
	} else {
		log.Println("Geth is not installed")
		fmt.Fprint(w, "error")
	}
}

//Verifying if IPFS dir exists
func verifyIpfsDir() bool{
	userprofile := getHome()
	//Check if the .ipfs folder was created in user's PC
	if _, err := os.Stat(userprofile + "/.ipfs/version"); !os.IsNotExist(err) {
		return true
	}
	return false
}

//Getting the HOME directory
func getHome() string{
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	return usr.HomeDir
}

//It receives the campaign id, user, ipns, campaign address, message, signature, and provider
func addCampaignIntoProvider(path string, w http.ResponseWriter) {
	query := path[len("queryAddCampaignIntoProvider="):]
	query = strings.Replace(query, QUESTION_MARK, "?", -1)
	query = strings.Replace(query, HASHTAG, "#", -1)
	query = strings.Replace(query, DOUBLEQUOTE, "\"", -1)
	query = strings.Replace(query, QUOTE, "'", -1)

	parts := strings.Split(query, THE_AND)
	id := parts[0]
	user := parts[1]
	ipns := parts[2]
	pkey := parts[3]
	message := parts[4]
	signature := parts[5]
	address := parts[6]
	provider := parts[7]

	resp, err := http.Get(provider + "/addCampaign?id=" + id + "&user=" + user + "&ipns=" + ipns + "&pkey=" + pkey + "&message=" + message + "&signature=" + signature + "&address=" + address)
	if err == nil {
		bodyString := ""
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes,_ := ioutil.ReadAll(resp.Body)
			bodyString = string(bodyBytes)
		}
		fmt.Fprint(w, bodyString)
	} else {
		fmt.Fprint(w, "error")
	}
}

//It checks if the informed address belongs to the informed user
func checkUser(path string, w http.ResponseWriter) {
	query := path[len("queryCheckUser="):]
	parts := strings.Split(query, THE_AND)
	id := parts[0]
	pkey := parts[1]
	provider := parts[2]

	resp, err := http.Get(provider + "/checkUser?id=" + id + "&pkey=" + pkey)
	if err == nil {
		bodyString := ""
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes,_ := ioutil.ReadAll(resp.Body)
			bodyString = string(bodyBytes)
		}
		fmt.Fprint(w, bodyString)
	} else {
		fmt.Fprint(w, "error")
	}
}

//Getting user's ip in order to figure out the enode
//We have to obtain it via third party services
func  whatIsMyIP(w http.ResponseWriter){
	rand.Seed(time.Now().UnixNano())
	urls := []string{
		"https://api.ipify.org/?format=json&callback=",
		"https://jsonip.com/?callback=",
		"https://ipinfo.io/json",
		"https://ipapi.co/json",
	}

	resp, err := http.Get(urls[rand.Intn(len(urls))])
	defer resp.Body.Close()

	if err == nil && resp.StatusCode == http.StatusOK{
		bodyBytes,_ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Fprint(w, bodyString)
	} else {
		fmt.Fprint(w, "error")
	}
}

//It calls geth in order to get the enode
func whatIsMyEnode(w http.ResponseWriter){
	cmd := exec.Command(executables[runtime.GOOS]["geth"], "attach")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprint(w, "error")
	}
	args := "admin.nodeInfo.enode"

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, args)
	}()

	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	if strings.LastIndex(outstr,"enode://") > 0{
		begin := strings.LastIndex(outstr,"enode://")
		end := strings.LastIndex(outstr,":30303")
		outstr = outstr[begin: end + 6]
	}

	fmt.Fprint(w, outstr)
}

//Insert or overwrite a new profile (person or campaign)
func addProfile(path string, w http.ResponseWriter){
	query := path[len("queryAddProfile="):]
	query = strings.Replace(query, QUESTION_MARK, "?", -1)
	query = strings.Replace(query, HASHTAG, "#", -1)
	query = strings.Replace(query, DOUBLEQUOTE, "\"", -1)
	query = strings.Replace(query, QUOTE, "'", -1)
	parts := strings.Split(query, THE_AND)
	dir := parts[0]
	content := parts[1]
	filename := parts[2]

	// Do not save kroot profile
	if strings.Index(parts[0],"kroot") >= 0 {
		return
	}

	//Profiles should be saved in the HOME directory
	os.Mkdir(getHome() + "/.kantcoin/profiles/" + dir, 0700)
	f, err := os.Create(getHome() + "/.kantcoin/profiles/" + dir + "/" + filename)
	if err == nil {
		_, err := f.WriteString(content)
		if err == nil{
			log.Println("New profile: " + dir)

			//Executing IPFS
			_, err = exec.LookPath(executables[runtime.GOOS]["ipfs"])
			if err == nil {
				//The error used to verify whether the profile was inserted in IPFS or not.
				var err1 *appError
				if runtime.GOOS == "windows"{
					err1 = newIPFSKeyWin(dir)
					if err1 == nil{
						address, err1 := newIPFSPageWin(dir, w, false)
						if err1 == nil {
							err1 = publishIPFSWin(address, dir)
						}
					}

				} else if runtime.GOOS == "linux"{
					err1 = newIPFSKeyLin(dir)
					if err1 == nil {
						address, err1 := newIPFSPageLin(dir, w, false)
						if err1 == nil {
							err1 = publishIPFSLin(address, dir)
						}
					}
				}
				if err1 == nil{
					fmt.Fprint(w, "complete")
				} else {
					fmt.Fprint(w, "error")
				}
			} else {
				fmt.Fprint(w, "error")
				log.Println("IPFS not installed")
			}
		}
		f.Close()
	} else {
		fmt.Fprint(w, "error")
		log.Println("File was not created")
	}
}

//The error used to verify whether the profile was inserted in IPFS or not.
type appError struct {
	Error   error
	Message string
}

//We need to use the cmd when the os is the windows
func newIPFSPageWin(dir string, w http.ResponseWriter, show bool) (string, *appError){
	cmd := exec.Command("cmd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	filepath,_ := os.Getwd()
	args := "( " + filepath + executables[runtime.GOOS]["ipfs_backslash"] + " add -r \"" + getHome() + "\\.kantcoin\\profiles\\" + dir + "\" )"

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, args)
	}()

	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	//Getting the ipfs address of the directory
	if len(outstr) >= 52{
		begin := strings.LastIndex(outstr,"added")
		outstr = outstr[begin + 6: begin + 6 + 46]
	}

	if err == nil {
		log.Println("IPFS page was inserted")
		if (show){
			fmt.Fprint(w, outstr)
		}
	} else {
		log.Println("IPFS insertion error")
		log.Println(err)
		if (show){
			fmt.Fprint(w, "error")
		}
		return "", &appError{err, "IPFS insertion error"}
	}
	return outstr, nil
}

//Before publishing we need to create a key (WINDOWS)
func newIPFSKeyWin(dir string) *appError{
	cmd := exec.Command("cmd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	filepath,_ := os.Getwd()
	args := "( " + filepath + executables[runtime.GOOS]["ipfs_backslash"] + " key gen --type=rsa --size=2048 " + dir + " )"

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, args)
	}()

	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	if err == nil {
		if strings.Index(outstr,"refusing to overwrite") > 0{
			log.Println("IPFS key already exists")
		} else {
			log.Println("IPFS key was created")
		}
	} else {
		log.Println("IPFS key creation error")
		log.Println(err)
		return &appError{err, "IPFS key creation error"}
	}
	//No error
	return nil
}

//Publish some page with some key(dir) (WINDOWS)
func publishIPFSWin(address, dir string) *appError{
	cmd := exec.Command("cmd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	filepath,_ := os.Getwd()
	args := "( " + filepath + executables[runtime.GOOS]["ipfs_backslash"] + " name publish --key=" + dir + " " + address + " )"

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, args)
	}()

	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	if err == nil {
		log.Println("Directory " + dir + " was published on IPFS")
	} else {
		log.Println("Publishing error with the directory " + dir)
		return &appError{err, "Publishing error with the directory " + dir}
	}
	return nil
}

//Adding new IPFS page on Linux
func newIPFSPageLin(dir string, w http.ResponseWriter, show bool) (string, *appError){
	out, err := exec.Command(executables[runtime.GOOS]["ipfs_backslash"],"add","-r", getHome() + "/.kantcoin/profiles/" + dir).Output()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	//Getting the ipfs address of the directory
	if (len(outstr) >=52){
		begin := strings.LastIndex(outstr,"added")
		outstr = outstr[begin + 6: begin + 6 + 46]
	}

	if err == nil {
		log.Println("IPFS page was inserted")
		if (show){
			fmt.Fprint(w, outstr)
		}
	} else {
		log.Println("IPFS insertion error")
		if (show){
			fmt.Fprint(w, "IPFS insertion error")
		}
		return "", &appError{err, "IPFS insertion error"}
	}

	return outstr, nil
}

//Before publishing we need to create a key (LINUX)
func newIPFSKeyLin(dir string) *appError{
	out, err := exec.Command(executables[runtime.GOOS]["ipfs_backslash"],"key","gen", "--type=rsa", "--size=2048" + dir).Output()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	if err == nil {
		if strings.Index(outstr,"refusing to overwrite") > 0{
			log.Println("IPFS key already exists")
		} else {
			log.Println("IPFS key was created")
		}
		log.Println(outstr)
	} else {
		log.Println("IPFS key creation error")
		log.Println(err)
		return &appError{err, "IPFS key creation error"}
	}
	return nil
}

func publishIPFSLin(address, dir string) *appError{
	out, err := exec.Command(executables[runtime.GOOS]["ipfs_backslash"],"name","publish", "--key=" + dir, address).Output()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	if err == nil {
		log.Println("Directory " + dir + " was published on IPFS")
	} else {
		log.Println("Publishing error with the directory " + dir)
		log.Println(err)
		return &appError{err, "Publishing error with the directory " + dir}
	}
	return nil
}

//Check wheter the profile exists or not, returning 'true' or 'false'
func profileExists(path string, w http.ResponseWriter){
	userprofile := getHome()
	query := path[len("queryProfileExists="):]
	//Template profile
	if strings.Index(query,"kroot") >= 0{
		fmt.Fprint(w, "true")
		return
	}
	if _, err := os.Stat(userprofile + "/.kantcoin/profiles/" + query); !os.IsNotExist(err) {
		log.Println("Profile " + query + " opened")
		fmt.Fprint(w, "true")
	} else {
		log.Println("Profile " + query + " does not exist")
		fmt.Fprint(w, "false")
	}
}

//Returns the profile html content
func getProfile(path string, w http.ResponseWriter){
	userprofile := getHome()
	query := path[len("queryGetProfile="):]
	file := ""

	//Firstly, check if this request is for a profile...
	if strings.Index(query,"profile") > 0 {
		//Initial page
		if strings.Index(query,"kroot") == 0 {
			file = "website/templates/kroot/profile"
			content, err := ioutil.ReadFile(file)
			if err == nil {
				fmt.Fprint(w, string(content))
				return;
			}
		} else {
			//A profile that has been already created
			file = userprofile + "/.kantcoin/profiles/" + query
			content, err := ioutil.ReadFile(file)
			if err == nil {
				fmt.Fprint(w, string(content))
				return;
			} else {
				//If this profile does not exist, show the initial page
				file = "website/templates/kroot/profile"
				content, err := ioutil.ReadFile(file)
				if err == nil {
					fmt.Fprint(w, string(content))
					return;
				}
			}
		}
	} else if strings.Index(query,"data") > 0 {  //...or for user data
		//There is no default user data
		file = userprofile + "/.kantcoin/profiles/" + query
		content, err := ioutil.ReadFile(file)
		if err == nil {
			fmt.Fprint(w, string(content))
			return;
		}
	}

	fmt.Fprint(w, "error")
}

//Call the login provider passing the campaign id
func getCampaign(path string, w http.ResponseWriter){
	query := path[len("queryGetCampaign="):]
	parts := strings.Split(query, THE_AND)
	id := parts[0]
	provider := parts[1]

	if (len(parts) != 2){
		fmt.Fprint(w, "error")
		return
	}

	resp, err := http.Get(provider + "/getCampaign?id=" + id)
	if err == nil {
		bodyString := ""
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes,_ := ioutil.ReadAll(resp.Body)
			bodyString = string(bodyBytes)
		}
		fmt.Fprint(w, bodyString)
	} else {
		fmt.Fprint(w, "error")
	}
}

//In order to avoid this error: 'Failed to write genesis block: database already contains an incompatible genesis block'
func removeBlockchainDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

//This function receives the same params received by the setBlockchain function, then it compares the content of the genesis file and the params
func verifyBlockchain(path string, w http.ResponseWriter) {
	query := path[len("queryVerifyBlockchain="):]
	query = strings.Replace(query, QUESTION_MARK, "?", -1)
	query = strings.Replace(query, HASHTAG, "#", -1)
	query = strings.Replace(query, DOUBLEQUOTE, "\"", -1)
	query = strings.Replace(query, QUOTE, "'", -1)

	//Obtaining the chainid and the address to be used in a specific campaign
	parts := strings.Split(query, THE_AND)
	chainid := parts[0]
	address := parts[1]
	//enode := parts[2]
	directory := parts[3]
	nonce := parts[4]
	delete_dir_if_different := parts[5]

	//Composing the genesis.json file
	genesis := "{ \"config\": {" +
		"\"chainId\": " + chainid + "," +
		"\"homesteadBlock\": 0," +
		"\"eip155Block\": 0," +
		"\"eip158Block\": 0," +
		"\"byzantiumBlock\": 0" +
		"}," +
		"\"difficulty\": \"200\"," +
		"\"gasLimit\": \"3100000000\"," +
		"\"nonce\": \""+ nonce + "\"," +
		"\"alloc\": {" +
		"\"" + address + "\": { \"balance\": \"10000000000000000000\" }" +
		"}}"

	content, err := ioutil.ReadFile(getHome() + "/.kantcoin/blockchains/" + directory + "/genesis.json")

	if err == nil {
		if string(content) == genesis {
			fmt.Fprint(w, "true")
		} else {
			if delete_dir_if_different == "true" {
				//Removing these directories to avoid this error: Failed to write genesis block: database already contains an incompatible genesis block
				removeBlockchainDir(getHome() + "/.kantcoin/blockchains/" + directory + "/geth")
				removeBlockchainDir(getHome() + "/.kantcoin/blockchains/" + directory + "/keystore")
			}
			fmt.Fprint(w, "false")
		}
	} else {
		fmt.Fprint(w, "error")
	}
}

//Call: geth --networkid "1151985..." etc
func runBlockchain(path string, w http.ResponseWriter) {
	/*if (time.Now().Second() < geth_init_time + 3){
		//Just waiting a few seconds before running geth
		time.Sleep(1 * time.Second)
	}*/
	query := path[len("queryRunBlockchain="):]
	parts := strings.Split(query, THE_AND)
	id := parts[0]
	address := parts[1]
	dir := parts[2]

	int_id,_ := strconv.Atoi(id)
	//If the campaign has not changed, keep the geth process running
	if int_id == last_id{
		log.Println("Old geth process was kept")
		fmt.Fprint(w, "complete")
	} else {
		if geth_cmd != nil && geth_cmd.Process != nil {
			err := geth_cmd.Process.Kill()
			if err == nil{
				log.Println("Old geth process was killed")
			} else {
				log.Println("Old geth process was not killed")
			}
		}

		if runtime.GOOS == "windows" {
			runBlockchainWin(w, id, address, dir)
		} else if runtime.GOOS == "linux" {
			runBlockchainLin(w, id, address, dir)
		}

		last_id = int_id
	}
}

//Call geth through cmd
func runBlockchainWin( w http.ResponseWriter, id, address, dir string){
	geth_cmd = exec.Command("cmd")
	stdin, err := geth_cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	filepath,_ := os.Getwd()

	var bootnodes = ""
	if enodes != ""{
		bootnodes = " --bootnodesv5 " + enodes
	}

	args := "( " + filepath + executables[runtime.GOOS]["geth_backslash"] + bootnodes + " --datadir \"" + getHome() + "\\.kantcoin\\blockchains\\" + dir + "\" --networkid \"" + id + "\" --lightkdf --shh --rpc --rpcport 8545 --rpcapi \"admin,db,eth,net,web3,shh\" --rpccorsdomain \"*\" --mine --minerthreads=1 --unlock \"" + address + "\" --etherbase \"" + address + "\" --password \"" + getHome() + "\\.kantcoin\\blockchains\\" + dir + "\\pwd" + "\" )" //depois trocar essa porta 8545 por outra e mudar no init.js

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, args)
	}()

	err = geth_cmd.Start()

	if err == nil {
		log.Println("Geth running")
		fmt.Fprint(w, "complete")
	} else {
		log.Println("Geth not running")
		fmt.Fprint(w, "error")
	}

	//Removing the privatekey and password files
	go func(){
		duration := 12 * time.Second
		time.Sleep(duration)
		os.Remove(getHome() + "/.kantcoin/blockchains/" + dir + "/privkey")
		os.Remove(getHome() + "/.kantcoin/blockchains/" + dir + "/pwd")
	}()

}

//Call geth directly on Linux
func runBlockchainLin(w http.ResponseWriter, id, address, dir string){
	if enodes != ""{
		geth_cmd = exec.Command(executables[runtime.GOOS]["geth_backslash"],"--bootnodesv5", enodes,"--datadir", getHome() + "/.kantcoin/blockchains/" + dir, "--networkid", id, "--lightkdf", "--shh", "--rpc", "--rpcport","8545","--rpcapi", "admin,db,eth,net,web3,shh", "--rpccorsdomain", "*", "--mine", "--minerthreads=1", "--unlock", address, "--etherbase", address, "--password", getHome() + "/.kantcoin/blockchains/" + dir + "/pwd") //trocar porta 8545 por outra
	} else {
		geth_cmd = exec.Command(executables[runtime.GOOS]["geth_backslash"],"--datadir", getHome() + "/.kantcoin/blockchains/" + dir, "--networkid", id, "--lightkdf", "--shh", "--rpc", "--rpcport","8545","--rpcapi", "admin,db,eth,net,web3,shh", "--rpccorsdomain", "*", "--mine", "--minerthreads=1", "--unlock", address, "--etherbase", address, "--password", getHome() + "/.kantcoin/blockchains/" + dir + "/pwd") //trocar porta 8545 por outra
	}

	err := geth_cmd.Start()
	if err == nil {
		log.Println("Geth running")
		fmt.Fprint(w, "complete")
	} else {
		log.Println("Geth not running")
		fmt.Fprint(w, "error")
	}

	//Removing the privatekey and password files
	go func(){
		duration := 12 * time.Second
		time.Sleep(duration)
		os.Remove(getHome() + "/.kantcoin/blockchains/" + dir + "/privkey")
		os.Remove(getHome() + "/.kantcoin/blockchains/" + dir + "/pwd")
	}()

}

//Verify if IPFS has already been initialized
func isIPFSRunning() bool{
	resp, err := http.Get("http://localhost:8080/ipfs/QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG/readme")
	if err == nil && resp.StatusCode == http.StatusOK{
		log.Println("IPFS has already been initialized")
		return true
	}
	return false
}

//Before calling runBlockchainLin/Win, create the password file to unlock the main account
func createPwdFile(path string, w http.ResponseWriter) {
	query := path[len("queryCreatePwdFile="):]
	parts := strings.Split(query, THE_AND)
	dir := parts[0]
	password := parts[1]

	f, err := os.Create(getHome() + "/.kantcoin/blockchains/" + dir + "/pwd")
	if err == nil{
		_,err2 := f.WriteString(password)
		f.Close()
		if err2 == nil{
			fmt.Fprint(w, "complete")
		} else {
			fmt.Fprint(w, "error")
		}
	} else {
		fmt.Fprint(w, "error")
	}
}

//Create new file with a private key (and another with a password) to be imported with the command: geth account import
func createPrivateKeyFile(directory, privkey, password string) bool{
 	f, err := os.Create(getHome() + "/.kantcoin/blockchains/" + directory + "/privkey")
	if err == nil {
		_,err2 := f.WriteString(privkey)
		f.Close()
		if err2 == nil{
			f2, err3 := os.Create(getHome() + "/.kantcoin/blockchains/" + directory + "/pwd")
			if err3 == nil{
				_,err4 := f2.WriteString(password)
				f2.Close()
				if err4 == nil{
					return true
				}
			}
		}
	}
	return false
}

//Insert a new account with the command geth account import. In order to do that, create new privatekey and password file.
func insertAccountIntoBlockchain(path string, w http.ResponseWriter) {
	query := path[len("queryInsertAccountIntoBlockchain="):]
	parts := strings.Split(query, THE_AND)
	dir := parts[0]
	privkey := parts[1]
	password := parts[2]

	if createPrivateKeyFile(dir, privkey, password){
		if runtime.GOOS == "windows" {
			newAccountWin(w, dir)
		} else if runtime.GOOS == "linux" {
			newAccountLin(w, dir)
		}
	}
}

//Call the command: geth account import [WINDOWS]
func newAccountWin(w http.ResponseWriter, dir string){
	cmd := exec.Command("cmd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	filepath,_ := os.Getwd()
	args := "( " + filepath + executables[runtime.GOOS]["geth_backslash"] + " --datadir \"" + getHome() + "\\.kantcoin\\blockchains\\" + dir + "\" account import \"" + getHome() + "\\.kantcoin\\blockchains\\" + dir + "\\privkey\" --password \"" + getHome() + "\\.kantcoin\\blockchains\\" + dir + "\\pwd\" )"

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, args)
	}()

	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	if err == nil {
		log.Println("Account inserted")
		fmt.Fprint(w, "complete")
	} else {
		log.Println("Account not inserted")
		fmt.Fprint(w, "error")
	}
}

//Call the command: geth account import [LINUX]
func newAccountLin(w http.ResponseWriter, dir string) {
	_, err := exec.LookPath(executables[runtime.GOOS]["geth"])
	if err == nil {
		out, err := exec.Command(executables[runtime.GOOS]["geth"], "--datadir", getHome() + "/.kantcoin/blockchains/" + dir, "account", "import", getHome() + "/.kantcoin/blockchains/"+ dir + "/privkey", "--password", getHome() + "/.kantcoin/blockchains/" + dir + "/pwd").Output()
		outstr := string(out)
		if (verbose) {
			log.Println("----------------------------------Command response ----------------------------------")
			log.Println(outstr)
			log.Println("------------------------------End of command response--------------------------------")
		}

		if err == nil {
			log.Println("Account inserted")
			fmt.Fprint(w, "complete")
		} else {
			log.Println("Account not inserted")
			fmt.Fprint(w, "error")
		}
	} else {
		log.Println("Geth is not installed")
		fmt.Fprint(w, "error")
	}

}

//Creating a key with the name received
func addIPNS(path string, w http.ResponseWriter){
	id := path[len("queryAddIPNSKey="):]
	_, err := exec.LookPath(executables[runtime.GOOS]["ipfs"])
	if err == nil {
		if runtime.GOOS == "windows"{
			addIPFSKeyWin(id, w)
		} else if runtime.GOOS == "linux"{
			newIPFSKeyLin(id)
		}
	} else {
		fmt.Fprint(w, "error")
	}
}

//Creating a key (WINDOWS)
func addIPFSKeyWin(dir string, w http.ResponseWriter){
	cmd := exec.Command("cmd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	filepath,_ := os.Getwd()
	args := "( " + filepath + executables[runtime.GOOS]["ipfs_backslash"] + " key gen --type=rsa --size=2048 " + dir + " )"

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, args)
	}()

	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	if err == nil {
		if strings.Index(outstr,"refusing to overwrite") > 0{
			fmt.Fprint(w, "error")
		} else {
			fmt.Fprint(w, "complete")
		}
	} else {
		fmt.Fprint(w, "error")
	}
}

//Creating a key (LINUX)
func addIPFSKeyLin(dir string, w http.ResponseWriter){
	out, err := exec.Command(executables[runtime.GOOS]["ipfs_backslash"],"key","gen", "--type=rsa", "--size=2048" + dir).Output()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	if err == nil {
		if strings.Index(outstr,"refusing to overwrite") > 0{
			fmt.Fprint(w, "error")
		} else {
			fmt.Fprint(w, "complete")
		}
	} else {
		log.Println(outstr)
		fmt.Fprint(w, "error")
	}
}

//Call the "ipfs key list -l" to obtain the ipfn address
func getIPNS(path string, w http.ResponseWriter){
	id := path[len("queryGetIPNS="):]
	_, err := exec.LookPath(executables[runtime.GOOS]["ipfs"])
	if err == nil {
		if runtime.GOOS == "windows"{
			getIPNSWin(w, id)
		} else if runtime.GOOS == "linux"{
			getIPNSLin(w, id)
		}
	} else {
		fmt.Fprint(w, "error")
	}
}

//Call the "ipfs key list -l" to obtain the ipfn address on WINDOWS
func getIPNSWin(w http.ResponseWriter, id string) string{
	resp := ""
	cmd := exec.Command("cmd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	filepath,_ := os.Getwd()
	args := "( " + filepath + executables[runtime.GOOS]["ipfs_backslash"] + " key list -l  )"

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, args)
	}()

	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	if err == nil {
		lines := StringToLines(outstr)
		resp = findIPNS(lines, id)
		if (w != nil){
			fmt.Fprint(w, resp)
		}
	} else {
		if (w != nil){
			fmt.Fprint(w, "error")
		}
	}
	return resp
}

//Call the "ipfs key list -l" to obtain the ipfn address on LINUX
func getIPNSLin(w http.ResponseWriter, id string) string{
	resp := ""
	out, err := exec.Command(executables[runtime.GOOS]["ipfs_backslash"],"key", "list", "-l").Output()
	outstr := string(out)
	if (verbose) {
		log.Println("----------------------------------Command response ----------------------------------")
		log.Println(outstr)
		log.Println("------------------------------End of command response--------------------------------")
	}

	if err == nil{
		lines := StringToLines(outstr)
		resp = findIPNS(lines, id)
		if (w != nil){
			fmt.Fprint(w, resp)
		}
	} else {
		if (w != nil){
			fmt.Fprint(w, "error")
		}
	}
	return resp
}

//Check if the .kantcoin directory was created. If it was not created, then create it
func initDirectory(){
	userprofile := getHome()
	_, err := exec.LookPath(userprofile + "/.kantcoin")
	if err != nil {
		os.Mkdir(userprofile + "/.kantcoin",0700)
		os.Mkdir(userprofile + "/.kantcoin/profiles",0700)
		os.Mkdir(userprofile + "/.kantcoin/blockchains",0700)
	}
}

//Load the executables map
func loadExecutables(){
	add(executables, "windows","ipfs", "ipfs/ipfs.exe")
	add(executables, "windows","ipfs_backslash", "\\ipfs\\ipfs.exe")
	add(executables, "windows","geth", "geth/geth.exe")
	add(executables, "windows","geth_backslash", "\\geth\\geth.exe")
	add(executables, "linux","ipfs", "ipfs/ipfs")
	add(executables, "linux","ipfs_backslash", "ipfs/ipfs")
	add(executables, "linux","geth", "geth/geth")
	add(executables, "linux","geth_backslash", "geth/geth")
}

//Trasforming a string in an array
func StringToLines(s string) []string {
	var lines []string

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

//Look each line to find the ipns of a certain campaign
func findIPNS(lines []string, id string) string{
	for i:=0; i < len(lines); i = i + 1{
		if strings.Index(lines[i], id) > 0{
			return lines[i][:46]
		}
	}
	return "error"
}