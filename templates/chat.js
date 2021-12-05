
window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }


    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };

        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                if (messages[i].trim().length === 0) {
                    break;
                }
                var msgObj = JSON.parse(messages[i]);
                console.log(msgObj);
                var item = document.createElement("div");
                item.setAttribute('data-from', msgObj.User.DisplayName);
                item.setAttribute('data-id', msgObj.ID);

                var color = msgObj.User.Color;
                if (!color) {
                    color = 'blue';
                }
                var subSpan = document.createElement("span");
                subSpan.setAttribute('className', 'meta');
                subSpan.setAttribute('style', 'color: ' + color)

                var nameSpan = document.createElement('span');
                nameSpan.setAttribute('className', 'name')
                nameSpan.innerHTML = msgObj.Channel + " " + msgObj.User.DisplayName + " ";
                subSpan.appendChild(nameSpan)
                item.appendChild(subSpan);

                subSpan = document.createElement('span');
                subSpan.setAttribute('className', 'message');
                subSpan.innerHTML = msgObj.Message;
                item.appendChild(subSpan);

                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};