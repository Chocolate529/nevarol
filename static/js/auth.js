document.addEventListener("DOMContentLoaded", () => {
  const loginLink = document.getElementById("loginLink");
  const logoutLink = document.getElementById("logoutLink");
  const accountLink = document.getElementById("accountLink");

  const user = JSON.parse(localStorage.getItem("user"));

  // --- Navbar visibility setup ---
  if (user) {
    if (loginLink) loginLink.classList.add("d-none");
    if (logoutLink) logoutLink.classList.remove("d-none");
    if (accountLink) accountLink.classList.remove("d-none");
  } else {
    if (loginLink) loginLink.classList.remove("d-none");
    if (logoutLink) logoutLink.classList.add("d-none");
    if (accountLink) accountLink.classList.add("d-none");
  }

  // --- Logout ---
  if (logoutLink) {
    logoutLink.addEventListener("click", (e) => {
      e.preventDefault();
      Swal.fire({
        title: "Log out?",
        icon: "warning",
        showCancelButton: true,
        confirmButtonText: "Yes, log out",
        cancelButtonText: "Cancel"
      }).then(result => {
        if (result.isConfirmed) {
          localStorage.removeItem("user");
          Swal.fire({
            title: "Logged out",
            icon: "success",
            timer: 1500,
            showConfirmButton: false
          }).then(() => {
            window.location.href = "/";
          });
        }
      });
    });
  }

  // --- Login / Register page ---
  const loginForm = document.getElementById("loginForm");
  const registerBtn = document.getElementById("registerBtn");

  if (loginForm) {
    // Login form submit
    loginForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const email = document.getElementById("email").value.trim();
      const password = document.getElementById("password").value.trim();

      const storedUser = JSON.parse(localStorage.getItem("user"));
      if (storedUser && storedUser.email === email && storedUser.password === password) {
        Swal.fire({
          title: "Welcome back!",
          icon: "success",
          timer: 1500,
          showConfirmButton: false
        }).then(() => {
          window.location.href = "/store";
        });
      } else {
        Swal.fire("Invalid credentials!", "Please check your email and password.", "error");
      }
    });

    // Register button click
    if (registerBtn) {
      registerBtn.addEventListener("click", () => {
        const email = document.getElementById("email").value.trim();
        const password = document.getElementById("password").value.trim();

        if (!email || !password) {
          Swal.fire("Enter all fields!", "", "warning");
          return;
        }

        const user = { email, password, date: new Date().toLocaleDateString() };
        localStorage.setItem("user", JSON.stringify(user));

        Swal.fire({
          title: "Registered!",
          text: "You can now log in with your credentials.",
          icon: "success"
        });
      });
    }
  }

  // --- Account page logic ---
  const accountInfo = document.getElementById("accountInfo");
  if (accountInfo) {
    const currentUser = JSON.parse(localStorage.getItem("user"));

    if (!currentUser) {
      Swal.fire({
        title: "Not logged in!",
        text: "Please log in to access your account.",
        icon: "warning",
        confirmButtonText: "Go to Login"
      }).then(() => {
        window.location.href = "/login";
      });
      return;
    }

    accountInfo.innerHTML = `
      <p><strong>Email:</strong> ${currentUser.email}</p>
      <p><strong>Registered on:</strong> ${currentUser.date || "N/A"}</p>
    `;
  }
});
