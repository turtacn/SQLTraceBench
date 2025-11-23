package reports

const HTMLTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SQLTraceBench Report</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 1000px; margin: 0 auto; padding: 20px; }
        h1, h2 { color: #2c3e50; border-bottom: 2px solid #eee; padding-bottom: 10px; }
        .summary-card { background: #f8f9fa; border: 1px solid #ddd; padding: 20px; border-radius: 8px; margin-bottom: 20px; display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; }
        .metric-box { text-align: center; }
        .metric-val { font-size: 24px; font-weight: bold; color: #007bff; }
        .metric-label { color: #666; font-size: 14px; }
        .pass { color: #28a745; font-weight: bold; }
        .fail { color: #dc3545; font-weight: bold; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background-color: #f2f2f2; }
        tr:nth-child(even) { background-color: #f9f9f9; }
        .chart-container { position: relative; height: 300px; width: 100%; margin-bottom: 40px; }
    </style>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body>
    <h1>SQLTraceBench Validation Report</h1>

    <div class="summary-card">
        <div class="metric-box">
            <div class="metric-val" style="color: {{if .Result.Pass}}#28a745{{else}}#dc3545{{end}}">{{if .Result.Pass}}PASS{{else}}FAIL{{end}}</div>
            <div class="metric-label">Status</div>
        </div>
        <div class="metric-box">
            <div class="metric-val">{{printf "%.2f" .Result.BaseMetrics.QPS}}</div>
            <div class="metric-label">Base QPS</div>
        </div>
        <div class="metric-box">
            <div class="metric-val">{{printf "%.2f" .Result.CandidateMetrics.QPS}}</div>
            <div class="metric-label">Candidate QPS</div>
        </div>
        <div class="metric-box">
            <div class="metric-val">{{printf "%.2f" .Metadata.Threshold}}%</div>
            <div class="metric-label">Threshold</div>
        </div>
    </div>

    <p>{{.Result.Reason}}</p>

    <h2>Latency Distribution</h2>
    <div class="chart-container">
        <canvas id="latencyChart"></canvas>
    </div>

    <h2>Detailed Metrics</h2>
    <table>
        <thead>
            <tr>
                <th>Metric</th>
                <th>Base</th>
                <th>Candidate</th>
                <th>Delta</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td>Total Queries</td>
                <td>{{.Result.BaseMetrics.QueriesExecuted}}</td>
                <td>{{.Result.CandidateMetrics.QueriesExecuted}}</td>
                <td>{{sub .Result.CandidateMetrics.QueriesExecuted .Result.BaseMetrics.QueriesExecuted}}</td>
            </tr>
            <tr>
                <td>QPS</td>
                <td>{{printf "%.2f" .Result.BaseMetrics.QPS}}</td>
                <td>{{printf "%.2f" .Result.CandidateMetrics.QPS}}</td>
                <td>{{printf "%.2f" (subf .Result.CandidateMetrics.QPS .Result.BaseMetrics.QPS)}}</td>
            </tr>
            <tr>
                <td>P50 Latency</td>
                <td>{{.Result.BaseMetrics.P50}}</td>
                <td>{{.Result.CandidateMetrics.P50}}</td>
                <td>-</td>
            </tr>
            <tr>
                <td>P90 Latency</td>
                <td>{{.Result.BaseMetrics.P90}}</td>
                <td>{{.Result.CandidateMetrics.P90}}</td>
                <td>-</td>
            </tr>
            <tr>
                <td>P99 Latency</td>
                <td>{{.Result.BaseMetrics.P99}}</td>
                <td>{{.Result.CandidateMetrics.P99}}</td>
                <td>-</td>
            </tr>
            <tr>
                <td>Errors</td>
                <td>{{.Result.BaseMetrics.Errors}}</td>
                <td>{{.Result.CandidateMetrics.Errors}}</td>
                <td>{{sub .Result.CandidateMetrics.Errors .Result.BaseMetrics.Errors}}</td>
            </tr>
        </tbody>
    </table>

    <script>
        const ctx = document.getElementById('latencyChart').getContext('2d');
        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: ['P50', 'P90', 'P99'],
                datasets: [{
                    label: 'Base',
                    data: [
                        {{.Result.BaseMetrics.P50.Milliseconds}},
                        {{.Result.BaseMetrics.P90.Milliseconds}},
                        {{.Result.BaseMetrics.P99.Milliseconds}}
                    ],
                    backgroundColor: 'rgba(54, 162, 235, 0.5)',
                    borderColor: 'rgba(54, 162, 235, 1)',
                    borderWidth: 1
                }, {
                    label: 'Candidate',
                    data: [
                        {{.Result.CandidateMetrics.P50.Milliseconds}},
                        {{.Result.CandidateMetrics.P90.Milliseconds}},
                        {{.Result.CandidateMetrics.P99.Milliseconds}}
                    ],
                    backgroundColor: 'rgba(255, 99, 132, 0.5)',
                    borderColor: 'rgba(255, 99, 132, 1)',
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'Latency (ms)'
                        }
                    }
                }
            }
        });
    </script>
</body>
</html>
`
