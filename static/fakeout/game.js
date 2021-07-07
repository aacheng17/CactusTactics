import * as networking from '../networking.js';
import { AVATARS } from '../avatars/avatars.js';
import { COLORS, name } from '../landing.js';
import { appendDataLog, setChatboxNotification } from '../ingame-utility.js';

var phase = 0;
var ingameLeft = document.getElementById("ingame-left");
var endgame = document.getElementById("endgame");
var players = document.getElementById("players");
var chatLog = document.getElementById("chat-log");
var promptText = document.getElementById("prompt-text");
var promptSubmit = document.getElementById("prompt-submit");
var promptForm = document.getElementById("prompt-form");
var promptField = document.getElementById("prompt-field");
var promptWaiting = document.getElementById("prompt-waiting");
var choices = document.getElementById("choices");
var choicesWaiting = document.getElementById("choices-waiting");

export function initMain(conn) {
    promptForm.onsubmit = function (e) {
        e.preventDefault();
        if (!conn) {
            return false;
        }
        if (!promptField.value.trim()) {
            return false;
        }
        networking.send(conn, "a" + promptField.value);
        promptField.value = "";
        return false;
    };

    conn.onmessage = function (evt) {
        if (name === undefined) {
            return false;
        }
        var messages = evt.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            var m = messages[i];
            console.log(m);
            var messageType = m.charAt(0);
            var data = networking.decode(m.substring(1,m.length));
            switch (messageType) {
            case '0':
                endgame.innerText = "end game";
                break;
            case '1':
                var item = networking.decodeToDiv(data[0]);
                appendDataLog(chatLog, item);
                var width = (window.innerWidth > 0) ? window.innerWidth : screen.width;
                if (!ingameLeft.classList.contains("ingame-left-expanded") && width <= 800) {
                    setChatboxNotification(1);
                }
                break;
            case '2':
                endgame.innerText = "new game";
                var item = networking.decodeToDiv(data[0]);
                appendDataLog(chatLog, item);
                break;
            case '3':
                while (players.firstChild) {
                    players.removeChild(players.firstChild);
                }
                for (let j = 0; j < data.length; j+=4) {
                    var player = document.createElement("div");
                    player.className = "player";
                    var playerInfo = document.createElement("div");
                    playerInfo.classList.add("player-info");
                    var text = "<b>" + data[j] + "</b>" + ": " + data[j+3].toString() + " points";
                    playerInfo.innerHTML = text;
                    player.appendChild(playerInfo);
                    var svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
                    svg.classList.add("player-avatar");
                    svg.setAttribute("width", "50px");
                    svg.setAttribute("height", "50px");
                    svg.setAttribute("viewBox", "0 0 1000 1000");
                    svg.setAttribute("fill", COLORS[data[j+2]]);
                    var path = document.createElementNS("http://www.w3.org/2000/svg", "path");
                    svg.appendChild(path);
                    path.setAttribute("d", AVATARS[data[j+1]]);
                    player.appendChild(svg);
                    players.appendChild(player);
                }
                break;
            case 'f':
                var item = document.createElement("div");
                item.innerText = "Winner: " + data[0] + " " + data[1] + " points\nBest word: " + data[2] + " " + data[3] + " " + data[4] + " points";
                appendDataLog(chatLog, item);
                break;
            case 'a':
                while (promptText.firstChild) {
                    promptText.removeChild(promptText.firstChild);
                }
                promptText.appendChild(networking.decodeToDiv(data[0]));
                break;
            case 'b':
                var err = parseInt(data[0]);
                switch (err) {
                case -1:
                    promptWaiting.innerText = "Your answer is too close to the actual answer. Pick something else"
                    break;
                default:
                    promptWaiting.innerText = "Response submitted. Waiting for other players."
                    promptField.disabled = true;
                    promptSubmit.disabled = true;
                }
                break;
            case 'c':
                promptField.disabled = true;
                promptSubmit.disabled = true;
                promptWaiting.innerText = "";
                for (let j = 0; j < data.length; j++) {
                    let d = data[j];
                    var item = document.createElement("button");
                    item.innerText = d;
                    item.onclick = function() {
                        networking.send(conn, "b" + j.toString());
                    };
                    choices.append(item);
                }
                break;
            case 'd':
                var err = parseInt(data[0]);
                switch (err) {
                case -1:
                    choicesWaiting.innerText = "You can't pick your own answer."
                    break;
                default:
                    choicesWaiting.innerText = "Choice submitted. Waiting for other players."
                    for (let j=0; j<choices.childElementCount; j++) {
                        choices.children[j].disabled = true;
                    }
                }
                break;
            case 'e':
                choices.classList.add("revealed");
                var curChild = choices.firstChild;
                for (let j=0; j<data.length/2; j++) {
                    var faker = data[j*2];
                    var fakedOut = data[j*2+1];
                    var item = document.createElement("p");
                    if (faker === "") {
                        item.innerText += "ACTUAL ANSWER";
                        if (fakedOut !== "") {
                            item.innerText += " picked by " + fakedOut;
                        }
                    } else {
                        item.innerText += faker;
                        if (fakedOut !== "") {
                            item.innerText += " faked out " + fakedOut;
                        }
                    }
                    choices.insertBefore(item, curChild);
                    curChild = curChild.nextElementSibling;
                }
                choicesWaiting.innerText = "";
                var item = document.createElement("button");
                item.innerText = "Continue";
                item.onclick = function() {
                    promptField.disabled = false;
                    promptSubmit.disabled = false;
                    choices.classList.remove("revealed");
                    while (choices.firstChild) {
                        choices.removeChild(choices.firstChild);
                    }
                    while (choicesWaiting.firstChild) {
                        choicesWaiting.removeChild(choicesWaiting.firstChild);
                    }
                    networking.send(conn, "c");
                };
                choicesWaiting.appendChild(item);
                break;
            }
        }
    };
}