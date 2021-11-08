const AudioContext = window.AudioContext || window.webkitAudioContext;
var audioContext;

export var sounds;

export function initAudio() {
    audioContext = new AudioContext();
    audioContext.resume();
    sounds = {
        "bubble": gen("bubble"),
        "tap": gen("tap"),
        "click1": gen("click1"),
        "click2": gen("click2"),
        "correct": gen("correct"),
        "fanfare": gen("fanfare"),
        "start": gen("start"),
    }
}

function gen(audioName) {
    let sound = new Audio("./static/audio/" + audioName + ".ogg");
    let source = audioContext.createMediaElementSource(sound);
    source.connect(audioContext.destination);
    return sound;
}

export function playAudio(sound) {
    sounds[sound].play();
}