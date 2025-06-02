// Проверка аутентификации при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    checkAuth();
    initNavigation();
});

function initNavigation() {
    // Обработчик выхода
    const logoutBtn = document.getElementById('logout');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', function(e) {
            e.preventDefault();
            logout();
        });
    }

    // Подсветка активной ссылки в навигации
    const currentPath = window.location.pathname.split('/').pop() || 'index.html';
    const navLinks = document.querySelectorAll('.nav-links a');

    navLinks.forEach(link => {
        const linkPath = link.getAttribute('href').split('/').pop();
        if (currentPath === linkPath ||
            (currentPath === '' && linkPath === 'index.html')) {
            link.classList.add('active');
        }

        // Добавляем обработчики для плавных переходов
        link.addEventListener('click', function(e) {
            if (this.getAttribute('href').startsWith('http')) return;

            e.preventDefault();
            const href = this.getAttribute('href');
            document.body.style.opacity = '0.5';
            setTimeout(() => {
                window.location.href = href;
            }, 200);
        });
    });
}

// Проверка аутентификации
function checkAuth() {
    const token = localStorage.getItem('token');
    const isLoginPage = window.location.pathname.includes('login.html') ||
        window.location.pathname.includes('register.html');

    if (!token && !isLoginPage) {
        window.location.href = 'login.html';
    } else if (token && isLoginPage) {
        window.location.href = 'index.html';
    }
}

// Выход из системы
function logout() {
    localStorage.removeItem('token');
    window.location.href = 'login.html';
}

// Загрузка статистики для главной страницы
async function loadDashboardStats() {
    try {
        // Загрузка общего количества активов
        const assets = await getAllAssets();
        document.getElementById('total-assets').textContent = assets.length;

        // Загрузка общего количества сотрудников
        const employees = await getAllEmployees();
        document.getElementById('total-employees').textContent = employees.length;

        // Загрузка последних перемещений
        const transfers = await getRecentTransfers();
        const transfersList = document.getElementById('recent-transfers');
        transfersList.innerHTML = '';

        transfers.slice(0, 5).forEach(transfer => {
            const li = document.createElement('li');
            li.innerHTML = `
                <strong>${transfer.asset_name}</strong> → 
                ${transfer.to_location}<br>
                <small>${new Date(transfer.transfer_date).toLocaleDateString()}</small>
            `;
            transfersList.appendChild(li);
        });
    } catch (error) {
        console.error('Ошибка загрузки статистики:', error);
        showError('Не удалось загрузить статистику');
    }
}

function showError(message) {
    const errorElement = document.createElement('div');
    errorElement.className = 'alert alert-danger';
    errorElement.textContent = message;
    document.body.prepend(errorElement);

    setTimeout(() => {
        errorElement.remove();
    }, 5000);
}

// Инициализация страницы
if (document.getElementById('total-assets')) {
    loadDashboardStats();
}