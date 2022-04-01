import * as networking from '../networking.js';
import { AVATARS } from '../avatars/avatars.js';
import { COLORS, name } from '../landing.js';
import { appendDataLog, setChatboxNotification } from '../ingame-utility.js';
import { initTitles, initHowToPlays } from '../importantStrings.js';
import { playAudio } from '../audio.js';
import { en } from './enum.js';
import { Noneable } from '../noneable.js';

const ingameLeft = document.getElementById("ingame-left");
const endgame = document.getElementById("end-game-button");
const players = document.getElementById("players");
const chatLog = document.getElementById("chat-log");
const scoreToWin = document.getElementById("score-to-win");
const deckSelection = document.getElementById("deck-selection");
const startGame = document.getElementById("start-game-button");
const promptText = document.getElementById("prompt-text");
const promptSubmit = document.getElementById("prompt-submit");
const promptForm = document.getElementById("prompt-form");
const promptField = document.getElementById("prompt-field");
const promptWaiting = document.getElementById("prompt-waiting");
const choices = document.getElementById("choices");
const choicesWaiting = document.getElementById("choices-waiting");
const gameoverText = document.getElementById("gameover-text");
const gameoverContinue = document.getElementById("gameover-continue");

const newgame = new Noneable(document.getElementById("new-game"));
const prompt = new Noneable(document.getElementById("prompt"));
const gameover = new Noneable(document.getElementById("gameover"));

var handlers = {};

export function initMain(conn) {
    initTitles("Fakeout");
    initHowToPlays("Rules\nTry to fool others by filling in the fact.\nThen try to pick the correct fill-in yourself.\n\nScoring\n50 points for faking out someone else\n100 points for guessing the correct answer.");

    scoreToWin.onchange = function (e) {
        networking.send(conn, en.ToServerCode.SCORE_TO_WIN + scoreToWin.value.toString());
    }

    endgame.onclick = function (e) {
        networking.send(conn, en.ToServerCode.END_GAME);
    }

    startGame.onclick = function (e) {
        networking.send(conn, en.ToServerCode.START_GAME);
    }

    gameoverContinue.onclick = function (e) {
        gameover.hide();
        newgame.show();
    }

    promptForm.onsubmit = function (e) {
        e.preventDefault();
        if (!conn) {
            return false;
        }
        if (!promptField.value.trim()) {
            return false;
        }
        networking.send(conn, en.ToServerCode.RESPONSE + promptField.value);
        promptField.value = "";
        return false;
    };

    handlers[en.ToClientCode.IN_MEDIA_RES] = (data) => {
        switch (data[0].charAt(0)) {
        case en.Phase.PREGAME:
            newgame.show();
            break;
        default:
            prompt.show();
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
            var text = "<b>" + data[j] + "</b>" + ": " + data[j+3].toString() + " points<br/>Fakeouts: " + data[j+4];
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

    //PREGAME
    handlers[en.ToClientCode.DECK_OPTIONS] = (data) => {
        for (let i=0; i<data.length; i++) {
            var item = document.createElement("button");
            item.id = "deck_selection_button" + i;
            item.classList.add("deck-selection-button");
            item.innerText = data[i];
            item.onclick = () => {
                networking.send(conn, en.ToServerCode.DECK_SELECTION + i.toString());
            };
            deckSelection.appendChild(item);
        }
    }

    handlers[en.ToClientCode.SCORE_TO_WIN] = (data) => {
        playAudio("click2");
        scoreToWin.value = data[0];
    }
    
    handlers[en.ToClientCode.DECK_SELECTION] = (data) => {
        let i = 0;
        let selectedButton = parseInt(data[0][0]);
        deckSelection.childNodes.forEach(button => {
            if (i === selectedButton) {
                button.classList.add("selected-deck-button");
            } else {
                button.classList.remove("selected-deck-button");
            }
            i++;
        });
    }
    
    handlers[en.ToClientCode.START_GAME] = (data) => {
        playAudio("start");
        newgame.hide();
        gameover.hide();
        prompt.show();
    }

    //PLAY
    handlers[en.ToClientCode.PROMPT] = (data) => {
        while (promptText.firstChild) {
            promptText.removeChild(promptText.firstChild);
        }
        promptText.appendChild(networking.decodeToDiv(data[0]));
    }
    
    handlers[en.ToClientCode.CHOICE_RESPONSE] = (data) => {
        var err = parseInt(data[0]);
        switch (err) {
        case -1:
            promptWaiting.innerText = "Your answer is too close to the actual answer. Pick something else"
            break;
        default:
            playAudio("bubble");
            promptWaiting.innerText = "Response submitted. Waiting for other players."
            promptField.disabled = true;
            promptSubmit.disabled = true;
        }
    }

    handlers[en.ToClientCode.CHOICES] = (data) => {
        playAudio("whoosh");
        promptField.disabled = true;
        promptSubmit.disabled = true;
        promptWaiting.innerText = "";
        for (let j = 0; j < data.length; j++) {
            let d = data[j];
            var item = document.createElement("button");
            item.innerText = d;
            item.onclick = function() {
                playAudio("click3");
                networking.send(conn, en.ToServerCode.CHOICE + j.toString());
            };
            choices.append(item);
        }
    }

    handlers[en.ToClientCode.CHOICES_RESPONSE] = (data) => {
        var err = parseInt(data[0]);
        switch (err) {
        case -1:
            choicesWaiting.innerText = "You can't pick your own answer."
            break;
        default:
            playAudio("dink2");
            choicesWaiting.innerText = "Choice submitted. Waiting for other players."
            for (let j=0; j<choices.childElementCount; j++) {
                choices.children[j].disabled = true;
            }
        }
    }

    handlers[en.ToClientCode.RESULTS] = (data) => {
        playAudio("blupblup");
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
            playAudio("click3");
            promptField.disabled = false;
            promptSubmit.disabled = false;
            choices.classList.remove("revealed");
            while (choices.firstChild) {
                choices.removeChild(choices.firstChild);
            }
            while (choicesWaiting.firstChild) {
                choicesWaiting.removeChild(choicesWaiting.firstChild);
            }
            networking.send(conn, en.ToServerCode.PROMPT_REQUEST);
        };
        choicesWaiting.appendChild(item);
    }
    
    handlers[en.ToClientCode.WINNERS] = (data) => {
        prompt.hide();
        playAudio("fanfare");
        gameoverText.innerText = "Winner: " + data[0] + " " + data[1] + " points\nMost fakeouts: " + data[2] + " " + data[3] + " fakeouts";
        gameover.show();
        while (choices.firstChild) {
            choices.removeChild(choices.firstChild);
        }
        while (choicesWaiting.firstChild) {
            choicesWaitingremoveChild(choicesWaiting.firstChild);
        }
    }
    
    handlers[en.ToClientCode.END_GAME] = (data) => {
        var item = networking.decodeToDiv(data[0]);
        appendDataLog(chatLog, item);
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