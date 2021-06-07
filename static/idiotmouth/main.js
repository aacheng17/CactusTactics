import * as networking from './../networking.js';
import { initCollapsible } from './../collapsible.js';

// MESSAGE TYPES (CLIENT TO SERVER)
// when using conn.send(), the first character in the string represents the message type
// 0 means regular chat message, 1 means setting player name, 2 means skip, 3 means game end / restart, 4 means what

window.onload = function () {
    initCollapsible();
    var col1 = document.getElementsByClassName("howtoplay-text");
    for (var i = 0; i < col1.length; i++) {
        col1[i].innerText = "Rules\nTry to think of a word that starts with the first letter and ends with the second letter before your opponents.\nWords must be at least 3 letters long.\n\nScoring\nThe more rare the letter combination, the more points it's worth (up to 100).\nEach word gets a length bonus multiplier as well."
    }

    var conn, name;
    var landing = document.getElementById("landing");
    var nameForm = document.getElementById("name-form");
    var nameField = document.getElementById("name-field");
    var ingame = document.getElementById("ingame");
    var ingameHowtoplayButton = document.getElementById("ingame-howtoplay-button");
    var ingameHowtoplay = document.getElementById("ingame-howtoplay");
    var endgame = document.getElementById("endgame");
    var players = document.getElementById("players");
    var startLetter = document.getElementById("start-letter");
    var endLetter = document.getElementById("end-letter");
    var skip = document.getElementById("skip")
    var promptExtraText = document.getElementById("prompt-extra-text");
    var chatLog = document.getElementById("chat-log");
    var chatForm = document.getElementById("chat-form");
    var chatField = document.getElementById("chat-field");

    function appendChatLog(item) {
        var doScroll = chatLog.scrollTop > chatLog.scrollHeight - chatLog.clientHeight - 1;
        chatLog.appendChild(item);
        if (doScroll) {
            chatLog.scrollTop = chatLog.scrollHeight - chatLog.clientHeight;
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
            networking.send(conn, "1" + name);
            e.preventDefault();
            landing.parentNode.removeChild(landing);
            ingame.style.visibility = "visible";
        }
    }

    ingameHowtoplayButton.addEventListener("click", function() {
        if (ingameHowtoplay.style.maxHeight){
            ingameHowtoplay.style.maxHeight = null;
        } else {
            ingameHowtoplay.style.maxHeight = ingameHowtoplay.scrollHeight + "px";
        } 
    });

    endgame.onclick = function (e) {
        if (!conn) {
            return false;
        }
        networking.send(conn, "3");
        endgame.innerText = endgame.innerText === "end game" ? "new game" : "end game";
    }

    skip.onclick = function (e) {
        if (!conn) {
            return false;
        }
        networking.send(conn, "2");
    }
    
    chatForm.onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!chatField.value.trim()) {
            return false;
        }
        networking.send(conn, "0" + chatField.value);
        chatField.value = "";
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
            appendChatLog(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var m = messages[i];
                var messageType = m.charAt(0);
                var data = networking.decode(m.substring(1,m.length));
                switch (messageType) {
                case '0':
                    var item = document.createElement("div");
                    item.appendChild(networking.decodeToHTML(data[0]));
                    appendChatLog(item);
                    break;
                case '1':
                    while (players.firstChild) {
                        players.removeChild(players.firstChild);
                    }
                    for (var j = 0; j < data.length; j+=4) {
                        var player = document.createElement("div");
                        player.id = "player";
                        var text = data[j] + " - " + data[j+1].toString() + " points";
                        if (data[j+2] != "") {
                            text += "\nBest word: " + data[j+2] + " " + data[j+3].toString() + " points";
                        }
                        player.innerText = text;
                        players.appendChild(player);
                    }
                    break
                case '2':
                    var sl = data[0].toUpperCase();
                    var el = data[1].toUpperCase();
                    startLetter.innerText = sl;
                    endLetter.innerText = el;
                    chatField.placeholder = "a word that starts with " + sl + " and ends with " + el;
                    promptExtraText.innerText = "Worth " + String(data[2]) + " points. There are " + String(data[3]) + " possible words.";
                    break
                case '3':
                    var item = document.createElement("div");
                    item.innerText = "Winner: " + data[0] + " " + data[1] + " points\nBest word: " + data[2] + " " + data[3] + " " + data[4] + " points";
                    appendChatLog(item);
                    break;
                case '4':
                    while (chatLog.firstChild) {
                        chatLog.removeChild(chatLog.firstChild);
                    }
                    break
                case '5':
                    var item = document.createElement("div");
                    item.classList.add("score-message");
                    var message = networking.decodeToHTML(data[0]);
                    item.appendChild(message);
                    var word = networking.decodeToHTML(data[1]);
                    item.appendChild(word);
                    var what = document.createElement("button");
                    what.classList.add("what-button");
                    what.innerText = "What?";
                    what.onclick = function (e) {
                        if (!conn) {
                            return false;
                        }
                        networking.send(conn, "4" + what.previousElementSibling.innerText);
                    }
                    item.appendChild(what);
                    appendChatLog(item);
                    break
                }
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendChatLog(item);
    }

    ingame.style.visibility = "hidden";
};