import * as networking from '../networking.js';
import { AVATARS } from '../avatars/avatars.js';
import { COLORS, name } from '../landing.js';
import { appendDataLog, setChatboxNotification } from '../ingame-utility.js';
import { initTitles, initHowToPlays } from '../importantStrings.js';
import { playAudio } from '../audio.js';
import { en } from './enum.js';
import { Noneable } from '../noneable.js';
import { GameOption } from '../gameoption.js';

const ingameLeft = document.getElementById("ingame-left");
const endGameButton = document.getElementById("end-game-button");
const players = document.getElementById("players");
const chatLog = document.getElementById("chat-log");
const minWordLengthDiv = document.getElementById("min-word-length-div");
const scoreToWinDiv = document.getElementById("score-to-win-div");
const chaosMode = document.getElementById("chaos-mode");
const startGameButton = document.getElementById("start-game-button");
const winners = document.getElementById("winners");
const newGameButton = document.getElementById("new-game-button");
const topboxWord = document.getElementById("topbox-word");
const topboxBackspaceButton = document.getElementById("topbox-backspace-button");
const topboxClearButton = document.getElementById("topbox-clear-button");
const topboxSubmitButton = document.getElementById("topbox-submit-button");
const topboxAddLetterButton = document.getElementById("topbox-add-letter-button");

const newgame = new Noneable(document.getElementById("new-game"));
const topboxInfo = new Noneable(document.getElementById("topbox-info"));
const endgame = new Noneable(document.getElementById("end-game"));
const gamebox = new Noneable(document.getElementById("gamebox"));

var handlers = {};

const selectedLetterObjects = [];

function toggleGameboxLetter(gameboxLetterElement, n) {
    if (gameboxLetterElement.innerText === '') {
        return;
    }
    if (gameboxLetterElement.classList.contains("selected-letter")) {
        gameboxLetterElement.classList.remove("selected-letter")

        const letterObjectIndex = selectedLetterObjects.findIndex(o => o[0] === gameboxLetterElement);
        const letterObject = selectedLetterObjects[letterObjectIndex];
        letterObject[1].parentNode.removeChild(letterObject[1]);

        selectedLetterObjects.splice(letterObjectIndex, 1);
    } else {
        gameboxLetterElement.classList.add("selected-letter")

        const newTopboxElement = document.createElement("h3");
        newTopboxElement.innerText = gameboxLetterElement.innerText;
        newTopboxElement.className = "topbox-letter";
        newTopboxElement.onclick = () => toggleGameboxLetter(gameboxLetterElement);
        topboxWord.appendChild(newTopboxElement);

        selectedLetterObjects.push([gameboxLetterElement, newTopboxElement, n]);
    }
}

const backspace = () => {
    if (selectedLetterObjects.length === 0) {
        return;
    }
    toggleGameboxLetter(selectedLetterObjects[selectedLetterObjects.length - 1][0]);
}

topboxBackspaceButton.onclick = backspace;

topboxClearButton.onclick = () => {
    while (selectedLetterObjects.length > 0) {
        backspace();
    }
}

function initGameboxLetters() {
    for (let i = 0; i < 20; i++) {
        const element = document.createElement("h2");
        element.className = "gamebox-letter";
        element.innerText = '';
        element.onclick = () => toggleGameboxLetter(element, i);
        gamebox.element.appendChild(element);
    }
}

async function setGameboxLetters(newLetters) {
    for (let i = 0; i < newLetters.length; i++) {
        const existingLetterElement = gamebox.element.children.item(i);
        const newLetter = newLetters[i];
        if (existingLetterElement.innerText === newLetter) {
            continue;
        }
        if (newLetter === ' ') {
            if (existingLetterElement.classList.contains("gamebox-letter-visible")) {
                await new Promise(r => setTimeout(r, 10));
                existingLetterElement.classList.remove("gamebox-letter-visible");
            }
        } else {
            if (!existingLetterElement.classList.contains("gamebox-letter-visible")) {
                await new Promise(r => setTimeout(r, 10));
                existingLetterElement.classList.add("gamebox-letter-visible");
            }
        }
        existingLetterElement.innerText = newLetter;
    }
}

