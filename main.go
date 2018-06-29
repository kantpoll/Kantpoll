/**
 * Kantcoin Project
 * https://kantcoin.org
 * This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
 * If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

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
	"archive/zip"
	"github.com/Azure/azure-storage-blob-go/2016-05-31/azblob"
	"context"
	"net/url"
	"github.com/Azure/azure-pipeline-go/pipeline"
	"golang.org/x/net/proxy"
	"github.com/afocus/captcha"
	"image/color"
	"image"
	"bytes"
	"image/png"
	"encoding/base64"
	"syscall"
	"github.com/skratchdot/open-golang/open"
	"github.com/asticode/go-astilectron"
	"github.com/lxn/win"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"crypto/cipher"
	"crypto/aes"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
	"github.com/peterbourgon/diskv"
	"encoding/json"
)

/******************** Constants ********************/

const THE_AND = "--and--"
const QUESTION_MARK = "--question--"
const HASHTAG = "--hashtag--"
const QUOTE = "--quote--"
const DOUBLEQUOTE = "--doublequote--"
const BACKSLASH = "--backslash--"
const STORE_ENODE = "StoreEnode"
const STORE_GROUP = "StoreGroup"
const REGISTER_VOTER = "RegisterVoter"
const SEND_VOTE = "SendVote"
const ENTER_GROUP = "EnterGroup"
const CONFIRMATION = "Confirmation"
const TOR2WEB_FORM_HTML = "<html>" +
	"<head><title>Captcha</title><meta http-equiv='Content-Type' content='text/html; charset=UTF-8' /><meta name='viewport' content='width=device-width, initial-scale=1.0' />" +
	"<style>" +
	"h1 {  height: 100px;  width: 100%;  font-size: 18px;  background: #18aa8d;  color: white;  line-height: 150%;  border-radius: 3px 3px 0 0;  box-shadow: 0 2px 5px 1px rgba(0, 0, 0, 0.2);} " +
	"form {  box-sizing: border-box;  width: 280px;  margin: 50px auto 0;  box-shadow: 2px 2px 5px 1px rgba(0, 0, 0, 0.2);  padding-bottom: 40px;  border-radius: 3px;} " +
	"form h1 {  box-sizing: border-box;  padding: 20px;} " +
	"input {  height: 50px; margin: 40px 25px;  width: 220px;  display: block;  border: none;  padding: 10px 0;  border-bottom: solid 1px #1abc9c;  -webkit-transition: all 0.3s cubic-bezier(0.64, 0.09, 0.08, 1); transition: all 0.3s cubic-bezier(0.64, 0.09, 0.08, 1);  background: -webkit-linear-gradient(top, rgba(255, 255, 255, 0) 96%, #1abc9c 4%);  background: linear-gradient(to bottom, rgba(255, 255, 255, 0) 96%, #1abc9c 4%);  background-position: -220px 0;  background-size: 220px 100%;  background-repeat: no-repeat;  color: #0e6252;} " +
	"input:focus, input:valid {  box-shadow: none;  outline: none;  background-position: 0 0;} " +
	"input:focus::-webkit-input-placeholder, input:valid::-webkit-input-placeholder {  color: #1abc9c;  font-size: 11px;  -webkit-transform: translateY(-20px); transform: translateY(-20px);  visibility: visible !important;} " +
	"#data_input {display: none; visibility: hidden;} " +
	"#opener_input {display: none; visibility: hidden;} " +
	"#img_div {text-align: center; width: 100%;}" +
	"button {  border: none;  background: #1abc9c;  cursor: pointer;  border-radius: 3px;  padding: 6px;  width: 220px;  color: white;  margin-left: 25px;  box-shadow: 0 3px 6px 0 rgba(0, 0, 0, 0.2);} " +
	"button:hover {  -webkit-transform: translateY(-3px);  -ms-transform: translateY(-3px);  transform: translateY(-3px);  box-shadow: 0 6px 6px 0 rgba(0, 0, 0, 0.2);} " +
	"input:focus::-webkit-input-placeholder, input:valid::-webkit-input-placeholder {  color: #1abc9c;  font-size: 11px;  -webkit-transform: translateY(-20px);  transform: translateY(-20px);  visibility: visible !important;} " +
	"</style>" +
	"<script type='text/javascript'>" +
	"var locale = ''; " +
	"sessionStorage.setItem('opener', '[[opener]]'); " +
	"var send_text = 'SEND'; " +
	"var form_title_text = 'Complete the captcha before interacting with the campaign'; " +
	"if (navigator.language){ locale = navigator.language.substring(0,2).toLowerCase();} " +
	"if (locale == 'pt'){send_text = 'ENVIAR'; form_title_text = 'Complete o captcha antes de interagir com a campanha';} else " +
	"if (locale == 'es'){send_text = 'ENVIAR'; form_title_text = 'Completa el captcha antes de interactuar con la campaña';} else " +
	"if (locale == 'fr'){send_text = 'ENVOYER'; form_title_text = \"Complétez le captcha avant d'interagir avec la campagne\";} " +
	"window.addEventListener('message', function (ev){if ('[[opener]]'.startsWith(ev.origin)){ var json_data = JSON.parse(ev.data); data_input.value = json_data.data;}}, false); " +
	"window.addEventListener('load', function(ev){window.opener.postMessage('load', '[[opener]]'); send_button.innerHTML = send_text; form_title.innerHTML = form_title_text; opener_input.value = '[[opener]]';}); " +
	"</script></head>" +
	"<body>" +
	"<form action='[[query_type]]' method='post'>" +
	"<h1 id='form_title'>Complete the captcha before interacting with the campaign</h1>" +
	"<input id='data_input' type='text' name='data' value='[[data_value]]' />" +
	"<div id='img_div'><img src='[[cap_uri]]' /></div>" +
	"<input placeholder='Captcha' type='text' name='captcha' required='' autocomplete='off' /><br>" +
	"<input id='opener_input' type='text' name='opener' /><br>" +
	"<button id='send_button'>SEND</button>" +
	"</form>" +
	"</body></html>"

const TOR2WEB_FORM_HTML_LOCAL = "<html>" +
	"<head><title>Captcha</title><meta http-equiv='Content-Type' content='text/html; charset=UTF-8' /><meta name='viewport' content='width=device-width, initial-scale=1.0' />" +
	"<style>" +
	"h1 {  height: 100px;  width: 100%;  font-size: 18px;  background: #18aa8d;  color: white;  line-height: 150%;  border-radius: 3px 3px 0 0;  box-shadow: 0 2px 5px 1px rgba(0, 0, 0, 0.2);} " +
	"form {  box-sizing: border-box;  width: 280px;  margin: 50px auto 0;  box-shadow: 2px 2px 5px 1px rgba(0, 0, 0, 0.2);  padding-bottom: 40px;  border-radius: 3px;} " +
	"form h1 {  box-sizing: border-box;  padding: 20px;} " +
	"input {  height: 50px; margin: 40px 25px;  width: 220px;  display: block;  border: none;  padding: 10px 0;  border-bottom: solid 1px #1abc9c;  -webkit-transition: all 0.3s cubic-bezier(0.64, 0.09, 0.08, 1); transition: all 0.3s cubic-bezier(0.64, 0.09, 0.08, 1);  background: -webkit-linear-gradient(top, rgba(255, 255, 255, 0) 96%, #1abc9c 4%);  background: linear-gradient(to bottom, rgba(255, 255, 255, 0) 96%, #1abc9c 4%);  background-position: -220px 0;  background-size: 220px 100%;  background-repeat: no-repeat;  color: #0e6252;} " +
	"input:focus, input:valid {  box-shadow: none;  outline: none;  background-position: 0 0;} " +
	"input:focus::-webkit-input-placeholder, input:valid::-webkit-input-placeholder {  color: #1abc9c;  font-size: 11px;  -webkit-transform: translateY(-20px); transform: translateY(-20px);  visibility: visible !important;} " +
	"#data_input {display: none; visibility: hidden;} " +
	"#opener_input {display: none; visibility: hidden;} " +
	"#img_div {text-align: center; width: 100%;}" +
	"button {  border: none;  background: #1abc9c;  cursor: pointer;  border-radius: 3px;  padding: 6px;  width: 220px;  color: white;  margin-left: 25px;  box-shadow: 0 3px 6px 0 rgba(0, 0, 0, 0.2);} " +
	"button:hover {  -webkit-transform: translateY(-3px);  -ms-transform: translateY(-3px);  transform: translateY(-3px);  box-shadow: 0 6px 6px 0 rgba(0, 0, 0, 0.2);} " +
	"input:focus::-webkit-input-placeholder, input:valid::-webkit-input-placeholder {  color: #1abc9c;  font-size: 11px;  -webkit-transform: translateY(-20px);  transform: translateY(-20px);  visibility: visible !important;} " +
	"</style>" +
	"<script type='text/javascript'>" +
	"var locale = ''; " +
	"sessionStorage.setItem('opener', '[[opener]]'); " +
	"var send_text = 'SEND'; " +
	"var change_port = '[[opener]]'.split('1985')[0] + '1988'; " +
	"var form_title_text = 'Complete the captcha before interacting with the campaign'; " +
	"if (navigator.language){ locale = navigator.language.substring(0,2).toLowerCase();} " +
	"if (locale == 'pt'){send_text = 'ENVIAR'; form_title_text = 'Complete o captcha antes de interagir com a campanha';} else " +
	"if (locale == 'es'){send_text = 'ENVIAR'; form_title_text = 'Completa el captcha antes de interactuar con la campaña';} else " +
	"if (locale == 'fr'){send_text = 'ENVOYER'; form_title_text = \"Complétez le captcha avant d'interagir avec la campagne\";} " +
	"window.addEventListener('load', function(ev){send_button.innerHTML = send_text; form_title.innerHTML = form_title_text; the_form.addEventListener('submit', sendURL)}); " +
	"function sendURL(ev){ev.preventDefault(); var request = new XMLHttpRequest(); request.open('GET', change_port + '/querySendTorRequest=[[onion_plus_query_type]]data=' + data_input.value + '&captcha=' + captcha_input.value + '&opener=[[opener]]', true); request.addEventListener('load', showResponse); request.send()}" +
	"function showResponse(){document.body.innerHTML = this.responseText;}" +
	"</script></head>" +
	"<body>" +
	"<form id='the_form'>" +
	"<h1 id='form_title'>Complete the captcha before interacting with the campaign</h1>" +
	"<input id='data_input' type='text' name='data' value='[[data_value]]' />" +
	"<div id='img_div'><img src='[[cap_uri]]' /></div>" +
	"<input id='captcha_input' placeholder='Captcha' type='text' name='captcha' required='' autocomplete='off' /><br>" +
	"<input id='opener_input' type='text' name='opener' /><br>" +
	"<button id='send_button'>SEND</button>" +
	"</form>" +
	"</body></html>"

const DONE_HTML = "<html><head><title>Captcha</title><meta http-equiv='Content-Type' content='text/html; charset=UTF-8' /><meta name='viewport' content='width=device-width, initial-scale=1.0' />" +
	"<script type='text/javascript'>if (!'[[opener]]'.startsWith('[[open')){window.addEventListener('load', function(ev){window.opener.postMessage('done', '[[opener]]');});} </script></head>" +
	"<body><div style='font-family:\"Courier New\",monospace;text-align:center;font-size: 14px;'><br><br>════════════════════════════════<br>Done - Hecho - Pronto - Terminé.<br>════════════════════════════════<br>Now, close the window.<br>Ahora, cierra la ventana.<br>Agora, feche a janela.<br>Maintenant, fermez la fenêtre.<br>════════════════════════════════<br></div>" +
	"</body></html>"

/******************** Structs ********************/

type CampaignIPNS struct{
	ipns string
	id string
}

