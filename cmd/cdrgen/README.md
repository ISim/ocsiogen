# cdrgen

Generátor fiktivních hovorů (CDR) ve formátu YATE ústředny (
https://docs.yate.ro/wiki/CDR_File_Module)


## Generování

```
go run ./cmd/cdrgen/ -start-at 2019-06-20 -finish-at 2019-06-24 -records-per-day 10
```

Dalším parametrem je `-msisdn-file <soubor>`, který obsahuje seznam A-čísel pro generování
např.

```
420499816415;2019-04-01T00:00:00Z;
420775095760;2019-04-01T00:00:00Z;2019-05-01T00:00:00Z
420318621441;2019-04-01T00:00:00Z;
420737769653;2019-04-01T00:00:00Z;
420737109631;2019-04-01T00:00:00Z;
```

Generátor pracuje s tabulkou směrů (destinací) uloženou staticky v `prefixes.go` a generuje
národní, mezinárodní hovory + hovory na speciální čísla v rozložení cca

* 5% speciální čísla
* 25% mezinárodní hovory
* 70% národní hovory


Jednotlivé hovory z konkrétného A-čísla se nepřekrývajív jeden okamžik může číslo mít pouze
jeden běžící hovor (pro generování většího počtu hovorů za den je tedy třeba mít dostatečně
mohutnou množinu čísel v souboru).

Pokud je generovaný soubor pro víkend, parametr `-records-per-day` se vydělí deseti.

Časy hovorů se generují náhodně s následujícím rozdělením:

* 1% - trvání 10minut - 3 hodiny
* 39% - trvání 3-9 minut
* 60% - trvání 0 - 2 minuty s normálním rozdělením


Výsledné soubory jsou generovány po hodinách s procentuálně nejvyšším zastoupením
v 9. a 14. hodině.