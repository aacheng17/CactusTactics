import * as networking from './../networking.js';
import { avatars } from './../avatars/avatars.js';
import { initCollapsible } from './../collapsible.js';

// MESSAGE TYPES (CLIENT TO SERVER)
// when using conn.send(), the first character in the string represents the message type
// 0 means regular game message, 1 means setting player name, 2 means skip, 3 means game end / restart, 4 means what

window.onload = function () {
    initCollapsible();
    var col1 = document.getElementsByClassName("howtoplay-text");
    for (var i = 0; i < col1.length; i++) {
        col1[i].innerText = "Rules\nTry to think of a word that starts with the first letter and ends with the second letter before your opponents.\nWords must be at least 3 letters long.\n\nScoring\nThe more rare the letter combination, the more points it's worth (up to 100).\nEach word gets a length bonus multiplier as well."
    }

    var conn, name;
    var avatarIndex = 0;
    var colorIndex = 0;
    var colors = ["chocolate", "crimson", "coral", "gold", "darkgreen", "springgreen", "turquoise", "cornflowerblue", "indigo", "orchid", "slategrey", "black"];
    var landing = document.getElementById("landing");
    var nameForm = document.getElementById("name-form");
    var nameField = document.getElementById("name-field");
    var avatarRandomize = document.getElementById("avatar-randomize");
    var avatarButtonLeft = document.getElementById("avatar-button-left");
    var avatarButtonRight = document.getElementById("avatar-button-right");
    var avatarButtonColorLeft = document.getElementById("avatar-button-color-left");
    var avatarButtonColorRight = document.getElementById("avatar-button-color-right");
    var avatarSvg = document.getElementById("avatar-svg");
    var avatarPath = document.getElementById("avatar-path");
    var ingame = document.getElementById("ingame");
    var ingameLeft = document.getElementById("ingame-left");
    var leftExpandButton = document.getElementById("ingame-left-expand-button");
    var ingameHowtoplayButton = document.getElementById("ingame-howtoplay-button");
    var ingameHowtoplay = document.getElementById("ingame-howtoplay");
    var endgame = document.getElementById("endgame");
    var players = document.getElementById("players");
    var chatLog = document.getElementById("chat-log");
    var chatForm = document.getElementById("chat-form");
    var chatField = document.getElementById("chat-field");
    var startLetter = document.getElementById("start-letter");
    var endLetter = document.getElementById("end-letter");
    var skip = document.getElementById("skip")
    var promptExtraText = document.getElementById("prompt-extra-text");
    var gameLog = document.getElementById("game-log");
    var gameForm = document.getElementById("game-form");
    var gameField = document.getElementById("game-field");

    function getRandomInt(max) {
        return Math.floor(Math.random() * max);
    }

    function randomizeAvatar() {
        avatarIndex = getRandomInt(avatars.length);
        setAvatarSvg();
        colorIndex = getRandomInt(colors.length);
        avatarSvg.style.fill = colors[colorIndex];
    }

    function setAvatarSvg() {
        avatarPath.setAttribute("d", avatars[avatarIndex]);
    }

    function appendDataLog(log, item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    if (nameForm !== null) {
        nameForm.onsubmit = function (e) {
            if (!conn) {
                return false;
            }
            if (!nameField.value.trim()) {
                return false;
            }
            name = nameField.value;
            networking.send(conn, "1" + name + "\t" + avatarIndex.toString() + "\t" + colorIndex);
            e.preventDefault();
            landing.parentNode.removeChild(landing);
            document.body.appendChild(ingame);
        }
    }

    avatarRandomize.onclick = function(e) {
        randomizeAvatar();
    }

    avatarButtonLeft.onclick = function(e) {
        avatarIndex--;
        if (avatarIndex < 0) avatarIndex = avatars.length - 1;
        setAvatarSvg();
    }

    avatarButtonRight.onclick = function(e) {
        avatarIndex++;
        if (avatarIndex >= avatars.length) avatarIndex = 0;
        setAvatarSvg();
    }
    
    avatarButtonColorLeft.onclick = function(e) {
        colorIndex--;
        if (colorIndex < 0) colorIndex = colors.length - 1;
        avatarSvg.style.fill = colors[colorIndex];
    }
    
    avatarButtonColorRight.onclick = function(e) {
        colorIndex++;
        if (colorIndex >= colors.length) colorIndex = 0;
        avatarSvg.style.fill = colors[colorIndex];
    }

    ingameHowtoplayButton.addEventListener("click", function() {
        var effected = ingameHowtoplay;
        if (effected.style.maxHeight){
            effected.style.maxHeight = null;
        } else {
            effected.style.maxHeight = effected.scrollHeight + "px";
        } 
    });

    leftExpandButton.addEventListener("click", function() {
        var effected = ingameLeft;
        if (effected.style.left === "-250px" || effected.style.left === ""){
            leftExpandButton.firstChild.innerText = "Less";
            effected.style.left = "0px";
            effected.style.boxShadow = "1px 1px black"
            effected.boxShadow
        } else {
            leftExpandButton.firstChild.innerText = "More";
            effected.style.left = "-250px";
            effected.style.boxShadow = null;
        } 
    });

    endgame.onclick = function (e) {
        if (!conn) {
            return false;
        }
        networking.send(conn, "3");
    }

    chatForm.onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!chatField.value.trim()) {
            return false;
        }
        networking.send(conn, "8" + chatField.value);
        chatField.value = "";
        return false;
    };

    skip.onclick = function (e) {
        if (!conn) {
            return false;
        }
        networking.send(conn, "2");
    }
    
    gameForm.onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!gameField.value.trim()) {
            return false;
        }
        networking.send(conn, "0" + gameField.value);
        gameField.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        var websocketprotocol = "ws://";
        if (document.location.protocol==="https:") {
            websocketprotocol = "wss://";
        }
        conn = new WebSocket(websocketprotocol + document.location.host + document.location.pathname);
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendDataLog(gameLog, item);
        };
        conn.onmessage = function (evt) {
            if (name === undefined) {
                return false;
            }
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var m = messages[i];
                var messageType = m.charAt(0);
                var data = networking.decode(m.substring(1,m.length));
                switch (messageType) {
                case '0':
                    var item = networking.decodeToDiv(data[0]);
                    appendDataLog(gameLog, item);
                    break;
                case '1':
                    while (players.firstChild) {
                        players.removeChild(players.firstChild);
                    }
                    for (var j = 0; j < data.length; j+=6) {
                        var player = document.createElement("div");
                        player.className = "player";
                        var playerInfo = document.createElement("div");
                        playerInfo.classList.add("player-info");
                        var text = "<b>" + data[j] + "</b>" + ": " + data[j+3].toString() + " points";
                        if (data[j+4] != "") {
                            text += "<br/>Best: " + data[j+4] + " (" + data[j+5].toString() + ")";
                        }
                        playerInfo.innerHTML = text;
                        player.appendChild(playerInfo);
                        var svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
                        svg.classList.add("player-avatar");
                        svg.setAttribute("width", "50px");
                        svg.setAttribute("height", "50px");
                        svg.setAttribute("viewBox", "0 0 1000 1000");
                        svg.setAttribute("fill", colors[data[j+2]]);
                        var path = document.createElementNS("http://www.w3.org/2000/svg", "path");
                        svg.appendChild(path);
                        path.setAttribute("d", avatars[data[j+1]]);
                        player.appendChild(svg);
                        players.appendChild(player);
                    }
                    break
                case '2':
                    var sl = data[0].toUpperCase();
                    var el = data[1].toUpperCase();
                    startLetter.innerText = sl;
                    endLetter.innerText = el;
                    gameField.placeholder = sl + "___" + el;
                    promptExtraText.innerText = "Worth " + String(data[2]) + " points. There are " + String(data[3]) + " possible words.";
                    break
                case '3':
                    var item = document.createElement("div");
                    item.innerText = "Winner: " + data[0] + " " + data[1] + " points\nBest word: " + data[2] + " " + data[3] + " " + data[4] + " points";
                    appendDataLog(gameLog, item);
                    break;
                case '4':
                    endgame.innerText = "end game";
                    while (gameLog.firstChild) {
                        gameLog.removeChild(gameLog.firstChild);
                    }
                    break
                case '5':
                    var item = document.createElement("div");
                    item.classList.add("score-message");
                    var message = networking.decodeToDiv(data[0]);
                    item.appendChild(message);
                    var what = document.createElement("button");
                    what.classList.add("what-button");
                    what.innerText = "What?";
                    what.onclick = function (e) {
                        if (!conn) {
                            return false;
                        }
                        if (what.previousElementSibling === undefined) {
                            return false;
                        }
                        if (what.previousElementSibling.children.length < 1) {
                            return false;
                        }
                        networking.send(conn, "4" + what.previousElementSibling.children[0].id);
                    }
                    item.appendChild(what);
                    appendDataLog(gameLog, item);
                    break
                case '6':
                    var children = gameLog.children;
                    var found = false;
                    var child = null;
                    for (var i = 0; i < children.length; i++) {
                        child = children[i];
                        var childsChildren = child.children;
                        if (childsChildren.length !== 2) continue;
                        var potentialWhatButton = childsChildren[childsChildren.length - 1];
                        if (potentialWhatButton.classList.contains("what-button")) {
                            if (potentialWhatButton.previousElementSibling.children[0].id === data[1]) {
                                found = true;
                                potentialWhatButton.parentNode.removeChild(potentialWhatButton);
                                break;
                            }
                        }
                    }
                    var doScroll = gameLog.scrollTop > gameLog.scrollHeight - gameLog.clientHeight - 1;
                    var item = networking.decodeToDiv(data[0]);
                    if (found) {
                        child.insertAdjacentElement('afterend', item);
                    } else {
                        gameLog.prepend(item);
                    }
                    if (doScroll) {
                        gameLog.scrollTop = gameLog.scrollHeight - gameLog.clientHeight;
                    }
                    break
                case '7':
                    endgame.innerText = "new game";
                    var item = networking.decodeToDiv(data[0]);
                    appendDataLog(gameLog, item);
                    break;
                case '8':
                    var item = networking.decodeToDiv(data[0]);
                    appendDataLog(chatLog, item);
                    break;      
                }
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendDataLog(gameLog, item);
    }

    ingame.parentNode.removeChild(ingame);
    randomizeAvatar();
};