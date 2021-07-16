import * as networking from '../networking.js';
import { AVATARS } from '../avatars/avatars.js';
import { COLORS, name } from '../landing.js';
import { appendDataLog, setChatboxNotification } from '../ingame-utility.js';

var ingameLeft = document.getElementById("ingame-left");
var endgame = document.getElementById("endgame");
var players = document.getElementById("players");
var chatLog = document.getElementById("chat-log");

export function initMain(conn) {    
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
            case '0':
                endgame.innerText = "end game";
                while (gameLog.firstChild) {
                    gameLog.removeChild(gameLog.firstChild);
                }
                break
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
                break
            case 'a':
                break;
            }
        }
    };
}