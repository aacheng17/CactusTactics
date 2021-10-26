import { initCollapsible } from './collapsible.js';
import { initLanding } from './landing.js';
import { initIngameLeft } from './ingame-left.js';
import { appendDataLog } from './ingame-utility.js';

window.onload = async function () {
    var gameName = window.location.pathname.slice(1);
    var firstSlashIndex = gameName.indexOf("/");
    if (firstSlashIndex !== -1) { gameName = gameName.slice(0, firstSlashIndex); }
    let { initMain } = await import('./' + gameName + '/game.js');
    
    var head = document.getElementsByTagName('head')[0];
    var style = document.createElement('link');
    style.href = './static/' + gameName + '/game.css';
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

    initCollapsible();
    initLanding(conn);
    initIngameLeft(conn);
    initMain(conn);
    ingame.parentNode.removeChild(ingame);
};