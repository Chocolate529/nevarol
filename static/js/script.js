

// Products will be loaded from backend
let products = [];

// Cart will be synced with backend
let cart = [];

let currentPage = 1;
const itemsPerPage = 6;
let currentFilter = "all";
let currentSearch = "";
let priceMin = 10;
let priceMax = 30;

// Load products from backend
async function loadProducts() {
  try {
    const response = await fetch('/api/products');
    const data = await response.json();
    if (data.ok && data.data) {
      products = data.data;
      renderProducts();
    }
  } catch (error) {
    console.error('Error loading products:', error);
  }
}

// Load cart from backend
async function loadCart() {
  try {
    const response = await fetch('/api/cart');
    const data = await response.json();
    if (data.ok && data.data) {
      cart = data.data;
      renderCart();
    } else {
      // User not authenticated, use empty cart
      cart = [];
      renderCart();
    }
  } catch (error) {
    console.error('Error loading cart:', error);
    cart = [];
    renderCart();
  }
}

// Render products with filter/search/pagination
function renderProducts() {
  const container = document.getElementById("product-list");
  const pagination = document.getElementById("pagination");
  if (!container) return;

  container.innerHTML = "";

  let filtered = products.filter(p => {
    const matchesFilter = currentFilter === "all" || p.type === currentFilter;
    const matchesSearch = p.name.toLowerCase().includes(currentSearch.toLowerCase());
    const matchesPrice = p.price >= priceMin && p.price <= priceMax;
    return matchesFilter && matchesSearch && matchesPrice;
  });

  // Pagination
  const totalPages = Math.ceil(filtered.length / itemsPerPage);
  if (currentPage > totalPages && totalPages > 0) currentPage = 1;
  const start = (currentPage - 1) * itemsPerPage;
  const paginated = filtered.slice(start, start + itemsPerPage);

  // Render products
  paginated.forEach(p => {
    const card = document.createElement("div");
    card.className = "col-md-4 mb-4";
    card.innerHTML = `
      <div class="card h-100 shadow-sm">
        <img src="./static/${p.image}" class="card-img-top" alt="${p.name}" onerror="this.src='./static/images/placeholder.jpg'">
        <div class="card-body d-flex flex-column">
          <h5 class="card-title">${p.name}</h5>
          <p class="card-text fw-bold">€${p.price.toFixed(2)}</p>
          <button class="btn btn-primary mt-auto" onclick="addToCart(${p.id})">Add to Cart</button>
        </div>
      </div>
    `;
    container.appendChild(card);
  });

  // Render pagination buttons
  if (pagination) {
    pagination.innerHTML = "";
    for (let i = 1; i <= totalPages; i++) {
      const li = document.createElement("li");
      li.className = `page-item ${i === currentPage ? "active" : ""}`;
      li.innerHTML = `<button class="page-link">${i}</button>`;
      li.addEventListener("click", () => {
        currentPage = i;
        renderProducts();
      });
      pagination.appendChild(li);
    }
  }
}

// --- Filters ---
document.addEventListener("DOMContentLoaded", () => {
  loadProducts();
  loadCart();
  initCartToggle();
  

  const searchInput = document.getElementById("search-input");
  const filterSelect = document.getElementById("filter-select");
  const priceMinInput = document.getElementById("price-min");
  const priceMaxInput = document.getElementById("price-max");
  const priceMinLabel = document.getElementById("price-min-label");
  const priceMaxLabel = document.getElementById("price-max-label");

  if (searchInput) {
    searchInput.addEventListener("input", () => {
      currentSearch = searchInput.value;
      currentPage = 1;
      renderProducts();
    });
  }

  if (filterSelect) {
    filterSelect.addEventListener("change", () => {
      currentFilter = filterSelect.value;
      currentPage = 1;
      renderProducts();
    });
  }

  if (priceMinInput && priceMaxInput) {
  function updatePriceLabels() {
    if (parseFloat(priceMinInput.value) > parseFloat(priceMaxInput.value)) {
      // prevent overlap
      const temp = priceMinInput.value;
      priceMinInput.value = priceMaxInput.value;
      priceMaxInput.value = temp;
    }
    priceMin = parseFloat(priceMinInput.value);
    priceMax = parseFloat(priceMaxInput.value);
    priceMinLabel.textContent = `€${priceMin.toFixed(2)}`;
    priceMaxLabel.textContent = `€${priceMax.toFixed(2)}`;
    currentPage = 1;
    renderProducts();
  }

  priceMinInput.addEventListener("input", updatePriceLabels);
  priceMaxInput.addEventListener("input", updatePriceLabels);
}
});

