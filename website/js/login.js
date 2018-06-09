/**
 * Kantcoin Project
 * https://kantcoin.org
 * This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
 * If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

/********************  jQuery configurations ********************/

$(document).ready(function(){
    $('.aniview').AniView();
    $('.carousel.carousel-slider').carousel({
        fullWidth: true,
        indicators: true
    });
    $('.modal').modal({dismissible:false});
    $('.parallax').parallax();
})

/******************** Constants ********************/

//Minimum password length
const PASSWORD_LENGTH = 8

/******************** Global variables ********************/

//Words and name obtained from the file
let file_data = {kantcoin_org_words : "", kantcoin_org_data : ""}
//This provider allows us to register vaults and check users
let login_provider = "https://login.kantcoin.com"
//Wallet with keys the size required by geth
let wallet = {}

/******************** Event listeners ********************/

/**
 *  Setting some variables
 */
window.addEventListener("load", function(){
    let locale = "en"
    if (navigator.language){
        locale = navigator.language.substring(0, 2).toLowerCase()
    }

    if(locale == 'pt'){
        klang = klang.portuguese
    } else if(locale == 'fr'){
        klang = klang.french
    } else if(locale == 'es'){
        klang = klang.spanish
    } else{
        klang = klang.english
    }

    login_about.innerHTML = klang.login_about
    login_download.innerHTML = klang.login_download
    login_title.innerHTML = klang.login_title
    login_message.innerHTML = klang.login_message
    login_create_button.innerHTML = klang.login_create_button
    login_open_button.innerHTML = klang.login_open_button
    modal1_title.innerHTML = klang.modal1_title
    modal1_file.innerHTML = klang.modal1_file
    modal1_password.innerHTML = klang.password
    modal1_user.innerHTML = klang.modal1_user
    modal1_open.text = klang.modal1_open
    modal2_title.innerHTML = klang.modal2_title
    modal2_file_name.innerHTML = klang.modal2_file_name
    modal2_save.text = klang.modal2_save
    modal4_title.innerHTML = klang.modal4_title
    modal4_user.innerHTML = klang.login
    modal4_password1.innerHTML = klang.modal4_password1
    modal4_password2.innerHTML = klang.modal4_password2
    modal4_confirm.text = klang.modal4_confirm
    modal5_title.innerHTML = klang.modal5_title
    modal5_login_provider.innerHTML = klang.modal5_login_provider
    modal5_send.text = klang.modal5_send
    privacy_link.innerHTML = klang.privacy
    motto_label.innerHTML = klang.motto
    resources_label.innerHTML = klang.resources
    contact_link.innerHTML = klang.contact
    conduct_link.innerHTML = klang.conduct
    license_link.innerHTML = klang.license
    for_whom_label.innerHTML = klang.for_whom_label
    polling_org_label.innerHTML = klang.polling_org_label
    unions_label.innerHTML = klang.unions_label
    parties_label.innerHTML = klang.parties_label
    universities_label.innerHTML = klang.universities_label
    schools_label.innerHTML = klang.schools_label
    municipalities_label.innerHTML = klang.municipalities_label
    many_more_label.innerHTML = klang.many_more_label
    easy_voters_label.innerHTML = klang.easy_voters_label
    easy_voters_text.innerHTML = klang.easy_voters_text
    promising_technologies_label.innerHTML = klang.promising_technologies_label
    promising_technologies_text.innerHTML = klang.promising_technologies_text
    simple_campaigns_label.innerHTML = klang.simple_campaigns_label
    acknowledgment_link.innerHTML = klang.acknowledgment

    //Setting the data-tooltip of the password input
    let aElementP1 = $('#password_div1')
    aElementP1.attr('data-tooltip', klang.min_chars)
    aElementP1.tooltip()

    //Setting the data-tooltip of the user input
    let aElementP2 = $('#user_div1')
    aElementP2.attr('data-tooltip', klang.hyphen)
    aElementP2.tooltip()

    if (detectIE()){
        buttons_div.innerHTML = browser_not_supported_div.innerHTML.replace("<!--[CDATA[","").replace("-->","")
        browser_not_supported_label.innerHTML = klang.browser_not_supported_label
    }
})

modal1_open.addEventListener("click", openVault)
modal4_confirm.addEventListener("click", createVault)
modal2_save.addEventListener("click", saveVault)
modal5_send.addEventListener("click", sendVault)

