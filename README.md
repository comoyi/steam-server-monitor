# steam-server-monitor

## Dependence
```
go get fyne.io/fyne/v2
go get fyne.io/fyne/v2/cmd/fyne

#go install fyne.io/fyne/v2/cmd/fyne@v2.2.3
```

## Usage

1.Edit Makefile replace <font-file> with a font file that supports Chinese then execute
```
#example
bundle-font:
	fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go ~/.local/share/fonts/HarmonyOS_Sans_SC_Regular.ttf
```

```
make bundle-font
```

2.Make and run

```
make
```
