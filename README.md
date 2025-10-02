# Leshy - Scanner Port

Leshy yek scanner port sari, ghavi va pishrafte ba zaban Go hast ke baraye skan kardan port-haye TCP va UDP estefade mishe. In barname baraye Termux va system-haye dige beine shode, az nokh-haye hamzaman baraye sorat besiar bala estefade mikone, va natayej ro dar JSON zakhire mikone. Mitune vaziyat port-ha (baz, baste, ya baz|filtered), baner, va noskhe service-ha ro ba payload-haye pishrafte (dns, ntp, snmp) neshun bede.

## Noskhe
- Noskhe jari: 0.2.8
- Tarikh enteshar: 02 Oct 2025
- Taghirat:
  - Refe kheta-ye compile 'undefined: tasnifKhta' dar leshy.go tavasot eslah be tasnifKheta va ta'kid bar import math/rand/v2 baraye rand.Uint16 (noskhe 0.2.8)
  - Refe kheta-ye compile 'undefined: rand.Uint16' ba ta'kid bar dorost budan import math/rand/v2 va check encoding file-haye dar WSL (noskhe 0.2.7)
  - Refe kheta-ye compile 'undefined: rand' va 'undefined: rand.Uint16' ba import math/rand/v2 dar leshy.go va payloads.go (noskhe 0.2.5)
  - Refe kheta-ye compile 'undefined: GirPayloadUDP' ba estefade az 'go build -o leshy *.go' (noskhe 0.2.3)
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
- Payload-haye pishrafte DNS/NTP/SNMP baraye UDP dar file payloads.go
- Rangi ba ANSI codes baraye khunayi
- Bedun vabastegi be ketabkhune-haye khareji
- Beinesazi baraye Termux ba -l
- Makhfi az firewall ba -f
- Ghat scan ba signal

## Nasb
Bayad Go (golang) rooye system nasb bashad (version 1.22 ya balatar baraye math/rand/v2). Dar Termux, Go ro ba in dastur nasb konid:

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
| `-b` | Khandan baner va noskhe | false |
| `-v` | Khoruji mofasal | false |
| `-l` | Kam masraf baraye termux | false |
| `-f` | Makhfi az firewall | false |
| `-p` | Protokol (tcp ya udp) | tcp |
| `-y` | Payload baraye udp (dns, ntp, snmp, ya none) | none |
| `-u` | File JSON | /sdcard/leshy_scan.json |
| `-V` | Neshan dadan noskhe barname | false |

### Mesal-ha
1. Skan TCP sari port-haye 1 ta 100:
   ```bash
   ./leshy -t 192.168.1.10 -m 1 -x 100
   ```

2. Skan UDP ba payload DNS va noskhe:
   ```bash
   ./leshy -t example.com -m 1 -x 1000 -p udp -y dns -b
   ```

3. Skan kam masraf ba makhfi baraye Termux:
   ```bash
   ./leshy -t 127.0.0.1 -m 1 -x 100 -l -f -b -p udp -y ntp
   ```

4. Didan noskhe:
   ```bash
   ./leshy -V
   ```

### Nemune khoruji
```bash
>> Scan 192.168.1.10 (udp)
>> Makhfi: port-haye tasadofi + dirang
>> 50%
>> 53/udp: baz
Baner: 9.11.4-3ubuntu5.1-Ubuntu
Noskhe: BIND 9.11.4
>> 123/udp: baz
Baner: ntpd 4.2.8p15@1.3728-o Wed Sep 23 14:46:23 UTC 2020
Noskhe: ntpd 4.2.8p15
>> Tamam!
>> 192.168.1.10 (udp)
>> Baz: 2
>> Zaman: 1234 ms
>> File: /sdcard/leshy_scan.json
```

## File khoruji
Natayej dar JSON zakhire mishan (pishfarz: `/sdcard/leshy_scan.json`):
```json
{
  "hadaf": "192.168.1.10",
  "shuru": "2025-10-02T15:35:00Z",
  "tamam": "2025-10-02T15:35:05Z",
  "zaman_ms": 5000,
  "port_ha": [
    {
      "port": 53,
      "protokol": "udp",
      "vaziyat": "baz",
      "baner": "9.11.4-3ubuntu5.1-Ubuntu",
      "noskhe": "BIND 9.11.4",
      "zaman_ms": 120
    },
    {
      "port": 123,
      "protokol": "udp",
      "vaziyat": "baz",
      "baner": "ntpd 4.2.8p15@1.3728-o Wed Sep 23 14:46:23 UTC 2020",
      "noskhe": "ntpd 4.2.8p15",
      "zaman_ms": 100
    },
    {
      "port": 161,
      "protokol": "udp",
      "vaziyat": "baz|filtered",
      "zaman_ms": 90
    }
  ]
}
```

## Ehtiyat
- **Etebar**: Faghat rooye system-haye khodetun skan konid. Skan bi ejaze ghanooni nist.
- **Termux**: Az `-l` baraye kam kardan masraf estefade konid. Bedun `-l`, skan ba maximum sorat ejra mishe.
- **Makhfi**: Ba `-f`, skan makhfi-tar mishe (port-haye tasadofi + dirang), ama sorat kamtar mishe.
- **UDP**: Payload-haye dns/ntp/snmp baraye UDP noskhe service ro behtar estekhraj mikonand, ama mitunan zaman bishtari begirand.
- **Dastresi**: Dar Termux, `/sdcard` bayad dastresi dashte bashe:
  ```bash
  termux-setup-storage
  ```
- **Compile**: Hatman az `go build -o leshy *.go` estefade konid ta hame file-haye .go (shamel leshy.go va payloads.go) compile shavand.
- **Go Version**: Az Go 1.22 ya balatar estefade konid ta az math/rand/v2 sazgari dashte bashad. File payloads.go bayad `import "math/rand/v2"` dashte bashad. Agar kheta-ye `undefined: rand.Uint16` ya `undefined: tasnifKhta` didid, file-haye payloads.go va leshy.go ro check konid va az dorost budan import math/rand/v2 va eslah tasnifKheta motmaen shavid.

## License
In project tahte ejaze-name [MIT License](LICENSE) montasher shode ast.

## Tamase ba ma
Sual ya pishnahad? Dar GitHub issue besazid ya ba <your-email> tamas begirid.