/**
 * It displays a JQCloud with all the open source projects used in this project
 */
acknowledgment_link.addEventListener("click", showThanks)

/**
 * Submiting modal forms in case of enter pressed
 */
kantcoin_org_password1.addEventListener("keypress", function (event) {
    let keyCode = event.keyCode
    if(keyCode == 13){
        modal1_open.click()
    }
})

file_name1.addEventListener("keypress", function (event) {
    let keyCode = event.keyCode
    if(keyCode == 13){
        modal2_save.click()
    }
})

kantcoin_org_password4.addEventListener("keypress", function (event) {
    let keyCode = event.keyCode
    if(keyCode == 13){
        modal4_confirm.click()
    }
})

/**
 * It reads the vault file (which contains the mnemonics, ekhash and login provider)
 */
words_file_button.addEventListener("change", function (event){
    let input = event.target
    let reader = new FileReader()

    //Cleaning these variables
    file_data.kantcoin_org_data = ""
    file_data.kantcoin_org_words = ""

    reader.addEventListener("load", function(){
        let text = reader.result
        let split = text.split("\r\n")

        //There are only two lines in this file
        if (split.length == 2){
            file_data.kantcoin_org_words = split[0]
            file_data.kantcoin_org_data = split[1]
        }
    })
    reader.readAsText(input.files[0])
})

/******************** Functions ********************/

/**
 * It opens the file that constais the user's mnemonics and name
 * Or the user may insert this data manually
 */
function openVault(){
    //It is necessary to generate the keys
    let password = kantcoin_org_password1.value
    //For the ekhash
    let user = kantcoin_org_user1.value

    //Checking if the user login, password and the words were given
    if (!file_data.kantcoin_org_words || !password || !user){
        Materialize.toast(klang.no_vault_opened, 3000, 'rounded')
        return
    }

    login_message.innerHTML = preloader_div.innerHTML.replace("<!--[CDATA[","").replace("-->","")

    //Cleaning the localStorage and the sessionStorage
    localStorage.setItem("kantcoin_org_words","")
    sessionStorage.setItem("kantcoin_org_key","")
    sessionStorage.setItem("kantcoin_org_aux_pubkey","")
    localStorage.setItem("kantcoin_org_login_provider","")
    localStorage.setItem("kantcoin_org_ekhash","")
    localStorage.setItem("kantcoin_tor_privkey","")
    localStorage.setItem("kantcoin_tor_pubkey","")

    //Setting the user in the sessionstorage
    sessionStorage.setItem("kantcoin_org_user",user)

    let words = file_data.kantcoin_org_words
    file_data.kantcoin_org_words = ""

    let m = new Mnemonic("english")

    let is_valid = m.check(words)

    //Checking if the mnemonics were valid
    if (is_valid){
        localStorage.setItem("kantcoin_org_words", words)
    } else {
        Materialize.toast(klang.browser_words_conflict, 4000, 'rounded')
        login_message.innerHTML = klang.login_message
        return
    }

    //Closing the modal
    $('#modal1').modal('close')

    let data = file_data.kantcoin_org_data
    //file_data.kantcoin_org_data = ""

    //Decrypting data
    decrypt(words, data).then(function (dec) {
        //Filling the localStorage
        loadFileData(dec)
        //Generating the private and public keys and putting them in the sessionStorage
        generateKey(words, password, user).then(function(){
            //Checking if the user and the password are correct (probably)
            if(!checkuserNPassword(user, sessionStorage.getItem("kantcoin_org_key"), localStorage.getItem("kantcoin_org_ekhash"))){
                Materialize.toast(klang.wrong_user_or_password, 3500, 'rounded')

                //Cleaning the fields and sessionStorage
                sessionStorage.setItem("kantcoin_org_key","")
                sessionStorage.setItem("kantcoin_org_user","")
                sessionStorage.setItem("kantcoin_org_aux_pubkey","")
                sessionStorage.setItem("kantcoin_org_wallet", "")

                login_message.innerHTML = klang.login_message
                return
            }

            window.location.href = (probablyTor() ? "http://" : "https://") + "kantcoin.org/home"
        })
    })
}

/**
 * It uses hashcodes to check if the password and user login are correct
 * @param {string} user
 * @param {string} key
 * @param {string} ekhash
 * @returns {boolean}
 */
