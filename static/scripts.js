function copyToClipboard(element) {
    /* Get the text field */
    var copyText = document.getElementById(element);

    /* Copy the text inside the text field */
    navigator.clipboard.writeText(copyText.innerHTML.trim());
}

function sendToClipboard(text) {
    navigator.clipboard.writeText(text);
}

document.getElementById('message').onkeyup = function () {
    var charCount = Math.max(0, 140 - document.getElementById('message').value.length);
    if (charCount > 0) {
        document.getElementById('countChars').innerHTML = charCount + " characters left before being able to send your message";
    } else {
        document.getElementById('countChars').innerHTML = "";
    }
};