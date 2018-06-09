/**
 * Kantcoin Project
 * https://kantcoin.org
 * This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
 * If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

if (localStorage.getItem("kantcoin_org_words")){
    window.location.href = (probablyTor() ? "http://" : "https://") + "kantcoin.org/home"
} else {
    window.location.href = (probablyTor() ? "http://" : "https://") + "kantcoin.org/login"
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
    if (navigator.userAgent.indexOf("Firefox") == -1){
        return false
    }
    if (window.innerWidth != window.screen.width || window.innerHeight != window.screen.height){
        return false
    }
    return true
}