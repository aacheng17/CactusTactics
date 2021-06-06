window.onload = function () {
    
    //collapsible start
    var coll = document.getElementsByClassName("collapsible-button");
    for (i = 0; i < coll.length; i++) {
        coll[i].addEventListener("click", function() {
            this.classList.toggle("collapsible-button-active");
            var content = this.nextElementSibling;
            if (content.style.maxHeight){
            content.style.maxHeight = null;
            } else {
            content.style.maxHeight = content.scrollHeight + "px";
            }
        });
    }
    //collapsible end

    col1 = document.getElementsByClassName("howtoplay-text");
    for (i = 0; i < col1.length; i++) {
        col1[i].innerText = "Rules\nTry to think of a word that starts with the first letter and ends with the second letter before your opponents.\nWords must be at least 3 letters long.\n\nScoring\nThe more rare the letter combination, the more points it's worth (up to 100).\nEach word gets a length bonus multiplier as well."
    }

    var conn, name;
    var landing = document.getElementById("landing");
    var nameForm = document.getElementById("name-form");
    var nameField = document.getElementById("name-field");
    var ingame = document.getElementById("ingame");
    var ingameHowtoplayButton = document.getElementById("ingame-howtoplay-button");
    var ingameHowtoplay = document.getElementById("ingame-howtoplay");
    var endgame = document.getElementById("endgame");
    var players = document.getElementById("players");
    var startLetter = document.getElementById("start-letter");
    var endLetter = document.getElementById("end-letter");
    var skip = document.getElementById("skip")
    var promptExtraText = document.getElementById("prompt-extra-text");
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    function send(s) {
        var toSend = "";
        for (var i = 0; i < s.length; i++) {
            toSend += s.replaceAll("\v", "");
            toSend += "\t" + s.replaceAll("\t", "");
        }
        conn.send(s);
    }

    function decode(s) {
        var ret = s.split("\t");
        ret.shift();
        return ret;
    }

    function decodeToHTML(s) {
        return decodeToHTMLHelper(s, 0)[0];
    }

    function decodeToHTMLHelper(s, i) {
        var ret = document.createElement("a");
        var tag = "";
        var state = 0; //0 normal, 1 in a start tag
        for (; i < s.length; i++) {
            if (state == 0) {
                if (s[i] !== "\v") {
                    ret.innerHTML += s[i];
                } else {
                    state = 1;
                }
            } else {
                if (s[i] === "\v") {
                    state = 0;
                    switch (tag) {
                    case "br/":
                        ret.appendChild(document.createElement("br"));
                        break;
                    case "b":
                        var item = document.createElement("b");
                        var result = decodeToHTMLHelper(s, i+1);
                        item.appendChild(result[0]);
                        ret.appendChild(item);
                        console.log(result[1]);
                        i = result[1];
                        break;
                    case "/":
                        return [ret, i];
                    default:
                        break;
                    }
                } else {
                    tag += s[i];
                }
            }
        }
        return [ret, i];
    }

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    if (nameForm !== null) {
        nameForm.onsubmit = function (e) {
            if (!conn) {
                return false;
            }
            if (!nameField.value.trim()) {
                return false;
            }
            name = nameField.value;
            conn.send("1" + name);
            e.preventDefault();
            landing.parentNode.removeChild(landing);
            ingame.style.visibility = "visible";
        }
    }

    ingameHowtoplayButton.addEventListener("click", function() {
        if (ingameHowtoplay.style.maxHeight){
            ingameHowtoplay.style.maxHeight = null;
        } else {
            ingameHowtoplay.style.maxHeight = ingameHowtoplay.scrollHeight + "px";
        } 
    });

    endgame.onclick = function (e) {
        if (!conn) {
            return false;
        }
        conn.send("3");
        endgame.innerText = endgame.innerText === "end game" ? "new game" : "end game";
    }

    skip.onclick = function (e) {
        if (!conn) {
            return false;
        }
        conn.send("2");
    }
    
    document.getElementById("chat-form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value.trim()) {
            return false;
        }
        conn.send("0" + msg.value);
        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        var websocketprotocol = "ws://";
        if (document.location.protocol==="https:") {
            websocketprotocol = "wss://";
        }
        conn = new WebSocket(websocketprotocol + document.location.host + document.location.pathname);
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var m = messages[i];
                var messageType = m.charAt(0);
                var data = decode(m.substring(1,m.length));
                switch (messageType) {
                case '0':
                    var item = document.createElement("div");
                    item.appendChild(decodeToHTML(data[0]));
                    appendLog(item);
                    break;
                case '1':
                    while (players.firstChild) {
                        players.removeChild(players.firstChild);
                    }
                    for (var j = 0; j < data.length; j+=4) {
                        var player = document.createElement("div");
                        player.id = "player";
                        var text = data[j] + " - " + data[j+1].toString() + " points";
                        if (data[j+2] != "") {
                            text += "\nBest word: " + data[j+2] + " " + data[j+3].toString() + " points";
                        }
                        player.innerText = text;
                        players.appendChild(player);
                    }
                    break
                case '2':
                    var sl = data[0].toUpperCase();
                    var el = data[1].toUpperCase();
                    startLetter.innerText = sl;
                    endLetter.innerText = el;
                    msg.placeholder = "a word that starts with " + sl + " and ends with " + el;
                    promptExtraText.innerText = "Worth " + String(data[2]) + " points. There are " + String(data[3]) + " possible words.";
                    break
                case '3':
                    var item = document.createElement("div");
                    item.innerText = "Winner: " + data[0] + " " + data[1] + " points\nBest word: " + data[2] + " " + data[3] + " " + data[4] + " points";
                    appendLog(item);
                    break;
                case '4':
                    while (log.firstChild) {
                        log.removeChild(log.firstChild);
                    }
                    break
                case '5':
                    var item = document.createElement("div");
                    item.classList.add("score-message");
                    var message = decodeToHTML(data[0]);
                    item.appendChild(message);
                    var word = decodeToHTML(data[1]);
                    item.appendChild(word);
                    var what = document.createElement("button");
                    what.classList.add("what-button");
                    what.innerText = "What?";
                    what.onclick = function (e) {
                        if (!conn) {
                            return false;
                        }
                        conn.send("4"+what.previousElementSibling.innerText);
                    }
                    item.appendChild(what);
                    appendLog(item);
                    break
                }
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }

    ingame.style.visibility = "hidden";
};