//#region Charger la map leaflet et openstreetmap

document.addEventListener("DOMContentLoaded", function () {
  const mapDiv = document.getElementById("map");

  // Récupère les coordonnées du dataset HTML
  const lat = parseFloat(mapDiv.dataset.lat) || 48.8566;
  const lon = parseFloat(mapDiv.dataset.lon) || 2.3522;

  const map = L.map("map").setView([lat, lon], 13);

  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution: "&copy; OpenStreetMap contributors",
  }).addTo(map);

  L.marker([lat, lon]).addTo(map).bindPopup("Votre position").openPopup();

  // Cree un point de depart de l'event sur la map
  // Variable pour stocker le marqueur de départ
  let startMarker = null;

  // Événement : double-clic sur la carte
  map.on("dblclick", function (e) {
    const lat = e.latlng.lat;
    const lng = e.latlng.lng;

    console.log("Position choisie :", lat, lng);

    // Supprimer le marqueur précédent s’il existe
    if (startMarker) {
      map.removeLayer(startMarker);
    }

    // Ajouter un nouveau marqueur rouge
    startMarker = L.marker([lat, lng], {
      icon: L.icon({
        iconUrl:
          "https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-red.png",
        shadowUrl:
          "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png",
        iconSize: [25, 41],
        iconAnchor: [12, 41],
        popupAnchor: [1, -34],
        shadowSize: [41, 41],
      }),
    })
      .addTo(map)
      .bindPopup("📍 Point de départ")
      .openPopup();
    // Récupération immédiate depuis le marqueur
    // Stocker les coordonnées dans deux variables distinctes
    const latitude = startMarker.getLatLng().lat;
    const longitude = startMarker.getLatLng().lng;
    // Mettre les coordonnées dans le champ input
    const input = document.getElementById("eventLocation");
    input.value = `Lat: ${latitude.toFixed(6)}, Lon: ${longitude.toFixed(6)}`;
  });
  //#endregion

  //#region Fonctions utilitaires
  // afficher ou cacher le mdp

  //#endregion
});
