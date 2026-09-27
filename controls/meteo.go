// Copyright (c) 2026 Murilo Gomes <profmugomes.com.br>. All Rights Reserved.

// Licensed under the PolyForm Strict License 1.0.0.
// See LICENSE.md for details.

package controls

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func BuscarClima(lat, lon string) (map[string]any, error) {
	sURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current_weather=true&hourly=temperature_2m,precipitation,snowfall&daily=temperature_2m_max,temperature_2m_min,precipitation_sum,snowfall_sum&timezone=auto",
		lat,
		lon,
	)

	resp, err := http.Get(sURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro Open-Meteo: %s", resp.Status)
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func ClimaAtual(data map[string]any) []string {
	current, ok := data["current_weather"].(map[string]any)
	if !ok {
		return []string{"Indisponível"}
	}

	temp, _ := current["temperature"].(float64)
	wind, _ := current["windspeed"].(float64)
	time, _ := current["time"].(string)

	return []string{time, strconv.FormatFloat(temp, 'f', 2, 64) + " °C", strconv.FormatFloat(wind, 'f', 2, 64) + " mm\n"}
}

func PrevisaoHoraria(data map[string]any, hours int) string {
	hourly, ok := data["hourly"].(map[string]any)
	if !ok {
		return "Previsão horária indisponível"
	}

	times, _ := hourly["time"].([]any)
	temps, _ := hourly["temperature_2m"].([]any)
	rain, _ := hourly["precipitation"].([]any)

	var b strings.Builder
	b.WriteString("⏱ Previsão horária:\n")

	for i := 0; i < hours && i < len(times); i++ {
		t, _ := temps[i].(float64)
		r, _ := rain[i].(float64)

		b.WriteString(fmt.Sprintf(
			"%s | 🌡 %.1f °C | 🌧 %.1f mm\n",
			times[i], t, r,
		))
	}
	return b.String()
}

func PrevisaoDiaria(data map[string]any, days int) ([]string, string) {
	daily, ok := data["daily"].(map[string]any)
	if !ok {
		return nil, "Previsão diária indisponível"
	}

	dates := daily["time"].([]any)
	max := daily["temperature_2m_max"].([]any)
	min := daily["temperature_2m_min"].([]any)
	rain := daily["precipitation_sum"].([]any)
	snow := daily["snowfall_sum"].([]any)

	var c []string
	for i := 0; i < days && i < len(dates); i++ {
		c = append(c, fmt.Sprintf(
			"%s | %.1f °C / %.1f °C | %.1f mm | %.1f mm",
			dates[i],
			max[i],
			min[i],
			rain[i],
			snow[i],
		))
	}

	return c, ""
}

func ResumoClima(data map[string]any) string {
	current := ClimaAtual(data)
	hourly := PrevisaoHoraria(data, 6)
	daily, _ := PrevisaoDiaria(data, 5)

	return fmt.Sprintf(
		"%s\n\n%s\n%s",
		current,
		hourly,
		daily,
	)
}
