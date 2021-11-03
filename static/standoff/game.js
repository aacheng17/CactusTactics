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
var resultsDiv = document.getElementById("results-div");
var outcome = document.getElementById("outcome");
var results = document.getElementById("results");
var continueDiv = document.getElementById("continue");
var gameResultsHeader = document.getElementById("game-results-header");
var gameResults = document.getElementById("game-results");

export function initMain(conn) {
    initTitles("Standoff");
    initHowToPlays("Shoot someone, or yourself, or nobody.\n\nIf anyone shoots you while you shoot yourself, they die instead of you.\n\nLast person standing wins.");

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
                gameResultsHeader.innerText = "";
                while (gameResults.firstChild) {
                    gameResults.removeChild(gameResults.firstChild);
                }
                endgame.disabled = true;
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
                endgame.disabled = false;
                choices.innerText = "Waiting for new game...";
                break;
            case '3':
                while (players.firstChild) {
                    players.removeChild(players.firstChild);
                }
                for (let j = 0; j < data.length; j+=6) {
                    var player = document.createElement("div");
                    player.className = "player";
                    var playerInfo = document.createElement("div");
                    playerInfo.classList.add("player-info");
                    playerInfo.innerHTML = "<b>" + data[j] + "</b>" + " " + data[j+3] + "<br/>" + data[j+4] + " kills";
                    playerInfo
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
                    switch (data[j+5]) {
                        case "dotdotdot": playerStatus.classList.add("dotdotdot"); break;
                        case "ready": playerStatus.innerHTML = "&#10003;"; break;
                        case "none": break;
                    }
                    playerAvatarContainer.appendChild(playerStatus);
                    player.appendChild(playerAvatarContainer);
                    players.appendChild(player);
                }
                break;
            case 'a':
                endgame.disabled = true;
                roundText.innerText = `Round ${data[0]}`;
                resultsDiv.style.display = "none";
                outcome.innerText = "";
                results.innerText = "";
                while (continueDiv.firstChild) {
                    continueDiv.removeChild(continueDiv.firstChild);
                }
                while (choices.firstChild) {
                    choices.removeChild(choices.firstChild);
                }
                if (data.length > 2) {
                    choices.innerText = "Who to shoot?";
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
                    var item = document.createElement("button");
                    item.innerText = "nobody";
                    item.onclick = function() {
                        networking.send(conn, "a-2");
                    };
                    choices.appendChild(item);
                } else {
                    choicesWaiting.innerText = "You are " + data[1] + ". Waiting for other players...";
                }
                break;
            case 'b':
                let childs = choices.children;
                for (let i=0; i<childs.length; i++) {
                    childs[i].disabled = true;
                }
                choicesWaiting.innerText = "Waiting for other players.";
                break;
            case 'c':
                choicesWaiting.innerText = "";
                resultsDiv.style.display = "flex";
                resultsDiv.style.flexDirection = "column";
                outcome.innerText = "Round Outcome";
                data.forEach(line => {
                    item = networking.decodeToDiv(line);
                    results.appendChild(item);
                });
                while (continueDiv.firstChild) {
                    continueDiv.removeChild(continueDiv.firstChild);
                }
                item = document.createElement("button");
                item.innerText = "Continue";
                item.onclick = function() {
                    networking.send(conn, "b");
                };
                continueDiv.appendChild(item);
                break;
            case 'd':
                while (continueDiv.firstChild) {
                    continueDiv.removeChild(continueDiv.firstChild);
                }
                choicesWaiting.innerText = "";
                gameResultsHeader.innerText = "Game Outcome"
                data.forEach(line => {
                    item = networking.decodeToDiv(line);
                    gameResults.appendChild(item);
                });
                break;
            }
        }
    };
}