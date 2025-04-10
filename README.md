# 🔁 NATS-to-Prometheus Exporter Bridge

A lightweight Go service that reads system metrics from **NATS JetStream** and exposes them at a Prometheus-compatible `/metrics` endpoint.  
Designed to work with agent-based systems like [Logs Exporter](https://github.com/gysosin/Logs_exporter).

---

## 🚀 Features

- 📨 Subscribes to NATS JetStream (`metrics` subject)
- 📊 Serves all received metrics at `/metrics` (Prometheus format)
- 🧠 Injects `system_name` label to support multi-agent scraping
- 🧹 TTL-based cache cleanup to avoid memory bloat
- ⚙️ Configurable via `config.json`
- 🐳 Docker & systemd ready
- ❤️ Built with simplicity, performance & scale in mind

---

## 📦 Example Output

```
windows_cpu_usage_percent{system_name="agent-A"} 12.3
windows_memory_bytes{type="used", system_name="agent-B"} 87654321
```

---

## ⚙️ Configuration

Edit `config.json`:

```json
{
  "listen_port": "2112",
  "nats_url": "nats://localhost:4222",
  "subject": "metrics"
}
```

---

## 🛠 Build & Run

### 🔨 Build Locally

```bash
go build -o exporter ./cmd/exporter
./exporter
```

---

### 🐳 Run with Docker

```bash
docker build -t nats-prom-exporter .
docker run -d -p 2112:2112 \
  -v $(pwd)/config.json:/etc/exporter/config.json \
  --name exporter nats-prom-exporter
```

---

### 🛡️ Run as Systemd Service (Linux)

```bash
sudo cp exporter /usr/local/bin/exporter
sudo cp config.json /etc/exporter/config.json
sudo cp exporter.service /etc/systemd/system/
sudo systemctl daemon-reexec
sudo systemctl enable exporter
sudo systemctl start exporter
```

---

## 📈 Prometheus Integration

Add to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: "nats_exporter"
    static_configs:
      - targets: ["localhost:2112"]
```

Then query with:

```promql
windows_memory_bytes{system_name="agent-A"}
avg(windows_cpu_usage_percent) by (system_name)
```

---

## ✅ Roadmap

- [ ] `/exporter_status` endpoint for bridge diagnostics
- [ ] Built-in Prometheus metrics for the exporter itself
- [ ] Support for NATS TLS & JWT authentication
- [ ] Per-agent scrape stats and last seen info
- [ ] Optional persistent cache layer (BoltDB, Redis)

---

## 💡 Created By

Made with ❤️ by **Abhishek Thakur**  
Go build cool things. Monitor them smarter. 🚀

---

## 📄 License

MIT — free to use, fork, improve, and share.

```
