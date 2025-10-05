


const products = [
  { id: 1, name: "Polyurethane Wheel Ø80mm", price: 19.99, type: "polyurethane", image: "images/wheel1.jpg" },
  { id: 2, name: "Nylon Wheel Ø70mm", price: 14.50, type: "nylon", image: "images/wheel2.jpg" },
  { id: 3, name: "Rubber Coated Wheel Ø90mm", price: 22.00, type: "rubber", image: "images/wheel3.jpg" },
  { id: 4, name: "Polyurethane Wheel Ø100mm", price: 25.00, type: "polyurethane", image: "images/wheel4.jpg" },
  { id: 5, name: "Nylon Wheel Ø85mm", price: 17.80, type: "nylon", image: "images/wheel5.jpg" },
  { id: 6, name: "Rubber Wheel Ø75mm", price: 16.20, type: "rubber", image: "images/wheel6.jpg" },
  { id: 7, name: "Polyurethane Wheel Ø110mm", price: 28.40, type: "polyurethane", image: "images/wheel7.jpg" },
  { id: 8, name: "Nylon Heavy Duty Ø95mm", price: 20.00, type: "nylon", image: "images/wheel8.jpg" },
  { id: 9, name: "Rubber Shock-Absorb Ø100mm", price: 27.50, type: "rubber", image: "images/wheel9.jpg" },
  { id: 10, name: "Polyurethane Silent Ø90mm", price: 23.90, type: "polyurethane", image: "images/wheel10.jpg" }
];




let cart = JSON.parse(localStorage.getItem("cart")) || [];

let currentPage = 1;
const itemsPerPage = 6;
let currentFilter = "all";
let currentSearch = "";
let priceMin = 10;
let priceMax = 30;



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
  if (currentPage > totalPages) currentPage = 1;
  const start = (currentPage - 1) * itemsPerPage;
  const paginated = filtered.slice(start, start + itemsPerPage);

  // Render products
  paginated.forEach(p => {
    const card = document.createElement("div");
    card.className = "col-md-4 mb-4";
    card.innerHTML = `
      <div class="card h-100 shadow-sm">
        <img src="${p.image}" class="card-img-top" alt="${p.name}">
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

// --- Filters ---
document.addEventListener("DOMContentLoaded", () => {
  renderProducts();
  renderCart();
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

// --- CART FUNCTIONS (same as before) ---
function addToCart(id) {
 const item = cart.find(i => i.id === id);
  if (item) {
    item.qty++;
    notie.alert({ type: 'info', text: `${item.name} quantity updated (x${item.qty})`, time: 2 });
  } else {
    const product = products.find(p => p.id === id);
    cart.push({ ...product, qty: 1 });
    notie.alert({ type: 'success', text: `${product.name} added to cart!`, time: 2 });
  }
  saveCart();
  renderCart();

  // Animate cart badge
  const badge = document.getElementById("cart-count");
  if (badge) {
    badge.classList.remove("cart-animate");
    void badge.offsetWidth;
    badge.classList.add("cart-animate");
  }
}
function removeFromCart(id) {
  const removed = cart.find(i => i.id === id);
  if (!removed) return;

  swal.fire({
    title: "Remove Item?",
    text: `Are you sure you want to remove "${removed.name}" from the cart?`,
    icon: "warning",
    showCancelButton: true,
    confirmButtonText: "Yes, remove it",
    cancelButtonText: "Cancel"
  }).then((result) => {
    if (result.isConfirmed) {
      // ✅ Update data and UI only if confirmed
      cart = cart.filter(item => item.id !== id);
      saveCart();
      renderCart();

      swal.fire({
        title: "Removed!",
        text: `"${removed.name}" was removed from your cart.`,
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
    } else if (result.dismiss === Swal.DismissReason.cancel) {
      // ❌ Cancel pressed — do nothing, no UI change
      // Optionally: a subtle info message or no alert at all
      // swalTheme.fire("Cancelled", "Item was not removed.", "info");
    }
  });
}

function decreaseQty(id) {
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
  const item = cart.find(i => i.id === id);
  if (!item) return;

  if (item.qty > 1) {
    item.qty--;
    Toast.fire({
      icon: "success",
      title: "Item removed succesfully"
    });
    saveCart();
    renderCart();
  } else {
    Swal.fire({
      title: "Remove Item?",
      text: `Quantity is 1. Do you want to remove "${item.name}" from the cart?`,
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Yes, remove it",
      cancelButtonText: "Cancel"
    }).then((result) => {
      if (result.isConfirmed) {
        cart = cart.filter(i => i.id !== id);
        saveCart();
        renderCart();

        Swal.fire({
          title: "Removed!",
          text: `"${item.name}" was removed.`,
          icon: "success",
          timer: 1500,
          showConfirmButton: false
        });
      }
    });
  }
}

function increaseQty(id){
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
  const item = cart.find(i => i.id === id);
  item.qty++;
  Toast.fire({
      icon: "success",
      title: "Item added succesfully"
    });
  saveCart();
  renderCart();
}

function saveCart() {
  localStorage.setItem("cart", JSON.stringify(cart));
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
    total += item.price * item.qty;
    count += item.qty;

    const div = document.createElement("div");
    div.className = "d-flex justify-content-between align-items-center mb-2 border-bottom pb-2";
    div.innerHTML = `
      <div>
        <strong>${item.name}</strong><br>
        <small>€${item.price.toFixed(2)} × ${item.qty}</small>
      </div>
      <div>
        <button class="btn btn-sm btn-outline-secondary me-1" onclick="decreaseQty(${item.id})">-</button>
        <button class="btn btn-sm btn-outline-secondary me-1" onclick="increaseQty(${item.id})">+</button>
        <button class="btn btn-sm btn-danger" onclick="removeFromCart(${item.id})">&times;</button>
      </div>
    `;
    container.appendChild(div);
  });

  totalEl.textContent = `€${total.toFixed(2)}`;
  countEl.textContent = count;
}

function clearCart() {
  if (cart.length === 0) {
    swal.fire("Cart is already empty!", "", "info");
    return;
  }

  swal.fire({
    title: "Clear Cart?",
    text: "All items will be removed. This action cannot be undone.",
    icon: "warning",
    showCancelButton: true,
    confirmButtonText: "Yes, clear it",
    cancelButtonText: "Cancel"
  }).then((result) => {
    if (result.isConfirmed) {
      // ✅ Clear data only when confirmed
      cart = [];
      saveCart();
      renderCart();

      swal.fire({
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
    } else if (result.dismiss === Swal.DismissReason.cancel) {
      // ❌ Do nothing
      // swalTheme.fire("Cancelled", "Your cart was not cleared.", "info");
    }
  });
}

function initCartToggle() {
  const openCartBtn = document.getElementById("open-cart");
  const cartPanel = document.getElementById("cart-panel");

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
