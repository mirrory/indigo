/* mirrory */
// Runs a user command on the server.
function runCommand(command) {
    // don't check command content here
    // except for input sanitization
    // and that it's a legit valid command maybe in some ways
    // just construct a POST /cmd request with the command data
    // construct request
    // send request async
    // await get response
    // return response
    return "I understand...";
}
// Shows the welcome message and enables the command line
function welcome() {
    var output = document.querySelector('#output');
    output.innerHTML += "Welcome to Project Indigo. There is a void here." + "\n";
    var cmd = document.querySelector('#cmd');
    cmd.addEventListener('keydown', function (e) {
        var output = document.querySelector('#output');
        if (e.key === "Enter") {
            output.innerHTML += runCommand(cmd.value) + "\n";
            output.scrollTop = output.scrollHeight;
            cmd.value = "";
        }
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
    var output = document.querySelector('#output');
    output.innerHTML += runCommand("move " + direction) + "\n";
    output.scrollTop = output.scrollHeight;
}
