<html>
    <body>
        <input id="input" type="text" />
        <button onclick="send()">Send</button>
        <pre id="output"></pre>
        <script>
            var input = document.getElementById("input");
            var output = document.getElementById("output");
            var socket = new WebSocket("ws://127.0.0.1:8080/main/my2");

            socket.onopen = function () {
                output.innerHTML += "Status: Connected\n";
            };

            socket.onmessage = function (e) {
                output.innerHTML += "Server: " + e.data + "\n";
            };

            function send() {
                socket.send(input.value);
                input.value = "";
            }
        </script>
    </body>
</html>