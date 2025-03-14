/* mirrory */
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
// Runs a user command on the server.
function runCommand(command) {
    return __awaiter(this, void 0, void 0, function* () {
        // don't check command content here
        // except for input sanitization
        // and that it's a legit valid command maybe in some ways
        // just construct a POST /cmd request with the command data
        // construct request
        // send request async
        // await get response
        // return response
        var flags = "abcd";
        if (command.indexOf(" ") > -1) {
            var firstSpace = command.indexOf(" ");
            flags = command.substring(firstSpace + 1, command.length);
            command = command.substring(0, firstSpace);
        }
        const response = yield fetch("http://192.168.56.101:8080/commands", {
            method: 'POST',
            body: '{"command": "' + command + '", "flags": "' + flags + '"}',
            headers: { 'Content-Type': 'application/json; charset=UTF-8' }
        });
        if (!response.ok) { /* fail */ }
        let asJSON = { "response": "??", "imagefile": "1.png" };
        if (response.body !== null) {
            asJSON = yield response.json();
        }
        let viewer = document.querySelector('#viewerimg');
        viewer.src = "/img/" + asJSON.imagefile;
        return asJSON.response;
    });
}
// Shows the welcome message and enables the command line
function welcome() {
    return __awaiter(this, void 0, void 0, function* () {
        let output = document.querySelector('#output');
        let cmd = document.querySelector('#cmd');
        // output.innerHTML += "Welcome to Project Indigo. There is a void here." + "\n";
        output.innerHTML += "Indigo is starting... please wait..." + "\n";
        output.innerHTML += (yield runCommand("welcome")) + "\n";
        output.scrollTop = output.scrollHeight;
        cmd.value = "";
        cmd.addEventListener('keydown', (e) => __awaiter(this, void 0, void 0, function* () {
            let output = document.querySelector('#output');
            if (e.key === "Enter") {
                output.innerHTML += (yield runCommand(cmd.value)) + "\n";
                output.scrollTop = output.scrollHeight;
                cmd.value = "";
            }
        }));
    });
}
// sends you to the help zone
function sendhelp() {
    window.location.href = "#helpzone";
}
// save your current session and download the file
function savefile() {
    // send a request to the server to write the current session to a file
    // may be a long running process that requires a loading bar or spinner
    // when done, file will download to user
}
// save your current session as a copy
function savefilecopy() {
    // exact same as the above except the save file has a new guid
    // the file would be safe to send to a friend if you want to send out a copy of your world that they can't modify on the server
    // since there are no accounts, server can't tell who is making modifications, only that the session matches a particular guid
}
// load an existing save file
function loadfile() {
    // pop up a file upload dialog
    // when file is received, send to server and populate db with it
}
// shortcut function to run a move command via a button
function move(direction) {
    return __awaiter(this, void 0, void 0, function* () {
        let output = document.querySelector('#output');
        output.innerHTML += (yield runCommand("move " + direction)) + "\n";
        output.scrollTop = output.scrollHeight;
    });
}
