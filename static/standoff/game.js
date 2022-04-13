import * as networking from '../networking.js';
import { AVATARS } from '../avatars/avatars.js';
import { COLORS, name } from '../landing.js';
import { appendDataLog, setChatboxNotification } from '../ingame-utility.js';
import { initTitles, initHowToPlays } from '../importantStrings.js';
import { playAudio } from '../audio.js';
import { en } from './enum.js';
import { Noneable } from '../noneable.js';

const ingameLeft = document.getElementById("ingame-left");
const endGameButton = document.getElementById("end-game-button");
const players = document.getElementById("players");
const chatLog = document.getElementById("chat-log");
const roundText = document.getElementById("round-text");
const startGameButton = document.getElementById("start-game-button");
const choices = document.getElementById("choices");
const choicesWaiting = document.getElementById("choices-waiting");
const outcome = document.getElementById("outcome");
const results = document.getElementById("results");
const gameoverHeader = document.getElementById("gameover-header");
const continueButton = document.getElementById("gameover-continue-button");
const gameoverResults = document.getElementById("gameover-results");

const startGameDiv = new Noneable(document.getElementById("start-game-div"));
const choicesContainer = new Noneable(document.getElementById("choices-container"));
const resultsDiv = new Noneable(document.getElementById("results-div"));
const resultsContinueButton = new Noneable(document.getElementById("results-continue-button"), true);
const roundTextDiv = new Noneable(document.getElementById("round-text-div"));
const gameover = new Noneable(document.getElementById("gameover"));

var handlers = {};

export function initMain(conn) {
    conn = conn;
    initTitles("Standoff");
    initHowToPlays("Shoot someone, or yourself, or nobody.\n\nIf anyone shoots you while you shoot yourself, they die instead of you.\n\nLast person standing wins.");

    startGameButton.onclick = function (e) {
        networking.send(conn, en.ToServerCode.START_GAME);
    };

    resultsContinueButton.element.onclick = function (e) {
        networking.send(conn, en.ToServerCode.PROMPT_REQUEST);
    };

    continueButton.onclick = function (e) {
        networking.send(conn, en.ToServerCode.START_GAME);
    };

    endGameButton.onclick = function (e) {
        networking.send(conn, en.ToServerCode.END_GAME);
    };

    handlers[en.ToClientCode.IN_MEDIA_RES] = (data) => {
        switch (data[0].charAt(0)) {
        case en.Phase.PREGAME:
            startGameDiv.show();
            break;
        case en.Phase.PLAY:
            SetChoicesWaitingStatus("spectating");
            choicesContainer.show();
            roundTextDiv.show();
            break;
        }
    }

    // ALL PHASES
    handlers[en.ToClientCode.PLAYERS] = (data) => {
        while (players.firstChild) {
            players.removeChild(players.firstChild);
        }
        for (let j = 0; j < data.length; j+=6) {
            var player = document.createElement("div");
            player.className = "player";
            var playerInfo = document.createElement("div");
            playerInfo.classList.add("player-info");
            playerInfo.innerHTML = "<b>" + data[j] + "</b>" + " " + data[j+3] + "<br/>" + data[j+4] + " wins";
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
    }

    handlers[en.ToClientCode.LOBBY_CHAT_MESSAGE] = (data) => {
        playAudio("glub");
        var item = networking.decodeToDiv(data[0]);
        appendDataLog(chatLog, item);
        var width = (window.innerWidth > 0) ? window.innerWidth : screen.width;
        if (!ingameLeft.classList.contains("ingame-left-expanded") && width <= 800) {
            setChatboxNotification(1);
        }
    }

    // PREGAME
    handlers[en.ToClientCode.START_GAME] = (data) => {
        playAudio("start");
        startGameDiv.hide();
        gameover.hide();
        choicesContainer.show();
        resultsContinueButton.show();
        roundTextDiv.show();
        while (gameoverResults.firstChild) {
            gameoverResults.removeChild(gameoverResults.firstChild);
        }
    }

    // PLAY
    handlers[en.ToClientCode.PROMPT] = (data) => {
        playAudio("dink2");
        roundText.innerText = `Round ${data[0]}`;
        resultsDiv.hide();
        outcome.innerText = "";
        results.innerText = "";
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
                    playAudio("click3")
                    networking.send(conn, en.ToServerCode.DECISION + clientId.toString());
                };
                choices.appendChild(item);
            }
            var item = document.createElement("button");
            item.innerText = "nobody";
            item.onclick = function() {
                playAudio("click3")
                networking.send(conn, en.ToServerCode.DECISION + "-2");
            };
            choices.appendChild(item);
        } else {
            SetChoicesWaitingStatus(data[1]);
        }
    }

    handlers[en.ToClientCode.DECISION_ACK] = (data) => {
        let childs = choices.children;
        for (let i=0; i<childs.length; i++) {
            childs[i].disabled = true;
        }
        choicesWaiting.innerText = "Waiting for other players.";
    }

    handlers[en.ToClientCode.RESULT] = (data) => {
        console.log("ayo");
        playAudio("whoosh");
        choicesWaiting.innerText = "";
        resultsDiv.show();
        outcome.innerText = "Round Outcome";
        while (results.firstChild) {
            results.removeChild(results.firstChild);
        }
        data.forEach(line => {
            var item = networking.decodeToDiv(line);
            results.appendChild(item);
        });
    }

    handlers[en.ToClientCode.WINNERS] = (data) => {
        playAudio("fanfare");
        choicesContainer.hide();
        roundTextDiv.hide();
        resultsContinueButton.hide();
        gameover.show();
        choicesWaiting.innerText = "";
        var item = document.createElement("p");
        let j = 1;
        if (data[0] === "1") {
            gameoverHeader.innerText = `Winner: ${data[1]}`
            item.innerHTML = `<b>${data[j]}</b> survived for ${data[j+1]} rounds.<br/>`;
            gameoverResults.appendChild(item);
            j += 3;
        } else {
            gameoverHeader.innerText = `Everyone died`
        }
        gameoverResults.appendChild(item);
        for (; j<data.length; j+=3) {
            var item = document.createElement("p");
            item.innerHTML = `<b>${data[j]}</b> survived for ${data[j+1]} rounds. Killed by: <b>${data[j+2]}</b><br/>`;
            gameoverResults.appendChild(item);
        }
    }

    handlers[en.ToClientCode.END_GAME] = (data) => {
        networking.send(conn, en.ToServerCode.END_GAME);
    }
    
    conn.onmessage = function (evt) {
        if (name === undefined) {
            return false;
        }
        const messages = evt.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            const m = messages[i];
            const toClientMessageCode = m.charAt(0);
            const data = networking.decode(m.substring(1,m.length));
            console.log("Got message: " + toClientMessageCode + " " + data);
            handlers[toClientMessageCode](data);
        }
    };
}

function SetChoicesWaitingStatus(status) {
    choicesWaiting.innerText = "You are " + status + ". Waiting for other players...";
};