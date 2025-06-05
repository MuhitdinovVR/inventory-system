// Auth module for inventory system
class Auth {
    constructor() {
        this.token = null;
        this.user = null;
    }

    async login(email, password) {
        try {
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Login failed');
            }

            if (data.token) {
                this.token = data.token;
                this.user = data.employee;
                localStorage.setItem('token', data.token);
                return true;
            }
            return false;
        } catch (error) {
            console.error('Login error:', error);
            throw error;
        }
    }

    logout() {
        this.token = null;
        this.user = null;
        localStorage.removeItem('token');
    }

    isAuthenticated() {
        return this.token !== null;
    }
}
window.auth = new Auth();
// Экспорт модуля
const authInstance = new Auth();
export default authInstance;