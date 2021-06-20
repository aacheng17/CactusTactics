export function appendDataLog(log, item) {
    if (log.firstChild != null) {
        if (log.firstChild.classList.contains("emptyLog")) {
            log.removeChild(log.firstChild);
        }
    }
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}

export function setChatboxNotification(n) {
    var chatboxNotification = document.getElementById("chatbox-notification");
    var s = "lobby chat";
    var t = chatboxNotification.innerText;
    if (n === 0) {
        chatboxNotification.innerText = s;
    } else if (n === 1) {
        if (t === s) {
            chatboxNotification.innerText = s + " (1)";
        } else {
            var c = parseInt(t.charAt(t.length-2));
            if (Number.isNaN(c)) return false;
            chatboxNotification.innerText = s + " (" + (c+1).toString() + ")";
        }
    }
}