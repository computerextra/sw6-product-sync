# Shopware 6 Product Sync

Ein kleines Stück Software, um Artikellisten von Wortmann und Kosatec mit einer Shopware Instanz zu Synchronisieren.

## Was wird benötigt

- [Go >=v1.23.3](https://go.dev/)

Das Programm wurde unter Windows 11 entwickelt. Eine Kompatibilität mit Linux Distributionen oder Mac OS ist weder vorgesehen noch getestet.

## Wie benutze ich das Skript.

Die Datei `.env-example` als `.env` Datei anlegen und alle Felder mit den korrekten Daten ausfüllen.
Die `config.?` Datei an die eigenen Bedürfnisse anpassen.

```sh
# Download all dependencies
go get .
# Build executeable
go build
```

Die Datei `?.exe` ausführen.

## Features

- Download von Preislisten der Hersteller _Wortmann_ und _Kosatec_
- NYI: Automatisches anlegen, pflegen und löschen der Artikel in Shopware 6
- NYI: Automatisches anlegen von Herstellern und Kategorien in Shopware 6
- NYI: Konfigurierbare Preise und Kategorien
- NYI: "Blacklisten" von Artikelnummern des Herstellers.
- NYI: "Blacklisten" von kompletten Kategorien des Herstellers.

## TODO

- [x] Download von Kosatec Preislisten
- [x] Download von Wortmann Preislisten
- [x] Download von Wortmann Bildern
- [x] Upload von Wortmann Bildern, damit eine URL fürs anlegen der Artikel generiert werden kann.
- [x] Löschen aller heruntergeladenen Dateien
- [x] Löschen allter hochgeladenen Bilder, da diese nicht mehr benötigt werden
- [x] Log erstellen
- [x] Log per Mail versenden
- [ ] Structs für die Artikellisten
- [ ] Kosatec Artikel sortieren
- [ ] Wortmann Artikel sortieren
- [ ] Config Datei erarbeiten (eventuell json/ini/eigenes Format?)
- [ ] Shop Artikel runterladen
- [ ] Shopartikel mit Artikellisten vergleichen
- [ ] Neue Liste für "Neue Artikel" erstellen
- [ ] Neue Liste für "Alte Artikel" erstellen
- [ ] Neue Liste für "EOL Artikel" erstellen
- [ ] "Neue Artikel" im Shop anlegen
- [ ] "Alte Artikel" im Shop anpassen
- [ ] "EOL Artikel" im Shop löschen
- [ ] Preis berechnung mit EK vom Händer mit Standard Aufschlag aus .env
- [ ] Artikel "blacklisten" nach Config File
- [ ] Preis berechnung mit besonderem Aufschlag mit Config json? Datei
  - [ ] besonderer Aufschlag als Prozent
  - [ ] besonderer Aufschlag als Absolut
  - [ ] besonderer Aufschlag nach Artikel Kategorie
  - [ ] besonderer Aufschlag nach Artikelnummer (bsp: AVM UVP Liste)
- [ ] Anlegen von Herstellern
- [ ] Anlegen und Pflegen von Kategorien
  - [ ] Standard Kategorien nach Config json? Datei
- [ ] Support für Intos Artikellisten einbauen
- [ ] Nicht lieferbare Artikel im Shop ausblenden
