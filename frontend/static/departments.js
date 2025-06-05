document.addEventListener('DOMContentLoaded', async function() {
    await loadDepartments();
    initDepartmentModal();

    document.getElementById('add-department-btn').addEventListener('click', function() {
        openDepartmentModal();
    });

    document.getElementById('department-form').addEventListener('submit', async function(e) {
        e.preventDefault();
        await saveDepartment();
    });
});

async function loadDepartments() {
    try {
        const departments = await getAllDepartments();
        const tableBody = document.getElementById('departments-table');
        tableBody.innerHTML = '';

        departments.forEach(department => {
            const row = document.createElement('tr');

            row.innerHTML = `
                <td>${department.id}</td>
                <td>${department.name}</td>
                <td>${department.location}</td>
                <td>${department.head_name || 'Не назначен'}</td>
                <td>
                    <button class="btn edit-department" data-id="${department.id}">Редактировать</button>
                    <button class="btn btn-danger delete-department" data-id="${department.id}">Удалить</button>
                    <button class="btn view-employees" data-id="${department.id}">Сотрудники</button>
                </td>
            `;

            tableBody.appendChild(row);
        });

        document.querySelectorAll('.edit-department').forEach(btn => {
            btn.addEventListener('click', function() {
                const departmentId = this.getAttribute('data-id');
                openDepartmentModal(departmentId);
            });
        });

        document.querySelectorAll('.delete-department').forEach(btn => {
            btn.addEventListener('click', async function() {
                const departmentId = this.getAttribute('data-id');
                if (confirm('Вы уверены, что хотите удалить этот отдел?')) {
                    try {
                        await deleteDepartment(departmentId);
                        await loadDepartments();
                    } catch (error) {
                        alert(error.message || 'Не удалось удалить отдел');
                    }
                }
            });
        });

        document.querySelectorAll('.view-employees').forEach(btn => {
            btn.addEventListener('click', function() {
                const departmentId = this.getAttribute('data-id');
                window.location.href = `../department-employees.html?id=${departmentId}`;
            });
        });

    } catch (error) {
        console.error('Ошибка загрузки отделов:', error);
        alert('Не удалось загрузить отделы');
    }
}

function initDepartmentModal() {
    const modal = document.getElementById('department-modal');
    const closeBtn = document.querySelector('#department-modal .close');

    closeBtn.addEventListener('click', function() {
        modal.style.display = 'none';
    });

    window.addEventListener('click', function(event) {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });
}

async function openDepartmentModal(departmentId = null) {
    const modal = document.getElementById('department-modal');
    const form = document.getElementById('department-form');
    const title = document.getElementById('department-modal-title');

    form.reset();

    // Загрузка потенциальных руководителей
    await loadPotentialHeads();

    if (departmentId) {
        title.textContent = 'Редактировать отдел';
        document.getElementById('department-id').value = departmentId;

        const department = await getDepartmentById(departmentId);
        document.getElementById('department-name').value = department.name;
        document.getElementById('department-location').value = department.location;
        document.getElementById('department-head').value = department.head_id || '';
    } else {
        title.textContent = 'Добавить отдел';
        document.getElementById('department-id').value = '';
    }

    modal.style.display = 'block';
}

async function loadPotentialHeads() {
    try {
        const employees = await getAllEmployees();
        const select = document.getElementById('department-head');

        // Сохраняем текущее значение
        const currentValue = select.value;

        // Очищаем и добавляем только опцию "Не назначен"
        select.innerHTML = '<option value="">Не назначен</option>';

        employees.forEach(employee => {
            const option = document.createElement('option');
            option.value = employee.id;
            option.textContent = `${employee.full_name} (${employee.position})`;
            select.appendChild(option);
        });

        // Восстанавливаем предыдущее значение
        if (currentValue) {
            select.value = currentValue;
        }
    } catch (error) {
        console.error('Ошибка загрузки сотрудников:', error);
    }
}

async function saveDepartment() {
    const modal = document.getElementById('department-modal');
    const departmentId = document.getElementById('department-id').value;

    const department = {
        name: document.getElementById('department-name').value,
        location: document.getElementById('department-location').value,
        head_id: document.getElementById('department-head').value || null
    };

    try {
        if (departmentId) {
            await updateDepartment(departmentId, department);
        } else {
            await createDepartment(department);
        }

        modal.style.display = 'none';
        await loadDepartments();
    } catch (error) {
        console.error('Ошибка сохранения отдела:', error);
        alert('Не удалось сохранить отдел');
    }
}