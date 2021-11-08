import * as networking from '../networking.js';
import { AVATARS } from '../avatars/avatars.js';
import { COLORS, name } from '../landing.js';
import { appendDataLog, setChatboxNotification } from '../ingame-utility.js';
import { initTitles, initHowToPlays } from '../importantStrings.js';
import { playAudio } from '../audio.js';

var ingameLeft = document.getElementById("ingame-left");
var endgame = document.getElementById("endgame");
var players = document.getElementById("players");
var chatLog = document.getElementById("chat-log");
var startLetter = document.getElementById("start-letter");
var endLetter = document.getElementById("end-letter");
var skip = document.getElementById("skip")
var promptExtraText = document.getElementById("prompt-extra-text");
var gameLog = document.getElementById("game-log");
var gameForm = document.getElementById("game-form");
var gameField = document.getElementById("game-field");

export function initMain(conn) {
    initTitles("Idiotmouth");
    initHowToPlays("Rules\nTry to think of a word that starts with the first letter and ends with the second letter before your opponents.\nWords must be at least 3 letters long.\n\nScoring\nThe more rare the letter combination, the more points it's worth (up to 100).\nEach word gets a length bonus multiplier as well.");
    
    skip.onclick = function (e) {
        if (!conn) {
            return false;
        }
        networking.send(conn, "b");
    }
    
    gameForm.onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!gameField.value.trim()) {
            return false;
        }
        networking.send(conn, "a" + gameField.value);
        gameField.value = "";
        return false;
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
            case 'a':
                playAudio("click2");
                var item = networking.decodeToDiv(data[0]);
                appendDataLog(gameLog, item);
                break;
            case '3':
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
                    var playerAvatarContainer = document.createElement("div");
                    playerAvatarContainer.classList.add("player-avatar-container");
                    var svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
                    svg.classList.add("player-avatar");
                    svg.setAttribute("width", "50px");
                    svg.setAttribute("height", "50px");
                    svg.setAttribute("viewBox", "0 0 1000 1000");
                    svg.setAttribute("fill", COLORS[data[j+2]]);
                    var path = document.createElementNS("http://www.w3.org/2000/svg", "path");
                    svg.appendChild(path);
                    path.setAttribute("d", AVATARS[data[j+1]]);
                    playerAvatarContainer.appendChild(svg);
                    var playerStatus = document.createElement("a")
                    playerStatus.classList.add("player-status")
                    //playerStatus.innerHTML = "&#10003;";
                    playerAvatarContainer.appendChild(playerStatus);
                    player.appendChild(playerAvatarContainer);
                    players.appendChild(player);
                }
                break
            case 'd':
                var sl = data[0].toUpperCase();
                var el = data[1].toUpperCase();
                startLetter.innerText = sl;
                endLetter.innerText = el;
                gameField.placeholder = sl + "___" + el;
                promptExtraText.innerText = "Worth " + String(data[2]) + " points. There are " + String(data[3]) + " possible words.";
                break
            case 'e':
                var item = document.createElement("div");
                item.innerText = "Winner: " + data[0] + " " + data[1] + " points\nBest word: " + data[2] + " " + data[3] + " " + data[4] + " points";
                appendDataLog(gameLog, item);
                break;
            case '0':
                playAudio("start");
                endgame.innerText = "end game";
                while (gameLog.firstChild) {
                    gameLog.removeChild(gameLog.firstChild);
                }
                break
            case 'c':
                playAudio("correct");
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
                    networking.send(conn, "c" + what.previousElementSibling.children[0].id);
                }
                item.appendChild(what);
                appendDataLog(gameLog, item);
                break
            case 'b':
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
            case '2':
                playAudio("fanfare");
                endgame.innerText = "new game";
                var item = networking.decodeToDiv(data[0]);
                appendDataLog(gameLog, item);
                break;
            case '1':
                playAudio("bubble");
                var item = networking.decodeToDiv(data[0]);
                appendDataLog(chatLog, item);
                var width = (window.innerWidth > 0) ? window.innerWidth : screen.width;
                if (!ingameLeft.classList.contains("ingame-left-expanded") && width <= 800) {
                    setChatboxNotification(1);
                }
                break;
            }
        }
    };
}