// --- CART FUNCTIONS with Backend API ---
async function addToCart(id) {
  try {
    const response = await fetch('/api/cart', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ product_id: id, quantity: 1 }),
    });

    const data = await response.json();

    if (data.ok) {
      await loadCart(); // Reload cart from backend
      const product = products.find(p => p.id === id);
      if (product) {
        notie.alert({ type: 'success', text: `${product.name} added to cart!`, time: 2 });
      }

      // Animate cart badge
      const badge = document.getElementById("cart-count");
      if (badge) {
        badge.classList.remove("cart-animate");
        void badge.offsetWidth;
        badge.classList.add("cart-animate");
      }
    } else {
      if (data.message === "Not authenticated") {
        Swal.fire({
          title: "Please log in",
          text: "You need to log in to add items to cart.",
          icon: "info",
          confirmButtonText: "Go to Login"
        }).then(() => {
          window.location.href = "/login";
        });
      } else {
        notie.alert({ type: 'error', text: data.message || 'Failed to add to cart', time: 3 });
      }
    }
  } catch (error) {
    console.error('Error adding to cart:', error);
    notie.alert({ type: 'error', text: 'Failed to add to cart', time: 3 });
  }
}

async function removeFromCart(itemId) {
  const removed = cart.find(i => i.id === itemId);
  if (!removed) return;

  const result = await Swal.fire({
    title: "Remove Item?",
    text: `Are you sure you want to remove "${removed.product.name}" from the cart?`,
    icon: "warning",
    showCancelButton: true,
    confirmButtonText: "Yes, remove it",
    cancelButtonText: "Cancel"
  });

  if (result.isConfirmed) {
    try {
      const response = await fetch(`/api/cart/${itemId}`, {
        method: 'DELETE',
      });

      const data = await response.json();

      if (data.ok) {
        await loadCart();

        Swal.fire({
          title: "Removed!",
          text: `"${removed.product.name}" was removed from your cart.`,
          icon: "success",
          timer: 1500,
          showConfirmButton: false
        });

        // Animate cart badge
        const badge = document.getElementById("cart-count");
        if (badge) {
          badge.classList.remove("cart-animate");
          void badge.offsetWidth;
          badge.classList.add("cart-animate");
        }
      } else {
        notie.alert({ type: 'error', text: data.message || 'Failed to remove from cart', time: 3 });
      }
    } catch (error) {
      console.error('Error removing from cart:', error);
      notie.alert({ type: 'error', text: 'Failed to remove from cart', time: 3 });
    }
  }
}