//Struct to receive the installation links from an http response
type InstallationLinks struct{
	UGW string `json:"url_geth_windows"`
	UGL string `json:"url_geth_linux"`
	UIW string `json:"url_ipfs_windows"`
	UIL string `json:"url_ipfs_linux"`
	UTW string `json:"url_tor_windows"`
	GF string `json:"geth_folder"`
	IF string `json:"ipfs_folder"`
}

/******************** Global variables ********************/

// Installation links
var url_geth_windows = "https://gethstore.blob.core.windows.net/builds/geth-windows-amd64-1.8.11-dea1ce05.zip"
var url_geth_linux = "https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.8.11-dea1ce05.tar.gz"
var url_ipfs_windows = "https://dist.ipfs.io/go-ipfs/v0.4.15/go-ipfs_v0.4.15_windows-amd64.zip"
var url_ipfs_linux = "https://dist.ipfs.io/go-ipfs/v0.4.15/go-ipfs_v0.4.15_linux-amd64.tar.gz"
var url_tor_windows = ""//"https://www.torproject.org/dist/torbrowser/7.5.6/tor-win32-0.3.3.7.zip"

//Geth and IPFS folders
var geth_folder = "geth-windows-amd64-1.8.11-dea1ce05"
var ipfs_folder = "go-ipfs"

//Flags
var (
	verbose bool
)

//Set the commands according to the os
var executables map [string] map [string] string = make(map [string] map [string] string)

//In order to know if an instance of geth is already running
//And to kill these commands on exit
var gethCmd, ipfsCmd, torCmd *exec.Cmd
//It receives the "exit" argument
var stdinGeth io.WriteCloser

//How many peers to connect
var how_many_enodes = "20"

//The last chainid, which is currently active
var lastId = 0

//Sliced messages received via Tor
var (
	registerVoterMessages []string
	sendVoteMessages      []string
	storeEnodeMessages    []string
	storeGroupMessages    []string
	enterGroupMessages    []string
	confirmMessages       []string
)

//Send requests via Tor
var torClient *http.Client

//The captcha generator
var theCaptcha *captcha.Captcha
//A map with half of the captcha as a key and half as a value
var capMap map[string] string

//The ipns of the current campaign
var currentIpns CampaignIPNS

//Logs to be displayed on browser
var logs string

//The astilectron gui
var (
	aelectron *astilectron.Astilectron
	window_electron *astilectron.Window
 	tray_electron *astilectron.Tray
	menu_tray_electron *astilectron.Menu
)

//It shows the progress of the installation
var (
	installerWindow *walk.MainWindow
	installerTE *walk.TextEdit
	installerLogs bool
)

//Info about the campaign and groups are cached in these variables
var (
	campaignInfo string
	campaignIPFSInfo string
	groupsInfo []string
	groupsMessages []string
)

//The language that the user selected on the system tray menu
var chosen_language = "en"

//Each voter has a secret in order to receive the group vote message in an encrypted message
var secretsDatabase *diskv.Diskv
//A database of voters and their group indexes
var votersGroupsDatabase *diskv.Diskv

//When exiting
var stopRecevingRequests bool

/******************** Functions ********************/

