document.addEventListener('DOMContentLoaded', async function() {
    await loadTransfers();
    initTransferModal();
    await loadFilters();

    document.getElementById('add-transfer-btn').addEventListener('click', function() {
        openTransferModal();
    });

    document.getElementById('transfer-form').addEventListener('submit', async function(e) {
        e.preventDefault();
        await saveTransfer();
    });

    document.getElementById('transfer-filters').addEventListener('submit', async function(e) {
        e.preventDefault();
        await applyFilters();
    });

    document.getElementById('transfer-filters').addEventListener('reset', async function() {
        await loadTransfers();
    });
});

async function loadTransfers(filters = {}) {
    try {
        let transfers;

        if (Object.keys(filters).length === 0) {
            transfers = await getAllTransfers();
        } else {
            // Здесь можно добавить логику фильтрации на клиенте
            // или сделать запрос к API с параметрами фильтрации
            transfers = await getAllTransfers(); // Временное решение
        }

        const tableBody = document.getElementById('transfers-table');
        tableBody.innerHTML = '';

        transfers.forEach(transfer => {
            const row = document.createElement('tr');

            row.innerHTML = `
                <td>${transfer.id}</td>
                <td>${transfer.asset_name}</td>
                <td>${transfer.employee_name}</td>
                <td>${transfer.from_location}</td>
                <td>${transfer.to_location}</td>
                <td>${new Date(transfer.transfer_date).toLocaleString()}</td>
                <td>
                    <button class="btn view-transfer" data-id="${transfer.id}">Просмотр</button>
                    <button class="btn btn-danger delete-transfer" data-id="${transfer.id}">Удалить</button>
                </td>
            `;

            tableBody.appendChild(row);
        });

        document.querySelectorAll('.view-transfer').forEach(btn => {
            btn.addEventListener('click', function() {
                const transferId = this.getAttribute('data-id');
                viewTransfer(transferId);
            });
        });

        document.querySelectorAll('.delete-transfer').forEach(btn => {
            btn.addEventListener('click', async function() {
                const transferId = this.getAttribute('data-id');
                if (confirm('Вы уверены, что хотите удалить это перемещение?')) {
                    await deleteTransfer(transferId);
                    await loadTransfers();
                }
            });
        });

    } catch (error) {
        console.error('Ошибка загрузки перемещений:', error);
        alert('Не удалось загрузить перемещения');
    }
}

async function loadFilters() {
    try {
        const [assets, employees, locations] = await Promise.all([
            getAllAssets(),
            getAllEmployees(),
            getAllLocations()
        ]);

        const assetSelect = document.getElementById('filter-asset');
        assets.forEach(asset => {
            const option = document.createElement('option');
            option.value = asset.id;
            option.textContent = asset.name;
            assetSelect.appendChild(option);
        });

        const employeeSelect = document.getElementById('filter-employee');
        employees.forEach(employee => {
            const option = document.createElement('option');
            option.value = employee.id;
            option.textContent = employee.full_name;
            employeeSelect.appendChild(option);
        });

        const locationSelect = document.getElementById('filter-location');
        locations.forEach(location => {
            const option = document.createElement('option');
            option.value = location.id;
            option.textContent = location.address;
            locationSelect.appendChild(option);
        });
    } catch (error) {
        console.error('Ошибка загрузки фильтров:', error);
    }
}

async function applyFilters() {
    const filters = {
        assetId: document.getElementById('filter-asset').value || null,
        employeeId: document.getElementById('filter-employee').value || null,
        locationId: document.getElementById('filter-location').value || null,
        date: document.getElementById('filter-date').value || null
    };

    // Здесь можно добавить логику фильтрации или сделать запрос к API
    await loadTransfers(filters);
}

function initTransferModal() {
    const modal = document.getElementById('transfer-modal');
    const closeBtn = document.querySelector('#transfer-modal .close');

    closeBtn.addEventListener('click', function() {
        modal.style.display = 'none';
    });

    window.addEventListener('click', function(event) {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });
}

async function openTransferModal() {
    const modal = document.getElementById('transfer-modal');
    const form = document.getElementById('transfer-form');
    const title = document.getElementById('transfer-modal-title');

    form.reset();
    title.textContent = 'Добавить перемещение';
    document.getElementById('transfer-id').value = '';

    // Загрузка данных для формы
    await loadTransferFormData();

    // Установка текущей даты и времени
    const now = new Date();
    const formattedDate = now.toISOString().slice(0, 16);
    document.getElementById('transfer-date').value = formattedDate;

    modal.style.display = 'block';
}

async function loadTransferFormData() {
    try {
        const [assets, employees, locations] = await Promise.all([
            getAllAssets(),
            getAllEmployees(),
            getAllLocations()
        ]);

        const assetSelect = document.getElementById('transfer-asset');
        assetSelect.innerHTML = '';
        assets.forEach(asset => {
            const option = document.createElement('option');
            option.value = asset.id;
            option.textContent = asset.name;
            assetSelect.appendChild(option);
        });

        const employeeSelect = document.getElementById('transfer-employee');
        employeeSelect.innerHTML = '';
        employees.forEach(employee => {
            const option = document.createElement('option');
            option.value = employee.id;
            option.textContent = employee.full_name;
            employeeSelect.appendChild(option);
        });

        const fromSelect = document.getElementById('transfer-from');
        const toSelect = document.getElementById('transfer-to');
        fromSelect.innerHTML = '';
        toSelect.innerHTML = '';
        locations.forEach(location => {
            const optionFrom = document.createElement('option');
            optionFrom.value = location.id;
            optionFrom.textContent = location.address;
            fromSelect.appendChild(optionFrom);

            const optionTo = document.createElement('option');
            optionTo.value = location.id;
            optionTo.textContent = location.address;
            toSelect.appendChild(optionTo);
        });
    } catch (error) {
        console.error('Ошибка загрузки данных для формы:', error);
        alert('Не удалось загрузить данные для формы');
    }
}

async function saveTransfer() {
    const modal = document.getElementById('transfer-modal');
    const transferId = document.getElementById('transfer-id').value;

    const transfer = {
        asset_id: parseInt(document.getElementById('transfer-asset').value),
        employee_id: parseInt(document.getElementById('transfer-employee').value),
        from_location_id: parseInt(document.getElementById('transfer-from').value),
        to_location_id: parseInt(document.getElementById('transfer-to').value),
        transfer_date: document.getElementById('transfer-date').value,
        notes: document.getElementById('transfer-notes').value
    };

    try {
        if (transfer.from_location_id === transfer.to_location_id) {
            throw new Error('Местоположения "Откуда" и "Куда" не должны совпадать');
        }

        await createTransfer(transfer);

        modal.style.display = 'none';
        await loadTransfers();
    } catch (error) {
        console.error('Ошибка сохранения перемещения:', error);
        alert(error.message || 'Не удалось сохранить перемещение');
    }
}

async function viewTransfer(transferId) {
    try {
        const transfer = await getTransferById(transferId);
        alert(`Подробности перемещения:\n
Актив: ${transfer.asset_name}\n
Сотрудник: ${transfer.employee_name}\n
Откуда: ${transfer.from_location}\n
Куда: ${transfer.to_location}\n
Дата: ${new Date(transfer.transfer_date).toLocaleString()}\n
Примечания: ${transfer.notes || 'нет'}`);
    } catch (error) {
        console.error('Ошибка просмотра перемещения:', error);
        alert('Не удалось загрузить данные перемещения');
    }
}