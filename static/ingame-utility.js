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