func main() {
	//Loading the flags
	verboseFlag := flag.Bool("verbose", false, "Activate the verbose mode")
	flag.Parse()
	verbose = *verboseFlag

	//The map with executable paths
	loadExecutables()

	//Loading voter's secrets
	initDatabases()

	//Checking if the .kantcoin subdirectories were created. If it was not created, then create it
	initDirectory()

	//Creating other necessary directories if they do not exist
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp",0770)
	}

	//Checking if these directories exist
	_, err_geth := os.Stat("geth")
	_, err_ipfs := os.Stat("ipfs")
	_, err_tor := os.Stat("tor")
	_, err_electron := os.Stat("vendor")
	_, err_website := os.Stat("website")

	if os.IsNotExist(err_website) {
		if _, err := os.Stat("website.zip"); err == nil {
			_, err2 := unzip("website.zip", ".")
			if err2 != nil{
				Println("Error unzipping website.zip")
			}
		} else {
			Println("File website.zip not found")
		}
	}

	if os.IsNotExist(err_geth) || os.IsNotExist(err_ipfs) || os.IsNotExist(err_tor) || os.IsNotExist(err_electron){
		getInstallationLinks()

		//Preparing the Installer Window
		screen_width := int(win.GetSystemMetrics(win.SM_CXSCREEN))
		screen_height := int(win.GetSystemMetrics(win.SM_CYSCREEN))
		width := 830
		height := 450
		icon, _ := walk.NewIconFromFile("website/imgs/logo.ico")
		MainWindow{
			AssignTo: &installerWindow,
			Title:   "Completing the installation",
			MinSize: Size{width, height},
			Font: Font{Family: "Courier"},
			Layout:  VBox{},
			Icon: icon,
			Children: []Widget{
				TextEdit{AssignTo: &installerTE, ReadOnly: true, VScroll:true},
			},
		}.Create()

		installerWindow.SetX((screen_width - width)/2)
		installerWindow.SetY((screen_height - height)/2)

		go func(){
			duration := 3 * time.Second
			time.Sleep(duration)

			installerLogs = true

			Println("╔--------------------------------------------------------╗")
			Println("| Installation - Instalación - Instalação - Installation |")
			Println("╚--------------------------------------------------------╝")

			if os.IsNotExist(err_geth) {
				os.Mkdir("geth",0770)
				Println("╔---------------------------------------------------╗")
				Println("| Installing Geth. Wait a few minutes.              |")
				Println("| Instalando el Geth. Espera unos minutos.          | ")
				Println("| Instalando o Geth. Aguarde alguns minutos.        | ")
				Println("| Installation de Geth. Attends quelques minutes.   | ")
				Println("╚---------------------------------------------------╝")
				installGeth()
			}
			if os.IsNotExist(err_ipfs) {
				os.Mkdir("ipfs",0770)
				Println("╔---------------------------------------------------╗")
				Println("| Installing IPFS. Wait a few minutes.              | ")
				Println("| Instalando el IPFS. Espera unos minutos.          | ")
				Println("| Instalando o IPFS. Aguarde alguns minutos.        | ")
				Println("| Installation de IPFS. Attends quelques minutes.   | ")
				Println("╚---------------------------------------------------╝")
				installIPFS()
			}
			if os.IsNotExist(err_tor) {
				os.Mkdir("tor",0770)
				Println("╔--------------------------------------------╗")
				Println("| Installing Tor. Wait a moment.             | ")
				Println("| Instalando el Tor. Espera un momento.      | ")
				Println("| Instalando o Tor. Aguarde um momento.      | ")
				Println("| Installation de Tor. Attendez un moment.   | ")
				Println("╚--------------------------------------------╝")
				installTor()
			}

			if os.IsNotExist(err_electron) {
				Println("╔-------------------------------------------------╗")
				Println("| Installing Electron. Wait a moment.             | ")
				Println("| Instalando el Electron. Espera un momento.      | ")
				Println("| Instalando o Electron. Aguarde um momento.      | ")
				Println("| Installation de Electron. Attendez un moment.   | ")
				Println("╚-------------------------------------------------╝")

				aelectron, _ = astilectron.New(astilectron.Options{
					AppName: "Kantcoin",
					AppIconDefaultPath: "website/imgs/icon32_32.png",
					//AppIconDarwinPath:  "icon",
					BaseDirectoryPath: ".",
				})

				aelectron.On(astilectron.EventNameAppCrash, func(e astilectron.Event) (deleteListener bool) {
					Println("Astilectron has crashed")
					return
				})

				aelectron.On(astilectron.EventNameAppErrorAccept, func(e astilectron.Event) (deleteListener bool) {
					Println("ErrorAccept event on Astilectron")
					return
				})

				// Start astilectron
				aelectron.Start()
			}

			Println("╔------------------------╗")
			Println("| Installation complete. |")
			Println("| Instalación completa.  |")
			Println("| Instalação completa.   |")
			Println("| Installation complète. |")
			Println("╚------------------------╝")

			duration = 2 * time.Second
			time.Sleep(duration)
			installerWindow.Close()

			installerLogs = false
		}()

		installerWindow.Run()
	}

	Println("Logs (only in english):")

	//Setting the captcha
	theCaptcha = captcha.New()
	capMap = make(map[string]string)
	fontContenrs, err := ioutil.ReadFile("website/fonts/comic.ttf")
	if err != nil {
		Println(err.Error())
	}
	err = theCaptcha.AddFontFromBytes(fontContenrs)
	if err != nil {
		Println(err.Error())
	}
	theCaptcha.SetSize(110, 40)
	theCaptcha.SetDisturbance(captcha.HIGH)
	theCaptcha.SetFrontColor(color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	theCaptcha.SetBkgColor(color.RGBA{255, 255, 255, 255})

	//Handler to interact with the IPFS and Geth
	changeHeaderThenServe := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			//When exiting
			if (stopRecevingRequests){
				return
			}
			//Adding Content Security Policy and CORS
			w.Header().Add("Content-Security-Policy", "script-src http://localhost:1985 http://127.0.0.1:1985 ; style-src http://localhost:1985 http://127.0.0.1:1985 'unsafe-inline'; child-src 'none'; object-src 'none'; form-action http://localhost:1985 http://127.0.0.1:1985 ; connect-src http: https: data: ; worker-src http://localhost:1985 http://127.0.0.1:1985 ")
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:1985 http://127.0.0.1:1985")

			path := r.URL.Path[1:]
			if strings.Index(path, "query") == 0 {
				if strings.Index(path, "queryVerbose=") == 0{
					choice := path[len("queryVerbose="):]
					if choice == "true"{
						verbose = true
					} else if choice == "false"{
						verbose = false
					}
					return
				} else if strings.Index(path, "queryMyIP=") == 0{
					whatIsMyIP(w)
					return
				} else if strings.Index(path, "queryEnode=") == 0{
					whatIsMyEnode(path, w)
					return
				} else if strings.Index(path, "queryAddProfile=") == 0{
					addProfile(path, w)
					return
				} else if strings.Index(path, "queryAddPeer=") == 0{
					addPeer(path, w)
					return
				} else if strings.Index(path, "queryProfileExists=") == 0{
					profileExists(path, w)
					return
				} else if strings.Index(path, "queryGetProfile=") == 0{
					getProfile(path, w)
					return
				} else if strings.Index(path, "querySetBlockchain=") == 0{
					setBlockchain(path, w)
					return
				} else if strings.Index(path, "queryGetIPNS=") == 0{
					getIPNS(path, w)
					return
				} else if strings.Index(path, "queryAddIPNSKey=") == 0{
					addIPNS(path, w)
					return
				} else if strings.Index(path, "queryVerifyBlockchain=") == 0{
					verifyBlockchain(path, w)
					return
				} else if strings.Index(path, "queryRunBlockchain=") == 0{
					runBlockchain(path, w)
					return
				} else if strings.Index(path, "queryInsertAccountIntoBlockchain=") == 0{
					insertAccountIntoBlockchain(path, w)
					return
				} else if strings.Index(path, "queryCheckUser=") == 0{
					checkUser(path, w)
					return
				} else if strings.Index(path, "queryCreatePwdFile=") == 0{
					createPwdFile(path, w)
					return
				} else if strings.Index(path, "queryGetHiddenServiceHostname=") == 0{
					getHiddenServiceHostname(w)
					return
				} else if strings.Index(path, "queryGetRegisterVoterMessages=") == 0{
					fmt.Fprint(w, registerVoterMessages)
					registerVoterMessages = registerVoterMessages[:0]
					return
				} else if strings.Index(path, "queryGetSendVoteMessages=") == 0{
					fmt.Fprint(w, sendVoteMessages)
					sendVoteMessages = sendVoteMessages[:0]
					return
				} else if strings.Index(path, "queryGetStoreEnodeMessages=") == 0{
					fmt.Fprint(w, storeEnodeMessages)
					storeEnodeMessages = storeEnodeMessages[:0]
					return
				} else if strings.Index(path, "queryGetStoreGroupMessages=") == 0{
					fmt.Fprint(w, storeGroupMessages)
					storeGroupMessages = storeGroupMessages[:0]
					return
				} else if strings.Index(path, "queryGetEnterGroupMessages=") == 0{
					fmt.Fprint(w, enterGroupMessages)
					enterGroupMessages = enterGroupMessages[:0]
					return
				} else if strings.Index(path, "queryGetConfirmMessages=") == 0{
					fmt.Fprint(w, confirmMessages)
					confirmMessages = confirmMessages[:0]
					return
				} else if strings.Index(path, "queryStoreGroupInfo=") == 0{
					storeGroupInfo(path, w)
					return
				} else if strings.Index(path, "queryStoreGroupMessage=") == 0{
					storeGroupMessage(path, w)
					return
				} else if strings.Index(path, "queryStoreVoter=") == 0{
					storeVoter(path, w)
					return
				} else if strings.Index(path, "queryStoreVoterSecret=") == 0{
					storeVoterSecret(path, w)
					return
				} else if strings.Index(path, "queryStoreCampaignInfo=") == 0{
					storeCampaignInfo(path, w)
					return
				} else if strings.Index(path, "queryStoreCampaignIPFSInfo=") == 0{
					storeCampaignIPFSInfo(path, w)
					return
				} else if strings.Index(path, "queryGetChosenLanguage=") == 0{
					fmt.Fprint(w, chosen_language)
					return
				} else if strings.Index(path, "queryCleanVariables=") == 0{
					campaignInfo = campaignInfo[:0]
					campaignIPFSInfo = campaignIPFSInfo[:0]
					groupsInfo = nil
					votersGroupsDatabase.EraseAll()
					secretsDatabase.EraseAll()
					groupsMessages = groupsMessages[:0]
					return
				}
			}

			// Serve with the actual handler
			h.ServeHTTP(w, r)
		}
	}

	//Handler to interact with Tor requests
	torHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			//When exiting
			if (stopRecevingRequests){
				return
			}

			//Adding Content Security Policy and CORS
			w.Header().Add("Content-Security-Policy", "script-src 'unsafe-inline' ;")
			w.Header().Add("Access-Control-Allow-Origin", "*")

			path := r.URL.Path[1:]
			if strings.Index(path, "query") == 0 {
				if strings.Index(path, "queryStoreEnodeHTML=") == 0{
					generateHTML(path, w, "/queryStoreEnodePost=")

				} else if strings.Index(path, "queryStoreGroupHTML=") == 0{
					generateHTML(path, w, "/queryStoreGroupPost=")

				} else if strings.Index(path, "queryRegisterVoterHTML=") == 0{
					generateHTML(path, w, "/queryRegisterVoterPost=")

				} else if strings.Index(path, "querySendVoteHTML=") == 0{
					generateHTML(path, w, "/querySendVotePost=")

				} else if strings.Index(path, "queryEnterGroupHTML=") == 0{
					generateHTML(path, w, "/queryEnterGroupPost=")

				} else if strings.Index(path, "queryConfirmHTML=") == 0{
					generateHTML(path, w, "/queryConfirmPost=")

				} else if strings.Index(path, "queryStoreEnodeHTMLLocal=") == 0{
					generateHTMLLocal(path, w, "/queryStoreEnodeGet=")

				} else if strings.Index(path, "queryStoreGroupHTMLLocal=") == 0{
					generateHTMLLocal(path, w, "/queryStoreGroupGet=")

				} else if strings.Index(path, "queryRegisterVoterHTMLLocal=") == 0{
					generateHTMLLocal(path, w, "/queryRegisterVoterGet=")

				} else if strings.Index(path, "querySendVoteHTMLLocal=") == 0{
					generateHTMLLocal(path, w, "/querySendVoteGet=")

				} else if strings.Index(path, "queryEnterGroupHTMLLocal=") == 0{
					generateHTMLLocal(path, w, "/queryEnterGroupGet=")

				} else if strings.Index(path, "queryConfirmHTMLLocal=") == 0{
					generateHTMLLocal(path, w, "/queryConfirmGet=")

				} else if strings.Index(path, "queryStoreEnodePost=") == 0{
					receivePostMessage(w, r, STORE_ENODE)

				} else if strings.Index(path, "queryStoreGroupPost=") == 0{
					receivePostMessage(w, r, STORE_GROUP)

				} else if strings.Index(path, "queryRegisterVoterPost=") == 0{
					receivePostMessage(w, r, REGISTER_VOTER)

				} else if strings.Index(path, "querySendVotePost=") == 0{
					receivePostMessage(w, r, SEND_VOTE)

				} else if strings.Index(path, "queryEnterGroupPost=") == 0{
					receivePostMessage(w, r, ENTER_GROUP)

				} else if strings.Index(path, "queryConfirmPost=") == 0{
					receivePostMessage(w, r, CONFIRMATION)

				} else if strings.Index(path, "queryStoreEnodeGet=") == 0{
					receiveGetMessage(path, w, STORE_ENODE)

				} else if strings.Index(path, "queryStoreGroupGet=") == 0{
					receiveGetMessage(path, w, STORE_GROUP)

				} else if strings.Index(path, "queryRegisterVoterGet=") == 0{
					receiveGetMessage(path, w, REGISTER_VOTER)

				} else if strings.Index(path, "querySendVoteGet=") == 0{
					receiveGetMessage(path, w, SEND_VOTE)

				} else if strings.Index(path, "queryEnterGroupGet=") == 0{
					receiveGetMessage(path, w, ENTER_GROUP)

				} else if strings.Index(path, "queryConfirmGet=") == 0{
					receiveGetMessage(path, w, CONFIRMATION)

				} else if strings.Index(path, "queryGetCampaignInfo=") == 0{
					getCampaignInfo(w)

				} else if strings.Index(path, "queryGetCampaignIPFSInfo=") == 0{
					getCampaignIPFSInfo(w)

				} else if strings.Index(path, "queryGetGroupInfo=") == 0{
					getGroupInfo(path, w)

				} else if strings.Index(path, "queryMyGroupIndex=") == 0{
					myGroupIndex(path, w)

				} else if strings.Index(path, "queryGetVoteMessage=") == 0{
					getVoteMessage(path, w)

				} else if strings.Index(path, "querySendTorRequest=") == 0{
					//This query should be here, since we want to avoid scripts from onion addresses accessing variables from localhost:1985
					sendTorRequest(path, w)
				}
			}
		}
	}

	//Initializing IPFS
	if !isIPFSRunning() {
		initIPFS(0)
	}

	//Initializing Tor to provide the hidden service
	initTor()

	//Creating a client to send requests via Tor
	initTorClient()

	go func(){
		//Listen requests made through the Tor network
		http.ListenAndServe(":1988", torHandler())
	}()

	go func(){
		//It could be already started during the installation
		if aelectron == nil{
			aelectron, _ = astilectron.New(astilectron.Options{
				AppName: "Kantcoin",
				AppIconDefaultPath: "website/imgs/icon32_32.png",
				//AppIconDarwinPath:  "icon",
				BaseDirectoryPath: ".",
			})

			// Start astilectron
			aelectron.Start()
		}

		//Opening the gui for the user
		width := int(win.GetSystemMetrics(win.SM_CXSCREEN))
		height := int(win.GetSystemMetrics(win.SM_CYSCREEN))
		window_electron, _ = aelectron.NewWindow("http://localhost:1985/home.html", &astilectron.WindowOptions{
			Center: astilectron.PtrBool(true),
			Height: astilectron.PtrInt(height),
			Width:  astilectron.PtrInt(width),
			Icon: astilectron.PtrStr("website/imgs/icon32_32.png"),
			MessageBoxOnClose: &astilectron.MessageBoxOptions{
				Message: "This will shut down the server as well.\nEsto apagará el servidor también.\nIsto desligará o servidor também.\nCela fermera le serveur aussi.",
				Title: "Confirm",
				Type: astilectron.MessageBoxTypeWarning,
				Buttons: []string{"OK", "Cancel"},
				ConfirmID: astilectron.PtrInt(0),
				CancelID: astilectron.PtrInt(1),
				},
			WebPreferences: &astilectron.WebPreferences{
				AllowRunningInsecureContent: astilectron.PtrBool(false),
				NodeIntegration: astilectron.PtrBool(false),
				NodeIntegrationInWorker: astilectron.PtrBool(false),
				ExperimentalFeatures: astilectron.PtrBool(false),
	 			ExperimentalCanvasFeatures: astilectron.PtrBool(false),
				ContextIsolation: astilectron.PtrBool(true),
				WebSecurity: astilectron.PtrBool(true),
				},
		})
		window_electron.Create()
		window_electron.Maximize()

		//Placing an icon in the notification area
		//systray.Run(onReady, onExit)
		tray_electron = aelectron.NewTray(&astilectron.TrayOptions{
			Image:   astilectron.PtrStr("website/imgs/icon32_32.png"),
			Tooltip: astilectron.PtrStr("Kantcoin"),
		})

		// Create tray
		tray_electron.Create()

		// New tray menu
		menu_tray_electron = tray_electron.NewMenu([]*astilectron.MenuItemOptions{
			{
				Label: astilectron.PtrStr("Restart"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					width := int(win.GetSystemMetrics(win.SM_CXSCREEN))
					height := int(win.GetSystemMetrics(win.SM_CYSCREEN))
					window_electron_aux, _ := aelectron.NewWindow("http://localhost:1985/home.html", &astilectron.WindowOptions{
						Center: astilectron.PtrBool(true),
						Height: astilectron.PtrInt(height),
						Width:  astilectron.PtrInt(width),
						Icon: astilectron.PtrStr("website/imgs/icon32_32.png"),
						MessageBoxOnClose: &astilectron.MessageBoxOptions{
							Message: "This will shut down the server as well.\nEsto apagará el servidor también.\nIsto desligará o servidor também.\nCela fermera le serveur aussi.",
							Title: "Confirm",
							Type: astilectron.MessageBoxTypeWarning,
							Buttons: []string{"OK", "Cancel"},
							ConfirmID: astilectron.PtrInt(0),
							CancelID: astilectron.PtrInt(1),
						},
						WebPreferences: &astilectron.WebPreferences{
							AllowRunningInsecureContent: astilectron.PtrBool(false),
							NodeIntegration: astilectron.PtrBool(false),
							NodeIntegrationInWorker: astilectron.PtrBool(false),
							ExperimentalFeatures: astilectron.PtrBool(false),
							ExperimentalCanvasFeatures: astilectron.PtrBool(false),
							ContextIsolation: astilectron.PtrBool(true),
							WebSecurity: astilectron.PtrBool(true),
						},
					})
					window_electron_aux.Create()
					window_electron_aux.Maximize()
					window_electron.Destroy()
					window_electron = window_electron_aux
					return
				},
			},
			{
				Label: astilectron.PtrStr("Language"),
				SubMenu: []*astilectron.MenuItemOptions{
					{Checked: astilectron.PtrBool(true), Label: astilectron.PtrStr("English"), Type: astilectron.MenuItemTypeRadio,
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							chosen_language = "en"
							return
						},
					},
					{Label: astilectron.PtrStr("Español"), Type: astilectron.MenuItemTypeRadio,
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							chosen_language = "es"
							return
						},
					},
					{Label: astilectron.PtrStr("Português"), Type: astilectron.MenuItemTypeRadio,
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							chosen_language = "pt"
							return
						},
					},
					{Label: astilectron.PtrStr("Français"), Type: astilectron.MenuItemTypeRadio,
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							chosen_language = "fr"
							return
						},
					},
				},
			},
			{
				Label: astilectron.PtrStr("Console"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					window_electron.OpenDevTools()
					return
				},
			},
			{
				Label: astilectron.PtrStr("Files"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					open.Run(getHome() + "/.kantcoin")
					return
				},
			},
			{
				Label: astilectron.PtrStr("Logs"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					err := ioutil.WriteFile(getHome() + "/.kantcoin/logs.txt", []byte(logs), 0700)
					if err == nil {
						open.Run(getHome() + "/.kantcoin/logs.txt")
					}
					return
				},
			},
			{
				Label: astilectron.PtrStr("Help"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					open.Run("https://sourceforge.net/p/kantcoin/wiki/Help/")
					return
				},
			},
			{
				Label: astilectron.PtrStr("Version"),
				SubMenu: []*astilectron.MenuItemOptions{
					{Label: astilectron.PtrStr("v0.2.1"), Type: astilectron.MenuItemTypeNormal,},
				},
			},
			{
				Label: astilectron.PtrStr("Exit"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					exit()
					return
				},
			},
		})

		// Create the menu
		menu_tray_electron.Create()

		menu_tray_electron.On(astilectron.EventNameTrayEventClicked, func(e astilectron.Event) (deleteListener bool) {
			window_electron.Focus()
			return
		})

		defer exit()

		//Waiting for events
		aelectron.Wait()
	}()

	//It provides the web pages and services
	http.ListenAndServe("localhost:1985", changeHeaderThenServe(http.FileServer(http.Dir("./website"))))
}

