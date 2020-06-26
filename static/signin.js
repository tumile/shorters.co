document.addEventListener("DOMContentLoaded", function () {
  const title = document.querySelector(".sign-in h2");
  const form = document.querySelector(".sign-in-form");
  const [emailInput, otpInput] = document.querySelectorAll(
    ".sign-in-form input"
  );
  const btn = document.querySelector(".sign-in-form button");
  const msg = document.querySelector(".sign-in-form p");

  form.onsubmit = onSubmitEmail;

  function onSubmitEmail(event) {
    event.preventDefault();
    const email = emailInput.value.trim();
    if (email) {
      fetch("/signin", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      })
        .then(() => {
          emailInput.disabled = true;
          otpInput.hidden = false;
          title.innerText = "Check your email for a one-time password!";
          btn.innerText = "Sign in";
          form.onsubmit = onSubmitOTP;
        })
        .catch((err) => {
          console.error(err);
          msg.innerText = err;
        });
    }
    return false;
  }

  function onSubmitOTP(event) {
    event.preventDefault();
    const email = emailInput.value.trim();
    const otp = otpInput.value.trim();
    if (email && otp) {
      fetch("/verify", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, otp }),
      })
        .then((res) => {
          if (!res.ok) throw res;
          return res.json();
        })
        .then((data) => {
          localStorage.setItem("token", data.token);
          window.location.replace("http://localhost:8080");
        })
        .catch((err) => {
          console.error(err);
          msg.innerText = "Hmm something was wrong...";
        });
    }
    return false;
  }
});
