// Configuration
const API_BASE = window.location.origin;
const UPDATE_INTERVAL = 30000; // 30 seconds
let currentTimeRange = 300; // 5 minutes default
let chart = null;

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    initChart();
    setupEventListeners();
    startPolling();
});

// Setup event listeners
function setupEventListeners() {
    document.querySelectorAll('.time-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            document.querySelectorAll('.time-btn').forEach(b => b.classList.remove('active'));
            e.target.classList.add('active');
            currentTimeRange = parseInt(e.target.dataset.seconds);
            fetchHistory();
        });
    });
}

// Initialize Chart.js
function initChart() {
    const ctx = document.getElementById('usageChart').getContext('2d');
    chart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [
                {
                    label: 'CPU使用率',
                    data: [],
                    borderColor: 'rgb(255, 99, 132)',
                    backgroundColor: 'rgba(255, 99, 132, 0.1)',
                    tension: 0.4,
                    fill: true
                },
                {
                    label: 'RAM使用率',
                    data: [],
                    borderColor: 'rgb(54, 162, 235)',
                    backgroundColor: 'rgba(54, 162, 235, 0.1)',
                    tension: 0.4,
                    fill: true
                },
                {
                    label: 'GPU使用率',
                    data: [],
                    borderColor: 'rgb(75, 192, 192)',
                    backgroundColor: 'rgba(75, 192, 192, 0.1)',
                    tension: 0.4,
                    fill: true,
                    hidden: true
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    labels: {
                        color: '#fff',
                        font: { size: 14 }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    max: 100,
                    ticks: {
                        color: '#fff',
                        callback: (value) => value + '%'
                    },
                    grid: {
                        color: 'rgba(255, 255, 255, 0.1)'
                    }
                },
                x: {
                    ticks: {
                        color: '#fff',
                        maxTicksLimit: 10
                    },
                    grid: {
                        color: 'rgba(255, 255, 255, 0.1)'
                    }
                }
            }
        }
    });
}

// Start polling for updates
function startPolling() {
    fetchLatest();
    fetchHistory();
    setInterval(() => {
        fetchLatest();
        fetchHistory();
    }, UPDATE_INTERVAL);
}

// Fetch latest metrics
async function fetchLatest() {
    try {
        const response = await fetch(`${API_BASE}/metrics/latest`);
        if (!response.ok) throw new Error('Failed to fetch latest metrics');

        const data = await response.json();
        updateUI(data);
        updateStatus(true);
    } catch (error) {
        console.error('Error fetching latest metrics:', error);
        updateStatus(false);
    }
}

// Fetch historical metrics
async function fetchHistory() {
    try {
        const response = await fetch(`${API_BASE}/metrics/history?seconds=${currentTimeRange}`);
        if (!response.ok) throw new Error('Failed to fetch history');

        const data = await response.json();
        updateChart(data.samples);
    } catch (error) {
        console.error('Error fetching history:', error);
    }
}

// Update UI with latest metrics
function updateUI(data) {
    // CPU
    document.getElementById('cpuUsageBadge').textContent = `${(data.cpu.usage * 100).toFixed(1)}%`;
    document.getElementById('cpuModel').textContent = data.cpu.model || '-';
    document.getElementById('cpuCores').textContent = data.cpu.cores || '-';
    document.getElementById('cpuThreads').textContent = data.cpu.threads || '-';

    if (data.cpu.frequency_mhz) {
        document.getElementById('cpuFreqRow').style.display = 'flex';
        document.getElementById('cpuFreq').textContent = `${data.cpu.frequency_mhz.toFixed(0)} MHz`;
    }

    if (data.cpu.load_avg && data.cpu.load_avg.length > 0) {
        document.getElementById('cpuLoadRow').style.display = 'flex';
        const [load1, load5, load15] = data.cpu.load_avg;
        document.getElementById('cpuLoad').textContent =
            `1分: ${load1.toFixed(2)}, 5分: ${load5.toFixed(2)}, 15分: ${load15.toFixed(2)}`;
    }

    // RAM
    document.getElementById('ramUsageBadge').textContent = `${(data.ram.usage * 100).toFixed(1)}%`;
    document.getElementById('ramTotal').textContent = `${data.ram.total_mb.toLocaleString()} MB`;
    document.getElementById('ramUsed').textContent = `${data.ram.used_mb.toLocaleString()} MB`;
    document.getElementById('ramFree').textContent = `${data.ram.free_mb.toLocaleString()} MB`;

    // GPU
    if (data.gpu) {
        document.getElementById('gpuCard').style.display = 'block';
        document.getElementById('gpuUsageBadge').textContent = `${(data.gpu.util * 100).toFixed(1)}%`;
        document.getElementById('gpuVendor').textContent = data.gpu.vendor.toUpperCase();
        document.getElementById('gpuModel').textContent = data.gpu.model || '-';

        if (data.gpu.temperature_c) {
            document.getElementById('gpuTempRow').style.display = 'flex';
            document.getElementById('gpuTemp').textContent = `${data.gpu.temperature_c.toFixed(1)}°C`;
        }

        if (data.gpu.vram_total_mb) {
            document.getElementById('gpuVramRow').style.display = 'flex';
            document.getElementById('gpuVram').textContent =
                `${data.gpu.vram_used_mb.toLocaleString()} / ${data.gpu.vram_total_mb.toLocaleString()} MB`;
        }

        // Show GPU dataset in chart
        chart.data.datasets[2].hidden = false;
    }

    // Update timestamp
    const timestamp = new Date(data.timestamp);
    document.getElementById('lastUpdate').textContent = timestamp.toLocaleString('ja-JP');
}

// Update chart with historical data
function updateChart(samples) {
    if (!samples || samples.length === 0) return;

    const labels = samples.map(s => {
        const date = new Date(s.timestamp);
        return date.toLocaleTimeString('ja-JP', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
    });

    const cpuData = samples.map(s => (s.cpu.usage * 100).toFixed(1));
    const ramData = samples.map(s => (s.ram.usage * 100).toFixed(1));
    const gpuData = samples.map(s => s.gpu ? (s.gpu.util * 100).toFixed(1) : null);

    chart.data.labels = labels;
    chart.data.datasets[0].data = cpuData;
    chart.data.datasets[1].data = ramData;
    chart.data.datasets[2].data = gpuData;
    chart.update('none'); // Update without animation for smoother updates
}

// Update connection status
function updateStatus(isConnected) {
    const statusDot = document.getElementById('statusDot');
    const statusText = document.getElementById('statusText');

    if (isConnected) {
        statusDot.classList.remove('error');
        statusText.textContent = '接続中';
    } else {
        statusDot.classList.add('error');
        statusText.textContent = '接続エラー';
    }
}
