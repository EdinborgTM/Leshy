# Leshy - Port Scanner

Leshy yek scanner port sari, ghavi va sadeh ba zaban Go hast ke baraye skan kardan port-haye TCP estefade mishe. In barname baraye Termux va system-haye dige beine shode, az nokh-haye hamzaman baraye sorat besiar bala estefade mikone, va natayej ro dar JSON zakhire mikone. Mitune vaziyat port-ha (baz, baste, ya filtered), banner, va version service-ha ro neshun bede.

## Vizhegi-ha
- Skan TCP ba sorat besiar bala (bedun timeout dar halat adi)
- Multi-threading ba nokh-haye auto
- Khoruji JSON ba version service-ha
- Khandan banner va version (optional)
- Rangi ba ANSI codes baraye khunayi
- Bedun vabastegi be ketabkhune-haye khareji
- Beinesazi baraye Termux ba -l
- Evasion mode ba -f baraye makhfi az firewall
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
./leshy -t <IP ya hostname> [options]
```

### Parametrha
| Parametr | Tozih | Pishfarz |
|----------|-------|----------|
| `-t` | Hadaf (IP ya hostname) - lazem | - |
| `-m` | Kamtarin port | 1 |
| `-x` | Bishtarin port | 1024 |
| `-r` | Tedad nokh (0 = auto) | 0 |
| `-o` | Timeout (ms, faghat ba -l) | 1000 |
| `-b` | Khandan banner va version | false |
| `-v` | Khoruji mofasal | false |
| `-l` | Kam masraf baraye Termux | false |
| `-f` | Evasion mode (makhfi az firewall) | false |
| `-u` | File JSON | /sdcard/leshy_scan.json |

### Mesal-ha
1. Skan sari port-haye 1 ta 100:
   ```bash
   ./leshy -t 192.168.1.10 -m 1 -x 100
   ```

2. Skan ba banner, version, va evasion:
   ```bash
   ./leshy -t example.com -m 1 -x 1000 -b -f
   ```

3. Skan kam masraf baraye Termux ba version:
   ```bash
   ./leshy -t 127.0.0.1 -m 1 -x 100 -l -b
   ```

### Nemune khoruji
```bash
>> Scan 192.168.1.10
>> Evasion mode: randomized ports + delay
>> 50%
>> 22: open
Banner: SSH-2.0-OpenSSH_8.0
Version: SSH-2.0-OpenSSH_8.0
>> 80: open
Banner: Apache/2.4.52 (Ubuntu)
Version: Apache/2.4.52 (Ubuntu)
>> Tamam!
>> 192.168.1.10
>> Baz: 2
>> Zaman: 1234 ms
>> File: /sdcard/leshy_scan.json
```

## File khoruji
Natayej dar JSON zakhire mishan (pishfarz: `/sdcard/leshy_scan.json`):
```json
{
  "target": "192.168.1.10",
  "started": "2025-09-29T12:37:00Z",
  "finished": "2025-09-29T12:37:05Z",
  "elapsed_ms": 5000,
  "ports": [
    {
      "port": 22,
      "state": "open",
      "banner": "SSH-2.0-OpenSSH_8.0",
      "version": "SSH-2.0-OpenSSH_8.0",
      "elapsed_ms": 120
    },
    {
      "port": 80,
      "state": "open",
      "banner": "Apache/2.4.52 (Ubuntu)",
      "version": "Apache/2.4.52 (Ubuntu)",
      "elapsed_ms": 100
    },
    {
      "port": 81,
      "state": "closed",
      "elapsed_ms": 90
    }
  ]
}
```

## Ehtiyat
- **Etebar**: Faghat rooye system-haye khodetun skan konid. Skan bi ejaze ghanooni nist.
- **Termux**: Az `-l` baraye kam kardan masraf estefade konid. Bedun `-l`, skan ba maximum sorat ejra mishe.
- **Evasion**: Ba `-f`, skan makhfi-tar mishe (randomized ports + delay), ama sorat kamtar mishe.
- **Version**: Baraye version service-ha, `-b` estefade konid (masalan Apache ya OpenSSH).
- **Dastresi**: Dar Termux, `/sdcard` bayad dastresi dashte bashe:
  ```bash
  termux-setup-storage
  ```

## License
In project tahte ejaze-name [MIT License](LICENSE) montasher shode ast.

## Tamase ba ma
Sual ya pishnahad? Dar GitHub issue besazid ya ba tamas begirid.