/*
 * This database stores voters' secrets, which are used to receive group vote messages
 */
func initDatabases() {
	// Simplest transform function: put all the data files into the base dir.
	flatTransform := func(s string) []string { return []string{} }

	// Initialize a new diskv store, rooted at "secretsdb", with a 20MB cache.
	secretsDatabase = diskv.New(diskv.Options{
		BasePath:     "secretsdb",
		Transform:    flatTransform,
		CacheSizeMax: 20 * 1024 * 1024,
	})

	// Initialize a new diskv store, rooted at "votersdb", with a 20MB cache.
	votersGroupsDatabase = diskv.New(diskv.Options{
		BasePath:     "votersdb",
		Transform:    flatTransform,
		CacheSizeMax: 20 * 1024 * 1024,
	})
}

/**
 * It generates the HTML page to be displayed to voters when they are not using a local server
 */
func generateHTML(path string, w http.ResponseWriter, queryType string){
	if !verifyParams(path, "opener"){
		fmt.Fprint(w, "error")
		return
	}

	opener := path[strings.LastIndex(path, "opener=") + len ("opener="):]
	_, capImg := genCaptcha()
	capUri := getCaptchaURI(capImg)
	html := strings.Replace(TOR2WEB_FORM_HTML, "[[cap_uri]]", capUri, -1)
	html = strings.Replace(html, "[[opener]]", opener, -1)
	html = strings.Replace(html, "[[query_type]]", queryType, -1)

	fmt.Fprint(w, html)
}

/**
 * It generates the HTML page to be displayed to voters when they are using a local server
 */
func generateHTMLLocal(path string, w http.ResponseWriter, queryType string){
	if !verifyParams(path, "data=", "&onion_address=", "&opener="){
		fmt.Fprint(w, "error")
		return
	}

	data := path[strings.LastIndex(path, "data=") + len("data="): strings.LastIndex(path, "&onion_address=")]
	onion := path[strings.LastIndex(path,"&onion_address=") + len("&onion_address="): strings.LastIndex(path,"&opener=")]
	opener := path[strings.LastIndex(path, "&opener=") + len ("&opener="):]
	_, capImg := genCaptcha()
	capUri := getCaptchaURI(capImg)
	html := strings.Replace(TOR2WEB_FORM_HTML_LOCAL, "[[cap_uri]]", capUri, -1)
	html = strings.Replace(html, "[[opener]]", opener, -1)
	html = strings.Replace(html, "[[onion_plus_query_type]]", onion + queryType, -1)
	html = strings.Replace(html, "[[data_value]]", data, -1)

	fmt.Fprint(w, html)
}

/**
 * It generates a captcha image and string, and stores this string in a map
 */
func genCaptcha()(string, image.Image){
	img, str := theCaptcha.Create(6, captcha.NUM)
	capMap[str[0:3]] = str[3:6]

	return str, img
}

/**
 * It returns an URI representing some image
 */
func getCaptchaURI(img image.Image) string{
	out := new(bytes.Buffer)
	err := png.Encode(out, img)
	base64Img := ""

	if err != nil {
		Println("Can't encode captcha image")
		return ""
	} else {
		base64Img = base64.StdEncoding.EncodeToString(out.Bytes())
	}

	return "data:image/png;base64," + base64Img
}

/**
 * Verifying if the captcha string provided is present in the captcha map
 */
func verifyCaptcha(str string) bool{
	if len(str) != 6{
		return false
	}
	if val, ok := capMap[str[0:3]]; ok {
		if val == str[3:6]{
			delete(capMap, str[0:3])
			return true
		}
	}
	return false
}

/**
 * This function receives GET messages from Tor2web pages
 */
func receiveGetMessage(path string, w http.ResponseWriter, messageType string){
	if !verifyParams(path, "data=", "&captcha=", "&opener="){
		fmt.Fprint(w, "error")
		return
	}

	queryType := ""
	capStr := path[strings.LastIndex(path, "&captcha=") + len("&captcha="): strings.LastIndex(path, "&opener=")]
	if verifyCaptcha(capStr){
		if messageType == STORE_ENODE{
			queryType = "/queryStoreEnodeGet="
			insertStoreEnodeMessage(path, w)
		} else if messageType == STORE_GROUP{
			queryType = "/queryStoreGroupGet="
			insertStoreGroupMessage(path, w)
		} else if messageType == REGISTER_VOTER{
			queryType = "/queryRegisterVoterGet="
			insertRegisterVoterMessage(path, w)
		} else if messageType == SEND_VOTE{
			queryType = "/querySendVoteGet="
			insertSendVoteMessage(path, w)
		} else if messageType == ENTER_GROUP{
			queryType = "/queryEnterGroupGet="
			insertEnterGroupMessage(path, w)
		} else if messageType == CONFIRMATION{
			queryType = "/queryConfirmGet="
			insertConfirmMessage(path, w)
		}
	} else {
		opener := path[strings.LastIndex(path, "&opener=") + len ("&opener="):]
		_, capImg := genCaptcha()
		capUri := getCaptchaURI(capImg)
		html := strings.Replace(TOR2WEB_FORM_HTML_LOCAL, "[[cap_uri]]", capUri, -1)
		html = strings.Replace(html, "[[onion_plus_query_type]]", queryType, -1)
		html = strings.Replace(html, "[[opener]]", opener, -1)

		data := path[strings.LastIndex(path, "data=") + len("data="): strings.LastIndex(path, "&captcha=")]
		html = strings.Replace(html, "[[data_value]]", data, -1)
		fmt.Fprint(w, html)
	}
}

/**
 * This function receives POST messages from Tor2web pages
 */
func receivePostMessage(w http.ResponseWriter, req *http.Request, messageType string) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		Println(err.Error())
		return
	}

	defer closeRequestBody(req)

	bodyStr := string(body)

	if !verifyParams(bodyStr, "data=", "&captcha=", "&opener="){
		fmt.Fprint(w, "error")
		return
	}

	data := bodyStr[len("data="): strings.LastIndex(bodyStr, "&captcha=")]
	capStr := bodyStr[strings.LastIndex(bodyStr, "&captcha=") + len("&captcha="): strings.LastIndex(bodyStr, "&opener=")]
	opener := bodyStr[strings.LastIndex(bodyStr, "&opener=") + len ("&opener="):]

	if !verifyCaptcha(capStr){
		queryType := ""
		if messageType == STORE_ENODE{
			queryType = "/queryStoreEnodePost="
		} else if messageType == STORE_GROUP{
			queryType = "/queryStoreGroupPost="
		} else if messageType == REGISTER_VOTER{
			queryType = "/queryRegisterVoterPost="
		} else if messageType == SEND_VOTE{
			queryType = "/querySendVotePost="
		} else if messageType == ENTER_GROUP{
			queryType = "/queryEnterGroupPost="
		} else if messageType == CONFIRMATION{
			queryType = "/queryConfirmPost="
		}

		_, capImg := genCaptcha()
		capUri := getCaptchaURI(capImg)
		html := strings.Replace(TOR2WEB_FORM_HTML, "[[cap_uri]]", capUri, -1)
		html = strings.Replace(html, "[[query_type]]", queryType, -1)
		html = strings.Replace(html, "[[opener]]", opener, -1)
		html = strings.Replace(html, "[[data_value]]", data, -1)
		fmt.Fprint(w, html)
		return
	}

	if verbose {
		if len(data) < 2000{
			Println("Received message has invalid size")
		} else {
			Println(messageType + " message received")
		}
	}

	if messageType == STORE_ENODE{
		storeEnodeMessages = append(storeEnodeMessages, data)
	} else if messageType == STORE_GROUP{
		storeGroupMessages = append(storeGroupMessages, data)
	} else if messageType == REGISTER_VOTER{
		registerVoterMessages = append(registerVoterMessages, data)
	} else if messageType == SEND_VOTE{
		sendVoteMessages = append(sendVoteMessages, data)
	} else if messageType == ENTER_GROUP{
		enterGroupMessages = append(enterGroupMessages, data)
	} else if messageType == CONFIRMATION{
		confirmMessages = append(confirmMessages, data)
	}

	unescaped_opener, _ := url.QueryUnescape(opener)
	done := strings.Replace(DONE_HTML, "[[opener]]", unescaped_opener, -1)
	fmt.Fprint(w, done)
}

/**
 * It informs the .onion address for the voters to send messages
 */
func getHiddenServiceHostname(w http.ResponseWriter) {
	dat, err := ioutil.ReadFile(getHome() + "/.kantcoin/tor/Data/HiddenService/hostname")
	if err == nil {
		fmt.Fprint(w, string(dat))
	} else {
		fmt.Fprint(w, "error")
	}
}

/**
 * Put a message into a slice
 */
func insertStoreEnodeMessage (path string, w http.ResponseWriter) {
	query := path[strings.LastIndex(path, "data=") + len("data="): strings.LastIndex(path, "&captcha=")]
	opener := path[strings.LastIndex(path, "&opener=") + len ("&opener="):]

	storeEnodeMessages = append(storeEnodeMessages, query)
	if verbose{
		Println(STORE_ENODE + " message received")
	}

	unescaped_opener, _ := url.QueryUnescape(opener)
	done := strings.Replace(DONE_HTML, "[[opener]]", unescaped_opener, -1)
	fmt.Fprint(w, done)
}

/**
 *Put a message into a slice
 */
func insertStoreGroupMessage (path string, w http.ResponseWriter) {
	query := path[strings.LastIndex(path, "data=") + len("data="): strings.LastIndex(path, "&captcha=")]
	opener := path[strings.LastIndex(path, "&opener=") + len ("&opener="):]

	storeGroupMessages = append(storeGroupMessages, query)
	if verbose {
		Println(STORE_GROUP + " message received")
	}

	unescaped_opener, _ := url.QueryUnescape(opener)
	done := strings.Replace(DONE_HTML, "[[opener]]", unescaped_opener, -1)
	fmt.Fprint(w, done)
}

/**
 * Put a message into a slice
 */
func insertRegisterVoterMessage (path string, w http.ResponseWriter) {
	query := path[strings.LastIndex(path, "data=") + len("data="): strings.LastIndex(path, "&captcha=")]
	opener := path[strings.LastIndex(path, "&opener=") + len ("&opener="):]

	registerVoterMessages = append(registerVoterMessages, query)
	if verbose {
		Println(REGISTER_VOTER + " message received")
	}

	unescaped_opener, _ := url.QueryUnescape(opener)
	done := strings.Replace(DONE_HTML, "[[opener]]", unescaped_opener, -1)
	fmt.Fprint(w, done)
}

/**
 * Put a message into a slice
 */
