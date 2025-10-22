document.addEventListener("DOMContentLoaded", () => {
  const burger = document.querySelector(".burger");
  const nav = document.querySelector(".nav-links");
  const dropdown = document.querySelector(".dropdown");
  const dropbtn = document.querySelector(".dropbtn");

  console.log(
    "Nombre de boutons trouvés :",
    document.querySelectorAll(".nav-links li a").length
  );

  // Burger menu
  if (burger && nav) {
    burger.addEventListener("click", () => {
      nav.classList.toggle("nav-active");
      burger.classList.toggle("toggle");
    });
  }

  // Dropdown
  if (dropbtn && dropdown) {
    dropbtn.addEventListener("click", (e) => {
      e.preventDefault();
      dropdown.classList.toggle("active");
    });

    document.addEventListener("click", (e) => {
      if (!dropdown.contains(e.target) && !dropbtn.contains(e.target)) {
        dropdown.classList.remove("active");
      }
    });
  }

  // Fermer le menu mobile au clic sur un lien
  document.querySelectorAll(".nav-links li a").forEach((link) => {
    link.addEventListener("click", (e) => {
      if (!link.classList.contains("dropbtn")) {
        nav.classList.remove("nav-active");
        burger.classList.remove("toggle");
      }
    });
  });
});
