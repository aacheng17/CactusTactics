export const MESSAGESEP = "\n"
export const DELIM = "\t"
export const TAG = "\v"

export function send(conn, s) {
    var toSend = "";
    for (var i = 0; i < s.length; i++) {
        toSend += s.replaceAll(TAG, "");
        toSend += DELIM + s.replaceAll(DELIM, "");
    }
    conn.send(s);
}

export function decode(s) {
    var ret = s.split(DELIM);
    ret.shift();
    return ret;
}

export function decodeToHTML(s) {
    return decodeToHTMLHelper(s, 0)[0];
}

export function decodeToHTMLHelper(s, i) {
    var ret = document.createElement("a");
    var tag = "";
    var state = 0; //0 normal, 1 in a tag
    for (; i < s.length; i++) {
        if (state == 0) {
            if (s[i] !== TAG) {
                ret.innerHTML += s[i];
            } else {
                state = 1;
            }
        } else {
            if (s[i] === TAG) {
                switch (tag) {
                case "br/":
                    ret.appendChild(document.createElement("br"));
                    break;
                case "b":
                    var item = document.createElement("b");
                    var result = decodeToHTMLHelper(s, i+1);
                    item.appendChild(result[0]);
                    ret.appendChild(item);
                    i = result[1];
                    break;
                case "/":
                    return [ret, i];
                default:
                    break;
                }
                tag = "";
                state = 0;
            } else {
                tag += s[i];
            }
        }
    }
    return [ret, i];
}