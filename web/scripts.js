import { SignMessage } from './lib.js'

export function copyToClipboard(element) {
  /* Get the text field */
  var copyText = document.getElementById(element);

  /* Copy the text inside the text field */
  navigator.clipboard.writeText(copyText.innerHTML.trim());
}
document.getElementById('message').onkeyup = function () {
  var charCount = Math.max(0, 140 - document.getElementById('message').value.length);
  if (charCount > 0) {
    document.getElementById('countChars').innerHTML = charCount + " characters left before being able to send your message";
  } else {
    document.getElementById('countChars').innerHTML = "";
  }
};

console.log("hello world")

document.getElementById('name').onkeyup = function () {
  console.log("clicking on name")
  const sec = document.getElementById('public-key').value
  const pub = document.getElementById('secret-key').value
  const name = document.getElementById('name').value
  const msg = document.getElementById('message').value
  const msgID = document.getElementById('answer-to').value

  const signature = SignMessage(sec, pub, name, msg, msgID,)
  console.log("signature from", name, signature)
  document.getElementById('signature').value = signature

}