function checkuserNPassword(user, key, ekhash){
    let hash = "" + hashCode(key + user)
    hash = hash.substr(ekhash.length > 4 ? 4 : 0)
    if (hash == ekhash){
        return true
    }
    return false
}

/**
 * It displays a JQCloud with all the open source projects used in this project
 */
function showThanks(){
    let words = [
        {text: "Kantcoin", weight: 8, link: "https://github.com/kantcoin"},
        {text: "Ethereum", weight: 6, link: "https://www.ethereum.org"},
        {text: "Tor", weight: 6, link: "https://www.torproject.org/projects/torbrowser.html"},
        {text: "Tor2web", weight: 2, link: "https://www.torproject.org/projects/torbrowser.html"},
        {text: "IPFS", weight: 6, link: "https://ipfs.io"},
        {text: "NSIS", weight: 4, link: "http://nsis.sourceforge.net/Main_Page"},
        {text: "URS", weight: 6, link: "https://github.com/monero-project/urs"},
        {text: "Ethereumjs", weight: 2, link: "https://github.com/ethereumjs/ethereumjs-tx"},
        {text: "Open-golang", weight: 2, link: "https://github.com/skratchdot/open-golang/"},
        {text: "Materialize", weight: 4, link: "http://materializecss.com"},
        {text: "Walk", weight: 2, link: "https://github.com/lxn/walk"},
        {text: "Diskv", weight: 2, link: "https://github.com/peterbourgon/diskv"},
        {text: "Electron", weight: 2, link: "https://github.com/electron/electron"},
        {text: "Aniview", weight: 2, link: "https://github.com/jjcosgrove/jquery-aniview"},
        {text: "Animatecss", weight: 2, link: "https://github.com/daneden/animate.css/"},
        {text: "Golang", weight: 4, link: "https://golang.org"},
        {text: "MaterialNote", weight: 4, link: "https://github.com/Cerealkillerway/materialNote"},
        {text: "NTRU", weight: 2, link: "https://github.com/NTRUOpenSourceProject/ntru-crypto"},
        {text: "Web3.js", weight: 4, link: "https://github.com/ethereum/web3.js/"},
        {text: "Sheet JS", weight: 4, link: "https://github.com/SheetJS/js-xlsx"},
        {text: "ntru.js", weight: 4, link: "https://github.com/cyph/ntru.js"},
        {text: "JQCloud", weight: 2, link: "https://github.com/lucaong/jQCloud"},
        {text: "Account Kit", weight: 6, link: "https://developers.facebook.com/docs/accountkit"},
        {text: "Bip39", weight: 2, link: "https://github.com/bitcoinjs/bip39"},
        {text: "ethers.js", weight: 2, link: "https://github.com/ethers-io/ethers.js/"},
        {text: "BitcoinJS", weight: 2, link: "https://github.com/bitcoinjs/bitcoinjs-lib"},
        {text: "JQuery", weight: 2, link: "https://jquery.com"},
        {text: "cryptocoinjs", weight: 2, link: "https://github.com/cryptocoinjs"},
        {text: "Astilectron", weight: 4, link: "https://github.com/asticode/go-astilectron"}
    ]

    acknowledgment_div.innerHTML = acknowledgment_html.innerHTML.replace("<!--[CDATA[","").replace("-->","").replace("[[thanks]]", klang.thanks_community)

    sleep(600).then(function () {
        $('#word_cloud').jQCloud(words)
    })
}

/**
 * This file has only the mnemonics, the login provider and a hash to verify if the user name matches with the password
 * @param {string} data_text
 */
function loadFileData(data_text){
    let jsonObj = JSON.parse(data_text)
    localStorage.setItem("kantcoin_org_login_provider",jsonObj.login_provider)
    localStorage.setItem("kantcoin_org_ekhash",jsonObj.ekhash)
    localStorage.setItem("kantcoin_org_secrets_base", jsonObj.secrets_base)
}

/**
 * It generates the main key and the aux public key
 * @param {string} words
 * @param {string} password
 * @param {string} user
 * @returns {Object}
 */
