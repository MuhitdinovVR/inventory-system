// Проверка аутентификации при загрузке страницы

document.addEventListener('DOMContentLoaded', function() {
    
    if (window.performance && window.performance.navigation.type === 1) {
        // Принудительное обновление данных при перезагрузке страницы
        localStorage.setItem('cacheBuster', new Date().getTime());
    }

    checkAuth();
    initNavigation();
    // Удаляем дублирующиеся обработчики
    document.querySelectorAll('script[src*="auth.js"]').forEach(script => {
        if (script.getAttribute('data-processed')) {
            script.remove();
        } else {
            script.setAttribute('data-processed', 'true');
        }
    });
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
    const currentPath = window.location.pathname.replace('.html', '');
    const navLinks = document.querySelectorAll('.nav-links a');

    navLinks.forEach(link => {
        const linkPath = link.getAttribute('href').replace('.html', '');
        if (currentPath === linkPath ||
            (currentPath === '/' && linkPath === '/index')) {
            link.classList.add('active');
        }

        // Добавляем обработчики для плавных переходов
        link.addEventListener('click', function(e) {
            if (this.getAttribute('href').startsWith('http')) return;

            e.preventDefault();
            const href = this.getAttribute('href').replace('.html', '');
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
    const path = window.location.pathname.replace('.html', '');
    const isLoginPage = path === '/login' || path === '/register';

    if (!token && !isLoginPage) {
        window.location.href = '/login';
    } else if (token && isLoginPage) {
        window.location.href = '/';
    }
}

// Выход из системы
function logout() {
    localStorage.removeItem('token');
    window.location.href = '/login';
}