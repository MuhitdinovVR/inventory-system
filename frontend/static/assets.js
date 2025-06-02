document.addEventListener('DOMContentLoaded', async function() {
    // Загрузка данных при открытии страницы
    await loadAssets();

    // Инициализация модального окна
    initAssetModal();

    // Обработчик кнопки добавления актива
    document.getElementById('add-asset-btn').addEventListener('click', function() {
        openAssetModal();
    });

    // Обработчик формы актива
    document.getElementById('asset-form').addEventListener('submit', async function(e) {
        e.preventDefault();
        await saveAsset();
    });
});

// Загрузка списка активов
async function loadAssets() {
    try {
        const assets = await getAllAssets();
        const tableBody = document.getElementById('assets-table');
        tableBody.innerHTML = '';

        assets.forEach(asset => {
            const row = document.createElement('tr');

            row.innerHTML = `
                <td>${asset.id}</td>
                <td>${asset.name}</td>
                <td>${asset.category}</td>
                <td>${asset.status}</td>
                <td>${asset.location}</td>
                <td>${asset.department || 'Не назначен'}</td>
                <td>
                    <button class="btn edit-asset" data-id="${asset.id}">Редактировать</button>
                    <button class="btn btn-danger delete-asset" data-id="${asset.id}">Удалить</button>
                </td>
            `;

            tableBody.appendChild(row);
        });

        // Добавление обработчиков для кнопок редактирования и удаления
        document.querySelectorAll('.edit-asset').forEach(btn => {
            btn.addEventListener('click', function() {
                const assetId = this.getAttribute('data-id');
                openAssetModal(assetId);
            });
        });

        document.querySelectorAll('.delete-asset').forEach(btn => {
            btn.addEventListener('click', async function() {
                const assetId = this.getAttribute('data-id');
                if (confirm('Вы уверены, что хотите удалить этот актив?')) {
                    await deleteAsset(assetId);
                    await loadAssets();
                }
            });
        });

    } catch (error) {
        console.error('Ошибка загрузки активов:', error);
        alert('Не удалось загрузить активы');
    }
}

// Инициализация модального окна
function initAssetModal() {
    const modal = document.getElementById('asset-modal');
    const closeBtn = document.querySelector('.close');

    closeBtn.addEventListener('click', function() {
        modal.style.display = 'none';
    });

    window.addEventListener('click', function(event) {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });
}

// Открытие модального окна для добавления/редактирования актива
async function openAssetModal(assetId = null) {
    const modal = document.getElementById('asset-modal');
    const form = document.getElementById('asset-form');
    const title = document.getElementById('modal-title');

    // Очистка формы
    form.reset();

    // Загрузка справочников
    await loadReferenceData();

    if (assetId) {
        // Редактирование существующего актива
        title.textContent = 'Редактировать актив';
        document.getElementById('asset-id').value = assetId;

        const asset = await getAssetById(assetId);
        document.getElementById('asset-name').value = asset.name;
        document.getElementById('asset-category').value = asset.category;
        document.getElementById('asset-date').value = asset.acquisition_date;
        document.getElementById('asset-cost').value = asset.cost;
        document.getElementById('asset-status').value = asset.status_id;
        document.getElementById('asset-location').value = asset.current_location_id;
        document.getElementById('asset-department').value = asset.department_id || '';
    } else {
        // Добавление нового актива
        title.textContent = 'Добавить актив';
        document.getElementById('asset-id').value = '';
    }

    modal.style.display = 'block';
}

// Загрузка справочных данных (статусы, местоположения, отделы)
async function loadReferenceData() {
    try {
        const [statuses, locations, departments] = await Promise.all([
            getAllStatuses(),
            getAllLocations(),
            getAllDepartments()
        ]);

        const statusSelect = document.getElementById('asset-status');
        statusSelect.innerHTML = '';
        statuses.forEach(status => {
            const option = document.createElement('option');
            option.value = status.id;
            option.textContent = status.name;
            statusSelect.appendChild(option);
        });

        const locationSelect = document.getElementById('asset-location');
        locationSelect.innerHTML = '';
        locations.forEach(location => {
            const option = document.createElement('option');
            option.value = location.id;
            option.textContent = location.address;
            locationSelect.appendChild(option);
        });

        const departmentSelect = document.getElementById('asset-department');
        departmentSelect.innerHTML = '<option value="">Не назначен</option>';
        departments.forEach(department => {
            const option = document.createElement('option');
            option.value = department.id;
            option.textContent = department.name;
            departmentSelect.appendChild(option);
        });

    } catch (error) {
        console.error('Ошибка загрузки справочников:', error);
        alert('Не удалось загрузить справочные данные');
    }
}

// Сохранение актива
async function saveAsset() {
    const modal = document.getElementById('asset-modal');
    const assetId = document.getElementById('asset-id').value;

    const asset = {
        name: document.getElementById('asset-name').value,
        category: document.getElementById('asset-category').value,
        acquisition_date: document.getElementById('asset-date').value,
        cost: parseFloat(document.getElementById('asset-cost').value),
        status_id: parseInt(document.getElementById('asset-status').value),
        current_location_id: parseInt(document.getElementById('asset-location').value),
        department_id: document.getElementById('asset-department').value || null
    };

    try {
        if (assetId) {
            // Обновление существующего актива
            await updateAsset(assetId, asset);
        } else {
            // Создание нового актива
            await createAsset(asset);
        }

        modal.style.display = 'none';
        await loadAssets();
    } catch (error) {
        console.error('Ошибка сохранения актива:', error);
        alert('Не удалось сохранить актив');
    }
}