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

export function decodeToDiv(s) {
    return decodeToDivHelper(s, 0, "div")[0];
}

export function decodeToDivHelper(s, i, tagName) {
    var ret = document.createElement(tagName);
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
                var tagData = tag.split(" ", -1);
                var tagName = tagData[0];
                var tagId = tagData[1];
                switch (tagName) {
                case "br":
                    ret.appendChild(document.createElement("br"));
                    break;
                case "/":
                    return [ret, i];
                default:
                    if (["p", "b"].includes(tagName)) {
                        var result = decodeToDivHelper(s, i+1, tagName);
                        var item = result[0];
                        if (tagId !== "") item.id = tagId;
                        tagData.slice(2).forEach(className => item.classList.add(className))
                        ret.appendChild(item);
                        i = result[1];
                        break;
                    }
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