const burger = document.querySelector(".burger");
const nav = document.querySelector(".nav-links");

burger.addEventListener("click", () => {
  nav.classList.toggle("nav-active"); // ouvre/ferme le menu
  burger.classList.toggle("toggle"); // anime le burger
});
document.querySelectorAll(".nav-links li a").forEach((link) => {
  link.addEventListener("click", () => {
    nav.classList.remove("nav-active");
    burger.classList.remove("toggle");
  });
});
