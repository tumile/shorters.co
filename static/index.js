document.addEventListener("DOMContentLoaded", function () {
  let user = null;
  const token = localStorage.getItem("token");
  if (token) {
    const payload = parseJwt(token);
    user = payload.sub;
  }
  if (user) {
    document.getElementById("user").innerText = user;
    document.querySelector(".nav a").hidden = true;
  }

  const form = document.querySelector(".shorten-form");
  const urlInput = document.querySelector(".shorten-form input");
  const msg = document.getElementById(".shorten-form p");

  form.onsubmit = function (event) {
    event.preventDefault();
    const url = urlInput.value.trim();
    const validUrl = url.match(
      /^http(s)?:\/\/[\w\-._~:/?#[\]@!$&'()\\*+,;=]+$/
    );
    const shortersUrl = url.match(/^http(s)?:\/\/shorters\.co(\/)?.*/);
    if (!validUrl) {
      msg.innerText = "Please enter a valid URL";
    } else if (shortersUrl) {
      msg.innerText = "It's already a Shorters' URL";
    } else {
      fetch("/", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url }),
      })
        .then((res) => res.json())
        .then((data) => (urlInput.value = `http://localhost:8080/${data.key}`))
        .catch((err) => (msg.innerText = err));
    }
    return false;
  };

  function parseJwt(token) {
    const base64Url = token.split(".")[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join("")
    );
    return JSON.parse(jsonPayload);
  }
});