async function updateQuantity(itemId, newQty) {
  try {
    const response = await fetch(`/api/cart/${itemId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ quantity: newQty }),
    });

    const data = await response.json();

    if (data.ok) {
      await loadCart();
      return true;
    } else {
      notie.alert({ type: 'error', text: data.message || 'Failed to update quantity', time: 3 });
      return false;
    }
  } catch (error) {
    console.error('Error updating quantity:', error);
    notie.alert({ type: 'error', text: 'Failed to update quantity', time: 3 });
    return false;
  }
}

async function decreaseQty(itemId) {
  const Toast = Swal.mixin({
    toast: true,
    position: "top",
    showConfirmButton: false,
    timer: 1500,
    timerProgressBar: true,
    didOpen: (toast) => {
      toast.onmouseenter = Swal.stopTimer;
      toast.onmouseleave = Swal.resumeTimer;
    }
  });
  
  const item = cart.find(i => i.id === itemId);
  if (!item) return;

  if (item.quantity > 1) {
    const success = await updateQuantity(itemId, item.quantity - 1);
    if (success) {
      Toast.fire({
        icon: "success",
        title: "Quantity decreased"
      });
    }
  } else {
    const result = await Swal.fire({
      title: "Remove Item?",
      text: `Quantity is 1. Do you want to remove "${item.product.name}" from the cart?`,
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Yes, remove it",
      cancelButtonText: "Cancel"
    });

    if (result.isConfirmed) {
      await removeFromCart(itemId);
    }
  }
}

async function increaseQty(itemId) {
  const Toast = Swal.mixin({
    toast: true,
    position: "top",
    showConfirmButton: false,
    timer: 1500,
    timerProgressBar: true,
    didOpen: (toast) => {
      toast.onmouseenter = Swal.stopTimer;
      toast.onmouseleave = Swal.resumeTimer;
    }
  });
  
  const item = cart.find(i => i.id === itemId);
  if (!item) return;

  const success = await updateQuantity(itemId, item.quantity + 1);
  if (success) {
    Toast.fire({
      icon: "success",
      title: "Quantity increased"
    });
  }
}

function renderCart() {
  const container = document.getElementById("cart-items");
  const totalEl = document.getElementById("cart-total");
  const countEl = document.getElementById("cart-count");

  if (!container) return;

  container.innerHTML = "";
  let total = 0;
  let count = 0;

  cart.forEach(item => {
    total += item.product.price * item.quantity;
    count += item.quantity;

    const div = document.createElement("div");
    div.className = "d-flex justify-content-between align-items-center mb-2 border-bottom pb-2";
    div.innerHTML = `
      <div>
        <strong>${item.product.name}</strong><br>
        <small>€${item.product.price.toFixed(2)} × ${item.quantity}</small>
      </div>
      <div>
        <button class="btn btn-sm btn-outline-secondary me-1" onclick="decreaseQty(${item.id})">-</button>
        <button class="btn btn-sm btn-outline-secondary me-1" onclick="increaseQty(${item.id})">+</button>
        <button class="btn btn-sm btn-danger" onclick="removeFromCart(${item.id})">&times;</button>
      </div>
    `;
    container.appendChild(div);
  });

  if (totalEl) totalEl.textContent = `€${total.toFixed(2)}`;
  if (countEl) countEl.textContent = count;
}

async function clearCart() {
  if (cart.length === 0) {
    Swal.fire("Cart is already empty!", "", "info");
    return;
  }

  const result = await Swal.fire({
    title: "Clear Cart?",
    text: "All items will be removed. This action cannot be undone.",
    icon: "warning",
    showCancelButton: true,
    confirmButtonText: "Yes, clear it",
    cancelButtonText: "Cancel"
  });

  if (result.isConfirmed) {
    try {
      const response = await fetch('/api/cart', {
        method: 'DELETE',
      });

      const data = await response.json();

      if (data.ok) {
        await loadCart();

        Swal.fire({
          title: "Cleared!",
          text: "Your cart has been emptied.",
          icon: "success",
          timer: 1500,
          showConfirmButton: false
        });

        const badge = document.getElementById("cart-count");
        if (badge) {
          badge.classList.remove("cart-animate");
          void badge.offsetWidth;
          badge.classList.add("cart-animate");
        }
      } else {
        notie.alert({ type: 'error', text: data.message || 'Failed to clear cart', time: 3 });
      }
    } catch (error) {
      console.error('Error clearing cart:', error);
      notie.alert({ type: 'error', text: 'Failed to clear cart', time: 3 });
    }
  }
}

function initCartToggle() {
  const openCartBtn = document.getElementById("open-cart");
  const cartPanel = document.getElementById("cart-panel");
  
  if (openCartBtn && cartPanel) {
    openCartBtn.addEventListener("click", () => {
      const bsOffcanvas = new bootstrap.Offcanvas(cartPanel);
      bsOffcanvas.show();
    });
  }
}

async function checkoutCart() {
  if (cart.length === 0) {
    Swal.fire("Cart is empty!", "Add some items before checking out.", "info");
    return;
  }

  const result = await Swal.fire({
    title: "Confirm Order",
    text: "Do you want to place this order?",
    icon: "question",
    showCancelButton: true,
    confirmButtonText: "Yes, place order",
    cancelButtonText: "Cancel"
  });

  if (result.isConfirmed) {
    try {
      const response = await fetch('/api/orders', {
        method: 'POST',
      });

      const data = await response.json();

      if (data.ok) {
        await loadCart(); // Cart should be empty now

        Swal.fire({
          title: "Order Placed!",
          text: `Your order #${data.data.id} has been placed successfully.`,
          icon: "success",
          confirmButtonText: "OK"
        }).then(() => {
          // Close the cart panel
          const cartPanel = document.getElementById("cart-panel");
          if (cartPanel) {
            const bsOffcanvas = bootstrap.Offcanvas.getInstance(cartPanel);
            if (bsOffcanvas) {
              bsOffcanvas.hide();
            }
          }
        });
      } else {
        Swal.fire("Error", data.message || "Failed to place order", "error");
      }
    } catch (error) {
      console.error('Error creating order:', error);
      Swal.fire("Error", "Failed to place order", "error");
    }
  }
}


  if (openCartBtn && cartPanel) {
    const bsOffcanvas = new bootstrap.Offcanvas(cartPanel);
    openCartBtn.addEventListener("click", () => bsOffcanvas.show());
  }
}


function checkoutCart() {
  if (cart.length === 0) {
    Swal.fire("Your cart is empty!", "Add some wheels before checking out.", "info");
    return;
  }

  Swal.fire({
    title: "Proceed to Checkout?",
    text: "You will be redirected to the checkout page to complete your order.",
    icon: "question",
    showCancelButton: true,
    confirmButtonColor: "#28a745",
    cancelButtonColor: "#6c757d",
    confirmButtonText: "Yes, checkout"
  }).then((result) => {
    if (result.isConfirmed) {
      // 🔥 Placeholder: redirect to a real checkout/payment page
      Swal.fire("Redirecting...", "Please wait while we take you to checkout.", "success");
      setTimeout(() => {
        window.location.href = "/checkout"; 
      }, 1500);
    }
  });
}