function generateKey(words, password, user){
    //Using bitcoinjsb to generate a privatekey from the mnemonics and the password
    //This key is used to sign and verify ring signatures
    let keys = bitcoinjsb.keypairsFromMnemonic(words, password, 3)
    let privkey = bs58.decode(keys[0].keyPair.toWIF()).toString("hex")

    let aux_signingkey = new ethers.SigningKey("0x" + privkey)
    let pubkey = aux_signingkey.publicKey.substring(2)
    let address = aux_signingkey.address.substring(2)

    let keyjson = "{\"address\":\"" + address + "\",\"privkey\":\"" + privkey + "\",\"pubkey\":\"" + pubkey + "\"}"

    //Generating the aux_pkey
    let privkey2 = bs58.decode(keys[1].keyPair.toWIF()).toString("hex")
    let aux_signingkey2 = new ethers.SigningKey("0x" + privkey2)
    //This public key is used to generate the directory
    let aux_pkey = aux_signingkey2.publicKey.substring(2)

    //The key is stored in window.sessionStorage for user security (less time exposed)
    sessionStorage.setItem("kantcoin_org_key",keyjson)
    sessionStorage.setItem("kantcoin_org_aux_pubkey",aux_pkey)

    //Generating new wallet from the mnemonics, user, and password
    let promise = ethers.Wallet.fromBrainWallet(words, user + password).then(function(the_wallet) {
        wallet = the_wallet
        sessionStorage.setItem("kantcoin_org_wallet", JSON.stringify(wallet))
    })

    return promise
}

/**
 * It creates a new vault, but doesn't save it into a file
 */
function createVault() {
    let password3 = kantcoin_org_password3.value
    let password4 = kantcoin_org_password4.value
    let user = kantcoin_org_user3.value

    //Checking if empty
    if (!password3 || !password4 || !user){
        Materialize.toast(klang.empty_fields, 2000, 'rounded')
        return
    }

    //Checking if passwords match
    if (password3 != password4){
        Materialize.toast(klang.different_passwords, 2000, 'rounded')
        return
    }

    //Verifying the password length
    if (password3.length < PASSWORD_LENGTH){
        Materialize.toast(klang.password_too_small, 2000, 'rounded')
        return
    }

    login_message.innerHTML = preloader_div.innerHTML.replace("<!--[CDATA[","").replace("-->","")

    let m = new Mnemonic("english")

    // Generating new mnemonics
    let words = m.generate(128)

    //Using bitcoinjsb to generate a privatekey from the mnemonics and the password
    //This key is used to sign and verify ring signatures
    let keys = bitcoinjsb.keypairsFromMnemonic(words, password3, 3)
    let privkey = bs58.decode(keys[0].keyPair.toWIF()).toString("hex")

    let aux_signingkey = new ethers.SigningKey("0x" + privkey)
    let pubkey = aux_signingkey.publicKey.substring(2)
    let address = aux_signingkey.address.substring(2)

    //This will be used to generate all voter votes
    let keyjson = "{\"address\":\"" + address + "\",\"privkey\":\"" + privkey + "\",\"pubkey\":\"" + pubkey + "\"}"

    //Generating the aux_pkey
    let privkey2 = bs58.decode(keys[1].keyPair.toWIF()).toString("hex")
    let aux_signingkey2 = new ethers.SigningKey("0x" + privkey2)
    //This public key is used to generate the directory
    let aux_pkey = aux_signingkey2.publicKey.substring(2)

    //Setting the local/sessionStorage variables
    localStorage.setItem("kantcoin_org_words", words)
    sessionStorage.setItem("kantcoin_org_user", user)
    sessionStorage.setItem("kantcoin_org_key",keyjson)
    sessionStorage.setItem("kantcoin_org_aux_pubkey",aux_pkey)

    //Creating Ekhash in order to check user login and password
    let ekhash = "" + hashCode(keyjson + user)
    ekhash = ekhash.substr(ekhash.length > 4 ? 4 : 0)

    //Setting the ekhash in the localStorage
    localStorage.setItem("kantcoin_org_ekhash", ekhash)

    //Setting the random secrets_base
    localStorage.setItem("kantcoin_org_secrets_base", randomString(45))

    //Generating new wallet from the mnemonics, user, and password
    ethers.Wallet.fromBrainWallet(words, user + password3).then(function(the_wallet) {
        wallet = the_wallet
        sessionStorage.setItem("kantcoin_org_wallet", JSON.stringify(wallet))

        login_message.innerHTML = klang.login_message
        $('#modal5').modal("open")
    })

    //Cleaning the fields
    kantcoin_org_password3.value = ""
    kantcoin_org_password4.value = ""
    kantcoin_org_user3.value = ""
}

/**
 *  It sends the public key and the user id to the login provider
 */