export function initMain(conn) {
    initTitles("Aaranagrams");
    initHowToPlays("Rules\nTake turns generating letters. Assemble words to score points.\n\nScoring\nThe more rare the letter combination, the more points it's worth (up to 100).\nEach word gets a length bonus multiplier as well.");
    initGameboxLetters();

    endGameButton.onclick = async function (e) {
        networking.send(conn, en.ToServerCode.END_GAME);
    }

    const minWordLength = new GameOption(minWordLengthDiv, "min-word-length", "Minimum word length:", 8, 1, 1);
    minWordLength.left.onclick = function (e) {
        minWordLength.decrement();
        networking.send(conn, en.ToServerCode.MIN_WORD_LENGTH + minWordLength.getValue().toString());
    }

    minWordLength.right.onclick = function (e) {
        minWordLength.increment();
        networking.send(conn, en.ToServerCode.MIN_WORD_LENGTH + minWordLength.getValue().toString());
    }
    
    const scoreToWin = new GameOption(scoreToWinDiv, "score-to-win", "Score to win:", 50000, 500, 500);
    scoreToWin.left.onclick = function (e) {
        scoreToWin.decrement();
        networking.send(conn, en.ToServerCode.SCORE_TO_WIN + scoreToWin.getValue().toString());
    }

    scoreToWin.right.onclick = function (e) {
        scoreToWin.increment();
        networking.send(conn, en.ToServerCode.SCORE_TO_WIN + scoreToWin.getValue().toString());
    }

    minWordLengthDiv.onchange = function (e) {
        networking.send(conn, en.ToServerCode.MIN_WORD_LENGTH + minWordLengthDiv.value.toString());
    }

    scoreToWin.onchange = function (e) {
        networking.send(conn, en.ToServerCode.SCORE_TO_WIN + scoreToWin.value.toString());
    }

    chaosMode.onclick = function (e) {
        networking.send(conn, en.ToServerCode.CHAOS_MODE + (chaosMode.checked ? "1" : "0"));
    }

    startGameButton.onclick = function (e) {
        networking.send(conn, en.ToServerCode.START_GAME);
    }
    
    topboxSubmitButton.onclick = function () {
        let stringToSend = "";
        selectedLetterObjects.forEach(letterObject => stringToSend += letterObject[2])
        networking.send(conn, en.ToServerCode.GAME_MESSAGE + stringToSend);
    };

    topboxAddLetterButton.onclick = function () {
        networking.send(conn, en.ToServerCode.LETTER);
    }

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
            topboxInfo.show();
            gamebox.show();
            break;
        }
    }
    
    // ALL PHASES
    handlers[en.ToClientCode.PLAYERS] = (data) => {
        while (players.firstChild) {
            players.removeChild(players.firstChild);
        }
        for (var j = 0; j < data.length; j+=7) {
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
            if (data[j+6]) {
                var playerStatus = document.createElement("a")
                playerStatus.classList.add("player-status")
                playerStatus.classList.add("dotdotdot");
                playerAvatarContainer.appendChild(playerStatus);
            }
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
    handlers[en.ToClientCode.MIN_WORD_LENGTH] = (data) => {
        playAudio("click2");
        minWordLength.setValue(data[0]);
    }
    
    handlers[en.ToClientCode.SCORE_TO_WIN] = (data) => {
        playAudio("click2");
        scoreToWin.setValue(data[0]);
    }
    
    handlers[en.ToClientCode.CHAOS_MODE] = (data) => {
        playAudio("click2");
        chaosMode.checked = data[0] === "1";
    }
    
    handlers[en.ToClientCode.START_GAME] = (data) => {
        topboxInfo.show();
        gamebox.show();
        newgame.hide();
        endgame.hide();
        playAudio("start");
    }
    
    //PLAY
    handlers[en.ToClientCode.YOUR_TURN] = (data) => {
        topboxAddLetterButton.disabled = data[0] !== "1";
    }

    handlers[en.ToClientCode.LETTERS] = (data) => {
        playAudio("click3");
        setGameboxLetters(data[0]);
    }

    handlers[en.ToClientCode.GAME_MESSAGE] = (data) => {
        /*playAudio("click3");
        var item = networking.decodeToDiv(data[0]);
        appendDataLog(gameLog, item, true);*/
    }
    
    handlers[en.ToClientCode.PROMPT] = (data) => {
        /*playAudio("correct");
        var sl = data[0].toUpperCase();
        var el = data[1].toUpperCase();
        startLetter.innerText = sl;
        endLetter.innerText = el;
        gameField.placeholder = sl + "___" + el;
        promptExtraText.innerText = "Worth " + String(data[2]) + " points. There are " + String(data[3]) + " possible words.";*/
    }
    
    // TODO: Bring WHAT handlers back in once game is working
    handlers[en.ToClientCode.MESSAGE_WITH_WHAT] = (data) => {
        /*var item = document.createElement("div");
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
        item.appendChild(what);*/
    }
    
    handlers[en.ToClientCode.WHAT_RESPONSE] = (data) => {
        /*playAudio("dink2");
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
        }*/
    }
    
    handlers[en.ToClientCode.END_GAME] = (data) => {
        topboxInfo.hide();
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