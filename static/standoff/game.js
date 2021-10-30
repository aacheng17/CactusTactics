import * as networking from '../networking.js';
import { AVATARS } from '../avatars/avatars.js';
import { COLORS, name } from '../landing.js';
import { appendDataLog, setChatboxNotification } from '../ingame-utility.js';
import { initTitles, initHowToPlays } from '../importantStrings.js';

var ingameLeft = document.getElementById("ingame-left");
var endgame = document.getElementById("endgame");
var players = document.getElementById("players");
var chatLog = document.getElementById("chat-log");
var roundText = document.getElementById("round-text");
var choices = document.getElementById("choices");
var choicesWaiting = document.getElementById("choices-waiting");
var results = document.getElementById("results");

export function initMain(conn) {
    initTitles("Standoff");
    initHowToPlays("Shoot someone, or shoot yourself.\n\nIf anyone shoots you while you shoot yourself, they die instead of you.\n\nLast person standing wins.");

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
                for (let j = 0; j < data.length; j+=5) {
                    var player = document.createElement("div");
                    player.className = "player";
                    var playerInfo = document.createElement("div");
                    playerInfo.classList.add("player-info");
                    var text = "<b>" + data[j] + "</b>" + " "
                    if (data[j+3].toString() === "false") {
                        text += "spectating";
                    } else if (data[j+4].toString() === "true") {
                        text += "alive";
                    } else {
                        text += "dead";
                    }
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
            case 'a':
                results.innerText = "";
                roundText.innerText = `Round ${data[0]}\nWho to shoot?`;
                while (choices.firstChild) {
                    choices.removeChild(choices.firstChild);
                }
                for (let j = 1; j < data.length; j+=2) {
                    let clientId = data[j];
                    let clientName = data[j+1]
                    var item = document.createElement("button");
                    item.innerText = clientName;
                    item.onclick = function() {
                        networking.send(conn, "a" + clientId.toString());
                    };
                    choices.appendChild(item);
                }
                break;
            case 'b':
                let childs = choices.children;
                for (let i=0; i<childs.length; i++) {
                    childs[i].disabled = true;
                }
                choicesWaiting.innerText = "Waiting for other players."
                break;
            case 'c':
                choicesWaiting.innerText = "";
                results.innerText = "Outcome:\n"
                var item = networking.decodeToDiv(data[0]);
                results.appendChild(item);
                item = document.createElement("button");
                item.innerText = "Continue";
                item.onclick = function() {
                    networking.send(conn, "b");
                };
                results.appendChild(item);
                break;
            case 'd':
                choicesWaiting.innerText = "";
                var item = networking.decodeToDiv(data[0]);
                results.appendChild(item);
                break;
            }
        }
    };
}