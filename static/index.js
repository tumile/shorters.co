const urlInput = document.getElementById("url-input");
const list = document.getElementById("url-list");
const error = document.getElementById("error");
const customBtn = document.getElementById("custom-btn");
const customField = document.getElementById("custom-field");
const customInput = document.getElementById("custom-input");
new ClipboardJS(".copy-btn");

let customToggled = false;

customBtn.addEventListener("click", () => {
  customField.style.display = customToggled ? "none" : "flex";
  customToggled = !customToggled;
});

let urls = localStorage.getItem("urls");
urls = urls ? JSON.parse(urls) : [];

if (urls.length > 0) {
  renderList();
} else {
  renderEmptyList();
}

function onSubmit(e) {
  e.preventDefault();
  clearError();
  const url = urlInput.value.trim();
  if (!url || !validateUrl(url)) {
    showInputError();
    return;
  }
  const prev = urls.find((i) => i.url === url);
  if (prev) {
    urls = urls.filter((i) => i.url !== url);
    urls.unshift(prev);
    renderList();
    clearInput();
    return;
  }
  if (customToggled) {
    const customUrl = customInput.value.trim();
    shortenUrl(url, customUrl);
  } else {
    shortenUrl(url);
  }
}

function shortenUrl(url, key) {
  fetch(key ? "/store" : "/shorten", {
    method: "POST",
    headers: {
      "content-type": "application/json",
    },
    body: JSON.stringify({ url, key }),
  })
    .then((res) => {
      if (!res.ok) {
        throw res;
      }
      return res.json();
    })
    .then((res) => {
      if (urls.length === 0) {
        clearList();
      }
      urls.unshift(res);
      list.prepend(renderListItem(res));
      clearInput();
    })
    .catch((err) => {
      showError(
        err.status === 400
          ? "Request malformed"
          : err.status === 409
          ? "Custom URL already exists"
          : "Uh oh... Please try again later"
      );
    });
}

function renderEmptyList() {
  const li = document.createElement("li");
  li.className = "list-item";
  li.innerText = "You don't have any shortened URLs yet!";
  list.append(li);
}

function renderList() {
  clearList();
  urls.map((item) => list.append(renderListItem(item)));
}

function renderListItem(item) {
  const { url, shortUrl } = item;
  const li = document.createElement("li");
  li.className = "list-item";
  li.innerHTML = `<div class="short-url">
  <a href="${shortUrl}" target="_blank" class="nes-text is-primary">${shortUrl}</a>
  <button type="button" class="nes-btn is-success copy-btn" data-clipboard-text="${shortUrl}">Copy</button>
</div>
<div class="original-url">
  <a href="${url}" target="_blank" class="nes-text is-disabled">${url}</a>
</div>`;
  return li;
}

function clearList() {
  list.innerHTML = "";
}

function clearInput() {
  urlInput.value = "";
  customInput.value = "";
}

function clearError() {
  error.style.display = "none";
  urlInput.classList.remove("is-error");
}

function showError(msg) {
  error.innerText = msg;
  error.style.display = "block";
}

function showInputError() {
  urlInput.classList.add("is-error");
  setTimeout(() => urlInput.classList.remove("is-error"), 2000);
}

function validateUrl(value) {
  return /^(?:(?:(?:https?|ftp):)?\/\/)(?:\S+(?::\S*)?@)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)(?:\.(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})))(?::\d{2,5})?(?:[/?#]\S*)?$/i.test(
    value
  );
}

window.addEventListener("beforeunload", () => {
  localStorage.setItem("urls", JSON.stringify(urls));
});
