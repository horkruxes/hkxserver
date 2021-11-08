function copyToClipboard(element) {
    /* Get the text field */
    var copyText = document.getElementById(element);

    /* Copy the text inside the text field */
    navigator.clipboard.writeText(copyText.innerHTML.trim());
}

