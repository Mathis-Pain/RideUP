const burger = document.querySelector(".burger");
const nav = document.querySelector(".nav-links");

// Burger menu pour mobile
burger.addEventListener("click", () => {
  nav.classList.toggle("nav-active");
  burger.classList.toggle("toggle");
});

// Dropdown pour profil
const dropdown = document.querySelector(".dropdown");
const dropbtn = document.querySelector(".dropbtn");

if (dropbtn && dropdown) {
  dropbtn.addEventListener("click", (e) => {
    e.preventDefault();
    dropdown.classList.toggle("active");
  });

  // Fermer le dropdown si on clique ailleurs
  document.addEventListener("click", (e) => {
    if (!dropdown.contains(e.target)) {
      dropdown.classList.remove("active");
    }
  });
}

// Fermer le menu mobile quand on clique sur un lien (SAUF le dropbtn)
document.querySelectorAll(".nav-links li a").forEach((link) => {
  link.addEventListener("click", (e) => {
    // Ne ferme pas si c'est le bouton dropdown
    if (!link.classList.contains("dropbtn")) {
      nav.classList.remove("nav-active");
      burger.classList.remove("toggle");
    }
  });
});
