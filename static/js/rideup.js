document.addEventListener("DOMContentLoaded", () => {
  const buttons = document.querySelectorAll(".join-btn");
  console.log("Nombre de boutons trouvés :", buttons.length);

  buttons.forEach((btn) => {
    btn.addEventListener("click", async () => {
      const eventId = btn.dataset.eventId;
      const action = btn.textContent.trim().toLowerCase(); // "rejoindre", "annuler" ou "supprimer"
      console.log("Action détectée :", action, "pour event", eventId);

      try {
        // 🔹 Suppression d’un événement (propriétaire)
        if (action === "supprimer") {
          if (!confirm("Voulez-vous vraiment supprimer cet événement ?"))
            return;

          const response = await fetch("/RideUp", {
            method: "POST",
            headers: { "Content-Type": "application/x-www-form-urlencoded" },
            body: `event_id=${eventId}&action=delete`,
          });

          if (!response.ok)
            throw new Error("Erreur serveur lors de la suppression");

          const data = await response.json();

          if (data.success) {
            // ✅ Recharge la page pour actualiser la liste
            window.location.href = "/RideUp";
          } else {
            alert("Impossible de supprimer cet événement.");
          }

          return; // on sort ici
        }

        // 🔹 Gestion du join / leave
        const actionType = action === "rejoindre" ? "join" : "leave";
        const response = await fetch("/JoinEvent", {
          method: "POST",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: `event_id=${eventId}&action=${actionType}`,
        });

        if (!response.ok) throw new Error("Erreur serveur");
        const data = await response.json();
        console.log("Réponse serveur :", data);

        // ✅ Recharge la page après rejoindre / annuler
        window.location.href = "/RideUp";
      } catch (err) {
        console.error("Erreur :", err);
        alert("Une erreur est survenue, veuillez réessayer.");
      }
    });
  });
});
