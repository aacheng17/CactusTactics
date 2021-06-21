import { initCollapsible } from './collapsible.js';
import { initLanding } from './landing.js';
import { initIngameLeft } from './ingame-left.js';
import { appendDataLog } from './ingame-utility.js';
import { initTitles, initHowToPlays } from './importantStrings.js';

window.onload = async function () {
    var gameName = window.location.pathname.slice(1);
    let { initMain } = await import('./' + gameName + '/game.js');
    
    var head = document.getElementsByTagName('head')[0];
    var style = document.createElement('link');
    style.href = './static/' + gameName + '/idiotmouth.css';
    style.type = 'text/css';
    style.rel = 'stylesheet';
    head.append(style);
    
    var conn;
    var ingame = document.getElementById("ingame");
    var gameLog = document.getElementById("game-log");

    if (window["WebSocket"]) {
        var websocketprotocol = "ws://";
        if (document.location.protocol==="https:") {
            websocketprotocol = "wss://";
        }
        conn = new WebSocket(websocketprotocol + document.location.host + document.location.pathname);
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendDataLog(gameLog, item);
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendDataLog(gameLog, item);
    }

    initTitles("Idiotmouth");
    initHowToPlays("Rules\nTry to think of a word that starts with the first letter and ends with the second letter before your opponents.\nWords must be at least 3 letters long.\n\nScoring\nThe more rare the letter combination, the more points it's worth (up to 100).\nEach word gets a length bonus multiplier as well.");
    initCollapsible();
    initLanding(conn);
    initIngameLeft(conn);
    initMain(conn);
    ingame.parentNode.removeChild(ingame);
};