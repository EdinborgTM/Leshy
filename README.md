# Leshy - Scanner Port

Leshy yek scanner port sari, ghavi va pishrafte ba zaban Go hast ke baraye skan kardan port-haye TCP va UDP estefade mishe. In barname baraye Termux va system-haye dige beine shode, az nokh-haye hamzaman baraye sorat besiar bala estefade mikone, va natayej ro dar JSON zakhire mikone. Mitune vaziyat port-ha (baz, baste, ya baz|filtered), baner, va noskhe service-ha ro ba payload-haye pishrafte (dns, version.bind, ntp, snmp, test) neshun bede.

## Noskhe
- Noskhe jari: 0.2.16
- Tarikh enteshar: 03 Oct 2025
- Taghirat:
  - Ezafe kardan flag -w baraye timeout UDP puya va log-haye debug bishtar dar girBaner (noskhe 0.2.16)
  - Ezafe kardan decodeDNSResponse baraye estekhraj TXT record az pasokh-haye DNS (noskhe 0.2.15)
  - Behbude estekhraj noskhe baraye version.bind va eslah payload dns (noskhe 0.2.15)
  - Ezafe kardan decode ASCII baraye baner-ha va payload version.bind baraye DNS (noskhe 0.2.14)
  - Afzayesh timeout UDP va TCP be 3 sanie (noskhe 0.2.14)
  - Refe kheta-ye compile 'start declared and not used' dar girBaner (noskhe 0.2.13)
  - Refe moshkel zaman_ms (0 budan) ba mohasebe dorost dar kargar (noskhe 0.2.12)
  - Afzayesh timeout TCP be 2 sanie va ezafe kardan log debug baraye girBaner (noskhe 0.2.12)
  - Ezafe kardan payload 'test' sadeh baraye UDP (noskhe 0.2.12)
  - Behbude tashkhis vaziyat port-haye UDP ba tasnifKheta va afzayesh timeout be 2 sanie baraye baner (noskhe 0.2.11)
  - Refe moshkel khali budan baner va noskhe dar skan UDP ba afzayesh timeout va tasnif vaziyat (noskhe 0.2.11)
  - Refe kheta-ye compile 'undefined: rand.Uint16' ba estefade az math/rand be jaye math/rand/v2 (noskhe 0.2.10)
  - Refe kheta-ye compile 'undefined: tasnifKhta' dar leshy.go tavasot eslah be tasnifKheta (noskhe 0.2.8)
  - Refe kheta-ye compile 'undefined: rand' va 'undefined: rand.Uint16' ba import math/rand/v2 (noskhe 0.2.5)
  - Refe kheta-ye compile 'undefined: GirPayloadUDP' ba 'go build -o leshy *.go' (noskhe 0.2.3)
  - Ezafe shodan file payloads.go baraye modularity payload-haye UDP
  - Behbude vaziyat UDP ba 'baz|filtered' baraye timeout-ha
  - Refe kheta 'ctxTimeout declared and not used' dar skan UDP (noskhe 0.2.2)
  - Refe kheta-haye context dar skan TCP/UDP (noskhe 0.2.1)
  - Behbude modiriyat kheta-ha va baste shodan etesal-ha
  - Ezafe shodan timeout pishfarz baraye UDP
  - Behbude tasadofi-sazi port-ha ba -f
  - Hadaksar 100 nokh baraye system-haye kam-masraf
  - Payam-haye kheta-ye karbarpasand-tar

## Vizhegi-ha
- Skan TCP va UDP ba sorat besiar bala (bedun timeout dar halat adi)
- Nokh-haye hamzaman ba tedad auto (hadaksar 100)
- Khoruji JSON ba noskhe service-ha
- Khandan baner va noskhe (optional)
- Payload-haye pishrafte DNS/VERSION.BIND/NTP/SNMP/TEST baraye UDP dar file payloads.go
- Rangi ba ANSI codes baraye khunayi
- Bedun vabastegi be ketabkhune-haye khareji
- Beinesazi baraye Termux ba -l
- Makhfi az firewall ba -f
- Ghat scan ba signal
- Timeout puya baraye UDP ba -w

## Nasb
Bayad Go (golang) rooye system nasb bashad (version 1.18 ya balatar baraye math/rand). Dar Termux, Go ro ba in dastur nasb konid:

```bash
pkg install golang
```

1. Code ro clone konid:
   ```bash
   git clone https://github.com/EdinborgTM/Leshy.git
   cd Leshy
   ```

