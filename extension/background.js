var current_tab = {}
const PATH_TO_INDEX_HTML = "index.html" //http://localhost:1985

if (chrome && chrome.omnibox){
    chrome.omnibox.onInputEntered.addListener(omniSearch)

    //Probably the most relevant campaigns will be recent ones
    chrome.omnibox.onInputChanged.addListener(
        function(text, suggest)
        {
            text = text.replace(" ", "");
            var suggestions = []
            suggestions.push({ content: "" + Math.floor(Date.now() / 100000000) + text, description: "" + Math.floor(Date.now() / 100000000) + "..." })
            suggestions.push({ content: "custom_provider=" + text, description: "Custom provider" })
            suggest(suggestions)
        }
    )
}

if (chrome && chrome.tabs){
    chrome.tabs.onRemoved.addListener(function (tabId, removeInfo) {
        if (tabId == current_tab.id){
            current_tab = {}
        }
    })
}

if (chrome && chrome.browserAction){
    chrome.browserAction.onClicked.addListener(function(activeTab){
        var newURL = PATH_TO_INDEX_HTML

        if (current_tab && current_tab.id){
            chrome.tabs.update(current_tab.id, {selected: true, url: newURL})
        } else {
            chrome.tabs.create({ url: newURL },function (tab) {
                current_tab = tab
            })
        }
    })
}


//Searching via omnibox
function omniSearch (query) {
    if (query.indexOf("custom_provider=") == 0){
        var newURL = PATH_TO_INDEX_HTML + "?" + query
    } else {
        var newURL = PATH_TO_INDEX_HTML + "?q=" + query
    }

    if (current_tab && current_tab.id){
        chrome.tabs.update(current_tab.id, {selected: true, url: newURL})
    } else {
        chrome.tabs.getSelected(function (tab) {
            chrome.tabs.update(tab.id, {url: newURL})
            current_tab = tab
        })
    }
}