package utils

import (
	"RideUP/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func ReverseGeocodeCoord(address string) (float64, float64, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", address)
	params.Set("format", "json")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var result models.NominatimResponseCoord
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}

	if len(result) == 0 {
		return 0, 0, fmt.Errorf("adresse introuvable")
	}

	lat, err := strconv.ParseFloat(result[0].Lat, 64)
	if err != nil {
		return 0, 0, err
	}
	lon, err := strconv.ParseFloat(result[0].Lon, 64)
	if err != nil {
		return 0, 0, err
	}

	return lat, lon, nil
}