2. Compile konid (hame file-haye .go ro shamel mishavad):
   ```bash
   go build -o leshy *.go
   ```

3. Dar Termux, dastresi be /sdcard bedid:
   ```bash
   termux-setup-storage
   ```

## Estefade
Baraye ejra:
```bash
./leshy -t <IP ya hostname> [guzineha]
```

Baraye didan noskhe barname:
```bash
./leshy -V
```

### Guzineha
| Guzine | Tozih | Pishfarz |
|--------|-------|----------|
| `-t` | Hadaf (IP ya hostname) - lazem | - |
| `-m` | Kamtarin port | 1 |
| `-x` | Bishtarin port | 1024 |
| `-r` | Tedad nokh (0 = auto) | 0 |
| `-o` | Timeout (ms, faghat ba -l) | 1000 |
| `-w` | Timeout UDP (ms, baraye baner) | 3000 |
| `-b` | Khandan baner va noskhe | false |
| `-v` | Khoruji mofasal | false |
| `-l` | Kam masraf baraye termux | false |
| `-f` | Makhfi az firewall | false |
| `-p` | Protokol (tcp ya udp) | tcp |
| `-y` | Payload baraye udp (dns, version.bind, ntp, snmp, test, ya none) | none |
| `-u` | File JSON | /sdcard/leshy_scan.json |
| `-V` | Neshan dadan noskhe barname | false |

### Mesal-ha
1. Skan TCP sari port-haye 1 ta 100:
   ```bash
   ./leshy -t 192.168.1.10 -m 1 -x 100
   ```

2. Skan UDP ba payload version.bind va noskhe:
   ```bash
   ./leshy -t 8.8.8.8 -m 53 -x 53 -p udp -y version.bind -b
   ```

3. Skan UDP ba timeout puya:
   ```bash
   ./leshy -t pool.ntp.org -m 123 -x 123 -p udp -y ntp -b -w 5000
   ```

4. Skan kam masraf ba makhfi baraye Termux:
   ```bash
   ./leshy -t 127.0.0.1 -m 1 -x 100 -l -f -b -p udp -y test
   ```

5. Didan noskhe:
   ```bash
   ./leshy -V
   ```

### Nemune khoruji
```bash
>> Scan 8.8.8.8 (udp)
>> 53/udp: baz
Baner: Google Public DNS
Noskhe: Google Public DNS
>> Tamam!
>> 8.8.8.8 (udp)
>> Baz: 1
>> Zaman: 120 ms
>> File: ./r.json
```

## File khoruji
Natayej dar JSON zakhire mishan (pishfarz: `/sdcard/leshy_scan.json`):
```json
{
  "hadaf": "8.8.8.8",
  "shuru": "2025-10-03T15:35:00Z",
  "tamam": "2025-10-03T15:35:05Z",
  "zaman_ms": 120,
  "port_ha": [
    {
      "port": 53,
      "protokol": "udp",
      "vaziyat": "baz",
      "baner": "Google Public DNS",
      "noskhe": "Google Public DNS",
      "zaman_ms": 120
    }
  ]
}
```

## Ehtiyat
- **Etebar**: Faghat rooye system-haye khodetun skan konid. Skan bi ejaze ghanooni nist.
- **Termux**: Az `-l` baraye kam kardan masraf estefade konid. Bedun `-l`, skan ba maximum sorat ejra mishe.
- **Makhfi**: Ba `-f`, skan makhfi-tar mishe (port-haye tasadofi + dirang), ama sorat kamtar mishe.
- **UDP**: Payload-haye dns/version.bind/ntp/snmp/test baraye UDP noskhe service ro behtar estekhraj mikonand. Timeout baraye UDP va TCP be 3 sanie afzayesh yaft.
- **Dastresi**: Dar Termux, `/sdcard` bayad dastresi dashte bashe:
  ```bash
  termux-setup-storage
  ```
- **Compile**: Hatman az `go build -o leshy *.go` estefade konid ta hame file-haye .go (shamel leshy.go va payloads.go) compile shavand.
- **Go Version**: Az Go 1.18 ya balatar estefade konid ta az math/rand sazgari dashte bashad.

## License
In project tahte ejaze-name [MIT License](LICENSE) montasher shode ast.

## Tamase ba ma
Sual ya pishnahad? Dar GitHub issue besazid ya ba <your-email> tamas begirid.
