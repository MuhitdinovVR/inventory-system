// Auth module for inventory system
class Auth {
    constructor() {
        this.token = null;
        this.user = null;
    }

    async login(username, password) {
        try {
            const response = await fetch('/api/auth', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
            });

            const data = await response.json();

            if (data.status === 'auth_ok') {
                this.token = data.token;
                this.user = data.user;
                return true;
            }
            return false;
        } catch (error) {
            console.error('Login error:', error);
            return false;
        }
    }

    logout() {
        this.token = null;
        this.user = null;
    }

    isAuthenticated() {
        return this.token !== null;
    }
}

// Экспорт модуля
const auth = new Auth();
export default auth;