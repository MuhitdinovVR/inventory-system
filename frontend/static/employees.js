document.addEventListener('DOMContentLoaded', async function() {
    await loadEmployees();
    initEmployeeModal();

    document.getElementById('add-employee-btn').addEventListener('click', function() {
        openEmployeeModal();
    });

    document.getElementById('employee-form').addEventListener('submit', async function(e) {
        e.preventDefault();
        await saveEmployee();
    });
});

async function loadEmployees() {
    try {
        const employees = await getAllEmployees();
        const tableBody = document.getElementById('employees-table');
        tableBody.innerHTML = '';

        employees.forEach(employee => {
            const row = document.createElement('tr');

            row.innerHTML = `
                <td>${employee.id}</td>
                <td>${employee.full_name}</td>
                <td>${employee.position}</td>
                <td>${employee.email}</td>
                <td>${employee.role}</td>
                <td>${employee.department || 'Не назначен'}</td>
                <td>
                    <button class="btn edit-employee" data-id="${employee.id}">Редактировать</button>
                    <button class="btn btn-danger delete-employee" data-id="${employee.id}">Удалить</button>
                </td>
            `;

            tableBody.appendChild(row);
        });

        document.querySelectorAll('.edit-employee').forEach(btn => {
            btn.addEventListener('click', function() {
                const employeeId = this.getAttribute('data-id');
                openEmployeeModal(employeeId);
            });
        });

        document.querySelectorAll('.delete-employee').forEach(btn => {
            btn.addEventListener('click', async function() {
                const employeeId = this.getAttribute('data-id');
                if (confirm('Вы уверены, что хотите удалить этого сотрудника?')) {
                    await deleteEmployee(employeeId);
                    await loadEmployees();
                }
            });
        });

    } catch (error) {
        console.error('Ошибка загрузки сотрудников:', error);
        alert('Не удалось загрузить сотрудников');
    }
}

function initEmployeeModal() {
    const modal = document.getElementById('employee-modal');
    const closeBtn = document.querySelector('#employee-modal .close');

    closeBtn.addEventListener('click', function() {
        modal.style.display = 'none';
    });

    window.addEventListener('click', function(event) {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });
}

async function openEmployeeModal(employeeId = null) {
    const modal = document.getElementById('employee-modal');
    const form = document.getElementById('employee-form');
    const title = document.getElementById('employee-modal-title');

    form.reset();
    document.getElementById('employee-password').value = '';

    // Загрузка отделов
    await loadDepartmentsForSelect();

    if (employeeId) {
        title.textContent = 'Редактировать сотрудника';
        document.getElementById('employee-id').value = employeeId;

        const employee = await getEmployeeById(employeeId);
        document.getElementById('employee-fullname').value = employee.full_name;
        document.getElementById('employee-position').value = employee.position;
        document.getElementById('employee-email').value = employee.email;
        document.getElementById('employee-role').value = employee.role;
        document.getElementById('employee-department').value = employee.department_id || '';
    } else {
        title.textContent = 'Добавить сотрудника';
        document.getElementById('employee-id').value = '';
    }

    modal.style.display = 'block';
}

async function loadDepartmentsForSelect() {
    try {
        const departments = await getAllDepartments();
        const select = document.getElementById('employee-department');

        // Сохраняем текущее значение
        const currentValue = select.value;

        // Очищаем и добавляем только опцию "Не назначен"
        select.innerHTML = '<option value="">Не назначен</option>';

        departments.forEach(department => {
            const option = document.createElement('option');
            option.value = department.id;
            option.textContent = department.name;
            select.appendChild(option);
        });

        // Восстанавливаем предыдущее значение
        if (currentValue) {
            select.value = currentValue;
        }
    } catch (error) {
        console.error('Ошибка загрузки отделов:', error);
    }
}

async function saveEmployee() {
    const modal = document.getElementById('employee-modal');
    const employeeId = document.getElementById('employee-id').value;
    const password = document.getElementById('employee-password').value;

    const employee = {
        full_name: document.getElementById('employee-fullname').value,
        position: document.getElementById('employee-position').value,
        email: document.getElementById('employee-email').value,
        role: document.getElementById('employee-role').value,
        department_id: document.getElementById('employee-department').value || null
    };

    if (password) {
        employee.password_hash = password;
    }

    try {
        if (employeeId) {
            await updateEmployee(employeeId, employee);
        } else {
            if (!password) {
                throw new Error('Для нового сотрудника необходимо указать пароль');
            }
            await createEmployee(employee);
        }

        modal.style.display = 'none';
        await loadEmployees();
    } catch (error) {
        console.error('Ошибка сохранения сотрудника:', error);
        alert(error.message || 'Не удалось сохранить сотрудника');
    }
}