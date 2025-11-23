# Common Issues & Solutions

## 1. Installation Issues

### Issue: Binary not found after installation
**Symptom**: `sql_trace_bench: command not found`

**Solution**:
```bash
# Ensure Go bin is in PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Or reinstall
go install github.com/yourusername/sql-trace-bench@latest
```

---

## 2. Generation Issues

### Issue: Out of memory during generation
**Symptom**: `fatal error: out of memory`

**Solution**:
1. Reduce batch size in config:
   ```yaml
   generation:
     batch_size: 100  # Reduce from 1000
   ```
2. Increase Docker memory limit:
   ```yaml
   # docker-compose.yml
   services:
     app:
       mem_limit: 4g  # Increase memory
   ```

---

## 3. Validation Issues

### Issue: K-S test always fails
**Symptom**: `Kolmogorov-Smirnov test failed: D=0.15`

**Possible Causes**:
*   Generated data distribution differs significantly from original.
*   Sample size too small.

**Solution**:
1. Check model training:
   ```bash
   sql_trace_bench train --validate
   ```
2. Adjust significance level in `configs/validation.yaml`:
   ```yaml
   validation:
     ks_test:
       alpha: 0.10  # More lenient threshold
   ```

---

## 4. Database Issues

### Issue: Connection refused to PostgreSQL
**Symptom**: `dial tcp 127.0.0.1:5432: connect: connection refused`

**Solution**:
```bash
# Check if PostgreSQL is running
docker-compose ps

# Restart PostgreSQL
docker-compose restart postgres

# Check logs
docker-compose logs postgres
```

---

## 5. Performance Issues

### Issue: Generation too slow
**Symptom**: <10 traces/sec throughput

**Solution**:
1. Enable parallel generation:
   ```yaml
   generation:
     concurrency: 10  # Adjust based on CPU cores
   ```
2. Use faster model:
   ```bash
   sql_trace_bench generate -m markov  # Instead of LSTM
   ```

---

## 6. API Issues

### Issue: 500 Internal Server Error
**Symptom**: API returns `{"error": "internal server error"}`

**Solution**:
1. Check logs:
   ```bash
   docker-compose logs app | grep ERROR
   ```
2. Enable debug mode:
   ```yaml
   # config.yaml
   server:
     log_level: debug
   ```

---

## 7. Metrics Issues

### Issue: Metrics not showing in Grafana
**Symptom**: Empty charts in dashboard

**Solution**:
1. Verify Prometheus is scraping:
   ```bash
   curl http://localhost:9090/api/v1/targets
   ```
2. Check metrics endpoint:
   ```bash
   curl http://localhost:9091/metrics | grep benchmark
   ```
3. Restart Grafana:
   ```bash
   docker-compose restart grafana
   ```

---

## 8. Docker Issues

### Issue: Port already in use
**Symptom**: `bind: address already in use`

**Solution**:
```bash
# Find process using port
lsof -i :8080

# Kill process or change port in docker-compose.yml
ports:
  - "8081:8080"  # Use different host port
```

---

## Getting Help

If your issue is not listed:
1. Check [GitHub Issues](https://github.com/yourusername/sql-trace-bench/issues)
2. Join our Discord community.
3. Contact support: [support@example.com](mailto:support@example.com)
