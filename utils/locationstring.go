package utils

import (
	"RideUP/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func ReverseGeocodeSimple(lat, lon float64) (*models.SimpleAddress, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json", lat, lon)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "RideUpApp/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data models.NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	addr := &models.SimpleAddress{
		Numero:     data.Address.HouseNumber,
		Rue:        data.Address.Road,
		CodePostal: data.Address.Postcode,
		Ville:      data.Address.City,
		Pays:       data.Address.Country,
	}

	return addr, nil
}