func insertSendVoteMessage (path string, w http.ResponseWriter) {
	query := path[strings.LastIndex(path, "data=") + len("data="): strings.LastIndex(path, "&captcha=")]
	opener := path[strings.LastIndex(path, "&opener=") + len ("&opener="):]

	sendVoteMessages = append(sendVoteMessages, query)
	if verbose {
		Println(SEND_VOTE + " message received")
	}

	unescaped_opener, _ := url.QueryUnescape(opener)
	done := strings.Replace(DONE_HTML, "[[opener]]", unescaped_opener, -1)
	fmt.Fprint(w, done)
}

/**
 * Put a message into a slice
 */
func insertEnterGroupMessage (path string, w http.ResponseWriter) {
	query := path[strings.LastIndex(path, "data=") + len("data="): strings.LastIndex(path, "&captcha=")]
	opener := path[strings.LastIndex(path, "&opener=") + len ("&opener="):]

	enterGroupMessages = append(enterGroupMessages, query)
	if verbose {
		Println(ENTER_GROUP + " message received")
	}

	unescaped_opener, _ := url.QueryUnescape(opener)
	done := strings.Replace(DONE_HTML, "[[opener]]", unescaped_opener, -1)
	fmt.Fprint(w, done)
}

/**
 * Put a message into a slice
 */
func insertConfirmMessage (path string, w http.ResponseWriter) {
	query := path[strings.LastIndex(path, "data=") + len("data="): strings.LastIndex(path, "&captcha=")]
	opener := path[strings.LastIndex(path, "&opener=") + len ("&opener="):]

	confirmMessages = append(confirmMessages, query)
	if verbose {
		Println(CONFIRMATION + " message received")
	}

	unescaped_opener, _ := url.QueryUnescape(opener)
	done := strings.Replace(DONE_HTML, "[[opener]]", unescaped_opener, -1)
	fmt.Fprint(w, done)
}

/**
 * It is important to execute this command before the user call "geth init", for the Firewall authorization
 */
func firstTimeGeth(){
	_, err := exec.LookPath(executables[runtime.GOOS]["geth"])

	if err == nil {
		cmd := exec.Command(executables[runtime.GOOS]["geth"],"console") //"--rpc", "--rpcport", "8545", "--rpccorsdomain", "\"*\"",
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, _ := cmd.Output()
		outstr := string(out)
		cmd.Process.Kill()

		if verbose {
			Println("----------------------------------Command response ----------------------------------")
			Println(outstr)
			Println("------------------------------End of command response--------------------------------")
		}

		go func(){
			duration := 25 * time.Second
			time.Sleep(duration)
			cmd.Process.Kill()
		}()
	} else {
		Println("Geth is not installed.")
	}
}

/**
 * Download Geth and then save the file on the due directory
 */
func installGeth(){
	urlStr := ""
	fileStr := ""
	if runtime.GOOS == "windows" {
		//It will change from time to time
		urlStr = url_geth_windows
		fileStr = "temp/geth.zip"
	} else if runtime.GOOS == "linux" {
		urlStr = url_geth_linux
		fileStr = "temp/geth.tar.gz"
	}

	u, _ := url.Parse(urlStr)
	blobURL := azblob.NewBlobURL(*u, azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))

	contentLength := int64(0) // Used for progress reporting to report the total number of bytes being downloaded.

	// NewGetRetryStream creates an intelligent retryable stream around a blob; it returns an io.ReadCloser.
	rs := azblob.NewDownloadStream(context.Background(),
		// We pass more tha "blobUrl.GetBlob" here so we can capture the blob's full
		// content length on the very first internal call to Read.
		func(ctx context.Context, blobRange azblob.BlobRange, ac azblob.BlobAccessConditions, rangeGetContentMD5 bool) (*azblob.GetResponse, error) {
			get, err := blobURL.GetBlob(ctx, blobRange, ac, rangeGetContentMD5)
			if err == nil && contentLength == 0 {
				// If 1st successful Get, record blob's full size for progress reporting
				contentLength = get.ContentLength()
			}
			return get, err
		},
		azblob.DownloadStreamOptions{})

	// NewResponseBodyStream wraps the GetRetryStream with progress reporting; it returns an io.ReadCloser.
	stream := pipeline.NewResponseBodyProgress(rs,
		func(bytesTransferred int64) {
			if verbose {
				fmt.Printf("Downloaded %d of %d bytes.\n", bytesTransferred, contentLength)	
			}
		})
	defer stream.Close() // The client must close the response body when finished with it

	file, err := os.Create(fileStr) // Create the file to hold the downloaded blob contents.
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	written, err := io.Copy(file, stream) // Write to the file by reading from the blob (with intelligent retries).
	if err != nil {
		log.Fatal(err)
	} else {
		_, err := unzip(fileStr, "geth")
		if err == nil {
			Println("╔----------------------------------╗")
			Println("| Done - Hecho - Pronto - Terminé. |")
			Println("╚----------------------------------╝")

			//Running geth for the first time in order that the firewall asks for permission
			if runtime.GOOS == "windows" {
				firstTimeGeth()
			} else if runtime.GOOS == "linux" {
				//Not yet available
			}
		} else {
			Println("Failed to install Geth. Can not unzip the file.")
		}
	}
	_ = written // Avoid compiler's "declared and not used" error

}

/**
 * Writing a new genesis.json file and starting a new geth instance
 */
func setBlockchain(path string, w http.ResponseWriter){
	query := path[len("querySetBlockchain="):]
	query = strings.Replace(query, QUESTION_MARK, "?", -1)
	query = strings.Replace(query, HASHTAG, "#", -1)
	query = strings.Replace(query, DOUBLEQUOTE, "\"", -1)
	query = strings.Replace(query, QUOTE, "'", -1)
	query = strings.Replace(query, BACKSLASH, "\\", -1)

	//Obtaining the chainid and the address to be used in a specific campaign
	parts := strings.Split(query, THE_AND)
	if len(parts) != 5 {
		return
	}
	chainid := parts[0]
	address := parts[1]
	how_many_enodes = parts[2]
	directory := parts[3]
	nonce := parts[4]

	//Creating the directory where the data of this specific blockchain will be placed
	if _, err := os.Stat(getHome() + "/.kantcoin/blockchains/" + directory); os.IsNotExist(err) {
		os.Mkdir(getHome() + "/.kantcoin/blockchains/" + directory,0700)
	}

	//Composing the genesis.json file
	genesis := "{ \"config\": {" +
		"\"chainId\": " + chainid + "," +
		"\"homesteadBlock\": 0," +
		"\"eip155Block\": 0," +
		"\"eip158Block\": 0," +
		"\"byzantiumBlock\": 0" +
		"}," +
		"\"difficulty\": \"20\"," +
		"\"gasLimit\": \"3100000000\"," +
		"\"nonce\": \""+ nonce + "\"," +
		"\"alloc\": {" +
		"\"" + address + "\": { \"balance\": \"10000000000000000000\" }" +
		"}}"


	//Writing genesis string into the file
	data := []byte(genesis)
	err := ioutil.WriteFile(getHome() + "/.kantcoin/blockchains/" + directory + "/genesis.json", data, 0700)
	if err != nil {
		fmt.Fprint(w, "error")
		Println("Error while creating the genesis.json file")
		return
	}

	//Initializing Geth
	initGeth(w, directory)
}

/**
 * Initialize the IPFS service to provide the users' pages
 */
func initIPFS(times int){
	if times == 2 {
		Println("Limit of tries reached")
		return
	}
	_, err := exec.LookPath(executables[runtime.GOOS]["ipfs"])
	existDir := verifyIpfsDir()

	if err == nil && existDir {
		//This command 'daemon' provides user access to IPFS sites
		ipfsCmd = exec.Command(executables[runtime.GOOS]["ipfs"], "daemon")
		ipfsCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		err := ipfsCmd.Start()
		if err == nil {
			Println("IPFS daemon has started")
		} else {
			Println("IPFS error")
		}
	} else if err == nil {
		//First we need to configure the IPFS, and then init it
		cmd := exec.Command(executables[runtime.GOOS]["ipfs"], "init")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, err := cmd.Output()
		outstr := string(out)
		cmd.Process.Kill()

		if verbose{
			Println("----------------------------------Command response ----------------------------------")
			Println(outstr)
			Println("------------------------------End of command response--------------------------------")
		}

		if err == nil {
			Println("IPFS was initialized")
			initIPFS(times + 1)
		} else {
			Println("IPFS was not initialized")
		}

	} else {
		Println("IPFS is not installed.")
	}
}

/**
 * Download IPFS and then save the file on the due directory
 */
func installIPFS(){
	link := ""
	file := ""
	if runtime.GOOS == "windows" {
		//It will change from time to time
		link = url_ipfs_windows
		file = "temp/ipfs.zip"
	} else if runtime.GOOS == "linux" {
		link = url_ipfs_linux
		file = "temp/ipfs.tar.gz"
	}

	out, err := os.Create(file)
	defer out.Close()
	if err == nil {
		resp, err2 := http.Get(link)
		defer closeResponseBody(resp)

		if err2 == nil {
			_, err3 := io.Copy(out, resp.Body)
			if err3 == nil {
				_, err4 := unzip(file, "ipfs")
				if err4 == nil {
					Println("╔----------------------------------╗")
					Println("| Done - Hecho - Pronto - Terminé. |")
					Println("╚----------------------------------╝")
				} else {
					Println("Failed to install IPFS. Can not unzip the file.")
				}
			} else {
				Println("Failed to install IPFS. Can not copy the file.")
			}
		} else {
			Println("Failed to install IPFS. Can not download IPFS.")
		}
	} else {
		Println("Failed to install IPFS. Can not create the file.")
	}
}

/**
 * Initialize the Tor to allow voters to send messages
 */
func initTor(){
	_, err := exec.LookPath(executables[runtime.GOOS]["tor"])

	if err == nil {
		//This command 'daemon' provides user access to IPFS sites
		absolutePath := getHome() + "\\.kantcoin\\tor\\Data\\Tor\\torrc"

		torCmd = exec.Command(executables[runtime.GOOS]["tor"], "-f", absolutePath)
		torCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		err := torCmd.Start()
		if err == nil {
			Println("Tor has started")
		} else {
			Println("Tor error")
		}
	} else {
		Println("Tor is not installed.")
	}
}

/**
 * Download Tor and then save the file on the due directory
 * Also create a new torrc file
 */
func installTor(){
	link := ""
	file := ""
	if runtime.GOOS == "windows" {
		//It will change from time to time
		link = url_tor_windows
		file = "temp/tor.zip"
	} else if runtime.GOOS == "linux" {
		//Not yet available
	}

	out, err := os.Create(file)
	defer out.Close()
	if err == nil {
		resp, err2 := http.Get(link)
		defer closeResponseBody(resp)

		if err2 == nil {
			_, err3 := io.Copy(out, resp.Body)
			if err3 == nil {
				//The executable file should be placed with the main program
				_, err4 := unzip(file, "tor")

				//The Data dir should be placed in the Home dir // to avoid being requested for admin privileges
				os.Mkdir(getHome() + "/.kantcoin/tor", 0700)
				_, err5 := unzip(file, getHome() + "/.kantcoin/tor")
				os.Mkdir (getHome() + "/.kantcoin/tor/Data/HiddenService/", 0700 )

				//Creating the torrc file
				torrcStr := "# Tor plus Kantcoin hidden service \r\n" +
							"DataDirectory " + getHome() + "\\.kantcoin\\tor\\Data\\Tor \r\n" +
							"GeoIPFile " + getHome() + "\\.kantcoin\\tor\\Data\\Tor\\geoip \r\n" +
							"GeoIPv6File " + getHome() + "\\.kantcoin\\tor\\Data\\Tor\\geoip6 \r\n" +
							"HiddenServiceDir " + getHome() + "\\.kantcoin\\tor\\Data\\HiddenService \r\n" +
							"HiddenServicePort 80 127.0.0.1:1988 \r\n"

				f, _ := os.Create(getHome() + "/.kantcoin/tor/Data/Tor/torrc")
				defer f.Close()
				f.WriteString(torrcStr)

				if err4 == nil && err5 == nil{
					Println("╔----------------------------------╗")
					Println("| Done - Hecho - Pronto - Terminé. |")
					Println("╚----------------------------------╝")
				} else {
					Println("Failed to install Tor. Can not unzip the file.")
				}
			} else {
				Println("Failed to install Tor. Can not copy the file.")
			}
		} else {
			Println("Failed to install Tor. Can not download IPFS.")
		}
	} else {
		Println("Failed to install Tor. Can not create the file.")
	}
}

