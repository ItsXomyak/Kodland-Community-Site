// Функция для проверки авторизации и обновления навигации
function updateNavigation() {
	const token = localStorage.getItem('token')
	const authLink = document.getElementById('authLink')
	const profileNavItem = document.getElementById('profileNavItem')
	const userName = document.getElementById('userName')

	if (token) {
		// Если пользователь авторизован
		if (authLink) authLink.parentElement.classList.add('d-none')
		if (profileNavItem) profileNavItem.classList.remove('d-none')

		// Получаем информацию о пользователе из токена
		try {
			const payload = JSON.parse(atob(token.split('.')[1]))
			if (userName && payload.email) {
				userName.textContent = payload.email.split('@')[0] // Показываем имя пользователя до @
			}
		} catch (error) {
			console.error('Ошибка при разборе токена:', error)
		}
	} else {
		// Если пользователь не авторизован
		if (authLink) authLink.parentElement.classList.remove('d-none')
		if (profileNavItem) profileNavItem.classList.add('d-none')
	}
}

// Обновляем навигацию при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
	updateNavigation()
})
