document.addEventListener("DOMContentLoaded", () => {
  const buttons = document.querySelectorAll(".join-btn");
  console.log("Nombre de boutons trouvés :", buttons.length);
  const mapDiv = document.getElementById("map");

  // Récupère les coordonnées du dataset HTML
  const lat = parseFloat(mapDiv.dataset.lat) || 48.8566;
  const lon = parseFloat(mapDiv.dataset.lon) || 2.3522;

  // Initialisation de la carte Leaflet
  const map = L.map("map").setView([lat, lon], 13);

  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution: "&copy; OpenStreetMap contributors",
  }).addTo(map);

  // Marqueur de position actuelle (bleu)
  L.marker([lat, lon]).addTo(map).bindPopup("📍 Votre position").openPopup();

  // Récupérer et afficher les événements disponibles sur la carte
  try {
    const eventsData = mapDiv.dataset.events;
    console.log("🔍 Data brute:", eventsData);

    if (eventsData) {
      const events = JSON.parse(eventsData);
      console.log("✅ Events parsés:", events);

      // 🟢 BOUCLE AVEC ICÔNE CONDITIONNELLE
      events.forEach((event) => {
        console.log("📌 Traitement event:", event);

        if (event.latitude && event.longitude) {
          // ✅ CHOIX DE L'ICÔNE SELON user_joined
          const iconColor = event.user_joined ? "violet" : "red";
          const iconUrl = `https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-${iconColor}.png`;

          console.log(
            `🎨 Event ${event.id}: ${
              event.user_joined ? "REJOINT (violet)" : "NON REJOINT (rouge)"
            }`
          );

          const marker = L.marker([event.latitude, event.longitude], {
            icon: L.icon({
              iconUrl: iconUrl,
              iconSize: [25, 41],
              iconAnchor: [12, 41],
              popupAnchor: [1, -34],
              shadowUrl:
                "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
              shadowSize: [41, 41],
            }),
          }).addTo(map);

          // 🟢 Conversion correcte UTC → Paris
          let formattedDate = "Date non définie";

          if (event.start_datetime) {
            const date = new Date(event.start_datetime);

            formattedDate = date
              .toLocaleString("fr-FR", {
                timeZone: "Europe/Paris",
                day: "2-digit",
                month: "2-digit",
                year: "numeric",
                hour: "2-digit",
                minute: "2-digit",
              })
              .replace(",", " à");

            formattedDate = `le ${formattedDate}`;
          }

          const popupContent = `
            <div class="event-popup">
              <h4>${event.title || "Sans titre"}</h4>
              <p><strong>Créateur:</strong> ${
                event.creator_name || "Inconnu"
              }</p>
              <p><strong>Date:</strong> ${formattedDate}</p>
              <p><strong>Lieu:</strong> ${event.address || "Non précisé"}</p>
              ${event.description ? `<p>${event.description}</p>` : ""}
              <p><strong>Participants:</strong> ${event.participants || 0}</p>
              <button class="join-btn-popup" data-event-id="${event.id}">
                ${event.user_joined ? "Annuler" : "Rejoindre"}
              </button>
            </div>
          `;

          marker.bindPopup(popupContent);
        } else {
          console.warn("⚠️ Event sans coordonnées:", event);
        }
      });

      console.log(`✅ ${events.length} marqueurs ajoutés à la carte`);
    } else {
      console.warn("⚠️ Aucune donnée d'événements trouvée");
    }
  } catch (error) {
    console.error("❌ Erreur lors du chargement des événements:", error);
  }

  // Gestion des clics sur les boutons (existant + nouveaux dans les popups)
  document.addEventListener("click", async (e) => {
    if (
      e.target.classList.contains("join-btn") ||
      e.target.classList.contains("join-btn-popup")
    ) {
      const btn = e.target;
      const eventId = btn.dataset.eventId;
      const action = btn.textContent.trim().toLowerCase();
      console.log("Action détectée :", action, "pour event", eventId);

      try {
        // 🔹 Suppression d'un événement (propriétaire)
        if (action === "supprimer") {
          if (!confirm("Voulez-vous vraiment supprimer cet événement ?"))
            return;

          const response = await fetch("/RideUp", {
            method: "POST",
            headers: {"Content-Type": "application/x-www-form-urlencoded"},
            body: `event_id=${eventId}&action=delete`,
          });

          if (!response.ok)
            throw new Error("Erreur serveur lors de la suppression");

          const data = await response.json();

          if (data.success) {
            window.location.href = "/RideUp";
          } else {
            alert("Impossible de supprimer cet événement.");
          }

          return;
        }

        // 🔹 Gestion du join / leave
        const actionType = action === "rejoindre" ? "join" : "leave";
        const response = await fetch("/JoinEvent", {
          method: "POST",
          headers: {"Content-Type": "application/x-www-form-urlencoded"},
          body: `event_id=${eventId}&action=${actionType}`,
        });

        if (!response.ok) throw new Error("Erreur serveur");
        const data = await response.json();
        console.log("Réponse serveur :", data);

        window.location.href = "/RideUp";
      } catch (err) {
        console.error("Erreur :", err);
        alert("Une erreur est survenue, veuillez réessayer.");
      }
    }
  });
});