/**
 * These links are often updated, so it is interesting to store them online
 */
func getInstallationLinks(){
	resp, err := http.Get("http://kantcoin.org/links")
	body, _ := ioutil.ReadAll(resp.Body)

	installationLinks := InstallationLinks{}

	defer closeResponseBody(resp)

	if err == nil{
		err2 := json.Unmarshal(body, &installationLinks)
		if err2 == nil{
			url_geth_windows = installationLinks.UGW
			url_geth_linux = installationLinks.UGL
			url_ipfs_windows = installationLinks.UIW
			url_ipfs_linux = installationLinks.UIL
			url_tor_windows = installationLinks.UTW
			geth_folder = installationLinks.GF
			ipfs_folder = installationLinks.IF
		}
	}
}

/**
 * Initialize Geth node (all OSs)
 */
func initGeth(w http.ResponseWriter, directory string){
	_, err := exec.LookPath(executables[runtime.GOOS]["geth"])
	if err == nil {
		//filepath,_ := os.Getwd()
		datadir := getHome() + string(os.PathSeparator) + ".kantcoin" + string(os.PathSeparator) + "blockchains" + string(os.PathSeparator) + directory
		genesis :=  getHome() + string(os.PathSeparator) + ".kantcoin" + string(os.PathSeparator) + "blockchains" + string(os.PathSeparator) + directory + string(os.PathSeparator) +"genesis.json"
		execFile, _ := os.Getwd()
		execFile += string(os.PathSeparator) + executables[runtime.GOOS]["geth_backslash"]

		cmd := exec.Command(execFile , "--datadir", datadir, "init", genesis)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err := cmd.Run()
		cmd.Process.Kill()

		if err == nil {
			Println("Geth has started")
			fmt.Fprint(w, "complete")
		} else {
			Println("Geth has not started")
			fmt.Fprint(w, "error")
		}
	} else {
		Println("Geth is not installed")
		fmt.Fprint(w, "error")
	}
}

/**
 * Verifying if IPFS dir exists
 */
func verifyIpfsDir() bool{
	userprofile := getHome()
	//Check if the .ipfs folder was created in user's PC
	if _, err := os.Stat(userprofile + "/.ipfs/version"); !os.IsNotExist(err) {
		return true
	}
	return false
}

/**
 * Getting the HOME directory
 */
func getHome() string{
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	return usr.HomeDir
}

/**
 * It checks if the informed address belongs to the informed user
 */
func checkUser(path string, w http.ResponseWriter) {
	query := path[len("queryCheckUser="):]
	parts := strings.Split(query, THE_AND)
	if len(parts) != 3 {
		return
	}
	id := parts[0]
	pkey := parts[1]
	provider := parts[2]

	resp, err := http.Get(provider + "/checkUser?id=" + id + "&pkey=" + pkey)
	if err == nil {
		bodyString := ""
		defer closeResponseBody(resp)
		if resp.StatusCode == http.StatusOK {
			bodyBytes,_ := ioutil.ReadAll(resp.Body)
			bodyString = string(bodyBytes)
		}
		fmt.Fprint(w, bodyString)
	} else {
		fmt.Fprint(w, "error")
	}
}

/**
 * Getting user's ip in order to figure out the enode
 * We have to obtain it via third party services
 */
func  whatIsMyIP(w http.ResponseWriter){
	rand.Seed(time.Now().UnixNano())
	urls := []string{
		"https://api.ipify.org/?format=json&callback=",
		"https://jsonip.com/?callback=",
		"https://ipinfo.io/json",
		"https://ipapi.co/json",
	}

	resp, err := http.Get(urls[rand.Intn(len(urls))])
	defer closeResponseBody(resp)

	if err == nil && resp.StatusCode == http.StatusOK{
		bodyBytes,_ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Fprint(w, bodyString)
	} else {
		fmt.Fprint(w, "error")
	}
}

/**
 * It calls geth in order to get the enode
 */
func whatIsMyEnode(path string, w http.ResponseWriter){
	rpcPort := path[len("queryEnode="):]

	if _, err := strconv.Atoi(rpcPort); err != nil {
		fmt.Fprint(w, "error")
		return
	}

	cmd := exec.Command(executables[runtime.GOOS]["geth"], "attach", "http://localhost:" + rpcPort)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

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
	if verbose {
		Println("----------------------------------Command response ----------------------------------")
		Println(outstr)
		Println("------------------------------End of command response--------------------------------")
	}

	if strings.LastIndex(outstr,"enode://") > 0{
		begin := strings.LastIndex(outstr,"enode://")
		end := strings.LastIndex(outstr,":30")
		outstr = outstr[begin: end + 6]
	}

	fmt.Fprint(w, outstr)

	go func(){
		duration := 3 * time.Second
		time.Sleep(duration)
		cmd.Process.Kill()
	}()
}

/**
 * It calls geth in order to get the enode
 */
func addPeer(path string, w http.ResponseWriter){
	parts := strings.Split(path[len("queryAddPeer="):], THE_AND)
	peer := parts[0]
	rpcPort := parts[1]

	if peer == ""{
		fmt.Fprint(w, "error")
		return
	}

	if _, err := strconv.Atoi(rpcPort); err != nil {
		fmt.Fprint(w, "error")
		return
	}

	cmd := exec.Command(executables[runtime.GOOS]["geth"], "attach", "http://localhost:" + rpcPort)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprint(w, "error")
	}
	args := "admin.addPeer('" + peer + "')"

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, args)
	}()

	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if verbose {
		Println("----------------------------------Command response ----------------------------------")
		Println(outstr)
		Println("------------------------------End of command response--------------------------------")
	}

	fmt.Fprint(w, "complete")

	go func(){
		duration := 3 * time.Second
		time.Sleep(duration)
		cmd.Process.Kill()
	}()
}

/**
 * Insert or overwrite a new profile (person or campaign)
 */
func addProfile(path string, w http.ResponseWriter){
	query := path[len("queryAddProfile="):]
	query = strings.Replace(query, QUESTION_MARK, "?", -1)
	query = strings.Replace(query, HASHTAG, "#", -1)
	query = strings.Replace(query, DOUBLEQUOTE, "\"", -1)
	query = strings.Replace(query, QUOTE, "'", -1)
	query = strings.Replace(query, BACKSLASH, "\\", -1)
	parts := strings.Split(query, THE_AND)
	if len(parts) != 3{
		fmt.Fprint(w, "error")
		return
	}

	dir := parts[0]
	content := parts[1]
	filename := parts[2]

	// Do not save kroot profile
	if strings.Index(parts[0],"kroot") >= 0 {
		fmt.Fprint(w, "error")
		return
	}

	//Profiles should be saved in the HOME directory
	os.Mkdir(getHome() + "/.kantcoin/profiles/" + dir, 0700)
	f, err := os.Create(getHome() + "/.kantcoin/profiles/" + dir + "/" + filename)
	if err == nil {
		_, err := f.WriteString(content)
		if err == nil{
			Println("New file: " + dir + "/" + filename)

			//Executing IPFS
			_, err = exec.LookPath(executables[runtime.GOOS]["ipfs"])
			if err == nil {
				newIPFSKey(dir)

				//The error used to verify whether the profile was inserted in IPFS or not.
				var err1 *appError
				address, err1 := newIPFSPage(dir, w, false)
				if err1 == nil {
					err1 = publishIPFS(address, dir)
					if err1 == nil{
						fmt.Fprint(w, "complete")
						//After publishing the file, try to pin it to local storage.
						pinIPFS(address)
					} else {
						fmt.Fprint(w, "error")
					}
				}
			} else {
				fmt.Fprint(w, "error")
				Println("IPFS not installed")
			}
		}
		f.Close()
	} else {
		fmt.Fprint(w, "error")
		Println("File was not created")
	}
}

/**
 * The error used to verify whether the profile was inserted in IPFS or not.
 */
type appError struct {
	Error   error
	Message string
}

/**
 * Before publishing we need to (try to) create a key (all OSs)
 * If a key already exists, IPFS will return a error message saying "refusing to overwrite". It does not affects the ongoing procedures.
 */
func newIPFSKey(dir string){
	execFile, _ := os.Getwd()
	execFile += string(os.PathSeparator) + executables[runtime.GOOS]["ipfs_backslash"]

	cmd := exec.Command(execFile,"key","gen", "--type=rsa", "--size=2048", dir)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
	cmd.Process.Kill()
}

/**
 * Publish some page with some key(dir) (all OSs)
 */
func publishIPFS(address, dir string) *appError{
	execFile, _ := os.Getwd()
	execFile += string(os.PathSeparator) + executables[runtime.GOOS]["ipfs_backslash"]

	cmd := exec.Command(execFile,"name","publish", "--key=" + dir, address)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	outstr := string(out)
	cmd.Process.Kill()

	if verbose {
		Println("----------------------------------Command response ----------------------------------")
		Println(outstr)
		Println("------------------------------End of command response--------------------------------")
	}

	if err == nil {
		Println("Directory " + dir + " was published on IPFS")
	} else {
		Println("Publishing error with the directory " + dir)
		Println(err.Error())
		return &appError{err, "Publishing error with the directory " + dir}
	}
	return nil
}

/**
 * It stores an IPFS object from a given path locally to disk. (all OSs)
 */
func pinIPFS(address string){
	execFile, _ := os.Getwd()
	execFile += string(os.PathSeparator) + executables[runtime.GOOS]["ipfs_backslash"]

	cmd := exec.Command(execFile,"pin","add", address)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	outstr := string(out)
	cmd.Process.Kill()

	if verbose {
		Println("----------------------------------Command response ----------------------------------")
		Println(outstr)
		Println("------------------------------End of command response--------------------------------")
	}

	if err == nil {
		Println(address + "was pinned to local storage")
	} else {
		Println("Pinning error with the address " + address)
		Println(err.Error())
	}
}

/**
 * Adding new IPFS page (all OSs)
 */
func newIPFSPage(dir string, w http.ResponseWriter, show bool) (string, *appError){
	execFile, _ := os.Getwd()
	execFile += string(os.PathSeparator) + executables[runtime.GOOS]["ipfs_backslash"]
	directory := getHome() + string(os.PathSeparator) + ".kantcoin" + string(os.PathSeparator) + "profiles" + string(os.PathSeparator) + dir

	cmd := exec.Command(execFile,"add","-r", directory)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	outstr := string(out)
	cmd.Process.Kill()

	if verbose {
		Println("----------------------------------Command response ----------------------------------")
		Println(outstr)
		Println("------------------------------End of command response--------------------------------")
	}

	//Getting the ipfs address of the directory
	if len(outstr) >=52{
		begin := strings.LastIndex(outstr,"added")
		outstr = outstr[begin + 6: begin + 6 + 46]
	}

	if err == nil {
		Println("IPFS page was inserted")
		if show{
			fmt.Fprint(w, outstr)
		}
	} else {
		Println("IPFS insertion error")
		if show{
			fmt.Fprint(w, "IPFS insertion error")
		}
		return "", &appError{err, "IPFS insertion error"}
	}

	return outstr, nil
}

/**
 * Check wheter the profile exists or not, returning 'true' or 'false'
 */
