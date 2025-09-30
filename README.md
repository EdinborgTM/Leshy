# Leshy - Scanner Port

Leshy yek scanner port sari, ghavi va pishrafte ba zaban Go hast ke baraye skan kardan port-haye TCP va UDP estefade mishe. In barname baraye Termux va system-haye dige beine shode, az nokh-haye hamzaman baraye sorat besiar bala estefade mikone, va natayej ro dar JSON zakhire mikone. Mitune vaziyat port-ha (baz, baste, ya filtered), baner, va noskhe service-ha ro ba payload-haye ghavi (dns, ntp, snmp) neshun bede.

## Vizhegi-ha
- Skan TCP va UDP ba sorat besiar bala (bedun timeout dar halat adi)
- Nokh-haye hamzaman ba tedad auto
- Khoruji JSON ba noskhe service-ha
- Khandan baner va noskhe (optional)
- Payload-haye pishrafte DNS/NTP/SNMP baraye UDP
- Rangi ba ANSI codes baraye khunayi
- Bedun vabastegi be ketabkhune-haye khareji
- Beinesazi baraye Termux ba -l
- Makhfi az firewall ba -f
- Ghat scan ba signal

## Nasb
Bayad Go (golang) rooye system nasb bashad (version 1.16 ya balatar). Dar Termux, Go ro ba in dastur nasb konid:

```bash
pkg install golang
```

1. Code ro clone konid:
   ```bash
   git clone https://github.com/EdinborgTM/Leshy.git
   cd Leshy
   ```

2. Compile konid:
   ```bash
   go build -o leshy leshy.go
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
| `-l` | Kam masraf baraye Termux | false |
| `-f` | Makhfi az firewall | false |
| `-p` | Protokol (tcp ya udp) | tcp |
| `-y` | Payload baraye udp (dns, ntp, snmp, ya none) | none |
| `-u` | File JSON | /sdcard/leshy_scan.json |

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
  "shuru": "2025-09-29T12:37:00Z",
  "tamam": "2025-09-29T12:37:05Z",
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
      "vaziyat": "baste",
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

## License
In project tahte ejaze-name [MIT License](LICENSE) montasher shode ast.
