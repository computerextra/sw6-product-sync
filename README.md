# Shopware 6 Product Sync

Ein kleines Stück Software, um Artikellisten von Wortmann und Kosatec mit einer Shopware Instanz zu Synchronisieren.

## Was wird benötigt

- [Go >=v1.23.3](https://go.dev/)

Das Programm wurde unter Windows 11 entwickelt. Eine Kompatibilität mit Linux Distributionen oder Mac OS ist weder vorgesehen noch getestet.

## Wie benutze ich das Skript.

Die Datei `.env-example` als `.env` Datei anlegen und alle Felder mit den korrekten Daten ausfüllen.
Die `config.toml` Datei an die eigenen Bedürfnisse anpassen.

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
- Konfigurierbare Preise und Kategorien
- "Blacklisten" von Artikelnummern des Herstellers.
- "Blacklisten" von kompletten Kategorien des Herstellers.

## TODO

- [x] Download von Kosatec Preislisten
- [x] Download von Wortmann Preislisten
- [x] Download von Wortmann Bildern
- [x] Upload von Wortmann Bildern, damit eine URL fürs anlegen der Artikel generiert werden kann.
- [x] Löschen aller heruntergeladenen Dateien
- [x] Löschen allter hochgeladenen Bilder, da diese nicht mehr benötigt werden
- [x] Log erstellen
- [x] Log per Mail versenden
- [x] Structs für die Artikellisten
- [x] Kosatec Artikel sortieren
- [x] Wortmann Artikel sortieren
- [x] Config Datei erarbeiten
- [x] Shop Artikel runterladen
- [x] Shopartikel mit Artikellisten vergleichen
- [x] Neue Liste für "Neue Artikel" erstellen
- [x] Neue Liste für "Alte Artikel" erstellen
- [x] Neue Liste für "EOL Artikel" erstellen
- [x] "Neue Artikel" im Shop anlegen
- [x] "Alte Artikel" im Shop anpassen
- [x] Preis berechnung mit EK vom Händer mit Standard Aufschlag aus .env
- [x] Artikel "blacklisten" nach Config File
- [x] Preis berechnung mit besonderem Aufschlag mit Config json? Datei
  - [x] besonderer Aufschlag als Prozent
  - [x] besonderer Aufschlag als Absolut
  - [x] besonderer Aufschlag nach Artikel Kategorie
  - [x] besonderer Aufschlag nach Artikelnummer (bsp: AVM UVP Liste)
- [x] Anlegen von Herstellern
- [x] Anlegen und Pflegen von Kategorien
  - [x] Standard Kategorien nach Config Datei
- [x] Nicht lieferbare Artikel im Shop ausblenden
- [ ] Support für Intos Artikellisten einbauen
- [ ] Bilder hochladen
- [ ] Bilder zu Produkten hinzufügen
- [ ] Produkt Cover setzen
- [ ] "EOL Artikel" im Shop löschen