function sendVault(){
    let user = sessionStorage.getItem("kantcoin_org_user")

    if(!user || !wallet.address){
        Materialize.toast(klang.no_vault_opened, 3000, 'rounded')
        return
    }

    if (login_provider_input.text){
        login_provider = login_provider_input.text
    }

    let x = (screen.width / 2) - 220
    let y = (screen.height / 2) - 300

    child = window.open("", "_blank", "width=440,height=600,top=" + y + ",left=" + x +
        ",resizable=no,status=no,menubar=no,scrollbars=no,titlebar=no,toolbar=no")
    if (child.location.href){
        child.location.href = login_provider + "/newUser?pkey=" + wallet.address + "&id=" + user
    } else {
        child.location = login_provider + "/newUser?pkey=" + wallet.address + "&id=" + user
    }

    $('#modal2').modal('open')
}

/**
 * It generates a file with the user's mnemonics and a hash
 * @param {string} filename
 */
function saveVault(){
    //Getting the mnemonics from local storage
    let words = localStorage.getItem("kantcoin_org_words")
    //The data to be saved on the file
    let data = ""
    //The data must be encrypted
    let enc = ""
    //The chosen file
    let file = file_name1.value

    //The user login and the key are necessary to generate the ekhash
    let user =  sessionStorage.getItem("kantcoin_org_user")
    let key = sessionStorage.getItem("kantcoin_org_key")

    //Checking if there are words to be saved and if the file name was provided
    if (!words || !file || !user || !key){
        Materialize.toast(klang.no_vault_saved, 2000, 'rounded')

        sleep(500).then(function(){
            $('#modal2').modal('open')
        })
        return
    }

    //Creating Ekhash in order to check user login and password
    let ekhash = "" + hashCode(key + user)
    ekhash = ekhash.substr(ekhash.length > 4 ? 4 : 0)

    //The data should be formatted as a JSON string
    data = '{'
        + '"login_provider":' + '"' + login_provider + '",'
        + '"ekhash":' + '"' + ekhash + '",'
        + '"secrets_base":' + '"' + localStorage.getItem("kantcoin_org_secrets_base") + '"'
        + '}'

    //Encrypting the data
    encrypt(words, data).then(function(value){
        enc = value

        //Removing the prefix in order to difficult file location by this string
        //Saving the mnemonics and data
        let blob = new Blob([words + "\r\n" + enc], {type: "text/plain;charset=utf-8"})
        let objectUrl = URL.createObjectURL(blob)

        file_saver.href = objectUrl
        file_saver.download = file
        file_saver.click()

        file_saver.href = ''
        file_saver.download = ''

        //Cleaning this field
        file_name1.value = ""

        window.location.href = (probablyTor() ? "http://" : "https://") + "kantcoin.org/home"
    })
}

/**
 * This function tells if the user is probably using the Tor browser
 * @returns {boolean}
 */
function probablyTor(){
    if (new Date().getTimezoneOffset() != 0){
        return false
    }
    if (navigator.plugins.length > 0){
        return false
    }
    if (!navigator.userAgent.startsWith("Mozilla")){
        return false
    }
    if (window.innerWidth != window.screen.width || window.innerHeight != window.screen.height){
        return false
    }
    return true
}

/******************** Tools ********************/

/**
 * Java-style hashcode
 * @param {string} str
 * @returns {number}
 */
function hashCode(str) {
    let hash = 0, i, chr
    if (str.length === 0) return hash
    for (i = 0; i < str.length; i++) {
        chr   = str.charCodeAt(i)
        hash  = ((hash << 5) - hash) + chr
        hash |= 0 // Convert to 32bit integer
    }
    //return only positive numbers
    return hash + 2147483648
}

/**
 * Sleep/Delay
 * @returns {Promise}
 */
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

/**
 * It generates a random string with N chars
 * @param {number} n
 * @returns {string}
 */
function randomString(n) {
    let text = ""
    let possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

    for (let i = 0; i < n; i++)
        text += possible.charAt(Math.floor(Math.random() * possible.length))

    return text
}

/**
 * It returns version of IE or false, if browser is not Internet Explorer
 */
