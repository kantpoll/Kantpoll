/**
 * Kantcoin Project
 * https://kantcoin.org
 * This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
 * If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

importScripts("../urs/urs.min.js")

/**
 * It receives a new vote to be processed
 * @param {Event} event
 */
onmessage = function(event){
    let vote_json = JSON.parse(event.data)

    if (vote_json.signature.indexOf(vote_json.first_number) == 0){
        let response = urs.roda(vote_json.vote_message, "", vote_json.pubkeys, "", vote_json.signature, false)
        vote_json.response = response
        postMessage(JSON.stringify(vote_json))
    }
}