# Leshy - Port Scanner

Leshy yek scanner port sadeh, ghavi va sari ba zaban Go hast ke baraye skan kardan port-haye TCP estefade mishe. In barname az nokh-haye hamzaman (multi-threading) baraye afzayesh sorat estefade mikone va natayej ro be surat JSON zakhire mikone. Ba estefade az in barname mitunid vaziyat port-ha (baz, baste, ya filtered) ro check konid va dar surat niaz banner service-ha ro bekhunid.

## Vizhegi-ha
- Skan kardan port-haye TCP ba sorat bala
- Poshtibani az multi-threading baraye afzayesh efficiency
- Khoruji JSON baraye tahlil asun
- Poshtibani az khandan banner (optional)
- Rangi kardan khoruji baraye khunayi behtar (ba ANSI codes)
- Bedun vabastegi be ketabkhune-haye khareji
- Poshtibani az signal handling baraye ghat scan

## Nasb
Baraye estefade az Leshy, bayad Go (golang) rooye system nasb bashad (version 1.16 ya balatar). Hich vabastegi khareji nadare, pas faghat kafi ast code ro download konid va compile konid.

1. Code ro clone ya download konid:
   ```bash
   git clone https://github.com/EdinborgTM/Leshy.git
   cd Leshy
   ```

2. Barname ro compile konid:
   ```bash
   go build -o leshy leshy.go
   ```

## Estefade
Baraye ejra, az dastur zir estefade konid:
```bash
./leshy --target <IP ya hostname> [options]
```

### Parametrha
| Parametr | Tozih | Pishfarz |
|----------|-------|----------|
| `--target` | Hadaf (IP ya hostname) - lazem | - |
| `--min` | Kamtarin port baraye skan | 1 |
| `--max` | Bishtarin port baraye skan | 1024 |
| `--threads` | Tedad nokh-haye hamzaman | 100 |
| `--timeout` | Timeout etesal (milli sanie) | 800 |
| `--banner` | Khandan banner bad az etesal | false |
| `--verbose` | Nashun dadan khoruji mofasal (hamin port-ha) | false |
| `--out` | File khoruji JSON | scans/leshy_scan.json |

### Mesal-ha
1. Skan kardan port-haye 1 ta 100 rooye 192.168.1.10:
   ```bash
   ./leshy --target 192.168.1.10 --min 1 --max 100
   ```

2. Skan ba banner va khoruji mofasal:
   ```bash
   ./leshy --target example.com --min 1 --max 1000 --banner --verbose
   ```

3. Skan ba timeout bishtar va thread kamtar:
   ```bash
   ./leshy --target 127.0.0.1 --min 1 --max 65535 --threads 50 --timeout 1000
   ```

### Nemune khoruji
```bash
>> Shuru scan rooye 192.168.1.10 (192.168.1.10) [port-ha: 1-100]
>> 50/100 port scan shod (50%)
>> Port 22: open (banner: SSH-2.0-OpenSSH_8.0)
>> Port 80: open
>> Scan tamam shod!
>> Hadaf: 192.168.1.10 (192.168.1.10)
>> Port-haye scan shode: 100
>> Port-haye baz: 2
>> Zaman kol: 1234 milli sanie
>> File khoruji: scans/leshy_scan.json
```

## File khoruji
Natayej scan dar yek file JSON zakhire mishan (pishfarz: `scans/leshy_scan.json`). Sakhtar file JSON be in surat ast:
```json
{
  "target": "192.168.1.10",
  "started": "2025-09-29T12:37:00Z",
  "finished": "2025-09-29T12:37:05Z",
  "elapsed_ms": 5000,
  "ports": [
    {"port": 22, "state": "open", "banner": "SSH-2.0-OpenSSH_8.0", "elapsed_ms": 120},
    {"port": 80, "state": "open", "elapsed_ms": 100},
    {"port": 81, "state": "closed", "elapsed_ms": 90}
  ]
}
```

## Ehtiyat
- **Etebar**: Faghat rooye hadaf-hayi ke ejaze skan daran (masalan system-haye khodetun) estefade konid. Skan kardan bi ejaze ghanooni nist.
- **Timeout**: Agar timeout khili kam bashe, natayej gheyre daghigh mishan. Pishnahad: 500-1000 milli sanie.
- **Threads**: Tedad thread-haye bishtar sorat ro afzayesh mide, ama mitune be system ya shabake feshar biyare.

## License
In project tahte ejaze-name [MIT License](LICENSE) montasher shode ast.

## Tamase ba ma
Agar sual ya pishnahadi darid, mitunid dar GitHub issue besazid ya ba email <your-email> tamas begirid.