function detectIE() {
    var ua = window.navigator.userAgent

    var msie = ua.indexOf('MSIE ')
    if (msie > 0) {
        // IE 10 or older => return version number
        return parseInt(ua.substring(msie + 5, ua.indexOf('.', msie)), 10)
    }

    var trident = ua.indexOf('Trident/')
    if (trident > 0) {
        // IE 11 => return version number
        var rv = ua.indexOf('rv:')
        return parseInt(ua.substring(rv + 3, ua.indexOf('.', rv)), 10)
    }

    var edge = ua.indexOf('Edge/')
    if (edge > 0) {
        // Edge (IE 12+) => return version number
        return parseInt(ua.substring(edge + 5, ua.indexOf('.', edge)), 10)
    }

    // other browser
    return false
}

/********************  AES encryption tools ********************/

/**
 * Encodes a utf8 string as a byte array.
 * @param {string} str
 * @returns {Uint8Array}
 */
function str2buf(str) {
    if (window.TextEncoder) {
        return new TextEncoder('utf-8').encode(str)
    }
    var utf8 = encodeURIComponent(str)
    var result = new Uint8Array(utf8.length);
    for (var i = 0; i < utf8.length; i++) {
        result[i] = utf8.charCodeAt(i)
    }
    return result
}

/**
 * Decodes a byte array as a utf8 string.
 * @param {Uint8Array} buffer
 * @returns {string}
 */
function buf2str(buffer) {
    if (window.TextDecoder) {
        return new TextDecoder("utf-8").decode(buffer)
    }
    var result = ""
    for (var i = 0; i < buffer.length; i++) {
        result += String.fromCharCode(buffer[i])
    }
    return result
}

/**
 * Decodes a string of hex to a byte array.
 * @param {string} hexStr
 * @returns {Uint8Array}
 */
function hex2buf(hexStr) {
    return new Uint8Array(hexStr.match(/.{2}/g).map(h => parseInt(h, 16)))
}

/**
 * Encodes a byte array as a string of hex.
 * @param {Uint8Array} buffer
 * @returns {string}
 */
function buf2hex(buffer) {
    return Array.prototype.slice
        .call(new Uint8Array(buffer))
        .map(x => [x >> 4, x & 15])
        .map(ab => ab.map(x => x.toString(16)).join(""))
        .join("")
}

/**
 * Given a passphrase, this generates a crypto key
 * using `PBKDF2` with SHA256 and 1000 iterations.
 * If no salt is given, a new one is generated.
 * The return value is an array of `[key, salt]`.
 * @param {string} passphrase
 * @param {UInt8Array} salt [salt=random bytes]
 * @returns {Promise<[CryptoKey,UInt8Array]>}
 */
function deriveKey(passphrase, salt) {
    salt = salt || crypto.getRandomValues(new Uint8Array(8))
    return crypto.subtle
        .importKey("raw", str2buf(passphrase), "PBKDF2", false, ["deriveKey"])
        .then(key =>
            crypto.subtle.deriveKey(
                { name: "PBKDF2", salt, iterations: 1000, hash: "SHA-256" },
                key,
                { name: "AES-GCM", length: 256 },
                false,
                ["encrypt", "decrypt"],
            ),
        )
        .then(key => [key, salt])
}

/**
 * Given a passphrase and some plaintext, this derives a key
 * (generating a new salt), and then encrypts the plaintext with the derived
 * key using AES-GCM. The ciphertext, salt, and iv are hex encoded and joined
 * by a "-". So the result is `"salt-iv-ciphertext"`.
 * @param {string} passphrase
 * @param {string} plaintext
 * @returns {Promise<string>}
 */
function encrypt(passphrase, plaintext) {
    const iv = crypto.getRandomValues(new Uint8Array(12))
    const data = str2buf(plaintext)
    return deriveKey(passphrase).then(([key, salt]) =>
        crypto.subtle
            .encrypt({ name: "AES-GCM", iv }, key, data)
            .then(ciphertext => `${buf2hex(salt)}-${buf2hex(iv)}-${buf2hex(ciphertext)}`),
    )
}

/**
 * Given a key and ciphertext (in the form of a string) as given by `encrypt`,
 * this decrypts the ciphertext and returns the original plaintext
 * @param {string} passphrase
 * @param {string} saltIvCipherHex
 * @returns {Promise<String>}
 */
function decrypt(passphrase, saltIvCipherHex) {
    const [salt, iv, data] = saltIvCipherHex.split("-").map(hex2buf)
    return deriveKey(passphrase, salt)
        .then(([key]) => crypto.subtle.decrypt({ name: "AES-GCM", iv }, key, data))
        .then(v => buf2str(new Uint8Array(v)))
}