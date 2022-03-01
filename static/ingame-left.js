import * as networking from './networking.js';
import { appendDataLog, setChatboxNotification } from './ingame-utility.js';

var ingameLeft = document.getElementById("ingame-left");
var ingameLeftClickableRegion = document.getElementById("ingame-left-clickable-region");
var leftExpandButton = document.getElementById("ingame-left-expand-button");
var moreGames = document.getElementById("more-games");
var ingameHowtoplayButton = document.getElementById("ingame-howtoplay-button");
var ingameHowtoplay = document.getElementById("ingame-howtoplay");
var chatLog = document.getElementById("chat-log");
var chatForm = document.getElementById("chat-form");
var chatField = document.getElementById("chat-field");
import { en } from './enum.js';

export function initIngameLeft(conn) {
    window.addEventListener('click', function(e){
        var width = (window.innerWidth > 0) ? window.innerWidth : screen.width;
        if (width <= 800) {
            if (ingameLeft.classList.contains("ingame-left-expanded")) {
                if (!ingameLeft.contains(e.target)) {
                    collapseLeft();
                }
            }
        }
    });

    ingameLeftClickableRegion.addEventListener("click", function() {
        if (!ingameLeft.classList.contains("ingame-left-expanded")) {
            expandLeft();
        } else {
            collapseLeft();
        }
    });

    moreGames.href = "https://cactustactics.herokuapp.com/";

    ingameHowtoplayButton.addEventListener("click", function() {
        var effected = ingameHowtoplay;
        if (effected.style.maxHeight) {
            effected.style.maxHeight = null;
        } else {
            effected.style.maxHeight = effected.scrollHeight + "px";
        } 
    });
    
    chatForm.onsubmit = function (e) {
        if (!conn) {
            return false;
        }
        if (!chatField.value.trim()) {
            return false;
        }
        networking.send(conn, en.ToServerCode.LOBBY_CHAT_MESSAGE + chatField.value);
        e.preventDefault();
        chatField.value = "";
    };
    
    var emptyLog = document.createElement("div");
    emptyLog.classList.add("emptyLog");
    appendDataLog(chatLog, emptyLog);
}

function expandLeft() {
    leftExpandButton.firstChild.innerText = "Collapse";
    ingameLeft.classList.add("ingame-left-expanded");
    setChatboxNotification(0);
}

function collapseLeft() {
    leftExpandButton.firstChild.innerText = "Expand";
    ingameLeft.classList.remove("ingame-left-expanded");
    ingameLeft.style.boxShadow = null;
}