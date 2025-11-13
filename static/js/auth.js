document.addEventListener("DOMContentLoaded", () => {
  const loginLink = document.getElementById("loginLink");
  const logoutLink = document.getElementById("logoutLink");
  const accountLink = document.getElementById("accountLink");

  // Check if user is logged in via backend
  checkAuthStatus();

  async function checkAuthStatus() {
    try {
      const response = await fetch('/api/user');
      if (response.ok) {
        const data = await response.json();
        if (data.ok && data.data) {
          // User is logged in
          if (loginLink) loginLink.classList.add("d-none");
          if (logoutLink) logoutLink.classList.remove("d-none");
          if (accountLink) accountLink.classList.remove("d-none");
          return;
        }
      }
    } catch (error) {
      console.log('Not authenticated');
    }
    
    // User is not logged in
    if (loginLink) loginLink.classList.remove("d-none");
    if (logoutLink) logoutLink.classList.add("d-none");
    if (accountLink) accountLink.classList.add("d-none");
  }

  // --- Logout ---
  if (logoutLink) {
    logoutLink.addEventListener("click", async (e) => {
      e.preventDefault();
      const result = await Swal.fire({
        title: "Log out?",
        icon: "warning",
        showCancelButton: true,
        confirmButtonText: "Yes, log out",
        cancelButtonText: "Cancel"
      });
      
      if (result.isConfirmed) {
        try {
          const response = await fetch('/api/logout', { method: 'POST' });
          const data = await response.json();
          
          if (data.ok) {
            await Swal.fire({
              title: "Logged out",
              icon: "success",
              timer: 1500,
              showConfirmButton: false
            });
            window.location.href = "/";
          } else {
            throw new Error(data.message);
          }
        } catch (error) {
          Swal.fire("Error", error.message || "Failed to logout", "error");
        }
      }
    });
  }

  // --- Login / Register page ---
  const loginForm = document.getElementById("loginForm");
  const registerBtn = document.getElementById("registerBtn");

  if (loginForm) {
    // Login form submit
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const email = document.getElementById("email").value.trim();
      const password = document.getElementById("password").value.trim();

      try {
        const response = await fetch('/api/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ email, password }),
        });

        const data = await response.json();

        if (data.ok) {
          await Swal.fire({
            title: "Welcome back!",
            icon: "success",
            timer: 1500,
            showConfirmButton: false
          });
          window.location.href = "/store";
        } else {
          Swal.fire("Invalid credentials!", data.message || "Please check your email and password.", "error");
        }
      } catch (error) {
        Swal.fire("Error", "Failed to login. Please try again.", "error");
      }
    });

    // Register button click
    if (registerBtn) {
      registerBtn.addEventListener("click", async () => {
        const email = document.getElementById("email").value.trim();
        const password = document.getElementById("password").value.trim();

        if (!email || !password) {
          Swal.fire("Enter all fields!", "", "warning");
          return;
        }

        if (password.length < 6) {
          Swal.fire("Password too short!", "Password must be at least 6 characters long.", "warning");
          return;
        }

        try {
          const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
          });

          const data = await response.json();

          if (data.ok) {
            await Swal.fire({
              title: "Registered!",
              text: "You can now log in with your credentials.",
              icon: "success"
            });
            // Clear password field for security
            document.getElementById("password").value = "";
          } else {
            Swal.fire("Registration failed!", data.message || "Please try again.", "error");
          }
        } catch (error) {
          Swal.fire("Error", "Failed to register. Please try again.", "error");
        }
      });
    }
  }

  // --- Account page logic ---
  const accountInfo = document.getElementById("accountInfo");
  if (accountInfo) {
    fetch('/api/user')
      .then(response => response.json())
      .then(data => {
        if (!data.ok || !data.data) {
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

        const user = data.data;
        accountInfo.innerHTML = `
          <p><strong>Email:</strong> ${user.email}</p>
          <p><strong>Registered on:</strong> ${new Date(user.created_at).toLocaleDateString()}</p>
        `;

        // Load orders
        loadOrders();
      })
      .catch(() => {
        window.location.href = "/login";
      });
  }

  async function loadOrders() {
    try {
      const response = await fetch('/api/orders');
      const data = await response.json();
      
      if (data.ok && data.data) {
        const ordersContainer = document.getElementById("ordersContainer");
        if (ordersContainer) {
          if (data.data.length === 0) {
            ordersContainer.innerHTML = '<p class="text-muted">No orders yet.</p>';
          } else {
            ordersContainer.innerHTML = data.data.map(order => `
              <div class="card mb-3">
                <div class="card-body">
                  <h5>Order #${order.id}</h5>
                  <p><strong>Date:</strong> ${new Date(order.created_at).toLocaleDateString()}</p>
                  <p><strong>Total:</strong> €${order.total_price.toFixed(2)}</p>
                  <p><strong>Status:</strong> <span class="badge bg-${order.status === 'pending' ? 'warning' : 'success'}">${order.status}</span></p>
                </div>
              </div>
            `).join('');
          }
        }
      }
    } catch (error) {
      console.error('Error loading orders:', error);
    }
  }
});
