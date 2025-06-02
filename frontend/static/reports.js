document.addEventListener('DOMContentLoaded', async function() {
    // Инициализация вкладок
    initTabs();

    // Загрузка отчетов
    await loadStatusReport();
    await loadDepartmentReport();
    await loadTransfersReport();
    await loadInventoryReport();

    // Обработчик фильтра перемещений
    document.getElementById('transfers-filter').addEventListener('submit', async function(e) {
        e.preventDefault();
        await loadTransfersReport();
    });
});

function initTabs() {
    const tabBtns = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-content');

    tabBtns.forEach(btn => {
        btn.addEventListener('click', function() {
            const tabId = this.getAttribute('data-tab');

            // Убираем активный класс у всех кнопок и контента
            tabBtns.forEach(b => b.classList.remove('active'));
            tabContents.forEach(c => c.classList.remove('active'));

            // Добавляем активный класс текущей кнопке и контенту
            this.classList.add('active');
            document.getElementById(tabId).classList.add('active');
        });
    });
}

async function loadStatusReport() {
    try {
        const report = await getAssetsByStatusReport();
        const tableBody = document.getElementById('status-table');
        tableBody.innerHTML = '';

        report.forEach(item => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${item.status}</td>
                <td>${item.count}</td>
                <td>${item.total_cost.toFixed(2)} ₽</td>
            `;
            tableBody.appendChild(row);
        });

        // Создание графика
        const ctx = document.getElementById('status-chart').getContext('2d');
        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: report.map(item => item.status),
                datasets: [{
                    label: 'Количество активов',
                    data: report.map(item => item.count),
                    backgroundColor: 'rgba(54, 162, 235, 0.5)',
                    borderColor: 'rgba(54, 162, 235, 1)',
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });
    } catch (error) {
        console.error('Ошибка загрузки отчета по статусам:', error);
    }
}

async function loadDepartmentReport() {
    try {
        const report = await getDepartmentCostsReport();
        const tableBody = document.getElementById('department-table');
        tableBody.innerHTML = '';

        report.forEach(item => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${item.department}</td>
                <td>${item.count}</td>
                <td>${item.total_cost.toFixed(2)} ₽</td>
                <td>${item.average_cost.toFixed(2)} ₽</td>
            `;
            tableBody.appendChild(row);
        });

        // Создание графика
        const ctx = document.getElementById('department-chart').getContext('2d');
        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: report.map(item => item.department),
                datasets: [{
                    label: 'Общая стоимость',
                    data: report.map(item => item.total_cost),
                    backgroundColor: 'rgba(75, 192, 192, 0.5)',
                    borderColor: 'rgba(75, 192, 192, 1)',
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });
    } catch (error) {
        console.error('Ошибка загрузки отчета по отделам:', error);
    }
}

async function loadTransfersReport() {
    try {
        const fromDate = document.getElementById('transfers-from').value;
        const toDate = document.getElementById('transfers-to').value;

        let transfers;

        if (fromDate && toDate) {
            transfers = await getTransfersReport(fromDate, toDate);
        } else {
            // Если даты не указаны, загружаем последние 30 дней
            const to = new Date();
            const from = new Date();
            from.setDate(from.getDate() - 30);

            transfers = await getTransfersReport(from.toISOString().split('T')[0], to.toISOString().split('T')[0]);

            // Устанавливаем даты в поля ввода
            document.getElementById('transfers-from').value = from.toISOString().split('T')[0];
            document.getElementById('transfers-to').value = to.toISOString().split('T')[0];
        }

        const tableBody = document.getElementById('transfers-report-table');
        tableBody.innerHTML = '';

        transfers.forEach(transfer => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${new Date(transfer.transfer_date).toLocaleDateString()}</td>
                <td>${transfer.asset_name}</td>
                <td>${transfer.from_location}</td>
                <td>${transfer.to_location}</td>
                <td>${transfer.employee_name}</td>
            `;
            tableBody.appendChild(row);
        });
    } catch (error) {
        console.error('Ошибка загрузки отчета по перемещениям:', error);
    }
}

async function loadInventoryReport() {
    try {
        const report = await getInventoryReport();

        // Основные показатели
        document.getElementById('total-assets-count').textContent = report.total_assets;
        document.getElementById('total-assets-value').textContent = `${report.total_value.toFixed(2)} ₽`;
        document.getElementById('report-date').textContent = new Date(report.generated_at).toLocaleString();

        // Графики распределения
        createDistributionChart('inventory-status-chart', report.by_status, 'Статусы');
        createDistributionChart('inventory-location-chart', report.by_location, 'Местоположения');
        createDistributionChart('inventory-department-chart', report.by_department, 'Отделы');

        // Последние перемещения
        const transfersTable = document.getElementById('recent-transfers-table');
        transfersTable.innerHTML = '';

        report.recent_transfers.forEach(transfer => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${new Date(transfer.transfer_date).toLocaleDateString()}</td>
                <td>${transfer.asset_name}</td>
                <td>${transfer.from_location}</td>
                <td>${transfer.to_location}</td>
            `;
            transfersTable.appendChild(row);
        });
    } catch (error) {
        console.error('Ошибка загрузки инвентаризационного отчета:', error);
    }
}

function createDistributionChart(canvasId, data, title) {
    const ctx = document.getElementById(canvasId).getContext('2d');

    // Преобразуем объект в массив для Chart.js
    const labels = Object.keys(data);
    const values = Object.values(data);

    // Генерация цветов
    const backgroundColors = labels.map((_, i) => {
        const hue = (i * 360 / labels.length) % 360;
        return `hsla(${hue}, 70%, 50%, 0.7)`;
    });

    new Chart(ctx, {
        type: 'pie',
        data: {
            labels: labels,
            datasets: [{
                data: values,
                backgroundColor: backgroundColors,
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            plugins: {
                title: {
                    display: true,
                    text: title
                }
            }
        }
    });
}