const API_BASE_URL = window.location.origin;

async function makeRequest(url, method = 'GET', body = null) {
    const headers = {
        'Content-Type': 'application/json',
    };

    // Получаем токен из localStorage
    const token = localStorage.getItem('token');
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const options = {
        method,
        headers,
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    const response = await fetch(`${API_BASE_URL}/api${url}`, options); // Добавлен /api

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Ошибка запроса');
    }

    return response.json();
}

// Auth
async function login(email, password) {
    return makeRequest('/api/auth/login', 'POST', { email, password });
}

async function register(employee) {
    return makeRequest('/api/auth/register', 'POST', employee);
}


// Активы
async function getAllAssets() {
    return makeRequest('/assets');
}

async function getAssetById(id) {
    return makeRequest(`/assets/${id}`);
}

async function createAsset(asset) {
    return makeRequest('/assets', 'POST', asset);
}

async function updateAsset(id, asset) {
    return makeRequest(`/assets/${id}`, 'PUT', asset);
}

async function deleteAsset(id) {
    return makeRequest(`/assets/${id}`, 'DELETE');
}

async function getAssetTransfers(id) {
    return makeRequest(`/assets/${id}/transfers`);
}

// Сотрудники
async function getAllEmployees() {
    return makeRequest('/employees');
}

async function getEmployeeById(id) {
    return makeRequest(`/employees/${id}`);
}

async function createEmployee(employee) {
    return makeRequest('/employees', 'POST', employee);
}

async function updateEmployee(id, employee) {
    return makeRequest(`/employees/${id}`, 'PUT', employee);
}

async function deleteEmployee(id) {
    return makeRequest(`/employees/${id}`, 'DELETE');
}

// Отделы
async function getAllDepartments() {
    return makeRequest('/departments');
}

async function getDepartmentById(id) {
    return makeRequest(`/departments/${id}`);
}

async function createDepartment(department) {
    return makeRequest('/departments', 'POST', department);
}

async function updateDepartment(id, department) {
    return makeRequest(`/departments/${id}`, 'PUT', department);
}

async function deleteDepartment(id) {
    return makeRequest(`/departments/${id}`, 'DELETE');
}

async function getDepartmentEmployees(id) {
    return makeRequest(`/departments/${id}/employees`);
}

// Перемещения
async function getAllTransfers() {
    return makeRequest('/transfers');
}

async function getTransferById(id) {
    return makeRequest(`/transfers/${id}`);
}

async function createTransfer(transfer) {
    return makeRequest('/transfers', 'POST', transfer);
}

async function getRecentTransfers() {
    const transfers = await getAllTransfers();
    return transfers.sort((a, b) => new Date(b.transfer_date) - new Date(a.transfer_date));
}

// Отчеты
async function getAssetsByStatusReport() {
    return makeRequest('/reports/assets-by-status');
}

async function getTransfersReport(from, to) {
    return makeRequest(`/reports/transfers?from=${from}&to=${to}`);
}

async function getDepartmentCostsReport() {
    return makeRequest('/reports/department-costs');
}

// Вспомогательные данные
async function getAllStatuses() {
    return makeRequest('/assets/statuses');
}

async function getAllLocations() {
    return makeRequest('/locations');
}