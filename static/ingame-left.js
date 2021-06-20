import * as networking from './networking.js';
import { appendDataLog } from './ingame-utility.js';

var ingameLeft = document.getElementById("ingame-left");
var ingameLeftClickableRegion = document.getElementById("ingame-left-clickable-region");
var leftExpandButton = document.getElementById("ingame-left-expand-button");
var ingameHowtoplayButton = document.getElementById("ingame-howtoplay-button");
var ingameHowtoplay = document.getElementById("ingame-howtoplay");
var endgame = document.getElementById("endgame");
var chatLog = document.getElementById("chat-log");
var chatForm = document.getElementById("chat-form");
var chatField = document.getElementById("chat-field");
var chatboxNotification = document.getElementById("chatbox-notification");
var newMessages = 0;

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

    ingameHowtoplayButton.addEventListener("click", function() {
        var effected = ingameHowtoplay;
        if (effected.style.maxHeight) {
            effected.style.maxHeight = null;
        } else {
            effected.style.maxHeight = effected.scrollHeight + "px";
        } 
    });
    
    endgame.onclick = function (e) {
        if (!conn) {
            return false;
        }
        networking.send(conn, "2");
    }
    
    chatForm.onsubmit = function (e) {
        if (!conn) {
            return false;
        }
        if (!chatField.value.trim()) {
            return false;
        }
        networking.send(conn, "1" + chatField.value);
        e.preventDefault();
        chatField.value = "";
    };
    
    var emptyLog = document.createElement("div");
    emptyLog.classList.add("emptyLog");
    emptyLog.innerText = "Nobody's said anything yet.";
    appendDataLog(chatLog, emptyLog);
}

function expandLeft() {
    leftExpandButton.firstChild.innerText = "Collapse";
    ingameLeft.classList.add("ingame-left-expanded");
    newMessages = 0;
    renderChatboxNotification();
}

function collapseLeft() {
    leftExpandButton.firstChild.innerText = "Expand";
    ingameLeft.classList.remove("ingame-left-expanded");
    ingameLeft.style.boxShadow = null;
}

function renderChatboxNotification() {
    var s = "lobby chat";
    if (newMessages > 0) {
        s += " (" + newMessages.toString() + ")";
    }
    chatboxNotification.innerText = s;
}