func profileExists(path string, w http.ResponseWriter){
	userprofile := getHome()
	query := path[len("queryProfileExists="):]
	//Template profile
	if strings.Index(query,"kroot") >= 0{
		fmt.Fprint(w, "true")
		return
	}
	if _, err := os.Stat(userprofile + "/.kantcoin/profiles/" + query); !os.IsNotExist(err) {
		Println("Profile " + query + " opened")
		fmt.Fprint(w, "true")
	} else {
		Println("Profile " + query + " does not exist")
		fmt.Fprint(w, "false")
	}
}

/**
 * Returns the profile html content
 */
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
				return
			}
		} else {
			//A profile that has been already created
			file = userprofile + "/.kantcoin/profiles/" + query
			content, err := ioutil.ReadFile(file)
			if err == nil {
				fmt.Fprint(w, string(content))
				return
			} else {
				//If this profile does not exist, show the initial page
				file = "website/templates/kroot/profile"
				content, err := ioutil.ReadFile(file)
				if err == nil {
					fmt.Fprint(w, string(content))
					return
				}
			}
		}
	} else if strings.Index(query,"data") > 0 {  //...or for user data
		//There is no default user data
		file = userprofile + "/.kantcoin/profiles/" + query
		content, err := ioutil.ReadFile(file)
		if err == nil {
			fmt.Fprint(w, string(content))
			return
		}
	}

	fmt.Fprint(w, "error")
}

/**
 * In order to avoid this error: 'Failed to write genesis block: database already contains an incompatible genesis block'
 */
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

/**
 * This function receives the same params received by the setBlockchain function, then it compares the content of the genesis file and the params
 */
func verifyBlockchain(path string, w http.ResponseWriter) {
	query := path[len("queryVerifyBlockchain="):]
	query = strings.Replace(query, QUESTION_MARK, "?", -1)
	query = strings.Replace(query, HASHTAG, "#", -1)
	query = strings.Replace(query, DOUBLEQUOTE, "\"", -1)
	query = strings.Replace(query, QUOTE, "'", -1)
	query = strings.Replace(query, BACKSLASH, "\\", -1)

	//Obtaining the chainid and the address to be used in a specific campaign
	parts := strings.Split(query, THE_AND)
	if len(parts) != 6 {
		return
	}
	chainid := parts[0]
	address := parts[1]
	//enode := parts[2]
	directory := parts[3]
	nonce := parts[4]
	deleteDirIfDifferent := parts[5]

	//Composing the genesis.json file
	genesis := "{ \"config\": {" +
		"\"chainId\": " + chainid + "," +
		"\"homesteadBlock\": 0," +
		"\"eip155Block\": 0," +
		"\"eip158Block\": 0," +
		"\"byzantiumBlock\": 0" +
		"}," +
		"\"difficulty\": \"20\"," +
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
			if deleteDirIfDifferent == "true" {
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

/**
 * Call: geth --networkid "1151985..." etc
 */
func runBlockchain(path string, w http.ResponseWriter) {
	query := path[len("queryRunBlockchain="):]
	parts := strings.Split(query, THE_AND)
	if len(parts) != 4 {
		return
	}
	id := parts[0]
	address := parts[1]
	dir := parts[2]
	role := parts[3]

	intId,_ := strconv.Atoi(id)
	//If the campaign has not changed, keep the geth process running
	if intId == lastId{
		Println("Old geth process was kept")
		fmt.Fprint(w, "complete")
	} else {
		if gethCmd != nil && gethCmd.Process != nil {
			err := gethCmd.Process.Kill()
			if err == nil{
				Println("Old geth process was killed")
			} else {
				Println("Old geth process was not killed")
			}
		}

		runGeth(w, id, address, dir, role)
		lastId = intId
	}
}

/**
 * It can be used by the creator of the campaign or another node (all OSs)
 */
func runGeth(w http.ResponseWriter, id, address, dir, role string){
	if len(role) != 3{
		return
	}
	execFile, _ := os.Getwd()
	execFile += string(os.PathSeparator) + executables[runtime.GOOS]["geth_backslash"]
	dataDir := getHome() + string(os.PathSeparator) + ".kantcoin" + string(os.PathSeparator) + "blockchains" + string(os.PathSeparator) + dir
	pwdFile := getHome() + string(os.PathSeparator) + ".kantcoin" + string(os.PathSeparator) + "blockchains" + string(os.PathSeparator) + dir + string(os.PathSeparator) + "pwd"
	privkeyFile := getHome() + string(os.PathSeparator) + ".kantcoin" + string(os.PathSeparator) + "blockchains" + string(os.PathSeparator) + dir + string(os.PathSeparator) + "privkey"

	rpcPort := "8" + role
	if role == "001" { //the creator
		gethCmd = exec.Command(execFile,"--datadir", dataDir, "--networkid", id, "--syncmode", "full", "--nodiscover", "--port", "30001", "--maxpeers", how_many_enodes, "--ipcdisable", "--rpc", "--rpcaddr", "localhost", "--rpcport", rpcPort, "--rpcapi", "admin,personal,net,eth,web3", "--rpccorsdomain", "http://localhost:1985,http://127.0.0.1:1985", "--mine", "--minerthreads=1", "--etherbase", address, "--unlock", address, "--password", pwdFile, "console")
	} else if strings.Index(role, "0") == 0 || strings.Index(role, "1") == 0 || strings.Index(role, "2") == 0 { //group chairpersons
		gethCmd = exec.Command(execFile,"--datadir", dataDir, "--networkid", id, "--syncmode", "full", "--nodiscover", "--port", "30" + role, "--maxpeers", how_many_enodes, "--ipcdisable", "--rpc", "--rpcaddr", "localhost", "--rpcport", rpcPort, "--rpcapi", "admin,personal,net,eth,web3", "--rpccorsdomain", "http://localhost:1985,http://127.0.0.1:1985", "--unlock", address, "--password", pwdFile, "console")
	} else if strings.Index(role, "3") == 0 || strings.Index(role, "4") == 0 || strings.Index(role, "5") == 0 ||
			strings.Index(role, "6") == 0 || strings.Index(role, "7") == 0 || strings.Index(role, "8") == 0 || strings.Index(role, "9") == 0{ //observers
		gethCmd = exec.Command(execFile,"--datadir", dataDir, "--networkid", id, "--syncmode", "full", "--nodiscover", "--port", "30" + role, "--maxpeers", how_many_enodes, "--ipcdisable", "--rpc", "--rpcaddr", "localhost", "--rpcport", rpcPort, "--rpcapi", "admin,personal,net,eth,web3", "--rpccorsdomain", "http://localhost:1985,http://127.0.0.1:1985", "--unlock", address, "--password", pwdFile, "console")
	}

	gethCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stdinGeth, _ = gethCmd.StdinPipe()
	err := gethCmd.Start()

	if err == nil {
		Println("Geth running")
		fmt.Fprint(w, "complete")
	} else {
		Println("Geth not running")
		fmt.Fprint(w, "error")
	}

	//Removing the privatekey and password files
	go func(){
		duration := 60 * time.Second
		time.Sleep(duration)

		os.Remove(privkeyFile)
		os.Remove(pwdFile)
	}()
}

/**
 * Verify if IPFS has already been initialized
 */
func isIPFSRunning() bool{
	resp, err := http.Get("http://localhost:8080/ipfs/QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG/readme")
	if err == nil && resp.StatusCode == http.StatusOK{
		Println("IPFS has already been initialized")
		return true
	}
	return false
}

/**
 * Before calling runBlockchainLin/Win, create the password file to unlock the main account
 */
func createPwdFile(path string, w http.ResponseWriter) {
	query := path[len("queryCreatePwdFile="):]
	parts := strings.Split(query, THE_AND)
	if len(parts) != 2 {
		return
	}
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

/**
 * Create new file with a private key (and another with a password) to be imported with the command: geth account import
 */
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

/**
 * Insert a new account with the command geth account import. In order to do that, create new privatekey and password file.
 */
func insertAccountIntoBlockchain(path string, w http.ResponseWriter) {
	query := path[len("queryInsertAccountIntoBlockchain="):]
	parts := strings.Split(query, THE_AND)
	if len(parts) != 3 {
		return
	}
	dir := parts[0]
	privkey := parts[1]
	password := parts[2]

	if createPrivateKeyFile(dir, privkey, password){
		newAccount(w, dir)
	}
}

/**
 * Call the command: geth account import (all OSs)
 */
func newAccount(w http.ResponseWriter, dir string) {
	_, err := exec.LookPath(executables[runtime.GOOS]["geth"])
	if err == nil {
		execFile, _ := os.Getwd()
		execFile += string(os.PathSeparator) + executables[runtime.GOOS]["geth_backslash"]
		dataDir := getHome() + string(os.PathSeparator) + ".kantcoin" + string(os.PathSeparator) + "blockchains" + string(os.PathSeparator) + dir
		pwdFile := getHome() + string(os.PathSeparator) + ".kantcoin" + string(os.PathSeparator) + "blockchains" + string(os.PathSeparator) + dir + string(os.PathSeparator) + "pwd"
		privkeyFile := getHome() + string(os.PathSeparator) + ".kantcoin" + string(os.PathSeparator) + "blockchains" + string(os.PathSeparator) + dir + string(os.PathSeparator) + "privkey"

		cmd := exec.Command(execFile, "--datadir", dataDir, "account", "import", privkeyFile, "--password", pwdFile)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, err := cmd.Output()
		outstr := string(out)
		cmd.Process.Kill()

		if verbose {
			Println("----------------------------------Command response ----------------------------------")
			Println(outstr)
			Println("------------------------------End of command response--------------------------------")
		}

		if err == nil {
			Println("Account inserted")
			fmt.Fprint(w, "complete")
		} else {
			Println("Account not inserted")
			fmt.Fprint(w, "error")
		}
	} else {
		Println("Geth is not installed")
		fmt.Fprint(w, "error")
	}
}

/**
 * Creating a key with the name received
 */
func addIPNS(path string, w http.ResponseWriter){
	id := path[len("queryAddIPNSKey="):]
	_, err := exec.LookPath(executables[runtime.GOOS]["ipfs"])
	if err == nil {
		newIPFSKey(id)
		fmt.Fprint(w, "complete")
	} else {
		fmt.Fprint(w, "error")
	}
}

/**
 * Call the "ipfs key list -l" to obtain the ipns address
 */
func getIPNS(path string, w http.ResponseWriter){
	id := path[len("queryGetIPNS="):]

	if currentIpns.id == id{
		fmt.Fprint(w, currentIpns.ipns)
		return
	}
	_, err := exec.LookPath(executables[runtime.GOOS]["ipfs"])
	if err == nil {
		ipns := ""
		execFile, _ := os.Getwd()
		execFile += string(os.PathSeparator) + executables[runtime.GOOS]["ipfs_backslash"]

		cmd := exec.Command(execFile,"key", "list", "-l")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, err := cmd.Output()
		outStr := string(out)
		cmd.Process.Kill()

		if verbose {
			Println("----------------------------------Command response ----------------------------------")
			Println(outStr)
			Println("------------------------------End of command response--------------------------------")
		}

		if err == nil{
			lines := stringToLines(outStr)
			ipns = findIPNS(lines, id)

			currentIpns.id = id
			currentIpns.ipns = ipns
			if w != nil{
				fmt.Fprint(w, ipns)
			}
		} else {
			if w != nil{
				fmt.Fprint(w, "error")
			}
		}
	} else {
		fmt.Fprint(w, "error")
	}
}

/**
 * Check if the .kantcoin directory was created. If it was not created, then create it
 */
func initDirectory(){
	userprofile := getHome()
	if _, err := os.Stat(userprofile + "/.kantcoin"); os.IsNotExist(err) {
		os.Mkdir(userprofile + "/.kantcoin",0700)
	}
	if _, err := os.Stat(userprofile + "/.kantcoin/profiles"); os.IsNotExist(err) {
		os.Mkdir(userprofile + "/.kantcoin/profiles",0700)
	}
	if _, err := os.Stat(userprofile + "/.kantcoin/blockchains"); os.IsNotExist(err) {
		os.Mkdir(userprofile + "/.kantcoin/blockchains",0700)
	}
}

/**
 * Load the executables map
 */
func loadExecutables(){
	add(executables, "windows","ipfs", "ipfs/" +ipfs_folder+ "/ipfs.exe")
	add(executables, "windows","ipfs_backslash", "\\ipfs\\" +ipfs_folder+ "\\ipfs.exe")
	add(executables, "windows","geth", "geth/" +geth_folder+ "/geth.exe")
	add(executables, "windows","geth_backslash", "\\geth\\" +geth_folder+ "\\geth.exe")
	add(executables, "windows","tor", "tor/Tor/tor.exe")
	add(executables, "linux","ipfs", "ipfs/" +ipfs_folder+ "/ipfs")
	add(executables, "linux","ipfs_backslash", "ipfs/" +ipfs_folder+ "/ipfs")
	add(executables, "linux","geth", "geth/" +geth_folder+ "/geth")
	add(executables, "linux","geth_backslash", "geth/" +geth_folder+ "/geth")
}

/**
 * This function builds the executables map
 */
func add(m map[string]map[string]string, os, cmd, value string) {
	mm, ok := m[os]
	if !ok {
		mm = make(map[string]string)
		m[os] = mm
	}
	mm[cmd] = value
}

/**
 * Trasforming a string in an slice
 */
func stringToLines(s string) []string {
	var lines []string

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

/**
 * Look each line to find the ipns of a certain campaign
 */
func findIPNS(lines []string, id string) string{
	for i:=0; i < len(lines); i = i + 1{
		if strings.Index(lines[i], id) > 0{
			return lines[i][:46]
		}
	}
	return "error"
}

/**
 * Unzip will un-compress a zip archive,
 * moving all files and folders to an output directory
 */
func unzip(src, dest string) ([]string, error) {
	var fileNames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return fileNames, err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return fileNames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fPath := filepath.Join(dest, f.Name)
		fileNames = append(fileNames, fPath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fPath, os.ModePerm)

		} else {
			// Make File
			var fdir string
			if lastIndex := strings.LastIndex(fPath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fPath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return fileNames, err
			}
			f, err := os.OpenFile(
				fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fileNames, err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return fileNames, err
			}

		}
	}
	return fileNames, nil
}

/**
 * It allows to dispatch HTTP requests via Tor in Go
 */
func initTorClient(){
	// Create a transport that uses Tor Browser's SocksPort.  If
	// talking to a system tor, this may be an AF_UNIX socket, or
	// 127.0.0.1:9050 instead.
	tbProxyURL, err := url.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		Println("Failed to parse proxy URL")
	}

	// Get a proxy Dialer that will create the connection on our
	// behalf via the SOCKS5 proxy.  Specify the authentication
	// and re-create the dialer/transport/client if tor's
	// IsolateSOCKSAuth is needed.
	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		Println("Failed to obtain proxy dialer")
	}

	// Make a http.Transport that uses the proxy dialer, and a
	// http.Client that uses the transport.
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	torClient = &http.Client{Transport: tbTransport}
}

