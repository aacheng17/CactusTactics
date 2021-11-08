const AudioContext = window.AudioContext || window.webkitAudioContext;
var audioContext;

export var sounds;

export function initAudio() {
    audioContext = new AudioContext();
    audioContext.resume();
    sounds = {};
    ["bikebell", "blupblup", "bubble", "click1", "click2", "click3", "correct", "ding1", "dink2", "fanfare", "glub", "start", "tap", "whoosh"].forEach(s => sounds[s] = gen(s));
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