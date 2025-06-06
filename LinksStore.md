# Links Store - Хранилище ссылок

Сервис хранит RSS-ссылки в MySQL и отдаёт их по HTTP.  
Адрес по умолчанию — **http://localhost:8080**

| Метод  | URL                          | Описание              |
|--------|-----------------------------|-----------------------|
| GET    | `/api/v1/links`             | список всех лент      |
| GET    | `/api/v1/links/{id}`        | одна лента по ID      |
| POST   | `/api/v1/links`             | добавить новую        |
| DELETE | `/api/v1/links/{id}`        | удалить ленту (опц.)  |

---

## Добавить ссылку

### Linux / macOS

```bash
curl -X POST http://localhost:8080/api/v1/links \
     -H "Content-Type: application/json" \
     -d '{"url":"ссылка","label":"название"}'
```


### Windows (PowerShell)

```bash
$hdr = @{ "Content-Type" = "application/json" }
Invoke-RestMethod -Uri http://localhost:8080/api/v1/links `
                  -Method POST -Headers $hdr `
                  -Body '{"url":"ссылка","label":"название"}'
```