/**
 * Sending a request using the Tor client
 */
func sendTorRequest(path string, w http.ResponseWriter){
	query := path[len("querySendTorRequest="):]

	resp, err := torClient.Get("http://" + query)
	if err != nil {
		fmt.Fprint(w, "error")
		return
	}
	defer closeResponseBody(resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(w,"error")
		return
	}
	fmt.Fprint(w, string(body))
}

/**
 * It prints the voter's group index on the page
 */
func myGroupIndex(path string, w http.ResponseWriter) {
	voterAddr := path[len("queryMyGroupIndex="):]

	if len(voterAddr) == 0 {
		fmt.Fprint(w, "[index_error]")
		return
	}

	if group, err := votersGroupsDatabase.Read(voterAddr); err == nil {
		fmt.Fprint(w, string(group))
	} else {
		fmt.Fprint(w, "-1")
	}
}

/**
 * It prints the voter's group index on the page
 */
func getVoteMessage(path string, w http.ResponseWriter) {
	if !verifyParams(path, "voter=", "&group="){
		fmt.Fprint(w, "error")
		return
	}

	voterAddr := path[len("queryGetVoteMessage=voter="): strings.LastIndex(path,"&group=")]
	groupStr := path[strings.LastIndex(path,"&group=") + len("&group="):]
	group, err := strconv.Atoi(groupStr)

	if len(groupsMessages) <= group{
		fmt.Fprint(w, "error")
		return
	}

	if vg, err := votersGroupsDatabase.Read(voterAddr); err != nil || string(vg) != groupStr{
		fmt.Fprint(w, "error")
		return
	}

	if err == nil{
		if secret, err := secretsDatabase.Read(voterAddr); err == nil {
			encryptedMessages := AESEncrypt(string(secret), groupsMessages[group])

			fmt.Fprint(w, encryptedMessages)
		} else {
			fmt.Fprint(w, "error")
		}
	}
}

/**
 * It stores the voter in the voters map with an empty (-2) group
 */
func storeVoter(path string, w http.ResponseWriter){
	voterAddr := path[len("queryStoreVoter="):]
	votersGroupsDatabase.Write(voterAddr, []byte("-2"))
	fmt.Fprint(w, "complete")
}

/**
 * It stores the information of one group into the groupsInfo slice
 */
func storeGroupInfo(path string, w http.ResponseWriter) {
	index := path[len("queryStoreGroupInfo=index="): strings.LastIndex(path,"&voters=")]
	votersStr := path[strings.LastIndex(path, "&voters=") + len("&voters="): strings.LastIndex(path,"&info=")]
	info := path[strings.LastIndex(path, "&info=") + len("&info="):]

	voters := strings.Split(votersStr, ",")
	for _, voter := range voters {
		//Voters can change their group (in case of fraud)
		if voter != "0x0000000000000000000000000000000000000000"{
			if group, err := votersGroupsDatabase.Read(voter); err != nil || string(group) != index {
				votersGroupsDatabase.Write(voter, []byte(index))
			}
		}
	}

	if indexNumber, _ := strconv.Atoi(index); indexNumber >= len(groupsInfo){
		groupsInfo = append(groupsInfo, info)
	} else {
		groupsInfo[indexNumber] = info
	}
	fmt.Fprint(w, "complete")
}

/**
 * It returns the information of certain group
 */
func getGroupInfo(path string, w http.ResponseWriter) {
	if !verifyParams(path, "index="){
		fmt.Fprint(w, "[group_info_error]")
		return
	}

	if len(path) > len("queryGetGroupInfo=index=") && strings.LastIndex(path, "queryGetGroupInfo=index=") == 0{
		index, err := strconv.Atoi(path[len("queryGetGroupInfo=index="):])

		if err != nil || index >= len(groupsInfo){
			fmt.Fprint(w, "[group_info_error]")
		} else if index >= 0{
			fmt.Fprint(w, groupsInfo[index])
		} else {
			fmt.Fprint(w, "[group_info_error]")
		}
	} else {
		fmt.Fprint(w, "[group_info_error]")
	}
}

/**
 * It stores the group's vote message
 */
func storeGroupMessage(path string, w http.ResponseWriter) {
	index, _ := strconv.Atoi(path[len("queryStoreGroupMessage=group="): strings.LastIndex(path,"&message=")])
	message := path[strings.LastIndex(path, "&message=") + len("&message="):]

	if index >= len(groupsMessages){
		groupsMessages = append(groupsMessages, message)
	} else if index >= 0{
		groupsMessages[index] = message
	}
	fmt.Fprint(w, "complete")
}

/**
 * It caches the information of the Campaign
 */
func storeCampaignInfo(path string, w http.ResponseWriter) {
	info := path[len("queryStoreCampaignInfo="):]
	campaignInfo = info
	fmt.Fprint(w, "complete")
}

/**
 * It caches the Campaign's IPFS data
 */
func storeCampaignIPFSInfo(path string, w http.ResponseWriter) {
	info := path[len("queryStoreCampaignIPFSInfo="):]
	campaignIPFSInfo = info
	fmt.Fprint(w, "complete")
}

/**
 * This secret is used to send the group vote message to voters
 */
func storeVoterSecret(path string, w http.ResponseWriter) {
	voter := path[len("queryStoreVoterSecret=voter="): strings.LastIndex(path, "&secret=")]
	secret := path[strings.LastIndex(path, "&secret=") + len("&secret="):]

	//votersSecretsMap[voter] = secret
	secretsDatabase.Write(voter, []byte(secret))
	fmt.Fprint(w, "complete")
}

/**
 * It returns the information of the Campaign
 */
func getCampaignInfo(w http.ResponseWriter) {
	if campaignInfo == ""{
		fmt.Fprint(w, "[campaign_info_error]")
	} else {
		fmt.Fprint(w, campaignInfo)
	}
}

/**
 * It returns the Campaign's IPFS data
 */
func getCampaignIPFSInfo(w http.ResponseWriter) {
	if campaignIPFSInfo == ""{
		fmt.Fprint(w, "[campaign_ipfs_info_error]")
	} else {
		fmt.Fprint(w, campaignIPFSInfo)
	}
}

/**
 * It closes the response body if it is not null (to avoid runtime errors)
 */
func closeResponseBody(resp *http.Response){
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

/**
 * It closes the request body if it is not null (to avoid runtime errors)
 */
func closeRequestBody(req *http.Request){
	if req != nil && req.Body != nil {
		req.Body.Close()
	}
}

/**
 * Terminating the program
 */
func exit(){
	stopRecevingRequests = true

	if ipfsCmd != nil && ipfsCmd.Process != nil {
		cmd := exec.Command(executables[runtime.GOOS]["ipfs"], "shutdown")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Run()
		if ipfsCmd != nil && ipfsCmd.Process != nil {
			ipfsCmd.Process.Kill()
		}
		Println("IPFS killed")
	}

	if torCmd != nil && torCmd.Process != nil {
		torCmd.Process.Kill()
		Println("Tor killed")
	}

	if gethCmd != nil && gethCmd.Process != nil {
		args := "exit"
		io.WriteString(stdinGeth, args)
		gethCmd.CombinedOutput()

		go func(){
			stdinGeth.Close()
			duration := 2 * time.Second
			time.Sleep(duration)
			gethCmd.Process.Kill()

			//The last statements
			Println("Geth killed")
			Println("Ciao")
			ioutil.WriteFile(getHome() + "/.kantcoin/logs.txt", []byte(logs), 0700)
		}()
	}

	theCaptcha = nil
	aelectron.Close()

	duration := 3 * time.Second
	time.Sleep(duration)
	os.Exit(0)
}

/**
 * Print logs in a new line
 */
func Println(str string){
	logs += time.Now().Format("2006/01/02 15:04:05") + " " + str + "\r\n"
	if installerLogs && installerTE != nil{
		installerTE.SetText(logs)
	}
}

/**
 * Functions used to encrypt the vote message
 */
func deriveKey(passphrase string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		// http://www.ietf.org/rfc/rfc2898.txt
		// Salt.
		rand.Read(salt)
	}
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}

func AESEncrypt(passphrase, plaintext string) string {
	key, salt := deriveKey(passphrase, nil)
	iv := make([]byte, 12)
	// http://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38d.pdf
	// Section 8.2
	rand.Read(iv)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return hex.EncodeToString(salt) + "-" + hex.EncodeToString(iv) + "-" + hex.EncodeToString(data)
}

/**
 * It verifies if all these parameters are present in the path string, and in the given order
 */
func verifyParams(path string, params... string) bool{
	for i, param := range params {
		if strings.LastIndex(path, param) == -1{
			return false
		}
		if i > 0 && strings.LastIndex(path, params[i - 1]) > strings.LastIndex(path, param){
			return false
		}
	}
	return true
}