document.addEventListener("DOMContentLoaded", () => {
  const mapDiv = document.getElementById("map");
  const lat = parseFloat(mapDiv.dataset.lat) || 48.8566;
  const lon = parseFloat(mapDiv.dataset.lon) || 2.3522;
  const currentUserId = parseInt(mapDiv.dataset.userId);

  // initilalisation de la map coordonée et zoom
  const map = L.map("map").setView([lat, lon], 8);

  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution: "&copy; OpenStreetMap contributors",
  }).addTo(map);

  L.marker([lat, lon]).addTo(map).bindPopup("📍 Votre position").openPopup();

  //  STOCKAGE DES MARQUEURS POUR Y ACCÉDER PLUS TARD
  const markers = {};

  try {
    // dataset genere automatiquement par le navigateur data-events → dataset.event
    const eventsData = mapDiv.dataset.events;

    if (eventsData) {
      const events = JSON.parse(eventsData);

      events.forEach((event) => {
        if (event.latitude && event.longitude) {
          let iconColor;
          if (event.created_by === currentUserId) {
            iconColor = "green";
          } else if (event.user_joined) {
            iconColor = "violet";
          } else {
            iconColor = "red";
          }

          const iconUrl = `https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-${iconColor}.png`;

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

          // STOCKE LE MARQUEUR AVEC L'ID DE L'ÉVÉNEMENT
          markers[event.id] = marker;

          let formattedStart = "Non défini";
          if (event.start_time) {
            const date = new Date(event.start_time);
            formattedStart = date.toLocaleString("fr-FR", {
              day: "2-digit",
              month: "2-digit",
              year: "numeric",
              hour: "2-digit",
              minute: "2-digit",
            });
          }

          let popupContent = `
            <div style="max-width: 300px;">
              <h3 style="margin: 0 0 10px 0; color: #333;">${event.title}</h3>
              <p style="margin: 5px 0;"><strong>📅 Date :</strong> ${formattedStart}</p>
              <p style="margin: 5px 0;"><strong>📍 Lieu :</strong> ${event.address}</p>
              <p style="margin: 5px 0;"><strong>👤 Organisateur :</strong> ${event.creator_name}</p>
          `;

          if (event.description) {
            popupContent += `<p style="margin: 10px 0 5px 0;"><strong>Description :</strong></p><p style="margin: 0;">${event.description}</p>`;
          }

          popupContent += `<p style="margin: 10px 0 5px 0;"><strong>👥 Participants :</strong> ${event.participants}</p>`;

          if (event.created_by !== currentUserId) {
            const buttonText = event.user_joined ? "QUITTER" : "REJOINDRE";
            const buttonClass = event.user_joined ? "joined" : "";
            popupContent += `
              <button class="join-btn-popup ${buttonClass}" data-event-id="${
              event.id
            }" style="
                margin-top: 10px;
                padding: 8px 16px;
                background-color: ${event.user_joined ? "#dc3545" : "#007bff"};
                color: white;
                border: none;
                border-radius: 4px;
                cursor: pointer;
                font-size: 14px;
              ">${buttonText}</button>
            `;
          }

          popupContent += `</div>`;

          marker.bindPopup(popupContent);
        }
      });
    }
  } catch (error) {
    console.error("❌ Erreur lors du parsing des events:", error);
  }

  // ✅ NOUVELLE FONCTIONNALITÉ : CLIC SUR UNE CARD
  document.querySelectorAll(".card").forEach((card) => {
    card.addEventListener("click", (e) => {
      // Ignore le clic si c'est un bouton
      if (
        e.target.classList.contains("join-btn") ||
        e.target.classList.contains("delete-btn") ||
        e.target.closest(".join-btn") ||
        e.target.closest(".delete-btn")
      ) {
        return;
      }

      const eventId = parseInt(card.dataset.eventId);
      const cardLat = parseFloat(card.dataset.lat);
      const cardLon = parseFloat(card.dataset.lon);

      console.log(
        `🎯 Card cliquée: Event ${eventId} à [${cardLat}, ${cardLon}]`
      );

      if (markers[eventId] && !isNaN(cardLat) && !isNaN(cardLon)) {
        // ✅ CENTRE LA CARTE AVEC OFFSET (pour ne pas cacher le marqueur sous la popup)
        map.setView([cardLat, cardLon], 16, {
          animate: true,
          duration: 0.5,
        });

        // ✅ OUVRE LA POPUP APRÈS UN LÉGER DÉLAI (pour que le centrage soit visible)
        setTimeout(() => {
          markers[eventId].openPopup();
        }, 300);

        // ✅ SCROLL VERS LA CARTE (optionnel)
        document.getElementById("map").scrollIntoView({
          behavior: "smooth",
          block: "center",
        });
      }
    });
  });

  // ✅ GESTION DES BOUTONS
  document.addEventListener("click", async (e) => {
    if (
      e.target.classList.contains("join-btn") ||
      e.target.classList.contains("delete-btn") ||
      e.target.classList.contains("join-btn-popup")
    ) {
      const btn = e.target;
      const eventId = btn.dataset.eventId;
      const action = btn.textContent.trim().toLowerCase();
      console.log("Action détectée :", action, "pour event", eventId);

      // ✅ CENTRAGE SUR LA CARTE AVANT L'ACTION
      const card = btn.closest(".card");
      if (card) {
        const cardLat = parseFloat(card.dataset.lat);
        const cardLon = parseFloat(card.dataset.lon);

        if (markers[eventId] && !isNaN(cardLat) && !isNaN(cardLon)) {
          map.setView([cardLat, cardLon], 16, {
            animate: true,
            duration: 0.5,
          });

          setTimeout(() => {
            markers[eventId].openPopup();
          }, 300);

          document.getElementById("map").scrollIntoView({
            behavior: "smooth",
            block: "center",
          });
        }
      }

      try {
        if (action.includes("supprimer")) {
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

        const actionType = action.includes("rejoindre") ? "join" : "leave";
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
