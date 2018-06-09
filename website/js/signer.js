/**
 * Kantcoin Project
 * https://kantcoin.org
 * This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
 * If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

importScripts("../urs/urs.min.js")

/**
 * It signs a vote
 * @param {Event} event
 */
onmessage = function(event) {
    let arguments_json = JSON.parse(event.data)

    //Arguments of "roda": verify-text, sign-text, keyring, keypair, signature, blind
    let signature = urs.roda("", arguments_json.vote_message, arguments_json.pubkeys, arguments_json.keypair, "", false)

    postMessage(signature)
}


