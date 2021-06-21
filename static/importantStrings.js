function setStrings(className, text) {
    var elements = document.getElementsByClassName(className);
    for (var i = 0; i < elements.length; i++) {
        elements[i].innerText = text;
    }    
}

export function initTitles(text) {
    setStrings("title", text);
}

export function initHowToPlays(text) {
    setStrings("howtoplay-text", text);
}