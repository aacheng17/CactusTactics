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
const newGameMinWordLength = document.getElementById("new-game-min-word-length");
const newGameScoreToWin = document.getElementById("new-game-score-to-win");
const startGameButton = document.getElementById("start-game-button");
const startLetter = document.getElementById("start-letter");
const endLetter = document.getElementById("end-letter");
const skip = document.getElementById("skip")
const promptExtraText = document.getElementById("prompt-extra-text");
const gameLog = document.getElementById("game-log");
const gameForm = document.getElementById("game-form");
const gameField = document.getElementById("game-field");
const winners = document.getElementById("winners");
const newGameButton = document.getElementById("new-game-button");

const newgame = new Noneable(document.getElementById("new-game"));
const prompt = new Noneable(document.getElementById("prompt"));
const endgame = new Noneable(document.getElementById("end-game"));
const gamebox = new Noneable(document.getElementById("gamebox"));

var handlers = {};

export function initMain(conn) {
    initTitles("Idiotmouth");
    initHowToPlays("Rules\nTry to think of a word that starts with the first letter and ends with the second letter before your opponents.\nWords must be at least 3 letters long.\n\nScoring\nThe more rare the letter combination, the more points it's worth (up to 100).\nEach word gets a length bonus multiplier as well.");

    endGameButton.onclick = function (e) {
        networking.send(conn, en.ToServerCode.END_GAME);
    }

    newGameMinWordLength.onchange = function (e) {
        networking.send(conn, en.ToServerCode.GAMERULE_MIN_WORD_LENGTH + newGameMinWordLength.value.toString());
    }

    newGameScoreToWin.onchange = function (e) {
        networking.send(conn, en.ToServerCode.GAMERULE_SCORE_TO_WIN + newGameScoreToWin.value.toString());
    }

    startGameButton.onclick = function (e) {
        networking.send(conn, en.ToServerCode.START_GAME);
    }

    skip.onclick = function (e) {
        networking.send(conn, en.ToServerCode.VOTE_SKIP);
    }
    
    gameForm.onsubmit = function () {
        if (!gameField.value.trim()) {
            return false;
        }
        networking.send(conn, en.ToServerCode.GAME_MESSAGE + gameField.value);
        gameField.value = "";
        return false;
    };

    newGameButton.onclick = function(e) {
        playAudio("click2");
        newgame.show();
        endgame.hide();
        gamebox.hide();
    }

    handlers[en.ToClientCode.IN_MEDIA_RES] = (data) => {
        console.log(data[0].charAt(0));
        switch (data[0].charAt(0)) {
        case en.Phase.PREGAME:
            newgame.show();
            break;
        case en.Phase.PLAY:
            prompt.show();
            gamebox.show();
            break;
        }
    }
    
    // ALL PHASES
    handlers[en.ToClientCode.PLAYERS] = (data) => {
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
    
    //PREGAME
    handlers[en.ToClientCode.GAMERULE_MIN_WORD_LENGTH] = (data) => {
        playAudio("click2");
        newGameMinWordLength.value = data[0];
    }
    
    handlers[en.ToClientCode.GAMERULE_SCORE_TO_WIN] = (data) => {
        playAudio("click2");
        newGameScoreToWin.value = data[0];
    }
    
    handlers[en.ToClientCode.START_GAME] = (data) => {
        prompt.show();
        gamebox.show();
        newgame.hide();
        endgame.hide();
        playAudio("start");
        while (gameLog.firstChild) {
            gameLog.removeChild(gameLog.firstChild);
        }
    }
    
    //PLAY
    handlers[en.ToClientCode.GAME_MESSAGE] = (data) => {
        playAudio("click3");
        var item = networking.decodeToDiv(data[0]);
        appendDataLog(gameLog, item);
    }
    
    handlers[en.ToClientCode.PROMPT] = (data) => {
        playAudio("correct");
        var sl = data[0].toUpperCase();
        var el = data[1].toUpperCase();
        startLetter.innerText = sl;
        endLetter.innerText = el;
        gameField.placeholder = sl + "___" + el;
        promptExtraText.innerText = "Worth " + String(data[2]) + " points. There are " + String(data[3]) + " possible words.";
    }
    
    handlers[en.ToClientCode.MESSAGE_WITH_WHAT] = (data) => {
        var item = document.createElement("div");
        item.classList.add("score-message");
        var message = networking.decodeToDiv(data[0]);
        item.appendChild(message);
        var what = document.createElement("button");
        what.classList.add("what-button");
        what.innerText = "What?";
        what.onclick = function (e) {
            if (what.previousElementSibling === undefined) {
                return false;
            }
            if (what.previousElementSibling.children.length < 1) {
                return false;
            }
            networking.send(conn, en.ToServerCode.WHAT + what.previousElementSibling.children[0].id);
        }
        item.appendChild(what);
        appendDataLog(gameLog, item);
    }
    
    handlers[en.ToClientCode.WHAT_RESPONSE] = (data) => {
        playAudio("dink2");
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
    }
    
    handlers[en.ToClientCode.END_GAME] = (data) => {
        prompt.hide();
        endgame.show();
        playAudio("fanfare");
        winners.innerHTML = "Winner: <b>" + data[0] + "</b> (" + data[1] + " points)<br/>Best word: <b>" + data[2] + "</b> - " + data[3] + " (" + data[4] + " points)";
    }

    conn.onmessage = function (evt) {
        if (name === undefined) {
            return false;
        }
        const messages = evt.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            const m = messages[i];
            const toClientMessageCode = m.charAt(0);
            console.log(toClientMessageCode);
            const data = networking.decode(m.substring(1,m.length));
            console.log(data);
            handlers[toClientMessageCode](data);
        }
    };
}