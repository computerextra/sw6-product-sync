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

Die Datei `sw6-product-sync.exe` ausführen.

## CLI-Argumente

`-h` - Zeigt die Hilfe-Seite
`-endless` - Startet das Programm in Dauerschleife, nur mit `-wait` zu verwenden
`-wait` - Setzt ein Timeout in Stunden für die wiederholungen im Endlos Modus
Bsp: `./sw6-product-sync.exe -endless -wait=4` - Startet das Programm im Endlos Modus alle 4 Stunden
Ohne Befehle startet das Program einmalig und beendet danach.

### Danger CLI

`-delete-products` - Löscht alle Produkte

## Features

- Download von Preislisten der Hersteller _Wortmann_ und _Kosatec_
- Automatisches anlegen, pflegen und löschen der Artikel in Shopware 6
- Automatisches anlegen von Herstellern und Kategorien in Shopware 6
- Konfigurierbare Preise und Kategorien
- "Blacklisten" von Artikelnummern des Herstellers.
- "Blacklisten" von kompletten Kategorien des Herstellers.

## TODO

- [ ] "Neue Artikel" im Shop anlegen (Kann optimiert werden)
- [ ] "Alte Artikel" im Shop anpassen (Kann optimiert werden)
- [ ] Separator für CSV Dateien über Config einlesen (Wird als string eingelesen, muss aber Rune sein.)
- [ ] Support für Intos Artikellisten einbauen
- [-] Bilder hochladen
- [-] Bilder zu Produkten hinzufügen
- [-] Produkt Cover setzen
- [-] "EOL Artikel" im Shop löschen
