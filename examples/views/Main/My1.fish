<html>
    <body>
        <input id="input" type="text" />
        <button onclick="send()">Send</button>
        <pre id="output"></pre>
        <script>
            var input = document.getElementById("input");
            var output = document.getElementById("output");
            var socket = new WebSocket("wss://127.0.0.1:8081/main/my2");

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