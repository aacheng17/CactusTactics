export function initHowToPlays(text) {
    var howtoplays = document.getElementsByClassName("howtoplay-text");
    for (var i = 0; i < howtoplays.length; i++) {
        howtoplays[i].innerText = text;
    }
}