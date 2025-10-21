document.addEventListener("DOMContentLoaded", () => {
  const buttons = document.querySelectorAll(".join-btn");
  console.log("Nombre de boutons trouvés :", buttons.length);

  buttons.forEach((btn) => {
    btn.addEventListener("click", async () => {
      console.log("Bouton cliqué !");
      const eventId = btn.dataset.eventId;
      const action = btn.textContent.trim() === "Rejoindre" ? "join" : "leave";

      try {
        const response = await fetch("/JoinEvent", {
          method: "POST",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: `event_id=${eventId}&action=${action}`,
        });

        if (!response.ok) throw new Error("Erreur serveur");
        const data = await response.json();
        console.log("Réponse serveur :", data);

        // 🔹 1. Met à jour le texte du bouton
        btn.textContent = data.joined ? "Annuler" : "Rejoindre";

        // 🔹 2. Met à jour le nombre de participants dans la carte
        const card = btn.closest(".card"); // remonte jusqu’à la carte parente
        const counter = card.querySelector(".participants-count"); // trouve le <span>
        if (counter) {
          counter.textContent = data.participants; // remplace par la nouvelle valeur
        }
      } catch (err) {
        console.error("Erreur :", err);
      }
    });
  });
});
