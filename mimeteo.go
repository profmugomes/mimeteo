// Copyright (c) 2026 Murilo Gomes <profmugomes.com.br>. All Rights Reserved.

// Licensed under the PolyForm Perimeter License 1.0.1.
// See LICENSE.md for details.

package main

import (
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	c "profmugomes/mimeteo/controls"

	"github.com/profmugomes/mgcolumnview/v2"
	"github.com/profmugomes/mgsettings/v4"
	"github.com/profmugomes/mgsmartflow/v2"
)

const VERSION_APP = "2.0.0"

var a fyne.App

func main() {
	sIcon := fyne.NewStaticResource("mimeteo.png", resourceIconPngData)

	a = app.NewWithID("br.com.profmugomes.mimeteo")
	a.Settings().SetTheme(&darkTheme{})
	a.SetIcon(sIcon)

	frmMain := a.NewWindow("MiMeteo")
	frmMain.CenterOnScreen()
	frmMain.SetFixedSize(true)
	frmMain.Resize(fyne.NewSize(600, 400))

	mnuSobre := fyne.NewMenu("Sobre",
		fyne.NewMenuItem(
			"Verificar Atualização", func() {
				sURL, _ := url.Parse("https://github.com/profmugomes/mimeteo/releases")
				a.OpenURL(sURL)
			},
		),
		fyne.NewMenuItem(
			"Apoie MiMeteo", func() {
				sURL, _ := url.Parse("https:///github.com/sponsors/profmugomes")
				a.OpenURL(sURL)
			},
		),
		fyne.NewMenuItem("Sobre MiMeteo", func() {
			showAbout()
		}),
	)

	frmMain.SetMainMenu(fyne.NewMainMenu(mnuSobre))

	config, _ := mgsettings.Load("mimeteo", true)

	flow := mgsmartflow.New()

	lblEstado := widget.NewLabel("Estado")
	lblEstado.TextStyle = fyne.TextStyle{Bold: true}

	lblCidade := widget.NewLabel("Cidade")
	lblCidade.TextStyle = fyne.TextStyle{Bold: true}

	cboCidade := widget.NewSelect(nil, nil)

	lstEstado, _ := c.ListarEstados()
	cboEstado := widget.NewSelect(lstEstado, func(s string) {
		sData := strings.Split(s, " - ")
		if len(sData) > 0 {
			lstCidades, _ := c.ListarMunicipiosPorUF(sData[0])
			cboCidade.SetOptions(lstCidades)
			config.SetString("estado", s)
			config.Save()
		}
	})

	cboCidade.OnChanged = func(s string) {
		config.SetString("cidade", s)
		config.Save()
	}

	flow.AddColumn(lblEstado, lblCidade)
	flow.AddColumn(cboEstado, cboCidade)

	flow.Gap(cboEstado, 7, 17)

	cboEstado.SetSelected(config.GetString("estado", ""))
	cboCidade.SetSelected(config.GetString("cidade", ""))

	cvTempAtual := mgcolumnview.NewColumnView(
		[]string{"Horário", "Temperatura", "Vento"},
		[]float32{100, 140, 200}, false,
	)

	cvTempProximo := mgcolumnview.NewColumnView(
		[]string{"Horário", "Temperatura", "Vento", "Neve"},
		[]float32{100, 140, 95, 200}, false,
	)
	btnVerificar := widget.NewButton("Verificar", func() {
		estado := strings.Split(config.GetString("estado", "SP"), " - ")
		if len(estado) > 0 && cboCidade.Selected != "" {
			lat, log, _ := c.BuscarLocalizacaoPorNome(cboCidade.Selected, estado[0])
			tempo, _ := c.BuscarClima(lat, log)
			sClimaAtual := c.ClimaAtual(tempo)
			sClimaDataAtual := strings.Split(sClimaAtual[0], "T")
			cvTempAtual.RemoveAll()
			cvTempAtual.AddRow([]string{sClimaDataAtual[0], sClimaAtual[1], sClimaAtual[2]})

			sClimaProximo, _ := c.PrevisaoDiaria(tempo, 6)
			cvTempProximo.RemoveAll()

			for _, row := range sClimaProximo {
				cvTempProximo.AddRow(strings.Split(row, " | "))
			}
		}
	})

	flow.AddRow(container.NewHBox(layout.NewSpacer(), btnVerificar, layout.NewSpacer()))
	flow.AddRow(widget.NewSeparator())

	tabs := container.NewAppTabs(
		container.NewTabItem("Temperatura Atual", cvTempAtual),
		container.NewTabItem("Temperatura (Próximos Dias)", cvTempProximo),
	)

	flow.AddRow(tabs)
	flow.Resize(tabs, frmMain.Canvas().Size().Width, frmMain.Canvas().Size().Height-7)

	frmMain.SetContent(flow.Container)
	frmMain.ShowAndRun()
}
