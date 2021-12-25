if ("serviceWorker" in navigator && window.location.href !== "http://127.0.0.1:3000/") {
  navigator.serviceWorker
    .register("/sw.js", { scope: "/" })
    .then((reg) => {
      console.log("Service Worker Registered", reg);
    })
    .catch((err) => {
      console.error("Service Worker Registeration Error", err);
    });

  navigator.serviceWorker.ready
    .then((reg) => {
      console.log("Service Worker Ready", reg);
    })
    .catch((err) => console.error("Service Worker Not Ready: Error", err));
}
