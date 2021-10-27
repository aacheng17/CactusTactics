var sounds = {
    "bubble": gen("bubble"),
    "tap": gen("tap"),
}

function gen(audioName) {
    return new Audio("./static/audio/" + audioName + ".mp3");
}

export function playAudio(audioName) {
    sounds[audioName].play();
}